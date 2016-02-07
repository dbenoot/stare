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


// This file contains the functions for archiving web pages and galleries


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

