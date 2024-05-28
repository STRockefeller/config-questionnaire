package questionnaire

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

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
	if err := updateModelValues(val, typ, strFieldValues, boolFieldValues); err != nil {
		return model, err
	}

	return val.Interface().(configModel), nil
}

// getFormItems generates form items from a struct's fields. It checks for a "questionnaire" tag to customize field names, creates form items based on field types (string, int, bool), and enqueues their values in respective queues. Unsupported field types result in an error.
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
		// Check for the questionnaire tag and use it if present
		questionnaireTags, ok := field.Tag.Lookup("questionnaire")
		fieldName := field.Name
		if ok && questionnaireTags != "" {
			tagParts := strings.Split(questionnaireTags, ";")
			for _, tagPart := range tagParts {
				if strings.Contains(tagPart, "title:") {
					fieldName = strings.Split(tagPart, ":")[1]
					break
				}
			}
		}

		switch field.Type.Kind() {
		case reflect.String:
			var fieldValue string
			formItems = append(formItems, huh.NewInput().Title(fieldName+"(string)").Value(&fieldValue))
			strFieldValues.Enqueue(&fieldValue)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var fieldValue string // Use string to capture input, convert to int later
			formItems = append(formItems, huh.NewInput().Title(fieldName+"(int)").Value(&fieldValue))
			strFieldValues.Enqueue(&fieldValue)
		case reflect.Bool:
			var fieldValue bool
			formItems = append(formItems, huh.NewConfirm().Title(fieldName+"(bool)").Value(&fieldValue))
			boolFieldValues.Enqueue(&fieldValue)
		case reflect.Struct:
			// Recursively get form items for nested struct
			nestedItems, nestedStrValues, nestedBoolValues, err := getFormItems(field.Type)
			if err != nil {
				return nil, collection.NewQueue[*string](), collection.NewQueue[*bool](), err
			}
			formItems = append(formItems, nestedItems...)
			for !nestedStrValues.IsEmpty() {
				strFieldValues.Enqueue(nestedStrValues.Dequeue())
			}
			for !nestedBoolValues.IsEmpty() {
				boolFieldValues.Enqueue(nestedBoolValues.Dequeue())
			}
		default:
			return nil, collection.NewQueue[*string](), collection.NewQueue[*bool](), fmt.Errorf("unsupported field type: %s", field.Type.Kind())
		}
	}

	return formItems, strFieldValues, boolFieldValues, nil
}

// updateModelValues updates the model's fields with the values from the form.
func updateModelValues(val reflect.Value, typ reflect.Type, strFieldValues collection.Queue[*string], boolFieldValues collection.Queue[*bool]) error {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)
		switch field.Type.Kind() {
		case reflect.String:
			fieldVal.SetString(*strFieldValues.Dequeue())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intVal, err := strconv.Atoi(*strFieldValues.Dequeue())
			if err != nil {
				return fmt.Errorf("invalid integer value for %s: %v", field.Name, err)
			}
			fieldVal.SetInt(int64(intVal))
		case reflect.Bool:
			fieldVal.SetBool(*boolFieldValues.Dequeue())
		case reflect.Struct:
			if err := updateModelValues(fieldVal, field.Type, strFieldValues, boolFieldValues); err != nil {
				return err
			}
		}
	}
	return nil
}
