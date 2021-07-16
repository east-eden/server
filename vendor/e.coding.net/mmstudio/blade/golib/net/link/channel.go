package link

import (
	"sync"
)

type ChannelKey interface{}

type Channel interface {
	Len() int
	Fetch(callback func(Session))
	Get(key ChannelKey) Session
	Put(key ChannelKey, session Session) Session
	Remove(key ChannelKey) bool
	FetchAndRemove(callback func(Session))
	Close()
}

func NewChannel() Channel {
	return &_channel{
		sessions: make(map[ChannelKey]Session),
	}
}

type _channel struct {
	mutex    sync.RWMutex
	sessions map[ChannelKey]Session

	// channel state
	State interface{}
}

func (channel *_channel) Len() int {
	channel.mutex.RLock()
	defer channel.mutex.RUnlock()
	return len(channel.sessions)
}

func (channel *_channel) Fetch(callback func(session Session)) {
	channel.mutex.RLock()
	defer channel.mutex.RUnlock()
	for _, session := range channel.sessions {
		callback(session)
	}
}

func (channel *_channel) Get(key ChannelKey) Session {
	channel.mutex.RLock()
	defer channel.mutex.RUnlock()
	session, _ := channel.sessions[key]
	return session
}

func (channel *_channel) Put(key ChannelKey, session Session) (oldSession Session) {
	channel.mutex.Lock()
	defer channel.mutex.Unlock()
	if session, exists := channel.sessions[key]; exists {
		oldSession = session
		channel.remove(key, session)
	}
	session.AddCloseCallback(channel, key, func() {
		channel.Remove(key)
	})
	channel.sessions[key] = session
	return
}

func (channel *_channel) remove(key ChannelKey, session Session) {
	session.RemoveCloseCallback(channel, key)
	delete(channel.sessions, key)
}

func (channel *_channel) Remove(key ChannelKey) bool {
	channel.mutex.Lock()
	defer channel.mutex.Unlock()
	session, exists := channel.sessions[key]
	if exists {
		channel.remove(key, session)
	}
	return exists
}

func (channel *_channel) FetchAndRemove(callback func(Session)) {
	channel.mutex.Lock()
	defer channel.mutex.Unlock()
	for key, session := range channel.sessions {
		session.RemoveCloseCallback(channel, key)
		delete(channel.sessions, key)
		callback(session)
	}
}

func (channel *_channel) Close() {
	channel.mutex.Lock()
	defer channel.mutex.Unlock()
	for key, session := range channel.sessions {
		channel.remove(key, session)
	}
}
