package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var oldFS = flag.String("old", "", "the name of the old filesystem")
var newFS = flag.String("new", "", "the name of the new filesystem")

type treeFS struct {
	name  string
	child []*treeFS
	visit bool
}

func (t *treeFS) Add(filename string, visit bool) bool {
	if !filepath.IsAbs(filename) {
		log.Printf("%s is not absolute\n", filename)
		return false
	}
	add := false
	tokens := strings.Split(filename, "/")[1:]
	for _, token := range tokens {
		prev := t
		for _, node := range t.child {
			if node.name == token {
				t = node
				break
			}
		}
		if prev == t {
			newNode := &treeFS{name: token}
			t.child = append(t.child, newNode)
			t = newNode
			add = true;
			t.visit = visit
		} else {
			t.visit = true
		}
	}
	return add
}

func (t *treeFS) PrintNotVisited(prefix string) {
	if t.child == nil && t.visit == false {
		fmt.Printf("REMOVED %s\n", prefix + t.name)
	}
	for _, node := range t.child {
		node.PrintNotVisited(prefix + t.name + "/")
	}
}

func main() {
	flag.Parse()
	if *oldFS == "" || *newFS == "" {
		fmt.Printf("usage %s --old file name --new file name\n", os.Args[0])
		os.Exit(1)
	}

	file1, err := os.Open(*oldFS)
	if err != nil {
		log.Fatal(err)
	}
	file2, err := os.Open(*newFS)
	if err != nil {
		log.Fatal(err)
	}

	fs := treeFS{}

	scanner := bufio.NewScanner(file1)
	for scanner.Scan() {
		fileName := scanner.Text()
		fs.Add(fileName, false)
	}
	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}
	scanner = bufio.NewScanner(file2)
	for scanner.Scan() {
		fileName := scanner.Text()
		success := fs.Add(fileName, true)
		if success {
			fmt.Printf("ADDED %s\n", fileName)
		}
	}
	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fs.PrintNotVisited("")
}
