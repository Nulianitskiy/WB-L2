package pattern

import "fmt"

/*
	Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/State_pattern

Состояние — это поведенческий паттерн проектирования, который позволяет объектам менять поведение
в зависимости от своего состояния. Извне создаётся впечатление, что изменился класс объекта.

Плюсы:
- Избавляет от множества больших условных операторов машины состояний.
- Концентрирует в одном месте код, связанный с определённым состоянием.
- Упрощает код контекста.

Минусы:
- Может неоправданно усложнить код, если состояний мало и они редко меняются.
*/

// Общий интерфейс партия Китай для различных состояний

type XiJinping interface {
	alert() string
}

// Реализует оповещение от состояния партия Китай миска рис

type XiJinpingReaction struct {
	state XiJinping
}

func (x *XiJinpingReaction) alert() string {
	return x.state.alert()
}

// Устанавливаем состояние партия
func (x *XiJinpingReaction) setState(state XiJinping) {
	x.state = state
}

func newXiJinpingState() *XiJinpingReaction {
	return &XiJinpingReaction{state: &ExaltState{}}
}

// ExaltState реализует состояние превознесения тебя партия Китай
type ExaltState struct {
}

func (e *ExaltState) alert() string {
	return "Партия миска рис гордится тобой, + 15 SOCIAL CREDIT"
}

// HateState реализует состояние ненависти к тебе партия Китай
type HateState struct {
}

func (h *HateState) alert() string {
	return "Партия считать ты злой бургер враг, готовься к нефритовому стержню, - 30000000 SOCIAL CREDIT Удар!"
}

func StatePattern() {
	xi := newXiJinpingState()
	fmt.Println(xi.alert())

	xi.setState(&HateState{})
	fmt.Println(xi.alert())
}
