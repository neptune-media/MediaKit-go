package stats

import "sync"

type Stats struct {
	counterLock sync.Mutex
	counters    map[string]*Counter
}

func (s *Stats) GetValuesMap() map[string]any {
	s.counterLock.Lock()
	defer s.counterLock.Unlock()

	m := make(map[string]any)
	for k, v := range s.counters {
		m[k] = v.Value()
	}

	return m
}

func (s *Stats) DeleteCounter(name string) {
	s.counterLock.Lock()
	defer s.counterLock.Unlock()

	delete(s.counters, name)
}

func (s *Stats) GetCounter(name string) *Counter {
	s.counterLock.Lock()
	defer s.counterLock.Unlock()

	if c, ok := s.counters[name]; ok {
		return c
	}

	return nil
}

func (s *Stats) SetCounter(name string, counter *Counter) {
	s.counterLock.Lock()
	defer s.counterLock.Unlock()

	s.counters[name] = counter
}

type Counter struct {
	lock  sync.Mutex
	value float64
}

func (c *Counter) Add(delta float64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.value += delta
}

func (c *Counter) Increment() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.value++
}

func (c *Counter) Value() float64 {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.value
}

type Gauge struct {
	lock  sync.Mutex
	value float64
}

func (g *Gauge) Set(value float64) {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.value = value
}

func (g *Gauge) Value() float64 {
	g.lock.Lock()
	defer g.lock.Unlock()
	return g.value
}
