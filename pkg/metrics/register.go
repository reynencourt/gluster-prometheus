package metrics

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/reynencourt/gluster-prometheus/pkg/conf"
	"github.com/reynencourt/gluster-prometheus/pkg/glusterutils"
	"github.com/reynencourt/gluster-prometheus/pkg/glusterutils/glusterconsts"
	"github.com/sirupsen/logrus"
)

type glusterMetric struct {
	name string
	fn   func(glusterutils.GInterface) error
}

var glusterMetrics []glusterMetric

func registerMetric(name string, fn func(glusterutils.GInterface) error) {
	glusterMetrics = append(glusterMetrics, glusterMetric{name: name, fn: fn})
}

func getDefaultGlusterdDir(mgmt string) string {
	if mgmt == glusterconsts.MgmtGlusterd2 {
		return defaultGlusterd2Workdir
	}
	return defaultGlusterd1Workdir
}

func InitGluterMetrics(clusterLabel string, configPath string, metrics []string) error {
	for _, metric := range metrics {
		switch metric {
		case "gluster_brick":
			registerMetric("gluster_brick", brickUtilization)
		case "gluster_brick_status":
			registerMetric("gluster_brick_status", brickStatus)
		case "gluster_peer_counts":
			registerMetric("gluster_peer_counts", peerCounts)
		case "gluster_peer_info":
			registerMetric("gluster_peer_info", peerInfo)
		case "gluster_ps":
			registerMetric("gluster_ps", ps)
		case "gluster_volume_heal":
			registerMetric("gluster_volume_heal", healCounts)
		case "gluster_volume_profile":
			registerMetric("gluster_volume_profile", profileInfo)
		case "gluster_volume_counts":
			registerMetric("gluster_volume_counts", volumeCounts)
		case "gluster_volume_status":
			registerMetric("gluster_volume_status", volumeInfo)
		default:
			return errors.New(fmt.Sprintf("metric '%s' not found", metric))
		}
	}

	f, err := os.Stat(configPath)
	if err != nil {
		logrus.WithError(err).Error("could not stat the config file ", configPath)
	}

	if f.IsDir() {
		logrus.WithError(err).Error("config file given is a path", configPath)
	}

	exporterConf, err := conf.LoadConfig(configPath)
	if err != nil {
		return err
	}

	if exporterConf.GlusterdWorkdir == "" {
		exporterConf.GlusterdWorkdir =
			getDefaultGlusterdDir(exporterConf.GlusterMgmt)
	}

	gluster := glusterutils.MakeGluster(exporterConf)

	for _, m := range glusterMetrics {
		if collectorConf, ok := exporterConf.CollectorsConf[m.name]; ok {
			if !collectorConf.Disabled {
				go func(m glusterMetric, gi glusterutils.GInterface) {
					for {
						// exporter's config will have proper Cluster ID set
						clusterID = clusterLabel
						err := m.fn(gi)
						interval := defaultInterval
						if collectorConf.SyncInterval > 0 {
							interval = time.Duration(collectorConf.SyncInterval)
						}
						if err != nil {
							logrus.WithError(err).WithFields(logrus.Fields{
								"name": m.name,
							}).Debug("failed to export metric")
						}
						time.Sleep(time.Second * interval)
					}
				}(m, gluster)
			}
		}
	}

	if len(glusterMetrics) == 0 {
		return errors.New("no Metrics registered")
	}

	return nil
}
