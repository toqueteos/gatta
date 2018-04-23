package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/toqueteos/gatta"
)

func main() {
	var key string

	flag.StringVar(&key, "key", "", "Subscription Key, accessible from Azure's Translator Text API > Keys")
	flag.Parse()

	if len(key) != 32 {
		fmt.Println("Invalid -key provided, it should be 32 characters long")
		os.Exit(1)
	}

	tr, err := gatta.New(key)
	fatalIf(err)

	resp, err := tr.Translate("Hello World!", "es")
	fatalIf(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fatalIf(err)

	fmt.Println(string(body))
}

func fatalIf(err error) {
	if err != nil {
		file, line := caller()
		fmt.Println(file, line, err)
		os.Exit(1)
	}
}

func caller() (file string, line int) {
	var ok bool
	_, file, line, ok = runtime.Caller(2)
	if !ok {
		file = "???"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return
}
