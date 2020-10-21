package nat

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddMapping(t *testing.T) {
	nat_table = nil

	AddMapping("1.1.1.1", "80", "2.2.2.2", "80")
	assert.Equal(t, len(nat_table), 1)

	AddMapping("2.2.2.2", "80", "3.3.3.3", "80")
	assert.Equal(t, len(nat_table), 2)

	AddMapping("2.2.2.2", "80", "1.1.1.1", "60")
	assert.Equal(t, len(nat_table), 2)

	exp := map[string]string{
		"1.1.1.1/80": "2.2.2.2/80",
		"2.2.2.2/80": "1.1.1.1/60",
	}

	eq := reflect.DeepEqual(exp, nat_table)
	if !eq {
		t.Errorf("Incorrect mapping. Expected %v Got %v", exp, nat_table)
	}
}

func TestGetMapping(t *testing.T) {
	nat_table = nil
	AddMapping("1.1.1.1", "80", "2.2.2.2", "80")
	AddMapping("2.2.2.2", "80", "3.3.3.3", "80")
	AddMapping("2.2.2.2", "60", "3.3.3.3", "70")
	AddMapping("3.3.3.3", "50", "4.4.4.4", "60")

	assert.Equal(t, len(nat_table), 4)

	ip, port, err := GetMapping("2.2.2.2", "80")
	assert.Equal(t, nil, err)
	assert.Equal(t, "80", port)
	assert.Equal(t, "3.3.3.3", ip)

	ip, port, err = GetMapping("2.2.2.2", "60")
	assert.Equal(t, nil, err)
	assert.Equal(t, "70", port)
	assert.Equal(t, "3.3.3.3", ip)

	ip, port, err = GetMapping("4.4.4.4", "60")
	assert.Equal(t, fmt.Errorf("Not Found"), err)
	assert.Equal(t, "", port)
	assert.Equal(t, "", ip)
}
