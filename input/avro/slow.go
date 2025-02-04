//go:build !fast_count

package avro

import (
	"errors"

	ac "github.com/takanoriyanagitani/go-avro-count"
	ah "github.com/takanoriyanagitani/go-avro-count/input/avro/hamba"
	util "github.com/takanoriyanagitani/go-avro-count/util"
)

var ErrNotImplemented error = errors.New("not implemented")

var StdinToAvroToCount util.IO[ac.Count] = ah.StdinToAvroToCount

var FilenameToAvroToCount func(string) util.IO[ac.Count] = func(
	_ string,
) util.IO[ac.Count] {
	return util.Err[ac.Count](ErrNotImplemented)
}
