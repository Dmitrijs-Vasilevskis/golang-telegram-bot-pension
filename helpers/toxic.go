package helpers

func IsToxic(message string) bool {
	balance := 0

	for _, char := range message {

		if char == '(' {
			balance++
		}

		if char == ')' {
			if balance == 0 {

				return true
			}
			balance--
		}
	}

	return false
}
