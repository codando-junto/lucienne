resource "vultr_instance" "web_server" {
  plan        = "vc2-1c-2gb"
  region      = "sea"
  image_id    = "docker"
  ssh_key_ids = [vultr_ssh_key.deploy_key.id]
}

resource "tls_private_key" "generated_deploy_key" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "vultr_ssh_key" "deploy_key" {
  name    = "deploy_key"
  ssh_key = tls_private_key.generated_deploy_key.public_key_openssh
}

resource "local_sensitive_file" "private_key" {
  content  = tls_private_key.generated_deploy_key.private_key_openssh
  filename = "${path.module}/deploy_key.pem"
}

