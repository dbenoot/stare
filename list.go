/* TODO

- list all pages
- list all galleries

*/

package main

import (
    "fmt"
    )

func sourcelist() {
    fmt.Println("PAGES")
    list("pages/*.html")
    
    fmt.Println("GALLERIES")
    listdir("gallery")
    
    return
    }