package api

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTreeInsert(t *testing.T) {
	tree := Tree{}

	testCases := []struct {
		name  string
		value int
		err   error
	}{
		{
			name:  "insert root",
			value: 10,
		},
		{
			name:  "insert root right node",
			value: 11,
		},
		{
			name:  "insert root right node's right",
			value: 12,
		},
		{
			name:  "insert root left node",
			value: 9,
		},
		{
			name:  "insert root left node's left",
			value: 8,
		},
		{
			name:  "insert existing",
			value: 10,
			err:   errors.New("this node value already exists"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tree.Insert(tc.value)
			assert.Equal(t, tc.err, err)

			assert.Equal(t, true, tree.Search(tc.value))
		})
	}
}

func TestTreeSearch(t *testing.T) {
	tree := Tree{}
	values := []int{1, 2, 3, 4, 7, 99}

	for _, value := range values {
		assert.NoError(t, tree.Insert(value))
		assert.Equal(t, true, tree.Search(value))
	}

	assert.Equal(t, false, tree.Search(404))

	tree = Tree{}
	assert.Equal(t, false, tree.Search(000))
}

func TestTreeDelete(t *testing.T) {
	tree := Tree{}
	values := []int{1, 2, 3, 4, 7, 99, 101, 222, 8}

	for _, value := range values {
		assert.NoError(t, tree.Insert(value))
		assert.Equal(t, true, tree.Search(value))
	}

	assert.Error(t, tree.Delete(404))

	for i := len(values) - 1; i >= 0; i-- {
		assert.NoError(t, tree.Delete(values[i]))
	}

	tree = Tree{}
	assert.Error(t, tree.Delete(000))
}
