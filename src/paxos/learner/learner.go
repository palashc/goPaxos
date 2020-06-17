package learner

import "sync"

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

	*ret = true
	l.value = value

	return nil
}
