#!/bin/sh

#cat sample.d/sample.json | jsons2maps2avro > ./sample.d/sample.avro

cat sample.d/sample.avro | ./avro2count
