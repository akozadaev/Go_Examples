package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Подключено к серверу. Введите сообщение (нажмите Enter для отправки):")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if _, err := conn.Write([]byte(text + "\n")); err != nil {
			fmt.Println("Ошибка отправки:", err)
			return
		}

		response := make([]byte, 1024)
		n, err := conn.Read(response)
		if err != nil {
			fmt.Println("Ошибка чтения ответа:", err)
			return
		}
		fmt.Printf("Сервер: %s", string(response[:n]))
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка ввода:", err)
	}
}
