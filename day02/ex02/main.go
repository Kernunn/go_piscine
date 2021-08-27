package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) <= 1 {
		os.Exit(1)
	}

	var cmd *exec.Cmd
	if len(os.Args) >= 3 {
		cmd = exec.Command(os.Args[1], os.Args[2:]...)
	} else {
		cmd = exec.Command(os.Args[1])
	}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cmd.Args = append(cmd.Args, scanner.Text())
	}
	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}

	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
