package pattern

import "fmt"

/*
	Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern

Строитель (англ. Builder) — порождающий шаблон проектирования предоставляет способ создания составного объекта.
Отделяет конструирование сложного объекта от его представления так, что в результате одного и того
же процесса конструирования могут получаться разные представления.

Проблема:
    Инициализация очень сложного, большого объекта со множеством параметров инициализации.
    Использовать один конструктор с множеством параметров - плохо (телескопический конструктор - анти-паттерн)

Решение:
    Паттерн Строитель предлагает вынести конструирование объекта за пределы его собственного класса, поручив это дело
    отдельным объектам, называемым строителями.

Плюсы:
    - Позволяет создавать продукты пошагово.
    - Позволяет использовать один и тот же код для создания различных продуктов.
    - Изолирует сложный код сборки продукта от его основной бизнес-логики.

Минусы:
    - Усложняет код программы из-за введения дополнительных классов.
    - Клиент будет привязан к конкретным классам строителей,
	  так как в интерфейсе директора может не быть метода получения результата.
*/

// IBuilder интерфейс строителя
type IBuilder interface {
	setWindowType()
	setRoofType()
	setDoorType()
	getBuilding() Building
}

// Функция возвращающая интерфейс конкретного строителя
func getBuilder(builderType string) (IBuilder, error) {
	switch {
	case builderType == "Wooden":
		return newWoodenBuilder(), nil
	case builderType == "Brick":
		return newBrickBuilder(), nil
	default:
		return nil, fmt.Errorf("wrong type of Builder passed")
	}
}

type Building struct {
	windowType string
	roofType   string
	doorType   string
}

// WoodenBuilder конкретный строитель для деревянных домов
type WoodenBuilder struct {
	windowType string
	roofType   string
	doorType   string
}

// newWoodenBuilder Создаем строителя, который строит только из дерева
func newWoodenBuilder() *WoodenBuilder {
	return &WoodenBuilder{}
}

// Методы строителя, который строит из дерева
func (w *WoodenBuilder) setWindowType() {
	w.windowType = "Wooden windows"
}

func (w *WoodenBuilder) setRoofType() {
	w.roofType = "Thatched roof"
}

func (w *WoodenBuilder) setDoorType() {
	w.doorType = "Wooden door"
}

func (w *WoodenBuilder) getBuilding() Building {
	return Building{
		w.windowType,
		w.roofType,
		w.doorType,
	}
}

// BrickBuilder конкретный строитель для кирпичных домов
type BrickBuilder struct {
	windowType string
	roofType   string
	doorType   string
}

// newBrickBuilder Создаем строителя, который строит только из кирпича
func newBrickBuilder() *BrickBuilder {
	return &BrickBuilder{}
}

// Методы строителя, который строит из кирпича
func (b *BrickBuilder) setWindowType() {
	b.windowType = "Brick windows"
}

func (b *BrickBuilder) setRoofType() {
	b.roofType = "Tile roof"
}

func (b *BrickBuilder) setDoorType() {
	b.doorType = "Steel door"
}

func (b *BrickBuilder) getBuilding() Building {
	return Building{
		b.windowType,
		b.roofType,
		b.doorType,
	}
}

// Director структура директора, для управления строителями
type Director struct {
	builder IBuilder
}

// Создание директора и назначение строителя
func newDirector(b IBuilder) *Director {
	return &Director{
		builder: b,
	}
}

// Назначение строителя
func (d *Director) setBuilder(b IBuilder) {
	d.builder = b
}

// Строительство здания
func (d *Director) createBuilding() Building {
	d.builder.setWindowType()
	d.builder.setRoofType()
	d.builder.setDoorType()
	return d.builder.getBuilding()
}

func ToDoBuilderPattern() {
	woodenBuilder, err := getBuilder("Wooden")
	if err != nil {
		panic(err)
	}

	brickBuilder, err := getBuilder("Brick")
	if err != nil {
		panic(err)
	}

	director := newDirector(woodenBuilder)
	woodenBuilding := director.createBuilding()

	fmt.Printf("Wooden House Window Type: %s\n", woodenBuilding.windowType)
	fmt.Printf("Wooden House Roof Type: %s\n", woodenBuilding.roofType)
	fmt.Printf("Wooden House Door Floor: %s\n", woodenBuilding.doorType)

	director.setBuilder(brickBuilder)
	brickBuilding := director.createBuilding()

	fmt.Printf("Brick House Window Type: %s\n", brickBuilding.windowType)
	fmt.Printf("Brick House Roof Type: %s\n", brickBuilding.roofType)
	fmt.Printf("Brick House Door Floor: %s\n", brickBuilding.doorType)
}
