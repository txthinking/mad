package mad

import "os/exec"

func Install(ca string) error {
	cmd := exec.Command("sh", "-c", "sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain "+ca)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
