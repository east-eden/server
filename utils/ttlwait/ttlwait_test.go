package ttlwait

import (
	"fmt"
	"testing"
	"time"
)

func TestTTLWait(t *testing.T) {
	ttlwait := NewTTLWait(time.Second * 3)
	go func() {
		ttlwait.Done()
		fmt.Println("run new go func")
	}()
	ttlwait.Wait()
	fmt.Println("exit perfect")
}

func TestTTLWaitWithTimeout(t *testing.T) {
	ttlwait := NewTTLWait(time.Second * 3)
	go func() {
		fmt.Println("run new go func timeout")
	}()
	ttlwait.Wait()
	fmt.Println("exit perfect")
}

func TestTTLWaitGroupWithTimeout(t *testing.T) {
	group := NewTTLWaitGroup(time.Second*3, 3)

	go func() {
		defer group.Done()
		fmt.Println("run group 1")
	}()

	go func() {
		defer group.Done()
		fmt.Println("run group 2")
	}()

	timeout := group.Wait()
	fmt.Println("exit with timeout = ", timeout)
}
