package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var path = flag.String("p", ".", "path to the project root")
var re = flag.String("r", "", "if specific, only packages whose path matches this regex will be updated")
var interactive = flag.Bool("i", false, "run goden in interactive mode")

func main() {
	var pattern *regexp.Regexp
	var pkg *build.Package
	var err error

	flag.Parse()

	pkg, err = build.ImportDir(*path, 0)
	if nil != err {
		fmt.Println(err)
		return
	}

	if "" != *re {
		pattern, err = regexp.Compile(*re)
		if nil != err {
			fmt.Println(err)
			return
		}
	}

	pkgs := make([]string, 0, len(pkg.Imports))
	for _, p := range pkg.Imports {
		ss := strings.Split(p, "/")
		if !strings.Contains(ss[0], ".") {
			continue
		}
		if pattern != nil && !pattern.MatchString(p) {
			continue
		}
		pkgs = append(pkgs, p)
	}

	if len(pkgs) == 0 {
		fmt.Println("No dependency needs updating.")
		return
	}

	if *interactive {
		reader := bufio.NewReader(os.Stdin)
		if nil != err {
			fmt.Println(err)
			return
		}
		for _, p := range pkgs {
			var update bool
			var choice string
			for true {
				fmt.Printf("Whether to update %s?[Y/n]:", p)
				choice, err = reader.ReadString('\n')
				if nil != err {
					fmt.Println(err)
					return
				}
				if '\n' == choice[0] || 'y' == choice[0] || 'Y' == choice[0] {
					update = true
					break
				} else if 'n' == choice[0] || 'N' == choice[0] {
					update = false
					break
				}
			}
			if update {
				err = goGetUpdate(p)
				if nil != err {
					fmt.Println(err)
				}
			}
		}
	} else {
		err = goGetUpdate(pkgs...)
		if nil != err {
			fmt.Println(err)
			return
		}
	}
}

func goGetUpdate(pkgs ...string) error {
	args := make([]string, 0, len(pkgs)+2)
	args = append(args, "get", "-u")
	args = append(args, pkgs...)
	cmd := exec.Command("go", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	fmt.Println(strings.Join(cmd.Args, " "))
	return cmd.Run()
}
