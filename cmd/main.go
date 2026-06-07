package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	cyoa "github.com/umesshk/choose-your-adventure"
)

func main() {

	file_name := flag.String("json", "gopher.json", "Provide any valid JSON file containing CYOA story ")
	flag.Parse()

	json_file, err := os.Open(*file_name)

	if err != nil {
		fmt.Println("Error Opening File ")
		os.Exit(1)
	}

	d := json.NewDecoder(json_file)

	var story cyoa.Story

	if err := d.Decode(&story); err != nil {
		fmt.Println("Error Decoding File ", err)
		os.Exit(1)
	}

	fmt.Println(story)
}
