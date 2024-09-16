package config

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ahmdyaasiin/agria-backend/internal/pkg/query"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"log"
	"strings"
)

type CustomValidator struct {
	Validate *validator.Validate
	DB       *sqlx.DB
}

func NewValidator(DB *sqlx.DB) *CustomValidator {
	v := validator.New()
	customValidator := &CustomValidator{
		Validate: v,
		DB:       DB,
	}

	err := v.RegisterValidation("unique", customValidator.Unique)
	if err != nil {
		panic(err)
	}

	return customValidator
}

func (cv *CustomValidator) Unique(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().String()
	fieldName := fl.StructFieldName()

	if fieldName == "PhoneNumber" {
		fieldValue = strings.Replace(fieldValue, "+", "", 1)
		if strings.HasPrefix(fieldValue, "0") {
			fieldValue = "62" + fieldValue[1:]
		}
	}

	var exists bool
	q := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM users WHERE %s = ?)", query.ConvertToSnakeCase(fieldName))

	err := cv.DB.Get(&exists, q, fieldValue)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("Database error (validation): %v\n", err)
		return false
	}

	return !exists
}
