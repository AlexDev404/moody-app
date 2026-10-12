// Package forms reads and validates the app's forms.
package forms

import (
	"net/http"
	"strings"

	"github.com/AlexDev404/folio/validator"
)

// Mood is the check-in form on the home page.
type Mood struct {
	Mood string
}

func NewMood(r *http.Request) Mood {
	return Mood{Mood: strings.TrimSpace(r.PostFormValue("mood"))}
}

// Validate returns field errors, or nil if the form is valid.
func (f Mood) Validate() map[string]string {
	v := validator.NewValidator()
	v.Check(validator.NotBlank(f.Mood), "mood", "Mood is required")
	v.Check(validator.MinLength(f.Mood, 3), "mood", "Mood must be at least 3 characters")
	v.Check(validator.MaxLength(f.Mood, 100), "mood", "Mood must not exceed 100 characters")
	return errorsOrNil(v)
}

// Login is the login form.
type Login struct {
	Email    string
	Password string
}

func NewLogin(r *http.Request) Login {
	return Login{
		Email:    strings.TrimSpace(r.PostFormValue("email")),
		Password: r.PostFormValue("password"),
	}
}

// Register is the sign-up form.
type Register struct {
	Email           string
	Password        string
	PasswordConfirm string
}

func NewRegister(r *http.Request) Register {
	return Register{
		Email:           strings.TrimSpace(r.PostFormValue("email")),
		Password:        r.PostFormValue("password"),
		PasswordConfirm: r.PostFormValue("password_confirm"),
	}
}

// Validate returns field errors, or nil if the form is valid. Whether the
// email is already taken is checked separately against the database.
func (f Register) Validate() map[string]string {
	v := validator.NewValidator()
	v.Check(validator.NotBlank(f.Email), "email", "Email is required")
	v.Check(validator.MaxLength(f.Email, 255), "email", "Email must not exceed 255 characters")
	v.Check(f.Email == "" || validator.IsValidEmail(f.Email), "email", "Enter a valid email address")
	v.Check(f.Password != "", "password", "Password is required")
	v.Check(f.Password == "" || validator.MinLength(f.Password, 8), "password", "Password must be at least 8 characters long")
	// bcrypt only uses the first 72 bytes of a password.
	v.Check(len(f.Password) <= 72, "password", "Password must not exceed 72 bytes")
	v.Check(f.Password == f.PasswordConfirm, "password_confirm", "Passwords do not match")
	return errorsOrNil(v)
}

func errorsOrNil(v *validator.Validator) map[string]string {
	if v.ValidData() {
		return nil
	}
	return v.Errors
}
