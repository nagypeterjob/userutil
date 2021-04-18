// +build !windows

package userutil

//
// Open, Parse, Read are influenced by:
// https://github.com/tailscale/tailscale/blob/main/util/lineread/lineread.go
//

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func Open(name string, fn func(string) error) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return Read(f, fn)
}

func Parse(content io.Reader, fn func(string) error) error {
	return Read(content, fn)
}

func Read(r io.Reader, fn func(string) error) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if err := fn(scanner.Text()); err != nil {
			return err
		}
	}
	return scanner.Err()
}

func toFilter(obj string) bool {
	for _, filter := range filters {
		if strings.HasPrefix(obj, filter) {
			return true
		}
	}
	return false
}
