package main

import (
	"context"
	"fmt"
	"log"

	ac "github.com/takanoriyanagitani/go-avro-count"
	util "github.com/takanoriyanagitani/go-avro-count/util"
)

import (
	ia "github.com/takanoriyanagitani/go-avro-count/input/avro"
)

var stdin2avro2count util.IO[ac.Count] = ia.StdinToAvroToCount

func CountToStdout(c ac.Count) util.IO[util.Void] {
	return func(_ context.Context) (util.Void, error) {
		fmt.Printf("%v\n", c)
		return util.Empty, nil
	}
}

var stdin2count2stdout util.IO[util.Void] = util.Bind(
	stdin2avro2count,
	CountToStdout,
)

func main() {
	_, e := stdin2count2stdout(context.Background())
	if nil != e {
		log.Printf("%v\n", e)
	}
}
