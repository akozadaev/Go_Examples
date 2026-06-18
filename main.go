package main

import "fmt"

type Celsius float64
type Fahrenheit float64

var temperature Celsius

func main() {
	temperature = 36.6
	fmt.Println(temperature)

	//Свойссства и назначение
	//1. Типобезопасность (Type Safety)
	{

		var c Celsius = 25.0
		var f Fahrenheit = 77.0
		fmt.Println(c)
		fmt.Println(f)
		c = 36.6
		fmt.Println(c)
		// ОШИБКА!:
		// c = f  // ошибка компиляции: cannot use f (type Fahrenheit) as type Celsius
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

func ToFahrenheit(c Celsius) Fahrenheit {
	return Fahrenheit(c*9.0/5.0 + 32.0)
}

func ToString(c Celsius) string {
	return fmt.Sprintf("%.2f°C", c)
}

// 2. Самодокументирование кода
//Когда вы видите:
//func SetTemperature(temp Celsius) — сразу понятно, что функция ожидает температуру в градусах Цельсия.
// 3.Добавление методов к типу

func (c Celsius) ToFahrenheit1() Fahrenheit {
	return Fahrenheit(c*9.0/5.0 + 32.0)
}

// 4. Контроль над валидацией и инициализацией
func NewCelsius(value float64) (Celsius, error) {
	if value < -273.15 {
		return 0, fmt.Errorf("температура ниже абсолютного нуля")
	}
	return Celsius(value), nil
}

//5. Упрощение рефакторинга и поддержки
//Если вы решите изменить внутреннее представление (например, хранить температуру в микроградусах как int64), интерфейс остается тем же:
