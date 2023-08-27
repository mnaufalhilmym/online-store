package authentity

import (
	"github.com/google/uuid"
	a "hilmy.dev/store/src/modules/account/account_entity"
)

type JWTPayload struct {
	ID         *uuid.UUID `json:"id,omitempty"`
	Role       *a.Role    `json:"role,omitempty"`
	Expiration *int64     `json:"exp,omitempty"`
}
