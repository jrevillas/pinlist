package models

import (
	"fmt"
	"strings"
)

func inQuery(query string, params []int64) string {
	var strParams = make([]string, len(params))
	for i, p := range params {
		strParams[i] = fmt.Sprint(p)
	}

	return fmt.Sprintf(query, "("+strings.Join(strParams, ", ")+")")
}
