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

// Render all production pages and galleries

package main

import (
        "fmt"
        "io/ioutil"
        "log"
        "os"
        "io"
        "path/filepath"
)

func init_site() {
    
    if checkSiteNotExist() == false {
        fmt.Println("Directory not empty.")
    } else {
        createConfigFile()
        createFolders()
        createTemplates()
    }
}

func checkSiteNotExist () bool {
    ok, err := IsDirEmpty("./")
    if err != nil {
         fmt.Println(err)
    }

    return ok
}

func createConfigFile() {
    fmt.Println("creating config file")
    
    output := "[general]\nmultiple_language_support = \nprimary_language = \nlanguages = "
    
    err := ioutil.WriteFile("config.ini", []byte(output), 0644)
    if err != nil {
            log.Fatalln(err)
    } 
    
}

func createFolders() {
    fmt.Println("creating folders")
    
	dirs := []string{
    	filepath.Join("archive","pages","blogs"),
    	filepath.Join("archive","pages","gallery"),
    	filepath.Join("pages","blogs"),
    	filepath.Join("pages","gallery"),
    	"src",
    	"templates",
	}
    
    for i := 0; i < len(dirs); i++ {
        os.MkdirAll(dirs[i] ,0755)
    }
    
    
}

func createTemplates() {
    fmt.Println("creating templates")
}

 func IsDirEmpty(name string) (bool, error) {
         f, err := os.Open(name)
         if err != nil {
                 return false, err
         }
         defer f.Close()

         _, err = f.Readdir(1)

         if err == io.EOF {
                 return true, nil
         }
         return false, err
 }