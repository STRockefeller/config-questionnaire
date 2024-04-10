package example

type Example struct {
	Name   string `questionnaire:"title:What's your name;"`
	Age    int
	HasPet bool
}
