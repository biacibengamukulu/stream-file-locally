package utils

import (
	"os"
	"strings"

	"github.com/google/uuid"
)

func IsValidUUID(s string) bool {
	_, errValidation := uuid.Parse(s)
	return errValidation == nil
}
func Clean1(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, " ", "")
	input = strings.ReplaceAll(input, "\n", "")
	input = strings.ReplaceAll(input, "\r", "")
	input = strings.ReplaceAll(input, "\t", "")
	return input
}
func CompareString(str1 string, str2 string) bool {
	return Clean1(str1) == Clean1(str2)
}
func Getenv(k, d string) string {
	v := os.Getenv(k)
	if v == "" {
		return d
	}
	return v
}
