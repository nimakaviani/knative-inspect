package config

import (
	"log"
)

type logger struct{}

func (o *logger) Infof(s string, params ...interface{}) {
	log.Printf(s, params...)
}

func (o *logger) Errorf(s string, params ...interface{}) {
	log.Printf(s, params...)
}

func (o *logger) Fatalf(s string, params ...interface{}) {
	log.Fatalf(s, params)
}
