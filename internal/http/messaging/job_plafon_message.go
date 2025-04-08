package messaging

import (
	"errors"
	"log"

	// "github.com/IlhamSetiaji/go-rabbitmq-utils/rabbitmq"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	jobResponse "github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-onboarding-be/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IJobPlafonMessage interface {
	SendCheckJobExistMessage(request request.CheckJobExistMessageRequest) (*jobResponse.CheckJobExistMessageResponse, error)
	SendFindJobByIDMessage(request request.SendFindJobByIDMessageRequest) (*jobResponse.SendFindJobByIDMessageResponse, error)
	SendFindJobLevelByIDMessage(request request.SendFindJobLevelByIDMessageRequest) (*jobResponse.SendFindJobLevelByIDMessageResponse, error)
	SendCheckJobByJobLevelMessage(request request.CheckJobByJobLevelRequest) (*jobResponse.CheckJobExistMessageResponse, error)
}

type JobPlafonMessage struct {
	Log *logrus.Logger
}

func NewJobPlafonMessage(log *logrus.Logger) IJobPlafonMessage {
	return &JobPlafonMessage{
		Log: log,
	}
}

func (m *JobPlafonMessage) SendCheckJobExistMessage(req request.CheckJobExistMessageRequest) (*jobResponse.CheckJobExistMessageResponse, error) {
	payload := map[string]interface{}{
		"job_id": req.ID,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "find_job_by_id",
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

	if errMsg, ok := resp.MessageData["error"].(string); ok && errMsg != "" {
		return nil, errors.New("[SendCheckJobExistMessage] " + errMsg)
	}

	exist := resp.MessageData["job_id"].(string) != ""

	return &jobResponse.CheckJobExistMessageResponse{
		JobID: uuid.MustParse(resp.MessageData["job_id"].(string)),
		Exist: exist,
	}, nil
}

func (m *JobPlafonMessage) SendFindJobByIDMessage(req request.SendFindJobByIDMessageRequest) (*jobResponse.SendFindJobByIDMessageResponse, error) {
	payload := map[string]interface{}{
		"job_id": req.ID,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "find_job_by_id",
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

	if errMsg, ok := resp.MessageData["error"].(string); ok && errMsg != "" {
		return nil, errors.New("[SendFindJobByIDMessage] " + errMsg)
	}

	return &jobResponse.SendFindJobByIDMessageResponse{
		JobID:     uuid.MustParse(resp.MessageData["job_id"].(string)),
		Name:      resp.MessageData["name"].(string),
		MidsuitID: resp.MessageData["midsuit_id"].(string),
	}, nil
}

func (m *JobPlafonMessage) SendFindJobLevelByIDMessage(req request.SendFindJobLevelByIDMessageRequest) (*jobResponse.SendFindJobLevelByIDMessageResponse, error) {
	payload := map[string]interface{}{
		"job_level_id": req.ID,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "find_job_level_by_id",
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

	if errMsg, ok := resp.MessageData["error"].(string); ok && errMsg != "" {
		return nil, errors.New("[SendFindJobLevelByIDMessage] " + errMsg)
	}

	return &jobResponse.SendFindJobLevelByIDMessageResponse{
		JobLevelID: uuid.MustParse(resp.MessageData["job_level_id"].(string)),
		Name:       resp.MessageData["name"].(string),
		Level:      resp.MessageData["level"].(float64),
		MidsuitID:  resp.MessageData["midsuit_id"].(string),
	}, nil
}

func (m *JobPlafonMessage) SendCheckJobByJobLevelMessage(req request.CheckJobByJobLevelRequest) (*jobResponse.CheckJobExistMessageResponse, error) {
	payload := map[string]interface{}{
		"job_id":       req.JobID,
		"job_level_id": req.JobLevelID,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "check_job_by_job_level",
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

	if errMsg, ok := resp.MessageData["error"].(string); ok && errMsg != "" {
		return nil, errors.New("[SendCheckJobByJobLevelMessage] " + errMsg)
	}

	exist := resp.MessageData["job_id"].(string) != ""

	return &jobResponse.CheckJobExistMessageResponse{
		JobID: uuid.MustParse(resp.MessageData["job_id"].(string)),
		Exist: exist,
	}, nil
}

func JobPlafonMessageFactory(log *logrus.Logger) IJobPlafonMessage {
	return NewJobPlafonMessage(log)
}
