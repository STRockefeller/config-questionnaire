package questionnaire

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/STRockefeller/collection"
	"github.com/charmbracelet/huh"
)

// GenerateAndRunQuestionnaire dynamically creates and runs a questionnaire based on the struct fields of the provided model.
func GenerateAndRunQuestionnaire[configModel any]() (configModel, error) {
	model := *new(configModel)
	ptr := reflect.New(reflect.TypeOf(model))
	val := ptr.Elem() // Obtain a reflect.Value that is addressable
	typ := val.Type()
	formItems, strFieldValues, boolFieldValues, err := getFormItems(typ)
	if err != nil {
		return model, err
	}
	form := huh.NewForm(huh.NewGroup(formItems...))
	if err := form.Run(); err != nil {
		return model, err
	}

	// Update model with form values
	for i := range formItems {
		field := typ.Field(i)
		switch field.Type.Kind() {
		case reflect.String:
			val.Field(i).SetString(*strFieldValues.Dequeue())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intVal, err := strconv.Atoi(*strFieldValues.Dequeue())
			if err != nil {
				return model, fmt.Errorf("invalid integer value for %s: %v", field.Name, err)
			}
			val.Field(i).SetInt(int64(intVal))
		case reflect.Bool:
			val.Field(i).SetBool(*boolFieldValues.Dequeue())
		}
	}

	return val.Interface().(configModel), nil
}

func getFormItems(typ reflect.Type) ([]huh.Field, collection.Queue[*string], collection.Queue[*bool], error) {
	// Ensure we're dealing with a struct
	if typ.Kind() != reflect.Struct {
		return nil, collection.NewQueue[*string](), collection.NewQueue[*bool](), fmt.Errorf("provided model must be a struct, got %s", typ.Kind())
	}

	var formItems []huh.Field
	strFieldValues := collection.NewQueue[*string]()
	boolFieldValues := collection.NewQueue[*bool]()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name

		switch field.Type.Kind() {
		case reflect.String:
			var fieldValue string
			formItems = append(formItems, huh.NewInput().Title(fieldName).Value(&fieldValue))
			strFieldValues.Enqueue(&fieldValue)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var fieldValue string // Use string to capture input, convert to int later
			formItems = append(formItems, huh.NewInput().Title(fieldName).Value(&fieldValue))
			strFieldValues.Enqueue(&fieldValue)
		case reflect.Bool:
			var fieldValue bool
			formItems = append(formItems, huh.NewConfirm().Title(fieldName).Value(&fieldValue))
			boolFieldValues.Enqueue(&fieldValue)
		// Add more types as needed
		default:
			return nil, collection.NewQueue[*string](), collection.NewQueue[*bool](), fmt.Errorf("unsupported field type: %s", field.Type.Kind())
		}
	}

	return formItems, strFieldValues, boolFieldValues, nil
}
