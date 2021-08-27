package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"unicode/utf8"
)

var line = flag.Bool("l", false, "for counting line")
var character = flag.Bool("m", false, "for counting character")
var word = flag.Bool("w", false, "for counting word")

func countPrint(fileName string, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.Open(fileName)
	if err != nil {
		log.Println(err)
		return
	}
	scanner := bufio.NewScanner(file)
	if *line {
		countLine := 0
		for scanner.Scan() {
			countLine++
		}
		if scanner.Err() != nil {
			log.Println(err)
			return
		}
		fmt.Printf("%d\t%s\n", countLine, fileName)
	} else if *word {
		scanner.Split(bufio.ScanWords)
		countWord := 0
		for scanner.Scan() {
			countWord++
		}
		if scanner.Err() != nil {
			log.Println(err)
			return
		}
		fmt.Printf("%d\t%s\n", countWord, fileName)
	} else if *character {
		countCharacter := 0
		for scanner.Scan() {
			countCharacter += utf8.RuneCountInString(scanner.Text())
			countCharacter++
		}
		if scanner.Err() != nil {
			log.Println(err)
			return
		}
		fmt.Printf("%d\t%s\n", countCharacter, fileName)
	}
}

func main() {
	flag.Parse()
	if flag.NFlag() > 1 {
		log.Fatalf("usage %s with only one of the flags\n", os.Args[0])
	}
	if flag.NFlag() == 0 {
		*word = true
	}
	wg := sync.WaitGroup{}
	for _, s := range flag.Args() {
		wg.Add(1)
		go countPrint(s, &wg)
	}
	wg.Wait()
}
