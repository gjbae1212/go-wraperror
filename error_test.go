package wraperror

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	assert := assert.New(t)

	test1Err := fmt.Errorf("[err] test")

	tests := map[string]struct {
		err    error
		output error
	}{
		"success-1": {err: nil, output: &wrapError{}},
		"success-2": {err: test1Err, output: &wrapError{current: test1Err}},
		"success-3": {err: &wrapError{current: test1Err}, output: &wrapError{current: test1Err}},
	}

	for _, t := range tests {
		wrapErr := Error(t.err)
		switch t.err.(type) {
		case *wrapError:
			assert.Equal(t.err.(*wrapError).current, wrapErr.current)
			assert.Equal(t.err.(*wrapError).child, wrapErr.child)
		default:
			assert.Equal(t.err, wrapErr.current)
			assert.Equal(nil, wrapErr.child)
		}

	}
}

func TestWrapError_Wrap(t *testing.T) {
	assert := assert.New(t)

	test1Err := fmt.Errorf("[err] test 1")
	test2Err := fmt.Errorf("[err] test 2")

	emptyErr := Error(nil)
	existErr := Error(test1Err)

	tests := map[string]struct {
		err    *wrapError
		input  error
		output *wrapError
	}{
		"empty": {err: emptyErr, input: test1Err, output: &wrapError{current: test1Err, child: emptyErr}},
		"exist": {err: existErr, input: test2Err, output: &wrapError{current: test2Err, child: existErr}},
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
		err   *wrapError
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

	tests := map[string]struct {
		err    *wrapError
		output string
	}{
		"success-1": {err: emptyErr, output: ""},
		"success-2": {err: existErr, output: "[err] test"},
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
		err    *wrapError
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
		err    *wrapError
		target error
		ok     bool
	}{
		"success 1": {err: chainExistErr, target: testExistOnlyErr, ok: true},
		"success 2": {err: chainExistErr, target: testDefaultErr, ok: true},
		"success 3": {err: chainEmptyErr, target: testDefaultErr, ok: true},
		"fail 1":    {err: &wrapError{}, target: testDefaultErr, ok: false},
		"fail 2":    {err: &wrapError{}, target: nil, ok: false},
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

	must := &wrapError{}
	var mustnot *os.PathError
	tests := map[string]struct {
		err    *wrapError
		target interface{}
		ok     bool
	}{
		"must success 1": {err: &wrapError{}, target: must, ok: true},
		"must success 2": {err: chainEmptyErr, target: must, ok: true},
		"must success 3": {err: chainExistErr, target: must, ok: true},
		"success 1":      {err: chainExistErr, target: testDefaultErr, ok: true},
		"success 2":      {err: chainExistErr, target: testExistOnlyErr, ok: true},
		"success 3":      {err: chainEmptyErr, target: testDefaultErr, ok: true},
		"success 4":      {err: chainEmptyErr, target: testExistOnlyErr, ok: true},
		"success 5":      {err: &wrapError{}, target: testExistOnlyErr, ok: true},
		"success 6":      {err: &wrapError{}, target: testDefaultErr, ok: true},
		"fail 1":         {err: &wrapError{}, target: mustnot, ok: false},
		"fail 2":         {err: chainEmptyErr, target: mustnot, ok: false},
		"fail 3":         {err: chainExistErr, target: mustnot, ok: false},
	}

	for _, t := range tests {
		switch a := t.target.(type) {
		case *wrapError:
			assert.Equal(t.ok, errors.As(t.err, &a))
		case *os.PathError:
			assert.Equal(t.ok, errors.As(t.err, &a))
		case error:
			assert.Equal(t.ok, errors.As(t.err, &a))
		}
	}
}
