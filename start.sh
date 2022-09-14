#!/bin/bash

./statuspage &
./simulator &

wait -n

exit $?