#!/bin/bash
set -e
# Configure Consul Client
cat <<EOF > /etc/consul.d/client.hcl
auto_encrypt = {
  tls = true
}
EOF

# Download and Install Nomad
curl -fL \
  "https://releases.hashicorp.com/nomad/${NOMAD_VERSION}/nomad_${NOMAD_VERSION}_linux_$( [ $(uname -m) = aarch64 ] && echo arm64 || echo amd64).zip" \
  | gunzip > /usr/local/bin/nomad
chmod u+x /usr/local/bin/nomad
nomad -version

mkdir -p /etc/nomad.d
mkdir -p /opt/nomad/data/plugins/
dnf install -y podman
systemctl enable podman
curl -fL \
  "https://releases.hashicorp.com/nomad-driver-podman/${NOMAD_PODMAN_VERSION}/nomad-driver-podman_${NOMAD_PODMAN_VERSION}_linux_amd64.zip" \
  | gunzip -> /opt/nomad/data/plugins/nomad-driver-podman
chmod u+x /opt/nomad/data/plugins/nomad-driver-podman
cat <<EOF > /etc/nomad.d/common.hcl
bind_addr = "{{ GetInterfaceIP \"tailscale0\" }}"
datacenter = "${CONSUL_DC}"
data_dir = "/opt/nomad/data"
advertise {
  http = "{{ GetInterfaceIP \"tailscale0\" }}"
  rpc  = "{{ GetInterfaceIP \"tailscale0\" }}"
  serf = "{{ GetInterfaceIP \"tailscale0\" }}"
}
plugin "nomad-driver-podman" {
  config {
    socket_path = "unix:///run/podman/podman.sock"
  }
}
EOF

cat <<'EOF' > /etc/systemd/system/nomad.service
[Unit]
Description=Nomad
Documentation=https://nomadproject.io/docs/
Wants=network-online.target consul.service
After=network-online.target consul.service sys-subsystem-net-devices-tailscale0.device tailscaled.service

[Service]
EnvironmentFile=-/etc/nomad.d/nomad.env
ExecReload=/bin/kill -HUP $MAINPID
ExecStart=/usr/local/bin/nomad agent -config /etc/nomad.d
KillMode=process
KillSignal=SIGINT
LimitNOFILE=65536
LimitNPROC=infinity
Restart=on-failure
RestartSec=2

TasksMax=infinity
OOMScoreAdjust=-1000

[Install]
WantedBy=multi-user.target

[Service]
TimeoutStopFailureMode=abort
EOF

modprobe br_netfilter
cat <<EOF > /etc/modules-load.d/10-bridge.conf
br_netfilter
EOF
modprobe br_netfilter

cat <<EOF > /etc/sysctl.d/10-bridge.conf
net.bridge.bridge-nf-call-arptables = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF
echo 1 > /proc/sys/net/bridge/bridge-nf-call-arptables
echo 1 > /proc/sys/net/bridge/bridge-nf-call-ip6tables
echo 1 > /proc/sys/net/bridge/bridge-nf-call-iptables

curl -L -o cni-plugins.tgz "https://github.com/containernetworking/plugins/releases/download/v${CNI_VERSION}/cni-plugins-linux-$( [ $(uname -m) = aarch64 ] && echo arm64 || echo amd64)"-v${CNI_VERSION}.tgz
mkdir -p /opt/cni/bin
tar -C /opt/cni/bin -xzf cni-plugins.tgz

systemctl enable nomad
