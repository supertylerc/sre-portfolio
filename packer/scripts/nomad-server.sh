#!/bin/bash
set -e
cat <<EOF > /etc/nomad.d/server.hcl
server {
  enabled          = true
  bootstrap_expect = 1
}
EOF
