package auth

import (
	"github.com/msteinert/pam"
)

func Authenticate(username, password string) error {
	t, err := pam.StartFunc("", username, func(s pam.Style, msg string) (string, error) {
		switch s {
		case pam.PromptEchoOff:
			return password, nil
		case pam.PromptEchoOn:
			return password, nil
		case pam.ErrorMsg:
			return "", nil
		case pam.TextInfo:
			return "", nil
		}
		return "", nil
	})
	if err != nil {
		return err
	}
	return t.Authenticate(0)
}
