#!/bin/bash
set -e
cat <<EOF > /etc/nomad.d/client.hcl
client {
  enabled           = true
  network_interface = "tailscale0"
  network_speed     = 1000
}
EOF
