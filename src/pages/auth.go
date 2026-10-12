package pages

import (
	"errors"
	"net/http"

	"baby-blog/database/models"
	"baby-blog/forms"

	"github.com/AlexDev404/folio"
)

func (p *Pages) LoginPage(c *folio.Ctx) error {
	if userID(c) != "" {
		c.Redirect(folio.SafeNext(c.Request, "/"), http.StatusSeeOther)
		return nil
	}
	c.Set("NextURL", folio.SafeNext(c.Request, ""))
	return nil
}

func (p *Pages) Login(c *folio.Ctx) error {
	form := forms.NewLogin(c.Request)
	user, err := p.Models.Users.Authenticate(form.Email, form.Password)
	if errors.Is(err, models.ErrInvalidCredentials) {
		c.Status(http.StatusUnprocessableEntity)
		c.Fail("Invalid email or password", nil)
		c.Set("Email", form.Email)
		return nil
	}
	if err != nil {
		return err
	}
	logIn(c, user.ID)
	return nil
}

func (p *Pages) RegisterPage(c *folio.Ctx) error {
	if userID(c) != "" {
		c.Redirect("/", http.StatusSeeOther)
		return nil
	}
	c.Set("NextURL", folio.SafeNext(c.Request, ""))
	return nil
}

func (p *Pages) Register(c *folio.Ctx) error {
	form := forms.NewRegister(c.Request)
	errs := form.Validate()
	if errs == nil {
		existing, err := p.Models.Users.GetByEmail(form.Email)
		if err != nil {
			return err
		}
		if existing != nil {
			errs = map[string]string{"email": "Email is already in use"}
		}
	}

	user := &models.User{Email: form.Email}
	if errs == nil {
		err := p.Models.Users.Insert(user, form.Password)
		if errors.Is(err, models.ErrDuplicateEmail) {
			errs = map[string]string{"email": "Email is already in use"}
		} else if err != nil {
			return err
		}
	}

	if errs != nil {
		c.Status(http.StatusUnprocessableEntity)
		c.Fail("", errs)
		c.Set("Email", form.Email)
		return nil
	}
	logIn(c, user.ID)
	return nil
}

// logIn starts a fresh session for userID and sends the visitor on to the
// page they originally asked for.
func logIn(c *folio.Ctx, userID string) {
	s := c.Session()
	s.Renew()
	s.Set(userKey, userID)
	c.Redirect(folio.SafeNext(c.Request, "/"), http.StatusSeeOther)
}

// Logout ends the session on the server and returns to the login page.
func Logout(w http.ResponseWriter, r *http.Request) {
	folio.GetSession(r).Clear()
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
