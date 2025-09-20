# Tiny Time Tracker

A small Terminal User Interface (TUI) application for tracking time and reporting it by weeks. Built with Go and the [Bubble Tea](https://github.com/charmbracelet/bubbletea) library as an excuse to explore TUI development! 🚀

## Features

- ⏱️ **Simple Time Tracking**: Start and stop timers with a single keypress
- 📊 **Weekly Reports**: View time tracking data organized by weeks
- 📝 **Last Interval Display**: See details of your most recent time session

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
./tiny_time_tracker
```

## Usage

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Space` | Start/Stop timer |
| `e` | Edit last interval (planned feature) |
| `Tab` | Switch between tabs (planned feature) |
| `q` or `Ctrl+C` | Quit application |

### How to Use

1. **Start Tracking**: Press `Space` to begin a new time session
2. **Stop Tracking**: Press `Space` again to end the current session
3. **View History**: The left panel shows your time tracking history
4. **Check Last Session**: When not running, the right panel displays details of your last completed session

## Planned Features

The following features are planned for future releases:

- ✏️ **Last Interval Editor**: Ability to edit the start/end times of the most recent session
- 📁 **CSV File**: Save/Load time tracking history to a CSV file. 