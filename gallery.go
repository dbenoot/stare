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

// Create the html for the galleries

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	// 	"sort"
	"strings"
	// 	"text/template"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
	"image/jpeg"
	"text/template"
)

func renderGalleries() {

	galleries := mapGalleries(filepath.Join("bodies", "galleries"))
	pictures := mapImg(filepath.Join("bodies", "galleries"))

	removeThumbs(pictures)

	for _, value := range galleries {
		// os.Remove(filepath.Join("bodies", "galleries", value, "gallery.html"))
		os.Remove(filepath.Join("bodies", "galleries", value, "gallery.jpg"))
		createSubGalleries(value)
	}

	for _, value := range pictures {
		resizePicture(value)
	}
}

func removeThumbs(p []string) {
	for _, value := range p {
		if strings.Contains(value, "_starethumb.jpg") {
			os.Remove(value)
		}
	}
}

func mapImg(path string) []string {

	// bodies := make(map[string]string)
	formats := []string{".jpg", ".JPG", ".jpeg", ".JPEG"}

	files := []string{}
	pictures := []string{}
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	for _, file := range files {
		if stringInSlice(filepath.Ext(file), formats) == true && strings.Contains(file, "_starethumb.jpg") == false {
			pictures = append(pictures, file)
		}
	}

	check(err)
	return pictures
}

func mapGalleries(path string) []string {

	var g []string

	files, err := ioutil.ReadDir(path)
	check(err)

	for _, f := range files {
		g = append(g, f.Name())
	}

	return g
}

func resizePicture(filename string) {

	// resize and create new thumbs

	// imgName := filepath.Base(filename)
	imgThumb := strings.TrimSuffix(filename, filepath.Ext(filename)) + "_starethumb.jpg"

	file, err := os.Open(filename)
	check(err)
	defer file.Close()

	img, err := jpeg.Decode(file)
	check(err)

	b := img.Bounds()
	imgWidth := b.Max.X
	imgHeight := b.Max.Y

	if imgWidth > 1000 {

		// resize to width 1000 using Lanczos resampling
		// and preserve aspect ratio
		m := resize.Resize(1000, 0, img, resize.Lanczos3)

		out, err := os.Create(filename)
		check(err)
		defer out.Close()

		jpeg.Encode(out, m, nil)
	}

	// Create the thumbnails

	n := img

	switch {
	case imgHeight <= 275:
		n, err = cutter.Crop(img, cutter.Config{
			Width:  350,
			Height: 275,
			Mode:   cutter.Centered,
		})
	default:
		tempn := resize.Resize(350, 0, img, resize.Lanczos3)
		n, err = cutter.Crop(tempn, cutter.Config{
			Width:  350,
			Height: 275,
			Mode:   cutter.Centered,
		})
	}

	out2, err := os.Create(imgThumb)
	check(err)
	defer out2.Close()

	// write resized image and thumbnails to file

	jpeg.Encode(out2, n, nil)
}

func createSubGalleries(path string) {

	// create subgallery html file

	images := mapImg(filepath.Join("bodies", "galleries", path))
	sgi, _ := template.ParseFiles("templates/subgallery_item.html")
	// newpath := filepath.Join("bodies", "galleries", path, "gallery.html")

	w := bytes.NewBufferString("")
	for key, value := range images {
		sgi.Execute(w, map[string]string{"Subimage": filepath.Base(value), "Subimagethumb": strings.TrimSuffix(filepath.Base(value), filepath.Ext(value)) + "_starethumb.jpg"})
		if key == 0 {
			createGalleryJpg(value)
		}
	}

	if _, err := os.Stat(filepath.Join("bodies", "galleries", path, "gallery.html")); err == nil {
		fmt.Println(filepath.Join("bodies", "galleries", path, "gallery.html"), "already exists. Skipped.")
	} else {
		createPage(filepath.Join("bodies", "galleries", path), "gallery", "\n\n"+w.String())
	}
}

func createGalleryJpg(n string) {
	dest := strings.TrimSuffix(n, filepath.Base(n)) + "gallery.jpg"
	copyfile(n, dest)
}
