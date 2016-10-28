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
	"log"
	"github.com/fsnotify/fsnotify"
	"path/filepath"
	"os"
	"strings"
)

func autorender_site() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Rename == fsnotify.Rename || event.Op&fsnotify.Remove == fsnotify.Remove  {
					//log.Println("modified file:", event.Name)
					//log.Println("event:", event)
					render_site()
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()
    
    folders := Subfolders("pages")
    folders = append(folders, Subfolders("src")...)
    folders = append(folders, Subfolders("templates")...)
    
	for _, folder := range folders {
		err = watcher.Add(filepath.Join(folder))
		if err != nil {
			log.Fatal(err)
		}
	}
	
	<-done
}

func Subfolders(path string) (paths []string) {
	filepath.Walk(path, func(newPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			name := info.Name()
			// skip folders that begin with a dot
			if ignoreDir(name) && name != "." && name != ".." {
				return filepath.SkipDir
			}
			paths = append(paths, newPath)
		}
		return nil
	})
	return paths
}

func ignoreDir(name string) bool {
	return strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_")
}