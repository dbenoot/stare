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
    "strconv"
    "strings"
)

func archive_page(pagename string) {
    fmt.Println("Archiving page " + pagename)
    
    var language string
    
    if site.multiLang == true {
        language = strings.Split(pagename, "/")[0]
        pagename = strings.Split(pagename, "/")[1]
    }
    
    path := filepath.Join("archive", language, "pages")
    fmt.Println(path)
    _, err := os.Stat(path) 
    if err != nil {
        os.MkdirAll(path, 0755)
    }
    
    move(filepath.Join("pages", language) +"/"+ pagename + ".html", path +"/"+ pagename + ".html")
}

func archive_blog(blogname string) {
    
    var blogId int
    var language string
    
    // Check whether the path exists and create if necessary
    
     if site.multiLang == true {
        language = strings.Split(blogname, "/")[0]
        blogname = strings.Split(blogname, "/")[1]
    }   
    
    path := "archive/pages/"+language+"/blogs/"
    
    
    _, err := os.Stat(path) 
    if err != nil {
        os.MkdirAll(path, 0755)
    }
    
    // Read all blog posts which contain the entered string
    
    blogposts, _ := filepath.Glob("pages/"+language+"/blogs/"+blogname+"*")

    // Select the correct blog post by assigning the correct blogId
    // If only 1 post is applicable, set the blogId
    // If more than 1 post is applicable, ask which post is applicable and set the blogId

    if len(blogposts) == 1 {
        blogId = 0
    } else {
    
        for i := 0; i < len(blogposts); i++ {
            fmt.Println(strconv.Itoa(i) + " - "+blogposts[i])
        }
        fmt.Println("Which blog post should be archived?")
        if _, err := fmt.Scanf("%d", &blogId); err != nil {
            fmt.Printf("%s\n", err)
        }
    }
    
    // Check that the blogId can exist and archive the correct blogId
    
    if blogId >= len(blogposts) {
        fmt.Println("Blog post does not exist.")
    } else {
        filename := strings.Split(blogposts[blogId],"/")[len(strings.Split(blogposts[blogId],"/"))-1]
        fmt.Println("Archiving blog post ", filename)
        move(blogposts[blogId], "archive/pages/"+language+"/blogs/"+filename)
    }
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

