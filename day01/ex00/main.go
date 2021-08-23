package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
)

type cake struct {
	Name     string      `xml:"name" json:"name"`
	Time     string      `xml:"stovetime" json:"time"`
	Products ingredients `xml:"ingredients" json:"-"`
	Items    []item      `xml:"-" json:"ingredients"`
}

type ingredients struct {
	Items []item `xml:"item"`
}

type item struct {
	Name  string `xml:"itemname" json:"ingredient_name"`
	Count string `xml:"itemcount" json:"ingredient_count"`
	Unit  string `xml:"itemunit" json:"ingredient_unit,omitempty"`
}

type recipes struct {
	Cake []cake `xml:"cake" json:"cake"`
}

type DBReader interface {
	Read([]byte) (*recipes, error)
}

type JSONReader struct{}

func (r JSONReader) Read(input []byte) (*recipes, error) {
	recipes := &recipes{}
	err := json.Unmarshal(input, recipes)
	if err != nil {
		return nil, fmt.Errorf("json unmarshalling failure: %v", err)
	}
	for i := range recipes.Cake {
		recipes.Cake[i].Products.Items = recipes.Cake[i].Items
	}
	return recipes, nil
}

type XMLReader struct{}

func (r XMLReader) Read(input []byte) (*recipes, error) {
	recipes := &recipes{}
	err := xml.Unmarshal(input, recipes)
	if err != nil {
		return nil, fmt.Errorf("xml unmarshalling failure: %v", err)
	}
	for i := range recipes.Cake {
		recipes.Cake[i].Items = recipes.Cake[i].Products.Items
	}
	return recipes, nil
}

var nameFile = flag.String("f", "", "file name")

func main() {
	flag.Parse()
	if *nameFile == "" {
		fmt.Printf("usage %s -f file name\n", os.Args[0])
		os.Exit(1)
	}

	var reader DBReader
	switch path.Ext(*nameFile) {
	case ".xml":
		reader = XMLReader{}
	case ".json":
		reader = JSONReader{}
	default:
		log.Fatal("unsupported file format")
	}

	file, err := os.ReadFile(*nameFile)
	if err != nil {
		log.Fatal(err)
	}
	recipes, err := reader.Read(file)
	if err != nil {
		log.Fatal(err)
	}

	var b []byte
	switch path.Ext(*nameFile) {
	case ".xml":
		b, err = json.MarshalIndent(recipes, "", "    ")
	case ".json":
		b, err = xml.MarshalIndent(recipes, "", "    ")
	default:
	}
	if err != nil {
		log.Fatalf("marshaling failure: %v\n", err)
	}
	fmt.Println(string(b))
}
