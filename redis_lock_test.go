package redislock

import (
	"github.com/go-redis/redis/v8"
	"log"
	"sync"
	"testing"
	"time"
)

/**
redis 6 is ok
*/
func NewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB

	})

	return rdb
}

func goTryLock(num int, r *redisLock, w *sync.WaitGroup) {

	defer w.Done()

	log.Println("goTryLock num=", num, " start to lock")
	r.Lock()
	log.Println("goTryLock num=", num, " get lock ok ")
	defer func() {
		log.Println("goTryLock num=", num, " ready to unlock ")
		r.Unlock()
	}()
	time.Sleep(time.Second * 2)
	log.Println("goTryLock num=", num, " sleep done ")
}

func TestLock(t *testing.T) {

	w := sync.WaitGroup{}
	r := NewRedisClient()
	rr, err := NewLock(r, "test-redis-lock")
	if err != nil {
		t.Fatal(err)
	}

	for i := 1; i < 10; i++ {
		w.Add(1)
		go goTryLock(i, rr, &w)
	}

	w.Wait()
	log.Println("all is done")
}
