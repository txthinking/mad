//go:build !darwin && !windows

package mad

import "errors"

func Install(ca string) error {
	return errors.New(`
We cannot automate the certificate installation for you.
Different Linux distributions have different installation methods, the following is an example for Ubuntu:

https://ubuntu.com/server/docs/security-trust-store

sudo apt-get install -y ca-certificates
sudo cp ~/.nami/bin/ca.pem /usr/local/share/ca-certificates/ca.crt
sudo update-ca-certificates

If you are using a different distribution, please refer to the relevant official documentation.
`)
}
