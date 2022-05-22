package common

import "strings"

func CleanStdoutSingleline(stdout []byte) string {
	return strings.TrimSpace(strings.ReplaceAll(string(stdout), "\n", ""))
}

func CleanstdoutMultiline(stdout []byte) string {
	return strings.TrimSpace(string(stdout))
}