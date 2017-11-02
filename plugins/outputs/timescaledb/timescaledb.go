package timescaledb

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/outputs"
)

type TimescaleDB struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
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

}

// Close will terminate the session to the backend, returning error if an issue arises
func (t *TimescaleDB) Close() error {
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

func (t *TimescaleDB) Write(metrics []Metric) string {
	if len(metrics) == 0 {
		return nil
	}

	for _, metric := range metrics {

	}
}

func init() {
	outputs.Add("timescaledb", func() telegraf.Output { return &TimescaleDB{} })
}
