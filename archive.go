package main

import (
    "fmt"
    /*"os"*/
    "path/filepath"
)

func archive_page(pagename string) {
    fmt.Println("Archiving page " + pagename)
    move("pages" + string(filepath.Separator) + pagename+".html", "archive" + string(filepath.Separator) + "pages" + string(filepath.Separator) + pagename + ".html")
}

func archive_gallery(galleryname string) {
    fmt.Printf("Archiving gallery %q", galleryname)
    move("gallery" + string(filepath.Separator) + galleryname, "archive" + string(filepath.Separator) + "gallery" + string(filepath.Separator) + galleryname)
}

