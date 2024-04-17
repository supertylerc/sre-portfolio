data "http" "public_ip" {
  url = "https://ipinfo.io/json"
  request_headers = {
    Accept = "application/json"
  }
}

module "vpc" {
  source = "terraform-do-modules/vpc/digitalocean"

  name        = var.labels["name"]
  environment = var.labels["environment"]
  region      = var.region
  ip_range    = "10.42.0.0/16"
  managedby   = var.labels["managed_by"]
}

module "spaces" {
  source = "terraform-do-modules/spaces/digitalocean"

  name          = var.spaces_name
  acl           = "private"
  force_destroy = false
  region        = var.spaces_region
  environment   = var.labels["environment"]
  managedby     = var.labels["managed_by"]
}

module "droplet" {
  source = "terraform-do-modules/droplet/digitalocean"

  droplet_count = var.droplet_count
  name          = var.labels["name"]
  environment   = var.labels["environment"]
  region        = var.region
  image_name    = var.image
  vpc_uuid      = module.vpc.id
  key_path      = var.ssh_public_key_path
  key_name      = var.ssh_key_name
  user_data     = file("user-data.sh")
  inbound_rules = []
  outbound_rule = var.outbound_rules
  managedby     = var.labels["managed_by"]
}

module "laptop_firewall" {
  source = "terraform-do-modules/firewall/digitalocean"

  name          = var.labels["name"]
  environment   = var.labels["environment"]
  allowed_ip    = [jsondecode(data.http.public_ip.body)["ip"]]
  allowed_ports = var.laptop_ports
  droplet_ids   = module.droplet.id
}
