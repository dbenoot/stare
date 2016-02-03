/* TODO

- add metadata to new page: 
    - date and time of creation
    - availability in the menu
    - order in the menu
    - draft
- add metadata to new gallery
    - date and time of creation
    - draft
- create gallery folder

*/

package main

import (
    "fmt"
    "os"
    "path/filepath"
    "time"
    //"io/ioutil"
    //"log"
    //"strings"
)

func create_page(pagename string) {

    fmt.Println("Creating page " + pagename)
    copyfile("." + string(filepath.Separator) + "templates" + string(filepath.Separator) + "page_template.html", "." + string(filepath.Separator) + "pages" + string(filepath.Separator) + pagename + ".html")

//    page := string("pages/"+pagename+".html")

    now := time.Now().Format(time.RFC1123)
    prepend("status          : in_draft\n------------------------------------------------------------------------", "pages/"+pagename+".html")
    
    var menuyn string
    fmt.Println("Present in menubar (y/n)")
            if _, err := fmt.Scanf("%s", &menuyn); err != nil {
            fmt.Printf("%s\n", err)
            return
        }
        if menuyn == "y" {
                var menuorder string
                    fmt.Println("Place in menubar (0-9)")
                    if _, err := fmt.Scanf("%s", &menuorder); err != nil {
                        fmt.Printf("%s\n", err)
                    return
                    }
                var menuname string
                    fmt.Println("Name of page in the menubar")
                    if _, err := fmt.Scanf("%s", &menuname); err != nil {
                        fmt.Printf("%s\n", err)
                    return
                    }
            prepend("menu name       : "+menuname, "pages/"+pagename+".html")
            prepend("menu order      : "+menuorder, "pages/"+pagename+".html")
            prepend("present in menu : y", "pages/"+pagename+".html")
        } else {
            prepend("present in menu : n\nmenu order      : nap\nmenu name       : nap", "pages/"+pagename+".html")
        }

    prepend("------------------------------------------------------------------------\ncreated on      : "+now, "pages/"+pagename+".html")
}

func create_gallery(galleryname string) {
    fmt.Println("Creating gallery " + galleryname)
    os.MkdirAll("pages" + string(filepath.Separator) + "gallery" + string(filepath.Separator) + galleryname ,0755)
}