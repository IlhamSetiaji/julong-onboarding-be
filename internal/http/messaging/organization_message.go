package messaging

import (
	"errors"
	"log"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	orgResponse "github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-onboarding-be/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IOrganizationMessage interface {
	SendFindOrganizationByIDMessage(request request.SendFindOrganizationByIDMessageRequest) (*orgResponse.SendFindOrganizationByIDMessageResponse, error)
	SendFindOrganizationLocationByIDMessage(request request.SendFindOrganizationLocationByIDMessageRequest) (*orgResponse.SendFindOrganizationLocationByIDMessageResponse, error)
	SendFindOrganizationStructureByIDMessage(request request.SendFindOrganizationStructureByIDMessageRequest) (*orgResponse.SendFindOrganizationStructureByIDMessageResponse, error)
	SendFindOrganizationLocationsPaginatedMessage(page int, pageSize int, search string, includedIDs []string, isNull bool, orgID string) (*orgResponse.OrganizationLocationPaginatedResponse, error)
	SendFindAllOrganizationMessage(includedIDs []string) (*[]orgResponse.OrganizationResponse, error)
	SendFindAllOrganizationLocationsMessage(includedIDs []string) (*[]orgResponse.OrganizationLocationResponse, error)
	SendFindAllOrgStructureChildrenIDsMessage(parentID string) (*[]string, error)
}

type OrganizationMessage struct {
	Log *logrus.Logger
}

func NewOrganizationMessage(log *logrus.Logger) IOrganizationMessage {
	return &OrganizationMessage{
		Log: log,
	}
}

func (m *OrganizationMessage) SendFindOrganizationByIDMessage(req request.SendFindOrganizationByIDMessageRequest) (*orgResponse.SendFindOrganizationByIDMessageResponse, error) {
	payload := map[string]interface{}{
		"organization_id": req.ID,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "find_organization_by_id",
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
		return nil, errors.New("[SendFindOrganizationByIDMessage] " + errMsg)
	}

	return &orgResponse.SendFindOrganizationByIDMessageResponse{
		OrganizationID:       resp.MessageData["organization_id"].(string),
		Name:                 resp.MessageData["name"].(string),
		OrganizationCategory: resp.MessageData["organization_category"].(string),
		OrganizationType:     resp.MessageData["organization_type"].(string),
		MidsuitID:            resp.MessageData["midsuit_id"].(string),
	}, nil
}

func (m *OrganizationMessage) SendFindAllOrganizationMessage(includedIDs []string) (*[]orgResponse.OrganizationResponse, error) {
	payload := map[string]interface{}{
		"included_ids": includedIDs,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "find_all_organization",
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
		return nil, errors.New("[SendFindAllOrganizationMessage] " + errMsg)
	}

	orgs := make([]orgResponse.OrganizationResponse, 0)
	for _, org := range resp.MessageData["organizations"].([]interface{}) {
		orgMap := org.(map[string]interface{})
		orgs = append(orgs, orgResponse.OrganizationResponse{
			ID:                 uuid.MustParse(orgMap["id"].(string)),
			OrganizationTypeID: uuid.MustParse(orgMap["organization_type_id"].(string)),
			Name:               orgMap["name"].(string),
		})
	}

	return &orgs, nil
}

func (m *OrganizationMessage) SendFindOrganizationLocationByIDMessage(req request.SendFindOrganizationLocationByIDMessageRequest) (*orgResponse.SendFindOrganizationLocationByIDMessageResponse, error) {
	payload := map[string]interface{}{
		"organization_location_id": req.ID,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "find_organization_location_by_id",
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
		return nil, errors.New("[SendFindOrganizationLocationByIDMessage] " + errMsg)
	}

	return &orgResponse.SendFindOrganizationLocationByIDMessageResponse{
		OrganizationLocationID: resp.MessageData["organization_location_id"].(string),
		Name:                   resp.MessageData["name"].(string),
	}, nil
}

