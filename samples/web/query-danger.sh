#!/bin/bash
curl http://localhost:8080 \
-H 'content-type: application/json; charset=utf-8' \
-d '{"content":"<h1>\"🎉 tada!!!\"</h2>","something":"👋 hello world 🌍"}'
