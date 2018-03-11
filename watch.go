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
   "github.com/go-fsnotify/fsnotify"
)

func watch() {

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
               log.Println("event:", event)
               if event.Op&fsnotify.Write == fsnotify.Write {
                   log.Println("modified file:", event.Name)
               }
               render_site()
           case err := <-watcher.Errors:
               log.Println("error:", err)
           }
       }
   }()

   err = watcher.Add("bodies")
   err = watcher.Add("src")
   err = watcher.Add("templates")
   if err != nil {
       log.Fatal(err)
   }
   <-done
}
