package icinga2

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/internal"
	"github.com/influxdata/telegraf/internal/tls"
	"github.com/influxdata/telegraf/plugins/inputs"
)

type Icinga2 struct {
	Server          string
	Filter          string
	Username        string
	Password        string
	ResponseTimeout internal.Duration
	tls.ClientConfig

	client *http.Client
}

type Result struct {
	Results []Object `json:"results"`
}

type Object struct {
	Attrs Attribute  `json:"attrs"`
	Name  string     `json:"name"`
	Joins struct{}   `json:"joins"`
	Meta  struct{}   `json:"meta"`
	Type  ObjectType `json:"type"`
}

type Attribute struct {
	CheckCommand string  `json:"check_command"`
	DisplayName  string  `json:"display_name"`
	Name         string  `json:"name"`
	State        float32 `json:"state"`
}

type ObjectType string

var sampleConfig = `
  ## Required Icinga2 server address (default: "https://localhost:5665")
  # server = "https://localhost:5665"
  
  ## Required Icinga2 object type ("services" or "hosts, default "services")
  # filter = "services"

  ## Credentials for basic HTTP authentication
  # username = "admin"
  # password = "admin"

  ## Maximum time to receive response.
  # response_timeout = "5s"

  ## Optional TLS Config
  # tls_ca = "/etc/telegraf/ca.pem"
  # tls_cert = "/etc/telegraf/cert.pem"
  # tls_key = "/etc/telegraf/key.pem"
  ## Use TLS but skip chain & host verification
  `

func (i *Icinga2) Description() string {
	return "Gather Icinga2 status"
}

func (i *Icinga2) SampleConfig() string {
	return sampleConfig
}

func (i *Icinga2) GatherStatus(acc telegraf.Accumulator, checks []Object) {
	for _, check := range checks {
		fields := make(map[string]interface{})
		tags := make(map[string]string)

		fields["name"] = check.Attrs.Name
		fields["state"] = check.Attrs.State

		tags["display_name"] = check.Attrs.DisplayName
		tags["check_command"] = check.Attrs.CheckCommand
		tags["source"] = i.Server

		acc.AddFields(fmt.Sprintf("icinga2_%s_status", i.Filter), fields, tags)
	}
}

func (i *Icinga2) createHttpClient() (*http.Client, error) {
	tlsCfg, err := i.ClientConfig.TLSConfig()
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsCfg,
		},
		Timeout: i.ResponseTimeout.Duration,
	}

	return client, nil
}

func (i *Icinga2) Gather(acc telegraf.Accumulator) error {
	if i.ResponseTimeout.Duration < time.Second {
		i.ResponseTimeout.Duration = time.Second * 5
	}

	if i.client == nil {
		client, err := i.createHttpClient()
		if err != nil {
			return err
		}
		i.client = client
	}

	url := fmt.Sprintf("%s/v1/objects/%s?attrs=name&attrs=display_name&attrs=state&attrs=check_command", i.Server, i.Filter)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(i.Username, i.Password)
	resp, err := i.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	result := Result{}
	json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}

	i.GatherStatus(acc, result.Results)

	return nil
}

func init() {
	inputs.Add("icinga2", func() telegraf.Input {
		return &Icinga2{
			Server: "https://localhost:5665",
			Filter: "services",
		}
	})
}
