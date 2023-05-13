# Terraform Provider Google Apps Script

This repo is made to handle Google Apps Script project 

## Requirements
- Golang

## Build provider

Run the following command to build the provider

```shell
$ go install
```

## Init
- On your Google Cloud Platform account, [enable](https://console.cloud.google.com/marketplace/product/google/script.googleapis.com?q=search&referrer=search) the apps script api
- Create an oauth consent screen [here](https://console.cloud.google.com/apis/credentials/consent?project=parcours-client-development)
- Create new oauth client id [here](https://console.cloud.google.com/apis/credentials?project=parcours-client-development)
- Download the json credentials file for your created oauth client id
- On your Google Apps Script [settings page](https://script.google.com/home/usersettings), enable API.
- Go into the init folder
- Download your oauth json file [here](https://console.cloud.google.com/apis/api/script.googleapis.com/credentials)
- Copy the downloaded file into `init/credentials.json`
- Run : `go quickstart.go`
- On the opened browser, grant access to the application
- Paste the code found in the redirected url in your terminal and press enter
- You can use the generated `token.json` file as token for the googleappscript provider
- Those operations has to be done only once

## Provider Configuration
```hcl
provider "googleappscript" {
  token = "generated_token_json_file"
  credentials = "downloaded_clients_oauth_credentials"
}
```