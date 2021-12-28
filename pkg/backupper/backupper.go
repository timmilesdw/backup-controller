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
	ConfigSpec config.Spec
}

func (b Backupper) Start() {
	logrus.Info("Started cron scheduler")
	c.Start()
}

func (b Backupper) ConfigureCron() {
	for _, backup := range b.ConfigSpec.Backups {
		backup := backup
		logrus.Infof("Setting backup task %s", backup.Name)
		var databases []config.Database
		for _, d := range backup.Databases {
			db, err := b.FindDatabase(d)
			if err != nil {
				logrus.Fatal(err)
			}
			databases = append(databases, *db)
		}
		storage, err := b.FindStorage(backup.Storage)
		if err != nil {
			logrus.Fatal(err)
		}
		c.AddFunc(
			backup.Schedule,
			func() {
				logrus.Printf("Starting task %s", backup.Name)
				for _, database := range databases {
					metrics.RunningBackups.WithLabelValues(backup.Name).Inc()
					err := b.BackupDatabase(database, *storage)
					if err != nil {
						metrics.RunningBackups.WithLabelValues(backup.Name).Dec()
						metrics.FailedBackups.WithLabelValues(backup.Name).Inc()
						logrus.Fatalf("Backup: [%s]: %s", backup.Name, err)
					}
					logrus.Infof(
						"Backup: [%s]: Database %s successfully backupped to %s",
						backup.Name,
						database.Name,
						storage.Name,
					)
					metrics.RunningBackups.WithLabelValues(backup.Name).Dec()
					metrics.SuccessfullBackups.WithLabelValues(backup.Name).Inc()
				}
			},
		)
	}
}

func (b Backupper) FindDatabase(d config.DatabaseElement) (*config.Database, error) {
	for _, database := range b.ConfigSpec.Databases {
		if database.Name == d.Name {
			return &database, nil
		}
	}
	return nil, fmt.Errorf("database %s not found", d.Name)
}

func (b Backupper) FindStorage(d config.StorageElement) (*config.Storage, error) {
	for _, storage := range b.ConfigSpec.Storages {
		if storage.Name == d.Name {
			return &storage, nil
		}
	}
	return nil, fmt.Errorf("storage %s not found", d.Name)
}

func (b Backupper) BackupDatabase(d config.Database, s config.Storage) error {
	storage := exporters.S3{
		Endpoint:     s.S3.Endpoint,
		Region:       s.S3.Region,
		Bucket:       s.S3.Bucket,
		AccessKey:    s.S3.AccessKey,
		ClientSecret: s.S3.ClientSecret,
		UseSSL:       false,
	}
	if d.Type == "postgres" {
		db := exporters.Postgres{
			Name:     d.Name,
			Host:     d.Host,
			Port:     d.Port,
			DB:       d.DB,
			Username: d.User,
			Password: d.Password,
		}
		err := db.Export().To(d.Name+d.DB+"_backups/", &storage)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("db type not found")
}
