package kvs

import (
	"context"
	"time"
)

const (
	defaultAutoSyncInterval = 5 * time.Minute
	defaultOperationTimeout = time.Second * time.Duration(10)
)

func init() {
	InstallRegistryOptionsWatchDog(func(cc *RegistryOptions) {
		cc.BaseStoreWriteOptionsInner = NewWriteOptions(
			WithWriteOptionTTL(cc.TTL),
			WithWriteOptionKeepAlive(cc.KeepAlive),
		)
		cc.BaseStoreReadOptionsInner = NewReadOptions(WithReadOptionConsistent(cc.UsingConsistentRead))
	})
}

//go:generate optiongen --option_with_struct_name=true --v=true
func RegistryOptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"TTL":                        time.Duration(time.Duration(30) * time.Second),
		"KeepAlive":                  true,
		"UsingConsistentRead":        false,
		"BaseStoreWriteOptionsInner": (*WriteOptions)(nil),
		"BaseStoreReadOptionsInner":  (*ReadOptions)(nil),
	}
}

//go:generate optiongen --option_with_struct_name=true --v=true
func WatchOptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		// Other options for implementations of the interface can be stored in a context
		"Context":       context.Context(nil),
		"Heartbeat":     time.Duration(time.Duration(10) * time.Second),
		"KVReadOptions": (*ReadOptions)(nil),
	}
}

//go:generate optiongen --option_with_struct_name=true --v=true
func GetOptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"Context":       context.Context(nil),
		"KVReadOptions": (*ReadOptions)(nil),
	}
}

//go:generate optiongen  --v=true
func StoreOptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"Endpoints":         []string{},
		"Name":              "",
		"ConnectionTimeout": time.Duration(time.Duration(20) * time.Second),
		"OperationTimeout":  time.Duration(defaultOperationTimeout),
		"AutoSyncInterval":  time.Duration(defaultAutoSyncInterval),
		"Username":          "",
		"Password":          "",
		"ETCDV3Config":      interface{}(nil),
	}
}

//go:generate optiongen --option_with_struct_name=true --v=true
func WriteOptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"TTL": time.Duration(0),
		// If true, the client will keep the lease alive in the background for stores that are allowing it.
		"KeepAlive":               true,
		"KeepAliveLostChanSetter": (func(chan struct{}))(nil),
	}
}

//go:generate optiongen --option_with_struct_name=true --v=true
func ReadOptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		// Consistent defines if the behavior of a Get operation is
		// linearizable or not. Linearizability allows us to 'see'
		// objects based on a real-time total order as opposed to
		// an arbitrary order or with stale values ('inconsistent'
		// scenario).
		"Consistent": false,
	}
}

//go:generate optiongen --option_with_struct_name=true --v=true
func LockOptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		// value to associate with the lock
		"Value": []byte(nil),
		// expiration ttl associated with the lock
		"TTL": time.Duration(0),
		// chan used to control and stop the session ttl renewal for the lock
		"RenewLock": chan struct{}(nil),
		//  If true, the value will be deleted when the lock is unlocked or expires
		"DeleteOnUnlock": false,
	}
}
