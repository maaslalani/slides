package model

import (
	"bufio"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/maaslalani/slides/internal/file"
	"github.com/maaslalani/slides/internal/navigation"
	"github.com/maaslalani/slides/internal/process"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/maaslalani/slides/internal/code"
	"github.com/maaslalani/slides/internal/meta"
	"github.com/maaslalani/slides/styles"
)

var (
	//go:embed tutorial.md
	slidesTutorial []byte
	tabSpaces      = strings.Repeat(" ", 4)
)

const (
	delimiter       = "\n---\n"
	pausedDelimiter = "\n+++\n"
)

// Model represents the model of this presentation, which contains all the
// state related to the current slides.
type Model struct {
	Slides   []Slide
	Page     int
	Part     int
	Author   string
	Date     string
	Theme    glamour.TermRendererOption
	Paging   string
	FileName string
	viewport viewport.Model
	buffer   string
	// VirtualText is used for additional information that is not part of the
	// original slides, it will be displayed on a slide and reset on page change
	VirtualText string
	Search      navigation.Search
}

// Slide represents a logical slide with optional paused parts.
type Slide struct {
	Parts []string
}

type fileWatchMsg struct{}

var fileInfo os.FileInfo

// Init initializes the model and begins watching the slides file for changes
// if it exists.
func (m Model) Init() tea.Cmd {
	if m.FileName == "" {
		return nil
	}
	fileInfo, _ = os.Stat(m.FileName)
	return fileWatchCmd()
}

func fileWatchCmd() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return fileWatchMsg{}
	})
}

// Load loads all of the content and metadata for the presentation.
func (m *Model) Load() error {
	var content string
	var err error

	if m.FileName != "" {
		content, err = readFile(m.FileName)
	} else {
		content, err = readStdin()
	}

	if err != nil {
		return err
	}

	content = strings.ReplaceAll(content, "\r", "")

	content = strings.TrimPrefix(content, strings.TrimPrefix(delimiter, "\n"))
	slides := strings.Split(content, delimiter)

	metaData, exists := meta.New().Parse(slides[0])
	// If the user specifies a custom configuration options
	// skip the first "slide" since this is all configuration
	if exists && len(slides) > 1 {
		slides = slides[1:]
	}

	m.Slides = splitSlides(slides)
	m.Author = metaData.Author
	m.Date = metaData.Date
	m.Paging = metaData.Paging
	if m.Theme == nil {
		m.Theme = styles.SelectTheme(metaData.Theme)
	}
	m.clampPosition()

	return nil
}

// Update updates the presentation model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		keyPress := msg.String()

		if m.Search.Active {
			switch msg.Type {
			case tea.KeyEnter:
				// execute current buffer
				if m.Search.Query() != "" {
					m.Search.Execute(&m)
				} else {
					m.Search.Done()
				}
				// cancel search
				return m, nil
			case tea.KeyCtrlC, tea.KeyEscape:
				// quit command mode
				m.Search.SetQuery("")
				m.Search.Done()
				return m, nil
			}

			var cmd tea.Cmd
			m.Search.SearchTextInput, cmd = m.Search.SearchTextInput.Update(msg)
			return m, cmd
		}

		switch keyPress {
		case "/":
			// Begin search
			m.Search.Begin()
			m.Search.SearchTextInput.Focus()
			return m, nil
		case "ctrl+n":
			// Go to next occurrence
			m.Search.Execute(&m)
		case "ctrl+e":
			// Run code blocks
			blocks, err := code.Parse(m.currentSlideContent())
			if err != nil {
				// We couldn't parse the code block on the screen
				m.VirtualText = "\n" + err.Error()
				return m, nil
			}
			var outs []string
			for _, block := range blocks {
				res := code.Execute(block)
				outs = append(outs, res.Out)
			}
			m.VirtualText = strings.Join(outs, "\n")
		case "y":
			blocks, err := code.Parse(m.currentSlideContent())
			if err != nil {
				return m, nil
			}
			for _, b := range blocks {
				_ = clipboard.WriteAll(b.Code)
			}
			return m, nil
		case "ctrl+c", "q":
			return m, tea.Quit
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			m.buffer = updateRepeatBuffer(m.buffer, keyPress)
			return m, nil
		case " ", "down", "j", "right", "l", "enter", "n", "pgdown":
			repeat := m.consumeRepeat()
			m.advanceSteps(repeat)
			return m, nil
		case "up", "k", "left", "h", "p", "pgup", "N":
			repeat := m.consumeRepeat()
			m.retreatSteps(repeat)
			return m, nil
		case "g":
			newState := navigation.Navigate(navigation.State{
				Buffer:      m.buffer,
				Page:        m.Page,
				TotalSlides: len(m.Slides),
			}, keyPress)
			m.buffer = newState.Buffer
			if newState.Page != m.Page {
				m.SetPageAndPart(newState.Page, 0)
			}
			return m, nil
		case "G":
			newState := navigation.Navigate(navigation.State{
				Buffer:      m.buffer,
				Page:        m.Page,
				TotalSlides: len(m.Slides),
			}, keyPress)
			m.buffer = newState.Buffer
			m.SetPageAndPart(newState.Page, 0)
			return m, nil
		default:
			newState := navigation.Navigate(navigation.State{
				Buffer:      m.buffer,
				Page:        m.Page,
				TotalSlides: len(m.Slides),
			}, keyPress)
			m.buffer = newState.Buffer
			m.SetPageAndPart(newState.Page, 0)
		}

	case fileWatchMsg:
		newFileInfo, err := os.Stat(m.FileName)
		if err == nil && newFileInfo.ModTime() != fileInfo.ModTime() {
			fileInfo = newFileInfo
			_ = m.Load()
			m.clampPosition()
		}
		return m, fileWatchCmd()
	}
	return m, nil
}

