output "public_ip" {
  description = "Public IP of the web server"
  value       = vultr_instance.web_server.main_ip
}
