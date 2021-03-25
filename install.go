// +build !darwin
// +build !windows

package mad

import "errors"

func Install(ca string) error {
	return errors.New("Unsupported your OS, PR welcome, https://github.com/txthinking/mad")
}
