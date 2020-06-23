package main

import (
	"github.com/gluster/gluster-prometheus/gluster-exporter/conf"
	"github.com/gluster/gluster-prometheus/pkg/glusterutils"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type glusterMetric struct {
	name string
	fn   func(glusterutils.GInterface) error
}

var glusterMetrics []glusterMetric

func registerMetric(name string, fn func(glusterutils.GInterface) error) {
	glusterMetrics = append(glusterMetrics, glusterMetric{name: name, fn: fn})
}

func InitGluterMetrics(configPath string) {
	registerMetric("gluster_brick", brickUtilization)
	registerMetric("gluster_brick_status", brickStatus)
	registerMetric("gluster_peer_counts", peerCounts)
	registerMetric("gluster_peer_info", peerInfo)
	registerMetric("gluster_ps", ps)
	registerMetric("gluster_volume_heal", healCounts)
	registerMetric("gluster_volume_profile", profileInfo)
	registerMetric("gluster_volume_counts", volumeCounts)
	registerMetric("gluster_volume_status", volumeInfo)

	f, err := os.Stat(configPath)
	if err != nil {
		logrus.WithError(err).Error("could not stat the config file ", configPath)
	}

	if f.IsDir() {
		logrus.WithError(err).Error("config file given is a path", configPath)
	}

	exporterConf, err := conf.LoadConfig(configPath)
	if err != nil {
		logrus.WithError(err).Error("Loading global config failed")
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
						clusterID = exporterConf.GlusterClusterID
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
		logrus.Error("No Metrics registered, Exiting..\n")

	}

}
