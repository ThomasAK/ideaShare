package routes

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHealth(t *testing.T) {
	res, _ := http.Get("http://localhost:3031/health")
	assert.Equal(t, 200, res.StatusCode, "Status code should be 200")
}
