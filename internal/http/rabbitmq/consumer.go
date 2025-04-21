package rabbitmq

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/service"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/usecase"
	"github.com/IlhamSetiaji/julong-onboarding-be/utils"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitConsumer(viper *viper.Viper, log *logrus.Logger) {
	// conn
	conn, err := amqp091.Dial(viper.GetString("rabbitmq.url"))
	if err != nil {
		log.Printf("ERROR: fail init consumer: %s", err.Error())
		os.Exit(1)
	}

	log.Printf("INFO: done init consumer conn")

	// create channel
	amqpChannel, err := conn.Channel()
	if err != nil {
		log.Printf("ERROR: fail create channel: %s", err.Error())
		os.Exit(1)
	}

	// create queue
	queue, err := amqpChannel.QueueDeclare(
		viper.GetString("rabbitmq.queue"), // channelname
		true,                              // durable
		false,                             // delete when unused
		false,                             // exclusive
		false,                             // no-wait
		nil,                               // arguments
	)
	if err != nil {
		log.Printf("ERROR: fail create queue: %s", err.Error())
		os.Exit(1)
	}

	// channel
	msgChannel, err := amqpChannel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.Printf("ERROR: fail create channel: %s", err.Error())
		os.Exit(1)
	}

	// consume
	for {
		select {
		case msg := <-msgChannel:
			// unmarshal
			docRply := &response.RabbitMQResponse{}
			docMsg := &request.RabbitMQRequest{}
			err = json.Unmarshal(msg.Body, docRply)
			if err != nil {
				log.Printf("ERROR: fail unmarshl: %s", msg.Body)
				msg.Nack(false, true)
				continue
			}
			log.Printf("INFO: received docRply: %v", docRply)

			err = json.Unmarshal(msg.Body, docMsg)
			if err != nil {
				log.Printf("ERROR: fail unmarshl: %s", msg.Body)
				msg.Nack(false, true)
				continue
			}
			log.Printf("INFO: received docMsg: %v", docMsg)

			// ack for message
			err = msg.Ack(true)
			if err != nil {
				log.Printf("ERROR: fail to ack: %s", err.Error())
			}

			// find waiting channel(with uid) and forward the reply to it
			if rchan, ok := utils.Rchans[docRply.ID]; ok {
				rchan <- *docRply
			}

			handleMsg(docMsg, log, viper)
		}
	}
}

