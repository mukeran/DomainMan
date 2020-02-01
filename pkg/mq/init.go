package mq

const (
	MessageQueueCapacity     = 50
	MessageQueueWorkerNumber = 8
)

var (
	MQ *MessageQueue
)

func Init() {
	MQ = NewMQ(MessageQueueCapacity, MessageQueueWorkerNumber)
}
