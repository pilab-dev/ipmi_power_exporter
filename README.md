# IPMI power consumption exporter for Prometheus

Now, you can enable and start the service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable ipmi_power_exporter.service
sudo systemctl start ipmi_power_exporter.service
```

```bash
sudo systemctl status ipmi_power_exporter.service
```
And to ensure that the service starts on boot, use:

```bash
sudo systemctl is-enabled ipmi_power_exporter.service
```
This will output enabled if the service is set to start on boot.

Remember to replace the ExecStart path in the systemd service file with the actual path where you placed the ipmi_power_exporter binary. Also, ensure that the user specified in the service file has the necessary permissions to execute the ipmitool command and access the IPMI interface.


---

### Install the Package

This package is build for archlinux primarly, if you need another packaging, leave me a message ;)

```bash
sudo pacman -U ipmi_power_exporter-0.1.0-1-any.pkg.tar.zst
```


```yaml
# Configuration for IPMI exporter
# ipmi:
#   # The IPMI host to connect to
#   host: "192.168.1.100"
#   # The username for IPMI authentication
#   user: "admin"
#   # The password for IPMI authentication
#   pass: "secret"

# Configuration for logging
log:
  # The log level (debug, info, warn, error)
  level: "info"
  # The log format (json, pretty)
  format: "pretty"

# metrics exporter port
server:
  address: "127.0.0.1"
  port: 9897

# Interval between calling IPMI for new data
collect:
  interval: 10s
```

Make sure to replace the example values with your actual IPMI host, username, and password. The log level and format can be set to debug, info, warn, error, and pretty, respectively.
