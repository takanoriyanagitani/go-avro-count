//go:build fast_count

package avro

import (
	ac "github.com/takanoriyanagitani/go-avro-count"
	util "github.com/takanoriyanagitani/go-avro-count/util"

	al "github.com/takanoriyanagitani/go-avro-count/input/avro/linkedin"
)

var StdinToAvroToCount util.IO[ac.Count] = al.StdinToAvroToCount
