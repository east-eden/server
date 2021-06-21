// +build !go1.13

package suturev4

func isErr(err error, target error) bool {
	return err == target
}
