package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

var file = flag.Bool("f", false, "find files")
var ext = flag.String("ext", "", "find files by extension")
var dir = flag.Bool("d", false, "find directories")
var symlinks = flag.Bool("sl", false, "find symlinks")
var all = false

func main() {
	flag.Parse()
	if !*file && !*dir && !*symlinks {
		all = true
	}
	err := filepath.Walk(flag.Arg(0), func(path string, info fs.FileInfo, err error) error {
		if info == nil || path == flag.Arg(0) {
			return err
		}
		e, ok := err.(*os.PathError)
		if ok && "permission denied" == e.Err.Error() {
			err = nil
			return err
		}
		if info.IsDir() && (*dir || all) {
			fmt.Println(path)
		}
		if info.Mode().IsRegular() && (*file || all) {
			if *file && *ext != "" {
				if filepath.Ext(path) == "." + *ext {
					fmt.Println(path)
				}
 			} else {
				fmt.Println(path)
			}
		}
		if info.Mode().Type() == fs.ModeSymlink && (*symlinks || all) {
			fmt.Printf("%s -> ", path)
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			if _, err = os.Stat(link); err != nil {
				fmt.Println("[broken]")
			} else {
				fmt.Println(link)
			}
		}
		return err
	})
	if err != nil {
		log.Fatal(err)
	}
}
