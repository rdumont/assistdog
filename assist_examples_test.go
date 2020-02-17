package assistdog_test

import (
	"reflect"

	"github.com/cucumber/godog/gherkin"
	"github.com/rdumont/assistdog"
)

func ExampleAssist_CreateInstance() {
	table := NewTable([][]string{
		{"Name", "John"},  //  | Name   | John |
		{"Height", "182"}, //  | Height | 182  |
	})

	assist := assistdog.NewDefault()
	result, err := assist.CreateInstance(new(Person), table)
	if err != nil {
		panic(err)
	}

	reflect.DeepEqual(result, &Person{
		Name:   "John",
		Height: 182,
	})
}

func ExampleAssist_CreateSlice() {
	table := NewTable([][]string{
		{"Name", "Height"}, // | Name | Height |
		{"John", "182"},    // | John | 182    |
		{"Mary", "170"},    // | Mary | 170    |
	})

	assist := assistdog.NewDefault()
	result, err := assist.CreateSlice(new(Person), table)
	if err != nil {
		panic(err)
	}

	reflect.DeepEqual(result, []*Person{
		{Name: "John", Height: 182},
		{Name: "Mary", Height: 170},
	})
}

func ExampleAssist_CompareToInstance() {
	table := NewTable([][]string{
		{"Name", "John"},  //  | Name   | John |
		{"Height", "182"}, //  | Height | 182  |
	})

	actual := &Person{
		Name:   "John",
		Height: 182,
	}

	assist := assistdog.NewDefault()
	err := assist.CompareToInstance(actual, table)
	if err != nil {
		panic(err)
	}
}

func ExampleAssist_CompareToSlice() {
	table := NewTable([][]string{
		{"Name", "Height"}, // | Name | Height |
		{"John", "182"},    // | John | 182    |
		{"Mary", "170"},    // | Mary | 170    |
	})

	actual := []*Person{
		{Name: "John", Height: 182},
		{Name: "Mary", Height: 170},
	}

	assist := assistdog.NewDefault()
	err := assist.CompareToSlice(actual, table)
	if err != nil {
		panic(err)
	}
}

type Person struct {
	Name   string
	Height int
}

func NewTable(src [][]string) *gherkin.DataTable {
	rows := make([]*gherkin.TableRow, len(src))
	for i, row := range src {
		cells := make([]*gherkin.TableCell, len(row))
		for j, value := range row {
			cells[j] = &gherkin.TableCell{Value: value}
		}

		rows[i] = &gherkin.TableRow{Cells: cells}
	}

	return &gherkin.DataTable{Rows: rows}
}
