#!/bin/bash
# hey -n 3000 -c 1000 -m POST \

hey -n 3000 -c 1000 -m POST \
-H 'content-type: application/json; charset=utf-8' \
-d '{"firstName":"Bob","LastName":"Morane"}' \
http://localhost:8080
