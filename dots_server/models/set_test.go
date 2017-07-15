package models_test

import (
	"reflect"
	"testing"

	"github.com/nttdots/go-dots/dots_server/models"
)

func TestNewSetInt(t *testing.T) {
	setInt := models.NewSetInt()

	setInt.Append(1)

	var expects interface{}

	expects = []int{1}

	if !reflect.DeepEqual(setInt.List(), expects) {
		t.Errorf("setInt.List() got %s, want %s", setInt.List(), expects)
	}

	expects = true

	if !reflect.DeepEqual(setInt.Include(1), expects) {
		t.Errorf("setInt.List() got %s, want %s", setInt.List(), expects)
	}

	expects = false

	if !reflect.DeepEqual(setInt.Include(2), expects) {
		t.Errorf("setInt.List() got %s, want %s", setInt.List(), expects)
	}

	setInt.Append(2)

	expects = []int{1, 2}

	if !reflect.DeepEqual(setInt.List(), expects) {
		t.Errorf("setInt.List() got %s, want %s", setInt.List(), expects)
	}

	expects = true

	if !reflect.DeepEqual(setInt.Include(2), expects) {
		t.Errorf("setInt.List() got %s, want %s", setInt.List(), expects)
	}

	expects = true

	if !reflect.DeepEqual(setInt.Include(2), expects) {
		t.Errorf("setInt.List() got %s, want %s", setInt.List(), expects)
	}

	expects = true

	if !reflect.DeepEqual(setInt.Include(2), expects) {
		t.Errorf("setInt.List() got %s, want %s", setInt.List(), expects)
	}

	expects = false

	setInt.Delete(1)

	if !reflect.DeepEqual(setInt.Include(1), expects) {
		t.Errorf("setInt.List() got %s, want %s", setInt.List(), expects)
	}

	expects = true

	setInt.AddList([]int{3, 4, 5})

	if !reflect.DeepEqual(setInt.Include(3), expects) {
		t.Errorf("setInt.List() got %s, want %s", setInt.List(), expects)
	}

	expects = true

	if !reflect.DeepEqual(setInt.Include(4), expects) {
		t.Errorf("setInt.List() got %s, want %s", setInt.List(), expects)
	}

	expects = true

	if !reflect.DeepEqual(setInt.Include(5), expects) {
		t.Errorf("setInt.List() got %s, want %s", setInt.List(), expects)
	}

}

func TestNewSetString(t *testing.T) {
	setString := models.NewSetString()

	setString.Append("abc")

	var expects interface{}

	expects = []string{"abc"}

	if !reflect.DeepEqual(setString.List(), expects) {
		t.Errorf("setString.List() got %s, want %s", setString.List(), expects)
	}

	expects = true

	if !reflect.DeepEqual(setString.Include("abc"), expects) {
		t.Errorf("setString.List() got %s, want %s", setString.List(), expects)
	}

	expects = false

	if !reflect.DeepEqual(setString.Include("def"), expects) {
		t.Errorf("setString.List() got %s, want %s", setString.List(), expects)
	}

	setString.Append("def")

	expects = []string{"abc", "def"}

	if !reflect.DeepEqual(setString.List(), expects) {
		t.Errorf("setString.List() got %s, want %s", setString.List(), expects)
		return
	}

	expects = true

	if !reflect.DeepEqual(setString.Include("def"), expects) {
		t.Errorf("setString.List() got %s, want %s", setString.List(), expects)
	}

	expects = true

	if !reflect.DeepEqual(setString.Include("def"), expects) {
		t.Errorf("setString.List() got %s, want %s", setString.List(), expects)
	}

	expects = true

	if !reflect.DeepEqual(setString.Include("def"), expects) {
		t.Errorf("setString.List() got %s, want %s", setString.List(), expects)
	}

	expects = false

	setString.Delete("abc")

	if !reflect.DeepEqual(setString.Include("abc"), expects) {
		t.Errorf("setString.List() got %s, want %s", setString.List(), expects)
	}

	expects = true

	setString.AddList([]string{"g", "h", "i"})

	if !reflect.DeepEqual(setString.Include("g"), expects) {
		t.Errorf("setString.List() got %s, want %s", setString.List(), expects)
	}

	expects = true

	if !reflect.DeepEqual(setString.Include("h"), expects) {
		t.Errorf("setString.List() got %s, want %s", setString.List(), expects)
	}

	expects = true

	if !reflect.DeepEqual(setString.Include("i"), expects) {
		t.Errorf("setString.List() got %s, want %s", setString.List(), expects)
	}

}
