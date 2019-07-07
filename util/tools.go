package util

import (
	"github.com/rs/xid"
)

func UuidProvide() string {
	return xid.New().String()
}
