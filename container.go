package di

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type container struct {
	mu        sync.RWMutex
	creators  map[string]func() any
	instances map[string]any
}

var globalContainer *container
var once sync.Once

func getContainer() *container {
	once.Do(func() {
		globalContainer = &container{
			mu:        sync.RWMutex{},
			creators:  make(map[string]func() any),
			instances: make(map[string]any),
		}
	})
	return globalContainer
}

func (c *container) injectable(typ reflect.Type, creator func() any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	selector := c.getSelector(typ)
	c.creators[selector] = creator
}

func (c *container) replace(typ reflect.Type, creator func() any, identifier ...string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	selector := c.getSelector(typ)
	instanceSelector := c.getSelector(typ, identifier...)
	c.creators[selector] = creator
	createdInstance := creator()
	c.instances[instanceSelector] = createdInstance
}

func (c *container) inject(typ reflect.Type, identifier ...string) (any, bool) {
	instanceSelector := c.getSelector(typ, identifier...)

	c.mu.RLock()
	if instance, ok := c.instances[instanceSelector]; ok {
		c.mu.RUnlock()
		return instance, true
	}
	c.mu.RUnlock()

	c.mu.RLock()
	if instance, ok := c.instances[instanceSelector]; ok {
		c.mu.RUnlock()
		return instance, true
	}
	c.mu.RUnlock()

	selector := c.getSelector(typ)
	creator, creatorExists := c.creators[selector]
	if !creatorExists {
		return nil, false
	}

	createdInstance := creator()
	c.mu.Lock()
	c.instances[instanceSelector] = createdInstance
	c.mu.Unlock()
	return createdInstance, true
}

func (c *container) destroy(typ reflect.Type, identifier ...string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	instanceSelector := c.getSelector(typ, identifier...)
	delete(c.instances, instanceSelector)
}

func (c *container) destroyAllMatching(match func(string) bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key := range c.instances {
		if match(key) {
			delete(c.instances, key)
		}
	}
	println("")
}

func (c *container) getSelector(typ reflect.Type, identifier ...string) string {
	typeName := typ.String()
	additionalKey := strings.Join(identifier, "_")
	return fmt.Sprintf("%s_%s", additionalKey, typeName)
}
