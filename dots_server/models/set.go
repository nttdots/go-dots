package models

import "github.com/nttdots/go-dots/dots_server/db_models"

type SetString map[string]struct{}

func (s SetString) Append(value string) {
	s[value] = struct{}{}
}

func (s SetString) List() []string {
	l := make([]string, 0)
	for ss, _ := range s {
		l = append(l, ss)
	}
	return l
}

func (s SetString) AddList(list []string) {

	for _, value := range list {
		s.Append(value)
	}
}

func (s SetString) Include(value string) bool {

	_, ok := s[value]
	return ok
}

func (s SetString) Delete(value string) {

	if s.Include(value) {
		delete(s, value)
	}
}

func (s SetString) ToInterfaceList() []interface {} {
	var array = make([]interface{}, 0)
	for _, v := range s.List() {
		array = append(array, v)
	}
	return array
}

func (s SetString) FromParameterValue(array []db_models.ParameterValue) {
	for _, v := range array {
		s.Append(v.StringValue)
	}
	return
}

func NewSetString() SetString {
	return make(SetString)
}

type SetInt map[int]struct{}

func (s SetInt) Append(value int) {
	s[value] = struct{}{}
}

func (s SetInt) List() []int {
	l := make([]int, 0)
	for ss, _ := range s {
		l = append(l, ss)
	}
	return l
}

func (s SetInt) AddList(list []int) {
	for _, value := range list {
		s.Append(value)
	}
}

func (s SetInt) Include(value int) bool {

	_, ok := s[value]
	return ok
}

func (s SetInt) Delete(value int) {

	if s.Include(value) {
		delete(s, value)
	}
}

func (s SetInt) FromParameterValue(array []db_models.ParameterValue) {
	for _, v := range array {
		s.Append(v.IntValue)
	}
	return
}

func (s SetInt) ToInterfaceList() []interface {} {
	var array = make([]interface{}, 0)
	for _, v := range s.List() {
		array = append(array, v)
	}
	return array
}

func NewSetInt() SetInt {
	return make(SetInt)
}

type Set interface {
	ToInterfaceList()	[]interface{}
	FromParameterValue([]db_models.ParameterValue)
}
