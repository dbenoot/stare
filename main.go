/* TODO

- refine list command output
- create render engine
- create creation engine

*/

package main

import (
    "fmt"
    "flag"
    "os"
)

func main() {

    renderCommand := flag.NewFlagSet("render", flag.ExitOnError)

	createCommand := flag.NewFlagSet("create", flag.ExitOnError)
	pageCreateFlag := createCommand.String("page", "", "Name of the new page.")
	galleryCreateFlag := createCommand.String("gallery", "", "Name of the new gallery.")

	archiveCommand := flag.NewFlagSet("archive", flag.ExitOnError)
	pageArchiveFlag := archiveCommand.String("page", "", "Name of the page to be archived.")
	gallerypageArchiveFlag := archiveCommand.String("gallery", "", "Name of the gallery to be archived.")

	postCommand := flag.NewFlagSet("post", flag.ExitOnError)
	pagePostFlag := postCommand.String("page", "", "Name of the page to be posted.")

	unpostCommand := flag.NewFlagSet("unpost", flag.ExitOnError)
	pageUnpostFlag := unpostCommand.String("page", "", "Name of the page to be unposted.")

    listCommand := flag.NewFlagSet("list", flag.ExitOnError)
	
	if len(os.Args) == 1 {
		fmt.Println("usage: stare <command> [<args>]")
		fmt.Println("The most commonly used stare commands are: \n")
		fmt.Println(" render      Renders the website.\n")
		fmt.Println(" list        Lists all pages and galleries.\n")
		fmt.Println(" create")
		fmt.Println("   -page     Creates a new page")
		fmt.Println("   -gallery  Create a new gallery\n")
		fmt.Println(" post")
		fmt.Println("   -page     Posts a page\n")
		fmt.Println(" unpost")
		fmt.Println("   -page     Unposts a page\n")
		fmt.Println(" archive     Archives a page.")
		fmt.Println("   -page     Archives a page")
		fmt.Println("   -gallery  Archives a gallery\n")
		return
	}

	switch os.Args[1] {
	case "render":
		renderCommand.Parse(os.Args[2:])
	case "create":
		createCommand.Parse(os.Args[2:])
	case "archive":
		archiveCommand.Parse(os.Args[2:])
	case "list":
	    listCommand.Parse(os.Args[2:])
	case "post":
		postCommand.Parse(os.Args[2:])
	case "unpost":
		unpostCommand.Parse(os.Args[2:])
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}

	if renderCommand.Parsed() {
		render_site()
	}

	if listCommand.Parsed() {
		sourcelist()
	}

	if createCommand.Parsed() {
		if *pageCreateFlag == "" && *galleryCreateFlag == "" {
			fmt.Println("Please provide the page name using -page option or the gallery name using -gallery.")
			return
		}
        if *pageCreateFlag != "" && *galleryCreateFlag == ""  {
            create_page(*pageCreateFlag)
            return
        }
        if *pageCreateFlag == "" && *galleryCreateFlag != "" {
            create_gallery(*galleryCreateFlag)
            return
        }
        if *pageCreateFlag != "" && *galleryCreateFlag != "" {
            create_page(*pageCreateFlag)
            create_gallery(*galleryCreateFlag)
            return
        } 
	}

	if archiveCommand.Parsed() {
		if *pageArchiveFlag == "" && *gallerypageArchiveFlag == "" {
			fmt.Println("Please provide the page name using -page option or the gallery name using -gallery.")
			return
		}

        if *pageArchiveFlag != "" && *gallerypageArchiveFlag == "" {
            archive_page(*pageArchiveFlag)
            return
        }
        
       if *pageArchiveFlag == "" && *gallerypageArchiveFlag != "" {
            archive_gallery(*gallerypageArchiveFlag)
            return
        }
        if *pageArchiveFlag != "" && *gallerypageArchiveFlag != "" {
            archive_page(*pageArchiveFlag)
            archive_gallery(*gallerypageArchiveFlag)
            return
        }
	}
	
	if postCommand.Parsed() {
		post_page(*pagePostFlag)
	}

	if unpostCommand.Parsed() {
		unpost_page(*pageUnpostFlag)
	}

    
}