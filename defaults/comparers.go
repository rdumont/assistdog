package defaults

import (
	"fmt"
	"strconv"
	"time"
)

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

func CompareInt32(raw string, actual interface{}) error {
	ai, ok := actual.(int32)
	if !ok {
		return fmt.Errorf("%v is not an int32", actual)
	}

	i, err := strconv.ParseInt(raw, 10, 32)
	if err != nil {
		return err
	}
	ei := int32(i)

	if ei != ai {
		return fmt.Errorf("expected %v, but got %v", ei, ai)
	}

	return nil
}

func CompareTime(raw string, actual interface{}) error {
	at, ok := actual.(time.Time)
	if !ok {
		return fmt.Errorf("%v is not time.Time", actual)
	}

	et, err := ParseTime(raw)
	if err != nil {
		return err
	}

	if et != at {
		return fmt.Errorf("expected %v, but got %v", et, at)
	}

	return nil
}
