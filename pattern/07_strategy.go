package pattern

import "fmt"

/*
	Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern

Стратегия — это поведенческий паттерн проектирования, который определяет семейство схожих алгоритмов и помещает
каждый из них в собственный класс, после чего алгоритмы можно взаимозаменять прямо во время исполнения программы.

Плюсы:
- Горячая замена алгоритмов на лету.
- Изолирует код и данные алгоритмов от остальных классов.
- Уход от наследования к делегированию.
- Реализует принцип открытости/закрытости.

Минусы:
- Усложняет программу за счёт дополнительных классов.
- Клиент должен знать, в чём состоит разница между стратегиями, чтобы выбрать подходящую.
*/

// Operator интерфейс стратегии
type Operator interface {
	apply(int, int) int
}

type Operation struct {
	Operator Operator
}

func (o *Operation) Operate(lvalue, rvalue int) int {
	return o.Operator.apply(lvalue, rvalue)
}

// Addiction конкретная стратегия для суммирования
type Addiction struct {
}

func (a *Addiction) apply(lval, rval int) int {
	return lval + rval
}

// Multiple конкретная стратегия для умножения
type Multiple struct {
}

func (m *Multiple) apply(lval, rval int) int {
	return lval * rval
}

func StrategyPattern() {
	// Передаем в контекст умножение
	mult := Operation{&Multiple{}}
	fmt.Println(mult.Operate(10, 2))
	// Передаем в контекст суммирование
	add := Operation{&Addiction{}}
	fmt.Println(add.Operate(10, 2))

}
