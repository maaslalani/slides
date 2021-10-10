package navigation

import (
	"strconv"
)

type repeatableFunction func(slide, totalSlides int) int

// State represents the current buffer, current slide, and true if virtual text
// should be cleared.
type State struct {
	Buffer string
	Slide int
	ClearVirtualText bool
}

// Navigate receives the current State and keyPress, and returns the new State.
func Navigate(currentState State, keyPress string, numSlides int) State {
	// Implementation

	switch keyPress {
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		newBuffer := keyPress

		if bufferIsNumeric(currentState.Buffer) {
			newBuffer = currentState.Buffer + keyPress
		}

		return State{
			Buffer:           newBuffer,
			Slide:            currentState.Slide,
			ClearVirtualText: false,
		}
	case "g":
		switch currentState.Buffer {
		case "g":
			return State {
				Buffer: "",
				Slide: navigateFirst(),
				ClearVirtualText: false,
			}
		default:
			return State {
				Buffer: "g",
				Slide: currentState.Slide,
				ClearVirtualText: false,
			}
		}
	case "G":
		if bufferIsNumeric(currentState.Buffer) {
			return State {
				Buffer: "",
				Slide: navigateSlide(currentState.Buffer, numSlides),
				ClearVirtualText: false,
			}
		} else {
			return State {
				Buffer: "",
				Slide: navigateLast(numSlides),
				ClearVirtualText: false,
			}
		}
	case " ", "down", "j", "right", "l", "enter", "n":
		return State {
			Buffer: "",
			Slide: navigateNext(currentState.Buffer, currentState.Slide, numSlides),
			ClearVirtualText: true,
		}
	case "up", "k", "left", "h", "p":
		return State{
			Buffer: "",
			Slide:  navigatePrevious(currentState.Buffer, currentState.Slide, numSlides),
			ClearVirtualText: true,
		}
	default:
		return State {
			Buffer: "",
			Slide: currentState.Slide,
			ClearVirtualText: false,
		}
	}
}

func bufferIsNumeric(buffer string) bool {
	_, err := strconv.Atoi(buffer)
	return err == nil
}

func navigateFirst() int {
	return 0
}

func navigateNext(buffer string, slide, numSlides int) int {
	return repeatableAction(func(slide, totalSlides int) int {
		if slide < totalSlides-1 {
			return slide + 1
		}

		return totalSlides - 1
	}, buffer, slide, numSlides)
}

func navigateSlide(buffer string, numSlides int) int {
	destinationSlide, _ := strconv.Atoi(buffer)
	destinationSlide -= 1

	if destinationSlide > numSlides -1 {
		return numSlides - 1
	}

	if destinationSlide < 0 {
		return 0
	}

	return destinationSlide
}

func navigatePrevious(buffer string, slide, totalSlides int) int {
	return repeatableAction(func(slide, totalSlides int) int {
		if slide > 0 {
			return slide - 1
		}

		return slide
	}, buffer, slide, totalSlides)
}

func navigateLast(numSlides int) int {
	return numSlides - 1
}

func repeatableAction(fn repeatableFunction, buffer string, slide, totalSlides int) int {
	if !bufferIsNumeric(buffer) {
		return fn(slide, totalSlides)
	}

	repeat, _ := strconv.Atoi(buffer)
	currentSlide := slide

	if repeat == 0 {
		// This is how behaviour works in Vim, so following principle of least astonishment.
		return fn(slide, totalSlides)
	}

	for i := 0; i < repeat; i++ {
		currentSlide = fn(currentSlide, totalSlides)
	}

	return currentSlide
}
