package kafka

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"time"
)

const TopicReputationEntry = "user.reputation.entry.v1"

type ReputationEntryEvent struct {
	EventID    string    `json:"event_id"`
	UserID     string    `json:"user_id"`
	Delta      int32     `json:"delta"`
	Reason     string    `json:"reason"`
	Source     string    `json:"source"`
	OccurredAt time.Time `json:"occurred_at"`
	TraceID    string    `json:"trace_id,omitempty"`
}

type Producer struct {
	w *kafka.Writer
}

func NewProducer(broker string) *Producer {
	return &Producer{
		w: &kafka.Writer{
			Addr:         kafka.TCP(broker),
			Topic:        TopicReputationEntry,
			Balancer:     &kafka.Hash{},
			RequiredAcks: kafka.RequireAll,
			Async:        false,
		},
	}
}
func (p *Producer) Close() error { return p.w.Close() }

func (p *Producer) PublishRegistration(ctx context.Context, userID, traceID string) error {
	ev := ReputationEntryEvent{
		EventID:    uuid.NewString(),
		UserID:     userID,
		Delta:      10,
		Reason:     "user_registered",
		Source:     "UserService",
		OccurredAt: time.Now().UTC(),
		TraceID:    traceID,
	}
	data, _ := json.Marshal(ev)

	msg := kafka.Message{
		Key:   []byte(userID),
		Value: data,
		Time:  time.Now(),
	}

	var attempt int
	for {
		if err := p.w.WriteMessages(ctx, msg); err != nil {
			attempt++
			if attempt >= 5 {
				return err
			}
			time.Sleep(time.Duration(attempt*200) * time.Millisecond)
			continue
		}
		return nil
	}
}
