#!/bin/bash

# Install Grafana

# Add Grafana's GPG key to apt
wget -q -O - https://packages.grafana.com/gpg.key | sudo apt-key add -

# Add Grafana's apt repository
sudo add-apt-repository "deb https://packages.grafana.com/oss/deb stable main"

# Update apt and install Grafana
sudo apt update
sudo apt install grafana

# Reload systemd and start Grafana
sudo systemctl daemon-reload
sudo systemctl start grafana-server

# Enable Grafana to start on boot
sudo systemctl enable grafana-server

# Check the status of Grafana
sudo systemctl status grafana-server

# Update grafana credentials using web on http://<your ip>:3000

# Disable user signup/registration
sudo sed -i 's/;allow_sign_up = true/allow_sign_up = false/' /etc/grafana/grafana.ini

# Enable anonymous access
sudo sed -i 's/;enabled = false/enabled = false/' /etc/grafana/grafana.ini

# Restart grafana server
sudo systemctl restart grafana-server
sudo systemctl status grafana-server

# Prometheus installation and configuration

# Update apt package index
sudo apt update

# Create necessary directories for Prometheus
sudo mkdir -p /etc/prometheus
sudo mkdir -p /var/lib/prometheus

# Download and extract Prometheus binary
wget https://github.com/prometheus/prometheus/releases/download/v2.31.0/prometheus-2.31.0.linux-amd64.tar.gz
tar -xvf prometheus-2.31.0.linux-amd64.tar.gz
cd prometheus-2.31.0.linux-amd64

# Move Prometheus binaries to /usr/local/bin/ directory
sudo mv prometheus promtool /usr/local/bin/

# Move console files to /etc/prometheus/ directory
sudo mv consoles/ console_libraries/ /etc/prometheus/

# Move Prometheus configuration file to /etc/prometheus/ directory
sudo mv prometheus.yml /etc/prometheus/prometheus.yml

# Verify Prometheus version
prometheus --version

# Create Prometheus system user and group
sudo groupadd --system prometheus
sudo useradd -s /sbin/nologin --system -g prometheus prometheus

# Set proper permissions for Prometheus directories
sudo chown -R prometheus:prometheus /etc/prometheus/ /var/lib/prometheus/
sudo chmod -R 775 /etc/prometheus/ /var/lib/prometheus/

# Configure and start Prometheus as a system service

# Create Prometheus service file in /etc/systemd/system/ directory
echo "[Unit]
Description=Prometheus
Wants=network-online.target
After=network-online.target

[Service]
User=prometheus
Group=prometheus
Restart=always
Type=simple
ExecStart=/usr/local/bin/prometheus \
    --config.file=/etc/prometheus/prometheus.yml \
    --storage.tsdb.path=/var/lib/prometheus/ \
    --web.console.templates=/etc/prometheus/consoles \
    --web.console.libraries=/etc/prometheus/console_libraries \
    --web.listen-address=0.0.0.0:9090

[Install]
WantedBy=multi-user.target" > /etc/systemd/system/prometheus.service

# Reload systemd daemon to load new service file
sudo systemctl daemon-reload

# Start and enable Prometheus service
sudo systemctl start prometheus
sudo systemctl enable prometheus

# Verify Prometheus service status and logs
sudo systemctl status prometheus
journalctl -u prometheus -f

# Configure firewall to allow incoming traffic on port 9090 for Prometheus

# Check the status of the firewall
sudo ufw status

# Allow incoming traffic on port 9090
sudo ufw allow 9090

# Install Node Exporter for Prometheus

# Download Node Exporter from the official website
cd /tmp
curl -LO https://github.com/prometheus/node_exporter/releases/download/v0.5.0/node_exporter-1.5.0.linux-amd64.tar.gz

# Extract Node Exporter and move it to the appropriate directory
tar -xvf node_exporter-1.5.0.linux-amd64.tar.gz
sudo mv node_exporter-1.5.0.linux-amd64/node_exporter /usr/local/bin/

