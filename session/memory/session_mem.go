package memory

import (
	"container/list"
	"sync"
	"time"
	"web/session"
)

type SessionDate struct {
	sid          string
	timeAccessed time.Time
	value        map[interface{}]interface{}
}

type ProviderDate struct {
	lock     sync.Mutex
	sessions map[string]*list.Element
	list     *list.List
}

func (sess *SessionDate) Set(key, value interface{}) error {
	sess.value[key] = value
	prvd.SessionUpdate(sess.sid)
	return nil
}

func (sess *SessionDate) Get(key interface{}) interface{} {
	prvd.SessionUpdate(sess.sid)
	if v, ok := sess.value[key]; ok {
		return v
	} else {
		return nil
	}
}

func (sess *SessionDate) Delete(key interface{}) error {
	delete(sess.value, key)
	prvd.SessionUpdate(sess.sid)
	return nil
}

func (sess *SessionDate) SessionID() string {
	return sess.sid
}

func (prvd *ProviderDate) SessionInit(sid string) (session.Session, error) {
	prvd.lock.Lock()
	defer prvd.lock.Unlock()
	v := make(map[interface{}]interface{})
	newSess := &SessionDate{sid: sid, timeAccessed: time.Now(), value: v}
	element := prvd.list.PushBack(newSess)
	prvd.sessions[sid] = element
	return newSess, nil
}

func (prvd *ProviderDate) SessionRead(sid string) (session.Session, error) {
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

func (prvd *ProviderDate) SessionGC(maxLifeTime int64) {
	prvd.lock.Lock()
	defer prvd.lock.Unlock()
	// maxLifeTime(time for GC) is constant, session expiration time occurs at the end of list
	for {
		element := prvd.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*SessionDate).timeAccessed.Unix() + maxLifeTime) < time.Now().Unix() {
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
		element.Value.(*SessionDate).timeAccessed = time.Now()
		prvd.list.MoveToFront(element)
		return nil
	}
	return nil
}

var prvd = &ProviderDate{sessions: make(map[string]*list.Element), list: list.New()}

func init() {
	session.Register("memory", prvd)
}
