package pattern

import "fmt"

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern

Цепочка обязанностей — это поведенческий паттерн проектирования,
который позволяет передавать запросы последовательно по цепочке обработчиков.
Каждый последующий обработчик решает, может ли он обработать запрос сам и стоит ли передавать запрос дальше по цепи.

Плюсы:
- Уменьшает зависимость между клиентом и обработчиками.
- Реализует принцип единственной обязанности(Single responsibility).
- Реализует принцип открытости/закрытости(Open-closed).

Минусы:
- Запрос может остаться никем не обработанным.

Реализуем паттерн взяв в пример регистрацию пользователя
*/

// User Структура пользователя, который будет отправлять запрос
type User struct {
	loginCorrect bool
	phoneCorrect bool
	emailCorrect bool
	registerDone bool
}

// Website Единый интерфейс под звенья цепочки обработчиков запроса
type Website interface {
	execute(*User)
	setNext(Website)
}

// Login Конкретный обработчик корректности логина
type Login struct {
	next Website
}

func (l *Login) execute(u *User) {
	if u.loginCorrect {
		fmt.Println("User already exist!")
		l.next.execute(u)
		return
	}
	fmt.Println("Correct login")
	u.loginCorrect = true
	l.next.execute(u)
}

func (l *Login) setNext(next Website) {
	l.next = next
}

// Phone Конкретный обработчик номера телефона
type Phone struct {
	next Website
}

func (p *Phone) execute(u *User) {
	if u.loginCorrect {
		fmt.Println("Phone number already exist!")
		p.next.execute(u)
		return
	}
	fmt.Println("Correct phone number")
	u.phoneCorrect = true
	p.next.execute(u)
}

func (p *Phone) setNext(next Website) {
	p.next = next
}

// Email Конкретный обработчик корректности почты
type Email struct {
	next Website
}

func (e *Email) execute(u *User) {
	if u.loginCorrect {
		fmt.Println("Email address already exist!")
		e.next.execute(u)
		return
	}
	fmt.Println("Correct email address")
	u.emailCorrect = true
	e.next.execute(u)
}

func (e *Email) setNext(next Website) {
	e.next = next
}

// Register Конкретный обработчик регистрации пользователя
type Register struct {
	next Website
}

func (r *Register) execute(u *User) {
	if u.registerDone {
		fmt.Println("Register done")
	}
	fmt.Println("Account have been registered!")
}

func (r *Register) setNext(next Website) {
	r.next = next
}

func ChainOfResPattern() {
	register := &Register{}

	login := &Login{}
	login.setNext(register)

	phone := &Phone{}
	phone.setNext(login)

	email := &Email{}
	email.setNext(phone)

	user := &User{}
	email.execute(user)
}
