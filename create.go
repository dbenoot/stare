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

// FUNCTIONS
//
//- create a page and add metadata to new page:
//    - date and time of creation
//    - availability in the menu
//    - order in the menu
//    - name in menu
//    - draft
//- create gallery folder

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var now = time.Now().Format(time.RFC1123)

func createPage(path string, pagename string, extra string) {

	file := filepath.Join(path, pagename+".html")

	if exists(file) {
		fmt.Println("Page already exists.")
	} else {

		f, err := os.Create(file)
		check(err)
		defer f.Close()

		content := "------------------------------------------------------------------------\ncreated on      : " + now + "\npresent in menu : n\nmenu order      : 0\nmenu name       : " + pagename + "\nstatus          : draft\n------------------------------------------------------------------------" + extra

		_, err = f.WriteString(content)
	}
}

func createGallery(galleryname string) {
	fmt.Println("Creating gallery " + galleryname)
	os.MkdirAll(filepath.Join("bodies", "galleries", galleryname), os.ModePerm)
}
