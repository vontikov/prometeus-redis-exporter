package parser

import (
	"testing"

	"io/ioutil"
	"os"

	"github.com/stretchr/testify/assert"
)

func TestFSMDictionaryMarshalling(t *testing.T) {
	assert := assert.New(t)

	dir, _ := os.Getwd()
	t.Log("PWD", dir)

	b, err := ioutil.ReadFile("../../test/data/info-everything.out")
	if err != nil {
		t.Fatal("input", err)
	}

	in := string(b)
	res := Parse(&in)
	assert.Equal(126, len(res))
}
