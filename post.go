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
	"strings"
)

var itemId int

func post(name []string) {

	bodies := mapBodies("bodies")
	pages := mapPages(bodies)
	for _, vn := range name {
		var items []string
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

func unpost(name []string) {

	bodies := mapBodies("bodies")
	pages := mapPages(bodies)
	for _, vn := range name {
		var items []string
		for _, vp := range pages {
			if strings.Contains(vp.path, vn) {
				items = append(items, vp.body_path)
			}
		}
		item := findItem(items)

		if len(item) > 0 {
			fmt.Println("Unposting", item)
			replaceInHeader(item, "status          : posted", "status          : ")
		}
	}
}
