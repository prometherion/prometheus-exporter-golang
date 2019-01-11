package queue

import (
	"fmt"
	"time"

	"github.com/adjust/rmq"
	"github.com/prometherion/prometheus-exporter-golang/app/signature"
	"github.com/sirupsen/logrus"
	"gopkg.in/redis.v3"
)

type Manager struct {
	connection rmq.Connection
	taskQueue  rmq.Queue
	queueName  string
}

func New(connectionString, tag string) (manager *Manager, err error) {
	// Testing connection
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    connectionString,
		DB:      0,
	})
	if _, err := c.Ping().Result(); err != nil {
		return nil, fmt.Errorf("cannot establish Redis connection (%s)", err.Error())
	}

	conn := rmq.OpenConnectionWithRedisClient(tag, c)

	return &Manager{
		connection: conn,
	}, nil
}

func (m Manager) Push(signature signature.TaskSignature) (err error) {
	var b []byte
	if b, err = signature.Bytes(); err != nil {
		return err
	}
	if ok := m.taskQueue.PublishBytes(b); ok == false {
		err = fmt.Errorf("cannot push task %s", signature.Id())
	}
	return err
}

func (m Manager) StartConsuming(prefetchLimit int, pollDuration time.Duration) (err error) {
	logrus.Infof("Starting consuming %d, polling every %.f second(s)", prefetchLimit, pollDuration.Seconds())
	if ok := m.taskQueue.StartConsuming(prefetchLimit, pollDuration); ok == false {
		err = fmt.Errorf("cannot start consuming due to error")
	}

	m.taskQueue.AddConsumer("tag", NewConsumer(m.queueName))
	return err
}

func (m *Manager) Open(queueName string) {
	m.queueName = queueName
	m.taskQueue = m.connection.OpenQueue(queueName)
}

