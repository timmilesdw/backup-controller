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
	Enabled  bool   `yaml:"enabled" validate:"required"`
	Port     int    `yaml:"port" validate:"required,lte=65535,gte=1024"`
	BasePath string `yaml:"basePath" validate:"required,startswith=/"`
}

type Metrics struct {
	Enabled  bool   `yaml:"enabled" validate:"required"`
	Port     int    `yaml:"port" validate:"required,lte=65535,gte=1024"`
	BasePath string `yaml:"basePath" validate:"required,startswith=/"`
	Path     string `yaml:"path" validate:"required,startswith=/"`
}

type Backupper struct {
	Storages  Storages  `yaml:"storages" validate:"required,dive"`
	Databases Databases `yaml:"databases" validate:"required,dive"`
	Cronjobs  []Cronjob `yaml:"cronjobs" validate:"required,dive"`
}

type Databases struct {
	Postgres []PostgresDatabase `yaml:"postgres"`
}

type PostgresDatabase struct {
	Name     string `yaml:"name" validate:"required"`
	Host     string `yaml:"host" validate:"required"`
	Port     string `yaml:"port" validate:"required"`
	DB       string `yaml:"db" validate:"required"`
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	Method   struct {
		Type    string   `yaml:"type" validate:"required"`
		Options []string `yaml:"options"`
	}
}
type Cronjob struct {
	Name     string          `yaml:"name" validate:"required"`
	Schedule string          `yaml:"schedule" validate:"required"`
	Database DatabaseElement `yaml:"database" validate:"required,dive"`
	Storage  StorageElement  `yaml:"storage" validate:"required"`
}

type DatabaseElement Element

type StorageElement Element

type Element struct {
	Name string `yaml:"name" validate:"required"`
	Type string `yaml:"type" validate:"required"`
}

type Storages struct {
	S3Storage []S3 `yaml:"s3"`
}

type S3 struct {
	Name         string `yaml:"name" validate:"required"`
	Endpoint     string `yaml:"endpoint" validate:"required"`
	Region       string `yaml:"region" validate:"required"`
	Bucket       string `yaml:"bucket" validate:"required"`
	AccessKey    string `yaml:"access-key" validate:"required"`
	ClientSecret string `yaml:"client-secret" validate:"required"`
}
