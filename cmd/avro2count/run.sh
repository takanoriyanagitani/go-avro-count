#!/bin/sh

#cat sample.d/sample.json | jsons2maps2avro > ./sample.d/sample.avro

cat sample.d/sample.avro | ./avro2count

export ENV_MULTI_THREAD=true
export GOMAXPROCS=4

printf \
	'%s\n' \
	./sample.d/sample.avro \
	./sample.d/sample.avro |
	ENV_STDIN_AS_FILENAMES=true ./avro2count
