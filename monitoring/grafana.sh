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

sudo rm -rf prometheus-2.31.0*

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
curl -LO https://github.com/prometheus/node_exporter/releases/download/v1.5.0/node_exporter-1.5.0.linux-amd64.tar.gz

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

sudo rm -rf node_exporter-1.5.0*

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

# Download Alert Manager
curl -LO https://github.com/prometheus/alertmanager/releases/download/v0.25.0/alertmanager-0.25.0.linux-amd64.tar.gz

# Extract Alert Manager
tar -xvf alertmanager-0.25.0.linux-amd64.tar.gz

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
cp alertmanager-example.yaml /etc/alertmanager/alertmanager.yml

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

sudo rm -rf alertmanager-0.25.0*

# Create directory for amtool configuration
sudo mkdir -p /etc/amtool

# Configure amtool to use Alertmanager running on localhost
echo "alertmanager.url: http://localhost:9093" > /etc/amtool/config.yml

# Verify amtool configuration
amtool config show

# Create alerting rules file with multiple alerting rules
cp rules-example.yaml /etc/prometheus/rules.yml

# Trigger test alert
curl -H "Content-Type: application/json" -XPOST http://localhost:9093/api/v2/alerts -d '[{"labels": {"alertname":"TestAlert"},"generatorURL":"http://localhost:9090/graph","annotations": {"summary":"TestAlert"}}]'
# Check alert was queued

curl http://localhost:9093/api/v2/alerts

# Configure targets in /etc/prometheus/prometheus.yml
cp prometheus-example.yaml /etc/prometheus/prometheus.yml

# Restart Prometheus to apply new configuration changes.
sudo systemctl restart prometheus
sudo systemctl status prometheus
journalctl -u prometheus -f

# Check Prometheus work correctly and all tne targets are reachable: http://45.250.253.76:9090/targets

# Configure journal logging to be persistent across reboots.
sudo tee <<EOF >/dev/null /etc/systemd/journald.conf
Storage=persistent
EOF

# Dashboard template to use for SUI validator in Grafana: https://grafana.com/grafana/dashboards/18297-sui-validator-dashboard-1-0
# How to add Data Source and Dashboard to Grafana: https://www.digitalocean.com/community/tutorials/how-to-monitor-mongodb-with-grafana-and-prometheus-on-ubuntu-20-04