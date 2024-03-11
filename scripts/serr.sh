#!/bin/bash

echo -n "got: "
input="$(</dev/stdin)"
# output to stderr
>&2 echo "$input"
exit 0