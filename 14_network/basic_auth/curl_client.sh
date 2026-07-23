curl -u admin:secret http://localhost:8080/admin
#
curl -u admin:secret -X POST \
  -H "Content-Type: application/json" \
  -d '{"username":"test","role":"admin"}' \
  http://localhost:8080/admin