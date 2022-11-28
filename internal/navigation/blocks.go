package navigation

import (
	"github.com/maaslalani/slides/internal/code"
	"strconv"
)

type Blocks struct {
	blocks  []code.Block
	nextIdx uint
	done    bool
}

func NewBlocks(blocks []code.Block) *Blocks {
	b := new(Blocks)
	b.blocks = blocks
	b.nextIdx = 0
	b.checkDone()

	return b
}

func (b *Blocks) ExecuteAll(m Model) {
	for _, block := range b.blocks {
		res := code.Execute(block)
		m.SetVirtualText(res.Out)
	}
	b.nextIdx = uint(len(b.blocks))
	b.checkDone()
}

func (b *Blocks) ExecuteNext(m Model) {
	if b.done {
		return
	}

	res := code.Execute(b.blocks[b.nextIdx])
	m.SetVirtualText(res.Out)

	b.nextIdx++
	b.checkDone()
}

func (b *Blocks) ExecuteIdx(idx uint, m Model) {
	if idx >= uint(len(b.blocks)) {
		m.SetVirtualText("no code block with index [" + strconv.Itoa(int(idx)) + "]")
		return
	}

	b.nextIdx = idx
	b.ExecuteNext(m)
}

func (b *Blocks) Done() bool {
	return b.done
}

func (b *Blocks) checkDone() {
	b.done = b.nextIdx >= uint(len(b.blocks))
}
