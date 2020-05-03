# Chefapi_client

The chefapi_client code is imported by other chefapi go modules.  It sets up a go-chef/chef Client 
with the correct base url, depending on the desired use, adds credentials and handles
TLS certificates.  

## Use
import "github/MarkGibbons/chefapi_client

Functions
* Client()  Creates a client. This client is suitable for global Chef Server API REST calls.
* Clients(Org string) Creates a client for Chef Server API REST calls for the specified organization.

Both Client and ClientOrg expect environment variables to be defined in order to specify 
the credentials needed to talk to the Chef Infra Server.

| Variable     | Purpose | Required |
|--------------|----------------------------------|------|
|CHEFAPICHEFUSER| Specify the Chef user name       | Yes |
|CHEFAPIKEYFILE| File containing the user's private Chef key| Yes |
|CHEFAPICHRURL| The Chef Server URL | Yes |
|CHEFAPICERTFILE| SSL Certificate for the Chef Server, used for self signed certificates | No |
