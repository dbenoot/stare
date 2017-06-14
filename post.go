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

// Post draft pages in production
// Unpost production pages to draft

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var itemId int

func post(name []string, path string) {
	var items []string
	bodies := mapBodies(path)
	pages := mapPages(bodies)
	for _, vn := range name {
		for _, vp := range pages {
			if strings.Contains(vp.path, vn) {
				items = append(items, vp.body_path)
			}
		}

		item := findItem(items)

		if len(item) > 0 {
			fmt.Println("Posting", item)
			replaceInHeader(item, "status          : ", "status          : posted")
		}

	}

}

func unpost(name []string, path string) {
	var items []string
	bodies := mapBodies(path)
	pages := mapPages(bodies)
	for _, vn := range name {
		for _, vp := range pages {
			if strings.Contains(vp.path, vn) {
				items = append(items, vp.body_path)
			}
		}
	}

	item := findItem(items)
	fmt.Println(items, ":", item)

	if len(item) > 0 {
		filename := strings.Split(item, "/")[len(strings.Split(item, "/"))-1]
		fmt.Println("Posting", filename)
		replaceInHeader(item, "status          : posted", "status          : draft")
	}
}

func getFiles(path string) []string {
	formats := []string{".html", ".HTML", ".md", ".MD"}
	allfiles := []string{}
	files := []string{}
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		allfiles = append(allfiles, path)
		return nil
	})
	check(err)

	for _, file := range allfiles {
		if stringInSlice(filepath.Ext(file), formats) == true {
			files = append(files, file)
		}
	}

	return files
}

func createContainsArray(n []string, f []string) []string {
	var of []string
	for _, a := range n {
		for _, b := range f {
			if b == a {
				of = append(of, b)
			}
		}
	}
	return of

}
