# GitViz 🎨

[![Go](https://img.shields.io/badge/Go-1.24.3-00ADD8?style=flat&logo=go&logoColor=white)](https://golang.org/)
[![Git](https://img.shields.io/badge/Git-F05032?style=flat&logo=git&logoColor=white)](https://git-scm.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Version](https://img.shields.io/badge/Version-1.0.1-blue.svg)](https://github.com/SwanHtetAungPhyo/gitviz/releases)
[![CLI](https://img.shields.io/badge/CLI-Tool-purple.svg)](https://github.com/SwanHtetAungPhyo/gitviz)

A beautiful and feature-rich Git repository visualizer built in Go. Transform your Git history into stunning visual representations with colorful output, detailed statistics, and multiple viewing modes.

## Features

✨ **Beautiful Terminal UI** - Styled with lipgloss for elegant borders and colors  
📊 **Repository Statistics** - Comprehensive stats including commits, file changes, and author contributions  
📈 **Timeline Visualization** - Weekly commit activity with visual bar charts  
🌳 **Commit Graph** - Visual representation of your Git history with branch information  
🎯 **Multiple Display Modes** - Compact and detailed views to suit your needs  
👥 **Author Analytics** - Top contributor analysis with commit counts  
🔄 **Merge Detection** - Special highlighting for merge commits  

## Installation

### Prerequisites
- Go 1.19 or higher
- Git repository to visualize

### Build from Source
```bash
git clone https://github.com/SwanHtetAungPhyo/gitviz.git
cd gitviz
go mod tidy
go build -o gitviz
```

### Install Dependencies
```bash
go get github.com/charmbracelet/lipgloss
go get github.com/fatih/color
go get github.com/go-git/go-git/v5
go get github.com/urfave/cli/v2
```

## Usage

Navigate to any Git repository and run GitViz:

```bash
# Basic usage - shows commit graph + statistics
./gitviz

# Limit number of commits
./gitviz --limit 20
./gitviz -n 20

# Compact view for quick overview
./gitviz --compact
./gitviz -c

# Show only statistics
./gitviz --stats
./gitviz -s

# Show commit timeline
./gitviz --timeline
./gitviz -t
```

## Command Line Options

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--limit` | `-n` | Limit number of commits to display | 50 |
| `--compact` | `-c` | Use compact display mode | false |
| `--stats` | `-s` | Show only repository statistics | false |
| `--timeline` | `-t` | Show commit timeline visualization | false |

## Display Modes

### 1. Default Mode
Shows a detailed commit graph with:
- Full commit hashes (7 characters)
- Author name and email
- Commit date and time
- Complete commit messages
- Files changed with names
- Parent commit information
- Branch information
- Merge commit detection

### 2. Compact Mode (`--compact`)
Streamlined view showing:
- Short commit hash
- Truncated commit message (50 chars)
- Author name
- Branch information
- Merge indicators

### 3. Statistics Only (`--stats`)
Comprehensive repository analytics:
- Total commit count
- Files changed summary
- Lines added/deleted
- Top 5 contributors by commit count

### 4. Timeline View (`--timeline`)
Visual timeline showing:
- Weekly commit activity
- Bar chart representation
- Commit count per week
- Chronological progression

## Output Features

### 🎨 Color Coding
- **Red**: Commit hashes
- **Green**: Author names
- **White**: Commit messages
- **Yellow**: Branch names
- **Blue**: File names
- **Cyan**: Metadata (dates, parents)
- **Magenta**: Merge commits

### 📋 Information Display
- **Merge Detection**: Special highlighting for merge commits
- **Branch Mapping**: Shows which branches contain each commit
- **File Changes**: Lists all modified files per commit
- **Parent Tracking**: Shows commit relationships
- **Author Statistics**: Ranked contributor analysis

## Examples

### Basic Repository Visualization
```bash
./gitviz --limit 10
```
Shows the last 10 commits with full details, followed by repository statistics.

### Quick Overview
```bash
./gitviz --compact --limit 20
```
Compact view of the last 20 commits in a single line format.

### Project Analytics
```bash
./gitviz --stats
```
Displays comprehensive repository statistics without the commit graph.

### Activity Timeline
```bash
./gitviz --timeline
```
Shows weekly commit activity with visual bar charts.

## Technical Details

### Dependencies
- **lipgloss**: Terminal styling and layout
- **color**: ANSI color support
- **go-git**: Pure Go Git implementation
- **cli**: Command-line interface framework

### Architecture
- **Visualizer**: Core engine for Git repository analysis
- **CommitInfo**: Structured commit data representation
- **SafeHashShort**: Utility for hash truncation
- **Multiple Display**: Flexible rendering system

### Performance
- Efficient Git object traversal
- Configurable commit limits
- Memory-conscious design
- Fast branch and reference loading

## Development

### Project Structure
```
gitviz/
├── main.go           # Main application entry point
├── go.mod           # Go module definition
├── go.sum           # Dependency checksums
└── README.md        # This file
```

### Key Components
- **NewVisualizer()**: Repository initialization
- **LoadCommits()**: Commit history parsing
- **LoadBranches()**: Reference and branch loading
- **DisplayGraph()**: Visual commit representation
- **DisplayStats()**: Statistics calculation and display
- **DisplayTimeline()**: Timeline visualization

## Version Information

- **Version**: 1.0.1
- **Author**: SwanHtet Aung Phyo
- **Title**: Computer Science Student at AGH, Poland - Backend Developer
- **Repository**: [github.com/SwanHtetAungPhyo/gitviz](https://github.com/SwanHtetAungPhyo/gitviz)

## Contributing

Contributions are welcome! Please feel free to submit issues, feature requests, or pull requests to [github.com/SwanHtetAungPhyo/gitviz](https://github.com/SwanHtetAungPhyo/gitviz).

### Development Setup
```bash
git clone https://github.com/SwanHtetAungPhyo/gitviz.git
cd gitviz
go mod tidy
go run main.go --help
```

## License

This project is open source. Please check the repository for license details.

---

**GitViz** - Transform your Git history into beautiful visualizations! 🚀