# Create a system user for Node Exporter
sudo useradd -rs /bin/false node_exporter

# Create a systemd service file for Node Exporter
echo "[Unit]
Description=Node Exporter
After=network.target

[Service]
User=node_exporter
Group=node_exporter
Type=simple
ExecStart=/usr/local/bin/node_exporter

[Install]
WantedBy=multi-user.target" > /etc/systemd/system/node_exporter.service

# Reload systemd daemon to load new service file
sudo systemctl daemon-reload

# Start and enable Node Exporter service
sudo systemctl start node_exporter
sudo systemctl enable node_exporter

# Verify Node Exporter service status and logs
sudo systemctl status node_exporter
journalctl -u node_exporter -f

# Access Prometheus in a web browser at http://<SERVER_PUBLIC_IP>:9090/targets to view available targets and their status

# install alert manager

# Install Alert Manager
cd /tmp

# Download Alert Manager
curl -LO https://github.com/prometheus/alertmanager/releases/download/v0.25.0/alertmanager-0.25.0.linux-amd64.tar.gz

# Extract Alert Manager
tar -xvf node_exporter-0.25.0.linux-amd64.tar.gz

# Move Alert Manager binary to /usr/local/bin/
sudo mv alertmanager-0.25.0.linux-amd64/alertmanager /usr/local/bin/

# Create user for Alertmanager
sudo useradd -rs /bin/false alertmanager

# Create systemd service file for Alert Manager
echo "[Unit]
Description=Alert Manager
Wants=network-online.target
After=network-online.target

[Service]
Type=simple
User=alertmanager
Group=alertmanager
ExecStart=/usr/local/bin/alertmanager \
  --config.file=/etc/alertmanager/alertmanager.yml \
  --storage.path=/data/alertmanager \
  --cluster.advertise-address="127.0.0.1:9093"

Restart=always

[Install]
WantedBy=multi-user.target" > /etc/systemd/system/alertmanager.service

# Create configuration file for Alert Manager to use slack integration
echo "global:
  resolve_timeout: 1m
  slack_api_url: 'https://hooks.slack.com/services/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'

route:
  receiver: 'slack-notifications'

receivers:
- name: 'slack-notifications'
  slack_configs:
  - channel: '#sui-alerts'
    send_resolved: true
    icon_url: https://avatars3.githubusercontent.com/u/3380462
    title: |-
     [{{ .Status | toUpper }}{{ if eq .Status "firing" }}:{{ .Alerts.Firing | len }}{{ end }}] {{ .CommonLabels.alertname }} for {{ .CommonLabels.job }}
     {{- if gt (len .CommonLabels) (len .GroupLabels) -}}
       {{" "}}(
       {{- with .CommonLabels.Remove .GroupLabels.Names }}
         {{- range $index, $label := .SortedPairs -}}
           {{ if $index }}, {{ end }}
           {{- $label.Name }}="{{ $label.Value -}}"
         {{- end }}
       {{- end -}}
       )
     {{- end }}
    text: >-
     {{ range .Alerts -}}
     *Alert:* {{ .Annotations.title }}{{ if .Labels.severity }} - `{{ .Labels.severity }}`{{ end }}

     *Description:* {{ .Annotations.description }}

     *Details:*
       {{ range .Labels.SortedPairs }} • *{{ .Name }}:* `{{ .Value }}`
       {{ end }}
     {{ end }}" > /etc/alertmanager/alertmanager.yml

# Reload systemd daemon
sudo systemctl daemon-reload

# Start and enable Alert Manager service
sudo systemctl start alertmanager
sudo systemctl enable alertmanager

# Verify Alert Manager service status and logs
sudo systemctl status alertmanager
journalctl -u alertmanager -f

# Access Alert Manager in a web browser at http://<SERVER_PUBLIC_IP>:9093 to view available rules and their status

# Install amtool

