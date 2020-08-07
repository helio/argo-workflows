// +build !windows

package util

import "io/ioutil"

// WriteTeriminateMessage writes the termination message if possible. Panics if an error happens
// during writing.
func WriteTeriminateMessage(message string) {
	err := ioutil.WriteFile("/dev/termination-log", []byte(message), 0644)
	if err != nil {
		panic(err)
	}
}
