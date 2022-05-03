package server

import (
	"bufio"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func handlerHelloWorld(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetStatusCode(200)
	ctx.Response.SetBodyString("Hello, World!")
}

func TestNew(t *testing.T) {
	s, err := New(12345, handlerHelloWorld)
	assert.NoError(t, err)
	go s.Run()

	resp, err := http.Get("http://localhost:12345")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	defer resp.Body.Close()

	r := bufio.NewReader(resp.Body)
	body := make([]byte, 13)
	n, err := r.Read(body)
	assert.NoError(t, err)
	assert.Greater(t, n, 0)

	assert.Equal(t, "Hello, World!", string(body))
	err = s.Stop()
	assert.NoError(t, err)

	resp, err = http.Get("http://localhost:12345")
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestNewGCP(t *testing.T) {
	os.Setenv("PORT", "12346")
	s, err := NewGCP(handlerHelloWorld)
	assert.NoError(t, err)
	go s.Run()

	resp, err := http.Get("http://localhost:12346")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	defer resp.Body.Close()

	s.Stop()
}

func TestNewGCPInvalidPort(t *testing.T) {
	os.Setenv("PORT", "123456789")
	_, err := NewGCP(handlerHelloWorld)
	assert.Error(t, err)
}

func TestNewGCPStringPort(t *testing.T) {
	os.Setenv("PORT", "thisisnotaport")
	_, err := NewGCP(handlerHelloWorld)
	assert.Error(t, err)
}

func TestNewGCPNoPort(t *testing.T) {
	os.Unsetenv("PORT")
	_, err := NewGCP(handlerHelloWorld)
	assert.Error(t, err)
}

func TestStartInvalid(t *testing.T) {
	s := &Server{
		listen:  ":asdf",
		errChan: make(chan error, 1),
	}
	s.start()
	err := <-s.errChan
	assert.Error(t, err)
}
