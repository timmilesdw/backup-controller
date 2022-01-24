package config

type Config struct {
	Logger    Logger    `yaml:"logger" validate:"required"`
	UI        UI        `yaml:"ui" validate:"required"`
	Metrics   Metrics   `yaml:"metrics" validate:"required"`
	Backupper Backupper `yaml:"backupper" validate:"required"`
}

type Logger struct {
	LogLevel string `yaml:"level" validate:"required,oneof='debug' 'info' 'warn' 'error' "`
}

type UI struct {
	Enabled  *bool  `yaml:"enabled" validate:"required"`
	Port     int    `yaml:"port" validate:"required,lte=65535,gte=1024"`
	BasePath string `yaml:"basePath" validate:"required,startswith=/"`
}

type Metrics struct {
	Enabled  *bool  `yaml:"enabled" validate:"required"`
	Port     int    `yaml:"port" validate:"required,lte=65535,gte=1024"`
	BasePath string `yaml:"basePath" validate:"required,startswith=/"`
	Path     string `yaml:"path" validate:"required,startswith=/"`
}

type Backupper struct {
	Enabled   *bool     `yaml:"enabled" validate:"required"`
	Storers   Storers   `yaml:"storers" validate:"required,dive"`
	Exporters Exporters `yaml:"exporters" validate:"required,dive"`
	Cronjobs  []Cronjob `yaml:"cronjobs" validate:"required,dive"`
}

type Exporters struct {
	Postgres []PostgresExporter `yaml:"postgres"`
}

type PostgresExporter struct {
	Name     string `yaml:"name" validate:"required"`
	Host     string `yaml:"host" validate:"required"`
	Port     string `yaml:"port" validate:"required"`
	DB       string `yaml:"db" validate:"required"`
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	Method   struct {
		Type    string   `yaml:"type" validate:"required,eq=pg_dump"`
		Options []string `yaml:"options"`
	}
}
type Cronjob struct {
	Name     string          `yaml:"name" validate:"required"`
	Schedule string          `yaml:"schedule" validate:"required"`
	Exporter ExporterElement `yaml:"exporter" validate:"required,dive"`
	Storers  StorerElement   `yaml:"storer" validate:"required,dive"`
}

type ExporterElement Element

type StorerElement Element

type Element struct {
	Name string `yaml:"name" validate:"required"`
	Type string `yaml:"type" validate:"required"`
}

type Storers struct {
	S3Storer []S3 `yaml:"s3"`
}

type S3 struct {
	Name         string `yaml:"name" validate:"required"`
	Endpoint     string `yaml:"endpoint" validate:"required"`
	Region       string `yaml:"region" validate:"required"`
	Bucket       string `yaml:"bucket" validate:"required"`
	AccessKey    string `yaml:"access-key" validate:"required"`
	ClientSecret string `yaml:"client-secret" validate:"required"`
}
