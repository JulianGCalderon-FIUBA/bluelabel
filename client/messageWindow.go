package main

type messageWindow struct {
	message string
}

func newWindow(message string) messageWindow {
	return messageWindow{
		message: message,
	}
}

func (window messageWindow) show(viewportHeight int, viewportWidth int) string {
	screenCenteringStyle := screenCenteringStyle(viewportHeight, viewportWidth)
	return screenCenteringStyle.Render(windowStyle().Render(window.message))
}
