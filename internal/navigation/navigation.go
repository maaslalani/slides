package navigation

import (
	"strconv"
)

type repeatableFunction func(slide, totalSlides int) int

// State represents the current buffer, current slide, and true if virtual text
// should be cleared.
type State struct {
	Buffer    string
	Slide     int
	NumSlides int
}

// Navigate receives the current State and keyPress, and returns the new State.
func Navigate(currentState State, keyPress string) State {
	switch keyPress {
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		newBuffer := keyPress

		if bufferIsNumeric(currentState.Buffer) {
			newBuffer = currentState.Buffer + keyPress
		}

		return State{
			Buffer:    newBuffer,
			Slide:     currentState.Slide,
			NumSlides: currentState.NumSlides,
		}
	case "g":
		switch currentState.Buffer {
		case "g":
			return State {
				Buffer: "",
				Slide: navigateFirst(),
				NumSlides: currentState.NumSlides,
			}
		default:
			return State {
				Buffer: "g",
				Slide: currentState.Slide,
				NumSlides: currentState.NumSlides,
			}
		}
	case "G":
		targetSlide := navigateLast(currentState.NumSlides)
		if bufferIsNumeric(currentState.Buffer) {
			targetSlide = navigateSlide(currentState.Buffer, currentState.NumSlides)
		}

		return State {
			Buffer: "",
			Slide: targetSlide,
			NumSlides: currentState.NumSlides,
		}
	case " ", "down", "j", "right", "l", "enter", "n":
		return State {
			Buffer: "",
			Slide: navigateNext(currentState),
			NumSlides: currentState.NumSlides,
		}
	case "up", "k", "left", "h", "p":
		return State{
			Buffer: "",
			Slide:  navigatePrevious(currentState),
			NumSlides: currentState.NumSlides,
		}
	default:
		return State {
			Buffer: "",
			Slide: currentState.Slide,
			NumSlides: currentState.NumSlides,
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

func navigateNext(state State) int {
	return repeatableAction(func(slide, totalSlides int) int {
		if slide < totalSlides-1 {
			return slide + 1
		}

		return totalSlides - 1
	}, state.Buffer, state.Slide, state.NumSlides)
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

func navigatePrevious(state State) int {
	return repeatableAction(func(slide, totalSlides int) int {
		if slide > 0 {
			return slide - 1
		}

		return slide
	}, state.Buffer, state.Slide, state.NumSlides)
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
