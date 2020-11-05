package assistdog

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/cucumber/godog"

	"github.com/rdumont/assistdog/defaults"
)

var defaultParsers = map[interface{}]ParseFunc{
	"":          defaults.ParseString,
	0:           defaults.ParseInt,
	time.Time{}: defaults.ParseTime,
}

var defaultComparers = map[interface{}]CompareFunc{
	"":          defaults.CompareString,
	0:           defaults.CompareInt,
	time.Time{}: defaults.CompareTime,
}

// ParseFunc parses a raw string value from a table into a given type.
// If it succeeds, it should return the parsed typed value. Otherwise, it should return an error
// describing why the value could not be parsed.
type ParseFunc func(raw string) (interface{}, error)

// CompareFunc compares a raw string value from a table to an actual typed value.
// If the values are considered a match, no error should be returned. Otherwise, an error that
// describes the differences should be returned.
type CompareFunc func(raw string, actual interface{}) error

// NewDefault creates a new Assist instance with all the default parsers and comparers.
func NewDefault() *Assist {
	a := new(Assist)
	for tp, p := range defaultParsers {
		a.RegisterParser(tp, p)
	}

	for tp, c := range defaultComparers {
		a.RegisterComparer(tp, c)
	}

	return a
}

// Assist provides utility methods to deal with Gherkin tables.
type Assist struct {
	lock      sync.RWMutex
	parsers   map[reflect.Type]ParseFunc
	comparers map[reflect.Type]CompareFunc
}

// RegisterParser registers a new value parser for a type.
// If a previous parser already exists for the given type, it will be replaced.
func (a *Assist) RegisterParser(i interface{}, parser ParseFunc) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.assertInit()
	a.parsers[reflect.TypeOf(i)] = parser
}

// RegisterComparer registers a new value comparer for a type.
// If a previous comparer already exists for the given type, it will be replaced.
func (a *Assist) RegisterComparer(i interface{}, comparer CompareFunc) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.assertInit()
	a.comparers[reflect.TypeOf(i)] = comparer
}

// RemoveParser removes the value parser for a type.
func (a *Assist) RemoveParser(i interface{}) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.assertInit()
	delete(a.parsers, reflect.TypeOf(i))
}

// RemoveComparer removes the value comparer for a type.
func (a *Assist) RemoveComparer(i interface{}) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.assertInit()
	delete(a.comparers, reflect.TypeOf(i))
}

// ParseMap takes a Gherkin table and returns a map that represents it.
// The table must have exactly two columns, where the first represents
// the key and the second represents the value.
func (a *Assist) ParseMap(table *godog.Table) (map[string]string, error) {
	if len(table.Rows) == 0 {
		return nil, fmt.Errorf("expected table to have at least one row")
	}

	if len(table.Rows[0].Cells) != 2 {
		return nil, fmt.Errorf("expected table to have exactly two columns")
	}

	result := map[string]string{}
	for _, row := range table.Rows {
		result[row.Cells[0].Value] = row.Cells[1].Value
	}

	return result, nil
}

// ParseSlice takes a Gherkin table and returns a slice of maps representing each row.
// The first row acts as a header and provides the keys.
func (a *Assist) ParseSlice(table *godog.Table) ([]map[string]string, error) {
	if len(table.Rows) < 2 {
		return nil, fmt.Errorf("expected table to have at least two rows")
	}

	if len(table.Rows[0].Cells) == 0 {
		return nil, fmt.Errorf("expected table to have at least one column")
	}

	fieldCells := table.Rows[0].Cells

	result := make([]map[string]string, len(table.Rows)-1)
	for i := 1; i < len(table.Rows); i++ {
		parsed := map[string]string{}
		for j := 0; j < len(fieldCells); j++ {
			parsed[fieldCells[j].Value] = table.Rows[i].Cells[j].Value
		}
		result[i-1] = parsed
	}

	return result, nil
}

// CreateInstance takes a type and a Gherkin table and returns an instance of
// that type filled with the table's parsed values.
// The table must have exactly two columns, where the first represents the field names
// and the second represents the values.
func (a *Assist) CreateInstance(tp interface{}, table *godog.Table) (interface{}, error) {
	tableMap, err := a.ParseMap(table)
	if err != nil {
		return nil, err
	}

	instance, errs := a.createInstance(tp, tableMap)
	if len(errs) != 0 {
		return nil, fmt.Errorf("failed to parse table as %v:\n- %v", reflect.TypeOf(tp), strings.Join(errs, "\n- "))
	}

	return instance.Interface(), nil
}

