package thermostat

import "fmt"

func padForLcd(str string) string {
	return fmt.Sprintf("%-16v", str)
}
