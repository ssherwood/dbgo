package randoms

import "github.com/google/uuid"

func RandomUUID() string {
	return uuid.NewString()
}
