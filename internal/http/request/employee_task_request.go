package request

import "mime/multipart"

type CreateEmployeeTaskRequest struct {
	CoverPath               *string                         `form:"cover_path" validate:"required"`
	EmployeeID              *string                         `form:"employee_id" validate:"required,uuid"`
	TemplateTaskID          *string                         `form:"template_task_id" validate:"omitempty"`
	SurveyTemplateID        *string                         `form:"survey_template_id" validate:"omitempty"`
	Name                    string                          `form:"name" validate:"required"`
	Priority                string                          `form:"priority" validate:"required,employee_task_priority_validation"`
	Description             string                          `form:"description" validate:"omitempty"`
	StartDate               string                          `form:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate                 string                          `form:"end_date" validate:"required,datetime=2006-01-02"`
	EmployeeTaskAttachments []EmployeeTaskAttachmentRequest `form:"employee_task_attachments" validate:"omitempty,dive"`
	EmployeeTaskChecklists  []EmployeeTaskChecklistRequest  `form:"employee_task_checklists" validate:"omitempty,dive"`
}

type UpdateEmployeeTaskRequest struct {
	ID                      *string                         `form:"id" validate:"required"`
	CoverPath               *string                         `form:"cover_path" validate:"required"`
	EmployeeID              *string                         `form:"employee_id" validate:"required,uuid"`
	TemplateTaskID          *string                         `form:"template_task_id" validate:"omitempty"`
	SurveyTemplateID        *string                         `form:"survey_template_id" validate:"omitempty"`
	Name                    string                          `form:"name" validate:"required"`
	Priority                string                          `form:"priority" validate:"required,employee_task_priority_validation"`
	Description             string                          `form:"description" validate:"omitempty"`
	StartDate               string                          `form:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate                 string                          `form:"end_date" validate:"required,datetime=2006-01-02"`
	VerifiedBy              *string                         `form:"verified_by" validate:"omitempty"`
	IsDone                  string                          `form:"is_done" validate:"omitempty,oneof=YES NO"`
	Proof                   *multipart.FileHeader           `form:"proof" validate:"omitempty"`
	ProofPath               *string                         `form:"proof_path" validate:"omitempty"`
	Status                  string                          `form:"status" validate:"omitempty,employee_task_status_validation"`
	Kanban                  string                          `form:"kanban" validate:"omitempty,employee_task_kanban_validation"`
	Notes                   string                          `form:"notes" validate:"omitempty"`
	EmployeeTaskAttachments []EmployeeTaskAttachmentRequest `form:"employee_task_attachments" validate:"omitempty,dive"`
	EmployeeTaskChecklists  []EmployeeTaskChecklistRequest  `form:"employee_task_checklists" validate:"omitempty,dive"`
}

type UpdateEmployeeTaskOnlyRequest struct {
	ID        *string `form:"id" validate:"required"`
	StartDate string  `form:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate   string  `form:"end_date" validate:"required,datetime=2006-01-02"`
}

type CreateEmployeeTasksForRecruitment struct {
	EmployeeID       string `json:"employee_id" validate:"required,uuid"`
	JoinedDate       string `json:"joined_date" validate:"required,datetime=2006-01-02"`
	OrganizationType string `json:"organization_type" validate:"required"`
}

type AdOrgId struct {
	ID         int    `json:"id" binding:"omitempty"`
	Identifier string `json:"identifier" binding:"required"`
}
type HcEmployeeId struct {
	ID         int    `json:"id" binding:"required"`
	Identifier string `json:"identifier" binding:"omitempty"`
}

type HcJobLevelId struct {
	ID         int    `json:"id" binding:"required"`
	Identifier string `json:"identifier" binding:"omitempty"`
}

type HcJobId struct {
	ID         int    `json:"id" binding:"required"`
	Identifier string `json:"identifier" binding:"midsuit"`
}

type HcOrgId struct {
	ID         int    `json:"id" binding:"required"`
	Identifier string `json:"identifier" binding:"omitempty"`
}

type SyncEmployeeTaskMidsuitRequest struct {
	AdOrgId          AdOrgId          `json:"AD_Org_ID" binding:"required"`
	Name             string           `json:"Name" binding:"required"`
	Category         TaskCategory     `json:"Category" binding:"required"`
	StartDate        string           `json:"StartDate" binding:"required"`
	EndDate          string           `json:"EndDate" binding:"required"`
	HCApproverID     HcApproverId     `json:"HC_Approver_ID" binding:"omitempty"`
	HCEmployeeID     HcEmployeeId     `json:"HC_Employee_ID" binding:"required"`
	HCJob2ID         HcJobId          `json:"HC_Job2_ID" binding:"omitempty"`
	HCJobLevel2ID    HcJobLevelId     `json:"HC_JobLevel2_ID" binding:"omitempty"`
	HCJobLevelID     HcJobLevelId     `json:"HC_JobLevel_ID" binding:"omitempty"`
	HCJobID          HcJobId          `json:"HC_Job_ID" binding:"omitempty"`
	HCOrg2ID         HcOrgId          `json:"HC_Org2_ID" binding:"omitempty"`
	HCOrgID          HcOrgId          `json:"HC_Org_ID" binding:"omitempty"`
	HCApproverUserID HcApproverUserId `json:"HC_ApproverUser_ID" binding:"omitempty"`
}

type TaskCategory struct {
	PropertyLabel string `json:"propertyLabel" binding:"omitempty"`
	Identifier    string `json:"identifier" binding:"required"`
	ModelName     string `json:"model-name" binding:"required"`
}

type HcApproverId struct {
	ID int `json:"id" binding:"required"`
}

type HcApproverUserId struct {
	ID int `json:"id" binding:"required"`
}

type SyncEmployeeTaskChecklistMidsuitRequest struct {
	AdOrgId      AdOrgId      `json:"AD_Org_ID" binding:"required"`
	Name         string       `json:"Name" binding:"required"`
	HCTaskID     HCTaskID     `json:"HC_Task_ID" binding:"required"`
	IsChecked    bool         `json:"IsChecked" binding:"required"`
	HCEmployeeID HcEmployeeId `json:"HC_Employee_ID" binding:"required"`
	ModelName    string       `json:"model-name" binding:"required"`
}

type HCTaskID struct {
	ID        int    `json:"id" binding:"required"`
	ModelName string `json:"model-name" binding:"required"`
}

type SyncEmployeeTaskAttachmentMidsuitRequest struct {
	Name string `json:"name" binding:"required"`
	Data string `json:"data" binding:"required"`
}
