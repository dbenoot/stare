/* TODO

- delete the current rendered folder
- resize pictures and create thumbnails (x_thumb.jpg)
- create gallery.html and the separate gallery pages, based on templates (gallery_template.html and subgallery_template.html)
- move all galleries from dist/gallery to dist/rendered/gallery
- move all pages from dist/pages to dist/rendered/pages, move index.html to dist/rendered, create gallery.html in dist/pages
- move all css/js/... from dist/src to dist/rendered/
- replace all placeholders with correct html (<<~~INDEX~~>>, <<~~PAGES~~>>, <<~~TITLE~~>>, <<~~NAVBAR~~>>, <<~~NAVLIST~~>>, <<~~FOOTER~~>>, <<~~JS~~>>, <<~~CSS~~>>, <<~~FONTS~~>>)
- order should be -> NAVBAR and FOOTER, then all the rest, as the NAVBAR and FOOTER can also contain <<~~X~~>> links

*/

package main

import (
        "io/ioutil"
        "log"
        "strings"
        "os"
        "path/filepath"
        "fmt"
        "github.com/go-ini/ini"
        "strconv"
)

func render_site() {
        fmt.Println("Rendering!")
        create_folder_structure ()
        copy_src ()
        render_pages()
        resize_pictures()
        
        /* Removing the temporary files */
        
        /*os.Remove("temp")*/
}

func create_folder_structure () {
        /* prepare folder structure */
	
        /*os.Remove("rendered")*/
        os.MkdirAll("temp/pages", 0755)
        os.MkdirAll("temp/pages/gallery", 0755)
        os.MkdirAll("rendered/pages", 0755)
        os.MkdirAll("rendered/pages/gallery", 0755)
}

func copy_src () {
        src_items, _ := filepath.Glob("src/*")
        i := 0
        for i < len(src_items) {
                copydir(src_items[i], "rendered/"+strings.Split(src_items[i], "/")[1])
                i += 1
        }
}

