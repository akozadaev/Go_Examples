#!/usr/bin/env bash
# Обновить пользователя с id=1
curl -s -X PUT -H "Content-Type: application/json" \
  -d '{"name":"Alexey Updated"}' \
  http://localhost:8080/users/1
echo
