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
	return 500 + rand.Intn(12)
}

func (r random5xx) Error() string {
	return r.kind()
}

type hangup struct{}

func (r hangup) kind() string {
	return "HUP"
}

type delay struct {}

func (r delay) kind() string {
	return "DELAY"
}

func (r delay) wait(url string) {
	d := time.Millisecond * time.Duration( rand.Intn(chaosConfig.MaxTimeout + 1) )
	log.Info().Println("Will delay", d, "for", url)
	time.Sleep(d)
}

type noError struct {}

func (r noError) kind() string {
	return "NOERR"
}

var randomErrors []randomError = []randomError{
	noError{},
	noError{},
	noError{},
	noError{},
	noError{},
	hangup{},
	noError{},
	noError{},
	noError{},
	noError{},
	random5xx{},
	noError{},
	noError{},
	noError{},
	noError{},
	noError{},
	random5xx{},
	noError{},
	noError{},
	noError{},
	noError{},
	delay{},
	noError{},
	noError{},
	noError{},
	noError{},
	random5xx{},
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

