package main

import (
	"bufio"
	"context"
	"fmt"
	"iter"
	"log"
	"os"
	"strconv"

	ac "github.com/takanoriyanagitani/go-avro-count"
	ia "github.com/takanoriyanagitani/go-avro-count/input/avro"
	util "github.com/takanoriyanagitani/go-avro-count/util"
)

func envVarByKey(key string) util.IO[string] {
	return func(_ context.Context) (string, error) {
		val, found := os.LookupEnv(key)
		switch found {
		case true:
			return val, nil
		default:
			return "", fmt.Errorf("env var %s missing", key)
		}
	}
}

var stdinAsFilenames util.IO[bool] = util.Bind(
	envVarByKey("ENV_STDIN_AS_FILENAMES"),
	util.Lift(strconv.ParseBool),
).Or(util.Of(false))

var filenames util.IO[iter.Seq[string]] = func(
	_ context.Context,
) (iter.Seq[string], error) {
	return func(yield func(string) bool) {
		var s *bufio.Scanner = bufio.NewScanner(os.Stdin)
		for s.Scan() {
			var filename string = s.Text()
			if !yield(filename) {
				return
			}
		}
	}, nil
}

var stdin2names2count util.IO[ac.Count] = util.Bind(
	filenames,
	ia.FilenamesToCountSingle,
)

var stdin2avro2count util.IO[ac.Count] = ia.StdinToAvroToCount

var stdin2count util.IO[ac.Count] = util.Bind(
	stdinAsFilenames,
	func(isFilenames bool) util.IO[ac.Count] {
		switch isFilenames {
		case true:
			return stdin2names2count
		default:
			return stdin2avro2count
		}
	},
)

func CountToStdout(c ac.Count) util.IO[util.Void] {
	return func(_ context.Context) (util.Void, error) {
		fmt.Printf("%v\n", c)
		return util.Empty, nil
	}
}

var stdin2count2stdout util.IO[util.Void] = util.Bind(
	stdin2count,
	CountToStdout,
)

func main() {
	_, e := stdin2count2stdout(context.Background())
	if nil != e {
		log.Printf("%v\n", e)
	}
}
