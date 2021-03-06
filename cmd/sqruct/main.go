package main

import (
	"io/ioutil"
	"log"
	"path"

	"github.com/oov/sqruct/gen"

	"gopkg.in/yaml.v2"
)

func main() {
	const filename = "sqruct.yml"
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}
	var sq gen.Sqruct
	err = yaml.Unmarshal(b, &sq)
	if err != nil {
		log.Fatalln(err)
	}

	err = sq.Export(path.Dir(filename))
	if err != nil {
		log.Fatalln(err)
	}
}
