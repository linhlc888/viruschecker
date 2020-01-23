package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
)

// These variables are initialized externally during the build. See the Makefile.
var GitCommit string
var GitLastTag string
var GitExactTag string

func printVersion() {
	if GitExactTag == "undefined" {
		GitExactTag = ""
	}

	version := GitLastTag

	if GitLastTag == "" {
		version = fmt.Sprintf("%s-dev-%.10s", version, GitCommit)
	}

	if GitCommit == "" {
		fmt.Println("Viruschecker version: unknown (not compiled with the makefile)")
	} else {
		fmt.Printf("Viruschecker version: %s\n", version)
	}
	cmd := exec.Command("clamscan", "--version")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	fmt.Printf("System version: %s/%s\n", runtime.GOARCH, runtime.GOOS)
	fmt.Printf("Golang version: %s\n", runtime.Version())
	fmt.Print(out.String())
}
