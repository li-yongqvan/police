#!/usr/bin/env bash
curl -s -w "\nHTTP:%{http_code}\n" -X POST http://127.0.0.1:8888/api/v1/demo-login \
  -H 'Content-Type: application/json' \
  -d '{"role":"student"}'
echo "---"
curl -s -w "\nHTTP:%{http_code}\n" -X POST http://127.0.0.1:8888/api/v1/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"demo_student","password":"demo123456"}'
