package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Error describes a validation error.
type Error struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

// UseJSONFieldNames registers a custom tag name fun to gin's validator that uses the json field name instead
// of struct field name when throwing a validation error.
func UseJSONFieldNames() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

// DescriptiveErrors converts ValidationErrors into a format that the client can work with.
func DescriptiveErrors(verr validator.ValidationErrors) []Error {
	errs := []Error{}

	for _, f := range verr {
		err := f.ActualTag()
		if f.Param() != "" {
			err = fmt.Sprintf("%s=%s", err, f.Param())
		}
		errs = append(errs, Error{Field: f.Field(), Reason: err})
	}

	return errs
}
