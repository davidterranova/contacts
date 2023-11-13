package eventsourcing

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"
)

type EventStream[T Aggregate] interface {
	Publisher[T]
	Subscriber[T]
}

type Publisher[T Aggregate] interface {
	Publish(ctx context.Context, events ...Event[T]) error
}

type SubscribeFn[T Aggregate] func(e Event[T])

type Subscriber[T Aggregate] interface {
	Subscribe(ctx context.Context, sub SubscribeFn[T])
}

type eventStream[T Aggregate] struct {
	stream      chan Event[T]
	subscribers []SubscribeFn[T]
	mtx         sync.RWMutex
	ctx         context.Context
}

func NewInMemoryPublisher[T Aggregate](ctx context.Context, buffer int) *eventStream[T] {
	p := &eventStream[T]{
		ctx:         ctx,
		stream:      make(chan Event[T], buffer),
		subscribers: make([]SubscribeFn[T], 0),
	}
	p.Run()

	return p
}

func (p *eventStream[T]) Publish(ctx context.Context, events ...Event[T]) error {
	for _, event := range events {
		log.Ctx(ctx).Debug().Str("type", event.EventType()).Interface("event", event).Msg("publishing event")
		p.stream <- event
	}

	return nil
}

func (p *eventStream[T]) Run() {
	go func() {
		for {
			select {
			case event := <-p.stream:
				p.mtx.RLock()
				for _, sub := range p.subscribers {
					sub(event)
				}
				p.mtx.RUnlock()
			case <-p.ctx.Done():
				close(p.stream)
				return
			}
		}
	}()
}

func (p *eventStream[T]) Subscribe(ctx context.Context, sub SubscribeFn[T]) {
	p.mtx.Lock()
	p.subscribers = append(p.subscribers, sub)
	p.mtx.Unlock()
}
