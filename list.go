//   Copyright 2016 The Stare Authors
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

// List of pages and galleries

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	// "log"
	"path/filepath"
	"strings"
)

var j int

func sourcelist() {
	fmt.Println("BODIES")
	j = 1
	list("bodies/*.html")

	fmt.Println("\nGALLERIES")
	listdir("bodies/galleries")

	return
}

func list(folder string) {
	filepath.Walk("bodies", checkStatus)
}

func listdir(folder string) {
	files, _ := ioutil.ReadDir(folder)

	i := 1
	for _, f := range files {
		fmt.Println("- ", f.Name())
		i++
	}
}

func visit(path string, f os.FileInfo, err error) error {
	fmt.Printf("Visited: %s\n", path)
	return nil
}

func checkStatus(file string, f os.FileInfo, err error) error {
	fi, err := os.Stat(file)
	check(err)

	if fi.IsDir() == false {

		input, err := ioutil.ReadFile(file)
		check(err)

		lines := strings.Split(string(input), "\n")

		for j := 1; j < 6; j++ {
			if strings.Contains(lines[j], "status          :") && strings.Contains(lines[j], "posted") == true {
				fmt.Println("- ", file, " \t posted")
			} else if strings.Contains(lines[j], "status          :") && strings.Contains(lines[j], "posted") == false {
				fmt.Println("- ", file, " \t draft")
			}
		}
	}

	return nil
}
