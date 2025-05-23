package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates a custom form struct, embeds a url.Values object
type Form struct {
	URLValues url.Values
	Errors    errors
}

// Valid returns true if there are no errors
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		URLValues: data,
		Errors:    errors(map[string][]string{}),
	}
}

// Required checks if form field is in post and not empty
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		x := f.URLValues.Get(field)
		// remove _ from field
		fieldDash := strings.Replace(field, "_", " ", -1)
		if strings.TrimSpace(x) == "" {
			f.Errors.Add(field, fmt.Sprintf("The %s field is required", fieldDash))
		}
	}
}

// Has checks if form field is in post and not empty
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	return x != ""
}

// MinLength checks for string minimum length
func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	x := r.Form.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

// Email checks for valid email address
func (f *Form) Email(field string, r *http.Request) bool {
	x := r.Form.Get(field)

	if !govalidator.IsEmail(x) {
		f.Errors.Add(field, fmt.Sprintf("Invalid email address: %s", x))
		return false
	}
	return true
}
