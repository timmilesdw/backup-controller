package metrics

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/timmilesdw/backup-controller/pkg/config"
)

var (
	RunningBackups = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "backup_controller_current_running_backups",
		Help: "Current backups running",
	}, []string{"name"})
	RegisteredBackups = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "backup_controller_registered_cronjobs",
		Help: "Registered Cronjobs",
	}, []string{"name", "databases", "schedule", "storage_name"})
	RegisteredDatabases = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "backup_controller_registered_exporters",
		Help: "Registered Databases",
	}, []string{"name", "type", "host", "db", "port"})
	RegisteredStorages = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "backup_controller_registered_storages",
		Help: "Registered Storages",
	}, []string{"name", "endpoint", "bucket", "region"})
	SuccessfullBackups = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "backup_controller_successful_backups",
		Help: "Number of successful backups",
	}, []string{"name"})
	FailedBackups = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "backup_controller_failed_backups",
		Help: "Number of failed backups",
	}, []string{"name"})
)

type MetricsServer struct {
	Route string
	Port  string
}

type Metrics struct {
	Gauges   Gauges
	Counters Counters
}

type Counters struct {
	RegisteredCronjobs  prometheus.Counter
	RegisteredDatabases prometheus.Counter
	RegisteredStorages  prometheus.Counter
	SuccessfullBackups  prometheus.Counter
	FailedBackups       prometheus.Counter
}

type Gauges struct {
	RunningBackups prometheus.Gauge
}

func RegisterMetrics(cfg config.Metrics) *MetricsServer {
	ms := &MetricsServer{
		Port:  ":" + strconv.Itoa(cfg.Port),
		Route: cfg.Path,
	}
	// for _, e := range exporters.Exporters {
	// 	RegisteredDatabases.WithLabelValues(
	// 		e.GetName(), e.GetType(),
	// 		e.GetHost(), e.GetDB(),
	// 		e.GetPort(),
	// 	).Add(1)
	// }
	// for _, s3 := range cfg.Storages {
	// 	RegisteredStorages.WithLabelValues(
	// 		s3.Name,
	// 		s3.S3.Endpoint,
	// 		s3.S3.Bucket,
	// 		s3.S3.Region,
	// 	).Add(1)
	// }
	// for _, backup := range cfg.Backups {
	// 	var databases []string
	// 	for _, dbName := range backup.Databases {
	// 		databases = append(databases, dbName.Name)
	// 	}
	// 	RegisteredBackups.WithLabelValues(
	// 		backup.Name, strings.Join(databases, ","),
	// 		backup.Schedule, backup.Storage.Name,
	// 	).Add(1)
	// }
	return ms
}

func (m *MetricsServer) Start() {
	http.Handle(m.Route, promhttp.Handler())
	err := http.ListenAndServe(m.Port, nil)
	if err != nil {
		logrus.Fatal(err)
	}
}
