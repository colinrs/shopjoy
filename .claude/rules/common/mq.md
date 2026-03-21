# MQ Rules (Message Queue)

Rules for message queue implementation and usage.

## MUST

| # | Rule | Rationale |
|---|------|-----------|
| 1 | Unified MQ technology stack | Avoid fragmentation, simplify operations |
| 2 | Physical isolation between environments | Prevent cross-environment data leakage |
| 3 | Queue/Exchange names must include environment prefix | Clear identification, prevent misrouting |
| 4 | Messages must contain unique `messageId` | Enable tracing, deduplication, debugging |
| 5 | Message timestamps must use UTC | Consistent time handling across regions |
| 6 | Message structure must include version number | Support backward compatibility, schema evolution |
| 7 | Consumers must support idempotent processing | Handle redelivery safely |
| 8 | ACK only after business logic completes | Prevent message loss on processing failure |
| 9 | Critical queues must configure dead-letter queue | Capture failed messages for analysis |

## SHOULD

| # | Rule | Rationale |
|---|------|-----------|
| 10 | Enable producer Confirm mode | Guarantee message delivery to broker |
| 11 | Set retry limits for failed messages | Prevent infinite retry loops |
| 12 | Control consumer concurrency appropriately | Balance throughput and resource usage |
| 13 | Process dead-letter queue regularly | Prevent queue buildup, identify issues |

## FORBIDDEN

| # | Rule | Consequence |
|---|------|-------------|
| 14 | Auto ACK | Message loss on consumer crash |
| 15 | Messages without `messageId` | Cannot trace, deduplicate, or debug |
| 16 | Swallowing exceptions silently | Hidden failures, data inconsistency |
| 17 | Single queue for multiple core businesses | Coupling, blocking, scaling issues |
| 18 | Sharing MQ across environments | Data leakage, environment pollution |

## Code Examples

### Message Structure

```go
// GOOD: Complete message structure
type Message struct {
    MessageID   string    `json:"message_id"`   // Unique identifier
    Version     string    `json:"version"`      // Schema version
    Timestamp   time.Time `json:"timestamp"`    // UTC time
    Type        string    `json:"type"`         // Message type
    Payload     any       `json:"payload"`      // Business data
}

func NewMessage(msgType string, payload any) *Message {
    return &Message{
        MessageID: uuid.New().String(),
        Version:   "1.0",
        Timestamp: time.Now().UTC(),
        Type:      msgType,
        Payload:   payload,
    }
}
```

### Idempotent Consumer

```go
// GOOD: Idempotent processing with manual ACK
func (c *Consumer) HandleMessage(ctx context.Context, msg *Message) error {
    // Check if already processed
    if c.repo.IsProcessed(ctx, msg.MessageID) {
        return nil // Idempotent: skip duplicate
    }

    // Process business logic
    if err := c.processOrder(ctx, msg.Payload); err != nil {
        return err // Don't ACK, will be redelivered
    }

    // Mark as processed
    if err := c.repo.MarkProcessed(ctx, msg.MessageID); err != nil {
        return err
    }

    return nil // ACK after success
}
```

### Dead Letter Queue Configuration

```go
// GOOD: Configure DLQ for critical queues
func setupQueue(ch *amqp.Channel, env string) error {
    queueName := fmt.Sprintf("%s.orders", env)
    dlqName := fmt.Sprintf("%s.orders.dlq", env)

    // Declare DLQ first
    _, err := ch.QueueDeclare(dlqName, true, false, false, false, nil)
    if err != nil {
        return err
    }

    // Declare main queue with DLQ routing
    _, err = ch.QueueDeclare(queueName, true, false, false, false, amqp.Table{
        "x-dead-letter-exchange":    "",
        "x-dead-letter-routing-key": dlqName,
    })
    return err
}
```

## Checklist

- [ ] Message has unique `messageId`
- [ ] Timestamp is UTC
- [ ] Message includes version number
- [ ] Consumer is idempotent
- [ ] Manual ACK after business logic
- [ ] Dead-letter queue configured for critical queues
- [ ] Queue name includes environment prefix
- [ ] Retry limit configured
