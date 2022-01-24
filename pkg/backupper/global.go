package backupper

import "github.com/timmilesdw/backup-controller/pkg/config"

var Exporters []Exporter
var Storers []Storer

func PopulateExporters(d config.Exporters) {
	// Populate with postgres databases
	for _, v := range d.Postgres {
		Exporters = append(Exporters, Postgres{
			Name: v.Name,
			Host: v.Host,
			Port: v.Port,
			// DB Name
			DB: v.DB,
			// Connection Username
			Username: v.User,
			// Password
			Password: v.Password,
			Method: struct {
				Type    string
				Options []string
			}{
				Type:    v.Method.Type,
				Options: v.Method.Options,
			},
		})
	}
}

func PopulateStorers(d config.Storers) {
	// Populate with s3 storers
	for _, v := range d.S3Storer {
		Storers = append(Storers, &S3{
			Endpoint:     v.Endpoint,
			Region:       v.Region,
			Bucket:       v.Bucket,
			AccessKey:    v.AccessKey,
			ClientSecret: v.ClientSecret,
			UseSSL:       false,
			Name:         v.Name,
		})
	}
}
