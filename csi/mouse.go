package csi

const (
	// See https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h2-Mouse-Tracking for definitions how
	// xterm handles mouse input.

	// reports mouse button presses and releases
	mouseTracking        = "1000"
	EnableMouseTracking  = CSI + "?" + mouseTracking + "h"
	DisableMouseTracking = CSI + "?" + mouseTracking + "l"

	// changes encoding of mouse events to use SGR notation
	mouseSGREncoding        = "1006"
	EnableMouseSGREncoding  = CSI + "?" + mouseSGREncoding + "h"
	DisableMouseSGREncoding = CSI + "?" + mouseSGREncoding + "l"

	// reports mouse highlighting, a.k.a. selection
	mouseHighlight        = "1001"
	EnableMouseHighlight  = CSI + "?" + mouseHighlight + "h"
	DisableMouseHighlight = CSI + "?" + mouseHighlight + "l"

	// reports mouse movement when a button is pressed
	mouseButtonEvent        = "1002"
	EnableMouseButtonEvent  = CSI + "?" + mouseButtonEvent + "h"
	DisableMouseButtonEvent = CSI + "?" + mouseButtonEvent + "l"

	// // reports any kind of movement
	// mouseMovement        = "1003"
	// EnableMouseMovement  = CSI + "?" + mouseMovement + "h"
	// DisableMouseMovement = CSI + "?" + mouseMovement + "l"
)
