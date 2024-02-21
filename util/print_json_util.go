package util

import (
	"encoding/json"
	"fmt"
	"log"
)

func PrintJSON(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	fmt.Println(string(jsonData))
}
