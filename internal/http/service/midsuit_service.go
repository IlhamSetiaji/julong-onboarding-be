package service

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/config"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/request"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type IMidsuitService interface {
	AuthOneStep() (*AuthOneStepResponse, error)
	SyncEmployeeTaskMidsuit(payload request.SyncEmployeeTaskMidsuitRequest, jwtToken string) (*string, error)
	SyncEmployeeTaskChecklistMidsuit(payload request.SyncEmployeeTaskChecklistMidsuitRequest, jwtToken string) (*string, error)
	SyncEmployeeTaskAttachmentMidsuit(midsuitID int, payload request.SyncEmployeeTaskAttachmentMidsuitRequest, jwtToken string) (*string, error)
	SyncUpdateEmployeeTaskMidsuit(midsuitID int, payload request.SyncEmployeeTaskMidsuitRequest, jwtToken string) (*string, error)
}

type MidsuitService struct {
	Viper *viper.Viper
	Log   *logrus.Logger
	DB    *gorm.DB
}

func NewMidsuitService(
	viper *viper.Viper,
	log *logrus.Logger,
	db *gorm.DB,
) IMidsuitService {
	return &MidsuitService{
		Viper: viper,
		Log:   log,
		DB:    db,
	}
}

func MidsuitServiceFactory(
	viper *viper.Viper,
	log *logrus.Logger,
) IMidsuitService {
	db := config.NewDatabase()
	return NewMidsuitService(viper, log, db)
}

type AuthOneStepResponse struct {
	UserID       int    `json:"userId"`
	Language     string `json:"language"`
	MenuTreeID   int    `json:"menuTreeId"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type SyncEmployeeTaskMidsuitResponse struct {
	ID int `json:"id"`
}

func (s *MidsuitService) AuthOneStep() (*AuthOneStepResponse, error) {
	payload := map[string]interface{}{
		"userName": s.Viper.GetString("midsuit.username"),
		// "password": s.Viper.GetString("midsuit.username") + "321!",
		"password": "JgiMidsuit123!",
		"parameters": map[string]interface{}{
			"clientId":       s.Viper.GetString("midsuit.client_id"),
			"roleId":         s.Viper.GetString("midsuit.role_id"),
			"organizationId": 0,
		},
	}

	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + s.Viper.GetString("midsuit.auth_endpoint")
	method := "POST"

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.AuthOneStep] Error when marshalling payload: " + err.Error())
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.AuthOneStep] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.AuthOneStep] Error when fetching response: " + err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.AuthOneStep] Error when fetching response: " + string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var authResponse AuthOneStepResponse
	if err := json.Unmarshal(bodyBytes, &authResponse); err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.AuthOneStep] Error when unmarshalling response: " + err.Error())
	}

	return &authResponse, nil
}

func (s *MidsuitService) SyncEmployeeTaskMidsuit(payload request.SyncEmployeeTaskMidsuitRequest, jwtToken string) (*string, error) {
	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + "/models/HC_Task"
	method := "POST"

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeTaskMidsuit] Error when marshalling payload: " + err.Error())
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeTaskMidsuit] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeTaskMidsuit] Error when fetching response: " + err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeTaskMidsuit] Error when fetching response haha: " + string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var syncResponse SyncEmployeeTaskMidsuitResponse
	if err := json.Unmarshal(bodyBytes, &syncResponse); err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeTaskMidsuit] Error when unmarshalling response: " + err.Error())
	}

	idStr := strconv.Itoa(syncResponse.ID)
	return &idStr, nil
}

func (s *MidsuitService) SyncEmployeeTaskChecklistMidsuit(payload request.SyncEmployeeTaskChecklistMidsuitRequest, jwtToken string) (*string, error) {
	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + "/models/HC_TaskChecklist"
	method := "POST"

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeTaskChecklistMidsuit] Error when marshalling payload: " + err.Error())
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeTaskChecklistMidsuit] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeTaskChecklistMidsuit] Error when fetching response: " + err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeTaskChecklistMidsuit] Error when fetching response haha: " + string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var syncResponse SyncEmployeeTaskMidsuitResponse
	if err := json.Unmarshal(bodyBytes, &syncResponse); err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeTaskChecklistMidsuit] Error when unmarshalling response: " + err.Error())
	}

	idStr := strconv.Itoa(syncResponse.ID)
	return &idStr, nil
}

func (s *MidsuitService) SyncEmployeeTaskAttachmentMidsuit(midsuitID int, payload request.SyncEmployeeTaskAttachmentMidsuitRequest, jwtToken string) (*string, error) {
	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + "/models/HC_Task/" + strconv.Itoa(midsuitID) + "/attachments"
	method := "POST"

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeTaskAttachmentMidsuit] Error when marshalling payload: " + err.Error())
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeTaskAttachmentMidsuit] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeTaskAttachmentMidsuit] Error when fetching response: " + err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeTaskAttachmentMidsuit] Error when fetching response haha: " + string(bodyBytes))
	}

	// bodyBytes, _ := io.ReadAll(res.Body)
	// var syncResponse SyncEmployeeTaskMidsuitResponse
	// if err := json.Unmarshal(bodyBytes, &syncResponse); err != nil {
	// 	s.Log.Error(err)
	// 	return nil, errors.New("[MidsuitService.SyncEmployeeTaskAttachmentMidsuit] Error when unmarshalling response: " + err.Error())
	// }

	// idStr := strconv.Itoa(syncResponse.ID)
	// return &idStr, nil
	message := "Success"
	return &message, nil
}

func (s *MidsuitService) SyncUpdateEmployeeTaskMidsuit(midsuitID int, payload request.SyncEmployeeTaskMidsuitRequest, jwtToken string) (*string, error) {
	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + "/models/HC_Task/" + strconv.Itoa(midsuitID)
	method := "PUT"

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeTaskMidsuit] Error when marshalling payload: " + err.Error())
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeTaskMidsuit] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeTaskMidsuit] Error when fetching response: " + err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeTaskMidsuit] Error when fetching response haha: " + string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var syncResponse SyncEmployeeTaskMidsuitResponse
	if err := json.Unmarshal(bodyBytes, &syncResponse); err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeTaskMidsuit] Error when unmarshalling response: " + err.Error())
	}

	idStr := strconv.Itoa(syncResponse.ID)
	return &idStr, nil
}
