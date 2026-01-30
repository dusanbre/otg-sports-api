package main

//	@title			OTG Sport API
//	@version		1.0
//	@description	REST API for accessing sports data synced from GoalServe.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	OTG Sports Support
//	@contact.email	support@otgsports.com

//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						X-API-Key
//	@description				API key for authentication

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Bearer token (use "Bearer <api_key>")

import (
	"github.com/dusanbre/otg-sports-api/cmd"
)

func main() {
	cmd.Execute()
}
