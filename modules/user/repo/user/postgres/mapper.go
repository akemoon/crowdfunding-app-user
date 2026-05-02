package postgres

import "fmt"

func MapRoleToDB(role string) (int, error) {
	switch role {
	case "user":
		return 1, nil
	case "moder":
		return 2, nil
	case "admin":
		return 3, nil
	default:
		return 0, fmt.Errorf("unknown role: %s", role)
	}
}
