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
        "io/ioutil"
        "log"
        "strings"
        "os"
        "path/filepath"
        "fmt"
        //"github.com/go-ini/ini"
        "github.com/nfnt/resize"
        "image/jpeg"
        "strconv"
        "time"
        "path"
)

var wd, _ = os.Getwd()


// define Site functions

func (site Site) createFolder () {
        
        // Remove previous rendering of the site

        os.RemoveAll(path.Join(wd, "rendered"))
        
        // Create directories for temporary files and newly rendered site
        
        if site.multiLang == true {
                for i := 0; i < len(site.languages); i++ {
                        os.MkdirAll("temp/"+site.pagedir+"/"+site.languages[i], 0755)
                        os.MkdirAll("temp/"+site.pagedir+"/"+site.languages[i]+"/"+site.blogdir, 0755)
                        os.MkdirAll("temp/"+site.gallerydir+"/"+site.languages[i], 0755)
                        os.MkdirAll("rendered/"+site.pagedir+"/"+site.languages[i], 0755)
                        os.MkdirAll("rendered/"+site.blogdir+"/"+site.languages[i], 0755)
                        os.MkdirAll("rendered/"+site.gallerydir+"/"+site.languages[i], 0755)
                }   
        } else {
                os.MkdirAll("temp/"+site.pagedir, 0755)
                os.MkdirAll("temp/"+site.pagedir+"/"+site.blogdir, 0755)
                os.MkdirAll("temp/"+site.gallerydir, 0755)
                os.MkdirAll("rendered/"+site.pagedir, 0755)
                os.MkdirAll("rendered/"+site.blogdir, 0755)
                os.MkdirAll("rendered/"+site.gallerydir, 0755)
        }
        
}

func (site Site) copySrc () {
        srcItems, _ := filepath.Glob(site.srcdir+"/*")
        i := 0
        for i < len(srcItems) {
                copydir(srcItems[i], "rendered/"+strings.Split(srcItems[i], "/")[1])
                i += 1
        }
}

func (site Site) copyFiles () {

        // move the pages, blogs and galleries to the temporary directory
        if site.multiLang == true {
                for i := 0; i < len(site.languages); i++ {
                        copydir(site.pagedir+"/"+site.languages[i], "temp/"+site.pagedir+"/"+site.languages[i])
                        copydir(site.pagedir+"/"+site.languages[i]+"/"+site.blogdir, "temp/"+site.pagedir+"/"+site.languages[i]+"/"+site.blogdir)
                        copydir(site.gallerydir, "temp/"+site.gallerydir+"/"+site.languages[i]) //TODO remove 'pages/' form gallery in site object
                }
        } else {
                copydir(site.pagedir, "temp/"+site.pagedir)
                copydir(site.pagedir+"/"+site.blogdir, "temp/"+site.pagedir+"/"+site.blogdir)
                copydir(site.gallerydir, "temp/"+site.gallerydir) //TODO remove 'pages/' form gallery in site object
        }
        // copy the navbar template to the temp folder 

        copyfile(site.templatedir+"/navbar_template.html", "temp/navbar.html")
}

