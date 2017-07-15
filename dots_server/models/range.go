package models

type Range interface {
	Start() interface{}
	End() interface{}
	IncludeRange(Range) bool
	Includes(interface{}) bool
}
