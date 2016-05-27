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
    "path/filepath"
    //"github.com/go-ini/ini"
    "os"
    //"strconv"
    "strings"
)

func addLanguage (language string) {

    // TODO
    //
    // 1. Create folder 'language' under site.pagedir
    // 2. Edit config.ini: 
    //      - add the provided string to the key languages
    
    path := filepath.Join(site.pagedir, language)
    templatepath := filepath.Join(site.templatedir , language)
    // err := 
    os.MkdirAll (path, 0755)
    os.MkdirAll (templatepath, 0755)
    // if err != nil {
    //     return err
    // }
    
    templates, _ := filepath.Glob(filepath.Join(site.templatedir,site.primaryLang,"/*.html"))
    for c := 0; c < len(templates); c++ {
        copyfile (templates[c], filepath.Join(templatepath, strings.Split(templates[c], "/")[len(strings.Split(templates[c], "/"))-1]))
    }
    
    var inilang string
    
    for lang := 0; lang < len(site.languages); lang++ {
        if lang == 0 {
            inilang = site.languages[lang]
        } else {
            inilang = inilang + ", " + site.languages[lang]
        }
    }
    
    inilang = inilang + ", " + language
    
    cfg.Section("general").NewKey("languages", inilang)
    cfg.SaveTo("config.ini")
}

func migLanguage (language string) {
    fmt.Println ("Migrating!")
    
    path := filepath.Join(site.pagedir, language)
    templatepath := filepath.Join(site.templatedir , language)
    // err := 
    os.MkdirAll (path, 0755)
    os.MkdirAll (templatepath, 0755)
    // if err != nil {
    //     return err
    // }
    pages, _ := filepath.Glob(site.pagedir+"/*.html")
    for c := 1; c < len(pages); c++ {
        copyfile (pages[c], filepath.Join(path, strings.Split(pages[c], "/")[len(strings.Split(pages[c], "/"))-1]))
        move(pages[c], filepath.Join("archive",pages[c]))
    }
    
    templates, _ := filepath.Glob(site.templatedir+"/*.html")
    for c := 0; c < len(templates); c++ {
        copyfile (templates[c], filepath.Join(templatepath, strings.Split(templates[c], "/")[len(strings.Split(templates[c], "/"))-1]))
        move(templates[c], filepath.Join("archive",templates[c]))
    }
    
    cfg.Section("general").NewKey("multiple_language_support", "y")
    cfg.Section("general").NewKey("primary_language", language)
    cfg.Section("general").NewKey("languages", language)
    cfg.SaveTo("config.ini")    
}