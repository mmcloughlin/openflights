// +build ignore

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	yaml "gopkg.in/yaml.v2"
)

type Field struct {
	Name string
	Doc  string
	Type string
}

type Schema struct {
	Name          string
	Doc           string
	CollectionDoc string `yaml:"collection_doc"`
	Fields        []*Field
}

func LoadSchema(filename string) (*Schema, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	schema := &Schema{}
	err = yaml.Unmarshal(b, &schema)
	if err != nil {
		return nil, err
	}

	return schema, nil
}

type formatter func(string) string

var formatters = map[string]formatter{
	"string": strconv.Quote,
}

func Output(w io.Writer, s *Schema, r io.Reader) error {
	fmt.Fprintf(w, "// Generated code. DO NOT EDIT.\n\n")

	fmt.Fprintf(w, "package openflights\n")

	// type definition
	if s.Doc != "" {
		fmt.Fprintf(w, "// %s\n", s.Doc)
	}
	fmt.Fprintf(w, "type %s struct {\n", s.Name)
	for _, f := range s.Fields {
		if f == nil {
			continue
		}
		if f.Doc != "" {
			fmt.Fprintf(w, "\t// %s\n", f.Doc)
		}
		fmt.Fprintf(w, "\t%s %s\n", f.Name, f.Type)
	}
	fmt.Fprint(w, "}\n")

	// data array
	if s.CollectionDoc != "" {
		fmt.Fprintf(w, "// %s\n", s.CollectionDoc)
	}
	fmt.Fprintf(w, "var %ss = []%s{\n", s.Name, s.Name)
	c := csv.NewReader(r)
	for {
		record, err := c.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		fmt.Fprint(w, "\t{")
		for i, f := range s.Fields {
			if f == nil {
				continue
			}
			v := record[i]
			if xform, ok := formatters[f.Type]; ok {
				v = xform(v)
			}
			fmt.Fprintf(w, "%s,", v)
		}
		fmt.Fprint(w, "},\n")
	}
	fmt.Fprint(w, "}\n")

	return nil
}

var (
	schemaFilename = flag.String("schema", "", "schema yaml file")
	dataFilename   = flag.String("data", "", "data filename")
	outputFilename = flag.String("output", "", "output filename")
)

func main() {
	flag.Parse()

	schema, err := LoadSchema(*schemaFilename)
	if err != nil {
		log.Fatal(err)
	}

	data, err := os.Open(*dataFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	out, err := os.Create(*outputFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	err = Output(out, schema, data)
	if err != nil {
		log.Fatal(err)
	}
}
