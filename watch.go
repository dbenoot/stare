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
	"fmt"
	"github.com/rjeczalik/notify"
)

func watch() {

	fmt.Println("Watching updates and rendering. Press ctrl+c to stop.")

	// Make the channel buffered to ensure no event is dropped. Notify will drop
	// an event if the receiver is not able to keep up the sending pace.

	c := make(chan notify.EventInfo, 10)

	// Set up a watchpoint listening on events within recursive working directory.
	// Dispatch all events separately to c.

	err := notify.Watch("bodies/...", c, notify.All)
	err = notify.Watch("src/...", c, notify.All)
	err = notify.Watch("templates/...", c, notify.All)
	check(err)
	defer notify.Stop(c)

	done := make(chan bool)

	go func() {
		for {
			select {
			case <-c:
				fmt.Printf(".")
				render_site()
			}
		}

	}()

	<-done
}