// CreateSlice takes a type and a Gherkin table and returns a slice of that type
// filled with each row as an instance.
// The first row acts as a header and provides the field names for each column.
func (a *Assist) CreateSlice(tp interface{}, table *godog.Table) (interface{}, error) {
	maps, err := a.ParseSlice(table)
	if err != nil {
		return nil, err
	}

	errs := []string{}
	slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(tp)), 0, len(maps))
	for i, row := range maps {
		instance, fieldErrors := a.createInstance(tp, row)
		if len(fieldErrors) > 0 {
			errs = append(errs, fmt.Sprintf("row %v:\n  - %v", i, strings.Join(fieldErrors, "\n  - ")))
			continue
		}

		slice = reflect.Append(slice, instance)
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to parse table as slice of %v:\n%v", reflect.TypeOf(tp), strings.Join(errs, "\n"))
	}

	return slice.Interface(), nil
}

// CompareToInstance compares an actual value to the expected fields from a Gherkin table.
func (a *Assist) CompareToInstance(actual interface{}, table *godog.Table) error {
	tableMap, err := a.ParseMap(table)
	if err != nil {
		return err
	}

	errs := a.compareToInstance(actual, tableMap)
	if len(errs) != 0 {
		return fmt.Errorf("comparison failed:\n- %v", strings.Join(errs, "\n- "))
	}

	return nil
}

// CompareToSlice compares an actual slice of values to the expected rows from a Gherkin table.
func (a *Assist) CompareToSlice(actual interface{}, table *godog.Table) error {
	maps, err := a.ParseSlice(table)
	if err != nil {
		return err
	}

	actualValue := reflect.ValueOf(actual)
	if actualValue.Kind() != reflect.Slice {
		return fmt.Errorf("actual value is not a slice")
	}

	errs := []string{}
	for i, row := range maps {
		rowErrs := a.compareToInstance(actualValue.Index(i).Interface(), row)
		if len(rowErrs) > 0 {
			errs = append(errs, fmt.Sprintf("row %v:\n  - %v", i, strings.Join(rowErrs, "\n  - ")))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("comparison failed:\n%v", strings.Join(errs, "\n"))
	}

	return nil
}

func (a *Assist) createInstance(tp interface{}, table map[string]string) (reflect.Value, []string) {
	errs := []string{}
	result := reflect.New(reflect.TypeOf(tp).Elem())
	sv := result.Elem()
	for fieldName, rawValue := range table {
		fv := sv.FieldByName(fieldName)
		if !fv.IsValid() {
			errs = append(errs, fmt.Sprintf("%v: field not found", fieldName))
			continue
		}

		if !fv.CanSet() {
			errs = append(errs, fmt.Sprintf("%v: cannot set value", fieldName))
			continue
		}

		parseField, ok := a.findParser(fv.Type())
		if !ok {
			errs = append(errs, fmt.Sprintf("%v: unrecognized type %v", fieldName, fv.Type()))
			continue
		}

		parsed, err := parseField(rawValue)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%v: %v", fieldName, err.Error()))
			continue
		}

		fv.Set(reflect.ValueOf(parsed))
	}

	return result, errs
}

func (a *Assist) compareToInstance(actual interface{}, table map[string]string) []string {
	errs := []string{}
	sv := reflect.ValueOf(actual).Elem()
	for fieldName, rawExpectedValue := range table {
		fv := sv.FieldByName(fieldName)
		if !fv.IsValid() {
			errs = append(errs, fmt.Sprintf("%v: field not found", fieldName))
			continue
		}

		compare, ok := a.findComparer(fv.Type())
		if !ok {
			errs = append(errs, fmt.Sprintf("%v: unrecognized type %v", fieldName, fv.Type()))
			continue
		}

		if err := compare(rawExpectedValue, fv.Interface()); err != nil {
			errs = append(errs, fmt.Sprintf("%v: %v", fieldName, err))
		}
	}

	return errs
}

func (a *Assist) findParser(tp reflect.Type) (ParseFunc, bool) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	p, ok := a.parsers[tp]
	return p, ok
}

func (a *Assist) findComparer(tp reflect.Type) (CompareFunc, bool) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	c, ok := a.comparers[tp]
	return c, ok
}

func (a *Assist) assertInit() {
	if a.parsers == nil {
		a.parsers = map[reflect.Type]ParseFunc{}
	}

	if a.comparers == nil {
		a.comparers = map[reflect.Type]CompareFunc{}
	}
}
