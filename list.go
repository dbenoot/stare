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
)

func sourcelist() {

	bodies := mapBodies("bodies")
	pages := mapPages(bodies)

	fmt.Println("BODIES")
	fmt.Println("   Name in menu\t  Posted?\t  File path")
	for _, value := range pages {
		fmt.Println("- ", value.menu_name, " \t ", value.posted, " \t ", value.path)
	}

	fmt.Println("\nGALLERIES")
	listdir("bodies/galleries")

	return
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
