package consul

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/errwrap"
)

// An array of inputs used as typed arguments and converted from their type into
// function objects that are dynamically constructed and executed.
type validatorInputs []interface{}

// validateDurationMin is the minimum duration to accept as input
type validateDurationMin string

// validateIntMin is the minimum integer value to accept as input
type validateIntMin int

// validateRegexp is a regexp pattern to use to validate schema input.
type validateRegexp string

// makeValidateionFunc takes the name of the attribute and a list of typed
// validator inputs in order to create a validation closure that calls each
// validator in serial until either a warning or error is returned from the
// first validation function.
func makeValidationFunc(name string, validators []interface{}) func(v interface{}, key string) (warnings []string, errors []error) {
	if len(validators) == 0 {
		return nil
	}

	fns := make([]func(v interface{}, key string) (warnings []string, errors []error), 0, len(validators))
	for _, v := range validators {
		switch u := v.(type) {
		case validateDurationMin:
			fns = append(fns, validateDurationMinFactory(name, string(u)))
		case validateIntMin:
			fns = append(fns, validateIntMinFactory(name, int(u)))
		case validateRegexp:
			fns = append(fns, validateRegexpFactory(name, string(u)))
		}
	}

	return func(v interface{}, key string) (warnings []string, errors []error) {
		for _, fn := range fns {
			warnings, errors = fn(v, key)
			if len(warnings) > 0 || len(errors) > 0 {
				break
			}
		}
		return warnings, errors
	}
}

func validateDurationMinFactory(name, minDuration string) func(v interface{}, key string) (warnings []string, errors []error) {
	dMin, err := time.ParseDuration(minDuration)
	if err != nil {
		return func(interface{}, string) (warnings []string, errors []error) {
			return nil, []error{
				errwrap.Wrapf(fmt.Sprintf("PROVIDER BUG: duration %q not valid: {{err}}", minDuration), err),
			}
		}
	}

	return func(v interface{}, key string) (warnings []string, errors []error) {
		d, err := time.ParseDuration(v.(string))
		if err != nil {
			errors = append(errors, errwrap.Wrapf(fmt.Sprintf("Invalid %s specified (%q): {{err}}", name), err))
		}

		if d < dMin {
			errors = append(errors, fmt.Errorf("Invalid %s specified: duration %q less than the required minimum %s", name, v.(string), dMin))
		}

		return warnings, errors
	}
}

func validateIntMinFactory(name string, min int) func(v interface{}, key string) (warnings []string, errors []error) {
	return func(v interface{}, key string) (warnings []string, errors []error) {
		if v.(int) < min {
			errors = append(errors, fmt.Errorf("Invalid %s specified: %d less than the required minimum %d", name, v.(int), min))
		}

		return warnings, errors
	}
}

func validateRegexpFactory(name string, reString string) func(v interface{}, key string) (warnings []string, errors []error) {
	re := regexp.MustCompile(reString)

	return func(v interface{}, key string) (warnings []string, errors []error) {
		if !re.MatchString(v.(string)) {
			errors = append(errors, fmt.Errorf("Invalid %s specified (%q): regexp failed to match string", name, v.(string)))
		}

		return warnings, errors
	}
}
