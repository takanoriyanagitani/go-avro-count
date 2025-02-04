//go:build fast_count

package avro

import (
	ac "github.com/takanoriyanagitani/go-avro-count"
	al "github.com/takanoriyanagitani/go-avro-count/input/avro/linkedin"
	util "github.com/takanoriyanagitani/go-avro-count/util"
)

var StdinToAvroToCount util.IO[ac.Count] = al.StdinToAvroToCount