func (site Site) renderPages(pages []string, shortLang string, longLang string) {

        // complete menu and menuName
        // define posted and draft pages
        
        menu, menuName := createMenu(pages)
        all_pages, draft_pages := definePages(pages)

        // add navbar to the pages and resolve the ties NAVACTIVE, NAVLINK, NAVITEM
        // cycling through all posted pages, then cycling through all menu items
        // adding navlinks as necessary and resolving the ties
        
        for i := 0; i < len(all_pages); i++ {
                inject_html(all_pages[i], "<<~~NAVBAR~~>>", "temp/navbar.html")

                create_navbar(all_pages[i], menu, menuName, false)

                // populate the header and footer tie
                
                inject_html(all_pages[i], "<<~~HEADER~~>>", site.templatedir+"/header_template.html")
                inject_html(all_pages[i], "<<~~FOOTER~~>>", site.templatedir+"/footer_template.html")
                
                // resolve ties CSS, JS, PAGE
                
                if strings.Split(all_pages[i],"/")[2] == "index.html" {
                        substitute(all_pages[i], "<<~~JS~~>>","js/")
                        substitute(all_pages[i], "<<~~CSS~~>>","css/")
                        substitute(all_pages[i], "<<~~PAGE~~>>","pages/")
                } else {
                        substitute(all_pages[i], "<<~~JS~~>>","../js/")
                        substitute(all_pages[i], "<<~~CSS~~>>","../css/")
                        substitute(all_pages[i], "<<~~PAGE~~>>","")
                }
                
                remove_header(all_pages[i])
                
                // Copy files to the correct location
                
                // LANG update : if site.multiLang == true && shortLang == site.primaryLang then
                if strings.Split(all_pages[i],"/")[2] == "index.html" {
                        copyfile(all_pages[i], "rendered/"+strings.Split(all_pages[i],"/")[2])
                } else {
                        copyfile(all_pages[i], "rendered/"+site.pagedir+"/"+strings.Split(all_pages[i],"/")[2])
                }
                // else move everything in same folder --> !! similar updates in create_navbar
        }

        // give an overview of rendered pages and blogs and draft pages and blogs

        fmt.Println("The following pages were rendered: ")
        
        for i := 0; i < len (all_pages); i++ {
                fmt.Println(all_pages[i])
                }

        if len(draft_pages) != 0 {
                fmt.Println("The following pages are still in draft and were not rendered: ")
                for i := 0; i < len (draft_pages); i++ {
                        fmt.Println(draft_pages[i])
                }
        }
}

func (site Site) renderBlogs(blogs []string, pages []string) {
        
        // create menu items for the navbar
        // define posted and draft blogs
        
        menu, menuName := createMenu(pages)
        all_blogs, draft_blogs := defineBlogs(blogs)
        
        // RENDER BLOG POSTS
        //
        // Functionality
        // 
        // - list of all blogs
        // - pagination
        // - taxonomy
        // - author
        // - shortlist with x titles
        // - sorting on date
        //
        
        fmt.Println(menu, menuName)
        
        fmt.Println("The following blog posts were rendered: ")
        
        for i := 0; i < len(all_blogs); i++ {
                fmt.Println(all_blogs[i])
        }
        
        if len(draft_blogs) != 0 {
                fmt.Println("The following blog posts are still in draft and were not rendered: ")
                for i := 0; i < len (draft_blogs); i++ {
                        fmt.Println(draft_blogs[i])
                }
        }
}

func (site Site) renderGalleries(dirs []os.FileInfo, pages []string) {

        // create menu items for the navbar
        // list all gallery directories in temp
        // define all_galleries and all_galleries_name
        
        menu, menuName := createMenu(pages)
        all_galleries, all_galleries_name := defineGalleries(dirs)
        
        // create gallery.html content and sub-gallery htmls
        
        if _, err := os.Stat(site.pagedir+"/gallery.html"); os.IsNotExist(err) {
                copyfile(site.templatedir+"/gallery_template.html", site.pagedir+"/gallery.html")
                
                now := time.Now().Format(time.RFC1123)
                prepend("status          : posted\n------------------------------------------------------------------------", "pages/gallery.html")    
                prepend("menu name       : gallery", "pages/gallery.html")
                prepend("menu order      : 10", "pages/gallery.html")
                prepend("present in menu : y", "pages/gallery.html")
                prepend("------------------------------------------------------------------------\ncreated on      : "+now, "pages/gallery.html")
        }
        
        // Loop over all images and do the following updates
        //
        // GALLERY.HTML
        //
        // SUBGALLERYLINK  = pages/gallery/filename.jpg
        // SUBGALLERYTHUMB = pages/gallery/filename_thumb.jpg
        // SUBGALLERYNAME  = all_galleries_name[i]
        //
        // SUBGALLERY.HTML
        //
        // loop over images and add them to the subgallery
        // SUBIMAGE      = galleryname/imagename
        // SUBIMAGETHUMB = galleryname/imagename_thumb
        
        for i := 0; i < len(all_galleries); i++ {

                inject_html("temp/"+site.pagedir+"/gallery.html", "<<~~GALLERYITEM~~>>", "templates/gallery_item.html")

                imagepath := "temp/"+site.gallerydir+"/"+all_galleries_name[i]+"/"
                renderpath := "rendered/"+site.gallerydir+"/"+all_galleries_name[i]+"/"
                
                copydir(imagepath, renderpath)
                
                images, _ := filepath.Glob(imagepath+"*")
                for a := 0; a < len(images); a++ {
                        substitute("temp/pages/gallery.html","<<~~SUBGALLERYLINK~~>>","gallery/"+all_galleries_name[i]+".html")
                        substitute("temp/pages/gallery.html","<<~~SUBGALLERYTHUMB~~>>",strings.Split(strings.Split(images[a],"temp/pages/")[1],".")[0]+"_thumb.jpg")
                        substitute("temp/pages/gallery.html","<<~~SUBGALLERYNAME~~>>",all_galleries_name[i])
                        
                        inject_html(all_galleries[i], "<<~~SUBGALLERYITEM~~>>", site.templatedir+"/subgallery_item.html")
                        substitute(all_galleries[i],"<<~~SUBIMAGE~~>>", strings.Split(images[a],"temp/pages/gallery/")[1])
                        substitute(all_galleries[i],"<<~~SUBIMAGETHUMB~~>>", strings.Split(strings.Split(images[a],"temp/pages/gallery/")[1],".")[0]+"_thumb.jpg")

                        resize_picture(images[a], renderpath)
                }

                // Remove final trailing tie
                
                substitute(all_galleries[i],"<<~~SUBGALLERYITEM~~>>","")

                // Insert the Gallery name as title
                
                substitute(all_galleries[i], "<<~~GALLERYTITLE~~>>",all_galleries_name[i])
                
                // inject header and navbar
                
                inject_html(all_galleries[i], "<<~~HEADER~~>>", site.templatedir+"/header_template.html")
                inject_html(all_galleries[i], "<<~~NAVBAR~~>>", "temp/navbar.html")
                
                // populate navbar with the correct links

                create_navbar(all_galleries[i], menu, menuName, true)

                // populate the footer tie
                
                inject_html(all_galleries[i], "<<~~FOOTER~~>>", site.templatedir+"/footer_template.html")                
                
                // resolve ties CSS, JS, PAGE
                
                substitute(all_galleries[i], "<<~~JS~~>>","../../js/")
                substitute(all_galleries[i], "<<~~CSS~~>>","../../css/")
                substitute(all_galleries[i], "<<~~PAGE~~>>","../")
                
                //remove_header(all_galleries[i])
                
                copyfile(all_galleries[i], "rendered/"+site.gallerydir+"/"+strings.Split(all_galleries[i],"/")[3])
                copydir("temp/"+site.gallerydir+"/"+all_galleries_name[i], "rendered/pages/gallery/"+all_galleries_name[i])
        }
        substitute("temp/"+site.pagedir+"/gallery.html","<<~~GALLERYITEM~~>>","")
        copyfile("temp/"+site.pagedir+"/gallery.html", "rendered/pages/gallery.html")
}


// define functions

func resize_picture (filename, output_folder string) {
        
        file_temp :=  strings.Split(filename, "/")
        img_name := file_temp[len(file_temp)-1]
        
        output_file := output_folder+img_name
        output_thumb := output_folder+strings.Split(img_name, ".")[0]+"_thumb.jpg"
        
        file, err := os.Open(filename)
        if err != nil {
                log.Fatal(err)
        }
        img, err := jpeg.Decode(file)
        if err != nil {
                log.Fatal(err)
        }
        file.Close()

        b := img.Bounds()
        imgWidth := b.Max.X
        //imgHeight := b.Max.Y    
        
        if imgWidth > 1000 {
                // resize to width 1000 using Lanczos resampling
                // and preserve aspect ratio
                m := resize.Resize(1000, 0, img, resize.Lanczos3)
        
                out, err := os.Create(output_file)
                if err != nil {
                        log.Fatal(err)
                }
                defer out.Close()
                
                jpeg.Encode(out, m, nil)
        }
        
        // Create the thumbnails
        
        n := resize.Resize(350, 0, img, resize.Lanczos3)
        
        out2, err := os.Create(output_thumb)
        if err != nil {
                log.Fatal(err)
        }
        defer out2.Close()
                
        // write resized image and thumbnails to file
        
        jpeg.Encode(out2, n, nil)
}

func prepend(text, file string){
        input, err := ioutil.ReadFile(file)
        if err != nil {
                log.Fatalln(err)
        }
        
        lines := strings.Split(string(input), "\n")
        new_lines := append([]string{text}, lines...)

        output := strings.Join(new_lines, "\n")
        err = ioutil.WriteFile(file, []byte(output), 0644)
        if err != nil {
                log.Fatalln(err)
        }    
}

func inject_html (file, tie, html_source_file string) {
        input, err := ioutil.ReadFile(file)
        if err != nil {
                log.Fatalln(err)
        }

        lines := strings.Split(string(input), "\n")
        
        html_input, err := ioutil.ReadFile(html_source_file)
        if err != nil {
                log.Fatalln(err)
        }
        
        s := string(html_input)
        
        for line := range lines {
                lines[line] = strings.Replace(lines[line], tie, s, -1)
        }
        
        output := strings.Join(lines, "\n")
        err = ioutil.WriteFile(file, []byte(output), 0644)
        if err != nil {
                log.Fatalln(err)
        } 
}

func create_navbar (page string, menu map[int64]string, menuName map[int64]string, galleryYN bool) {
        for j := 0; j < 10; j++ {
                if orig_link, ok := menu[int64(j)]; ok {        
                        var page_name string = menuName[int64(j)]
                        page_link := strings.Split(orig_link,"/")[len(strings.Split(orig_link,"/"))-1]
                        inject_nav_items(page, "<<~~NAVLIST~~>>", site.templatedir+"/navbar_item.html")
                        if page_link == strings.Split(page,"/")[2] {
                                substitute(page,"<<~~NAVACTIVE~~>>", "class=\"active\"")
                        } else {
                                substitute(page,"<<~~NAVACTIVE~~>>", "")
                        }
                        
                        if galleryYN == true {
                                if page_link == "index.html" {
                                        substitute(page, "<<~~NAVLINK~~>>","../../"+page_link)
                                } else {
                                        substitute(page, "<<~~NAVLINK~~>>","../"+page_link)
                                }
                        } else if strings.Split(page,"/")[2] == "index.html" {
                                if page_link == "index.html" {
                                        substitute(page, "<<~~NAVLINK~~>>",page_link)
                                } else {
                                        substitute(page, "<<~~NAVLINK~~>>",site.pagedir+"/"+page_link)
                                }
                        } else {
                                if page_link == "index.html" {
                                        substitute(page, "<<~~NAVLINK~~>>","../"+page_link)
                                } else {
                                        substitute(page, "<<~~NAVLINK~~>>",page_link)
                                }
                        }
                        substitute(page, "<<~~NAVITEM~~>>",page_name)
                }
        }
        substitute(page, "<<~~NAVLIST~~>>","")        
}

func inject_nav_items (file, tie, html_source_file string) {
        input, err := ioutil.ReadFile(file)
        if err != nil {
                log.Fatalln(err)
        }

        lines := strings.Split(string(input), "\n")
        
        html_input, err := ioutil.ReadFile(html_source_file)
        if err != nil {
                log.Fatalln(err)
        }
        
        s := string(html_input)
        
        for line := range lines {
                lines[line] = strings.Replace(lines[line], tie, s+"\n<<~~NAVLIST~~>>", -1)
                /*lines[line] = strings.Replace(lines[line], tie, s, -1)*/
        }
        
        output := strings.Join(lines, "\n")
        err = ioutil.WriteFile(file, []byte(output), 0644)
        if err != nil {
                log.Fatalln(err)
        } 
}

func remove_header (file string) {
        input, err := ioutil.ReadFile(file)
        if err != nil {
                log.Fatalln(err)
        }
        
        lines := strings.Split(string(input), "\n")
        i := 0
        for i <= 6 {
                lines = append(lines[:0], lines[1:]...)
                // lines[i] = strings.Replace(lines[line], tie, replacetext, -1)
                i += 1
        }
        
        output := strings.Join(lines, "\n")
        err = ioutil.WriteFile(file, []byte(output), 0644)
        if err != nil {
                log.Fatalln(err)
        }
}

