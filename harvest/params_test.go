package harvest

import (
	"net/url"
	"reflect"
	"testing"
)

func TestParamsClone(t *testing.T) {
	params := make(Params)
	params.Set("foo", "bar")

	clonedParams := params.Clone()

	if !reflect.DeepEqual(params, clonedParams) {
		t.Logf("Expected cloned params to equal %+#v, got %+#v\n", params, clonedParams)
		t.Fail()
	}

	// Test for deep copy
	clonedParams.Set("foo", "qux")

	foo := params.Get("foo")
	clonedFoo := clonedParams.Get("foo")

	if foo == clonedFoo {
		t.Logf("Expected value in 'foo' to equal 'qux', got %q\n", clonedFoo)
		t.Fail()
	}
	if reflect.DeepEqual(params, clonedParams) {
		t.Logf("Expected cloned params not to equal %+#v, got %+#v\n", params, clonedParams)
		t.Fail()
	}

	// Test uninitialized value
	var uninitializedParams Params
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Expected Clone not to panic, got %T:%v\n", r, r)
			t.Fail()
		}
	}()
	clone := uninitializedParams.Clone()

	if clone == nil {
		t.Logf("Expected clone not to be nil\n")
		t.Fail()
	}
}

func TestParamsAdd(t *testing.T) {
	params := make(Params)
	params.Add("foo", "bar")

	expected := Params(url.Values{"foo": []string{"bar"}})

	if !reflect.DeepEqual(expected, params) {
		t.Logf("Expected params to equal \n%+#v\n\tgot\n%+#v\n", expected, params)
		t.Fail()
	}

	// Add another value for the same key
	params.Add("foo", "qux")

	expected = Params(url.Values{"foo": []string{"bar", "qux"}})

	if !reflect.DeepEqual(expected, params) {
		t.Logf("Expected params to equal \n%+#v\n\tgot\n%+#v\n", expected, params)
		t.Fail()
	}

	// Test uninitialized value
	var uninitialized Params
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Expected Add not to panic, got %T:%v\n", r, r)
			t.Fail()
		}
	}()
	uninitialized.Add("foo", "bar")
}

func TestParamsSet(t *testing.T) {
	params := make(Params)
	params.Set("foo", "bar")

	expected := Params(url.Values{"foo": []string{"bar"}})

	if !reflect.DeepEqual(expected, params) {
		t.Logf("Expected params to equal \n%+#v\n\tgot\n%+#v\n", expected, params)
		t.Fail()
	}

	// Set another value for the same key
	params.Set("foo", "qux")

	expected = Params(url.Values{"foo": []string{"qux"}})

	if !reflect.DeepEqual(expected, params) {
		t.Logf("Expected params to equal \n%+#v\n\tgot\n%+#v\n", expected, params)
		t.Fail()
	}

	// Test uninitialized value
	var uninitialized Params
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Expected Set not to panic, got %T:%v\n", r, r)
			t.Fail()
		}
	}()
	uninitialized.Set("foo", "bar")
}

func TestParamsGet(t *testing.T) {
	params := make(Params)
	params.Set("foo", "bar")

	expected := "bar"
	actual := params.Get("foo")

	if expected != actual {
		t.Logf("Expected fetched value to equal %q, got%q\n", expected, actual)
		t.Fail()
	}

	// Test uninitialized value
	var uninitialized Params
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Expected Get not to panic, got %T:%v\n", r, r)
			t.Fail()
		}
	}()
	uninitialized.Get("foo")
}

func TestParamsDel(t *testing.T) {
	params := make(Params)
	params.Set("foo", "bar")

	params.Del("foo")

	if len(params) != 0 {
		t.Logf("Expected params length to equal 0, got%d\n", len(params))
		t.Fail()
	}

	// Test uninitialized value
	var uninitialized Params
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Expected Del not to panic, got %T:%v\n", r, r)
			t.Fail()
		}
	}()
	uninitialized.Del("foo")
}

func TestParamsEncode(t *testing.T) {
	params := make(Params)
	params.Set("foo", "bar")

	expected := url.Values{"foo": []string{"bar"}}.Encode()

	actual := params.Encode()

	if expected != actual {
		t.Logf("Expected encoded params to equal %q, got %q\n", expected, actual)
		t.Fail()
	}

	// Test uninitialized value
	var uninitialized Params
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Expected Encode not to panic, got %T:%v\n", r, r)
			t.Fail()
		}
	}()
	uninitialized.Encode()
}

func TestParamsValues(t *testing.T) {
	params := make(Params)
	params.Set("foo", "bar")

	expected := url.Values{"foo": []string{"bar"}}

	actual := params.Values()

	if !reflect.DeepEqual(expected, actual) {
		t.Logf("Expected encoded params to equal %+#v, got %+#v\n", expected, actual)
		t.Fail()
	}

	// Test uninitialized value
	var uninitialized Params
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Expected Values not to panic, got %T:%v\n", r, r)
			t.Fail()
		}
	}()
	uninitialized.Values()
}

func TestParamsMerge(t *testing.T) {
	var tests = []struct {
		initial       Params
		valuesToMerge url.Values
		resulting     Params
	}{
		{
			Params{"foo": []string{"bar"}},
			url.Values{"bar": []string{"X"}},
			Params{"foo": []string{"bar"}, "bar": []string{"X"}},
		},
		{
			Params{"foo": []string{"bar"}},
			url.Values{"foo": []string{"X"}},
			Params{"foo": []string{"bar", "X"}},
		},
	}
	for _, test := range tests {
		actual := test.initial.Merge(test.valuesToMerge)

		if !reflect.DeepEqual(test.resulting, actual) {
			t.Logf("Expected encoded params to equal \n%+#v\n\tgot:\n%+#v\n", test.resulting, actual)
			t.Fail()
		}
	}

	// Test uninitialized value
	var uninitialized Params
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Expected Merge not to panic, got %T:%v\n", r, r)
			t.Fail()
		}
	}()
	uninitialized.Merge(url.Values{})
}

func TestParamsForTimeframe(t *testing.T) {
	params := make(Params)

	timeframe := Timeframe{
		StartDate: Date(2010, 01, 01, nil),
		EndDate:   Date(2012, 01, 01, nil),
	}

	params.ForTimeframe(timeframe)

	to := params.Get("to")
	from := params.Get("from")

	if to == "" || from == "" {
		t.Logf("Expected timeframe to get serialized in to and from, was not\n")
		t.Fail()
	}
}
