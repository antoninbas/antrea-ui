package v1alpha1

type UpdatePassword struct {
	// base64-encoded password
	Password []byte `json:"password"`
}
