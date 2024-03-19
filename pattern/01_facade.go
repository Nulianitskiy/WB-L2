package pattern

import "strings"

/*
	Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern

Проблема:
Минимизировать зависимость подсистем некоторой сложной системы и обмен информацией между ними.

Решение:
Фасад — простой интерфейс для работы со сложным фреймворком. Фасад не имеет всей функциональности фреймворка,
но зато скрывает его сложность от клиентов.

Плюсы:
	- Изолирует клиентов от компонентов сложной подсистемы.
Минусы:
	- Фасад рискует стать божественным объектом, привязанным ко всем классам программы.
	(объект, который хранит в себе «слишком много» или делает «слишком много».)

Примером реализации фасада будет человек, которому доступно множество функций, но мы реализуем только базовые:
Вырастить дом, построить дерево и посадить сына :)
*/

// Man Наш фасад
type Man struct {
	tree  *Tree
	house *House
	son   *Son
}

func (man Man) ToDoFacade() string {
	result := []string{
		man.house.ToGrow(),
		man.tree.ToBuild(),
		man.son.ToPlant(),
	}
	return strings.Join(result, "\n")
}

// Tree имплементация дерева
type Tree struct {
}

func (t *Tree) ToBuild() string {
	return "Build a Tree"
}

// House имплементация дома
type House struct {
}

func (h *House) ToGrow() string {
	return "Grow House"
}

// Son имплементация сына
type Son struct {
}

func (s *Son) ToPlant() string {
	return "To Plant a Son :)"
}
