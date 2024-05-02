locals {
  nodes = provider::go::expandnodes(var.nodes)
}
