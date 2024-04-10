# Config Questionnaire

![GitHub License](https://img.shields.io/github/license/STRockefeller/config-questionnaire)![GitHub Top Language](https://img.shields.io/github/languages/top/STRockefeller/config-questionnaire)![GitHub Actions Status](https://img.shields.io/github/actions/workflow/status/STRockefeller/config-questionnaire/super-linter.yml)[![Go Report Card](https://goreportcard.com/badge/github.com/STRockefeller/go-linq)](https://goreportcard.com/report/github.com/STRockefeller/config-questionnaire)[![Coverage Status](https://coveralls.io/repos/github/STRockefeller/go-linq/badge.svg?branch=main)](https://coveralls.io/github/STRockefeller/config-questionnaire?branch=main)

A Go module for dynamically generating and running questionnaires based on the struct fields of provided config models.

## Overview

This module leverages reflection to inspect struct fields of any given config model and generates a questionnaire accordingly. It supports fields of type `string`, `int` (and its variants), and `bool`. The questionnaire is interactive and runs in the terminal using the [huh](https://github.com/charmbracelet/huh) library for form generation and input handling.

## Installation

To use this module in your Go project, ensure you have Go 1.22.1 or later, then run:

```bash
go get github.com/STRockefeller/config-questionnaire
```

## Usage

1. Define a config model as a Go struct. 
```go
package example

type Example struct {
    Name    string 
    Age     int    
    Married bool   
}
```

2. Use the `GenerateAndRunQuestionnaire` function to generate and run the questionnaire based on your model.

```go
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
```

## Dependencies

This module relies on several external libraries for its functionality:

- [github.com/STRockefeller/collection](https://github.com/STRockefeller/collection) for queue data structures.
- [github.com/charmbracelet/huh](https://github.com/charmbracelet/huh) for generating interactive forms.
- Various libraries from the [charmbracelet](https://github.com/charmbracelet) ecosystem for terminal UI components.

## License

This module is licensed under the MIT License - see the LICENSE file for details.
