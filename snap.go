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

// This file contains the functions to create a site snapshot

package main

import (
	// 	"fmt"
	"os"
	"path/filepath"
	// 	"strings"
	"time"
)

func snap() {
	t := time.Now()

	pwd := filepath.Join("snapshots", t.Format("20060102T1504"))

	if _, err := os.Stat(pwd); os.IsNotExist(err) {
		_ = os.MkdirAll(pwd, 0755)
	}

	err := copydir("rendered", pwd)
	check(err)

}
