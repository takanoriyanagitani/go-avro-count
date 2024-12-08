//go:build !fast_count

package avro

import (
	ac "github.com/takanoriyanagitani/go-avro-count"
	util "github.com/takanoriyanagitani/go-avro-count/util"

	ah "github.com/takanoriyanagitani/go-avro-count/input/avro/hamba"
)

var StdinToAvroToCount util.IO[ac.Count] = ah.StdinToAvroToCount
