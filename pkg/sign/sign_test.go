package sign_test

import (
	"fmt"
	"testing"

	"github.com/benderr/metrics/pkg/sign"
	"github.com/stretchr/testify/assert"
)

func ExampleNew() {
	v := sign.New("secret", []byte("123"))
	fmt.Println(v)
}

func TestNew(t *testing.T) {
	v := sign.New("secret", []byte("123"))
	assert.Equal(t, v, "77de38e4b50e618a0ebb95db61e2f42697391659d82c064a5f81b9f48d85ccd5")
}
