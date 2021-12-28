package exporters

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	// PGDumpCmd is the path to the `pg_dump` executable
	PGDumpCmd = "pg_dump"
)

// Postgres is an `Exporter` interface that backs up a Postgres database via the `pg_dump` command
type Postgres struct {
	Name string
	// DB Host (e.g. 127.0.0.1)
	Host string
	// DB Port (e.g. 5432)
	Port string
	// DB Name
	DB string
	// Connection Username
	Username string
	// Password
	Password string
	// Extra pg_dump options
	// e.g []string{"--inserts"}
	Options []string
}

// Export produces a `pg_dump` of the specified database, and creates a gzip compressed tarball archive.
func (x Postgres) Export() *ExportResult {
	result := &ExportResult{MIME: "application/x-tar"}
	result.Path = fmt.Sprintf(`bu_%v_%v.sql.tar.gz`, x.Name+x.DB, time.Now().Unix())
	options := append(x.dumpOptions(), "-v", "-Fc", fmt.Sprintf(`-f%v`, result.Path))
	cmd := exec.Command(PGDumpCmd, options...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "PGPASSWORD="+x.Password)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	logrus.Println(stderr.String())
	if err != nil {
		result.Error = makeErr(err, stderr.String())
	}
	logrus.Println(out.String())
	return result
}

func (x Postgres) dumpOptions() []string {
	options := x.Options

	if x.DB != "" {
		options = append(options, fmt.Sprintf(`-d%v`, x.DB))
	}

	if x.Host != "" {
		options = append(options, fmt.Sprintf(`-h%v`, x.Host))
	}

	if x.Port != "" {
		options = append(options, fmt.Sprintf(`-p%v`, x.Port))
	}

	if x.Username != "" {
		options = append(options, fmt.Sprintf(`-U%v`, x.Username))
	}

	return options
}
