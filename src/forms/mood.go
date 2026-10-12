package forms

import (
	"net/http"

	"github.com/AlexDev404/folio/validator"
)

type MoodForm struct {
	Mood string
}

func NewMoodForm(r *http.Request) MoodForm {
	return MoodForm{Mood: r.PostFormValue("mood")}
}

// Validate returns the field errors, or nil if the form is valid.
func (f MoodForm) Validate() map[string]string {
	v := validator.NewValidator()
	v.Check(validator.NotBlank(f.Mood), "mood", "Mood is required")
	v.Check(validator.MinLength(f.Mood, 3), "mood", "Mood must be at least 3 characters")
	v.Check(validator.MaxLength(f.Mood, 100), "mood", "Mood must not exceed 100 characters")
	if v.ValidData() {
		return nil
	}
	return v.Errors
}
