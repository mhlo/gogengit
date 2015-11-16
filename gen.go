// Copyright (c) 2015 Michael Boe. All rights reserved. Use of this source code
// is governed by the LICENSE file.

// Generate a file (by default, genver.go) that contains versioning info in it.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	VFILE   = flag.String("version-file", "VERSION", "name of file containing a version-string")
	IGFILE  = flag.String("ignore-file", "NOVERSION", "if named file exists, ignore the version-file, even if arg is present")
	OFILE   = flag.String("production", "genver.go", "name of produced file contain version info")
	PACKAGE = flag.String("package", "main", "name of package in produce file")
)

var tVersionGo = `package %s

var ver%s = "%s"
`

func main() {
	flag.Parse()
	// does VFILE not exist or IGFILE exist? If so, use that for versioning
	_, igErr := os.Stat(*IGFILE)
	versionInfo := ""
	useVersionFile := true

	if igErr != nil {
		if !os.IsNotExist(igErr) {
			fmt.Fprintf(os.Stderr,
				"ignore-file %s has a problem; file-information not available. Error=%v\n",
				*IGFILE, igErr,
			)
		}
	} else {
		useVersionFile = false
	}

	if useVersionFile {
		if _, verErr := os.Stat(*VFILE); verErr != nil {
			useVersionFile = false
		}
	}

	if useVersionFile {
		versionBytes, readErr := ioutil.ReadFile(*VFILE)
		if readErr != nil {
			log.Fatal("version-file", *VFILE, "exists but bad read:", readErr)
		}
		versionInfo = strings.TrimRight(string(versionBytes), "\n\r")
	} else {
		// if no version-file, use git
		// git log -n 1 --format="format: +%h %cd" HEAD
		var out bytes.Buffer
		cmd := exec.Command("git", "log", "-n", "1", "--format=format: +%h %cd", "HEAD")
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Fatal("cannot get version info from git. Command exited: " + err.Error())
		}
		versionInfo = strings.TrimRight(out.String(), "\n\r")
	}

	versionGo := fmt.Sprintf(tVersionGo, *PACKAGE, *PACKAGE, versionInfo)
	os.Remove(*OFILE)
	ofErr := ioutil.WriteFile(*OFILE, []byte(versionGo), 0444)
	if ofErr != nil {
		log.Fatal("cannot open file holding version production ", ofErr.Error)
	}
}
