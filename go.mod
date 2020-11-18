module github.com/reynencourt/gluster-prometheus

go 1.14

replace github.com/reynencourt/gluster-prometheus => github.com/reynencourt/gluster-prometheus v0.1.1

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/gluster/glusterd2 v0.0.0-20181211075249-a9044cb33d93
	github.com/pborman/uuid v1.2.0 // indirect
	github.com/prometheus/client_golang v1.8.0
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.15.0
	github.com/prometheus/procfs v0.2.0 // indirect
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/viper v1.7.0 // indirect
	github.com/stretchr/testify v1.6.1 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)

replace (
	github.com/prometheus/client_model => github.com/prometheus/client_model v0.2.0
    github.com/prometheus/common => github.com/prometheus/common v0.14.0
    github.com/prometheus/client_golang => github.com/prometheus/client_golang v1.8.0
    github.com/prometheus/prometheus => github.com/prometheus/client_golang v1.8.2-0.20200724121523-657ba532e42f // indirect
)