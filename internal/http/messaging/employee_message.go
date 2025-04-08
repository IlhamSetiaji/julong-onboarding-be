package messaging

import (
	"errors"
	"log"
	"time"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-onboarding-be/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IEmployeeMessage interface {
	SendFindEmployeeByIDMessage(req request.SendFindEmployeeByIDMessageRequest) (*response.EmployeeResponse, error)
}

type EmployeeMessage struct {
	Log *logrus.Logger
}

func NewEmployeeMessage(log *logrus.Logger) IEmployeeMessage {
	return &EmployeeMessage{
		Log: log,
	}
}

func (m *EmployeeMessage) SendFindEmployeeByIDMessage(req request.SendFindEmployeeByIDMessageRequest) (*response.EmployeeResponse, error) {
	payload := map[string]interface{}{
		"employee_id": req.ID,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "find_employee_by_id",
		MessageData: payload,
		ReplyTo:     "julong_onboarding",
	}

	log.Printf("INFO: document message: %v", docMsg)

	// create channel and add to rchans with uid
	rchan := make(chan response.RabbitMQResponse)
	utils.Rchans[docMsg.ID] = rchan

	// publish rabbit message
	msg := utils.RabbitMsgPublisher{
		QueueName: "julong_sso",
		Message:   *docMsg,
	}
	utils.Pchan <- msg

	// wait for reply
	resp, err := waitReply(docMsg.ID, rchan)
	if err != nil {
		return nil, err
	}

	log.Printf("INFO: response: %v", resp)

	if errMsg, ok := resp.MessageData["error"]; ok {
		return nil, errors.New("[EmployeeMessage.SendFindEmployeeByIDMessage] " + errMsg.(string))
	}

	employeeData := resp.MessageData["employee"].(map[string]interface{})
	employee := convertInterfaceToEmployeeResponse(employeeData)

	return employee, nil
}

func convertInterfaceToEmployeeResponse(data map[string]interface{}) *response.EmployeeResponse {
	// Extract values from the map
	id := data["id"].(string)
	organizationID := data["organization_id"].(string)
	name := data["name"].(string)
	endDate, _ := time.Parse("2006-01-02", data["end_date"].(string))
	retirementDate, _ := time.Parse("2006-01-02", data["retirement_date"].(string))
	email := data["email"].(string)
	mobilePhone := data["mobile_phone"].(string)
	employeeJob := data["employee_job"].(map[string]interface{})
	midsuitID := data["midsuit_id"].(string)

	return &response.EmployeeResponse{
		ID:             uuid.MustParse(id),
		OrganizationID: uuid.MustParse(organizationID),
		Name:           name,
		EndDate:        endDate,
		RetirementDate: retirementDate,
		Email:          email,
		MobilePhone:    mobilePhone,
		EmployeeJob:    employeeJob,
		MidsuitID:      midsuitID,
	}
}

func EmployeeMessageFactory(log *logrus.Logger) IEmployeeMessage {
	return NewEmployeeMessage(log)
}
