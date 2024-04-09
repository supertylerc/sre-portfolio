#!/bin/bash
set -e
# Base package updates/installs
dnf update -y
dnf install -y dnf-plugins-core curl jq net-tools gzip

# Install tailscale
dnf config-manager -y --add-repo https://pkgs.tailscale.com/stable/fedora/tailscale.repo
dnf -y install tailscale
systemctl enable tailscaled

# Download and Install consul
useradd consul -d /home/consul -s /bin/false
curl -fL \
  "https://releases.hashicorp.com/consul/${CONSUL_VERSION}/consul_${CONSUL_VERSION}_linux_$( [ $(uname -m) = aarch64 ] && echo arm64 || echo amd64).zip" \
  | gunzip -> /usr/local/bin/consul
chmod a+x /usr/local/bin/consul
consul -version
mkdir -p /etc/consul.d/tls
mkdir -p /opt/consul/data
chown consul:consul /opt/consul/data

cat <<'EOF' > /etc/systemd/system/consul.service
[Unit]
Description="HashiCorp Consul - A service mesh solution"
Documentation=https://www.consul.io/
Requires=network-online.target
After=network-online.target sys-subsystem-net-devices-tailscale0.device tailscaled.service
ConditionFileNotEmpty=/etc/consul.d/consul.hcl

[Service]
EnvironmentFile=-/etc/consul.d/consul.env
User=consul
Group=consul
ExecStart=/usr/local/bin/consul agent -config-dir=/etc/consul.d/
ExecReload=/bin/kill --signal HUP $MAINPID
KillMode=process
KillSignal=SIGTERM
Restart=on-failure
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target

[Service]
TimeoutStopFailureMode=abort
EOF
systemctl enable consul

cat <<EOF > /etc/consul.d/consul.hcl
protocol = 3
tls {
  defaults {
    verify_incoming = true
    verify_outgoing = true
    ca_file = "/etc/consul.d/tls/consul-agent-ca.pem"
  }
  internal_rpc {
    verify_server_hostname = true
  }
}
ports {
  https = 8501
}
data_dir = "/opt/consul/data"
log_level = "INFO"
ui_config {
  enabled = true
}
client_addr = "{{ GetInterfaceIP \"tailscale0\" }} 127.0.0.1"
recursors = ["1.1.1.1", "8.8.8.8"]
bind_addr = "{{ GetInterfaceIP \"tailscale0\" }}"
advertise_addr = "{{ GetInterfaceIP \"tailscale0\" }}"
# Development only
retry_join = ["consul-${CONSUL_DC}-0"]
EOF

# Setup systemd-resolved to forward *.consul. to consul
mkdir -p /etc/systemd/resolved.conf.d/
cat <<EOF > /etc/systemd/resolved.conf.d/consul.conf
[Resolve]
DNS=127.0.0.1:8600
DNSSEC=false
Domains=~consul
EOF
systemctl restart systemd-resolved
