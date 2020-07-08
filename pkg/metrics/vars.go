package metrics

import "time"

// These vars will be overwritten by the ones in the config file
var (
	defaultGlusterd1Workdir = ""
	defaultGlusterd2Workdir = ""
)

var (
	defaultInterval time.Duration = 5
	clusterID                     = ""
)
