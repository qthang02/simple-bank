package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskSendVerifyEmail = "Task: send_verify_email"

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option) error {
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshl task payload: %w", &err)
		}

		task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
		info, err := distributor.client.EnqueueContext(ctx, task)
		if err != nil {
			return fmt.Errorf("failed to enqueue task: %w", &err)
		}

		log.Info().Str("type", task.Type()).
			Bytes("payload", task.Payload()).
			Str("queue", info.Queue).
			Int("max_retry", info.MaxRetry).
			Msgf("enqueued task")
		
		return nil
}


func (processor *RedisTaskProcessor)	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
		var payload PayloadSendVerifyEmail
		if err := json.Unmarshal(task.Payload(), &payload); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", &err)
		}

		user, err := processor.store.GetUser(ctx, payload.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("user %s not found: %w", payload.Username, &err)
			}
			return fmt.Errorf("failed to get user %s: %w", payload.Username, &err)
		}

		// TODO: send email to user.Email

		log.Info().Str("type", task.Type()).
			Bytes("payload", task.Payload()).
			Str("email", user.Email).
			Msgf("processing task")

		return nil
}