package asyncq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/tenant"
	"github.com/hibiken/asynq"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
	QueueLow      = "low"
)

type TaskType string

const (
	TaskTypeOrderTimeout   TaskType = "order:timeout"
	TaskTypeStockRelease   TaskType = "stock:release"
	TaskTypePaymentCancel  TaskType = "payment:cancel"
	TaskTypePromotionStart TaskType = "promotion:start"
	TaskTypePromotionEnd   TaskType = "promotion:end"
	TaskTypeCouponExpire   TaskType = "coupon:expire"
	TaskTypeSendEmail      TaskType = "notification:email"
	TaskTypeSendSMS        TaskType = "notification:sms"
	TaskTypeSyncInventory  TaskType = "inventory:sync"
)

type Client struct {
	client *asynq.Client
}

func NewClient(redisAddr string) *Client {
	return &Client{
		client: asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr}),
	}
}

func (c *Client) Close() error {
	return c.client.Close()
}

type TaskPayload struct {
	TenantID shared.TenantID `json:"tenant_id"`
	Type     TaskType        `json:"type"`
	Data     json.RawMessage `json:"data"`
}

func (c *Client) Enqueue(ctx context.Context, taskType TaskType, tenantID shared.TenantID, payload interface{}, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal task payload: %w", err)
	}

	taskPayload := TaskPayload{
		TenantID: tenantID,
		Type:     taskType,
		Data:     data,
	}

	payloadBytes, err := json.Marshal(taskPayload)
	if err != nil {
		return nil, fmt.Errorf("marshal task: %w", err)
	}

	task := asynq.NewTask(string(taskType), payloadBytes)
	return c.client.Enqueue(task, opts...)
}

func (c *Client) EnqueueWithDelay(ctx context.Context, taskType TaskType, tenantID shared.TenantID, payload interface{}, delay time.Duration) (*asynq.TaskInfo, error) {
	return c.Enqueue(ctx, taskType, tenantID, payload, asynq.ProcessIn(delay))
}

func (c *Client) EnqueueWithDeadline(ctx context.Context, taskType TaskType, tenantID shared.TenantID, payload interface{}, deadline time.Time) (*asynq.TaskInfo, error) {
	return c.Enqueue(ctx, taskType, tenantID, payload, asynq.Deadline(deadline))
}

type Server struct {
	server *asynq.Server
	mux    *asynq.ServeMux
}

func NewServer(redisAddr string, concurrency int) *Server {
	if concurrency <= 0 {
		concurrency = 10
	}

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: concurrency,
			Queues: map[string]int{
				QueueCritical: 6,
				QueueDefault:  3,
				QueueLow:      1,
			},
		},
	)

	return &Server{
		server: srv,
		mux:    asynq.NewServeMux(),
	}
}

type Handler func(ctx context.Context, payload json.RawMessage) error

func (s *Server) RegisterHandler(taskType TaskType, handler Handler) {
	s.mux.HandleFunc(string(taskType), func(ctx context.Context, t *asynq.Task) error {
		var payload TaskPayload
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return fmt.Errorf("unmarshal task payload: %w", err)
		}

		ctx = tenant.WithContext(ctx, payload.TenantID)
		return handler(ctx, payload.Data)
	})
}

func (s *Server) Start() error {
	return s.server.Start(s.mux)
}

func (s *Server) Stop() {
	s.server.Stop()
}

func (s *Server) Shutdown() {
	s.server.Shutdown()
}

func QueueName(priority int) string {
	switch {
	case priority >= 8:
		return QueueCritical
	case priority >= 4:
		return QueueDefault
	default:
		return QueueLow
	}
}