// View renders the current slide in the presentation and the status bar which
// contains the author, date, and pagination information.
func (m Model) View() string {
	r, _ := glamour.NewTermRenderer(m.Theme, glamour.WithWordWrap(m.viewport.Width))
	slide := m.currentSlideContent()
	slide = code.HideComments(slide)
	slide, err := r.Render(slide)
	slide = strings.ReplaceAll(slide, "\t", tabSpaces)
	slide += m.VirtualText
	if err != nil {
		slide = fmt.Sprintf("Error: Could not render markdown! (%v)", err)
	}
	slide = styles.Slide.Render(slide)

	var left string
	if m.Search.Active {
		// render search bar
		left = m.Search.SearchTextInput.View()
	} else {
		// render author and date
		left = styles.Author.Render(m.Author) + styles.Date.Render(m.Date)
	}

	right := styles.Page.Render(m.paging())
	status := styles.Status.Render(styles.JoinHorizontal(left, right, m.viewport.Width))
	return styles.JoinVertical(slide, status, m.viewport.Height)
}

func (m *Model) paging() string {
	count := strings.Count(m.Paging, "%d")
	switch count {
	case 0:
		return m.Paging
	case 1:
		return fmt.Sprintf(m.Paging, m.Page+1)
	case 2:
		return fmt.Sprintf(m.Paging, m.Page+1, len(m.Slides))
	case 3:
		return fmt.Sprintf(m.Paging, m.Page+1, len(m.Slides), m.Part+1)
	default:
		// 4 or more: slide, totalSlides, part, totalParts
		totalParts := m.stepsInSlide(m.Page)
		return fmt.Sprintf(m.Paging, m.Page+1, len(m.Slides), m.Part+1, totalParts)
	}
}

