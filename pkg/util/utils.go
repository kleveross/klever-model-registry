package util

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/caicloud/nirvana/service"
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

func GetFormValueFromRequest(ctx context.Context, field string) string {
	return service.HTTPContextFrom(ctx).Request().FormValue(field)
}

func GetRequstFromContext(ctx context.Context) *http.Request {
	return service.HTTPContextFrom(ctx).Request()
}

func GetResponseFromContext(ctx context.Context) http.ResponseWriter {
	return service.HTTPContextFrom(ctx).ResponseWriter()
}
