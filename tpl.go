package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"text/template"

	"gopkg.in/yaml.v2"
)

var inplace = flag.Bool("i", false, "write files in-place")

func ExecuteTemplates(valuesIn io.Reader, tplFile string) error {
	tpl, err := template.ParseFiles(tplFile)
	if err != nil {
		return fmt.Errorf("Error parsing template(s): %v", err)
	}

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, valuesIn)
	if err != nil {
		return fmt.Errorf("Failed to read standard input: %v", err)
	}

	var values map[string]interface{}
	err = yaml.Unmarshal(buf.Bytes(), &values)
	if err != nil {
		return fmt.Errorf("Failed to parse standard input: %v", err)
	}

	out := os.Stdout
	if *inplace {
		out, err = os.OpenFile(tplFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer out.Close()
	}
	err = tpl.Execute(out, values)
	if err != nil {
		return fmt.Errorf("Failed to parse standard input: %v", err)
	}
	return nil
}

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "usage:\n\tgotpl template < values.yml\n\tgotpl -i template < values.yml\n")
		os.Exit(1)
	}

	for _, tplFile := range flag.Args() {
		err := ExecuteTemplates(os.Stdin, tplFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
