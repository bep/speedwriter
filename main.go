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
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"unicode"
)

func main() {

	const (
		helpText string = `Usage: speedwriter [FILE]
		
		An alternative to the above, is to pipe the text to stdin, e.g.:
			
			cat myfile.txt | speedwriter
			
		This makes it easy to make it look like you're coding like Linux Torvalds:
		
		    curl -s https://raw.githubusercontent.com/git/git/master/block-sha1/sha1.c | egrep -v "^(//|/\*| \*)" | tail -n +153 | speedwriter
			
		
		Options:
		
		    -help        display this help and exit
            -version     output version information and exit
		    
		`

		versionText string = `speedwriter 0.1`
	)

	var (
		fileReader io.Reader
		err        error
	)

	if flag.NArg() == 0 || flag.Arg(0) == "-" {
		stat, _ := os.Stdin.Stat()
		isPipe := (stat.Mode() & os.ModeCharDevice) == 0
		if !isPipe {
			log.Fatal("error: Nothing to read")
		}
		fileReader = os.Stdin
	} else {
		filename := flag.Arg(0)
		fileReader, err = os.Open(filename)

		if err != nil {
			log.Fatalf("error: Failed to read input file %s: %s", filename, err)
		}

		defer fileReader.(io.Closer).Close()
	}

	var (
		textReader = bufio.NewReader(fileReader)

		done = make(chan bool, 0)
		read = make(chan bool, 4) // buffer 4 is arbitrary, but should be plenty.
		stop = make(chan os.Signal, 1)
	)

	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)

	wg := &sync.WaitGroup{}

	go func() {
		<-stop
		done <- true
	}()

	wg.Add(2)

	// clear screen
	fmt.Print("\033[H\033[2J")

	go readTerm(done, read, wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			select {
			case <-read:
				for {
					r, _, err := textReader.ReadRune()
					if err != nil {
						done <- true
						return
					}
					if r == '\n' {
						fmt.Print("\r\n")
					} else {
						fmt.Print(string(r))

					}
					// to make it extra speedy, make spaces "automatic"
					if !unicode.IsSpace(r) {
						break
					}
				}
			case <-done:
				return
			}
		}
	}(wg)

	wg.Wait()

}

func init() {
	flag.Parse()
	log.SetPrefix("speedwriter")
}
