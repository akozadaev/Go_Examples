package person

import "fmt"

type Person struct {
	name string
	age  int
}

func (p Person) Hello() {
	//fmt.Println("Hello", p.Name)
}

func (p Person) Name() string {
	return p.name
}

func (p *Person) SetName(name string) error {
	if name == "" {
		return fmt.Errorf("не может быть пустым")
	}
	p.name = name
	return nil
}

func (p *Person) SetAge(age int) error {
	if age <= 0 {
		return fmt.Errorf("не может быть пустым")
	}
	p.age = age
	return nil
}

func NewPerson(name string, age int) (*Person, error) {
	p := &Person{}
	if err := p.SetName(name); err != nil {
		return nil, err
	}

	if err := p.SetAge(age); err != nil {
		return nil, err
	}

	return p, nil
}
