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


// TODO 

// rewrite remove_header so it cuts on the second line consisting of ------

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
        "github.com/russross/blackfriday"
        "image/jpeg"
        "strconv"
        //"time"
        "path"
        "github.com/oliamb/cutter"
)

var wd, _ = os.Getwd()
var pagesAllLang []string

// define Site functions

func (site Site) createFolder () {
        
        // Remove previous rendering of the site

        RemoveContentsLeaveGit(path.Join(wd, "rendered"))
        
        
        // Create directories for temporary files and newly rendered site
        
        if site.multiLang == true {
                
                os.MkdirAll("temp/"+site.gallerydir, 0755)
                os.MkdirAll("rendered/"+site.gallerydir, 0755)
                
                for i := 0; i < len(site.languages); i++ {
                        os.MkdirAll("temp/"+site.pagedir+"/"+site.languages[i], 0755)
                        os.MkdirAll("temp/"+site.pagedir+"/"+site.languages[i]+"/"+site.blogdir, 0755)
                        os.MkdirAll("rendered/"+site.pagedir+"/"+site.languages[i], 0755)
                        os.MkdirAll("rendered/"+site.blogdir+"/"+site.languages[i], 0755)
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
        for i := 0; i < len(srcItems); i++ {
                file, err := os.Open(srcItems[i])
                if err != nil {
                    // handle the error and return
                }
                defer file.Close()

                fi, err := file.Stat()
                if err != nil {
                    // handle the error and return
                }
                if fi.IsDir() {
                    // it's a directory
                    copydir(srcItems[i], "rendered/"+strings.Split(srcItems[i], "/")[1])
                } else {
                    // it's not a directory
                    copyfile(srcItems[i], "rendered/"+strings.Split(srcItems[i], "/")[1])
                    //fmt.Println("Reached this point - ", srcItems[i])
                }
        }
}

func (site Site) copyFiles () {

        // move the pages, blogs and galleries to the temporary directory
        if site.multiLang == true {
                for i := 0; i < len(site.languages); i++ {
                        copydir(site.pagedir+"/"+site.languages[i], "temp/"+site.pagedir+"/"+site.languages[i])
                        copydir(site.pagedir+"/"+site.languages[i]+"/"+site.blogdir, "temp/"+site.pagedir+"/"+site.languages[i]+"/"+site.blogdir)
                        copydir(site.gallerydir, "temp/"+site.pagedir+"/"+site.gallerydir)
                }
        } else {
                copydir(site.pagedir, "temp/"+site.pagedir)
                copydir(site.pagedir+"/"+site.blogdir, "temp/"+site.pagedir+"/"+site.blogdir)
                copydir(site.gallerydir, "temp/"+site.pagedir+"/"+site.gallerydir) 
        }
        // copy the navbar template to the temp folder 

        copyfile(site.templatedir+"/navbar_template.html", "temp/navbar.html")
}

func (site Site) renderPages(pages []string, blogs []string, language string) {

        // complete menu and menuName
        // define posted and draft pages

        menu, menuName := createMenu(pages, language)
        all_pages, draft_pages := definePages(pages)

        posted_blogs, _ := defineBlogs(blogs)
        author, title, time, filename, _, revTaxonomy := dissectBlogs(posted_blogs)
        
        //fmt.Println(taxonomy)
        //fmt.Println(revTaxonomy)
        //fmt.Println(title)
        
        // add navbar to the pages and resolve the ties NAVACTIVE, NAVLINK, NAVITEM
        // cycling through all posted pages, then cycling through all menu items
        // adding navlinks as necessary and resolving the ties
        
        for i := 0; i < len(all_pages); i++ {
                
                inject_html(all_pages[i], "<<~~NAVBAR~~>>", "temp/navbar.html","")
                
                create_navbar(all_pages[i], language, menu, menuName, false)

                // populate the header and footer tie
                
                inject_html(all_pages[i], "<<~~HEADER~~>>", site.templatedir+"/"+language+"/header_template.html","")
                inject_html(all_pages[i], "<<~~FOOTER~~>>", site.templatedir+"/"+language+"/footer_template.html","")
                
                // populate the blogpost items
                
                // insert blogtotal function to get # of blogs; then use this number to rebuild 2nd field in inject_html (blogcounter - blogtotal)!
                
                var originalBlogTotal = totalBlogs(all_pages[i])
                var extraCount = 0
                var blogTotal = originalBlogTotal

                if blogTotal > len(filename) {
                        blogTotal = len(filename)
                }
                
                for blogCounter := blogTotal; blogCounter > 0; blogCounter-- {
                        
                        number := originalBlogTotal-extraCount
                        newNumber := originalBlogTotal-extraCount-1
                        
                        if blogTotal-extraCount == 1 {
                                inject_html(all_pages[i], "<<~~BLOG:POSTS:"+strconv.Itoa(number)+"~~>>", site.templatedir+"/"+language+"/blogpost_template.html", "")
                                inject_html(all_pages[i], "<<~~BLOG:TITLES:"+strconv.Itoa(number)+"~~>>", site.templatedir+"/"+language+"/blogtitles_template.html", "")
                        } else {
                                inject_html(all_pages[i], "<<~~BLOG:POSTS:"+strconv.Itoa(number)+"~~>>", site.templatedir+"/"+language+"/blogpost_template.html", "<<~~BLOG:POSTS:"+strconv.Itoa(newNumber)+"~~>>")
                                inject_html(all_pages[i], "<<~~BLOG:TITLES:"+strconv.Itoa(number)+"~~>>", site.templatedir+"/"+language+"/blogtitles_template.html", "<<~~BLOG:TITLES:"+strconv.Itoa(newNumber)+"~~>>")
                        }
                        substitute(all_pages[i], "<<~~BLOGPOST:TITLE~~>>",title[extraCount])
                        substitute(all_pages[i], "<<~~BLOGPOST:AUTHOR~~>>",author[extraCount])
                        substitute(all_pages[i], "<<~~BLOGPOST:TIME~~>>",time[extraCount])
                        substitute(all_pages[i], "<<~~BLOGPOST:TAXONOMY~~>>",revTaxonomy[extraCount])
                        
                        renderTaxLink(revTaxonomy[extraCount])
                        
                        blogContent := loadBlogContent(filename[extraCount])
                        substitute(all_pages[i], "<<~~BLOGPOST:CONTENT~~>>",blogContent)
                        
                        extraCount++
                }
                
                // resolve ties CSS, JS, PAGE, IMG
                
                if site.multiLang == true && language == site.primaryLang {
                        if strings.Split(all_pages[i],"/")[len(strings.Split(all_pages[i],"/"))-1] == "index.html" {
                                substitute(all_pages[i], "<<~~JS~~>>","js/")
                                substitute(all_pages[i], "<<~~CSS~~>>","css/")
                                substitute(all_pages[i], "<<~~IMG~~>>","img/")
                                substitute(all_pages[i], "<<~~PAGE~~>>","pages/")
                        } else {
                                substitute(all_pages[i], "<<~~JS~~>>","../../js/")
                                substitute(all_pages[i], "<<~~CSS~~>>","../../css/")
                                substitute(all_pages[i], "<<~~IMG~~>>","../../img/")
                                substitute(all_pages[i], "<<~~PAGE~~>>","")
                        }
                } else if site.multiLang == true && language != site.primaryLang {
                        substitute(all_pages[i], "<<~~JS~~>>","../../js/")
                        substitute(all_pages[i], "<<~~CSS~~>>","../../css/")
                        substitute(all_pages[i], "<<~~IMG~~>>","../../img/")
                        substitute(all_pages[i], "<<~~PAGE~~>>","")
                } else {
                        if strings.Split(all_pages[i],"/")[len(strings.Split(all_pages[i],"/"))-1] == "index.html" {
                                substitute(all_pages[i], "<<~~JS~~>>","js/")
                                substitute(all_pages[i], "<<~~CSS~~>>","css/")
                                substitute(all_pages[i], "<<~~IMG~~>>","img/")
                                substitute(all_pages[i], "<<~~PAGE~~>>","pages/")
                        } else {
                                substitute(all_pages[i], "<<~~JS~~>>","../js/")
                                substitute(all_pages[i], "<<~~CSS~~>>","../css/")
                                substitute(all_pages[i], "<<~~IMG~~>>","../img/")
                                substitute(all_pages[i], "<<~~PAGE~~>>","")
                        }
                }

        }

        // give an overview of rendered pages and blogs and draft pages and blogs
        if len(draft_pages) != 0 {
                fmt.Println ("Draft pages were not rendered.")
        }
        // fmt.Println("The following pages were rendered: ")
        
        // for i := 0; i < len (all_pages); i++ {
        //         fmt.Println(all_pages[i])
        //         }

        // if len(draft_pages) != 0 {
        //         fmt.Println("The following pages are still in draft and were not rendered: ")
        //         for i := 0; i < len (draft_pages); i++ {
        //                 fmt.Println(draft_pages[i])
        //         }
        // }
}

func (site Site) addLangLinks (pages []string, language string) {
        
        all_pages, _ := definePages(pages)

        for i := 0; i < len(all_pages); i++ {

                pagename := strings.Split(all_pages[i],"/")[len(strings.Split(all_pages[i],"/"))-1]
                
                for j := 0; j < len(site.languages); j++ {
                        if language != site.languages[j] {
                                
                                inject_nav_items(all_pages[i], "<<~~LANGLIST~~>>", site.templatedir+"/"+language+"/langbar_item.html")
                                
                                if Contains(pagesAllLang, "temp/pages/"+site.languages[j]+"/"+pagename) == true && pagename == "index.html" && language != site.primaryLang && site.languages[j] == site.primaryLang {
                                        substitute(all_pages[i], "<<~~LANGLINK~~>>", "../../"+pagename)
                                        substitute(all_pages[i], "<<~~LANGITEM~~>>", site.languages[j])
                                } else if Contains(pagesAllLang, "temp/pages/"+site.languages[j]+"/"+pagename) == true && pagename == "index.html" && language == site.primaryLang && site.languages[j] != site.primaryLang {
                                        substitute(all_pages[i], "<<~~LANGLINK~~>>", site.pagedir+"/"+site.languages[j]+"/"+pagename)
                                        substitute(all_pages[i], "<<~~LANGITEM~~>>", site.languages[j])
                                } else if Contains(pagesAllLang, "temp/pages/"+site.languages[j]+"/"+pagename) == true {
                                        substitute(all_pages[i], "<<~~LANGLINK~~>>", "../"+site.languages[j]+"/"+pagename)
                                        substitute(all_pages[i], "<<~~LANGITEM~~>>", site.languages[j])
                                } else {
                                        substitute(all_pages[i], "<<~~LANGLINK~~>>", "../"+site.languages[j]+"/index.html")
                                        substitute(all_pages[i], "<<~~LANGITEM~~>>", site.languages[j])
                                }
                        }
                }
                
                substitute(all_pages[i], "<<~~LANGLIST~~>>","")
        }
}

func (site Site) copyRenderedPages(pages []string, language string) {
        
        all_pages, _ := definePages(pages)

        // copy files
        
        for i := 0; i < len(all_pages); i++ {
                
                // Remove the stare header before copying to the rendered folder
        
                remove_header(all_pages[i])
                
                // Copy files to the correct location
        
                if (site.multiLang == true && language == site.primaryLang) || site.multiLang == false {
                        if strings.Split(all_pages[i],"/")[len(strings.Split(all_pages[i],"/"))-1] == "index.html" {
                                copyfile(all_pages[i], "rendered/"+strings.Split(all_pages[i],"/")[len(strings.Split(all_pages[i],"/"))-1])
                        } else {
                                copyfile(all_pages[i], "rendered/"+site.pagedir+"/"+language+"/"+strings.Split(all_pages[i],"/")[len(strings.Split(all_pages[i],"/"))-1])
                        }
                } else {
                        copyfile(all_pages[i], "rendered/"+site.pagedir+"/"+language+"/"+strings.Split(all_pages[i],"/")[len(strings.Split(all_pages[i],"/"))-1])
                }
        }
}

func (site Site) renderGalleries(dirs []os.FileInfo, pages []string, language string) {

        // create menu items for the navbar
        // list all gallery directories in temp
        // define all_galleries and all_galleries_name
        
        menu, menuName := createMenu(pages, language)

        all_galleries, all_galleries_name := defineGalleries(dirs)
        
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

                inject_html("temp/"+site.pagedir+"/gallery.html", "<<~~GALLERYITEM~~>>", "templates/gallery_item.html", "")

                imagepath := "temp/"+site.pagedir+"/"+site.gallerydir+"/"+all_galleries_name[i]+"/"
                renderpath := "rendered/"+site.pagedir+"/"+site.gallerydir+"/"+all_galleries_name[i]+"/"
                
                copydir(imagepath, renderpath)
                
                images, _ := filepath.Glob(imagepath+"*")
                for a := 0; a < len(images); a++ {
                        substitute("temp/pages/gallery.html","<<~~SUBGALLERYLINK~~>>","gallery/"+all_galleries_name[i]+".html")
                        substitute("temp/pages/gallery.html","<<~~SUBGALLERYTHUMB~~>>",strings.Split(strings.Split(images[a],"temp/pages/")[1],".")[0]+"_thumb.jpg")
                        substitute("temp/pages/gallery.html","<<~~SUBGALLERYNAME~~>>",all_galleries_name[i])
                        
                        inject_html(all_galleries[i], "<<~~SUBGALLERYITEM~~>>", site.templatedir+"/subgallery_item.html", "")
                        substitute(all_galleries[i],"<<~~SUBIMAGE~~>>", strings.Split(images[a],"temp/pages/gallery/")[1])
                        substitute(all_galleries[i],"<<~~SUBIMAGETHUMB~~>>", strings.Split(strings.Split(images[a],"temp/pages/gallery/")[1],".")[0]+"_thumb.jpg")
                        substitute(all_galleries[i],"<<~~SUBGALLERYNAME~~>>",all_galleries_name[i])

                        resize_picture(images[a], renderpath)
                }

                // Remove final trailing tie
                
                substitute(all_galleries[i],"<<~~SUBGALLERYITEM~~>>","")

                // Insert the Gallery name as title
                
                substitute(all_galleries[i], "<<~~GALLERYTITLE~~>>",all_galleries_name[i])
                
                // inject header and navbar
                
                inject_html(all_galleries[i], "<<~~HEADER~~>>", site.templatedir+"/header_template.html", "")
                inject_html(all_galleries[i], "<<~~NAVBAR~~>>", "temp/navbar.html", "")
                
                // populate navbar with the correct links

                create_navbar(all_galleries[i], language, menu, menuName, true)

                // populate the footer tie
                
                inject_html(all_galleries[i], "<<~~FOOTER~~>>", site.templatedir+"/footer_template.html", "")                
                
                // resolve ties CSS, JS, PAGE
                
                substitute(all_galleries[i], "<<~~JS~~>>","../../js/")
                substitute(all_galleries[i], "<<~~CSS~~>>","../../css/")
                substitute(all_galleries[i], "<<~~PAGE~~>>","../")
                
                //remove_header(all_galleries[i])
                
                copyfile(all_galleries[i], "rendered/"+site.pagedir+"/"+site.gallerydir+"/"+strings.Split(all_galleries[i],"/")[3])
                copydir("temp/"+site.pagedir+"/"+site.gallerydir+"/"+all_galleries_name[i], "rendered/pages/gallery/"+all_galleries_name[i])
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
        imgHeight := b.Max.Y    
        
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
        n := img
        
        switch {
                case imgHeight <= 275:
                        n, err = cutter.Crop(img, cutter.Config{
                                Width: 350,
                                Height: 275,
                                Mode: cutter.Centered, 
                        })
                default:
                        tempn := resize.Resize(350, 0, img, resize.Lanczos3)
                        n, err = cutter.Crop(tempn, cutter.Config{
                                Width: 350,
                                Height: 275,
                                Mode: cutter.Centered,
                        })
        }
        
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

func inject_html (file, tie, html_source_file string, footer string) {
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
                if strings.Contains(lines[line], "<<~~BLOG:POSTS:") == true && len(footer) > 0 {
                        lines[line] = lines[line] +"\n"+ footer
                }
                lines[line] = strings.Replace(lines[line], tie, s, -1)
        }
        
        output := strings.Join(lines, "\n")
        err = ioutil.WriteFile(file, []byte(output), 0644)
        if err != nil {
                log.Fatalln(err)
        } 
}

func totalBlogs (file string) int {
        
        var a int
        
        input, err := ioutil.ReadFile(file)
        if err != nil {
                log.Fatalln(err)
        }

        lines := strings.Split(string(input), "\n")

        for line := range lines {
                if strings.Contains(lines[line],"<<~~BLOG:POSTS:") == true {
                        a, _ = strconv.Atoi(strings.Split(strings.Split(lines[line], "<<~~BLOG:POSTS:")[1], "~~>>")[0])
                        
                }
                
        }
        return a
}

func loadBlogContent (file string) string {
        input, err := ioutil.ReadFile(file)
        if err != nil {
                log.Fatalln(err)
        }
        
        lines := strings.Split(string(input), "\n")
        
        for i := 0; i <= 6; i++ {
                lines = append(lines[:0], lines[1:]...)
                
        }
        
        output := strings.Join(lines, "\n")
        
        html := blackfriday.MarkdownCommon([]byte(output))
        
        return string(html)
        
}

func create_navbar (page string, language string, menu map[int64]string, menuName map[int64]string, galleryYN bool) {
        for j := 0; j < 10; j++ {
                if orig_link, ok := menu[int64(j)]; ok {        
                        
                        var page_name string = menuName[int64(j)]
                        page_link := strings.Split(orig_link,"/")[len(strings.Split(orig_link,"/"))-1]
                        
                        inject_nav_items(page, "<<~~NAVLIST~~>>", site.templatedir+"/"+language+"/navbar_item.html")
                        
                        if page_link == strings.Split(page,"/")[len(strings.Split(page,"/"))-1] {
                                substitute(page,"<<~~NAVACTIVE~~>>", "class=\"active\"")
                        } else {
                                substitute(page,"<<~~NAVACTIVE~~>>", "")
                        }
                        
                        if site.multiLang == false {
                                if galleryYN == true {
                                        if page_link == "index.html" {
                                                substitute(page, "<<~~NAVLINK~~>>","../../"+page_link)
                                        } else {
                                                substitute(page, "<<~~NAVLINK~~>>","../"+page_link)
                                        }
                                } else if strings.Split(page,"/")[len(strings.Split(page,"/"))-1] == "index.html" {
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
                        } else if (site.multiLang == true && language == site.primaryLang) {
                                if galleryYN == true {
                                        if page_link == "index.html" {
                                                substitute(page, "<<~~NAVLINK~~>>","../../"+page_link)
                                        } else {
                                                substitute(page, "<<~~NAVLINK~~>>","../"+page_link)
                                        }
                                } else if strings.Split(page,"/")[len(strings.Split(page,"/"))-1] == "index.html" {
                                        if page_link == "index.html" {
                                                substitute(page, "<<~~NAVLINK~~>>",page_link)
                                        } else {
                                                substitute(page, "<<~~NAVLINK~~>>",site.pagedir+"/"+language+"/"+page_link)
                                        }
                                } else {
                                        if page_link == "index.html" {
                                                substitute(page, "<<~~NAVLINK~~>>","../../"+page_link)
                                        } else {
                                                substitute(page, "<<~~NAVLINK~~>>",page_link)
                                        }
                                }
                        } else if (site.multiLang == true && language != site.primaryLang) {
                                if galleryYN == true {
                                        if page_link == "index.html" {
                                                substitute(page, "<<~~NAVLINK~~>>","../../"+page_link)
                                        } else {
                                                substitute(page, "<<~~NAVLINK~~>>","../"+page_link)
                                        }
                                // } else if strings.Split(page,"/")[len(strings.Split(page,"/"))-1] == "index.html" {
                                //         if page_link == "index.html" {
                                //                 substitute(page, "<<~~NAVLINK~~>>",page_link)
                                //         } else {
                                //                 substitute(page, "<<~~NAVLINK~~>>",page_link)
                                //         }
                                } else {
                                        //if page_link == "index.html" {
                                        //        substitute(page, "<<~~NAVLINK~~>>",page_link)
                                        //} else {
                                                substitute(page, "<<~~NAVLINK~~>>",page_link)
                                        //}
                                }
                        } else {
                                if galleryYN == true {
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
                lines[line] = strings.Replace(lines[line], tie, s+"\n"+tie, -1)
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
        
        for i := 0; i <= 6; i++ {
                lines = append(lines[:0], lines[1:]...)
                
        }


        output := strings.Join(lines, "\n")
        err = ioutil.WriteFile(file, []byte(output), 0644)
        if err != nil {
                log.Fatalln(err)
        }
}

func createMenu (fileList []string, language string) (map[int64]string, map[int64]string) {

        // creates the values for the navbar based on a list of stare html input files and the information in the stare headers

        // declare variables
        
        menuItem := []string{}
        menu := make(map[int64]string)
        menuName := make(map[int64]string)
        
        // Read pages

        for i := 0; i < len(fileList); i++ {
        
                if site.gallery == false && fileList[i] == "temp/pages/gallery.html" {
                        
                } else {
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
                if site.gallery == false && fileList[i] == "temp/pages/gallery.html" {
                        
                } else {
                        input, err := ioutil.ReadFile(fileList[i])
                        if err != nil {
                                log.Fatalln(err)
                        }
        
                        lines := strings.Split(string(input), "\n")
                        
                        for j := 1; j < 6; j++  {
                                if strings.Contains(lines[j], "posted") == true {
                                        all_pages = append(all_pages, fileList[i])
                                        pagesAllLang = append(pagesAllLang, fileList[i])
                                } 
                                if strings.Contains(lines[j], "in_draft") == true {
                                        draft_pages = append(draft_pages, fileList[i])
                                }
                        }
                }
        }
        return all_pages, draft_pages
}

func defineBlogs (fileList []string) ([]string, []string) {
        
        posted_blogs := []string{}
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
                                posted_blogs = append(posted_blogs, fileList[i])
                        } 
                        if strings.Contains(lines[j], "in_draft") == true {
                                draft_blogs = append(draft_blogs, fileList[i])
                        }
                }
        }
        
        return posted_blogs, draft_blogs
}

func dissectBlogs (posted_blogs []string) (map[int]string, map[int]string, map[int]string, map[int]string, map[string]string, map[int]string) {
        
        author := make(map[int]string)
        title := make(map[int]string)
        time := make(map[int]string)
        filename := make(map[int]string)
        taxonomy := make(map[string]string)
        revTaxonomy := make(map[int]string)
        
        if _, err := os.Stat(path.Join(site.pagedir,site.blogdir,"taxonomies")); os.IsNotExist(err) {
                                        blogdir, _ := os.Stat(filepath.Join(site.pagedir,site.blogdir))
                                        err = os.MkdirAll(path.Join("temp",site.pagedir,site.blogdir,"taxonomies"), blogdir.Mode())
                                        // if err != nil {
                                        //         return err
                                        // }
                                }
        
        // Read pages

        for i := 0; i < len(posted_blogs); i++ {
        
                input, err := ioutil.ReadFile(posted_blogs[i])
                if err != nil {
                        log.Fatalln(err)
                }
        
                lines := strings.Split(string(input), "\n")
                
                filename[i] = posted_blogs[i]
                
                for k := 1; k < 6; k++ {
                        if strings.Contains(lines[k], "created by      : ") == true {
                                author[i] = strings.Split(lines[k], ": ")[1]
                        }
                        if strings.Contains(lines[k], "title           : ") == true {
                                title[i] = strings.Split(lines[k], ": ")[1]
                        }
                        if strings.Contains(lines[k], "created on      : ") == true {
                                time[i] = strings.Split(lines[k], ": ")[1]
                        }
                        if strings.Contains(lines[k], "taxonomies      : ") == true {
                                numTax := len(strings.Split(strings.Split(lines[k], ": ")[1], ";"))
                                Tx := strings.Split(strings.Split(lines[k], ": ")[1], ";")
                                
                                // create taxonomy
                                
                                for l := 0; l < numTax; l++ {
                                        if val, ok := taxonomy[Tx[l]]; ok {
                                            taxonomy[Tx[l]] = val + ";" + strconv.Itoa(i)
                                        } else {
                                                taxonomy[Tx[l]] = strconv.Itoa(i)
                                        }
                                }
                                
                                
                                
                                // create revTaxonomy
                                
                                for l := 0; l < numTax; l++ {
                                        if val, ok := revTaxonomy[i]; ok {
                                            revTaxonomy[i] = val + ";" + Tx[l]
                                        } else {
                                            revTaxonomy[i] = Tx[l]
                                        }
                                }
                                
                        }
                }
        }
        return author, title, time, filename, taxonomy, revTaxonomy
}

func renderTaxLink (taxlist string) string {
        
        var link string
        
        for i := 1; i <= len(strings.Split(taxlist, ";")); i++ {
                link = link + strings.Split(taxlist, ";")[i-1] + "  "
        
        }
        
        
        
        fmt.Println(link)
        return link
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

func Contains(list []string, elem string) bool { 
        for _, t := range list { if t == elem { return true } } 
        return false 
}

func RemoveContentsLeaveGit(dir string) error {
    d, err := os.Open(dir)
    if err != nil {
        return err
    }
    defer d.Close()
    names, err := d.Readdirnames(-1)
    if err != nil {
        return err
    }
    for _, name := range names {
        if name != ".git" {
                err = os.RemoveAll(filepath.Join(dir, name))
        }
        if err != nil {
            return err
        }
    }
    return nil
}

func render_site() {

        fmt.Println("Rendering!")
        site.createFolder()
        site.copySrc()
        site.copyFiles()
        
        if site.multiLang == false {
                pages, _ := filepath.Glob("temp/"+site.pagedir+"/*.html")
                blogs, _ := filepath.Glob("temp/"+site.pagedir+"/"+site.blogdir+"/*.md")
                dirs, _ := ioutil.ReadDir ("temp/"+site.pagedir+"/"+site.gallerydir+"/")
                
                if site.gallery == true {
                        site.renderGalleries(dirs, pages, "")
                }
                
                site.renderPages(pages, blogs, "")
                site.copyRenderedPages(pages, "")
                
        } else {
                for i := 0; i < len(site.languages); i++ {
                        
                        fmt.Println("Multilanguage site - rendering language ", site.languages[i])
                        
                        pages, _ := filepath.Glob("temp/"+site.pagedir+"/"+site.languages[i]+"/*.html")
                        blogs, _ := filepath.Glob("temp/"+site.pagedir+"/"+site.languages[i]+"/"+site.blogdir+"/*.md")
                        dirs, _ := ioutil.ReadDir ("temp/"+site.pagedir+"/"+site.gallerydir+"/")
                        if site.gallery == true {
                                site.renderGalleries(dirs, pages, site.languages[i])
                        }
                        
                        site.renderPages(pages, blogs, site.languages[i])

                }
                
                for i := 0; i < len(site.languages); i++ {
                        pages, _ := filepath.Glob("temp/"+site.pagedir+"/"+site.languages[i]+"/*.html")
                        site.addLangLinks(pages, site.languages[i])
                        site.copyRenderedPages(pages, site.languages[i])
                }
        }

        // Remove the temporary files 
        
        os.RemoveAll(path.Join(wd, "temp"))
}