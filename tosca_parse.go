package main

import (
	"log"
	// "net/http"

	"flag"
	"fmt"
	"github.com/owulveryck/toscalib"
	// "github.com/owulveryck/toscaviewer"
	"os"
	"io/ioutil"
	"path/filepath"
	
	// "gopkg.in/yaml.v2"
	"archive/zip"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/zipfs"
)

type toscaDefinition struct{
	toscalib.ServiceTemplateDefinition
}

func noop(source string, std *toscalib.ServiceTemplateDefinition) error {
	return nil
}

func (t *toscaDefinition) ParseVNFD(zipfile string, plan string) error {

	rc, err := zip.OpenReader(zipfile)
	if err != nil {
		return err
	}
	defer rc.Close()
	fs := zipfs.New(rc, zipfile)

	dirname := fmt.Sprintf("/%v", filepath.Dir(plan))
	base := filepath.Base(plan)
	ns := vfs.NameSpace{}
	ns.Bind("/", fs, dirname, vfs.BindReplace)

	// pass in a resolver that has the context of the virtual filespace
	// of the archive file to handle resolving imports
	return t.ParseSource(base, func(l string) ([]byte, error) {
		var r []byte
		rsc, err := ns.Open(l)
		if err != nil {
			return r, err
		}
		return ioutil.ReadAll(rsc)
	}, toscalib.ParserHooks{ParsedSTD: noop}) // TODO(kenjones): Add hooks as method parameter
}

func main() {
	
	var zipfile = flag.String("zipfile", "", "a VNFD zip file to process")
	var plan = flag.String("plan", "", "the plan file in VNFD zip")
	flag.Parse()

	var toscaTemplate toscaDefinition

	//err = yaml.Unmarshal(file, &toscaTemplate)
	err = toscaTemplate.ParseVNFD(*zipfile, *plan)
	if err != nil {
		log.Panic("error: ", err)
	}
	// router := toscaviewer.NewRouter(&toscaTemplate)

	// log.Println("connect here: http://localhost:8080/svg")
	// log.Fatal(http.ListenAndServe(":8080", router))

}
