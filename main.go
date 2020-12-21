package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"vs-ru/api"
)

func main() {
	filename := flag.String("filename", "array.json", "set a JSON filename to read to")

	fileBytes, err := ioutil.ReadFile(*filename)
	if err != nil {
		log.Fatalf("unable to read file %s: %v", *filename, err)
	}

	treeValues := []int{}
	err = json.Unmarshal(fileBytes, &treeValues)
	if err != nil {
		log.Fatalf("unable to unmarshal data from file %s: %v", *filename, err)
	}

	apiService, err := api.NewAPI(treeValues)
	if err != nil {
		log.Fatal("unable to start api service: ", err)
	}

	apiService.Start()
}
