#!/bin/bash

echo "i will dump a folder: $1 into: $2"

if [ ! -d "$1" ]; then
    >&2 echo "error $1 does not exist or is not a dir"
    exit 1
fi

tar -czf $2 $1

echo "done"
ls -loah $2
