package main

import (
    "os"
    "path/filepath"
)

func create_site(sitename string) {
    os.MkdirAll("." + string(filepath.Separator) + sitename + string(filepath.Separator) + "src" + string(filepath.Separator) + "gallery",0755)
}