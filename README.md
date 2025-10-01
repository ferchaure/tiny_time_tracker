# Tiny Time Tracker

A small Terminal User Interface (TUI) application for tracking time and reporting it by weeks. Built with Go and the [Bubble Tea](https://github.com/charmbracelet/bubbletea) library as an excuse to explore TUI development! 🚀

## Features

- ⏱️ **Simple Time Tracking**: Start and stop timers with a single keypress
- 📊 **Weekly Reports**: View time tracking data organized by weeks
- 📝 **Last Interval Display**: See details of your most recent time session
- ✏️ **Edit Mode**: Modify the start and end times of your last completed interval

## Installation

### Prerequisites

- Go 1.25.1 or later

### Build from Source

1. Clone the repository:
```bash
git clone https://github.com/ferchaure/tiny_time_tracker.git
cd tiny_time_tracker
```

2. Build the application:
```bash
go build -o tiny_time_tracker
```

3. Run the application:
```bash
./tiny_time_tracker -f data.csv
```

## Usage
### Arguments

- `-f <file>`: CSV filename to read/write. Defaults to `data.csv`.
- `-wd <weekday>`: Weekday to start the week (Sunday=0, Monday=1, ..., Saturday=6). Defaults to `1` (Monday).


### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Space` | Start/Stop timer |
| `e` | Edit last interval / Accept edit |
| `Tab` | Switch between input fields (in edit mode) |
| `q` or `Ctrl+C` | Quit application |

### How to Use

1. **Start Tracking**: Press `Space` to begin a new time session
2. **Stop Tracking**: Press `Space` again to end the current session
3. **View History**: The left panel shows your time tracking history
4. **Check Last Session**: When not running, the right panel displays details of your last completed session
5. **Edit Last Interval**: Press `e` to edit the start and end times of your last completed session
   - Use `Tab` to switch between the "From" and "End" input fields
   - Enter times in the format: `15:04:05 2006/01/02` (HH:MM:SS YYYY/MM/DD)
   - Press `e` again to accept your changes
