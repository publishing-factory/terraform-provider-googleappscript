terraform {
  required_providers {
    googleappscript = {
      source = "hashicorp.com/edu/googleappscript"
    }
  }
}

provider "googleappscript" {
  token = "generated_token"
  credentials = "your_google_oauth_credentials_json"
}


resource "googleappscript_project" "this" {
  title = "Last toto"
  parent_id = ""
}

output "test" {
  value = googleappscript_project.this
}