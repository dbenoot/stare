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


// FUNCTIONS
//
//- create a page and add metadata to new page: 
//    - date and time of creation
//    - availability in the menu
//    - order in the menu
//    - name in menu
//    - draft
//- create gallery folder

package main

import (
    "fmt"
    "os"
    "path/filepath"
    "time"
    "bufio"
    //"log"
    "strings"
)

var now = time.Now().Format(time.RFC1123)

func create_page(pagename string) {

    fmt.Println("Creating page " + pagename)
    copyfile("." + string(filepath.Separator) + "templates" + string(filepath.Separator) + "page_template.html", "." + string(filepath.Separator) + "pages" + string(filepath.Separator) + pagename + ".html")

    
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
                    reader := bufio.NewReader(os.Stdin)
                    menuname, _ = reader.ReadString('\n')
                    menuname = strings.TrimSpace(menuname)
 
            prepend("menu name       : "+menuname, "pages/"+pagename+".html")
            prepend("menu order      : "+menuorder, "pages/"+pagename+".html")
            prepend("present in menu : y", "pages/"+pagename+".html")
        } else {
            prepend("present in menu : n\nmenu order      : nap\nmenu name       : nap", "pages/"+pagename+".html")
        }

    prepend("------------------------------------------------------------------------\ncreated on      : "+now, "pages/"+pagename+".html")
}

func create_blog (blogName string) {
    fmt.Println("Creating blog " + blogName)
    if _, err := os.Stat("pages/blogs/"); os.IsNotExist(err) {
        os.MkdirAll("pages/blogs/", 0755)
    }
    
    filename := blogName+" - "+now+".md"
    
    os.Create("pages/blogs/"+filename)
    
    prepend("status          : in_draft\n------------------------------------------------------------------------", "pages/blogs/"+filename)
    prepend("taxonomies      : ", "pages/blogs/"+filename)
    prepend("created by      :", "pages/blogs/"+filename)
    prepend("------------------------------------------------------------------------\ncreated on      : "+now, "pages/blogs/"+filename)
}

func create_gallery(galleryname string) {
    fmt.Println("Creating gallery " + galleryname)
    os.MkdirAll("pages" + string(filepath.Separator) + "gallery" + string(filepath.Separator) + galleryname ,0755)
}