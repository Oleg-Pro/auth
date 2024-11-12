package user_saver

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"

	"github.com/Oleg-Pro/auth/internal/model"
)

func (s *service) UserSaveHandler(_ context.Context, msg *sarama.ConsumerMessage) error {
	userInfo := &model.UserInfo{}
	err := json.Unmarshal(msg.Value, userInfo)
	if err != nil {
		return err
	}

	log.Printf("New User Created: %#v", userInfo)
	return nil
}
