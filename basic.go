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

// Basic functions

package main

import (
	// "bytes"
	"fmt"
	"github.com/russross/blackfriday"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	// "text/template"
)

func mapBodies(path string) map[string]string {

	bodies := make(map[string]string)
	formats := []string{".html", ".HTML", ".md", ".MD"}

	files := []string{}
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	for _, file := range files {
		if stringInSlice(filepath.Ext(file), formats) == true {
			content, _ := ioutil.ReadFile(file)
			bodies[file] = string(content)
		}
	}

	check(err)
	return bodies
}

func mapPages(bodies map[string]string) map[int]Page {

	var content_temp string
	c := make(map[int]Page)
	i := 0
	for key, value := range bodies {
		if checkValidPage(value) == true {
			t := c[i]

			t.menu_present, t.menu_order, t.menu_name, t.posted, t.time, t.custom_header, content_temp = parsePage(value)

			t.filetype = strings.ToLower(filepath.Ext(key))

			t.body_path, _ = filepath.Rel("", key)

			// define path

			t.path = strings.Replace(t.body_path, "bodies"+string(filepath.Separator), "", 1)

			// define the relative path

			if strings.Contains(filepath.Dir(t.path), "pages") {
				t.rel_path = filepath.Join("..")
			} else if strings.Contains(filepath.Dir(t.path), "galleries") {
				t.rel_path = filepath.Join("..", "..")
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

			t.content = content_temp

			c[i] = t
			i++
		}
	}
	return c
}

func checkValidPage(input string) bool {

	var header string

	lines := strings.Split(string(input), "\n")

	if len(lines) >= 8 {

		for j := 1; j < 8; j++ {
			header = header + lines[j]
		}

		if strings.Contains(header, "status") == true && strings.Contains(header, "present in menu") == true && strings.Contains(header, "menu order") == true && strings.Contains(header, "menu name") == true && strings.Contains(header, "created on") == true {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func parsePage(input string) (bool, int, string, bool, string, string, string) {

	var menu_present, posted bool
	var menu_order int
	var menu_name, time, content, custom_header string

	lines := strings.Split(string(input), "\n")

	for j := 1; j < 8; j++ {
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
		if strings.Contains(lines[j], "custom template") == true {
			custom_header = strings.TrimSpace(strings.Split(lines[j], ":")[1])
		}
		if strings.Contains(lines[j], "created on") == true {
			time = strings.TrimSpace(strings.Split(lines[j], ":")[1])
		}
	}

	for i := 9; i < len(lines); i++ {
		content = content + lines[i] + "\n"
	}

	return menu_present, menu_order, menu_name, posted, time, custom_header, content

}

func move(inputname, outputname string) {
	err := os.Rename(inputname, outputname)

	if err != nil {
		fmt.Println("Page or gallery could not be moved. Please check standard folder structure.")
		return
	}
}

func copydir(source string, dest string) (err error) {

	// get properties of source dir
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// create dest dir

	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		sourcefilepointer := source + "/" + obj.Name()

		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			// create sub-directories - recursively
			err = copydir(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// perform copy
			err = copyfile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
	return
}

func movedir(source string, dest string) (err error) {

	// get properties of source dir
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// create dest dir

	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		sourcefilepointer := source + "/" + obj.Name()

		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			// create sub-directories - recursively
			err = copydir(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// perform copy
			err = copyfile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		}

	}

	// Remove the temporary files

	os.RemoveAll(source)

	return
}

func copyfile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}

	}

	return
}

func substitute(file, tie, replacetext string) {
	input, err := ioutil.ReadFile(file)
	check(err)

	lines := strings.Split(string(input), "\n")

	for line := range lines {
		lines[line] = strings.Replace(lines[line], tie, replacetext, -1)
	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(file, []byte(output), 0644)
	check(err)
}

func replaceInHeader(f, o, n string) {
	input, err := ioutil.ReadFile(f)
	check(err)

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, o) {
			lines[i] = n
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(f, []byte(output), 0644)
	check(err)
}

func substitute_in_header(file, o, n string) {
	input, err := ioutil.ReadFile(file)
	check(err)

	lines := strings.Split(string(input), "\n")

	for line := 0; line < 6; line++ {
		lines[line] = strings.Replace(lines[line], o, n, -1)
	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(file, []byte(output), 0644)
	check(err)
}

func findItem(items []string) (item string) {

	var itemId int

	if len(items) == 1 {
		itemId = 0
	} else {

		for i := 0; i < len(items); i++ {
			fmt.Println(strconv.Itoa(i) + " - " + items[i])
		}
		fmt.Println("Select the correct item:")
		if _, err := fmt.Scanf("%d", &itemId); err != nil {
			fmt.Printf("%s\n", err)
		}
	}

	if itemId >= len(items) {
		fmt.Println("Item does not exist.")
		return
	} else {
		return items[itemId]
	}

	return
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
