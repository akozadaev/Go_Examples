package main

import "fmt"

type Celsius float64
type Fahrenheit float64

var temperature Celsius

func main() {
	temperature = 36.6
	fmt.Println(temperature)

	// ---- 1. Типобезопасность ----
	{
		var c Celsius = 25.0
		var f Fahrenheit = 77.0
		fmt.Println(c)
		fmt.Println(f)
		c = 36.6
		fmt.Println(c)
		//c = f // ./main.go:22:7: cannot use f (variable of float64 type Fahrenheit) as Celsius value in assignment
		fahrenheit := ToFahrenheit(c)
		fahrenheit1 := c.ToFahrenheit1()
		fmt.Println(fahrenheit)
		fmt.Println(fahrenheit1)

		celsius, err := NewCelsius(-300)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(celsius)
	}
}

// ---- 2. Явные функции для преобразования ----
func ToFahrenheit(c Celsius) Fahrenheit {
	return Fahrenheit(c*9.0/5.0 + 32.0)
}

// ---- 3. Метод для типа Celsius ----
func (c Celsius) ToFahrenheit1() Fahrenheit {
	return Fahrenheit(c*9.0/5.0 + 32.0)
}

// ---- 4. Конструктор с валидацией ----
func NewCelsius(value float64) (Celsius, error) {
	if value < -273.15 {
		return 0, fmt.Errorf("температура ниже абсолютного нуля")
	}
	return Celsius(value), nil
}
