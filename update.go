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

package main

import (
    "fmt"
    //"os"
    //"path/filepath"
    //"strconv"
    //"strings"
)

func addLanguage (language string) {
    fmt.Println ("Adding!")
    
    // TODO
    //
    // 1. Create folder 'language' under site.pagedir
    // 2. Edit config.ini: 
    //      - add the provided string to the key languages
    
}

func migLanguage (language string) {
    fmt.Println ("Migrating!")
    
    // TODO
    //
    // 1. Create folder 'language' under site.pagedir and move all html files there
    // 2. Edit config.ini: 
    //      - update key multiple_language_support -> y
    //      - create key primary_language and complete with the provided string 
    //      - create key languages and add the provided string
    
    
}