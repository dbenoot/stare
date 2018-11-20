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

// TODO
//
// add responsive image breakpoints

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// define variables

var site = Site{
	pagedir:     "pages",
	blogdir:     "blogs",
	srcdir:      "src",
	gallerydir:  "gallery",
	templatedir: "templates",
}

// define Page, Nav, Site

type Page struct {
	filename      string
	filetype      string
	time          string
	menu_present  bool
	menu_order    int
	menu_name     string
	custom_header string
	posted        bool
	content       string
	body_path     string
	path          string
	navbar        string
	gallery       string
	output        string
	rel_path      string
	base_path     string
	index         string
}

type Nav struct {
	path      string
	name      string
	orig_key  int
	base_path string
	filename  string
}

type Gallery struct {
	link  string
	thumb string
	name  string
}

type Site struct {
	pagedir, blogdir, srcdir, gallerydir, templatedir string
}

func main() {

	createCommand := flag.NewFlagSet("create", flag.ExitOnError)
	pageCreateFlag := createCommand.String("page", "", "Name of the new page.")
	galleryCreateFlag := createCommand.String("gallery", "", "Name of the new gallery.")

	if len(os.Args) == 1 {
		fmt.Println("usage: stare <command> [<args>]")
		fmt.Println("The most commonly used stare commands are: ")
		fmt.Println(" init			Initialize a stare website.")
		fmt.Println(" render			Renders the website.")
		fmt.Println(" watch			Watches for changes and renders when a change happens.")
		fmt.Println(" gallery		Preprocesses the galleries. Run before render command.")
		fmt.Println(" create")
		fmt.Println("   -page		Creates a new page.")
		fmt.Println("   -gallery		Create a new gallery.")
		fmt.Println(" post			Posts a document.")
		fmt.Println(" unpost			Unposts a document.")
		fmt.Println(" list			Lists all documents.")
		fmt.Println(" archive		Archives a document.")
		fmt.Println(" unarchive		Unarchives a document.")
		fmt.Println(" snap          Create a snapshot of the current rendered version.")
		fmt.Println(" tag			Add a tag to a document. Syntax: stare tag [document] [tag].")
		return
	}

	switch os.Args[1] {
	case "init":
		init_site()
	case "render":
		fmt.Println("Rendering...")
		startTime := time.Now()
		render_site()
		endTime := time.Now()
		fmt.Println("Elapsed time:", endTime.Sub(startTime))
	case "gallery":
		fmt.Println("Creating the galleries...")
		startTime := time.Now()
		renderGalleries()
		endTime := time.Now()
		fmt.Println("Elapsed time:", endTime.Sub(startTime))
	case "create":
		createCommand.Parse(os.Args[2:])
	case "archive":
		archive(os.Args[2:])
	case "unarchive":
		unarchive(os.Args[2:])
	case "post":
		post(os.Args[2:])
	case "unpost":
		unpost(os.Args[2:])
	case "list":
		sourcelist()
	case "watch":
		watch()
	case "snap":
		snap()
	case "tag":
		tag(os.Args[2:])
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}

	if createCommand.Parsed() {
		if *pageCreateFlag == "" && *galleryCreateFlag == "" {
			fmt.Println("Please provide the page or gallery name using -page or -gallery parameter.")
		} else if *pageCreateFlag != "" {
			createPage(filepath.Join("bodies", "pages"), *pageCreateFlag, "")
		} else if *galleryCreateFlag != "" {
			createGallery(*galleryCreateFlag)
		}
	}

}
