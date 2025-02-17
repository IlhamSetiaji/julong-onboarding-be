package request

type CheckJobExistMessageRequest struct {
	ID string `json:"id"`
}

type SendFindJobByIDMessageRequest struct {
	ID string `json:"id"`
}

type SendFindJobLevelByIDMessageRequest struct {
	ID string `json:"id"`
}

type CheckJobByJobLevelRequest struct {
	JobID      string `json:"job_id"`
	JobLevelID string `json:"job_level_id"`
}
