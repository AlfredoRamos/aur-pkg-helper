package types

import "net/mail"

type GitConfig struct {
	Name  string
	Email string
}

func (c *GitConfig) IsValidName() bool {
	return len(c.Name) > 0
}

func (c *GitConfig) IsValidEmail() bool {
	if len(c.Email) < 1 {
		return false
	}

	_, err := mail.ParseAddress(c.Email)

	return err == nil
}

func (c *GitConfig) IsValid() bool {
	return c.IsValidName() && c.IsValidEmail()
}
