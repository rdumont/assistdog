package assistdog

import (
	"reflect"
	"testing"

	"github.com/rdumont/assistdog/defaults"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/stretchr/testify/assert"
)

type person struct {
	Name   string
	Height int
}

func TestRemoveParser(t *testing.T) {
	assist := &Assist{parsers: map[reflect.Type]ParseFunc{}}
	assist.parsers[reflect.TypeOf("")] = defaults.ParseString

	assist.RemoveParser("")

	assert.Len(t, assist.parsers, 0)
}

func TestRemoveComparer(t *testing.T) {
	assist := &Assist{comparers: map[reflect.Type]CompareFunc{}}
	assist.comparers[reflect.TypeOf("")] = defaults.CompareString

	assist.RemoveComparer("")

	assert.Len(t, assist.comparers, 0)
}

func TestCreateInstance(t *testing.T) {
	t.Run("successfully", func(t *testing.T) {
		table := buildTable([][]string{
			{"Name", "John"},
			{"Height", "182"},
		})

		result, err := NewDefault().CreateInstance(new(person), table)
		if err != nil {
			t.Error(err)
			return
		}

		typed := result.(*person)
		assert.Equal(t, "John", typed.Name)
		assert.Equal(t, 182, typed.Height)
	})

	t.Run("with extra field", func(t *testing.T) {
		table := buildTable([][]string{
			{"Name", "John"},
			{"Height", "182"},
			{"Age", "25"},
		})

		_, err := NewDefault().CreateInstance(new(person), table)
		if !assert.Error(t, err) {
			return
		}

		assert.Equal(t, `failed to parse table as *assistdog.person:
- Age: field not found`, err.Error())
	})

	t.Run("with invalid integer", func(t *testing.T) {
		table := buildTable([][]string{
			{"Name", "John"},
			{"Height", "nono"},
		})

		_, err := NewDefault().CreateInstance(new(person), table)
		if !assert.Error(t, err) {
			return
		}

		assert.Equal(t, `failed to parse table as *assistdog.person:
- Height: strconv.Atoi: parsing "nono": invalid syntax`, err.Error())
	})
}

func TestCreateSlice(t *testing.T) {
	t.Run("successfully", func(t *testing.T) {
		table := buildTable([][]string{
			{"Name", "Height"},
			{"John", "182"},
			{"Mary", "170"},
		})

		result, err := NewDefault().CreateSlice(new(person), table)
		if !assert.NoError(t, err) {
			return
		}

		typed := result.([]*person)
		if !assert.Len(t, typed, 2) {
			return
		}

		assert.Equal(t, "John", typed[0].Name)
		assert.Equal(t, 182, typed[0].Height)

		assert.Equal(t, "Mary", typed[1].Name)
		assert.Equal(t, 170, typed[1].Height)
	})

	t.Run("with invalid integer", func(t *testing.T) {
		table := buildTable([][]string{
			{"Name", "Height"},
			{"John", "nono"},
		})

		_, err := NewDefault().CreateSlice(new(person), table)
		if !assert.Error(t, err) {
			return
		}

		assert.Equal(t, `failed to parse table as slice of *assistdog.person:
row 0:
  - Height: strconv.Atoi: parsing "nono": invalid syntax`, err.Error())
	})
}

func TestCompareInstance(t *testing.T) {
	t.Run("successfully", func(t *testing.T) {
		table := buildTable([][]string{
			{"Name", "John"},
			{"Height", "182"},
		})

		actual := &person{
			Name:   "John",
			Height: 182,
		}

		err := NewDefault().CompareToInstance(actual, table)
		assert.NoError(t, err)
	})

	t.Run("with different value for int", func(t *testing.T) {
		table := buildTable([][]string{
			{"Name", "John"},
			{"Height", "900"},
		})

		actual := &person{
			Name:   "John",
			Height: 182,
		}

		err := NewDefault().CompareToInstance(actual, table)
		if !assert.Error(t, err) {
			return
		}

		assert.Equal(t, `comparison failed:
- Height: expected 900, but got 182`, err.Error())
	})

	t.Run("with different value for string", func(t *testing.T) {
		table := buildTable([][]string{
			{"Name", "Mary"},
			{"Height", "182"},
		})

		actual := &person{
			Name:   "John",
			Height: 182,
		}

		err := NewDefault().CompareToInstance(actual, table)
		if !assert.Error(t, err) {
			return
		}

		assert.Equal(t, `comparison failed:
- Name: expected Mary, but got John`, err.Error())
	})
}

func TestCompareMap(t *testing.T) {
	t.Run("successfully", func(t *testing.T) {
		table := buildTable([][]string{
			{"Name", "John"},
			{"Height", "182"},
		})

		actual := map[string]interface{}{
			"Name":   "John",
			"Height": 182,
		}

		err := NewDefault().CompareToMap(actual, table)
		assert.NoError(t, err)
	})

	t.Run("with different value for int", func(t *testing.T) {
		table := buildTable([][]string{
			{"Name", "John"},
			{"Height", "900"},
		})

		actual := map[string]interface{}{
			"Name":   "John",
			"Height": 182,
		}

		err := NewDefault().CompareToMap(actual, table)
		if !assert.Error(t, err) {
			return
		}

		assert.Equal(t, `comparison failed:
- Height: expected 900, but got 182`, err.Error())
	})

	t.Run("with different value for string", func(t *testing.T) {
		table := buildTable([][]string{
			{"Name", "Mary"},
			{"Height", "182"},
		})

		actual := map[string]interface{}{
			"Name":   "John",
			"Height": 182,
		}

		err := NewDefault().CompareToMap(actual, table)
		if !assert.Error(t, err) {
			return
		}

		assert.Equal(t, `comparison failed:
- Name: expected Mary, but got John`, err.Error())
	})
}
func TestCompareSlice(t *testing.T) {
	t.Run("successfully", func(t *testing.T) {
		table := buildTable([][]string{
			{"Name", "Height"},
			{"John", "182"},
			{"Mary", "170"},
		})

		actual := []*person{
			{Name: "John", Height: 182},
			{Name: "Mary", Height: 170},
		}

		err := NewDefault().CompareToSlice(actual, table)
		assert.NoError(t, err)
	})

	t.Run("with different value for int", func(t *testing.T) {
		table := buildTable([][]string{
			{"Name", "Height"},
			{"John", "182"},
			{"Mary", "1234"},
		})

		actual := []*person{
			{Name: "John", Height: 182},
			{Name: "Mary", Height: 170},
		}

		err := NewDefault().CompareToSlice(actual, table)
		if !assert.Error(t, err) {
			return
		}

		assert.Equal(t, `comparison failed:
row 1:
  - Height: expected 1234, but got 170`, err.Error())
	})

	t.Run("passing something other than a slice", func(t *testing.T) {
		table := buildTable([][]string{
			{"Name", "Height"},
			{"John", "182"},
			{"Mary", "1234"},
		})

		actual := &person{}

		err := NewDefault().CompareToSlice(actual, table)
		if !assert.Error(t, err) {
			return
		}

		assert.Equal(t, `actual value is not a slice`, err.Error())
	})
}

func buildTable(src [][]string) *gherkin.DataTable {
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
