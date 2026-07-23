#!/usr/bin/env bash
# Удалить пользователя с id=2
curl -s -o /dev/null -w "HTTP %{http_code}\n" -X DELETE \
  http://localhost:8080/users/3
