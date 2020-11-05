#!/bin/bash

function test
{
  go test -v $(go list ./... | grep -v vendor) --count 1 -covermode=atomic
}

CMD=$1
shift
$CMD $*
