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
    "fmt"
    "flag"
    "os"
)

func main() {

	createCommand := flag.NewFlagSet("create", flag.ExitOnError)
	pageCreateFlag := createCommand.String("page", "", "Name of the new page.")
	galleryCreateFlag := createCommand.String("gallery", "", "Name of the new gallery.")
	blogCreateFlag := createCommand.String("blog", "", "Name of the new gallery.")

	archiveCommand := flag.NewFlagSet("archive", flag.ExitOnError)
	pageArchiveFlag := archiveCommand.String("page", "", "Name of the page to be archived.")
	gallerypageArchiveFlag := archiveCommand.String("gallery", "", "Name of the gallery to be archived.")
	blogArchiveFlag := archiveCommand.String("blog", "", "Name of the blog post to be archived.")

	postCommand := flag.NewFlagSet("post", flag.ExitOnError)
	pagePostFlag := postCommand.String("page", "", "Name of the page to be posted.")
	blogPostFlag := postCommand.String("blog", "", "Name of the page to be posted.")

	unpostCommand := flag.NewFlagSet("unpost", flag.ExitOnError)
	pageUnpostFlag := unpostCommand.String("page", "", "Name of the page to be unposted.")
	blogUnpostFlag := unpostCommand.String("blog", "", "Name of the page to be unposted.")

	if len(os.Args) == 1 {
		fmt.Println("usage: stare <command> [<args>]")
		fmt.Println("The most commonly used stare commands are: \n")
		fmt.Println(" render      Renders the website.\n")
		fmt.Println(" list        Lists all pages, blog posts and galleries.\n")
		fmt.Println(" create")
		fmt.Println("   -page     Creates a new page")
		fmt.Println("   -gallery  Create a new gallery")
		fmt.Println("   -blog     Create a new blog post\n")
		fmt.Println(" post")
		fmt.Println("   -page     Posts a page")
		fmt.Println("   -blog     Posts a blog post\n")
		fmt.Println(" unpost")
		fmt.Println("   -page     Unposts a page")
		fmt.Println("   -blog     Unposts a blog post\n")
		fmt.Println(" archive")
		fmt.Println("   -page     Archives a page")
		fmt.Println("   -gallery  Archives a gallery")
		fmt.Println("   -blog     Archives a blog post\n")
		return
	}

	switch os.Args[1] {
	case "render":
		render_site()
	case "create":
		createCommand.Parse(os.Args[2:])
	case "archive":
		archiveCommand.Parse(os.Args[2:])
	case "list":
	    sourcelist()
	case "post":
		postCommand.Parse(os.Args[2:])
	case "unpost":
		unpostCommand.Parse(os.Args[2:])
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}

	if createCommand.Parsed() {
		if *pageCreateFlag == "" && *galleryCreateFlag == "" && *blogCreateFlag == "" {
			fmt.Println("Please provide the page, gallery or blog name using -page, -gallery or -blog parameter.")
		} else if *pageCreateFlag != "" {
        	create_page(*pageCreateFlag)
        	fmt.Println(*pageCreateFlag)
		} else if *galleryCreateFlag != "" {
        	create_gallery(*galleryCreateFlag)
		} else if *blogCreateFlag != "" {
        	create_blog(*blogCreateFlag)
        }	
    }

	if archiveCommand.Parsed() {
		if *pageArchiveFlag == "" && *gallerypageArchiveFlag == "" && *blogArchiveFlag == ""{
			fmt.Println("Please provide the page name using -page option or the gallery name using -gallery.")
		} else if *pageArchiveFlag != "" {
            archive_page(*pageArchiveFlag)
        } else if *gallerypageArchiveFlag != "" {
            archive_gallery(*gallerypageArchiveFlag)
        } else if *blogArchiveFlag != "" {
            archive_blog(*blogArchiveFlag)
        }
	}
	
	if postCommand.Parsed() {
		if *pagePostFlag == "" && *blogPostFlag == "" {
			fmt.Println("Please provide the page name using -page option or the blog post using -blog.")
		} else if *pagePostFlag != "" {
			post(*pagePostFlag, "pages/")
		} else if *blogPostFlag != "" {
			post(*blogPostFlag, "pages/blogs/")
		}
	}

	if unpostCommand.Parsed() {
		if *pageUnpostFlag == "" && *blogUnpostFlag == "" {
			fmt.Println("Please provide the page name using -page option or the blog post using -blog.")
		} else if *pageUnpostFlag != "" {
			unpost(*pageUnpostFlag, "pages/")
		} else if *blogUnpostFlag != "" {
			unpost(*blogUnpostFlag, "pages/blogs/")
		}
	}

}