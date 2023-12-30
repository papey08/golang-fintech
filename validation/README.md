# validation

Данный модуль позволяет осуществить проверку полей структуры, отмеченных 
соответсвующим тегом, на валидность.

## Поддерживаемые теги

| Тег                      | Описание                                               | Поддерживаемые типы |
|--------------------------|--------------------------------------------------------|---------------------|
| **len:arg**              | Длина строки строго равна ***arg***                    | string              |
| **in:arg1, arg2,…,argn** | Число/строка равны одному из ***arg***                 | string, int         |
| **min:arg**              | Число/длина строки не меньше ***arg***                 | string, int         |
| **max:arg**              | Число/длина строки не больше ***arg***                 | string, int         |
| **lenInterval:min,max**  | Длина строки не больше ***max*** и не меньше ***min*** | string              |

## Пример кода

```go
package main

import (
	"fmt"

	v "github.com/papey08/golang-fintech/validation"
)

type Person struct {
	Name string `validate:"max:50"`
	Age  int    `validate:"min:0"`
}

func main() {
	p1 := Person{
		Name: "papey08",
		Age:  17,
	}
	p2 := Person{
		Name: "papey08",
		Age:  -1,
	}

	// вызываем функцию Validate из модуля
	if err := v.Validate(p1); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("p1 is a correct Person struct")
	}

	if err := v.Validate(p2); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("p2 is a correct Person struct")
	}
}
```

### Результат

```text
p1 is a correct Person struct
value of field is not validate
```
