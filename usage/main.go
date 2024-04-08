package main

import (
	"fmt"

	questionnaire "github.com/STRockefeller/yaml-questionnaire"
	"github.com/STRockefeller/yaml-questionnaire/example"
)

func main() {
	e, err := questionnaire.GenerateAndRunQuestionnaire[example.Example]()
	if err != nil {
		panic(err)
	}
	fmt.Println(e)
}
