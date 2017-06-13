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

// TODO

// rewrite remove_header so it cuts on the second line consisting of ------

package main

import (
	// 	"flag"
	"bytes"
	"fmt"
	"github.com/russross/blackfriday"
	// "io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"
	// "time"
)

type Page struct {
	filename     string
	filetype     string
	time         string
	menu_present bool
	menu_order   int
	menu_name    string
	posted       bool
	content      string
	path         string
	navbar       string
	output       string
	rel_path     string
	base_path    string
	index        string
}

type Nav struct {
	path      string
	name      string
	orig_key  int
	base_path string
	filename  string
}

func render_site() {

	bodies := mapBodies("bodies")
	pages := mapPages(bodies)
	pages = createNavbar(pages)
	pages = createOutput(pages)
	writeOutput(pages)

	copySrc()

}

func mapBodies(path string) map[string]string {

	bodies := make(map[string]string)
	formats := []string{"html", "HTML", "md", "MD"}

	files := []string{}
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	for _, file := range files {
		if stringInSlice(strings.Split(file, ".")[len(strings.Split(file, "."))-1], formats) == true {
			content, _ := ioutil.ReadFile(file)
			bodies[file] = string(content)
		}
	}
	if err != nil {
		fmt.Println(err)
	}
	return bodies
}

func mapPages(bodies map[string]string) map[int]Page {

	var content_temp string
	c := make(map[int]Page)
	i := 0
	for key, value := range bodies {

		t := c[i]

		t.menu_present, t.menu_order, t.menu_name, t.posted, t.time, content_temp = parsePage(value)

		t.filetype = strings.ToLower(filepath.Ext(key))

		t.path, _ = filepath.Rel("", key)
		t.path = strings.Replace(t.path, "bodies"+string(filepath.Separator), "", 1)

		// define the relative path

		if strings.Contains(filepath.Dir(t.path), "pages") {
			t.rel_path = filepath.Join("..")
		} else {
			t.rel_path = filepath.Join(".")
		}

		// define the base_path

		t.base_path, t.filename = filepath.Split(t.path)

		// define the location of the index relative to the page

		t.index = filepath.Join(t.rel_path, "index.html")

		// Render md to html

		if t.filetype == ".md" {

			// change path to .html

			t.path = t.path[0:len(t.path)-len(t.filetype)] + ".html"

			// change output filename to .html

			t.filename = t.filename[0:len(t.filename)-len(t.filetype)] + ".html"

			// Render markdown to html

			content_temp2 := blackfriday.MarkdownCommon([]byte(content_temp))
			content_temp = string(content_temp2)
		}

		content, err := template.New("body").Parse(content_temp)
		check(err)
		w := bytes.NewBufferString("")
		content.Execute(w, map[string]string{"Css": filepath.Join(t.rel_path, "css") + string(filepath.Separator), "Js": filepath.Join(t.rel_path, "js") + string(filepath.Separator), "Index": t.index, "Img": filepath.Join(t.rel_path, "img") + string(filepath.Separator), "Page": filepath.Join(t.rel_path, "pages") + string(filepath.Separator)})

		t.content = w.String()

		fmt.Println(t.path, t.filename, t.filetype)

		c[i] = t
		i++
	}
	return c
}

func createNavbar(pages map[int]Page) map[int]Page {

	c := make(map[int]Nav)

	// get pages which should be present in the navbar

	for key, value := range pages {

		nav := c[value.menu_order]

		if value.menu_present {
			nav.name = value.menu_name
			nav.path = value.path
			nav.orig_key = key
			nav.base_path = value.base_path
			nav.filename = value.filename
		}

		c[value.menu_order] = nav
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
		j := 0
		for _, _ = range c {

			if pages[i].path == c[keys[j]].path {
				navact = "class=\"active\""
			} else {
				navact = ""
			}

			t.Execute(w, map[string]string{"Navactive": navact, "Navlink": filepath.Join(c[keys[j]].base_path, pages[i].rel_path, c[keys[j]].filename), "Navitem": c[keys[j]].name})

			j++
		}

		n.Execute(u, map[string]string{"Navlist": w.String(), "Index": pages[i].index})

		var tmp = pages[i]
		tmp.navbar = u.String()
		pages[i] = tmp

	}

	return pages
}

func createOutput(pages map[int]Page) map[int]Page {

	head, _ := template.ParseFiles("templates/header_template.html")
	foot, _ := template.ParseFiles("templates/footer_template.html")
	t, _ := template.ParseFiles("templates/page_template.html")

	i := 0
	for key, value := range pages {
		header := bytes.NewBufferString("")
		footer := bytes.NewBufferString("")
		w := bytes.NewBufferString("")

		head.Execute(header, map[string]string{"Css": filepath.Join(value.rel_path, "css") + string(filepath.Separator), "Js": filepath.Join(value.rel_path, "js") + string(filepath.Separator), "Index": value.index, "Img": filepath.Join(value.rel_path, "img") + string(filepath.Separator), "Page": filepath.Join(value.rel_path, "pages") + string(filepath.Separator)})
		foot.Execute(footer, map[string]string{"Css": filepath.Join(value.rel_path, "css") + string(filepath.Separator), "Js": filepath.Join(value.rel_path, "js") + string(filepath.Separator), "Index": value.index, "Img": filepath.Join(value.rel_path, "img") + string(filepath.Separator), "Page": filepath.Join(value.rel_path, "pages") + string(filepath.Separator)})

		t.Execute(w, map[string]string{"Header": header.String(), "Navbar": value.navbar, "Body": value.content, "Footer": footer.String(), "Css": filepath.Join(value.rel_path, "css") + string(filepath.Separator), "Js": filepath.Join(value.rel_path, "js") + string(filepath.Separator), "Index": value.index, "Img": filepath.Join(value.rel_path, "img") + string(filepath.Separator), "Page": filepath.Join(value.rel_path, "pages") + string(filepath.Separator)})

		var tmp = pages[key]
		tmp.output = w.String()
		pages[key] = tmp

		i++
	}

	return pages
}

func writeOutput(pages map[int]Page) {

	os.Mkdir(filepath.Join(".", "rendered"), os.ModePerm)

	for i := 0; i < len(pages); i++ {

		newPath := filepath.Join(".", "rendered", pages[i].path)

		err := os.MkdirAll(filepath.Dir(newPath), os.ModePerm)
		check(err)
		f, _ := os.Create(newPath)

		f.WriteString(pages[i].output)
	}
}

func parsePage(input string) (bool, int, string, bool, string, string) {

	var menu_present, posted bool
	var menu_order int
	var menu_name, time, content string

	lines := strings.Split(string(input), "\n")

	for j := 1; j < 6; j++ {
		if strings.Contains(lines[j], "posted") == true {
			posted = true
		}
		if strings.Contains(lines[j], "present in menu") == true && strings.TrimSpace(strings.Split(lines[j], ":")[1]) == "y" {
			menu_present = true
		}
		if strings.Contains(lines[j], "menu order") == true {
			menu_order, _ = strconv.Atoi(strings.TrimSpace(strings.Split(lines[j], ":")[1]))
		}
		if strings.Contains(lines[j], "menu name") == true {
			menu_name = strings.TrimSpace(strings.Split(lines[j], ":")[1])
		}
		if strings.Contains(lines[j], "created on") == true {
			time = strings.TrimSpace(strings.Split(lines[j], ":")[1])
		}
	}

	for i := 7; i < len(lines); i++ {
		content = content + lines[i] + "\n"
	}

	return menu_present, menu_order, menu_name, posted, time, content

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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
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
