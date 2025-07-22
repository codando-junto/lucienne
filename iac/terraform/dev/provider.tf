
terraform {
  required_providers {
    vultr = {
      source  = "vultr/vultr"
      version = "2.26.0"
    }
  }
}

# Configure the Vultr Provider
provider "vultr" {
  rate_limit  = 100
  retry_limit = 3
}

