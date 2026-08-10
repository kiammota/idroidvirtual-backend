package middlewares

import (
	"reflect"
)

func isEmptyJSON(data interface{}) bool {

	if data == nil {
		return true
	}

	value := reflect.ValueOf(data)

	switch value.Kind() {

	case reflect.Map:
		if value.Len() == 0 {
			return true
		}

		for _, key := range value.MapKeys() {
			item := value.MapIndex(key)

			if item.Interface() != nil {
				return false
			}
		}

		return true

	case reflect.Slice:
		return value.Len() == 0

	default:
		return false
	}
}
