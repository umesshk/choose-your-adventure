package main

import (
	"flag"
	"fmt"
	"github.com/umesshk/choose-your-adventure/internal"
	"os"
)

func main() {

	file_name := flag.String("json", "gopher.json", "Provide any valid JSON file containing CYOA story ")
	flag.Parse()

	json_file, err := os.Open(*file_name)

	if err != nil {
		fmt.Println("Error Opening File ")
		os.Exit(1)
	}

	story, err := internal.MakeStory(json_file)

	if err != nil {
		fmt.Println("Error Occured : ", err)
		os.Exit(1)
	}

	fmt.Println(story)

	internal.ServerPage()
}
