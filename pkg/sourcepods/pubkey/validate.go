package pubkey

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

type (
	//ValidationErrors are returned with a slice of all invalid fields
	ValidationErrors struct {
		Errors []ValidationError
	}
	//ValidationError knows for a given field the error
	ValidationError struct {
		Field string
		Error error
	}
)

func (e ValidationErrors) Error() string {
	return fmt.Sprintf("there are %d validation errors", len(e.Errors))
}

func validateCreate(pk *PubKey) error {
	var errs ValidationErrors

	if err := validateName(pk.Name); err != nil {
		errs.Errors = append(errs.Errors, ValidationError{
			Field: "name",
			Error: err,
		})
	}
	if err := validateContent([]byte(pk.Content)); err != nil {
		errs.Errors = append(errs.Errors, ValidationError{
			Field: "content",
			Error: err,
		})
	}

	if len(errs.Errors) > 0 {
		return errs
	}

	return nil
}

func validateName(name string) error {
	if ok := govalidator.IsAlphanumeric(name); !ok {
		return fmt.Errorf("name is not alphanumeric")
	}
	if ok := govalidator.IsByteLength(name, 4, 32); !ok {
		return fmt.Errorf("name is not between 4 and 32 characters long")
	}
	return nil
}

func validateContent(content []byte) error {
	_, _, _, _, err := ssh.ParseAuthorizedKey(content)
	if err != nil {
		return errors.Wrap(err, "public key not valid")
	}
	return nil
}
