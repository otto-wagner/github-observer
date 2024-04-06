package main

import (
	"github-listener/cmd"
)

// @title           Github listener
// @version         1.0
// @description     It is a listener for Github Events
// @termsOfService  http://swagger.io/terms/
// @host      		localhost:8443
// @BasePath  		/
func main() {
	cmd.Execute()
}
