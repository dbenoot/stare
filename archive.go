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

// This file contains the functions for archiving and unarchiving web pages and galleries

// TODO
//
// - archiving and unarchiving galleries

package main

import (
	"fmt"
	// 	"io/ioutil"
	// 	"os"
	"path/filepath"
	"strings"
)

func archive(name []string) {

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
			fmt.Println("Archiving page", item)
			move(item, filepath.Join("archive", item))
		}
	}
}

func unarchive(name []string) {

	bodies := mapBodies(filepath.Join("archive", "bodies"))
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
			fmt.Println("Unarchiving page", item)
			move(item, strings.Replace(item, "archive"+string(filepath.Separator), "", 1))
		}
	}
}

// func archive_gallery(galleryname string) {

// 	var galleries []string

// 	path := filepath.Join("archive", "pages", "gallery")
// 	_, err := os.Stat(path)
// 	if err != nil {
// 		os.MkdirAll(path, 0755)
// 	}

// 	g, _ := ioutil.ReadDir(filepath.Join(site.pagedir, site.gallerydir))

// 	for i := 0; i < len(g); i++ {
// 		if g[i].IsDir() == true {
// 			galleries = append(galleries, g[i].Name())
// 		}
// 	}

// 	gallery := findItem(galleries)

// 	fmt.Println("Archiving gallery", gallery)
// 	movedir(filepath.Join("pages", "gallery", gallery), filepath.Join(path, gallery))

// 	files, _ := ioutil.ReadDir(site.gallerydir)

// 	if len(files) == 0 {
// 		cfg.Section("general").NewKey("gallery", "")
// 		cfg.SaveTo("config.ini")
// 	}
// }

// func unarchive_gallery(galleryname string) {

// 	var galleries []string

// 	path := filepath.Join("archive", "pages", "gallery")
// 	_, err := os.Stat(path)
// 	if err != nil {
// 		os.MkdirAll(path, 0755)
// 	}

// 	g, _ := ioutil.ReadDir(filepath.Join(path))

// 	for i := 0; i < len(g); i++ {
// 		if g[i].IsDir() == true {
// 			galleries = append(galleries, g[i].Name())
// 		}
// 	}

// 	gallery := findItem(galleries)

// 	fmt.Println("Archiving gallery", gallery)
// 	movedir(filepath.Join(path, gallery), filepath.Join("pages", "gallery", gallery))

// 	if site.gallery == false {
// 		cfg.Section("general").NewKey("gallery", "y")
// 		cfg.SaveTo("config.ini")
// 	}
// }
