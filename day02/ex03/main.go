package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

var outputDir = flag.String("a", "", "output directory")

func main() {
	flag.Parse()
	wg := sync.WaitGroup{}
	for _, s := range flag.Args() {
		wg.Add(1)
		go tarArchiveAndCompress(s, &wg)
	}
	wg.Wait()
}

func tarArchiveAndCompress(fileForArchiving string, wg *sync.WaitGroup) {
	defer wg.Done()

	input, err := os.Open(fileForArchiving)
	if err != nil {
		log.Println(err)
		return
	}
	defer input.Close()

	stat, err := input.Stat()
	if err != nil {
		log.Println(err)
		return
	}
	timestamp := fmt.Sprintf("%d", stat.ModTime().Unix())
	var output *os.File
	if *outputDir != "" {
		output, err = os.Create(*outputDir + "/" + fileForArchiving + "_" + timestamp + ".tar.gz")
	} else {
		output, err = os.Create(fileForArchiving + "_" + timestamp + ".tar.gz")
	}
	if err != nil {
		log.Println(err)
		return
	}
	defer output.Close()

	gw := gzip.NewWriter(output)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	header, err := tar.FileInfoHeader(stat, stat.Name())
	if err != nil {
		log.Println(err)
		return
	}

	header.Name = fileForArchiving

	err = tw.WriteHeader(header)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = io.Copy(tw, input)
	if err != nil {
		log.Println(err)
		return
	}
}
