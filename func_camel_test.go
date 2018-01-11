package tmpl

import "testing"

func TestLowerCamel(t *testing.T) {
	ok := func(in, exp string) {
		act := funcLowerCamel(in)
		if act != exp {
			t.Errorf("not match: in=%q exp=%q act=%q", in, exp, act)
		}
	}
	ok("foo", "foo")
	ok("Foo", "foo")
	ok("fooBar", "fooBar")
	ok("FooBar", "fooBar")
	ok("fooBarBaz", "fooBarBaz")
	ok("FooBarBaz", "fooBarBaz")
}