func handleMsg(docMsg *request.RabbitMQRequest, log *logrus.Logger, viper *viper.Viper) {
	// switch case
	var msgData map[string]interface{}

	switch docMsg.MessageType {
	case "reply":
		log.Printf("INFO: received reply message")
		return
	case "send_mail":
		to, ok := docMsg.MessageData["to"].(string)
		if !ok {
			log.Errorf("Invalid request format: missing 'to'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'to'").Error(),
			}
			break
		}
		subject, ok := docMsg.MessageData["subject"].(string)
		if !ok {
			log.Errorf("Invalid request format: missing 'subject'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'subject'").Error(),
			}
			break
		}
		body, ok := docMsg.MessageData["body"].(string)
		if !ok {
			log.Errorf("Invalid request format: missing 'body'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'body'").Error(),
			}
			break
		}
		from, ok := docMsg.MessageData["from"].(string)
		if !ok {
			log.Errorf("Invalid request format: missing 'from'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'from'").Error(),
			}
			break
		}

		mailService := service.MailServiceFactory(log, viper)
		err := mailService.SendMail(service.MailData{
			From:    from,
			To:      []string{to},
			Subject: subject,
			Body:    body,
		})
		if err != nil {
			log.Errorf("ERROR: fail send mail: %s", err.Error())
			msgData = map[string]interface{}{
				"error": err.Error(),
			}
			break
		} else {
			log.Printf("INFO: success send mail")
		}

		msgData = map[string]interface{}{
			"message": "success",
		}
	case "find_user_profile_by_user_id":
		userID, ok := docMsg.MessageData["user_id"].(string)
		if !ok {
			log.Errorf("Invalid request format: missing 'user_id'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'user_id'").Error(),
			}
			break
		}
		userMessageFactory := messaging.UserMessageFactory(log)
		resp, err := userMessageFactory.FindUserProfileByUserIDMessage(userID)
		if err != nil {
			log.Errorf("ERROR: fail find user profile by user id: %s", err.Error())
			msgData = map[string]interface{}{
				"error": err.Error(),
			}
			break
		} else {
			log.Printf("INFO: success find user profile by user id")
		}

		msgData = map[string]interface{}{
			"user_profile": resp,
		}
	case "create_employee_tasks":
		employeeID, ok := docMsg.MessageData["employee_id"].(string)
		if !ok {
			log.Errorf("Invalid request format: missing 'employee_id'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'employee_id'").Error(),
			}
			break
		}
		joinedDate, ok := docMsg.MessageData["joined_date"].(string)
		if !ok {
			log.Errorf("Invalid request format: missing 'joined_date'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'joined_date'").Error(),
			}
			break
		}
		organizationType, ok := docMsg.MessageData["organization_type"].(string)
		if !ok {
			log.Errorf("Invalid request format: missing 'organization_type'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'organization_type'").Error(),
			}
			break
		}
		employeeMidsuitID, ok := docMsg.MessageData["employee_midsuit_id"].(string)
		if !ok {
			log.Errorf("Invalid request format: missing 'employee_midsuit_id'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'employee_midsuit_id'").Error(),
			}
			break
		}

		jobMidsuitID, ok := docMsg.MessageData["job_midsuit_id"].(string)
		if !ok {
			log.Errorf("Invalid request format: missing 'job_midsuit_id'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'job_midsuit_id'").Error(),
			}
			break
		}

		jobLevelMidsuitID, ok := docMsg.MessageData["job_level_midsuit_id"].(string)
		if !ok {
			log.Errorf("Invalid request format: missing 'job_level_midsuit_id'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'job_level_midsuit_id'").Error(),
			}
			break
		}

		orgMidsuitID, ok := docMsg.MessageData["org_midsuit_id"].(string)
		if !ok {
			log.Errorf("Invalid request format: missing 'org_midsuit_id'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'org_midsuit_id'").Error(),
			}
			break
		}

		orgStructureMidsuitID, ok := docMsg.MessageData["org_structure_midsuit_id"].(string)
		if !ok {
			log.Errorf("Invalid request format: missing 'org_structure_midsuit_id'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'org_structure_midsuit_id'").Error(),
			}
			break
		}

		templateTaskUseCaseFactory := usecase.EmployeeTaskUseCaseFactory(log, viper)
		err := templateTaskUseCaseFactory.CreateEmployeeTasksForRecruitment(&request.CreateEmployeeTasksForRecruitment{
			EmployeeID:            employeeID,
			JoinedDate:            joinedDate,
			OrganizationType:      organizationType,
			EmployeeMidsuitID:     employeeMidsuitID,
			JobMidsuitID:          jobMidsuitID,
			JobLevelMidsuitID:     jobLevelMidsuitID,
			OrgMidsuitID:          orgMidsuitID,
			OrgStructureMidsuitID: orgStructureMidsuitID,
		})
		if err != nil {
			log.Errorf("ERROR: fail create employee tasks: %s", err.Error())
			msgData = map[string]interface{}{
				"error": err.Error(),
			}
			break
		}

		msgData = map[string]interface{}{
			"message": "success",
		}
	case "count_kanban_progress_by_employee_id":
		employeeID, ok := docMsg.MessageData["employee_id"].(string)
		if !ok {
			log.Errorf("Invalid request format: missing 'employee_id'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'employee_id'").Error(),
			}
			break
		}
		parsedEmployeeUUID, err := uuid.Parse(employeeID)
		if err != nil {
			log.Errorf("Invalid request format: invalid 'employee_id'")
			msgData = map[string]interface{}{
				"error": errors.New("invalid 'employee_id'").Error(),
			}
			break
		}

		employeeTaskUseCaseFactory := usecase.EmployeeTaskUseCaseFactory(log, viper)
		resp, err := employeeTaskUseCaseFactory.CountKanbanProgressByEmployeeID(parsedEmployeeUUID)
		if err != nil {
			log.Errorf("ERROR: fail count kanban progress by employee id: %s", err.Error())
			msgData = map[string]interface{}{
				"error": err.Error(),
			}
			break
		} else {
			log.Printf("INFO: success count kanban progress by employee id")
		}

		msgData = map[string]interface{}{
			"employee_id": employeeID,
			"total_task":  resp.TotalTask,
			"to_do":       resp.ToDo,
			"in_progress": resp.InProgress,
			"need_review": resp.NeedReview,
			"completed":   resp.Completed,
		}
	default:
		log.Printf("Unknown message type, please recheck your type: %s", docMsg.MessageType)

		msgData = map[string]interface{}{
			"error": errors.New("unknown message type").Error(),
		}
	}
	// reply
	reply := response.RabbitMQResponse{
		ID: docMsg.ID,
		// MessageType: docMsg.MessageType,
		MessageType: "reply",
		MessageData: msgData,
	}
	msg := utils.RabbitMsgConsumer{
		QueueName: docMsg.ReplyTo,
		Reply:     reply,
	}
	utils.Rchan <- msg
}
