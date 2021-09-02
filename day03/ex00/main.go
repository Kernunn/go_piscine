package main

import (
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"strings"
)

func main() {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatal(err)
	}

	res, err := es.Index("places", strings.NewReader(`{
	"mappings": {
		"properties": {
			"name": {
				"type": "text"
			},
			"address": {
				"type": "text"
			},
			"phone": {
				"type": "text"
			},
			"location": {
				"type": "geo_point"
			}
		}
	}
}
`))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)
	defer res.Body.Close()
}
