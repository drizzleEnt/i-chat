# i-chat
A simple desktop GUI chat application written in Go using the Fyne toolkit.

## Features
- Clean, intuitive chat interface
- User authentication (login/register)
- Real-time messaging

## Getting Started

### Prerequisites
- Go 1.24.3 or later
- A graphical display (X11/Wayland on Linux)

### Build
```bash
go build ./...
```

### Run
```bash
go run main.go
```

## Project Structure
- `main.go` — application entry point
- `app/` — application container and lifecycle
- `ui/` — Fyne UI implementation

## Development
Run `go vet` and `go build ./...` before submitting changes. Test GUI changes locally with a display available.