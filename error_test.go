package wraperror

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapError_Current(t *testing.T) {
	assert := assert.New(t)

	testErr := fmt.Errorf("[err] test 1")

	emptyErr := Error(nil)
	existErr := Error(testErr)

	tests := map[string]struct {
		err    *WrapError
		output error
	}{
		"empty": {err: emptyErr, output: nil},
		"exist": {err: existErr, output: testErr},
	}

	for _, t := range tests {
		current := t.err.Current()
		assert.Equal(t.output, current)
	}
}

func TestWrapError_Child(t *testing.T) {
	assert := assert.New(t)

	test1Err := fmt.Errorf("[err] test 1")
	test2Err := fmt.Errorf("[err] test 2")

	emptyErr := Error(nil)
	existErr := Error(test1Err)

	tests := map[string]struct {
		err    *WrapError
		output error
	}{
		"empty": {err: emptyErr, output: nil},
		"exist": {err: existErr.Wrap(test2Err), output: existErr},
	}

	for _, t := range tests {
		child := t.err.Child()
		assert.Equal(t.output, child)
	}
}

func TestWrapError_Wrap(t *testing.T) {
	assert := assert.New(t)

	test1Err := fmt.Errorf("[err] test 1")
	test2Err := fmt.Errorf("[err] test 2")

	emptyErr := Error(nil)
	existErr := Error(test1Err)

	tests := map[string]struct {
		err    *WrapError
		input  error
		output *WrapError
	}{
		"empty": {err: emptyErr, input: test1Err, output: &WrapError{current: test1Err, child: emptyErr}},
		"exist": {err: existErr, input: test2Err, output: &WrapError{current: test2Err, child: existErr}},
	}

	for _, t := range tests {
		wrapErr := t.err.Wrap(t.input)
		assert.Equal(t.output.current, wrapErr.current)
		assert.Equal(t.output.child, wrapErr.child)
	}
}

func TestWrapError_Flatten(t *testing.T) {
	assert := assert.New(t)

	testDefaultErr := fmt.Errorf("[err] default test")
	testExistOnlyErr := fmt.Errorf("[err] sample test ")

	chainExistErr := Error(testExistOnlyErr)
	chainEmptyErr := Error(nil)

	chainEmptyErr = chainEmptyErr.Wrap(testDefaultErr)
	chainExistErr = chainExistErr.Wrap(testDefaultErr)
	for i := 0; i < 100; i++ {
		chainEmptyErr = chainEmptyErr.Wrap(fmt.Errorf("[err] test %d %w", i, &os.PathError{Err: fmt.Errorf("%d", i)}))
		chainExistErr = chainExistErr.Wrap(fmt.Errorf("[err] test %d", i))
	}
	mixErr := chainExistErr.Wrap(chainEmptyErr)

	tests := map[string]struct {
		err   *WrapError
		count int
	}{
		"chain empty": {err: chainEmptyErr, count: 301}, // 1 + 100 + Os.PathError(100) + Os.PathError.Err(100)
		"chain exist": {err: chainExistErr, count: 102},
		"mix":         {err: mixErr, count: 403},
	}

	for _, t := range tests {
		assert.Len(t.err.Flatten(), t.count)
	}
}

func TestWrapError_Error(t *testing.T) {
	assert := assert.New(t)

	emptyErr := Error(nil)
	existErr := Error(fmt.Errorf("[err] test"))
	wrappedErr := Error(fmt.Errorf("foo")).Wrap(existErr)

	tests := map[string]struct {
		err    *WrapError
		output string
	}{
		"success-1": {err: emptyErr, output: ""},
		"success-2": {err: existErr, output: "[err] test"},
		"success-3": {err: wrappedErr, output: "[err] test: foo"},
	}

	for _, t := range tests {
		assert.Equal(t.output, t.err.Error())
	}
}

func TestWrapError_Unwrap(t *testing.T) {
	assert := assert.New(t)

	test1Err := fmt.Errorf("[err] test 1")
	test2Err := fmt.Errorf("[err] test 2")

	emptyErr := Error(nil)
	existErr := Error(test1Err)

	chainEmptyErr := emptyErr.Wrap(nil)
	chainExistErr := existErr.Wrap(test2Err)

	tests := map[string]struct {
		err    *WrapError
		output error
	}{
		"empty":       {err: emptyErr, output: nil},
		"exist":       {err: existErr, output: nil},
		"chain empty": {err: chainEmptyErr, output: emptyErr},
		"chain exist": {err: chainExistErr, output: existErr},
	}

	for _, t := range tests {
		assert.Equal(t.output, t.err.Unwrap())
	}
}

