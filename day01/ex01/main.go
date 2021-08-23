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

var oldDB = flag.String("old", "", "the name of the old database")
var newDB = flag.String("new", "", "the name of the new database")

func main() {
	flag.Parse()
	if *oldDB == "" || *newDB == "" {
		fmt.Printf("usage %s --old file name --new file name\n", os.Args[0])
		os.Exit(1)
	}

	var readerOldDB DBReader
	var readerNewDB DBReader
	switch path.Ext(*oldDB) {
	case ".xml":
		readerOldDB = XMLReader{}
	case ".json":
		readerOldDB = JSONReader{}
	default:
		log.Fatal("unsupported file format")
	}
	switch path.Ext(*newDB) {
	case ".xml":
		readerNewDB = XMLReader{}
	case ".json":
		readerNewDB = JSONReader{}
	default:
		log.Fatal("unsupported file format")
	}

	oldDB, err := os.ReadFile(*oldDB)
	if err != nil {
		log.Fatal(err)
	}
	newDB, err := os.ReadFile(*newDB)
	if err != nil {
		log.Fatal(err)
	}

	recipes1, err := readerOldDB.Read(oldDB)
	if err != nil {
		log.Fatal(err)
	}
	recipes2, err := readerNewDB.Read(newDB)
	if err != nil {
		log.Fatal(err)
	}

	printDiff(recipes1, recipes2)
}

func printDiff(old, new *recipes) {
	oldCakes := make(map[string]cake)
	newCakes := make(map[string]cake)
	sameCakes := make([]string, 10)
	for _, c := range old.Cake {
		oldCakes[c.Name] = c
	}
	for _, c := range new.Cake {
		newCakes[c.Name] = c
	}
	for s := range newCakes {
		c, ok := oldCakes[s]
		if !ok {
			fmt.Printf("ADDED cake \"%s\"\n", s)
		} else {
			sameCakes = append(sameCakes, c.Name)
		}
	}
	for s := range oldCakes {
		if _, ok := newCakes[s]; !ok {
			fmt.Printf("REMOVED cake \"%s\"\n", s)
		}
	}

	for _, nameCake := range sameCakes {
		if oldCakes[nameCake].Time != newCakes[nameCake].Time {
			fmt.Printf("CHANGED cooking time for cake \"%s\" - \"%s\" instead of \"%s\"\n",
				nameCake, newCakes[nameCake].Time, oldCakes[nameCake].Time)
		}
		oldIngredients := make(map[string]item)
		newIngredients := make(map[string]item)
		sameIngredients := make([]string, 10)
		for _, item := range oldCakes[nameCake].Items {
			oldIngredients[item.Name] = item
		}
		for _, item := range newCakes[nameCake].Items {
			newIngredients[item.Name] = item
		}
		for nameItem := range newIngredients {
			item, ok := oldIngredients[nameItem]
			if !ok {
				fmt.Printf("ADDED ingredient \"%s\" for cake \"%s\"\n", nameItem, nameCake)
			} else {
				sameIngredients = append(sameIngredients, item.Name)
			}
		}
		for nameItem := range oldIngredients {
			if _, ok := newIngredients[nameItem]; !ok {
				fmt.Printf("REMOVED ingredient \"%s\" for cake \"%s\"\n", nameItem, nameCake)
			}
		}

		for _, nameItem := range sameIngredients {
			oldIngredient := oldIngredients[nameItem]
			newIngredient := newIngredients[nameItem]
			if oldIngredient.Unit != "" || newIngredient.Unit != "" {
				if oldIngredient.Unit != "" && newIngredient.Unit != "" {
					if oldIngredient.Unit != newIngredient.Unit {
						fmt.Printf("CHANGED unit for ingredient \"%s\" for cake  \"%s\" - " +
							"\"%s\" instead of \"%s\"\n", nameItem, nameCake, newIngredient.Unit,
							oldIngredient.Unit)
					}
				} else if oldIngredient.Unit != "" {
					fmt.Printf("ADDED unit \"%s\" for ingredient \"%s\" for cake  \"%s\"\n",
						oldIngredient.Unit, nameItem, nameCake)
				} else {
					fmt.Printf("REMOVED unit \"%s\" for ingredient \"%s\" for cake  \"%s\"\n",
						oldIngredient.Unit, nameItem, nameCake)
				}
			}
			if oldIngredient.Count != newIngredient.Count {
				fmt.Printf("CHANGED unit count for ingredient \"%s\" for cake  \"%s\" - " +
					"\"%s\" instead of \"%s\"\n", nameItem, nameCake, newIngredient.Count,
					oldIngredient.Count)
			}
		}
	}
}
