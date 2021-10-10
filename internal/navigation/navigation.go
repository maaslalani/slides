package navigation

import (
	"strconv"
)

type repeatableFunction func(slide, totalSlides int) int

// Navigate receives the current buffer, keyPress, current slide, and total number of slides.
// Navigate returns the new (updated) buffer, new current slide, and true if virtual text should be cleared.
// For example, if showing user slide 1 and there are 10 slides available, slide will be 0 and numSlides will be 10.
func Navigate(buffer, keyPress string, slide, numSlides int) (string, int, bool) {
	// Implementation

	switch keyPress {
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		if bufferIsNumeric(buffer) {
			return buffer + keyPress, slide, false
		} else {
			return keyPress, slide, false
		}
	case "g":
		switch buffer {
		case "g":
			return "", navigateFirst(), false
		default:
			return "g", slide, false
		}
	case "G":
		if bufferIsNumeric(buffer) {
			return "", navigateSlide(buffer, numSlides), false
		} else {
			return "", navigateLast(numSlides), false
		}
	case " ", "down", "j", "right", "l", "enter", "n":
		return "", navigateNext(buffer, slide, numSlides), true
	case "up", "k", "left", "h", "p":
		return "", navigatePrevious(buffer, slide, numSlides), true
	default:
		return "", slide, false
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
