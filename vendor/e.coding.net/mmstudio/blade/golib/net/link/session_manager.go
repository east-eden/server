package link

import (
	"sync"
)

type Manager struct {
	sms         ConcurrentMapUint64Session
	disposeOnce sync.Once
	disposeWait sync.WaitGroup
	disposed    bool
}

func NewManager() *Manager {
	manager := &Manager{}
	manager.sms = NewConcurrentMapUint64Session()
	return manager
}

func (m *Manager) Dispose() {
	m.disposeOnce.Do(func() {
		m.disposed = true
		for tuple := range m.sms.Iter() {
			_ = tuple.Val.Close("manager_dispose")
		}
		m.disposeWait.Wait()
	})
}

func (m *Manager) SessionCount() int {
	if m.disposed {
		return 0
	}
	return m.sms.Count()
}

func (m *Manager) Range(f func(Session) bool) {
	for tuple := range m.sms.Iter() {
		if !f(tuple.Val) {
			break
		}
	}
}

func (m *Manager) NewSession(tans Transporter, spec *Spec) Session {
	session := NewSession(tans, spec)
	session.SetManager(m)
	m.putSession(session)
	return session
}

func (m *Manager) GetSession(sessionID uint64) (Session, bool) {
	return m.sms.Get(sessionID)
}

func (m *Manager) putSession(session *_session) {
	if m.disposed {
		_ = session.Close("manager_disposed")
		return
	}
	m.sms.Set(session.id, session)
	m.disposeWait.Add(1)
}

func (m *Manager) delSession(session *_session) {
	if m.sms.Has(session.id) {
		m.sms.Remove(session.id)
		m.disposeWait.Done()
	}
}