func (m *OrganizationMessage) SendFindOrganizationStructureByIDMessage(req request.SendFindOrganizationStructureByIDMessageRequest) (*orgResponse.SendFindOrganizationStructureByIDMessageResponse, error) {
	payload := map[string]interface{}{
		"organization_structure_id": req.ID,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "find_organization_structure_by_id",
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
		return nil, errors.New("[SendFindOrganizationStructureByIDMessage] " + errMsg)
	}

	return &orgResponse.SendFindOrganizationStructureByIDMessageResponse{
		OrganizationStructureID: resp.MessageData["organization_structure_id"].(string),
		Name:                    resp.MessageData["name"].(string),
		MidsuitID:               resp.MessageData["midsuit_id"].(string),
	}, nil
}

func (m *OrganizationMessage) SendFindOrganizationLocationsPaginatedMessage(page int, pageSize int, search string, includedIDs []string, isNull bool, orgID string) (*orgResponse.OrganizationLocationPaginatedResponse, error) {
	m.Log.Infof("Included IDs: %v", includedIDs)
	payload := map[string]interface{}{
		"page":            page,
		"page_size":       pageSize,
		"search":          search,
		"included_ids":    includedIDs,
		"is_null":         isNull,
		"organization_id": orgID,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "find_organization_locations_paginated",
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
		return nil, errors.New("[SendFindOrganizationLocationsPaginatedMessage] " + errMsg)
	}

	orgLocs := make([]orgResponse.OrganizationLocationResponse, 0)
	for _, orgLoc := range resp.MessageData["organization_locations"].([]interface{}) {
		orgLocMap := orgLoc.(map[string]interface{})
		orgLocs = append(orgLocs, orgResponse.OrganizationLocationResponse{
			ID:               uuid.MustParse(orgLocMap["id"].(string)),
			OrganizationID:   uuid.MustParse(orgLocMap["organization_id"].(string)),
			OrganizationName: orgLocMap["organization_name"].(string),
			Name:             orgLocMap["name"].(string),
			CreatedAt:        orgLocMap["created_at"].(string),
			UpdatedAt:        orgLocMap["updated_at"].(string),
		})
	}

	return &orgResponse.OrganizationLocationPaginatedResponse{
		OrganizationLocations: orgLocs,
		Total:                 int64(resp.MessageData["total"].(float64)),
	}, nil
}

func (m *OrganizationMessage) SendFindAllOrganizationLocationsMessage(includedIDs []string) (*[]orgResponse.OrganizationLocationResponse, error) {
	payload := map[string]interface{}{
		"included_ids": includedIDs,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "find_all_organization_locations_by_ids",
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
		return nil, errors.New("[SendFindAllOrganizationLocationsMessage] " + errMsg)
	}

	orgLocs := make([]orgResponse.OrganizationLocationResponse, 0)
	for _, orgLoc := range resp.MessageData["organization_locations"].([]interface{}) {
		orgLocMap := orgLoc.(map[string]interface{})
		orgLocs = append(orgLocs, orgResponse.OrganizationLocationResponse{
			ID:               uuid.MustParse(orgLocMap["id"].(string)),
			OrganizationID:   uuid.MustParse(orgLocMap["organization_id"].(string)),
			OrganizationName: orgLocMap["organization_name"].(string),
			Name:             orgLocMap["name"].(string),
			CreatedAt:        orgLocMap["created_at"].(string),
			UpdatedAt:        orgLocMap["updated_at"].(string),
		})
	}

	return &orgLocs, nil
}

func (m *OrganizationMessage) SendFindAllOrgStructureChildrenIDsMessage(parentID string) (*[]string, error) {
	payload := map[string]interface{}{
		"parent_id": parentID,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "find_all_org_structure_children_ids",
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
		return nil, errors.New("[SendFindAllOrgStructureChildrenIDsMessage] " + errMsg)
	}

	ids := make([]string, 0)
	for _, id := range resp.MessageData["children_ids"].([]interface{}) {
		ids = append(ids, id.(string))
	}

	return &ids, nil
}

func OrganizationMessageFactory(log *logrus.Logger) IOrganizationMessage {
	return NewOrganizationMessage(log)
}
