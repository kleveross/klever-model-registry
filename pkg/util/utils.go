package util

import (
	"fmt"
	"strings"
	"time"

	utilrand "k8s.io/apimachinery/pkg/util/rand"
)

const (
	timeFormat = "20060102150405"
)

func RandomNameWithPrefix(prefix string) string {
	return fmt.Sprintf("%s-%s", strings.ToLower(prefix), RandomNameWithTS(8))
}

func RandomNameWithTS(n int) string {
	return strings.ToLower(
		fmt.Sprintf(
			"%s-%s",
			time.Now().Format(timeFormat),
			utilrand.String(n)))
}
