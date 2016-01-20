/* TODO

- delete the current rendered folder
- resize pictures and create thumbnails (x_thumb.jpg)
- create gallery.html and the separate gallery pages, based on templates (gallery_template.html and subgallery_template.html)
- move all galleries from dist/gallery to dist/rendered/gallery
- move all pages from dist/pages to dist/rendered/pages, move index.html to dist/rendered
- move all css/js/... from dist/src to dist/rendered/
- replace all placeholders with correct html (<<~~NAVBAR~~>>, <<~~FOOTER~~>>, <<~~JS~~>>, <<~~CSS~~>>, <<~~FONTS~~>>)

*/

package main

import (
        "io/ioutil"
        "log"
        "strings"
)

func render_site(file string) {
        input, err := ioutil.ReadFile(file)
        if err != nil {
                log.Fatalln(err)
        }

        lines := strings.Split(string(input), "\n")

        for i, line := range lines {
                if strings.Contains(line, "<<~~ NAVBAR ~~>>") {
                        lines[i] = "REPLACED"
                }
        }
        output := strings.Join(lines, "\n")
        err = ioutil.WriteFile(file, []byte(output), 0644)
        if err != nil {
                log.Fatalln(err)
        }
}