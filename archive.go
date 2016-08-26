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


// This file contains the functions for archiving and unarchiving web pages and galleries

// TODO
//
// - archiving and unarchiving galleries

package main

import (
    "fmt"
    "os"
    "io/ioutil"
    "path/filepath"
    "strings"
)

func archive_page(pagename string) {
    
    var language string
    
    if site.multiLang == true {
        language = strings.Split(pagename, "/")[0]
        pagename = strings.Split(pagename, "/")[1]
    }
    
    path := filepath.Join("archive", language, "pages")

    _, err := os.Stat(path) 
    if err != nil {
        os.MkdirAll(path, 0755)
    }
    
    pages, _ := filepath.Glob(site.pagedir+"/"+pagename+"*") 
    
    page := findItem(pages)
    
    if len(page) > 0 {
        filename := strings.Split(page,"/")[len(strings.Split(page,"/"))-1]
        fmt.Println("Archiving page", filename)
        move(page, "archive/pages/"+language+"/"+filename)
    }
    
}

func unarchive_page(pagename string) {
    
    var language string
    
    if site.multiLang == true {
        language = strings.Split(pagename, "/")[0]
        pagename = strings.Split(pagename, "/")[1]
    }
    
    path := filepath.Join("archive", language, "pages")

    pages, _ := filepath.Glob(path+"/"+pagename+"*") 
    
    page := findItem(pages)
    
    if len(page) > 0 {
        filename := strings.Split(page,"/")[len(strings.Split(page,"/"))-1]
        fmt.Println("Unarchiving page", filename)
        move(page, "pages/"+language+"/"+filename)
    }  
    
}

func archive_blog(blogname string) {
    
    var language string

    if site.multiLang == true {
        language = strings.Split(blogname, "/")[0]
        blogname = strings.Split(blogname, "/")[1]
    }   
    
    // Check whether the path exists and create if necessary
    
    path := "archive/pages/"+language+"/blogs/"

    _, err := os.Stat(path) 
    if err != nil {
        os.MkdirAll(path, 0755)
    }
    
    // Read all blog posts which contain the entered string
    
    blogposts, _ := filepath.Glob("pages/"+language+"/blogs/"+blogname+"*")

    // Select the correct blog in case of name ambiguity
    
    blogpost := findItem(blogposts)
    
    if len(blogpost) > 0 {
        filename := strings.Split(blogpost,"/")[len(strings.Split(blogpost,"/"))-1]
        fmt.Println("Archiving blog post", filename)
        move(blogpost, "archive/pages/"+language+"/blogs/"+filename)
    }
    
}

func unarchive_blog(blogname string) {

    var language string

    if site.multiLang == true {
        language = strings.Split(blogname, "/")[0]
        blogname = strings.Split(blogname, "/")[1]
    }   
    
    // Check whether the path exists and create if necessary
    
    path := "archive/pages/"+language+"/blogs/"

    // Read all blog posts which contain the entered string
    
    blogposts, _ := filepath.Glob(path+blogname+"*")

    // Select the correct blog in case of name ambiguity
    
    blogpost := findItem(blogposts)
    
    if len(blogpost) > 0 {
        filename := strings.Split(blogpost,"/")[len(strings.Split(blogpost,"/"))-1]
        fmt.Println("Unarchiving blog post", filename)
        move(blogpost, "pages/"+language+"/blogs/"+filename)
    }
    
}

//TODO add something like finditems for the galleries so incomplete names can be further selected

func archive_gallery(galleryname string) {
    
    path := filepath.Join("archive", "pages", "gallery")
    _, err := os.Stat(path) 
    if err != nil {
        os.MkdirAll(path, 0755)
    }

    fmt.Println("Archiving gallery", galleryname)
    movedir(filepath.Join("pages", "gallery", galleryname), filepath.Join(path, galleryname))
    
    files, _ := ioutil.ReadDir(site.gallerydir)
    
    if len(files) == 0 {
        cfg.Section("general").NewKey("gallery", "")
        cfg.SaveTo("config.ini")
    }
}

func unarchive_gallery(galleryname string) {
    path := filepath.Join("archive", "pages", "gallery")
    _, err := os.Stat(path) 
    if err != nil {
        os.MkdirAll(path, 0755)
    }

    fmt.Println("Archiving gallery", galleryname)
    movedir(filepath.Join(path, galleryname), filepath.Join("pages","gallery", galleryname))
    
    if site.gallery == false {
        cfg.Section("general").NewKey("gallery", "y")
        cfg.SaveTo("config.ini")
    }
}