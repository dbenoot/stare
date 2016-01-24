/* TODO

- delete the current rendered folder
- resize pictures and create thumbnails (x_thumb.jpg)
- create gallery.html and the separate gallery pages, based on templates (gallery_template.html and subgallery_template.html)
- move all galleries from dist/gallery to dist/rendered/gallery
- move all pages from dist/pages to dist/rendered/pages, move index.html to dist/rendered
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
	fmt.Println("Rendering!\n")
	
	/* prepare folder structure */
	
        /*os.Remove("rendered")*/
        os.MkdirAll("temp/pages", 0755)
        os.MkdirAll("temp/gallery", 0755)
        os.MkdirAll("rendered/pages", 0755)
        os.MkdirAll("rendered/gallery", 0755)
        
        /* move the pages to temp */
        
        copydir("./pages", "./temp/pages")
        copydir("./gallery", "./temp/gallery")
        
        /* load ini file */
        
        cfg, _ := ini.Load("config.ini")
        site_title := cfg.Section("general").Key("title").String()

        /* create navlist */
        
        item, _ := filepath.Glob("temp/pages/*.html")
        all_pages := []string{}
        menu_item := []string{}
        menu := make(map[int64]string)
        
        i := 0
        
        for i < len(item) {
        
                // check whether the page is posted
                
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
                                                // menu_item[menu_order] = strings.Split(item[i], "/")[2]
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
        
        // reversed_menu := make(map[string]int64)

        // for key,value := range menu {
        //         reversed_menu[value] = key
        // }

        /* copy the navbar template to the temp folder and add the correct amount of navitems to the navbar */

        copyfile("templates/navbar_template.html", "temp/navbar.html")        

        /* add navbar to the pages and resolve the ties NAVACTIVE, NAVLINK, NAVITEM */
        /* and resolve ties CSS, JS, PAGE */
        
        i = 0
        for i < len(all_pages) {
                substitute(all_pages[i], "<<~~TITLE~~>>",site_title)
                inject_html(all_pages[i], "<<~~NAVBAR~~>>", "temp/navbar.html")

                j := 0
                for j < len(menu_item) {
                        var orig_link string = menu[int64(j)]
                        page_link := strings.Split(orig_link,"/")[2]
                        page_name := strings.Split(strings.Split(orig_link,"/")[2], ".")[0]

                                inject_nav_items(all_pages[i], "<<~~NAVLIST~~>>", "templates/nav_item.html")
                                if page_link == strings.Split(all_pages[i],"/")[2] {
                                        substitute(all_pages[i],"<<~~NAVACTIVE~~>>", "class=\"active\"")
                                } else {
                                        substitute(all_pages[i],"<<~~NAVACTIVE~~>>", "")
                                }
                                
                                if strings.Split(all_pages[i],"/")[2] == "index.html" {
                                        substitute(all_pages[i], "<<~~NAVLINK~~>>","pages/"+page_link)
                                } else {
                                        substitute(all_pages[i], "<<~~NAVLINK~~>>",page_link)
                                }
                                substitute(all_pages[i], "<<~~NAVITEM~~>>",page_name)
                        j += 1
                }

                substitute(all_pages[i], "<<~~NAVLIST~~>>","")
                
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
                i += 1
                
        }

        /* resize pictures */
        
        /* create gallery.html and sub-galleries */
        
        
        /*os.Remove("temp")*/

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