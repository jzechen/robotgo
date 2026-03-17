//go:build !darwin

package robotgo

import "errors"

func typeByOSAScript(text string) error {
	return errors.New("unsupported platform")
}
