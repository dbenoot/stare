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
    "strconv"
    "strings"
    )

var itemId int

func post (name string, path string) {
    
    items, _ := filepath.Glob(path+"*"+name+"*")

    // Select the correct blog post by assigning the correct itemId
    // If only 1 item is applicable, set the itemId
    // If more than 1 item is applicable, ask which item is applicable and set the itemId

    if len(items) == 1 {
        itemId = 0
    } else {
    
        for i := 0; i < len(items); i++ {
            fmt.Println(strconv.Itoa(i) + " - "+items[i])
        }
        fmt.Println("Which item should be posted?")
        if _, err := fmt.Scanf("%d", &itemId); err != nil {
            fmt.Printf("%s\n", err)
        }
    }
    
    // Check that the blogId can exist and post the correct blogId
    
    if itemId >= len(items) {
        fmt.Println("Item does not exist.")
    } else {
        filename := strings.Split(items[itemId],"/")[len(strings.Split(items[itemId],"/"))-1]
        fmt.Println("Posting ", filename)
        substitute_in_header(items[itemId], "in_draft", "posted")
    }
    
}

func unpost (name string, path string) {
    
    items, _ := filepath.Glob(path+"*"+name+"*")

    // Select the correct blog post by assigning the correct itemId
    // If only 1 item is applicable, set the itenId
    // If more than 1 item is applicable, ask which item is applicable and set the itemId

    if len(items) == 1 {
        itemId = 0
    } else {
    
        for i := 0; i < len(items); i++ {
            fmt.Println(strconv.Itoa(i) + " - "+items[i])
        }
        fmt.Println("Which item should be posted?")
        if _, err := fmt.Scanf("%d", &itemId); err != nil {
            fmt.Printf("%s\n", err)
        }
    }
    
    // Check that the blogId can exist and post the correct blogId
    
    if itemId >= len(items) {
        fmt.Println("Item does not exist.")
    } else {
        filename := strings.Split(items[itemId],"/")[len(strings.Split(items[itemId],"/"))-1]
        fmt.Println("Unposting ", filename)
        substitute_in_header(items[itemId], "posted", "in_draft")
    }
}