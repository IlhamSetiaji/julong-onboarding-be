package utils

import (
	"crypto/rand"
	"math/big"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
)

var ResponseChannel = make(chan map[string]interface{}, 100)

var Rchans = make(map[string](chan response.RabbitMQResponse))

type RabbitMsgPublisher struct {
	QueueName string                  `json:"queueName"`
	Message   request.RabbitMQRequest `json:"message"`
}

type RabbitMsgConsumer struct {
	QueueName string                    `json:"queueName"`
	Reply     response.RabbitMQResponse `json:"reply"`
}

// channel to publish rabbit messages
var Pchan = make(chan RabbitMsgPublisher, 10)
var Rchan = make(chan RabbitMsgConsumer, 10)

func GenerateRandomIntToken(digits int) (int64, error) {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(digits)), nil).Sub(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(digits)), nil), big.NewInt(1))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}
	return n.Int64(), nil
}

func GenerateRandomStringToken(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return string(b)
}
