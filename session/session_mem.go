package session

import (
	"container/list"
	"sync"
	"time"
)

type SessionDate struct {
	sid          string
	timeAccessed time.Time
	value        map[interface{}]interface{}
}

type ProviderDate struct {
	lock    sync.Mutex
	session map[string]*list.Element
	list    *list.List
}

var prvd = &ProviderDate{list: list.New()}

func (sess *SessionDate) Set(key, value interface{}) error {
	sess.value[key] = value
	return nil
}

func (sess *SessionDate) Get(key interface{}) interface{} {
	if v, ok := sess.value[key]; ok {
		return v
	} else {
		return nil
	}
}

func (sess *SessionDate) Delete(key interface{}) error {
	delete(sess.value, key)
	return nil
}

func (sess *SessionDate) SessionID() string {
	return sess.sid
}

func (prvd *ProviderDate) SessionInit(sid string) (Session, error) {
	prvd.lock.Lock()
	defer prvd.lock.Unlock()
	v := make(map[interface{}]interface{}, 0)
	newSess := &SessionDate{sid: sid, timeAccessed: time.Now(), value: v}
	element := prvd.list.PushBack(newSess)
	prvd.sessions[sid] = element
	return newSess, nil
}

func (prvd *ProviderDate) SessionRead(sid string) (Session, error) {
	if element, ok := prvd.sessions[sid]; ok {
		return element.Value.(*SessionDate), nil
	} else {
		sess, err := prvd.SessionInit(sid)
		return sess, err
	}
}

func (prvd *ProviderDate) SessionDestory(sid string) error {
	if element, ok := prvd.sessions[sid]; ok {
		delete(prvd.sessions, sid)
		prvd.list.Remove(element)
		return nil
	}
	return nil
}

func (prvd *ProviderDate) SessinGC(maxLifeTime int64) {
	prvd.lock.Lock()
	defer prvd.lock.Unlock()

	for {
		element := prvd.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*SessionStore).timeAccessed.Unix() + maxLifeTime) < time.Now() {
			prvd.list.Remove(element)
		} else {
			break
		}
	}
}

func (prvd *ProviderDate) SessionUpdate(sid string) error {
	prvd.lock.Lock()
	defer prvd.lock.Unlock()
	if element, ok := prvd.sessions[sid]; ok {
		element.Value.(*SessionStore).timeAccessed = time.Now()
		prvd.list.MoveToFront(element)
		return nil
	}
	return nil
}
