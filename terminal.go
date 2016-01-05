// Copyright 2016 Bjørn Erik Pedersen <bjorn.erik.pedersen@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/pkg/term"
	"sync"
)

// TODO(bep): check different OSs
func readTerm(done, read chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	t, _ := term.Open("/dev/tty")
	defer t.Restore()
	defer t.Close()
	term.RawMode(t)

	for {
		select {
		case <-done:
			return
		default:
			b := make([]byte, 3)
			i, err := t.Read(b)

			if err != nil {
				return
			}

			if i == 1 && b[0] == 3 {
				// CTRL-C
				done <- true
				return
			}

			read <- true
		}
	}
}
