package example

import (
	"strconv"

	"github.com/charmbracelet/huh"
)

func ExampleQuestionnaire() (Example, error) {
	var (
		name   string
		age    string // was int
		hasPet bool
	)
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Name").Value(&name),
			huh.NewInput().Title("Age").Value(&age),
			huh.NewConfirm().Title("HasPet").Value(&hasPet),
		),
	)
	if err := form.Run(); err != nil {
		return Example{}, err
	}

	ageInt, err := strconv.Atoi(age)
	if err != nil {
		return Example{}, err
	}
	return Example{Name: name, Age: ageInt, HasPet: hasPet}, nil
}
