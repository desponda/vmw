#!/bin/bash
for (( ; ; ))
do
   echo "Hitting url"
   curl http://localhost:8080/200
   sleep 5
   curl http://localhost:8080/503
   sleep 5
   curl http://localhost:8080/metrics
   sleep 10
done
