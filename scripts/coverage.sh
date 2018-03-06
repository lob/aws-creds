#!/bin/bash

THRESHOLD=85

PERCENTS=$(cat ./coverage.out | grep "coverage" | tr -s " " | cut -d ' ' -f 3 | sed 's/%//' | awk -F. '{print $1}')

for PERCENT in $PERCENTS; do
  if (( $PERCENT < $THRESHOLD )); then
    exit 1
  fi
done
