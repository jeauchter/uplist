package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jeremyauchter/uplist/pkg/client"
	"github.com/jeremyauchter/uplist/pkg/config"
)

func main() {
	config := config.NewConfig()
	etsyApi := client.NewEtsyAPI(*config)

	if len(os.Args) == 2 && os.Args[1] == "authorize" {
		fmt.Println("Authorizing the app")
		status, err := etsyApi.AuthorizeApp()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(status)
	}

	status, err := etsyApi.RefreshToken()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(status)

}
