package defaults

import (
	"fmt"
	"strconv"
	"time"
)

var supportedTimeLayouts = []string{
	time.RFC822,
	time.RFC3339,
	time.RFC3339Nano,
}

func ParseString(raw string) (interface{}, error) {
	return raw, nil
}

func ParseInt(raw string) (interface{}, error) {
	return strconv.Atoi(raw)
}

func ParseInt32(raw string) (interface{}, error) {
	i, err := strconv.ParseInt(raw, 10, 32)
	if err != nil {
		return nil, err
	}
	return int32(i), nil
}

func ParseTime(raw string) (interface{}, error) {
	var fieldTime time.Time
	var err error
	for _, layout := range supportedTimeLayouts {
		fieldTime, err = time.Parse(layout, raw)
		if err != nil {
			continue
		}

		break
	}

	if err != nil {
		return nil, fmt.Errorf("unrecognized time format %v", raw)
	}

	return fieldTime, nil
}
