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
    /*os.MkdirAll("." + string(filepath.Separator) + pagename + string(filepath.Separator) + "src" + string(filepath.Separator) + "gallery",0755)
    os.MkdirAll("." + string(filepath.Separator) + pagename + string(filepath.Separator) + "src" + string(filepath.Separator) + "pages",0755)
    os.MkdirAll("." + string(filepath.Separator) + pagename + string(filepath.Separator) + "src" + string(filepath.Separator) + "templates",0755)
    os.MkdirAll("." + string(filepath.Separator) + pagename + string(filepath.Separator) + "rendered",0755)*/
    
    fmt.Println("Creating page " + pagename)
    copyfile("." + string(filepath.Separator) + "templates" + string(filepath.Separator) + "page_template.html", "." + string(filepath.Separator) + "pages" + string(filepath.Separator) + pagename + ".html")
    
    
    page := string("pages/"+pagename+".html")
    fmt.Println(page)
    
    now := time.Now().Format(time.RFC1123)
    prepend("in_draft\n~~>>", "pages/"+pagename+".html")
    
    var menuyn string
    fmt.Println("Present in menubar (y/n)")
            if _, err := fmt.Scanf("%s", &menuyn); err != nil {
            fmt.Printf("%s\n", err)
            return
        }
        if menuyn == "y" {
            prepend("in_menu", "pages/"+pagename+".html")
                var menuorder string
                    fmt.Println("Place in menubar (0-9)")
                    if _, err := fmt.Scanf("%s", &menuorder); err != nil {
                        fmt.Printf("%s\n", err)
                    return
                    }
                prepend("menu_order_"+menuorder, "pages/"+pagename+".html")
        

        }



    prepend("<<~~\n"+now, "pages/"+pagename+".html")

}

func create_gallery(galleryname string) {
    fmt.Println("Creating gallery " + galleryname)
    os.MkdirAll("." + string(filepath.Separator) + "gallery" + string(filepath.Separator) + galleryname ,0755)
}