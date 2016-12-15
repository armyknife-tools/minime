variable "region" {
  default = "us-central1"
}

variable "region_zone" {
  default = "us-central1-b"
}

variable "region_zone_2" {
  default = "us-central1-c"
}

variable "project_name" {
  description = "The ID of the Google Cloud project"
}

variable "credentials_file_path" {
  description = "Path to the JSON file used to describe your account credentials"
  default     = "~/.gcloud/Terraform.json"
}
