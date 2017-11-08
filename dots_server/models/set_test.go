package models_test

import (
	"reflect"
	"testing"

	"github.com/nttdots/go-dots/dots_server/models"
)

func assertDeepEqual(t *testing.T, msg string, actual interface{}, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Error("%s got %s, want %s", msg, actual, expected)
	}
}

func TestNewSetInt(t *testing.T) {
	setInt := models.NewSetInt()

	setInt.Append(1)

	assertDeepEqual(t, "setInt.List()", setInt.List(), []int{1})
	assertDeepEqual(t, "setInt.Include(1)", setInt.Include(1), true)
	assertDeepEqual(t, "setInt.Include(2)", setInt.Include(2), false)

	setInt.Append(2)

	{
		actual := setInt.List()
		expected1 := []int{1, 2}
		expected2 := []int{2, 1}

		if !reflect.DeepEqual(actual, expected1) && !reflect.DeepEqual(actual, expected2) {
			t.Errorf("setInt.List() got %s, want %s or %s", actual, expected1, expected2)
		}
	}

	assertDeepEqual(t, "setInt.Include(2)", setInt.Include(2), true)

	setInt.Delete(1)

	assertDeepEqual(t, "setInt.Include(1)", setInt.Include(1), false)

	setInt.AddList([]int{3, 4, 5})

	assertDeepEqual(t, "setInt.Include(3)", setInt.Include(3), true)
	assertDeepEqual(t, "setInt.Include(4)", setInt.Include(4), true)
	assertDeepEqual(t, "setInt.Include(5)", setInt.Include(5), true)
}

func TestNewSetString(t *testing.T) {
	setString := models.NewSetString()

	setString.Append("abc")

	assertDeepEqual(t, "setString.List()", setString.List(), []string{"abc"})
	assertDeepEqual(t, "setString.Include(\"abc\")", setString.Include("abc"), true)
	assertDeepEqual(t, "setString.Include(\"def\")", setString.Include("def"), false)

	setString.Append("def")

	{
		actual := setString.List()
		expected1 := []string{"abc", "def"}
		expected2 := []string{"def", "abc"}

		if !reflect.DeepEqual(actual, expected1) && !reflect.DeepEqual(actual, expected2) {
			t.Errorf("setString.List() got %s, want %s or %s", actual, expected1, expected2)
		}
	}

	assertDeepEqual(t, "setString.Include(\"def\")", setString.Include("def"), true)

	setString.Delete("abc")

	assertDeepEqual(t, "setString.Include(\"abc\")", setString.Include("abc"), false)

	setString.AddList([]string{"g", "h", "i"})

	assertDeepEqual(t, "setString.Include(\"g\")", setString.Include("g"), true)
	assertDeepEqual(t, "setString.Include(\"h\")", setString.Include("h"), true)
	assertDeepEqual(t, "setString.Include(\"i\")", setString.Include("i"), true)
}
