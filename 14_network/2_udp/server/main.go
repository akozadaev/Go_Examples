package main

import (
	"fmt"
	"net"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:8081")
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("UDP server listening on :8081")

	// Для учебного echo достаточно 1024. Если дейтаграмма больше len(buffer),
	// ReadFromUDP вернёт усечённые данные - в реальном коде буфер обычно делают больше.
	buffer := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Получили ошибку:", err)
			continue
		}
		msg := string(buffer[:n])
		fmt.Printf("From %s: %s", clientAddr, msg)

		// Отправляем ответ
		b, err := conn.WriteToUDP([]byte("Привет, привет!"), clientAddr)
		if err != nil {
			fmt.Println("Получили ошибку:", err)
			continue
		}
		fmt.Println(b)
	}
}