func render_pages() {

        /* move the pages to temp */
        
        copydir("./pages", "./temp/pages")
        copydir("./gallery", "./temp/pages/gallery")
        
        /* load ini file */
        
        cfg, _ := ini.Load("config.ini")
        site_title := cfg.Section("general").Key("title").String()

        /* create navlist */
        
        item, _ := filepath.Glob("temp/pages/*.html")
        all_pages := []string{}
        menu_item := []string{}
        menu := make(map[int64]string)
        
        // checking whether the page is posted
        // checking whether the page should be present in the menu        
        
        i := 0
        for i < len(item) {
        
                input, err := ioutil.ReadFile(item[i])
                if err != nil {
                        log.Fatalln(err)
                }
        
                lines := strings.Split(string(input), "\n")
                
                j := 1
                for j < 5  {
                        if lines[j] == "posted" {
                                all_pages = append(all_pages, item[i])
                                k := 1
                                for k < 5 {
                                        if lines[k] == "in_menu" {
                                                menu_order, _ := strconv.ParseInt(strings.Split(lines[2], "_", )[2], 0, 64)
                                                menu_item = append(menu_item, strings.Split(item[i], "/")[2])
                                                menu[menu_order] = item[i]
                                        }
                                        k += 1
                                }
                        }
                j += 1
                }
        i += 1
        }
        
        // copy the navbar template to the temp folder 

        copyfile("templates/navbar_template.html", "temp/navbar.html")        

        // add navbar to the pages and resolve the ties NAVACTIVE, NAVLINK, NAVITEM
        // Cycling through all posted pages, then cycling through all menu items
        // Adding navlinks as necessary and resolving the ties
        
        i = 0
        for i < len(all_pages) {
                substitute(all_pages[i], "<<~~TITLE~~>>",site_title)
                inject_html(all_pages[i], "<<~~NAVBAR~~>>", "temp/navbar.html")

                j := 0
                for j < len(menu_item) {
                        var orig_link string = menu[int64(j)]
                        page_link := strings.Split(orig_link,"/")[2]
                        page_name := strings.Split(strings.Split(orig_link,"/")[2], ".")[0]
                        inject_nav_items(all_pages[i], "<<~~NAVLIST~~>>", "templates/navbar_item.html")
                        if page_link == strings.Split(all_pages[i],"/")[2] {
                                substitute(all_pages[i],"<<~~NAVACTIVE~~>>", "class=\"active\"")
                        } else {
                                substitute(all_pages[i],"<<~~NAVACTIVE~~>>", "")
                        }
                        
                        if strings.Split(all_pages[i],"/")[2] == "index.html" {
                                if page_link == "index.html" {
                                        substitute(all_pages[i], "<<~~NAVLINK~~>>",page_link)
                                } else {
                                        substitute(all_pages[i], "<<~~NAVLINK~~>>","pages/"+page_link)
                                }
                        } else {
                                if page_link == "index.html" {
                                        substitute(all_pages[i], "<<~~NAVLINK~~>>","../"+page_link)
                                } else {
                                        substitute(all_pages[i], "<<~~NAVLINK~~>>",page_link)
                                }
                        }
                        substitute(all_pages[i], "<<~~NAVITEM~~>>",page_name)
                        j += 1
                }

                substitute(all_pages[i], "<<~~NAVLIST~~>>","")
                
                /* populate the footer tie */
                
                inject_html(all_pages[i], "<<~~FOOTER~~>>", "templates/footer_template.html")
                
                /* resolve ties CSS, JS, PAGE */
                
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
                
                if strings.Split(all_pages[i],"/")[2] == "index.html" {
                        copyfile(all_pages[i], "rendered/"+strings.Split(all_pages[i],"/")[2])
                } else {
                        copyfile(all_pages[i], "rendered/pages/"+strings.Split(all_pages[i],"/")[2])
                }
                i += 1
                
        }

        /* create gallery.html content and sub-gallery htmls */
        
        dirs, _ := ioutil.ReadDir ("temp/pages/gallery/")
        
        all_galleries := []string{}
        all_galleries_name := []string{}
        
        i = 0
        for i < len(dirs) {
                if dirs[i].IsDir() == true {
                        copyfile("." + string(filepath.Separator) + "templates" + string(filepath.Separator) + "subgallery_template.html", "." + string(filepath.Separator) + "temp" + string(filepath.Separator) + "pages" + string(filepath.Separator) + "gallery" + string(filepath.Separator) + dirs[i].Name() + ".html")
                        all_galleries = append(all_galleries, "temp/pages/gallery/"+dirs[i].Name()+".html")
                        all_galleries_name = append(all_galleries_name, dirs[i].Name())
                }
                i += 1
        }
        
        i = 0
        for i < len(all_galleries) {
                inject_html("temp/pages/gallery.html", "<<~~GALLERYITEM~~>>", "templates/gallery_item.html")
                
                // Loop over all images and do the following updates
                //
                // SUBGALLERYLINK  = pages/gallery/filename.jpg
                // SUBGALLERYTHUMB = pages/gallery/filename_thumb.jpg
                // SUBGALLERYNAME  = all_galleries_name[i]
                
                imagepath := "temp/pages/gallery/"+all_galleries_name[i]+"/"
                images, _ := filepath.Glob(imagepath+"*")
                fmt.Println(all_galleries[i])
                a := 0
                for a < len(images) {
                        substitute("temp/pages/gallery.html","<<~~SUBGALLERYLINK~~>>",strings.Split(all_galleries[i],"temp/pages/")[1] )
                        substitute("temp/pages/gallery.html","<<~~SUBGALLERYTHUMB~~>>",strings.Split(strings.Split(images[a],"temp/pages/")[1],".")[0]+"_thumb.jpg")
                        substitute("temp/pages/gallery.html","<<~~SUBGALLERYNAME~~>>",all_galleries_name[i])
                        a += 1
                }
                
                
                
                //
                
                substitute(all_galleries[i], "<<~~TITLE~~>>",site_title)
                substitute(all_galleries[i], "<<~~GALLERYTITLE~~>>",all_galleries_name[i])
                
                /* inject navbar */
                
                inject_html(all_galleries[i], "<<~~NAVBAR~~>>", "temp/navbar.html")
                
                /* populate navbar with the correct links */
                
                j := 0
                for j < len(menu_item) {
                        var orig_link string = menu[int64(j)]
                        page_link := strings.Split(orig_link,"/")[2]
                        page_name := strings.Split(strings.Split(orig_link,"/")[2], ".")[0]
                        inject_nav_items(all_galleries[i], "<<~~NAVLIST~~>>", "templates/navbar_item.html")
                        substitute(all_galleries[i],"<<~~NAVACTIVE~~>>", "")
                        if page_link == "index.html" {
                                substitute(all_galleries[i], "<<~~NAVLINK~~>>","../../"+page_link)
                        } else {
                                substitute(all_galleries[i], "<<~~NAVLINK~~>>","../"+page_link)
                        }
                        substitute(all_galleries[i], "<<~~NAVITEM~~>>",page_name)
                        j += 1
                }
                substitute(all_galleries[i], "<<~~NAVLIST~~>>","")

                /* create html */
                
                
                
                /* populate the footer tie */
                
                inject_html(all_galleries[i], "<<~~FOOTER~~>>", "templates/footer_template.html")                
                
                /* resolve ties CSS, JS, PAGE */
                
                substitute(all_galleries[i], "<<~~JS~~>>","../../js/")
                substitute(all_galleries[i], "<<~~CSS~~>>","../../css/")
                substitute(all_galleries[i], "<<~~PAGE~~>>","../")
                
                //remove_header(all_galleries[i])
                copyfile(all_galleries[i], "rendered/pages/gallery/"+strings.Split(all_galleries[i],"/")[3])
                copydir("temp/pages/gallery/"+all_galleries_name[i], "rendered/pages/gallery/"+all_galleries_name[i])
                i += 1
                
        }
        substitute("temp/pages/gallery.html","<<~~GALLERYITEM~~>>","")
        copyfile("temp/pages/gallery.html", "rendered/pages/gallery.html")
}

func resize_pictures () {}

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

func inject_last_nav_item (file, tie, html_source_file string) {
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
                /*lines[line] = strings.Replace(lines[line], tie, s+"\n<<~~NAVLIST~~>>", -1)*/
                lines[line] = strings.Replace(lines[line], tie, s, -1)
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
        for i <= 5 {
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