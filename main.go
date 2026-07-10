package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func Validate(v any) error {
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	// Проверяем, что передан указатель на структуру или структура
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct, got %s", val.Kind())
	}

	// Проходим по всем полям структуры
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// Извлекаем тег validate
		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}

		// Разбиваем правила по точке с запятой
		rules := strings.Split(tag, ";")

		for _, rule := range rules {
			rule = strings.TrimSpace(rule)
			if rule == "" {
				continue
			}

			// Разбираем правило (ключ=значение)
			parts := strings.Split(rule, "=")
			if len(parts) != 2 {
				return fmt.Errorf("invalid rule format: %s", rule)
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Проверяем в зависимости от типа поля и правила
			switch fieldValue.Kind() {
			case reflect.String:
				err := validateString(field.Name, fieldValue.String(), key, value)
				if err != nil {
					return err
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				err := validateInt(field.Name, int(fieldValue.Int()), key, value)
				if err != nil {
					return err
				}
			case reflect.Float32, reflect.Float64:
				err := validateFloat(field.Name, fieldValue.Float(), key, value)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func validateString(fieldName, value, rule, param string) error {
	switch rule {
	case "min":
		minLen, err := strconv.Atoi(param)
		if err != nil {
			return fmt.Errorf("invalid min value: %s", param)
		}
		// Используем руны для корректной работы с кириллицей и иероглифами
		if len([]rune(value)) < minLen {
			return fmt.Errorf("field %s length must be at least %d characters", fieldName, minLen)
		}
	case "max":
		maxLen, err := strconv.Atoi(param)
		if err != nil {
			return fmt.Errorf("invalid max value: %s", param)
		}
		// Используем руны для корректной работы с кириллицей и иероглифами
		if len([]rune(value)) > maxLen {
			return fmt.Errorf("field %s length must not exceed %d characters", fieldName, maxLen)
		}
	case "regexp":
		pattern := param
		matched, err := regexp.MatchString(pattern, value)
		if err != nil {
			return fmt.Errorf("invalid regexp pattern: %s", pattern)
		}
		if !matched {
			return fmt.Errorf("field %s does not match pattern %s", fieldName, pattern)
		}
	}
	return nil
}

func validateInt(fieldName string, value int, rule, param string) error {
	switch rule {
	case "min":
		minVal, err := strconv.Atoi(param)
		if err != nil {
			return fmt.Errorf("invalid min value: %s", param)
		}
		if value < minVal {
			return fmt.Errorf("field %s must be at least %d", fieldName, minVal)
		}
	case "max":
		maxVal, err := strconv.Atoi(param)
		if err != nil {
			return fmt.Errorf("invalid max value: %s", param)
		}
		if value > maxVal {
			return fmt.Errorf("field %s must not exceed %d", fieldName, maxVal)
		}
	}
	return nil
}

func validateFloat(fieldName string, value float64, rule, param string) error {
	switch rule {
	case "min":
		minVal, err := strconv.ParseFloat(param, 64)
		if err != nil {
			return fmt.Errorf("invalid min value: %s", param)
		}
		if value < minVal {
			return fmt.Errorf("field %s must be at least %f", fieldName, minVal)
		}
	case "max":
		maxVal, err := strconv.ParseFloat(param, 64)
		if err != nil {
			return fmt.Errorf("invalid max value: %s", param)
		}
		if value > maxVal {
			return fmt.Errorf("field %s must not exceed %f", fieldName, maxVal)
		}
	}
	return nil
}

type User struct {
	Name  string `validate:"min=3"`
	Age   int    `validate:"min=18;max=65"`
	Email string `validate:"regexp=^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"`
}

func main() {
	// Тест 1: имя слишком короткое
	if err := Validate(User{Name: "Ив", Age: 18, Email: "test@example.com"}); err != nil {
		fmt.Println("Validation error:", err)
	}

	// Тест 2: возраст превышает максимум
	if err := Validate(User{Name: "Иван", Age: 70, Email: "test@example.com"}); err != nil {
		fmt.Println("Validation error:", err)
	}

	// Тест 3: некорректный email
	if err := Validate(User{Name: "Иван", Age: 35, Email: "invalid email"}); err != nil {
		fmt.Println("Validation error:", err)
	}

	// Тест 4: все поля корректны
	if err := Validate(User{Name: "Иван", Age: 35, Email: "test@example.com"}); err != nil {
		fmt.Println("Validation error:", err)
	} else {
		fmt.Println("Validation passed successfully")
	}
}
