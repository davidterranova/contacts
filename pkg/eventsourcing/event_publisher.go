package eventsourcing

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

type EventStreamPublisher[T Aggregate] struct {
	eventStore EventStore[T]
	stream     Publisher[T]
	batchSize  int
}

func NewEventStreamPublisher[T Aggregate](eventStore EventStore[T], stream Publisher[T], batchSize int) *EventStreamPublisher[T] {
	return &EventStreamPublisher[T]{
		eventStore: eventStore,
		stream:     stream,
		batchSize:  batchSize,
	}
}

func (p *EventStreamPublisher[T]) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := p.processBatch(ctx)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("event stream publisher: failed to process batch")
			}
		}
		time.Sleep(30 * time.Second)
	}
}

func (p *EventStreamPublisher[T]) processBatch(ctx context.Context) error {
	events, err := p.eventStore.LoadUnpublished(ctx, p.batchSize)
	if err != nil {
		return err
	}

	if len(events) == 0 {
		return nil
	}

	err = p.stream.Publish(ctx, events...)
	if err != nil {
		return err
	}

	err = p.eventStore.MarkPublished(ctx, events...)
	if err != nil {
		return err
	}

	return nil
}
