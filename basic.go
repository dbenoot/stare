package main

import (
    "io"
    "io/ioutil"
    "os"
    "fmt"
    "path/filepath"
    )
    
func copy(inputname, outputname string) {
     readfile, err := os.Open(inputname)
     if err != nil {
         fmt.Println("Input file not available.\n")
     }
     defer readfile.Close()

     writefile, err := os.Create(outputname)
     if err != nil {
         fmt.Println("Output file could not be created. Please check that the standard folder structure is available.\n")
     }
     defer writefile.Close()

     // do the actual work
     io.Copy(writefile, readfile)
}

func move (inputname, outputname string) {
       err :=  os.Rename(inputname, outputname)

       if err != nil {
           fmt.Println("Page or gallery could not be moved. Please check standard folder structure.")
           return
       }
}

func list (folder string) {
    files, _ := filepath.Glob(folder)
    fmt.Println(files)
}

func listdir (folder string) {
    files, _ := ioutil.ReadDir(folder)
    for _, f := range files {
    fmt.Println(f.Name())
    }
}