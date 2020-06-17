package learner

import (
	"fmt"
	"sync"
)

type Learner struct {
	lock  sync.Mutex
	value string
	id    int
	addr  string
}

func NewLearner(id int, addr string) *Learner {

	l := &Learner{
		id:   id,
		addr: addr,
	}

	return l
}

func (l *Learner) Notify(value string, ret *bool) error {

	l.lock.Lock()
	defer l.lock.Unlock()

	fmt.Printf("[Learner:Notify] Got notified of value: %s\n", value)
	*ret = true
	l.value = value

	return nil
}
