package main

import "strconv"

type Type int

const (
	STRING Type = iota
	INT
	FLOAT
	VARIADIC
)

func Match(typs []Type, vals []string) ([]interface{}, Response) {
	if typs[len(typs)-1] == VARIADIC {
		if len(vals) < len(typs)-1 {
			return nil, R("Expected at least %d arguments, got %d", len(typs), len(vals))
		}
	} else if len(vals) != len(typs) {
		return nil, R("Expected %d arguments, got %d", len(typs), len(vals))
	}

	out := make([]interface{}, 0, len(typs))
	for i, val := range vals {
		var typ Type
		if typs[len(typs)-1] == VARIADIC && i >= len(typs)-1 {
			typ = typs[len(typs)-2]
		} else {
			typ = typs[i]
		}

		switch typ {
		case STRING:
			out = append(out, val)

		case INT:
			v, err := strconv.Atoi(val)
			if err != nil {
				return nil, R("%s", err.Error())
			}
			out = append(out, v)

		case FLOAT:
			v, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, R("%s", err.Error())
			}
			out = append(out, v)
		}
	}
	return out, Rg()
}
