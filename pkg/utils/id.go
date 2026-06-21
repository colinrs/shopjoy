package utils

import (
	"strconv"
)

// ParseInt64Slice converts a []string of int64-as-string to []int64.
// Used at the API boundary to convert frontend string IDs back to int64
// for the domain layer. Returns the first parse error encountered.
func ParseInt64Slice(ss []string) ([]int64, error) {
	ids := make([]int64, 0, len(ss))
	for _, s := range ss {
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// FormatInt64Slice converts []int64 to []string for JSON serialization.
// Used at the API boundary to avoid JS Number overflow with snowflake IDs.
func FormatInt64Slice(ids []int64) []string {
	result := make([]string, 0, len(ids))
	for _, id := range ids {
		result = append(result, strconv.FormatInt(id, 10))
	}
	return result
}
