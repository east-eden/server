package gls

import (
	"sync"

	"e.coding.net/mmstudio/blade/golib/gid"
)

// should use ConcurrentMap
var glsMap sync.Map

func init() {
}

type copiable interface {
	Copy() interface{}
}

// ResetGls reset the goroutine local storage for specified goroutine
func ResetGls(id gid.Type, initialValue map[interface{}]interface{}) {
	glsMap.Store(id, initialValue)
}

// DeleteGls remove goroutine local storage for specified goroutine
func DeleteGls(id gid.Type) {
	glsMap.Delete(id)
}

// GetGls get goroutine local storage for specified goroutine
// if the goroutine did not set gls, it will return nil
func GetGls(id gid.Type) map[interface{}]interface{} {
	val, ok := glsMap.Load(id)
	if ok {
		return val.(map[interface{}]interface{})
	}
	return nil
}

// WithGls setup and teardown the gls in the wrapper.
// go WithGls(func(){}) will add gls for the new goroutine.
// The gls will be removed once goroutine exit
func WithGls(f func()) func() {
	parentGls := GetGls(gid.Get())
	// parentGls can not be used in other goroutine, otherwise not thread safe
	// make a deep for child goroutine
	childGls := map[interface{}]interface{}{}
	for k, v := range parentGls {
		asCopiable, ok := v.(copiable)
		if ok {
			childGls[k] = asCopiable.Copy()
		} else {
			childGls[k] = v
		}
	}
	return func() {
		gid := gid.Get()
		ResetGls(gid, childGls)
		defer DeleteGls(gid)
		f()
	}
}

// WithEmptyGls works like WithGls, but do not inherit gls from parent goroutine.
func WithEmptyGls(f func()) func() {
	// do not inherit from parent gls
	return func() {
		goid := gid.Get()
		ResetGls(goid, make(map[interface{}]interface{}))
		defer DeleteGls(goid)
		f()
	}
}

// Get key from goroutine local storage
func Get(key interface{}) interface{} {
	glsMap := GetGls(gid.Get())
	if glsMap == nil {
		return nil
	}
	return glsMap[key]
}

func SetMutiple(kv ...interface{}) {
	glsMap := GetGls(gid.Get())
	if glsMap == nil {
		return
	}
	if len(kv)%2 == 1 {
		return
	}
	var k interface{}
	for i, s := range kv {
		if i%2 == 0 {
			k = s
			continue
		}
		glsMap[k] = s
	}
}

// Set key and element to goroutine local storage
func Set(key interface{}, value interface{}) {
	glsMap := GetGls(gid.Get())
	if glsMap == nil {
		return
	}
	glsMap[key] = value
}

// IsGlsEnabled test if the gls is available for specified goroutine
func IsGlsEnabled(id gid.Type) bool {
	return GetGls(id) != nil
}
