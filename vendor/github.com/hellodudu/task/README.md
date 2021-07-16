[![GoDoc](https://godoc.org/github.com/hellodudu/task?status.svg)](https://godoc.org/github.com/hellodudu/task)
[![Go Report Card](https://goreportcard.com/badge/github.com/hellodudu/task)](https://goreportcard.com/report/github.com/hellodudu/task)

# task
This simple package provides a concurrency and configued tasks handle queue.

## Example

```go
// new a tasker with options
tasker := NewTasker(1)
tasker.Init(
    WithContextDoneCb(func() {
        fmt.Println("tasker context done...")
    }),

    WithTimeout(time.Second*5),

    WithSleep(time.Millisecond*100),

    WithUpdateCb(func() {
        fmt.Println("tasker update...")
    }),
)

// add a simple task
tasker.Add(
    context.Background(),
    func(context.Context, ...interface{}) error {
        fmt.Println("task handled")
    },
)

// add task with parameters
tasker.Add(
    context.Background(),
    func(ctx context.Context, p ...interface{}) error {
        p1 := p[0].(int)
        p2 := p[1].(string)
        fmt.Println("task handled with parameters:", p1, p2)
    },
    1,
    "name",
)

// begin run
go tasker.Run(context.Background())

// stop
tasker.Stop()
```
