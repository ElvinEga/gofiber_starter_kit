package utils

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationResult holds the result of validation
type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Errors []ValidationError `json:"errors,omitempty"`
}

// Validator provides validation functionality
type Validator struct {
	// You can add custom validation functions here if needed
	customValidators map[string]func(interface{}) error
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		customValidators: make(map[string]func(interface{}) error),
	}
}

// AddCustomValidator adds a custom validation function
func (v *Validator) AddCustomValidator(name string, fn func(interface{}) error) {
	v.customValidators[name] = fn
}

// Validate performs validation on the given struct
func (v *Validator) Validate(s interface{}) ValidationResult {
	result := ValidationResult{
		Valid:  true,
		Errors: []ValidationError{},
	}

	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()

	// Iterate through all fields
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Get the json tag name for the field
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = strings.ToLower(fieldType.Name)
		}
		// Remove any json tag options (like ,omitempty)
		jsonTag = strings.Split(jsonTag, ",")[0]

		// Check required fields
		if v.isRequired(fieldType) && v.isEmpty(field) {
			result.Errors = append(result.Errors, ValidationError{
				Field:   jsonTag,
				Message: "This field is required",
			})
			result.Valid = false
			continue
		}

		// Skip validation for empty optional fields
		if v.isEmpty(field) {
			continue
		}

		// Validate field based on its type
		if err := v.validateField(field, fieldType); err != nil {
			result.Errors = append(result.Errors, ValidationError{
				Field:   jsonTag,
				Message: err.Error(),
			})
			result.Valid = false
		}
	}

	return result
}

// isRequired checks if a field is required based on struct tags
func (v *Validator) isRequired(field reflect.StructField) bool {
	// Check for required tag
	if _, ok := field.Tag.Lookup("required"); ok {
		return true
	}

	// Check form tag
	formTag := field.Tag.Get("form")
	return strings.Contains(formTag, "required")
}

// isEmpty checks if a field is empty
func (v *Validator) isEmpty(field reflect.Value) bool {
	switch field.Kind() {
	case reflect.String:
		return field.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() == 0
	case reflect.Float32, reflect.Float64:
		return field.Float() == 0
	case reflect.Struct:
		if field.Type() == reflect.TypeOf(time.Time{}) {
			return field.Interface().(time.Time).IsZero()
		}
		return false
	case reflect.Ptr:
		return field.IsNil()
	default:
		return false
	}
}

// validateField validates a single field based on its type and tags
func (v *Validator) validateField(field reflect.Value, fieldType reflect.StructField) error {
	switch field.Interface().(type) {
	case uuid.UUID:
		if field.Interface().(uuid.UUID) == uuid.Nil {
			return fmt.Errorf("invalid UUID format")
		}
	case string:
		if maxLen, ok := fieldType.Tag.Lookup("max"); ok {
			// Add string length validation
			_ = maxLen // You can implement this
		}
	case float64:
		if min, ok := fieldType.Tag.Lookup("min"); ok {
			// Add minimum value validation
			_ = min // You can implement this
		}
	case time.Time:
		// Add any specific time validation if needed
		if field.Interface().(time.Time).IsZero() {
			return fmt.Errorf("invalid date format")
		}
	}

	// Run custom validators if any
	if validatorName := fieldType.Tag.Get("validator"); validatorName != "" {
		if validator, ok := v.customValidators[validatorName]; ok {
			if err := validator(field.Interface()); err != nil {
				return err
			}
		}
	}

	return nil
}
