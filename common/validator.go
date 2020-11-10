package common

import (
    "fmt"
    "os"
    "reflect"
    "sync"

    "github.com/gin-gonic/gin/binding"
    "gopkg.in/go-playground/validator.v9"
)

type DefaultValidator struct {
    once     sync.Once
    validate *validator.Validate
}

var _ binding.StructValidator = &DefaultValidator{}

func (v *DefaultValidator) ValidateStruct(obj interface{}) error {
    if kindOfData(obj) == reflect.Struct {
        v.lazyinit()
        if err := v.validate.Struct(obj); err != nil {
            return err
        }
    }

    return nil
}

func (v *DefaultValidator) Engine() interface{} {
    v.lazyinit()
    return v.validate
}

func (v *DefaultValidator) lazyinit() {
    v.once.Do(func() {
        v.validate = validator.New()
        v.validate.SetTagName("binding")

        // add any custom validations etc. here
        for valName, valFunc := range validatorMapper {
            if err := v.validate.RegisterValidation(valName, valFunc); err != nil {
                fmt.Println(err)
                os.Exit(1)
            }
        }
    })
}

func kindOfData(data interface{}) reflect.Kind {
    value := reflect.ValueOf(data)
    valueType := value.Kind()

    if valueType == reflect.Ptr {
        valueType = value.Elem().Kind()
    }
    return valueType
}

func InitValidator() {
    binding.Validator = new(DefaultValidator)
}

var validatorMapper = map[string]func(fl validator.FieldLevel) bool{
    "check_mobile": CheckMobile,
}

func CheckMobile(fl validator.FieldLevel) bool {
    if code, ok := fl.Field().Interface().(string); ok {
        if len(code) != 11 {
            return false
        }
    }

    return true
}
