package redislock

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strings"
	"sync"
)

const StatusOk = 1
const StatusClose = 0
const unlock = "unlock"

type redisLock struct {
	status  int8
	cond    *sync.Cond
	c       *redis.Client
	closing chan struct{}
	key     string
}

func NewLock(rdb *redis.Client, key string) (r *redisLock, err error) {

	if rdb == nil {
		return nil, errors.New("redis.Client is nil")
	}

	key = strings.TrimSpace(key)
	if key == "" {
		return nil, errors.New("key is empty")
	}

	r = &redisLock{
		status:  StatusOk,
		c:       rdb,
		cond:    sync.NewCond(&sync.Mutex{}),
		key:     key,
		closing: make(chan struct{}, 1),
	}

	go r.loop()
	return r, nil
}

func (m *redisLock) loop() {
	ctx := context.Background()
	key := fmt.Sprintf("%v:sub", m.key)
	pubsub := m.c.Subscribe(ctx, key)

	defer pubsub.Close()

	ch := pubsub.Channel()
	for {
		select {
		case msg, ok := <-ch:

			if ok && msg.Payload == unlock {

				m.cond.Broadcast()
			}

		case _ = <-m.closing:
			fmt.Println("all things is closing")
			m.cond.Broadcast()
			break
		}
	}

}

func (m *redisLock) Lock() {
	if m.status == StatusClose {
		panic("lock status is wrong")
	}

	key := fmt.Sprintf("%v:lock", m.key)
	ctx := context.Background()

	for {
		if m.status == StatusClose {
			break
		}
		num, _ := m.c.Incr(ctx, key).Result()

		if num == 1 {
			//lock ok
			break
		}

		//wait
		m.cond.L.Lock()
		m.cond.Wait()
		m.cond.L.Unlock()
	}

}

func (m *redisLock) Unlock() {
	if m.status == StatusClose {
		panic("lock status is wrong")
	}

	key := fmt.Sprintf("%v:lock", m.key)
	ctx := context.Background()

	m.c.Del(ctx, key)

	subKey := fmt.Sprintf("%v:sub", m.key)
	m.c.Publish(ctx, subKey, unlock)

}

func (m *redisLock) Close() {
	m.status = StatusClose
	m.closing <- struct{}{}
}
