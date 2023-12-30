package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jeremyauchter/uplist/repository"
)

func main() {
	resources := repository.GetResources()
	jsonData, err := json.Marshal(resources)
	log.Println(string(jsonData))
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	resp, err := http.Post("http://localhost:8080/resources", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
}
