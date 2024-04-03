#!/bin/bash
curl http://localhost:8080/pouet \
-H 'content-type: application/json; charset=utf-8' \
-d '{"firstName":"Bob","LastName":"Morane"}'
