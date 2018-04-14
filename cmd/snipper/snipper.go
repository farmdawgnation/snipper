package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/farmdawgnation/snipper/pkg/processor"
)

var (
	showHelp bool
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.Parse()

	if showHelp == true || len(flag.Args()) < 2 {
		fmt.Println("snipper - snippet style transformers for YAML")
		fmt.Println("Usage: snipper template.yaml transformer.yaml [transformer.yaml [transformer.yaml ...]]")
	} else {
		args := flag.Args()
		argCount := len(args)
		parsedFiles := make([]map[interface{}]interface{}, argCount)

		for index, arg := range args {
			holdingMap := make(map[interface{}]interface{})
			data, err := ioutil.ReadFile(arg)
			check(err)

			err = yaml.Unmarshal([]byte(data), &holdingMap)
			check(err)

			parsedFiles[index] = holdingMap
		}

		template := parsedFiles[0]
		transformers := parsedFiles[1:argCount]

		for _, transformer := range transformers {
			template = processor.Process(template, transformer)
		}

		rendered, err := yaml.Marshal(&template)
		check(err)

		fmt.Println(string(rendered))
	}
}
