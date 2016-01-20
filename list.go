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