package response

import "time"

type GetUserResponse struct {
	ID       int64   `json:"id"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	GoogleID *string `json:"googleId"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GetAccountResponse struct {
	User GetUserResponse `json:"user"`
	Plan *PlanResponse   `json:"plan"`
}
