package pattern

import "fmt"

/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern

Посетитель — это поведенческий паттерн проектирования, который позволяет добавлять в программу новые операции,
не изменяя классы объектов, над которыми эти операции могут выполняться.

Применимость:
Когда вам нужно выполнить какую-то операцию над всеми элементами сложной структуры объектов, например, деревом.
Посетитель позволяет применять одну и ту же операцию к объектам различных классов.

Когда над объектами сложной структуры объектов надо выполнять некоторые не связанные между собой операции,
но вы не хотите «засорять» классы такими операциями.
Посетитель позволяет извлечь родственные операции из классов, составляющих структуру объектов,
поместив их в один класс-посетитель. Если структура объектов является общей для нескольких приложений,
то паттерн позволит в каждое приложение включить только нужные операции.

Когда новое поведение имеет смысл только для некоторых классов из существующей иерархии.
Посетитель позволяет определить поведение только для этих классов, оставив его пустым для всех остальных.

Плюсы:
- Упрощает добавление операций, работающих со сложными структурами объектов.
- Объединяет родственные операции в одном классе.
- Посетитель может накапливать состояние при обходе структуры элементов.

Минусы:
- Паттерн не оправдан, если иерархия элементов часто меняется.
- Может привести к нарушению инкапсуляции элементов.

*/

type Visitor interface {
	visitForSquare(*Square)
	visitForCircle(*Circle)
	visitForRectangle(*Rectangle)
}

// Shape общий интерфейс для всех фигур, с методами,
// которые мы добавим для дальнейшего взаимодействия с любым типом фигур
type Shape interface {
	getType() string
	accept(Visitor)
}

// Square реализация полей и методов интерфейса Shape у квадрата
type Square struct {
	side int
}

func (s *Square) accept(v Visitor) {
	v.visitForSquare(s)
}

func (s *Square) getType() string {
	return "Square"
}

// Circle реализация полей и методов интерфейса Shape у круга
type Circle struct {
	radius int
}

func (c *Circle) accept(v Visitor) {
	v.visitForCircle(c)
}

func (c *Circle) getType() string {
	return "Circle"
}

// Rectangle реализация полей и методов интерфейса Shape у прямоугольника
type Rectangle struct {
	a int
	b int
}

func (r *Rectangle) accept(v Visitor) {
	v.visitForRectangle(r)
}

func (r *Rectangle) getType() string {
	return "Rectangle"
}

// AreaCalculator Создаем конкретного посетителя и реализуем его интерфейс
type AreaCalculator struct {
	area int
}

func (a *AreaCalculator) visitForSquare(s *Square) {
	fmt.Println("Считаем площадь квадрата...")
}

func (a *AreaCalculator) visitForCircle(c *Circle) {
	fmt.Println("Считаем площадь круга...")
}

func (a *AreaCalculator) visitForRectangle(r *Rectangle) {
	fmt.Println("Считаем площадь прямоугольника...")
}

func VisitorPattern() {
	square := &Square{2}
	circle := &Circle{5}
	rectangle := &Rectangle{2, 4}

	areaCalculator := &AreaCalculator{}

	square.accept(areaCalculator)
	circle.accept(areaCalculator)
	rectangle.accept(areaCalculator)
}
