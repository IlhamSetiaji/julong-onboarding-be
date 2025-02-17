package messaging

import (
	"errors"
	"log"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/repository"
	"github.com/IlhamSetiaji/julong-onboarding-be/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IUserMessage interface {
	SendFindUserByIDMessage(request request.SendFindUserByIDMessageRequest) (*response.SendFindUserByIDResponse, error)
	SendGetUserMe(request request.SendFindUserByIDMessageRequest) (*response.SendGetUserMeResponse, error)
	FindUserProfileByUserIDMessage(userId string) (*response.UserProfileResponse, error)
}

type UserMessage struct {
	Log                   *logrus.Logger
	UserProfileRepository repository.IUserProfileRepository
}

func NewUserMessage(log *logrus.Logger, userProfileRepository repository.IUserProfileRepository) IUserMessage {
	return &UserMessage{
		Log:                   log,
		UserProfileRepository: userProfileRepository,
	}
}

func UserMessageFactory(log *logrus.Logger) IUserMessage {
	userProfileRepository := repository.UserProfileRepositoryFactory(log)
	return NewUserMessage(log, userProfileRepository)
}

func (m *UserMessage) SendFindUserByIDMessage(req request.SendFindUserByIDMessageRequest) (*response.SendFindUserByIDResponse, error) {
	payload := map[string]interface{}{
		"user_id": req.ID,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "find_user_by_id",
		MessageData: payload,
		ReplyTo:     "julong_manpower",
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
		return nil, errors.New("[SendFindUserByIDMessage] " + errMsg)
	}

	return &response.SendFindUserByIDResponse{
		ID:   resp.MessageData["user_id"].(string),
		Name: resp.MessageData["name"].(string),
	}, nil
}

func (m *UserMessage) SendGetUserMe(req request.SendFindUserByIDMessageRequest) (*response.SendGetUserMeResponse, error) {
	payload := map[string]interface{}{
		"user_id": req.ID,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "get_user_me",
		MessageData: payload,
		ReplyTo:     "julong_recruitment",
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
		return nil, errors.New("[SendGetUserMe] " + errMsg)
	}

	return &response.SendGetUserMeResponse{
		User: resp.MessageData,
	}, nil
}

func (m *UserMessage) FindUserProfileByUserIDMessage(userId string) (*response.UserProfileResponse, error) {
	parsedUserID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}
	userProfile, err := m.UserProfileRepository.FindByUserID(parsedUserID)
	if err != nil {
		return nil, err
	}

	if userProfile == nil {
		return nil, errors.New("user profile not found")
	}

	m.Log.Printf("INFO: user profile: %v", userProfile)

	return &response.UserProfileResponse{
		ID:     userProfile.ID,
		UserID: userProfile.UserID,
		Name:   userProfile.Name,
		Status: userProfile.Status,
	}, nil
}
