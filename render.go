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
                                                menu_item = append(menu_item, strings.Split(item[i], "/")[2])
                                        }
                                        k += 1
                                }
                        }
                j += 1
                }
        i += 1
        }
        
        // fmt.Println("allpages: ", all_pages)
        // fmt.Println("menu_item: ", menu_item)
        
        /* copy the navbar template to the temp folder and add the correct amount of navitems to the navbar */

        copyfile("templates/navbar_template.html", "temp/navbar.html")        
        i = 0
        for i < len(menu_item) {
                if i+1 < len(menu_item) {
                        inject_nav_items("temp/navbar.html", "<<~~NAVLIST~~>>", "templates/nav_item.html")
                } else {
                        inject_last_nav_item("temp/navbar.html", "<<~~NAVLIST~~>>", "templates/nav_item.html")
                }
                i += 1
        }

        /* add navbar to the pages and resolve the ties NAVACTIVE, NAVLINK, NAVITEM */ 
        
        i = 0
        for i < len(all_pages) {
                inject_html(all_pages[i], "<<~~NAVBAR~~>>", "temp/navbar.html")
                substitute(all_pages[i], "<<~~NAVACTIVE~~>>","replace1")
                substitute(all_pages[i], "<<~~NAVLINK~~>>","replace2")
                substitute(all_pages[i], "<<~~NAVITEM~~>>","replace3")
                substitute(all_pages[i], "<<~~TITLE~~>>",site_title)
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