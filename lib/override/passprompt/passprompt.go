package passprompt

import (
	"fmt"
	"io"
	"os"

	"github.com/howeyc/gopass"
	"github.com/zalando/go-keyring"
)

// Interface for some means of getting a password from the user (or another source)
type PasswordGetter interface {
	// Print the given prompt and retrieve a password. May return io.EOF if the
	// user cancels the prompt.
	GetPasswd(prompt string) (string, error)
}

// A default password getter using stdin/stderr
type PasswordPrompt struct{}

func (PasswordPrompt) GetPasswd(prompt string) (string, error) {
	passwd, err := gopass.GetPasswdPrompt(prompt, true, os.Stdin, os.Stderr)
	if err == io.EOF {
		return "", nil
	} else if err != nil {
		return "", err
	} else {
		return string(passwd), nil
	}
}

type LoginFunc func(string) (bool, error)

func Login(login LoginFunc, getter PasswordGetter, keyringService, keyringUser, initialPrompt, failPrefix string) error {
	keyringFirst := keyringService != ""
	prompt := initialPrompt
	for {
		var password string
		var err error
		if keyringFirst {
			keyringFirst = false
			password, err = keyringGet(keyringService, keyringUser)
			if err == errNotFound {
				continue
			} else if err != nil {
				return fmt.Errorf("keyring error: %w", err)
			}
		} else if getter != nil {
			password, err = getter.GetPasswd(prompt)
			if err != nil {
				return err
			} else if password == "" {
				return io.EOF
			}
		} else {
			return io.EOF
		}
		ok, err := login(password)
		if err != nil {
			return err
		} else if ok {
			if keyringService != "" {
				err := keyringSet(keyringService, keyringUser, password)
				if err != nil {
					return fmt.Errorf("keyring error: %w", err)
				}
			}
			return nil
		}
		prompt = failPrefix + initialPrompt
	}
}

var errNotFound = keyring.ErrNotFound

func keyringGet(service, user string) (string, error) {
	if keyring.Get == nil {
		return "", errNotFound
	}
	return keyring.Get(service, user)
}

func keyringSet(service, user, password string) error {
	if keyring.Set == nil {
		return errNotFound
	}
	return keyring.Set(service, user, password)
}

// func keyringGet(service, user string) (string, error) {
// 	return keyring.Get(service, user)
// }

// func keyringSet(service, user, password string) error {
// 	return keyring.Set(service, user, password)
// }

// // var errNotFound = errors.New("keyring support not available")

// // func keyringGet(service, user string) (string, error) {
// // 	return "", errNotFound
// // }

// // func keyringSet(service, user, password string) error {
// // 	return errNotFound
// // }
