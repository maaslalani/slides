package navigation

import (
	"strings"
	"time"
)

// KeyPressBuffer tracks a sequence of key-presses, each made
// within a configurable duration of the previous one.
type KeyPressBuffer struct {
	// Active is true if the buffer is still
	// waiting for next key-press.
	Active bool

	// onKeyPress is called when the buffer is active
	// and a key-press is a received. It is called with the
	// key-press history so far, including the just-received
	// key-press. It returns true if the just-received key-press
	// was handled; false, otherwise.
	onKeyPress func(string) bool

	// onDeactivate is called when the buffer gets deactivated.
	onDeactivate func()

	// History of key-presses, within this slide,
	// newer ones at the end.
	keyPressHistory []string

	// How long to wait for before deactivating buffer.
	waitDuration time.Duration
	waitTimer *time.Timer
}

func NewKeyPressBuffer(initialKeyPress string, waitDuration time.Duration, onKeyPress func(keyPressHistory string) bool, onDeactivate func()) *KeyPressBuffer {
	buf := new(KeyPressBuffer)

	buf.Active = true
	buf.onKeyPress = onKeyPress
	buf.onDeactivate = onDeactivate

	buf.waitDuration = waitDuration
	buf.keyPressHistory = append(buf.keyPressHistory, initialKeyPress)
	buf.waitTimer = time.AfterFunc(buf.waitDuration, func() {
		buf.Active = false
		buf.onDeactivate()
	})

	return buf
}

func ( buf *KeyPressBuffer) OnKeyPress(keyPress string) bool {
	// key-press buffer is not active; return key-press as unhandled.
	if !buf.Active {
		return false
	}

	// Reset timer as key-press buffer is still active. If the reset fails,
	// it means the timer just expired, so return key-press as unhandled.
	if !buf.waitTimer.Reset(buf.waitDuration) {
		return false
	}

	buf.keyPressHistory = append(buf.keyPressHistory, keyPress)
	return buf.onKeyPress(strings.Join(buf.keyPressHistory, " "))
}

