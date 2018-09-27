package promise

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	p := New(func(resolve, reject func(interface{})) {
	})
	assert.IsType(t, &Promise{}, p)
}

func TestResolve(t *testing.T) {
	p := Resolve("test").Await()
	assert.Equal(t, "test", p)
}

func TestReject(t *testing.T) {
	defer func() {
		err := recover()
		assert.NotEqual(t, nil, err)
	}()

	Reject("test").Await()
}

func TestAll(t *testing.T) {

	defer func() {
		err := recover()
		assert.Equal(t, "response p4", err)
	}()
	p1 := New(func(resolve, reject func(interface{})) {
		resolve("response p1")
	})
	p2 := New(func(resolve, reject func(interface{})) {
		resolve("response p2")
	})

	p3 := Resolve("response p3")

	p4 := Reject("response p4")

	response := All(p1, p2, p3).Await()

	assert.Equal(t, []interface{}{
		"response p1",
		"response p2",
		"response p3",
	}, response)

	All(p1, p2, p3, p4).Await()
}

func TestRace(t *testing.T) {

	p1 := New(func(resolve, reject func(interface{})) {
		time.Sleep(2 * time.Second)
		resolve("response p1")
	})

	p2 := New(func(resolve, reject func(interface{})) {
		time.Sleep(1 * time.Second)
		resolve("response p2")
	})

	assert.Equal(t, "response p2", Race(p1, p2).Await())
}

func TestFromFunction(t *testing.T) {
	testFunc := func() interface{} {
		time.Sleep(2 * time.Second)
		return "hi"
	}

	p1 := FromFunction(testFunc).Await()

	assert.Equal(t, "hi", p1)
}

func TestPromise_Then(t *testing.T) {
	p := Resolve(2).Then(func(value interface{}) interface{} {
		return value.(int) + 4
	})

	assert.Equal(t, 6, p.Await())

	p = p.Then(func(value interface{}) interface{} {
		return "This message has " + strconv.Itoa(value.(int)) + " words."
	})

	assert.Equal(t, "This message has 6 words.", p.Await())
}

func TestPromise_Catch(t *testing.T) {
	p := Reject("kaboom").Catch(func(reason interface{}) interface{} {
		return reason.(string) + " blocked by catch"
	}).Await()

	assert.Equal(t, "kaboom blocked by catch", p)
}

func TestPromise_Await(t *testing.T) {
	p := Resolve("such")
	assert.IsType(t, &Promise{}, p)

	assert.IsType(t, "", p.Await())
	assert.Equal(t, "such", p.Await())
}