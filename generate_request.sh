#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: $0 <number_of_requests>"
  exit 1
fi

if ! [[ "$1" =~ ^[0-9]+$ ]]; then
  echo "Number of requests must be a number"
  exit 1
fi

if [ "$1" -lt 1 ]; then
  echo "Number of requests must be greater than 0"
  exit 1
fi

NB_REQUESTS=$1
success_count=0

for ((i=1; i<=NB_REQUESTS; i++)); do
  status=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:8080/process")
  if [ "$status" -eq 200 ]; then
    echo "Request $i done (status $status)"
    success_count=$((success_count+1))
  else
    echo "Error on request $i (status $status)"
    sleep $((i/10))
  fi
done

success_rate=$(awk "BEGIN {printf \"%.2f\", (${success_count}/$NB_REQUESTS)*100}")
echo -e "\n\n=======\nSuccess rate: $success_rate% ($success_count/$NB_REQUESTS)\n======="
