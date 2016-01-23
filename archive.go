package main

import (
    "fmt"
    "os"
    "path/filepath"
)

func archive_page(pagename string) {
    fmt.Println("Archiving page " + pagename)
    
    path := "archive" + string(filepath.Separator) + "pages" + string(filepath.Separator)
    _, err := os.Stat(path) 
    if err != nil {
        os.MkdirAll(path, 0755)
    }
    
    move("pages" + string(filepath.Separator) + pagename+".html", path + pagename + ".html")
}

func archive_gallery(galleryname string) {
    fmt.Println("Archiving gallery ", galleryname)
    
    path := "archive" + string(filepath.Separator) + "gallery" + string(filepath.Separator)
    _, err := os.Stat(path) 
    if err != nil {
        os.MkdirAll(path, 0755)
    }

    move("gallery" + string(filepath.Separator) + galleryname, path + galleryname)
}

