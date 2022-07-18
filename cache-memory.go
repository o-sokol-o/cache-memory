package cachememory

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zhashkevych/scheduler"
)

var ()

type CacheMemory interface {
	Set(key string, value interface{})
	Get(key string) (interface{}, error)
	Delete(key string)
}

type cacheValue struct {
	t    time.Time
	data interface{}
}

type Cache struct {
	worker *scheduler.Scheduler
	ctx    context.Context
	lt     int64 // default life time
	cv     map[string]cacheValue
}

func New(lifeTimeSec int64) *Cache {
	c := Cache{
		ctx:    context.Background(),
		worker: scheduler.NewScheduler(),
		lt:     lifeTimeSec * 1000000000, // default life time nanosecond
		cv:     make(map[string]cacheValue)}

	c.worker.Add(c.ctx,
		func(ctx context.Context) {
			func() {
				for k, v := range c.cv {
					if time.Now().UnixNano()-v.t.UnixNano() >= c.lt {
						fmt.Println(k + " - key deleted")
						delete(c.cv, k)
					}
				}
			}()
		},
		time.Second)

	return &c
}

func (c Cache) Set(key string, value interface{}) {
	c.cv[key] = cacheValue{time.Now(), value}
}

func (c Cache) Get(key string) (interface{}, error) {
	if v, ok := c.cv[key]; ok {
		return v.data, nil
	}
	return nil, errors.New(key + " - key absent")
}

func (c Cache) Delete(key string) {
	fmt.Println(key + " - key deleted")
	delete(c.cv, key)
}

func (c Cache) Free() {
	c.worker.Stop()
}
