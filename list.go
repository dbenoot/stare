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

func sourcelist() {
    fmt.Println("PAGES")
    list("pages/*.html")
    
    fmt.Println("\nGALLERIES")
    listdir("gallery")
    
    return
    }
    
func list (folder string) {
    files, _ := filepath.Glob(folder)

    i := 0
    for i < len(files) {
           fmt.Println(i+1 , " - " , files[i])
           i += 1 
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