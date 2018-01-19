package defaults

import (
	"fmt"
	"strconv"
	"time"
)

// CompareString function
func CompareString(raw string, actual interface{}) error {
	as, ok := actual.(string)
	if !ok {
		return fmt.Errorf("%v is not a string", actual)
	}

	if as != raw {
		return fmt.Errorf("expected %v, but got %v", raw, as)
	}

	return nil
}

// CompareInt function
func CompareInt(raw string, actual interface{}) error {
	ai, ok := actual.(int)
	if !ok {
		return fmt.Errorf("%v is not an int", actual)
	}

	ei, err := strconv.Atoi(raw)
	if err != nil {
		return err
	}

	if ei != ai {
		return fmt.Errorf("expected %v, but got %v", ei, ai)
	}

	return nil
}

// CompareFloat32 function
func CompareFloat32(raw string, actual interface{}) error {
	as, ok := actual.(float32)
	if !ok {
		return fmt.Errorf("%v is not a float32", actual)
	}

	ei, err := strconv.ParseFloat(raw, 32)
	if err != nil {
		return err
	}

	if as != float32(ei) {
		return fmt.Errorf("expected %v, but got %v", raw, as)
	}

	return nil
}

// CompareTime function
func CompareTime(raw string, actual interface{}) error {
	at, ok := actual.(time.Time)
	if !ok {
		return fmt.Errorf("%v is not time.Time", actual)
	}

	et, err := ParseTime(raw)
	if err != nil {
		return err
	}

	if !at.Equal(et.(time.Time)) {
		return fmt.Errorf("expected %v, but got %v", et, at)
	}

	return nil
}

// CompareStringPointer function
func CompareStringPointer(raw string, actual interface{}) error {
	val, ok := actual.(*string)

	if !ok {
		return fmt.Errorf("%v is not *string", actual)
	}

	if val == nil {
		if raw == NilRawString {
			return nil
		}
		return fmt.Errorf("expected %v, but got nil", raw)
	}

	return CompareString(raw, *val)
}

// CompareIntPointer function
func CompareIntPointer(raw string, actual interface{}) error {
	val, ok := actual.(*int)

	if !ok {
		return fmt.Errorf("%v is not *int", actual)
	}

	if val == nil {
		if raw == NilRawString {
			return nil
		}
		return fmt.Errorf("expected %v, but got nil", raw)
	}

	return CompareInt(raw, *val)
}

// CompareFloat32Pointer function
func CompareFloat32Pointer(raw string, actual interface{}) error {
	val, ok := actual.(*float32)

	if !ok {
		return fmt.Errorf("%v is not *float32", actual)
	}

	if val == nil {
		if raw == NilRawString {
			return nil
		}
		return fmt.Errorf("expected %v, but got nil", raw)
	}

	return CompareFloat32(raw, *val)
}

// CompareTimePointer function
func CompareTimePointer(raw string, actual interface{}) error {
	val, ok := actual.(*time.Time)

	if !ok {
		return fmt.Errorf("%v is not *time.Time", actual)
	}

	if val == nil {
		if raw == NilRawString {
			return nil
		}
		return fmt.Errorf("expected %v, but got nil", raw)
	}
	return CompareTime(raw, *val)
}
