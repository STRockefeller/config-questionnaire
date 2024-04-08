package questionnaire

import (
	"reflect"
	"testing"

	"github.com/STRockefeller/collection"
	"github.com/charmbracelet/huh"
	"github.com/stretchr/testify/assert"
)

func TestGetFormItems_ValidStructType(t *testing.T) {
	// Arrange
	type TestStruct struct {
		Field1 string `yaml:"field1"`
		Field2 int    `yaml:"field2"`
		Field3 bool   `yaml:"field3"`
	}

	expectedFormItems := []huh.Field{
		huh.NewInput().Title("Field1").Value(new(string)),
		huh.NewInput().Title("Field2").Value(new(string)),
		huh.NewConfirm().Title("Field3").Value(new(bool)),
	}

	expectedStrFieldValues := collection.NewQueue[*string]()
	expectedStrFieldValues.Enqueue(new(string))
	expectedStrFieldValues.Enqueue(new(string))
	expectedBoolFieldValues := collection.NewQueue[*bool]()
	expectedBoolFieldValues.Enqueue(new(bool))

	// Act
	formItems, strFieldValues, boolFieldValues, err := getFormItems(reflect.TypeOf(TestStruct{}))

	// Assert
	assert.NoError(t, err)
	compareFormItems(t, expectedFormItems, formItems)
	assert.Equal(t, expectedStrFieldValues, strFieldValues)
	assert.Equal(t, expectedBoolFieldValues, boolFieldValues)
}

func compareFormItems(t *testing.T, expected, actual []huh.Field) {
	assert.Equal(t, len(expected), len(actual))
	for i, field := range expected {
		assert.Equal(t, field.GetKey(), actual[i].GetKey())
		assert.Equal(t, field.GetValue(), actual[i].GetValue())
	}
}

func TestGetFormItems_NonStructType(t *testing.T) {
	// Arrange
	var nonStructType int

	// Act
	formItems, strFieldValues, boolFieldValues, err := getFormItems(reflect.TypeOf(nonStructType))

	// Assert
	assert.Error(t, err)
	assert.Nil(t, formItems)
	assert.True(t, strFieldValues.IsEmpty())
	assert.True(t, boolFieldValues.IsEmpty())
}
