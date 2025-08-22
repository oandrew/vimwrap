package main

import (
	"flag"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
)

var prefixesFlag = flag.String("p", "", "allowed prefixes")

var vimFlag = flag.String("vim", "/usr/bin/vim", "path to vim")

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func hasAnyPrefix(s string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

func main() {
	flag.Parse()

	prefixes := []string{}
	if *prefixesFlag != "" {
		prefixes = strings.Split(*prefixesFlag, ",")
	}

	args := slices.Clone(flag.Args())

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	for i := range args {
		arg := args[i]
		if hasAnyPrefix(arg, prefixes) && filepath.IsLocal(arg) {
			newArg := filepath.Join(home, arg)
			// spew.Dump(newArg)
			if !exists(arg) && exists(newArg) {
				args[i] = newArg
			}
		}
	}

	args = append([]string{"vim"}, args...)
	syscall.Exec(*vimFlag, args, os.Environ())

}
