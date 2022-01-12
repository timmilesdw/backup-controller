package backupper

import (
	"fmt"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/timmilesdw/backup-controller/pkg/config"
	"github.com/timmilesdw/backup-controller/pkg/exporters"
	"github.com/timmilesdw/backup-controller/pkg/metrics"
)

var c = cron.New()

type Backupper struct {
	ConfigSpec config.Backupper
}

func (b Backupper) Start() {
	logrus.Info("Started cron scheduler")
	c.Start()
}

func (b Backupper) ConfigureCron() {
	for _, cronjob := range b.ConfigSpec.Cronjobs {
		cron := cronjob
		logrus.Infof("Setting backup task %s", cron.Name)
		exporter, err := b.FindExporter(cron.Exporter)
		if err != nil {
			logrus.Fatal(err)
		}
		storer, err := b.FindStorer(cron.Storers)
		if err != nil {
			logrus.Fatal(err)
		}
		c.AddFunc(
			cron.Schedule,
			func() {
				logrus.Printf("Starting task %s", cron.Name)
				metrics.RunningBackups.WithLabelValues(cron.Name).Inc()
				err := b.BackupDatabase(exporter, storer)
				if err != nil {
					metrics.RunningBackups.WithLabelValues(cron.Name).Dec()
					metrics.FailedBackups.WithLabelValues(cron.Name).Inc()
					logrus.Fatalf("Backup: [%s]: %s", cron.Name, err)
				}
				logrus.Infof(
					"Backup: [%s]: Database %s successfully backupped to %s",
					cron.Name,
					exporter.GetName(),
					storer.GetName(),
				)
				metrics.RunningBackups.WithLabelValues(cron.Name).Dec()
				metrics.SuccessfullBackups.WithLabelValues(cron.Name).Inc()
			},
		)
	}
}

func (b Backupper) FindExporter(d config.ExporterElement) (exporters.Exporter, error) {
	for _, e := range exporters.Exporters {
		if d.Type == e.GetType() {
			if d.Name == e.GetName() {
				return e, nil
			}
		}
	}
	return nil, fmt.Errorf("exporter %s not found", d.Name)
}

func (b Backupper) FindStorer(d config.StorerElement) (exporters.Storer, error) {
	for _, s := range exporters.Storers {
		if d.Type == s.GetType() {
			if d.Name == s.GetName() {
				return s, nil
			}
		}
	}
	return nil, fmt.Errorf("storer %s not found", d.Name)
}

func (b Backupper) BackupDatabase(e exporters.Exporter, s exporters.Storer) error {
	folderName := fmt.Sprintf("%v_%v_backups/", e.GetType(), e.GetName())
	err := e.Export().To(folderName, s)
	if err != nil {
		return err
	}
	return nil
}
