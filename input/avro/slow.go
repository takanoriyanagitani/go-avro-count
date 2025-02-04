//go:build !fast_count

package avro

import (
	ac "github.com/takanoriyanagitani/go-avro-count"
	ah "github.com/takanoriyanagitani/go-avro-count/input/avro/hamba"
	util "github.com/takanoriyanagitani/go-avro-count/util"
)

var StdinToAvroToCount util.IO[ac.Count] = ah.StdinToAvroToCount
