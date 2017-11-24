package trimet

import (
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func parseDuration(s string) (*Time, error) {
	if s == "" {
		return nil, nil
	}
	parts := strings.Split(s, ":")
	if len(parts) != 3 {
		return nil, errors.Errorf("expected 3 parts, found %d", len(parts))
	}

	var intParts [3]time.Duration
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		intParts[i] = time.Duration(n)
	}
	dur := intParts[0]*time.Hour + intParts[1]*time.Minute + intParts[2]*time.Second
	return (*Time)(&dur), nil
}

func parseFloat(s string, defaultValue float64) (float64, error) {
	if s == "" {
		return defaultValue, nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return f, nil
}

func parseInt(s string, defaultValue int) (int, error) {
	if s == "" {
		return defaultValue, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return n, nil
}

func parseNullableFloat(s string) (*float64, error) {
	if s == "" {
		return nil, nil
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &f, nil
}

func parseNullableInt(s string) (*int, error) {
	if s == "" {
		return nil, nil
	}

	n, err := strconv.Atoi(s)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &n, nil
}

func parseNullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
