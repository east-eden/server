// +build go1.13

package suturev4

import "errors"

func isErr(err error, target error) bool {
	return errors.Is(err, target)
}
