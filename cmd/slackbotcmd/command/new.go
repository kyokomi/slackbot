package command

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

func CmdNew(c *cli.Context) {
	log.SetFlags(log.Llongfile | log.LstdFlags)

	pkgName := c.String("pkg")
	if len(pkgName) <= 0 {
		fmt.Fprintf(os.Stderr, "not selected pkg name")
		os.Exit(2)
	}

	// mkdir
	if err := os.Mkdir(pkgName, 0755); err != nil {
		log.Fatalln(err)
	}

	// execute template
	var buf bytes.Buffer
	if err := PluginTmpl(&buf, pkgName); err != nil {
		log.Fatalln(err)
	}

	// go fmt
	srcBuf, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatalln(err)
	}

	// Write your code here
	filePath := filepath.Join(pkgName, pkgName+".go")
	if err := ioutil.WriteFile(filePath, srcBuf, 0644); err != nil {
		log.Fatalln(err)
	}

	// execute template
	buf.Reset()
	if err := PluginTestTmpl(&buf, pkgName); err != nil {
		log.Fatalln(err)
	}

	// go fmt
	srcBuf, err = format.Source(buf.Bytes())
	if err != nil {
		log.Fatalln(err)
	}

	// Write your code here
	filePath = filepath.Join(pkgName, pkgName+"_test.go")
	if err := ioutil.WriteFile(filePath, srcBuf, 0644); err != nil {
		log.Fatalln(err)
	}
}
