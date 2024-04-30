# k8s Lab on KVM/Libvirt/Qemu

This is a basic configuration of my home k8s lab, which runs on a 64GB RAM, 16
thread mini PC.

To use, you should have OpenTofu installed, export the `TF_VAR_join_token` env
var with a value obtained from `kubeadm token generate`, and modify the
`locals{}` in `main.tf`.  From there, the normal `terraform init`,
`terraform plan`, and `terraform apply` workflow can be followed.  Multiple
dirs for multiple labs can be leveraged as well as the lab configuration
itself is a module.  With a small bit of work, this could be updated to
support Tofu/Terraform workspaces to improve lab switching.

This lab will create a cluster with Cilium as the CNI, operating in
non-kube-proxy mode.

> NB: Right now, this assumes that it runs on the node that has
> libvirt/qemu/kvm.  Adjustments would be needed to support remote
> hypervisors.

> :warn: cloud-init is hard-coded in the module with my username and
> password.  You should definitely change those if you plan to use this.  An
> improvement would be to pass those values in as variables, but since this
> is for me, it's good enough for now.
