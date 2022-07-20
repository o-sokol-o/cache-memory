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
	Set(key string, value interface{})
	Get(key string) (interface{}, error)
	Delete(key string)
	Free()
}

func New(lifeTimeSec int) Cache {
	return NewCacheMem(lifeTimeSec)
}

//===================================================================

const defaultLifeTimeSec int = 60 // default life time second

type cacheValue struct {
	exp  time.Time
	data interface{}
}

type CacheMem struct {
	worker *scheduler.Scheduler
	ctx    context.Context
	lt     time.Duration // life time nanosecond
	cv     sync.Map      // normal map: cv     map[string]cacheValue
}

func NewCacheMem(lifeTimeSec int) *CacheMem {
	if lifeTimeSec < 1 {
		lifeTimeSec = defaultLifeTimeSec
	}

	c := CacheMem{
		ctx:    context.Background(),
		worker: scheduler.NewScheduler(),
		lt:     time.Duration(lifeTimeSec) * time.Second, // life time nanosecond
		// normal map: cv:     make(map[string]cacheValue),
	}

	c.worker.Add(c.ctx,
		func(_ context.Context) {
			func() {
				/* normal map:
				for key, val := range c.cv {
					if time.Now().After(val.exp) {
						fmt.Println(key + " - key deleted")
						// delete(c.cv, key)
						c.cv.Delete(key)
					}
				}
				*/

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

func (c *CacheMem) Set(key string, value interface{}) {
	c.cv.Store(key, cacheValue{time.Now().Add(c.lt), value}) // normal map: c.cv[key] = cacheValue{time.Now().Add(c.lt), value}
}

func (c *CacheMem) Get(key string) (interface{}, error) {
	// normal map: if v, ok := c.cv[key]; ok {
	if v, ok := c.cv.Load(key); ok {
		return v.(cacheValue).data, nil
	}
	return nil, errors.New(key + " - key absent")
}

func (c *CacheMem) Delete(key string) {
	fmt.Println(key + " - key deleted")
	c.cv.Delete(key) // normal map: delete(c.cv, key)
}

func (c *CacheMem) Free() {
	c.worker.Stop()
}
