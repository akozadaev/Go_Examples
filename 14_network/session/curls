# Первый запрос - устанавливает сессию и получает Set-Cookie
curl -c cookies.txt http://localhost:8080/

# Второй запрос - отправляет сохранённые cookies и получает значение из сессии
curl -b cookies.txt http://localhost:8080/