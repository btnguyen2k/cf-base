package cfbase

import (
	"reflect"
	"regexp"
)

var _typMapString = reflect.TypeOf(map[string]string{})
var _typMapAny = reflect.TypeOf(map[string]any{})

var _rexpContentDir = regexp.MustCompile(`^(\d+)-([\w-]+)$`)
