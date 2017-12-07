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

// Render all production pages and galleries

package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

func render_site() {

	RemoveContentsLeaveGit("rendered")

	bodies := mapBodies("bodies")
	pages := mapPages(bodies)
	pages = qPosted(pages)
	pages = createNavbar(pages)
	pages = createGalleryBody(pages)
	pages = createOutput(pages)
	writeOutput(pages)

	copySrc()
	copyGalleries()

}

func createNavbar(pages map[int]Page) map[int]Page {

	c := make(map[int]Nav)

	// get pages which should be present in the navbar and their values

	for key, value := range pages {

		nav := c[value.menu_order]
		if value.menu_present {
			nav.name = value.menu_name
			nav.path = value.path
			nav.orig_key = key
			nav.base_path = value.base_path
			nav.filename = value.filename

			// Check that the menu order is unique and if so, write to map

			if _, ok := c[value.menu_order]; ok {
				fmt.Println("Menu order in the posted menu page -", value.menu_name, "- is not unique (same as page -", c[value.menu_order].name, "-). Please correct.")
				os.Exit(2)
			} else {
				c[value.menu_order] = nav
			}
		}
	}

	// make the keys consecutive

	keys := getKeys(c)
	sort.Ints(keys)

	// add info in navbar_item and record answers in response

	n, _ := template.ParseFiles("templates/navbar_template.html")
	t, _ := template.ParseFiles("templates/navbar_item.html")
	var navact string

	for i := 0; i < len(pages); i++ {
		u := bytes.NewBufferString("")
		w := bytes.NewBufferString("")

		// initialize j as 0, and don't use the range c val; otherwise the navbar will not be in the correct order

		j := 0

		// iterate over c

		for _, _ = range c {

			// set class to active when page is the same as nav

			if pages[i].path == c[keys[j]].path {
				navact = "class=\"active\""
			} else {
				navact = ""
			}
			t.Execute(w, map[string]string{"Navactive": navact, "Navlink": filepath.Join(pages[i].rel_path, c[keys[j]].base_path, c[keys[j]].filename), "Navitem": c[keys[j]].name})

			j++
		}

		n.Execute(u, map[string]string{"Navlist": w.String(), "Index": pages[i].index})

		var tmp = pages[i]
		tmp.navbar = u.String()
		pages[i] = tmp
	}
	return pages
}

func createGalleryBody(pages map[int]Page) map[int]Page {

	c := make(map[int]Gallery)

	// get pages which should be present in the gallery

	for key, value := range pages {

		if strings.Contains(value.body_path, filepath.Join("bodies", "galleries")) {
			gal := c[key]
			gal.link = filepath.Join(value.path)
			gal.thumb = filepath.Join(value.base_path, "gallery.jpg")
			gal.name = value.menu_name
			c[key] = gal
		}

	}

	for key, pageValue := range pages {

		u := bytes.NewBufferString("")
		for _, value := range c {

			g, _ := template.ParseFiles("templates/gallery_item.html")
			g.Execute(u, map[string]string{"Subgallerylink": filepath.Join(pageValue.rel_path, value.link), "Subgallerythumb": filepath.Join(pageValue.rel_path, value.thumb), "Subgalleryname": value.name})
		}

		var tmp = pages[key]
		tmp.gallery = u.String()
		pages[key] = tmp

	}
	return pages
}

