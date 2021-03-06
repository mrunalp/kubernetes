/*
Copyright 2014 Google Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cache

import (
	"time"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/util"
	"github.com/golang/glog"
)

// Enumerator should be able to return the list of objects to be synced with
// one object at a time.
type Enumerator interface {
	Len() int
	Get(index int) (ID string, object interface{})
}

// GetFunc should return an enumerator that you wish the Poller to proccess.
type GetFunc func() (Enumerator, error)

// Poller is like Reflector, but it periodically polls instead of watching.
// This is intended to be a workaround for api objects that don't yet support
// watching.
type Poller struct {
	getFunc GetFunc
	period  time.Duration
	store   Store
}

// NewPoller constructs a new poller. Note that polling probably doesn't make much
// sense to use along with the FIFO queue. The returned Poller will call getFunc and
// sync the objects in 'store' with the returned Enumerator, waiting 'period' between
// each call. It probably only makes sense to use a poller if you're treating the
// store as read-only.
func NewPoller(getFunc GetFunc, period time.Duration, store Store) *Poller {
	return &Poller{
		getFunc: getFunc,
		period:  period,
		store:   store,
	}
}

// Run begins polling. It starts a goroutine and returns immediately.
func (p *Poller) Run() {
	go util.Forever(func() {
		e, err := p.getFunc()
		if err != nil {
			glog.Errorf("failed to list: %v", err)
			return
		}
		p.sync(e)
	}, p.period)
}

func (p *Poller) sync(e Enumerator) {
	current := p.store.Contains()
	for i := 0; i < e.Len(); i++ {
		id, object := e.Get(i)
		p.store.Update(id, object)
		current.Delete(id)
	}
	// Delete all the objects not found.
	for id := range current {
		p.store.Delete(id)
	}
}