func TestWrapError_Is(t *testing.T) {
	assert := assert.New(t)

	testDefaultErr := fmt.Errorf("[err] default test")
	testExistOnlyErr := fmt.Errorf("[err] sample test ")

	chainExistErr := Error(testExistOnlyErr)
	chainEmptyErr := Error(nil)

	chainEmptyErr = chainEmptyErr.Wrap(testDefaultErr)
	chainExistErr = chainExistErr.Wrap(testDefaultErr)
	for i := 0; i < 100; i++ {
		chainEmptyErr = chainEmptyErr.Wrap(fmt.Errorf("[err] test %d", i))
		chainExistErr = chainExistErr.Wrap(fmt.Errorf("[err] test %d", i))
	}

	tests := map[string]struct {
		err    *WrapError
		target error
		ok     bool
	}{
		"success 1": {err: chainExistErr, target: testExistOnlyErr, ok: true},
		"success 2": {err: chainExistErr, target: testDefaultErr, ok: true},
		"success 3": {err: chainEmptyErr, target: testDefaultErr, ok: true},
		"fail 1":    {err: &WrapError{}, target: testDefaultErr, ok: false},
		"fail 2":    {err: &WrapError{}, target: nil, ok: false},
		"fail 3":    {err: chainEmptyErr, target: testExistOnlyErr, ok: false},
	}

	for _, t := range tests {
		assert.Equal(t.ok, errors.Is(t.err, t.target))
	}
}

func TestWrapError_As(t *testing.T) {
	assert := assert.New(t)

	testDefaultErr := errors.New("[err] default test")
	testExistOnlyErr := errors.New("[err] sample test ")

	chainExistErr := Error(testExistOnlyErr)
	chainEmptyErr := Error(nil)

	chainEmptyErr = chainEmptyErr.Wrap(testDefaultErr)
	chainExistErr = chainExistErr.Wrap(testDefaultErr)
	for i := 0; i < 100; i++ {
		chainEmptyErr = chainEmptyErr.Wrap(fmt.Errorf("[err] test %d", i))
		chainExistErr = chainExistErr.Wrap(fmt.Errorf("[err] test %d", i))
	}

	must := &WrapError{}
	var mustnot *os.PathError
	tests := map[string]struct {
		err    *WrapError
		target interface{}
		ok     bool
	}{
		"must success 1": {err: &WrapError{}, target: must, ok: true},
		"must success 2": {err: chainEmptyErr, target: must, ok: true},
		"must success 3": {err: chainExistErr, target: must, ok: true},
		"success 1":      {err: chainExistErr, target: testDefaultErr, ok: true},
		"success 2":      {err: chainExistErr, target: testExistOnlyErr, ok: true},
		"success 3":      {err: chainEmptyErr, target: testDefaultErr, ok: true},
		"success 4":      {err: chainEmptyErr, target: testExistOnlyErr, ok: true},
		"success 5":      {err: &WrapError{}, target: testExistOnlyErr, ok: true},
		"success 6":      {err: &WrapError{}, target: testDefaultErr, ok: true},
		"fail 1":         {err: &WrapError{}, target: mustnot, ok: false},
		"fail 2":         {err: chainEmptyErr, target: mustnot, ok: false},
		"fail 3":         {err: chainExistErr, target: mustnot, ok: false},
	}

	for _, t := range tests {
		switch a := t.target.(type) {
		case *WrapError:
			assert.Equal(t.ok, errors.As(t.err, &a))
		case *os.PathError:
			assert.Equal(t.ok, errors.As(t.err, &a))
		case error:
			assert.Equal(t.ok, errors.As(t.err, &a))
		}
	}
}

func TestError(t *testing.T) {
	assert := assert.New(t)

	test1Err := fmt.Errorf("[err] test")

	tests := map[string]struct {
		err    error
		output error
	}{
		"success-1": {err: nil, output: &WrapError{}},
		"success-2": {err: test1Err, output: &WrapError{current: test1Err}},
		"success-3": {err: &WrapError{current: test1Err}, output: &WrapError{current: test1Err}},
	}

	for _, t := range tests {
		wrapErr := Error(t.err)
		switch t.err.(type) {
		case *WrapError:
			assert.Equal(t.err.(*WrapError).Current(), wrapErr.current)
			assert.Equal(t.err.(*WrapError).Child(), wrapErr.child)
		default:
			assert.Equal(t.err, wrapErr.current)
			assert.Equal(nil, wrapErr.child)
		}

	}
}

func TestFromError(t *testing.T) {
	assert := assert.New(t)

	test1Err := fmt.Errorf("[err] test")
	s1 := Error(test1Err)

	tests := map[string]struct {
		err    error
		output error
		ok     bool
	}{
		"fail":    {err: test1Err, output: nil, ok: false},
		"success": {err: s1, output: s1, ok: true},
	}

	for _, t := range tests {
		wrapErr, ok := FromError(t.err)
		assert.Equal(t.ok, ok)
		if t.ok {
			assert.Equal(wrapErr, t.output)
		}
	}
}
