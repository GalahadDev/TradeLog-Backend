package trades

import (
	"fmt"
	"strings"
)

func validateScreenshotPaths(paths []string, userID string) error {
	if len(paths) > 3 {
		return fmt.Errorf("máximo 3 screenshots por trade")
	}
	for _, p := range paths {
		if !strings.HasPrefix(p, "https://") && !strings.HasPrefix(p, userID+"/") {
			return fmt.Errorf("screenshot path inválido")
		}
	}
	return nil
}
