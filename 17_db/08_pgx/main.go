package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "ibs"
	password = "ibs"
	dbname   = "ibs"
)

func main() {
	ctx := context.Background()
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, dbname)

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v\n", err)
	}
	defer conn.Close(ctx)

	fmt.Println("Успешное подключение к PostgreSQL!")

	createTableSQL := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL
		);`

	_, err = conn.Exec(ctx, createTableSQL)
	if err != nil {
		log.Fatalf("Ошибка при создании таблицы: %v\n", err)
	}
	fmt.Println("Таблица 'users' создана или уже существует.")

	insertSQL := "INSERT INTO users (name, email) VALUES ($1, $2) ON CONFLICT DO NOTHING;"
	_, err = conn.Exec(ctx, insertSQL, "Иван Иванов", "ivan@example.com")
	if err != nil {
		log.Fatalf("Ошибка при вставке данных: %v\n", err)
	}
	fmt.Println("Одиночная вставка через Exec завершена.")

	if err := insertWithBatch(ctx, conn); err != nil {
		log.Fatalf("Ошибка Batch: %v\n", err)
	}

	if err := insertWithCopy(ctx, conn); err != nil {
		log.Fatalf("Ошибка CopyFrom: %v\n", err)
	}

	rows, err := conn.Query(ctx, "SELECT id, name, email FROM users ORDER BY id;")
	if err != nil {
		log.Fatalf("Ошибка при выполнении запроса: %v\n", err)
	}
	defer rows.Close()

	fmt.Println("\nСодержимое таблицы 'users':")
	for rows.Next() {
		var id int
		var name, email string
		err = rows.Scan(&id, &name, &email)
		if err != nil {
			log.Fatalf("Ошибка при чтении строки: %v\n", err)
		}
		fmt.Printf("ID: %d, Имя: %s, Email: %s\n", id, name, email)
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("Ошибка при итерации по строкам: %v\n", err)
	}

	fmt.Println("\n Пример завершён успешно!")
}

// insertWithBatch отправляет несколько SQL-запросов за один сетевой обмен.
func insertWithBatch(ctx context.Context, conn *pgx.Conn) error {
	users := []struct {
		name  string
		email string
	}{
		{name: "Мария Машина", email: "maria@example.com"},
		{name: "Пётр Петров", email: "petr@example.com"},
		{name: "Ольга Орлова", email: "olga@example.com"},
	}

	batch := &pgx.Batch{}
	for _, user := range users {
		user := user
		batch.Queue(`
			INSERT INTO users (name, email)
			VALUES ($1, $2)
			ON CONFLICT (email) DO UPDATE SET name = EXCLUDED.name
			RETURNING id
		`, user.name, user.email).QueryRow(func(row pgx.Row) error {
			var id int
			if err := row.Scan(&id); err != nil {
				return err
			}
			fmt.Printf("Batch: %s сохранён с ID=%d\n", user.name, id)
			return nil
		})
	}

	results := conn.SendBatch(ctx, batch)
	// Close вычитывает все ответы, вызывает QueryRow callbacks
	// и освобождает соединение для следующих запросов.
	if err := results.Close(); err != nil {
		return fmt.Errorf("выполнить batch: %w", err)
	}

	return nil
}

// insertWithCopy массово загружает строки через PostgreSQL COPY protocol.
func insertWithCopy(ctx context.Context, conn *pgx.Conn) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("начать транзакцию: %w", err)
	}
	defer tx.Rollback(ctx) // после Commit вернёт pgx.ErrTxClosed

	// Удаляем только строки CopyFrom-демо, чтобы пример можно было запускать повторно.
	if _, err := tx.Exec(ctx, `DELETE FROM users WHERE email LIKE '%@copy.example'`); err != nil {
		return fmt.Errorf("удалить демонстрационные строки: %w", err)
	}

	rows := [][]any{
		{"Анна Соколова", "anna@copy.example"},
		{"Олег Орлов", "oleg@copy.example"},
		{"Елена Волкова", "elena@copy.example"},
		{"Сергей Котов", "sergey@copy.example"},
		{"Ирина Белова", "irina@copy.example"},
	}

	count, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"users"},
		[]string{"name", "email"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return fmt.Errorf("выполнить CopyFrom: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("зафиксировать CopyFrom: %w", err)
	}

	fmt.Printf("CopyFrom: загружено строк: %d\n", count)
	return nil
}
