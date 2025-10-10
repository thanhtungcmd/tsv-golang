package direction

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"

	"github.com/99designs/gqlgen/graphql"
)

func Validate(
	ctx context.Context, obj any, next graphql.Resolver, required *bool, minLength *int, maxLength *int, pattern *string,
) (interface{}, error) {
	// Resolve the value (value after default/unmarshal)
	val, err := next(ctx)
	if err != nil {
		return nil, err
	}

	// Helper to get field name for error messages
	fieldName := "field"
	if pc := graphql.GetPathContext(ctx); pc != nil && *pc.Field != "" {
		fieldName = *pc.Field
	}

	// Check required
	if required != nil && *required {
		if val == nil {
			return nil, fmt.Errorf("field '%s' is required", fieldName)
		}
		// empty string check as well
		if s, ok := val.(string); ok && len(s) == 0 {
			return nil, fmt.Errorf("field '%s' is required", fieldName)
		}
	}

	// If value is nil, and not required, skip further checks
	if val == nil {
		return val, nil
	}

	rv := reflect.ValueOf(val)
	kind := rv.Kind()

	// If pointer, deref
	if kind == reflect.Ptr {
		if rv.IsNil() {
			return val, nil
		}
		rv = rv.Elem()
		kind = rv.Kind()
		val = rv.Interface()
	}

	// Validate string length and pattern
	if kind == reflect.String {
		s := rv.String()
		if minLength != nil {
			if len(s) < *minLength {
				return nil, fmt.Errorf("field '%s' length must be >= %d", fieldName, *minLength)
			}
		}
		if maxLength != nil {
			if len(s) > *maxLength {
				return nil, fmt.Errorf("field '%s' length must be <= %d", fieldName, *maxLength)
			}
		}
		if pattern != nil && *pattern != "" {
			re, err := regexp.Compile(*pattern)
			if err != nil {
				return nil, errors.New("server: invalid validation pattern")
			}
			if !re.MatchString(s) {
				return nil, fmt.Errorf("field '%s' does not match required format", fieldName)
			}
		}
		return val, nil
	}

	return val, nil
}
