package response

type SendFindUserByIDResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SendGetUserMeResponse struct {
	User map[string]interface{} `json:"user"`
}
