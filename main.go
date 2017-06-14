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
	"github.com/go-ini/ini"
	"os"
	"path/filepath"
	"time"
)

// define variables

var cfg, _ = ini.LooseLoad("config.ini")

var site = Site{
	pagedir:     "pages",
	blogdir:     "blogs",
	srcdir:      "src",
	gallerydir:  "gallery",
	templatedir: "templates",
	multiLang:   cfg.Section("general").Key("multiple_language_support").MustBool(),
	primaryLang: cfg.Section("general").Key("primary_language").String(),
	languages:   cfg.Section("general").Key("languages").Strings(","),
	gallery:     cfg.Section("general").Key("gallery").MustBool(),
}

// define Page, Nav, Site

type Page struct {
	filename     string
	filetype     string
	time         string
	menu_present bool
	menu_order   int
	menu_name    string
	posted       bool
	content      string
	body_path    string
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

type Site struct {
	pagedir, blogdir, srcdir, gallerydir, templatedir, primaryLang string
	languages                                                      []string
	multiLang, gallery                                             bool
}

func main() {

	createCommand := flag.NewFlagSet("create", flag.ExitOnError)
	pageCreateFlag := createCommand.String("page", "", "Name of the new page.")
	galleryCreateFlag := createCommand.String("gallery", "", "Name of the new gallery.")

	archiveCommand := flag.NewFlagSet("archive", flag.ExitOnError)
	pageArchiveFlag := archiveCommand.String("page", "", "Name of the page to be archived.")
	gallerypageArchiveFlag := archiveCommand.String("gallery", "", "Name of the gallery to be archived.")

	unarchiveCommand := flag.NewFlagSet("unarchive", flag.ExitOnError)
	pageUnarchiveFlag := unarchiveCommand.String("page", "", "Name of the page to be unarchived.")
	gallerypageUnarchiveFlag := unarchiveCommand.String("gallery", "", "Name of the gallery to be unarchived.")

	// postCommand := flag.NewFlagSet("post", flag.ExitOnError)
	// pagePostFlag := postCommand.String("page", "", "Name of the page to be posted.")

	// unpostCommand := flag.NewFlagSet("unpost", flag.ExitOnError)
	// pageUnpostFlag := unpostCommand.String("page", "", "Name of the page to be unposted.")

	if len(os.Args) == 1 {
		fmt.Println("usage: stare <command> [<args>]")
		fmt.Println("The most commonly used stare commands are: \n")
		fmt.Println(" init          Initialize a stare website.\n")
		fmt.Println(" render        Renders the website.\n")
		fmt.Println(" create")
		fmt.Println("   -page       Creates a new page")
		fmt.Println("   -gallery    Create a new gallery")
		fmt.Println(" post")
		fmt.Println("   -page       Posts a page")
		fmt.Println(" unpost")
		fmt.Println("   -page       Unposts a page")
		fmt.Println(" archive")
		fmt.Println("   -page       Archives a page")
		fmt.Println("   -gallery    Archives a gallery")
		fmt.Println(" unarchive")
		fmt.Println("   -page       Unarchives a page")
		fmt.Println("   -gallery    Unarchives a gallery")
		return
	}

	switch os.Args[1] {
	case "init":
		init_site()
	case "render":
		fmt.Println("Rendering!")
		startTime := time.Now()
		render_site()
		endTime := time.Now()
		fmt.Println("Elapsed time:", endTime.Sub(startTime))
	case "create":
		createCommand.Parse(os.Args[2:])
	case "archive":
		archiveCommand.Parse(os.Args[2:])
	case "unarchive":
		unarchiveCommand.Parse(os.Args[2:])
	case "post":
		post(os.Args[2:], "bodies/")
		//postCommand.Parse(os.Args[2:])
	case "unpost":
		unpost(os.Args[2:], "bodies/")
	case "list":
		sourcelist()
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}

	if createCommand.Parsed() {
		if *pageCreateFlag == "" && *galleryCreateFlag == "" {
			fmt.Println("Please provide the page or gallery name using -page or -gallery parameter.")
		} else if *pageCreateFlag != "" {
			createPage(filepath.Join("bodies", "pages"), *pageCreateFlag)
		} else if *galleryCreateFlag != "" {
			createGallery(*galleryCreateFlag)
		}
	}

	if archiveCommand.Parsed() {
		if *pageArchiveFlag == "" && *gallerypageArchiveFlag == "" {
			fmt.Println("Please provide the page or blog name using -page or -gallery option.")
		} else if *pageArchiveFlag != "" {
			archive_page(*pageArchiveFlag)
		} else if *gallerypageArchiveFlag != "" {
			archive_gallery(*gallerypageArchiveFlag)
		}
	}

	if unarchiveCommand.Parsed() {
		if *pageUnarchiveFlag == "" && *gallerypageUnarchiveFlag == "" {
			fmt.Println("Please provide the page or gallery name using -page, -blog or -gallery option.")
		} else if *pageUnarchiveFlag != "" {
			unarchive_page(*pageUnarchiveFlag)
		} else if *gallerypageUnarchiveFlag != "" {
			unarchive_gallery(*gallerypageUnarchiveFlag)
		}
	}
}
