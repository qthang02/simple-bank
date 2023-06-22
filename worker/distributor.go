package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface{
	DistributeTaskSendVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opt ...asynq.Option) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(residOpt asynq.RedisClientOpt) TaskDistributor  {
	client := asynq.NewClient(residOpt)
	return &RedisTaskDistributor{client: client}
}