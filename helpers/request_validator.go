package helpers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/thedevsaddam/govalidator"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"time"
)

func init() {
	govalidator.AddCustomRule("array", func(field string, rule string, message string, value interface{}) error {
		val := reflect.ValueOf(value)

		if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
			if message != "" {
				return fmt.Errorf(message)
			}
			return fmt.Errorf("the %s field must be an array", field)
		}
		return nil
	})

	govalidator.AddCustomRule("base64", func(field string, rule string, message string, value interface{}) error {
		val, ok := value.(string)
		if !ok {
			if message != "" {
				return fmt.Errorf(message)
			}
			return fmt.Errorf("the %s field must be a valid base64 encoded string", field)
		}

		_, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			if message != "" {
				return fmt.Errorf(message)
			}
			return fmt.Errorf("the %s field must be a valid base64 encoded string", field)
		}

		return nil
	})

	govalidator.AddCustomRule("uuid_array", func(field, rule, message string, value interface{}) error {
		uuidRegex := regexp.MustCompile(govalidator.UUID4)

		if uuids, ok := value.([]string); ok {
			for _, uuid := range uuids {
				if !uuidRegex.MatchString(uuid) {
					return fmt.Errorf("each item in %s must be a valid UUID", field)
				}
			}
		}
		return nil
	})

	govalidator.AddCustomRule("array_unique_uuid", func(field, rule, message string, value interface{}) error {
		if uuids, ok := value.([]string); ok {
			seenUUIDs := make(map[string]bool)
			for _, uuid := range uuids {
				if seenUUIDs[uuid] {
					return fmt.Errorf("each uuid in %s must be unique", field)
				}
				seenUUIDs[uuid] = true
			}
		}
		return nil
	})

	govalidator.AddCustomRule("valid_future_timestamp", func(field string, rule string, message string, value interface{}) error {
		strVal, ok := value.(string)

		if !ok {
			return fmt.Errorf("%s must be a string", field)
		}

		parsedTime, err := time.Parse(time.DateTime, strVal)

		if err != nil {
			return fmt.Errorf("%s must be a valid timestamp in the format YYYY-MM-DD HH:MM:SS", field)
		}

		if parsedTime.Before(time.Now()) {
			return fmt.Errorf("%s must be a future timestamp", field)
		}

		return nil
	})
}

func ValidateRequest(options govalidator.Options, method string) url.Values {
	var e url.Values

	v := govalidator.New(options)

	switch method {
	case "json":
		e = v.ValidateJSON()
		break
	case "struct":
		e = v.ValidateStruct()
		break
	case "query":
		e = v.Validate()
	}

	return e
}

func ReturnValidationErrors(w http.ResponseWriter, e url.Values) {
	err := map[string]interface{}{"message": "The given data was invalid", "errors": e}
	w.WriteHeader(http.StatusUnprocessableEntity)
	_ = json.NewEncoder(w).Encode(err)
}
