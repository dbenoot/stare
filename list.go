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


// List of pages and galleries


package main

import (
    "fmt"
    "io/ioutil"
    "path/filepath"
    )

var j int

func sourcelist() {
    fmt.Println("PAGES")
    j = 1
    if site.multiLang == true {
        for i := 1; i <= len(site.languages); i++ {
            list("pages/"+site.languages[i]+"/*.html")
        }
    } else {
        list("pages/*.html")
    }
    
    fmt.Println("\nBLOG POSTS")
    j = 1
    if site.multiLang == true {
        for i := 0; i < len(site.languages); i++ {
            list("pages/"+site.languages[i]+"/blogs/*.md")
        }
    } else {
        list("pages/blogs/*.md")
    }    
    
    
    fmt.Println("\nGALLERIES")
    listdir("pages/gallery")
    
    return
    }
    
func list (folder string) {
    files, _ := filepath.Glob(folder)

    
    for i:= 0; i < len(files); i++ {
           fmt.Println(j , " - " , files[i])
           j += 1 
    }    
    
    }

func listdir (folder string) {
    files, _ := ioutil.ReadDir(folder)
    
    i := 1
    for _, f := range files {
    fmt.Println(i, " - ", f.Name())
    i += 1
    }
    
    }