/* copy page from main_template file */

package main

import (
    "fmt"
    /*"os"
    "path/filepath"*/
)

func create_page(pagename string) {
    /*os.MkdirAll("." + string(filepath.Separator) + pagename + string(filepath.Separator) + "src" + string(filepath.Separator) + "gallery",0755)
    os.MkdirAll("." + string(filepath.Separator) + pagename + string(filepath.Separator) + "src" + string(filepath.Separator) + "pages",0755)
    os.MkdirAll("." + string(filepath.Separator) + pagename + string(filepath.Separator) + "src" + string(filepath.Separator) + "templates",0755)
    os.MkdirAll("." + string(filepath.Separator) + pagename + string(filepath.Separator) + "rendered",0755)*/
    
    fmt.Println("Creating page " + pagename)
/*    copy("pages" + string(filepath.Separator) + pagename+".html", "archive" + string(filepath.Separator) + "pages" + string(filepath.Separator) + pagename + ".html")*/

}

func create_gallery(galleryname string) {
    fmt.Println("Creating gallery " + galleryname)
}