// Package field provides types and parsing logic for individual cron expression fields.
package field

import (
	"fmt"
	"strconv"
	"strings"
)

// Kind represents which positional field in a cron expression this is.
type Kind int

const (
	Minute Kind = iota
	Hour
	DayOfMonth
	Month
	DayOfWeek
)

// bounds holds the inclusive min/max for a given field kind.
var bounds = map[Kind][2]int{
	Minute:     {0, 59},
	Hour:       {0, 23},
	DayOfMonth: {1, 31},
	Month:      {1, 12},
	DayOfWeek:  {0, 6},
}

// Field represents a parsed cron field with its resolved set of values.
type Field struct {
	Kind   Kind
	Values []int // sorted, deduplicated set of matched values
}

// Parse parses a single cron field token (e.g. "*", "1-5", "*/2", "1,3,5") for the given kind.
func Parse(token string, kind Kind) (*Field, error) {
	min, max := bounds[kind][0], bounds[kind][1]
	values, err := parseToken(token, min, max)
	if err != nil {
		return nil, fmt.Errorf("field %d: %w", kind, err)
	}
	return &Field{Kind: kind, Values: values}, nil
}

func parseToken(token string, min, max int) ([]int, error) {
	switch {
	case token == "*":
		return rangeSlice(min, max, 1), nil
	case strings.Contains(token, ","):
		return parseList(token, min, max)
	case strings.Contains(token, "/"):
		return parseStep(token, min, max)
	case strings.Contains(token, "-"):
		return parseRange(token, min, max)
	default:
		return parseLiteral(token, min, max)
	}
}

func parseList(token string, min, max int) ([]int, error) {
	parts := strings.Split(token, ",")
	seen := map[int]struct{}{}
	var result []int
	for _, p := range parts {
		vals, err := parseToken(strings.TrimSpace(p), min, max)
		if err != nil {
			return nil, err
		}
		for _, v := range vals {
			if _, ok := seen[v]; !ok {
				seen[v] = struct{}{}
				result = append(result, v)
			}
		}
	}
	return result, nil
}

func parseRange(token string, min, max int) ([]int, error) {
	parts := strings.SplitN(token, "-", 2)
	lo, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid range start %q", parts[0])
	}
	hi, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid range end %q", parts[1])
	}
	if lo < min || hi > max || lo > hi {
		return nil, fmt.Errorf("range %d-%d out of bounds [%d,%d]", lo, hi, min, max)
	}
	return rangeSlice(lo, hi, 1), nil
}

func parseStep(token string, min, max int) ([]int, error) {
	parts := strings.SplitN(token, "/", 2)
	step, err := strconv.Atoi(parts[1])
	if err != nil || step <= 0 {
		return nil, fmt.Errorf("invalid step %q", parts[1])
	}
	var lo, hi int
	if parts[0] == "*" {
		lo, hi = min, max
	} else if strings.Contains(parts[0], "-") {
		rng, err := parseRange(parts[0], min, max)
		if err != nil {
			return nil, err
		}
		lo, hi = rng[0], rng[len(rng)-1]
	} else {
		lo, err = strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid step base %q", parts[0])
		}
		hi = max
	}
	return rangeSlice(lo, hi, step), nil
}

func parseLiteral(token string, min, max int) ([]int, error) {
	v, err := strconv.Atoi(token)
	if err != nil {
		return nil, fmt.Errorf("invalid value %q", token)
	}
	if v < min || v > max {
		return nil, fmt.Errorf("value %d out of bounds [%d,%d]", v, min, max)
	}
	return []int{v}, nil
}

func rangeSlice(lo, hi, step int) []int {
	var s []int
	for i := lo; i <= hi; i += step {
		s = append(s, i)
	}
	return s
}
