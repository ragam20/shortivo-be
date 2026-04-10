package repo

import (
	"strings"
)

func parseSort(sort string, defaultSort string, sortMap map[string]string) (string, error) {
	orders := strings.Split(sort, ",")

	var sb strings.Builder
	first := true

	for _, order := range orders {
		if order == "" {
			continue
		}

		direction := "ASC"
		key := order

		if strings.HasPrefix(order, "-") {
			direction = "DESC"
			key = strings.TrimPrefix(order, "-")
		}

		column, ok := sortMap[key]
		if !ok {
			continue
		}

		if !first {
			sb.WriteString(", ")
		}

		sb.WriteString(column)
		sb.WriteString(" ")
		sb.WriteString(direction)

		first = false
	}

	// default
	if sb.Len() == 0 {
		return defaultSort, nil
	}

	return sb.String(), nil
}
