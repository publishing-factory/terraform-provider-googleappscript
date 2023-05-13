terraform {
  required_providers {
    googleappscript = {
      source = "hashicorp.com/edu/googleappscript"
    }
  }
}

provider "googleappscript" {
  token = "{\"access_token\":\"ya29.a0AWY7CklsRLIwhZ8b2aDPsFtP1W0qRKrC65EUfblD8m9ZGook8WNnvLR67wXcpLSTQflUM-BXz6X_Ql78RrwtCB1mKW0-PkYmbT7FLv1_6rYk-BbHk65zztvB3ShwOIxdgHPoUgxB1muvZE8E1V0qqtAkacSHaCgYKAfESARESFQG1tDrp1oIzJ77MQ37FnMnX5FgDtw0163\",\"token_type\":\"Bearer\",\"refresh_token\":\"1//03UwTiFg2EP7VCgYIARAAGAMSNwF-L9IrTHE4bnv1YzcjZJgMlbPxrFCEBA73wh5WMyj8k8CAvTcVFbo-K9ZfiP51QEzeElg1L-s\",\"expiry\":\"2023-05-13T16:22:12.881583+02:00\"}"
  credentials = "{\"installed\":{\"client_id\":\"555032422313-gc6n5o36h229l0g89qtop25rc4tkla6a.apps.googleusercontent.com\",\"project_id\":\"parcours-client-development\",\"auth_uri\":\"https://accounts.google.com/o/oauth2/auth\",\"token_uri\":\"https://oauth2.googleapis.com/token\",\"auth_provider_x509_cert_url\":\"https://www.googleapis.com/oauth2/v1/certs\",\"client_secret\":\"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\",\"redirect_uris\":[\"http://localhost\"]}}"
}


resource "googleappscript_project" "this" {
  title = "Last toto"
  parent_id = ""
}

output "test" {
  value = googleappscript_project.this
}