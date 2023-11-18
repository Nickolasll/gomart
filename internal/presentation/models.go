package presentation

type RegistrationPayload struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
