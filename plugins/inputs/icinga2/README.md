# Example Input Plugin

This plugin gather services & hosts status using Icinga2 Remote API.

The icinga2 plugin uses the icinga2 remote API to gather status on running
services and hosts. You can read Icinga2's documentation for their remote API
[here](https://docs.icinga.com/icinga2/latest/doc/module/icinga2/chapter/icinga2-api)

### Configuration:

```toml
# Description
[[inputs.icinga2]]
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
```

### Measurements & Fields:

- ll measurements have the following fields:
    - name (string)
    - state (int)

### Tags:

- All measurements have the following tags:
    - check_command
    - display_name
    - source

### Sample Queries:

```
SELECT * FROM "icinga2_services_status" WHERE state = 0 AND time > now() - 24h // Service with OK status
SELECT * FROM "icinga2_services_status" WHERE state = 1 AND time > now() - 24h // Service with WARNING status
SELECT * FROM "icinga2_services_status" WHERE state = 2 AND time > now() - 24h // Service with Critical status
SELECT * FROM "icinga2_services_status" WHERE state = 3 AND time > now() - 24h // Service with UNKNOWN status
```

### Example Output:

```
$ ./telegraf -config telegraf.conf -input-filter icinga2 -test
icinga2_hosts_status,display_name=router-fr.eqx.fr,check_command=hostalive-custom,host=test-vm,source=https://localhost:5665 name="router-fr.eqx.fr",state=0 1492021603000000000
```