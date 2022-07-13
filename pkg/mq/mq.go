package mq

import (
	"DomainMan/pkg/errors"
	"context"
	"log"
)

type MessageQueue struct {
	capacity     uint
	workerNumber uint
	ctx          context.Context
	cancel       context.CancelFunc
	channel      chan func()
}

func worker(mq *MessageQueue, id uint) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Worker %v experiencing panic. Error: %v\n", id, err)
		}
	}()
	for {
		select {
		case f := <-mq.channel:
			log.Printf("Worker %v received a job\n", id)
			f()
		case <-mq.ctx.Done():
			log.Printf("Stopped Worker %v\n", id)
			return
		}
	}
}

func (mq *MessageQueue) startWorker() {
	for i := uint(0); i < mq.workerNumber; i++ {
		go func(i uint) {
			log.Printf("Worker %v started\n", i)
			worker(mq, i)
		}(i)
	}
}

func (mq *MessageQueue) Push(f func()) error {
	if cap(mq.channel)-len(mq.channel) <= 0 {
		return errors.ErrMessageQueueFull
	}
	mq.channel <- f
	return nil
}

func (mq *MessageQueue) Stop() error {
	if mq.cancel == nil {
		return errors.ErrMessageQueueCancelled
	}
	mq.cancel()
	mq.cancel = nil
	return nil
}

func NewMQ(capacity uint, workerNumber uint) *MessageQueue {
	mq := &MessageQueue{}
	mq.capacity = capacity
	mq.workerNumber = workerNumber
	mq.ctx, mq.cancel = context.WithCancel(context.Background())
	mq.channel = make(chan func(), capacity)
	mq.startWorker()
	return mq
}
