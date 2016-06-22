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
var nowFile = time.Now().Format(time.RFC3339)

func createPage (pagename string, languagedir string) {

    copyfile("." + string(filepath.Separator) + "templates" + string(filepath.Separator) + "page_template.html", "." + string(filepath.Separator) + "pages" + string(filepath.Separator) + languagedir + string(filepath.Separator) + pagename + ".html")

    
    prepend("status          : in_draft\n------------------------------------------------------------------------", "pages/"+languagedir+"/"+pagename+".html")
    
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
 
            prepend("menu name       : "+menuname, "pages/"+languagedir+"/"+pagename+".html")
            prepend("menu order      : "+menuorder, "pages/"+languagedir+"/"+pagename+".html")
            prepend("present in menu : y", "pages/"+languagedir+"/"+pagename+".html")
        } else {
            prepend("present in menu : n\nmenu order      : nap\nmenu name       : nap", "pages/"+languagedir+"/"+pagename+".html")
        }

    prepend("------------------------------------------------------------------------\ncreated on      : "+now, "pages/"+languagedir+"/"+pagename+".html")

}

func createBlog (blogName string, languagedir string) {
    
    if _, err := os.Stat("pages/"+languagedir+"/blogs/"); os.IsNotExist(err) {
        os.MkdirAll(filepath.Join("pages",languagedir,"blogs"), 0755)
    }
    
    filename := nowFile+"-"+blogName+".md"
    
    os.Create("pages/"+languagedir+"/blogs/"+filename)
    
    prepend("status          : in_draft\n------------------------------------------------------------------------", "pages/"+languagedir+"/blogs/"+filename)
    prepend("taxonomies      : ", "pages/"+languagedir+"/blogs/"+filename)
    prepend("created by      : ", "pages/"+languagedir+"/blogs/"+filename)
    prepend("created on      : "+now, "pages/"+languagedir+"/blogs/"+filename)
    prepend("title           : ", "pages/"+languagedir+"/blogs/"+filename)
    prepend("------------------------------------------------------------------------", "pages/"+languagedir+"/blogs/"+filename)
}

func create_page (pageName string) {

    fmt.Println("Creating page " + pageName)
    if site.multiLang == true {
        for i := 0; i < len(site.languages); i++ {
            fmt.Println(site.languages[i])
            
            os.MkdirAll("pages" + string(filepath.Separator) + site.languages[i] ,0755)
            createPage(pageName, site.languages[i])
        }
    } else {
        languagedir := ""
        createPage(pageName, languagedir)
    }
}

func create_blog (blogName string) {

    fmt.Println("Creating blog " + blogName)
    if site.multiLang == true {
        for i := 0; i < len(site.languages); i++ {

            createBlog(blogName, site.languages[i])
        }
    } else {
        languagedir := ""
        createBlog(blogName, languagedir)
    }
}

func create_gallery(galleryname string) {
    fmt.Println("Creating gallery " + galleryname)
    os.MkdirAll("pages" + string(filepath.Separator) + "gallery" + string(filepath.Separator) + galleryname ,0755)
}