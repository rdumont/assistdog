package assistdog

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/DATA-DOG/godog/gherkin"
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

// NewDefault creates a new Assist instance with all the default parsers and comparers
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

type Assist struct {
	lock      sync.RWMutex
	parsers   map[reflect.Type]ParseFunc
	comparers map[reflect.Type]CompareFunc
}

func (a *Assist) RegisterParser(i interface{}, parser ParseFunc) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.assertInit()
	a.parsers[reflect.TypeOf(i)] = parser
}

func (a *Assist) RegisterComparer(i interface{}, comparer CompareFunc) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.assertInit()
	a.comparers[reflect.TypeOf(i)] = comparer
}

func (a *Assist) RemoveParser(i interface{}) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.assertInit()
	delete(a.parsers, reflect.TypeOf(i))
}

func (a *Assist) RemoveComparer(i interface{}) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.assertInit()
	delete(a.comparers, reflect.TypeOf(i))
}

func (a *Assist) ParseMap(table *gherkin.DataTable) (map[string]string, error) {
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

func (a *Assist) ParseSlice(table *gherkin.DataTable) ([]map[string]string, error) {
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

func (a *Assist) CreateInstance(tp interface{}, table *gherkin.DataTable) (interface{}, error) {
	tableMap, err := a.ParseMap(table)
	if err != nil {
		return nil, err
	}

	fieldErrors := []string{}

	instance := reflect.ValueOf(tp)
	sv := instance.Elem()
	for fieldName, rawValue := range tableMap {
		fv := sv.FieldByName(fieldName)
		if !fv.IsValid() {
			fieldErrors = append(fieldErrors, fmt.Sprintf("%v: field not found", fieldName))
			continue
		}

		if !fv.CanSet() {
			fieldErrors = append(fieldErrors, fmt.Sprintf("%v: cannot set value", fieldName))
			continue
		}

		parseField, ok := a.findParser(fv.Type())
		if !ok {
			fieldErrors = append(fieldErrors, fmt.Sprintf("%v: unrecognized type %v", fieldName, fv.Type()))
			continue
		}

		parsed, err := parseField(rawValue)
		if err != nil {
			fieldErrors = append(fieldErrors, fmt.Sprintf("%v: %v", fieldName, err.Error()))
			continue
		}

		fv.Set(reflect.ValueOf(parsed))
	}

	if len(fieldErrors) != 0 {
		return nil, fmt.Errorf("failed to parse table as %v:\n- %v", sv.Type(), strings.Join(fieldErrors, "\n- "))
	}

	return instance.Interface(), nil
}

func (a *Assist) CreateSlice(tp interface{}, table *gherkin.DataTable) (interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

func (a *Assist) CompareToInstance(actual interface{}, table *gherkin.DataTable) error {
	tableMap, err := a.ParseMap(table)
	if err != nil {
		return err
	}

	fieldErrors := []string{}
	sv := reflect.ValueOf(actual).Elem()
	for fieldName, rawExpectedValue := range tableMap {
		fv := sv.FieldByName(fieldName)
		if !fv.IsValid() {
			fieldErrors = append(fieldErrors, fmt.Sprintf("%v: field not found", fieldName))
			continue
		}

		compare, ok := a.findComparer(fv.Type())
		if !ok {
			fieldErrors = append(fieldErrors, fmt.Sprintf("%v: unrecognized type %v", fieldName, fv.Type()))
			continue
		}

		if err := compare(rawExpectedValue, fv.Interface()); err != nil {
			fieldErrors = append(fieldErrors, fmt.Sprintf("%v: %v", fieldName, err))
		}
	}

	if len(fieldErrors) != 0 {
		return fmt.Errorf("comparison failed:\n- %v", strings.Join(fieldErrors, "\n- "))
	}

	return nil
}

func (a *Assist) CompareToSlice(actual interface{}, table *gherkin.DataTable) error {
	return fmt.Errorf("not implemented")
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
