#!/bin/bash
set -e
cat <<EOF > /etc/consul.d/server.hcl
# Only 1 server for development
bootstrap_expect = 1
server = true
tls {
  defaults {
    cert_file = "/etc/consul.d/tls/${CONSUL_DC}-server-consul-0.pem"
    key_file = "/etc/consul.d/tls/${CONSUL_DC}-server-consul-0-key.pem"
  }
}
auto_encrypt {
  allow_tls = true
}
EOF

# Static certs for development
# TODO: Migrate certificates to blob storage and retrieve dynamically on new instance creation
# NB: You will need to pull the CA cert from the Consul server and upload it to all running client
#     instances.
mkdir /tmp/tls && pushd /tmp/tls
consul tls ca create
consul tls cert create -server -dc ${CONSUL_DC}
mv * /etc/consul.d/tls/
popd
chown -R consul:consul /etc/consul.d/tls