func createOutput(pages map[int]Page) map[int]Page {

	head, _ := template.ParseFiles("templates/header_template.html")
	foot, _ := template.ParseFiles("templates/footer_template.html")

	t, _ := template.ParseFiles("templates/page_template.html")

	i := 0
	for key, value := range pages {
		body, _ := template.New("body").Parse(value.content)
		header := bytes.NewBufferString("")
		footer := bytes.NewBufferString("")
		parsedContent := bytes.NewBufferString("")

		w := bytes.NewBufferString("")
		head.Execute(header, map[string]string{"Css": filepath.Join(value.rel_path, "css") + string(filepath.Separator), "Js": filepath.Join(value.rel_path, "js") + string(filepath.Separator), "Index": value.index, "Img": filepath.Join(value.rel_path, "img") + string(filepath.Separator), "Page": filepath.Join(value.rel_path, "pages") + string(filepath.Separator)})
		foot.Execute(footer, map[string]string{"Css": filepath.Join(value.rel_path, "css") + string(filepath.Separator), "Js": filepath.Join(value.rel_path, "js") + string(filepath.Separator), "Index": value.index, "Img": filepath.Join(value.rel_path, "img") + string(filepath.Separator), "Page": filepath.Join(value.rel_path, "pages") + string(filepath.Separator)})

		body.Execute(parsedContent, map[string]string{"Header": header.String(), "Navbar": value.navbar, "Gallery": value.gallery, "Body": value.content, "Footer": footer.String(), "Css": filepath.Join(value.rel_path, "css") + string(filepath.Separator), "Js": filepath.Join(value.rel_path, "js") + string(filepath.Separator), "Index": value.index, "Img": filepath.Join(value.rel_path, "img") + string(filepath.Separator), "Page": filepath.Join(value.rel_path, "pages") + string(filepath.Separator)})

		t.Execute(w, map[string]string{"Header": header.String(), "Navbar": value.navbar, "Gallery": value.gallery, "Body": parsedContent.String(), "Footer": footer.String(), "Css": filepath.Join(value.rel_path, "css") + string(filepath.Separator), "Js": filepath.Join(value.rel_path, "js") + string(filepath.Separator), "Index": value.index, "Img": filepath.Join(value.rel_path, "img") + string(filepath.Separator), "Page": filepath.Join(value.rel_path, "pages") + string(filepath.Separator)})

		var tmp = pages[key]
		tmp.output = w.String()
		pages[key] = tmp

		i++
	}

	return pages
}

func qPosted(pages map[int]Page) map[int]Page {
	c := make(map[int]Page)
	for key, value := range pages {
		if value.posted {
			c[key] = value
		}
	}
	return c
}

func writeOutput(pages map[int]Page) {

	os.Mkdir(filepath.Join(".", "rendered"), os.ModePerm)

	for i := 0; i < len(pages); i++ {

		// var f *File

		newPath := filepath.Join(".", "rendered", pages[i].path)
		err := os.MkdirAll(filepath.Dir(newPath), os.ModePerm)
		check(err)

		_, err = os.Stat(newPath)
		if err != nil {

			f, err := os.Create(newPath)
			check(err)
			defer f.Close()

			f.WriteString(pages[i].output)
		}
	}
}

func getKeys(mymap map[int]Nav) []int {
	keys := make([]int, len(mymap))

	i := 0
	for k := range mymap {
		keys[i] = k
		i++
	}

	return keys
}

func RemoveContentsLeaveGit(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		if name != ".git" {
			err = os.RemoveAll(filepath.Join(dir, name))
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func copySrc() {

	srcItems, _ := filepath.Glob("src" + string(filepath.Separator) + "*")

	for i := 0; i < len(srcItems); i++ {
		file, err := os.Open(srcItems[i])
		check(err)
		defer file.Close()

		fi, err := file.Stat()
		check(err)

		if fi.IsDir() {
			copydir(srcItems[i], "rendered"+string(filepath.Separator)+filepath.Base(srcItems[i]))
		} else {
			copyfile(srcItems[i], "rendered"+string(filepath.Separator)+filepath.Base(srcItems[i]))
		}
	}
}

func copyGalleries() {

	srcItems, _ := filepath.Glob(filepath.Join("bodies", "galleries", "*"))

	for i := 0; i < len(srcItems); i++ {
		file, err := os.Open(srcItems[i])
		check(err)
		defer file.Close()

		fi, err := file.Stat()
		check(err)

		if fi.IsDir() {
			copydir(srcItems[i], filepath.Join("rendered", "galleries", filepath.Base(srcItems[i])))
		} else {
			copyfile(srcItems[i], filepath.Join("rendered", "galleries", filepath.Base(srcItems[i])))
		}
	}
}
