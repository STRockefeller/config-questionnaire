package main

import (
	"fmt"

	questionnaire "github.com/STRockefeller/config-questionnaire"
	"github.com/STRockefeller/config-questionnaire/example"
)

func main() {
	e, err := questionnaire.GenerateAndRunQuestionnaire[example.Example]()
	if err != nil {
		panic(err)
	}
	fmt.Println(e)
}
