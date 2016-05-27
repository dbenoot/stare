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


// Post draft pages in production
// Unpost production pages to draft

package main

import (
    "path/filepath"
    "fmt"
    "strings"
    )

var itemId int

func post (name string, path string) {
    
    items, _ := filepath.Glob(path+"*"+name+"*")

    // Select the correct item in case of name ambiguity

    item := findItem(items)
    
    if len(item) > 0 {
        filename := strings.Split(item,"/")[len(strings.Split(item,"/"))-1]
        fmt.Println("Posting", filename)
        substitute_in_header(item, "in_draft", "posted")
    }
    
}

func unpost (name string, path string) {
    
    items, _ := filepath.Glob(path+"*"+name+"*")

    // Select the correct item in case of name ambiguity
    
    item := findItem(items)
    
    if len(item) > 0 {
        filename := strings.Split(item,"/")[len(strings.Split(item,"/"))-1]
        fmt.Println("Unposting", filename)
        substitute_in_header(item, "posted", "in_draft")
    }
}