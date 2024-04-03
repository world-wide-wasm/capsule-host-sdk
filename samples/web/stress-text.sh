#!/bin/bash
# hey -n 3000 -c 1000 -m POST \

hey -n 3000 -c 1000 -m POST \
-H 'content-type: text/plain; charset=utf-8' \
-d '😄 Jane Doe' \
http://localhost:8080

