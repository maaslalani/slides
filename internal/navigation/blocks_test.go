package navigation

import (
	"github.com/maaslalani/slides/internal/code"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBlocks_ExecuteAll(t *testing.T) {
	b := []code.Block{
		{},
		{},
		{},
	}
	m := mockModel{}

	blocks := NewBlocks(b)
	require.False(t, blocks.Done())

	blocks.ExecuteAll(&m)
	require.True(t, blocks.Done())
}

func TestBlocks_ExecuteNext(t *testing.T) {
	b := []code.Block{
		{
			Language: "javascript",
			Code:     `console.log("first block")`,
		},
		{
			Language: "javascript",
			Code:     `console.log("second block")`,
		},
		{
			Language: "javascript",
			Code:     `console.log("third block")`,
		},
	}
	m := mockModel{}

	blocks := NewBlocks(b)
	require.False(t, blocks.Done())

	blocks.ExecuteNext(&m)
	require.Equal(t, "first block\n", m.virtualText)
	require.False(t, blocks.Done())

	blocks.ExecuteNext(&m)
	require.Equal(t, "second block\n", m.virtualText)
	require.False(t, blocks.Done())

	blocks.ExecuteNext(&m)
	require.Equal(t, "third block\n", m.virtualText)
	require.True(t, blocks.Done())
}

func TestBlocks_ExecuteIdx(t *testing.T) {
	b := []code.Block{
		{
			Language: "javascript",
			Code:     `console.log("first block")`,
		},
		{
			Language: "javascript",
			Code:     `console.log("second block")`,
		},
		{
			Language: "javascript",
			Code:     `console.log("third block")`,
		},
	}
	m := mockModel{}

	blocks := NewBlocks(b)
	require.False(t, blocks.Done())

	blocks.ExecuteIdx(1, &m)
	require.Equal(t, "second block\n", m.virtualText)
	require.False(t, blocks.Done())

	blocks.ExecuteIdx(0, &m)
	require.Equal(t, "first block\n", m.virtualText)
	require.False(t, blocks.Done())

	blocks.ExecuteIdx(2, &m)
	require.Equal(t, "third block\n", m.virtualText)
	require.True(t, blocks.Done())
}
