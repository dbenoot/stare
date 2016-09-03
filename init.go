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

// Initialize the site

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
        
        fmt.Println("Congratulations! You have initiated your new website.")
        
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

    output := "[general]\nmultiple_language_support = \nprimary_language = \nlanguages = \ngallery = "
    
    err := ioutil.WriteFile("config.ini", []byte(output), 0644)
    if err != nil {
            log.Fatalln(err)
    } 
    
}

func createFolders() {
    for i := 0; i < len(dirs); i++ {
        os.MkdirAll(dirs[i] ,0755)
    }
}

func createTemplates() {
    for key, value := range templ {
        initTemplate(key, value)
    }
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
 
func initTemplate (filename string, filecontent string) {
    err := ioutil.WriteFile(filepath.Join("templates",filename), []byte(filecontent), 0644)
    if err != nil {
            log.Fatalln(err)
    } 
}

var templ = map[string]string {
    "blog_template.html": "<!DOCTYPE html>\n<html lang=\"en\">\n  <head>\n    <meta charset=\"utf-8\">\n    <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n    <title><<~~TITLE~~>></title>\n  </head>\n  \n  <body>\n    <<~~NAVBAR~~>>\n  </body>\n</html>",
    "blogpost_template.html": "                <div>\n                  <h2><<~~BLOGPOST:TITLE~~>></h2>\n                </div>\n                <div>\n                  <<~~BLOGPOST:AUTHOR~~>>\n                </div>\n                <div>\n                  <<~~BLOGPOST:TIME~~>>\n                </div>\n                <div>\n                  <<~~BLOGPOST:CONTENT~~>>\n                </div>",
    "blogtitles_template.html": "                <div>\n                  <strong><<~~BLOGPOST:TITLE~~>></strong> - <em><<~~BLOGPOST:AUTHOR~~>></em>\n                </div>",
    "footer_template.html": "",
    "gallery_item.html": "					<div>\n						<a href=\"<<~~SUBGALLERYLINK~~>>\">\n						<div><img src=\"<<~~SUBGALLERYTHUMB~~>>\"></div>\n						<div><h3><<~~SUBGALLERYNAME~~>></h3></div>\n						</a>\n					</div>\n<<~~GALLERYITEM~~>>",
    "gallery_template.html": "<!DOCTYPE html>\n<html lang=\"en\">\n  <head>\n    <meta charset=\"utf-8\">\n    <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n    <title><<~~TITLE~~>></title>\n  </head>\n  \n  <body>\n    \n<<~~NAVBAR~~>>\n\n 		<div>\n<<~~GALLERYITEM~~>>\n    </div>\n\n<<~~FOOTER~~>>\n\n  </body>\n</html>",
    "header_template.html": "",
    "langbar_item.html": "					<li><a href=\"<<~~LANGLINK~~>>\"><<~~LANGITEM~~>></a></li>",
    "navbar_item.html": "					<li <<~~NAVACTIVE~~>>><a href=\"<<~~NAVLINK~~>>\"><<~~NAVITEM~~>></a></li>",
    "navbar_template.html": "	<nav class=\"navbar navbar-default navbar-fixed-top\">\n		<div class=\"container-fluid\">\n			<div class=\"navbar-header\">\n				<button type=\"button\" class=\"navbar-toggle collapsed\" data-toggle=\"collapse\" data-target=\"#navbar\" aria-expanded=\"false\" aria-controls=\"navbar\">\n					<span class=\"sr-only\">Toggle navigation</span>\n					<span class=\"icon-bar\"></span>\n					<span class=\"icon-bar\"></span>\n					<span class=\"icon-bar\"></span>\n				</button>\n				<a class=\"navbar-brand\" href=\"#\">Stare Website</a>\n			</div>\n			<div id=\"navbar\" class=\"navbar-collapse collapse\">\n				<ul class=\"nav navbar-nav\">\n					<<~~NAVLIST~~>>\n				</ul>\n				<ul class=\"nav navbar-nav navbar-right\">\n					<li><a href=\"#\">Item 1</a></li>\n					<li><a href=\"#\">Item 2</a></li>\n					<li><a href=\"#\">Item 3</a></li>\n				</ul>\n			</div>\n		</div>\n	</nav>",
    "page_template.html": "<!DOCTYPE html>\n<html lang=\"en\">\n  <head>\n    <meta charset=\"utf-8\">\n    <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n    <title><<~~TITLE~~>></title>\n  </head>\n  \n  <body>\n    \n    <<~~NAVBAR~~>>\n\n  </body>\n</html>",
    "subgallery_item.html": "                <div><a href=\"<<~~SUBIMAGE~~>>\"><img src=\"<<~~SUBIMAGETHUMB~~>>\" width=\"350\" height=\"275\"></a></div>\n<<~~SUBGALLERYITEM~~>>",
    "subgallery_template.html": "<!DOCTYPE html>\n<html lang=\"en\">\n  <head>\n    <meta charset=\"utf-8\">\n    <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n    <title><<~~TITLE~~>></title>\n  </head>\n  \n  <body>\n    \n    <<~~NAVBAR~~>>\n    \n    <h1><<~~GALLERYTITLE~~>></h1>\n    <a href=\"../gallery.html\">Back</a>\n    \n    <<~~SUBGALLERYITEM~~>>\n    \n    <<~~FOOTER~~>>\n    \n  </body>\n</html>",
}

var dirs = []string{
	filepath.Join("archive","pages","blogs"),
	filepath.Join("archive","pages","gallery"),
	filepath.Join("pages","blogs"),
	filepath.Join("pages","gallery"),
	"src",
	"templates",
}