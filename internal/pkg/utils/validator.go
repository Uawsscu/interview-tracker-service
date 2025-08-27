package utils

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func ValidateRequest(request interface{}) map[string]string {
	if err := newValidator().Struct(request); err != nil {
		return validatorErrors(err)
	}
	return nil
}

// newValidator func for creating a new validator for model fields.
func newValidator() *validator.Validate {
	validate := validator.New()
	_ = validate.RegisterValidation("notEmpty", func(fl validator.FieldLevel) bool {
		return notEmpty(fl)
	})
	_ = validate.RegisterValidation("thaiLanguage", func(fl validator.FieldLevel) bool {
		return validateThaiLanguage(fl)
	})
	_ = validate.RegisterValidation("englishAlphabet", func(fl validator.FieldLevel) bool {
		return validateEnglishAlphabet(fl)
	})
	_ = validate.RegisterValidation("uuid", func(fl validator.FieldLevel) bool {
		return ValidateUUID(fl)
	})

	return validate
}

// validatorErrors func for showing validation errors for each invalid field.
func validatorErrors(err error) map[string]string {
	fields := map[string]string{}
	for _, err := range err.(validator.ValidationErrors) {
		fields[err.Field()] = err.Field() + " " + err.Tag() + " " + err.Param()
	}

	return fields
}

// notEmpty supports string, json, array and map types
func notEmpty(fl validator.FieldLevel) bool {
	field := fl.Field()
	fieldKind := field.Kind()
	if fieldKind == reflect.Array || fieldKind == reflect.Slice || fieldKind == reflect.Map {
		return field.Len() > 0
	}
	trimmedValue := strings.TrimSpace(field.Interface().(string))
	return len(trimmedValue) > 0
}

func validateThaiLanguage(fl validator.FieldLevel) bool {
	// Define a regex pattern for Thai language characters
	pattern := "([\u0E00-\u0E7F]+)"

	// Compile the regex pattern
	regex := regexp.MustCompile(pattern)

	// Use the regex to match the field's value
	return regex.MatchString(fl.Field().String())
}

func validateEnglishAlphabet(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	// Define a regular expression pattern that matches English alphabet characters only
	pattern := "^[a-zA-Z]+$"
	match, _ := regexp.MatchString(pattern, value)

	return match
}

func ValidateUUID(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() == reflect.String {
		id := field.String()
		if _, err := uuid.Parse(id); err != nil {
			return false
		}
	} else if field.Kind() == reflect.Slice {
		for i := 0; i < field.Len(); i++ {
			id := field.Index(i).String()
			if _, err := uuid.Parse(id); err != nil {
				return false
			}
		}
	} else {
		return false
	}

	return true
}

func ValidateStringParamUUID(uuid string) bool {
	pattern := `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`
	matched, _ := regexp.MatchString(pattern, uuid)
	return matched
}
