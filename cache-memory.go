package cachememory

// TODO:  Delete fmt.Print*

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/zhashkevych/scheduler"
)

type Cache interface {
	Set(key string, value interface{}, lifeTimeSec int) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	Free()
}

func New(lifeTimeSec int) Cache {
	return NewCacheMem(lifeTimeSec)
}

//================================  Implementation  ===================================

const defaultLifeTimeSec int = 60 // default life time second

type cacheValue struct {
	exp  time.Time
	data interface{}
}

type CacheMem struct {
	worker *scheduler.Scheduler
	ctx    context.Context
	lt     time.Duration // life time nanosecond
	cv     sync.Map
}

func NewCacheMem(lifeTimeSec int) *CacheMem {
	if lifeTimeSec < 1 {
		lifeTimeSec = defaultLifeTimeSec
	}

	c := CacheMem{
		ctx:    context.Background(),
		worker: scheduler.NewScheduler(),
		lt:     time.Duration(lifeTimeSec) * time.Second, // life time nanosecond
	}

	c.worker.Add(c.ctx,
		func(_ context.Context) {
			func() {

				c.cv.Range(func(key, value interface{}) bool {

					if time.Now().After(value.(cacheValue).exp) {
						fmt.Println(fmt.Sprint(key) + " - key deleted")
						c.cv.Delete(key)
					}
					return true
				})
			}()
		},
		time.Second)

	return &c
}

func (c *CacheMem) Set(key string, value interface{}, lifeTimeSec int) error {
	if lifeTimeSec == 0 {
		c.cv.Store(key, cacheValue{time.Now().Add(c.lt), value})
	} else {
		c.cv.Store(key, cacheValue{time.Now().Add(time.Duration(lifeTimeSec) * time.Second), value})
	}
	return nil
}

func (c *CacheMem) Get(key string) (interface{}, error) {

	if v, ok := c.cv.Load(key); ok {
		return v.(cacheValue).data, nil
	}
	return nil, errors.New(key + " - key absent")
}

func (c *CacheMem) Delete(key string) error {

	_, ok := c.cv.LoadAndDelete(key)
	if ok {
		fmt.Println(key + " - key deleted")
	} else {
		return errors.New(key + " - key absent")
	}
	return nil
}

func (c *CacheMem) Free() {
	c.worker.Stop()
}
