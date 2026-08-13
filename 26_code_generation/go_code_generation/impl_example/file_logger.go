package impl_example

// Этот файл демонстрирует, как использовать `impl` для генерации заглушек методов.
//
// Использование:
//   impl 'recv *FileLogger' impl_example.Logger
//
// Это создаст заглушки всех методов интерфейса Logger для структуры FileLogger.

import (
	"fmt"
	"os"
)

// FileLogger реализует интерфейс Logger для записи в файл.
// Заглушки методов были сгенерированы с помощью команды:
//
//	impl 'recv *FileLogger' impl_example.Logger
type FileLogger struct {
	Filename string
}

// Реализация методов интерфейса Logger (заглушки созданы через impl):

func (f *FileLogger) Log(msg string) {
	f.logToFile(fmt.Sprintf("[INFO] %s\n", msg))
}

func (f *FileLogger) Logf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	f.logToFile(fmt.Sprintf("[INFO] %s\n", msg))
}

func (f *FileLogger) Error(msg string) {
	f.logToFile(fmt.Sprintf("[ERROR] %s\n", msg))
}

func (f *FileLogger) Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	f.logToFile(fmt.Sprintf("[ERROR] %s\n", msg))
}

// logToFile - вспомогательный метод для записи в файл
func (f *FileLogger) logToFile(msg string) {
	file, err := os.OpenFile(f.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Ошибка открытия файла: %v\n", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(msg); err != nil {
		fmt.Printf("Ошибка записи в файл: %v\n", err)
	}
}
