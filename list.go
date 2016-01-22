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