# Copy amtool binary to /usr/local/bin directory
sudo cp alertmanager-0.25.0.linux-amd64/amtool /usr/local/bin/

# Create directory for amtool configuration
sudo mkdir -p /etc/amtool

# Configure amtool to use Alertmanager running on localhost
echo "alertmanager.url: http://localhost:9093" > /etc/amtool/config.yml

# Verify amtool configuration
amtool config show

# Create alerting rules file with multiple alerting rules
echo "groups:
  - name: sui-node
    rules:
    - alert: NodeDown
      expr: increase(uptime{}[5m]) == 0
      for: 1m
      labels:
        severity: critical
      annotations:
        summary: "Sui Node {{ $labels.instance }} uptime is stuck"
        description: "Sui Node {{ $labels.instance }} uptime is stuck"

    - alert: NodeCurrentRoundStuck
      expr: increase(current_round{}[5m]) == 0
      for: 1m
      labels:
        severity: critical
      annotations:
        summary: "Sui Node {{ $labels.instance }} current round is stuck"
        description: "Sui Node {{ $labels.instance }} current round is stuck"

    - alert: NodeLastCommittedRoundsStuck
      expr: increase(last_committed_round{}[5m]) == 0
      for: 1m
      labels:
        severity: critical
      annotations:
        summary: "Sui Node {{ $labels.instance }} last committed round is stuck"
        description: "Sui Node {{ $labels.instance }} last committed round is stuck"

    - alert: LastExecutedCheckpoint
      expr: increase(last_executed_checkpoint{}[5m]) == 0
      for: 1m
      labels:
        severity: critical
      annotations:
        summary: "Checkpoints are not being executed on {{ $labels.instance }}"
        description: "Checkpoints are not being executed on {{ $labels.instance }}"

  - name: system
    rules:
    - alert: HighCpuUsage
      expr: (100 - (avg by (instance) (irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)) > 80
      for: 1m
      labels:
        severity: critical
      annotations:
        summary: "High CPU usage on {{ $labels.instance }}"
        description: "CPU usage on {{ $labels.instance }} has been above 80% for the last 1 minute."

    - alert: HighMemoryUsage
      expr: (node_memory_Active_bytes / node_memory_MemTotal_bytes) > 0.8
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "High memory usage on {{ $labels.instance }}"
        description: "Memory usage on {{ $labels.instance }} has been above 80% for the last 5 minutes."

    - alert: HighDiskUsage
      expr: (node_filesystem_size_bytes{fstype="ext4"} - node_filesystem_free_bytes{fstype="ext4"}) / node_filesystem_size_bytes{fstype="ext4"} > 0.9
      for: 10m
      labels:
        severity: critical
      annotations:
        summary: "High disk usage on {{ $labels.instance }}"
        description: "Disk usage on {{ $labels.instance }} has been above 90% for the last 10 minutes."" > /etc/prometheus/rules.yml

# Configure targets in /etc/prometheus/prometheus.yml
echo "# my global config
global:
  scrape_interval: 15s # Set the scrape interval to every 15 seconds. Default is every 1 minute.
  evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1 minute.
  # scrape_timeout is set to the global default (10s).

# Alertmanager configuration
alerting:
  alertmanagers:
    - static_configs:
        - targets:
          - 127.0.0.1:9093

# Load rules once and periodically evaluate them according to the global 'evaluation_interval'.
rule_files:
  - rules.yml

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: "prometheus"

    # metrics_path defaults to '/metrics'
    # scheme defaults to 'http'.

    static_configs:
      - targets: ["localhost:9090"]

  - job_name: 'node_exporter_metrics'
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:9100']

  - job_name: 'validator_metrics'
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:9184']" > /etc/prometheus/prometheus.yml

# Restart Prometheus to apply new configuration changes.
sudo systemctl restart prometheus

# Configure journal logging to be persistent across reboots.
sudo tee <<EOF >/dev/null /etc/systemd/journald.conf
Storage=persistent
EOF