package request

type SendFindEmployeeByIDMessageRequest struct {
	ID string `json:"id" binding:"required"`
}