func readFile(path string) (string, error) {
	s, err := os.Stat(path)
	if err != nil {
		return "", errors.New("could not read file")
	}
	if s.IsDir() {
		return "", errors.New("can not read directory")
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	content := string(b)

	// Pre-process slides if the file is executable to avoid
	// unintentional code execution when presenting slides
	if file.IsExecutable(s) {
		// Remove shebang if file has one
		if strings.HasPrefix(content, "#!") {
			content = strings.Join(strings.SplitN(content, "\n", 2)[1:], "\n")
		}

		content = process.Pre(content)
	}

	return content, err
}

func readStdin() (string, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	if stat.Mode()&os.ModeNamedPipe == 0 && stat.Size() == 0 {
		return string(slidesTutorial), nil
	}

	reader := bufio.NewReader(os.Stdin)
	var b strings.Builder

	for {
		r, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		_, err = b.WriteRune(r)
		if err != nil {
			return "", err
		}
	}

	return b.String(), nil
}

// CurrentPage returns the current page the presentation is on.
func (m *Model) CurrentPage() int {
	return m.Page
}

// CurrentPart returns the current part the presentation is on.
func (m *Model) CurrentPart() int {
	return m.Part
}

// SetPage sets which page the presentation should render.
func (m *Model) SetPage(page int) {
	m.SetPageAndPart(page, 0)
}

// SetPageAndPart sets which page and part the presentation should render.
func (m *Model) SetPageAndPart(page int, part int) {
	if len(m.Slides) == 0 {
		m.Page = 0
		m.Part = 0
		m.VirtualText = ""
		return
	}
	if page < 0 {
		page = 0
	}
	if page >= len(m.Slides) {
		page = len(m.Slides) - 1
	}
	maxPart := m.stepsInSlide(page) - 1
	if part < 0 {
		part = 0
	}
	if part > maxPart {
		part = maxPart
	}
	if m.Page == page && m.Part == part {
		return
	}
	m.VirtualText = ""
	m.Page = page
	m.Part = part
}

// Pages returns all the slides in the presentation.
func (m *Model) Pages() []string {
	pages := make([]string, len(m.Slides))
	for i, slide := range m.Slides {
		pages[i] = strings.Join(slide.Parts, "\n")
	}
	return pages
}

// SlideParts returns the paused parts for a specific slide.
func (m *Model) SlideParts(page int) []string {
	if page < 0 || page >= len(m.Slides) {
		return nil
	}
	parts := m.Slides[page].Parts
	if len(parts) == 0 {
		return []string{""}
	}
	return parts
}

func splitSlides(slides []string) []Slide {
	parsed := make([]Slide, 0, len(slides))
	for _, slide := range slides {
		parts := strings.Split(slide, pausedDelimiter)
		if len(parts) == 0 {
			parts = []string{slide}
		}
		parsed = append(parsed, Slide{Parts: parts})
	}
	return parsed
}

func (m *Model) stepsInSlide(page int) int {
	if page < 0 || page >= len(m.Slides) {
		return 0
	}
	parts := m.Slides[page].Parts
	if len(parts) == 0 {
		return 1
	}
	return len(parts)
}

func (m *Model) currentSlideContent() string {
	if len(m.Slides) == 0 {
		return ""
	}
	parts := m.Slides[m.Page].Parts
	if len(parts) == 0 {
		return ""
	}
	part := m.Part
	if part < 0 {
		part = 0
	}
	if part >= len(parts) {
		part = len(parts) - 1
	}
	return strings.Join(parts[:part+1], "\n")
}

func (m *Model) clampPosition() {
	if len(m.Slides) == 0 {
		m.Page = 0
		m.Part = 0
		return
	}
	if m.Page < 0 {
		m.Page = 0
	}
	if m.Page >= len(m.Slides) {
		m.Page = len(m.Slides) - 1
	}
	maxPart := m.stepsInSlide(m.Page) - 1
	if m.Part < 0 {
		m.Part = 0
	}
	if m.Part > maxPart {
		m.Part = maxPart
	}
}

func (m *Model) advanceSteps(steps int) {
	if len(m.Slides) == 0 {
		return
	}
	if steps < 0 {
		m.retreatSteps(-steps)
		return
	}
	for i := 0; i < steps; i++ {
		parts := m.stepsInSlide(m.Page)
		if parts == 0 {
			return
		}
		if m.Part < parts-1 {
			m.SetPageAndPart(m.Page, m.Part+1)
			continue
		}
		if m.Page < len(m.Slides)-1 {
			m.SetPageAndPart(m.Page+1, 0)
			continue
		}
		break
	}
}

func (m *Model) retreatSteps(steps int) {
	if len(m.Slides) == 0 {
		return
	}
	if steps < 0 {
		m.advanceSteps(-steps)
		return
	}
	for i := 0; i < steps; i++ {
		if m.Part > 0 {
			m.SetPageAndPart(m.Page, m.Part-1)
			continue
		}
		if m.Page > 0 {
			prevPage := m.Page - 1
			m.SetPageAndPart(prevPage, m.stepsInSlide(prevPage)-1)
			continue
		}
		break
	}
}

func (m *Model) consumeRepeat() int {
	if !bufferIsNumeric(m.buffer) {
		m.buffer = ""
		return 1
	}
	repeat, _ := strconv.Atoi(m.buffer)
	m.buffer = ""
	if repeat == 0 {
		return 1
	}
	return repeat
}

func updateRepeatBuffer(buffer string, digit string) string {
	if bufferIsNumeric(buffer) {
		return buffer + digit
	}
	return digit
}

func bufferIsNumeric(buffer string) bool {
	_, err := strconv.Atoi(buffer)
	return err == nil
}