func createMenu (fileList []string) (map[int64]string, map[int64]string) {

        // creates the values for the navbar based on a list of stare html input files and the information in the stare headers

        // declare variables
        
        menuItem := []string{}
        menu := make(map[int64]string)
        menuName := make(map[int64]string)
        
        // Read pages

        for i := 0; i < len(fileList); i++ {
        
                input, err := ioutil.ReadFile(fileList[i])
                if err != nil {
                        log.Fatalln(err)
                }
        
                lines := strings.Split(string(input), "\n")
                
                for j := 1; j < 6; j++  {
                        if strings.Contains(lines[j], "posted") == true {
                                for k := 1; k < 6; k++ {
                                        if strings.Contains(lines[k], "present in menu : y") == true {
                                                menuOrder, _ := strconv.ParseInt(strings.Split(lines[3], ": ", )[1], 0, 64)
                                                menuItem = append(menuItem, strings.Split(fileList[i], "/")[2])
                                                menu[menuOrder] = fileList[i]
                                        }
                                        if strings.Contains(lines[k], "menu name       : ") == true {
                                                menuOrder, _ := strconv.ParseInt(strings.Split(lines[3], ": ", )[1], 0, 64)
                                                menuName[menuOrder] = strings.Split(lines[k], ": ")[1]
                                        }
                                }
                        } 
                }
        }
        return menu, menuName
}

func definePages (fileList []string) ([]string, []string) {

        // declare variables

        all_pages := []string{}
        draft_pages := []string{}

        // complete all_pages with all posted pages and draft_pages with all in draft pages
        // check whether the page is posted
        // check whether the page should be present in the menu
        
        for i := 0; i < len(fileList); i++ {
        
                input, err := ioutil.ReadFile(fileList[i])
                if err != nil {
                        log.Fatalln(err)
                }
        
                lines := strings.Split(string(input), "\n")
                
                for j := 1; j < 6; j++  {
                        if strings.Contains(lines[j], "posted") == true {
                                all_pages = append(all_pages, fileList[i])
                        } 
                        if strings.Contains(lines[j], "in_draft") == true {
                                draft_pages = append(draft_pages, fileList[i])
                        }
                }
        }
        return all_pages, draft_pages
}

func defineBlogs (fileList []string) ([]string, []string) {
        
        all_blogs := []string{}
        draft_blogs := []string{}
        // all_blogs_taxonomy := make(map[string]string)
        
                // complete all_blogs with all posted blogs and draft_blogs with all in draft blogs
        // check whether the page is in draft
        
        for i := 0; i < len(fileList); i++ {
        
                input, err := ioutil.ReadFile(fileList[i])
                if err != nil {
                        log.Fatalln(err)
                }
        
                lines := strings.Split(string(input), "\n")
                
                for j := 1; j < 6; j++  {
                        if strings.Contains(lines[j], "posted") == true {
                                all_blogs = append(all_blogs, fileList[i])
                        } 
                        if strings.Contains(lines[j], "in_draft") == true {
                                draft_blogs = append(draft_blogs, fileList[i])
                        }
                }
        }
        
        return all_blogs, draft_blogs
}

func defineGalleries (galleryDirs []os.FileInfo) ([]string, []string) {
        
        // declare variables
       
        all_galleries := []string{}
        all_galleries_name := []string{}
        
        // create all_galleries and all_galleries_name
        
        for i := 0; i < len(galleryDirs); i++ {
                if galleryDirs[i].IsDir() == true {
                        copyfile("." + string(filepath.Separator) + site.templatedir + string(filepath.Separator) + "subgallery_template.html", "." + string(filepath.Separator) + "temp" + string(filepath.Separator) + "pages" + string(filepath.Separator) + "gallery" + string(filepath.Separator) + galleryDirs[i].Name() + ".html")
                        all_galleries = append(all_galleries, "temp/pages/gallery/"+galleryDirs[i].Name()+".html")
                        all_galleries_name = append(all_galleries_name, galleryDirs[i].Name())
                }
        }
        
        return all_galleries, all_galleries_name
}

func render_site() {

        fmt.Println("Rendering!")
        site.createFolder()
        site.copySrc()
        
        if site.multiLang == false {
                site.copyFiles()
        } else {
                site.copyFilesLangSupp("nl")
        }
        // list all html pages in temp
        // list all blog posts in temp
        // list all gallery dirs in temp
        
        pages, _ := filepath.Glob("temp/"+site.pagedir+"/*.html")
        blogs, _ := filepath.Glob("temp/"+site.blogdir+"/*.md")
        dirs, _ := ioutil.ReadDir ("temp/"+site.gallerydir+"/")
        
        site.renderPages(pages, "nl", "nederlands")
        site.renderBlogs(blogs, pages)
        site.renderGalleries(dirs, pages)
        
        // Remove the temporary files 
        
        os.RemoveAll(path.Join(wd, "temp"))
}