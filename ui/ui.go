package ui

import (
	"log"
	"sync"
	"time"

	"github.com/jroimartin/gocui"
)

type UI struct {
	gui          *gocui.Gui
	currentView  int
	nextItem     int
	currentModal string
	consoleLog   string
	cursors      Cursors
	modalTimer   *time.Timer
	logTimer     *time.Timer
	mutex        *sync.Mutex
	fontpath     string
}

func NewUI(fontpath string) *UI {
	var err error

	ui := new(UI)
	ui.gui, err = gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Panicln(err)
	}

	ui.cursors = NewCursors()
	ui.fontpath = fontpath
	return ui
}

func (ui *UI) Init() {
	if err := ui.initGui(ui.gui); err != nil {
		log.Panicln(err)
	}
}

// Cursors stores the cursor position for a specific panel view.
// Used to restore mouse position when click is detected.
type Cursors map[string]struct{ x, y int }

func NewCursors() Cursors {
	return make(Cursors)
}

func (c Cursors) Restore(view *gocui.View) error {
	return view.SetCursor(c.Get(view.Name()))
}

func (c Cursors) Get(view string) (int, int) {
	if v, ok := c[view]; ok {
		return v.x, v.y
	}
	return 0, 0
}

func (c Cursors) Set(view string, x, y int) {
	c[view] = struct{ x, y int }{x, y}
}

// Loop starts the GUI loop
func (ui *UI) Loop() {
	if err := ui.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

// Close closes the app
func (ui *UI) Close() {
	ui.gui.Close()
}

// initGui initializes the GUI
func (ui *UI) initGui(g *gocui.Gui) error {
	// Default Panel settings
	ui.gui.Highlight = true
	ui.gui.InputEsc = false
	ui.gui.SelFgColor = gocui.ColorGreen

	// Mouse settings
	ui.gui.Cursor = true
	ui.gui.Mouse = true

	ui.currentView = ui.findViewByName(DIAGRAM_PANEL)
	ui.nextItem = 0

	// Set Layout function
	ui.gui.SetManager(ui)

	if err := keyHandlers.ApplyKeyBindings(ui, g); err != nil {
		return err
	}
	return nil
}
