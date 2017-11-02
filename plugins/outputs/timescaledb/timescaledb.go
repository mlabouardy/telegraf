package timescaledb

import (
	"database/sql"
	"fmt"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/outputs"
)

type TimescaleDB struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	Client   *sql.DB
}

var sampleConfig = `
  ## DNS name of the TimescaleDB server
  host = "opentsdb.example.com"

  ## Port of the TimescaleDB server
  port = 5432

  ## TimescaleDB credentials
  username = "postgres"
  password = "postgres"

  ## The target database for metrics.
  database = "telegraf" # required
`

// Connect initiates the connection to TimescaleDB server
func (t *TimescaleDB) Connect() error {
	if t.Host == "" || t.Database == "" {
		return fmt.Errorf("TimescaleDB host or database is not defined")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", t.Host, t.Port, t.Username, t.Password, t.Database)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return fmt.Errorf("TimescaleDB connection failed: %s", err)
	}

	t.Client = db

	return nil
}

// Close will terminate the session to the backend, returning error if an issue arises
func (t *TimescaleDB) Close() error {
	t.Client.Close()
	return nil
}

// SampleConfig returns the formatted sample configuration for the plugin
func (t *TimescaleDB) Description() string {
	return "Configuration for timescaledb server to send metrics to"
}

// Description returns the human-readable function definition of the plugin
func (t *TimescaleDB) SampleConfig() string {
	return sampleConfig
}

func (t *TimescaleDB) Write(metrics []telegraf.Metric) error {
	if len(metrics) == 0 {
		return nil
	}

	for _, metric := range metrics {
		fmt.Println(metric.Fields())
	}

	return nil
}

func init() {
	outputs.Add("timescaledb", func() telegraf.Output { return &TimescaleDB{} })
}
