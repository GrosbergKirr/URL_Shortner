package ms

import (
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/exp/slog"
)

func Err(err error) slog.Attr {
	return slog.Attr{Key: "error",
		Value: slog.StringValue(err.Error())}
}
