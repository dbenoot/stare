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


// Basic functions


package main

import (
    "io"
    "os"
    "fmt"
    "log"
    "strings"
    "io/ioutil"
    )
    
func move (inputname, outputname string) {
       err :=  os.Rename(inputname, outputname)

       if err != nil {
           fmt.Println("Page or gallery could not be moved. Please check standard folder structure.")
           return
       }
}

func copydir(source string, dest string) (err error) {

     // get properties of source dir
     sourceinfo, err := os.Stat(source)
     if err != nil {
         return err
     }

     // create dest dir

     err = os.MkdirAll(dest, sourceinfo.Mode())
     if err != nil {
         return err
     }

     directory, _ := os.Open(source)

     objects, err := directory.Readdir(-1)

     for _, obj := range objects {

         sourcefilepointer := source + "/" + obj.Name()

         destinationfilepointer := dest + "/" + obj.Name()


         if obj.IsDir() {
             // create sub-directories - recursively
             err = copydir(sourcefilepointer, destinationfilepointer)
             if err != nil {
                 fmt.Println(err)
             }
         } else {
             // perform copy
             err = copyfile(sourcefilepointer, destinationfilepointer)
             if err != nil {
                 fmt.Println(err)
             }
         }

     }
     return
 }
 
func copyfile(source string, dest string) (err error) {
     sourcefile, err := os.Open(source)
     if err != nil {
         return err
     }

     defer sourcefile.Close()

     destfile, err := os.Create(dest)
     if err != nil {
         return err
     }

     defer destfile.Close()

     _, err = io.Copy(destfile, sourcefile)
     if err == nil {
         sourceinfo, err := os.Stat(source)
         if err != nil {
             err = os.Chmod(dest, sourceinfo.Mode())
         }

     }

     return
 }
 
func substitute (file, tie, replacetext string) {
        input, err := ioutil.ReadFile(file)
        if err != nil {
                log.Fatalln(err)
        }

        lines := strings.Split(string(input), "\n")

        for line := range lines {
                lines[line] = strings.Replace(lines[line], tie, replacetext, -1)
        }
        
        output := strings.Join(lines, "\n")
        err = ioutil.WriteFile(file, []byte(output), 0644)
        if err != nil {
                log.Fatalln(err)
        } 
}

func substitute_in_header (file, tie, replacetext string) {
        input, err := ioutil.ReadFile(file)
        if err != nil {
                log.Fatalln(err)
        }

        lines := strings.Split(string(input), "\n")

        for line := 0; line < 6; line++ {
                lines[line] = strings.Replace(lines[line], tie, replacetext, -1)
        }
        
        output := strings.Join(lines, "\n")
        err = ioutil.WriteFile(file, []byte(output), 0644)
        if err != nil {
                log.Fatalln(err)
        } 
}
