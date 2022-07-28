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

type Resolution int

const (
	ResolutionDefault Resolution = iota
	ResolutionSeconds
	ResolutionMinutes
	ResolutionHours
	ResolutionDays
)

type Cache interface {
	Set(key string, value interface{}, lifeTime int) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	Free()
}

func New(rs Resolution) Cache {
	return NewCacheMem(rs)
}

//================================  Implementation  ===================================

type cacheValue struct {
	exp  time.Time
	data interface{}
}

type CacheMem struct {
	worker *scheduler.Scheduler
	ctx    context.Context
	lt     time.Duration // life time nanosecond
	rt     time.Duration
	cv     sync.Map
}

func NewCacheMem(rs Resolution) *CacheMem {

	var resolutionTime time.Duration
	var defaultLifeTime time.Duration
	switch rs {
	case ResolutionMinutes:
		resolutionTime = time.Minute
		defaultLifeTime = 60 * time.Minute
	case ResolutionHours:
		resolutionTime = time.Hour
		defaultLifeTime = 24 * time.Hour
	case ResolutionDays:
		resolutionTime = 24 * time.Hour
		defaultLifeTime = 3 * 24 * time.Hour
	default:
		resolutionTime = time.Second
		defaultLifeTime = 60 * time.Second
	}

	c := CacheMem{
		ctx:    context.Background(),
		worker: scheduler.NewScheduler(),
		lt:     defaultLifeTime,
		rt:     resolutionTime,
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
		resolutionTime)

	return &c
}

func (c *CacheMem) Set(key string, value interface{}, lifeTime int) error {
	if lifeTime == 0 {
		c.cv.Store(key, cacheValue{time.Now().Add(c.lt), value})
	} else {
		c.cv.Store(key, cacheValue{time.Now().Add(time.Duration(lifeTime) * c.rt), value})
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
