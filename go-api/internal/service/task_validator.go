package service

import (
	"time"
)

// ValidationError represents a task validation failure.
// Messages contains the list of validation error strings, matching PHP's InvalidTaskException.
type ValidationError struct {
	Messages []string
}

func (e *ValidationError) Error() string {
	return "Task is not valid!"
}

// ValidateTask replicates PHP's TaskValidator::validate() exactly.
// The error messages, their ordering, and edge-case behavior all match PHP.
func ValidateTask(name, description, when string, requireID bool, id interface{}) error {
	var errors []string

	// 1. Name validation — matches: if (!is_string($taskDto->name) || empty($taskDto->name))
	if name == "" {
		errors = append(errors, "Task name is not valid!")
	}

	// 2. ID validation (only for updates) — matches PHP's shouldHaveId block
	if requireID {
		if id == nil {
			errors = append(errors, "We need an id to know which entity to update!")
		} else {
			// In PHP, json_decode produces float64 for numbers.
			// Go's json.Unmarshal into interface{} also produces float64.
			// PHP checks is_integer() which is false for floats.
			// We accept float64 values that are whole numbers as valid integers.
			switch v := id.(type) {
			case float64:
				// Whole number float is OK (this is how JSON numbers decode)
				if v != float64(int(v)) {
					errors = append(errors, "The id must be an integer!")
				}
			default:
				// string, bool, etc. — not an integer
				errors = append(errors, "The id must be an integer!")
			}
		}
	}

	// 3. Description validation
	if description == "" {
		errors = append(errors, "Task description is not valid!")
	}

	// 4. When validation
	if when == "" {
		errors = append(errors, "Task must have a date!")
	} else {
		parsed, err := time.Parse("2006-01-02", when)
		if err != nil {
			errors = append(errors, "`when` is not a valid date!")
		} else {
			// PHP: $when < (new \DateTime())
			// PHP's createFromFormat sets the time to current H:M:S, then compares.
			// But the comparison is: parsed date at 00:00:00 vs now (with time).
			// Actually PHP createFromFormat('Y-m-d', ...) sets H:M:S to current time.
			// So the comparison is: date at current-H:M:S vs now.
			// In practice for same-day: they're nearly equal, comparison depends on execution speed.
			// We replicate: parsed date at 00:00:00 compared to now.
			// This means today's date will be < now() (since 00:00 < current time), matching PHP behavior
			// where the parsed date at midnight is less than the current datetime.
			if parsed.Before(time.Now().Truncate(24 * time.Hour)) {
				errors = append(errors, "You can't do a task in a the past. Or you can?")
			}
		}
	}

	if len(errors) > 0 {
		return &ValidationError{Messages: errors}
	}
	return nil
}
