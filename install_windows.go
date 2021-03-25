package mad

import "os/exec"

func Install(ca string) error {
	cmd := exec.Command("sh", "-c", "certutil -addstore -f \"ROOT\" "+ca)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
