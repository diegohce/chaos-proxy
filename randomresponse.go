package main

import (
	"math/rand"
	"time"
)

type randomError interface {
	kind() string
}

type random5xx struct {}

func (r random5xx) kind() string {
	return "5xx"
}
func (r random5xx) status() int {
	return 500 + rand.Intn(5)
}

func (r random5xx) Error() string {
	return r.kind()
}

type hangup struct{}

func (r hangup) kind() string {
	return "HUP"
}

type delay struct{}

func (r delay) kind() string {
	return "DELAY"
}

func (r delay) wait() {
	//d := time.Millisecond * time.Duration( rand.Intn(54000001) )
	d := time.Millisecond * time.Duration( rand.Intn(600001) )
	time.Sleep(d)
}

type noError struct {}

func (r noError) kind() string {
	return "NOERR"
}

var randomErrors []randomError = []randomError{
	random5xx{},
	hangup{},
	delay{},
	noError{},
	noError{},
	noError{},
	noError{},
	noError{},
	noError{},
	noError{},
	noError{},
	noError{},
	noError{},
	noError{},
	noError{},
	noError{},
	noError{},
}

func rollDices() randomError {
	return randomErrors[rand.Intn(len(randomErrors))]
}

