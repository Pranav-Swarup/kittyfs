package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

type item struct {
	Name  string // displayed in the list
	IsDir bool
	Path  string
}

// Implement list.Item interface
func (i item) Title() string       { return i.Name }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.Name }

type model struct {
	list             list.Model
	currentDir       string
	driveSelect      bool
	drives           []item
	width, height    int
	showExtendedHelp bool
	borderColor      lipgloss.Color
	highlightColor   lipgloss.Color
	paddingX         int
	paddingY         int
	maxItemWidth     int
}

type Config struct {
	BorderColor    string `json:"border_color"`
	HighlightColor string `json:"highlight_color"`
}

func toListItems(items []item) []list.Item {
	res := make([]list.Item, len(items))
	for i, v := range items {
		res[i] = v
	}
	return res
}

func getDrives() []item {
	letters := "CDEFGHIJKLMNOPQRSTUVWXYZ"
	drives := []item{}
	for _, l := range letters {
		drive := string(l) + ":\\"
		if _, err := os.Stat(drive); err == nil {
			drives = append(drives, item{
				Name:  drive,
				IsDir: true,
				Path:  drive,
			})
		}
	}
	return drives
}

func listDirectory(path string) []item {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return []item{}
	}
	var items []item
	for _, e := range entries {
		items = append(items, item{
			Name:  e.Name(),
			IsDir: e.IsDir(),
			Path:  filepath.Join(path, e.Name()),
		})
	}
	return items
}

func openFile(path string) {
	switch runtime.GOOS {
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", path).Start()
	case "darwin":
		exec.Command("open", path).Start()
	default:
		exec.Command("xdg-open", path).Start()
	}
}

func loadConfig() Config {
	data, err := os.ReadFile("config.json")
	if err != nil {
		return Config{BorderColor: "#FF69B4", HighlightColor: "#FF1493"}
	}
	var cfg Config
	json.Unmarshal(data, &cfg)
	return cfg
}

func saveConfig(cfg Config) {
	data, _ := json.MarshalIndent(cfg, "", "  ")
	os.WriteFile("config.json", data, 0644)
}

func calculateMaxWidth(items []item) int {
	maxWidth := 20
	for _, item := range items {
		if len(item.Name) > maxWidth {
			maxWidth = len(item.Name)
		}
	}
	// Add some buffer for the list styling
	return maxWidth + 10
}

func initialModel() model {
	cfg := loadConfig()

	items := getDrives()
	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(cfg.HighlightColor)).
		Bold(true)
	d.Styles.SelectedDesc = lipgloss.NewStyle().
		Foreground(lipgloss.Color(cfg.HighlightColor))

	maxWidth := calculateMaxWidth(items)
	// Ensure width fits the help bar (approximately 80 chars for help text)
	minWidthForHelp := 80
	if maxWidth < minWidthForHelp {
		maxWidth = minWidthForHelp
	}

	l := list.New(toListItems(items), d, maxWidth, 30)
	l.SetSize(maxWidth, 30)
	
	l.Title = "Select Drive =^..^="
	l.SetShowHelp(true)
	l.KeyMap.CursorUp.SetHelp("↑/↓", "files")
	l.KeyMap.CursorDown.SetHelp("←/→", "pages")
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("backspace"),
				key.WithHelp("backspace", "parent"),
			),
		}
	}

	return model{
		list:             l,
		currentDir:       "",
		driveSelect:      true,
		drives:           items,
		borderColor:      lipgloss.Color(cfg.BorderColor),
		highlightColor:   lipgloss.Color(cfg.HighlightColor),
		paddingX:         2,
		paddingY:         1,
		maxItemWidth:     maxWidth,
		showExtendedHelp: false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Calculate available space for content
		borderOverhead := 6 // border lines + title
		
		contentHeight := m.height - borderOverhead - (m.paddingY * 2)
		
		if contentHeight < 10 {
			contentHeight = 10
		}

		m.list.SetSize(m.maxItemWidth, contentHeight)

	case tea.KeyMsg:
		switch msg.String() {
		case "?":
			m.showExtendedHelp = !m.showExtendedHelp
			m.list.SetShowHelp(!m.showExtendedHelp)

		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			if m.driveSelect {
				// choose drive
				if i, ok := m.list.SelectedItem().(item); ok {
					m.currentDir = i.Path
					entries := listDirectory(i.Path)
					m.maxItemWidth = calculateMaxWidth(entries)
					m.list.SetItems(toListItems(entries))
					m.list.Title = fmt.Sprintf("Browsing %s", m.currentDir)
					m.driveSelect = false
					m.list.ResetFilter()
					
					// Update list size with new width
					contentHeight := m.height - 6 - (m.paddingY * 2)
					m.list.SetSize(m.maxItemWidth, contentHeight)
				}
				return m, nil
			} else {
				if i, ok := m.list.SelectedItem().(item); ok {
					if i.IsDir {
						entries := listDirectory(i.Path)
						m.currentDir = i.Path
						m.maxItemWidth = calculateMaxWidth(entries)
						m.list.SetItems(toListItems(entries))
						m.list.Title = fmt.Sprintf("Browsing %s", m.currentDir)
						m.list.ResetFilter()
						
						// Update list size with new width
						contentHeight := m.height - 6 - (m.paddingY * 2)
						m.list.SetSize(m.maxItemWidth, contentHeight)
					} else {
						openFile(i.Path)
					}
				}
			}
			return m, nil

		case "backspace":
			if !m.driveSelect {
				parent := filepath.Dir(m.currentDir)
				if parent == m.currentDir || parent == "" || len(m.currentDir) <= 3 {
					// go back to drive select
					m.maxItemWidth = calculateMaxWidth(m.drives)
					m.list.SetItems(toListItems(m.drives))
					m.list.Title = "Select Drive =^..^="
					m.driveSelect = true
					m.currentDir = ""
					m.list.ResetFilter()
					
					// Update list size
					contentHeight := m.height - 6 - (m.paddingY * 2)
					m.list.SetSize(m.maxItemWidth, contentHeight)
				} else {
					entries := listDirectory(parent)
					m.currentDir = parent
					m.maxItemWidth = calculateMaxWidth(entries)
					m.list.SetItems(toListItems(entries))
					m.list.Title = fmt.Sprintf("Browsing %s", m.currentDir)
					m.list.ResetFilter()
					
					// Update list size
					contentHeight := m.height - 6 - (m.paddingY * 2)
					m.list.SetSize(m.maxItemWidth, contentHeight)
				}
			}
			return m, nil

		case "t":
			themes := []struct {
				border, highlight string
			}{
				{"#FF69B4", "#FF1493"}, // Hot Pink
				{"#00CED1", "#00FFFF"}, // Turquoise/Cyan
				{"#9370DB", "#BA55D3"}, // Purple
				{"#FF6347", "#FF4500"}, // Tomato/Red
				{"#32CD32", "#7FFF00"}, // Lime Green
				{"#FFD700", "#FFA500"}, // Gold/Orange
				{"#4169E1", "#1E90FF"}, // Royal Blue
				{"#FF1493", "#FF69B4"}, // Deep Pink
				{"#00FA9A", "#00FF7F"}, // Spring Green
				{"#FF8C00", "#FF6347"}, // Dark Orange
				{"#8A2BE2", "#9400D3"}, // Blue Violet
				{"#DC143C", "#FF0000"}, // Crimson
				{"#00BFFF", "#87CEEB"}, // Deep Sky Blue
				{"#ADFF2F", "#FFFF00"}, // Green Yellow
				{"#FF00FF", "#DA70D6"}, // Magenta/Orchid
			}

			idx := 0
			for i, t := range themes {
				if t.border == string(m.borderColor) {
					idx = i
					break
				}
			}
			idx = (idx + 1) % len(themes)
			m.borderColor = lipgloss.Color(themes[idx].border)
			m.highlightColor = lipgloss.Color(themes[idx].highlight)

			// Update delegate with new colors
			d := list.NewDefaultDelegate()
			d.Styles.SelectedTitle = lipgloss.NewStyle().
				Foreground(m.highlightColor).
				Bold(true)
			d.Styles.SelectedDesc = lipgloss.NewStyle().
				Foreground(m.highlightColor)
			m.list.SetDelegate(d)

			// Save config
			saveConfig(Config{
				BorderColor:    string(m.borderColor),
				HighlightColor: string(m.highlightColor),
			})
			return m, nil
		
		case "o":
			// Open location in file explorer
			if !m.driveSelect {
				if i, ok := m.list.SelectedItem().(item); ok {
					// Open parent directory in file explorer
					parentDir := filepath.Dir(i.Path)
					if i.IsDir {
						// If it's a directory, open the directory itself
						parentDir = i.Path
					}
					switch runtime.GOOS {
					case "windows":
						exec.Command("explorer", parentDir).Start()
					case "darwin":
						exec.Command("open", parentDir).Start()
					default:
						exec.Command("xdg-open", parentDir).Start()
					}
				}
			}
			return m, nil

		}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func renderExtendedHelp() string {
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Padding(1, 0)
	
	helpText := `Extended Help: ( ? to close )
  ↑/↓ or j/k    - Navigate up/down
  enter         - Open file or enter folder
  o             - Open location in file explorer
  backspace     - Go back to parent folder
  /             - Filter/search items
  esc           - Clear filter
  t             - Change color theme
  q or ctrl+c   - Quit`
	
	return helpStyle.Render(helpText)
}

func (m model) View() string {
	// Small retro ASCII art title
	titleArt := "░█▄▀░▀█▀░▀█▀░▀█▀░█░█░█▀▀░█▀▀\n░█░█░░█░░░█░░░█░░░█░░█▀▀░▀▀█\n░▀░▀░▀▀▀░░▀░░░▀░░░▀░░▀░░░▀▀▀"

	titleStyle := lipgloss.NewStyle().
		Foreground(m.highlightColor).
		Bold(true).
		Align(lipgloss.Center)

	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.borderColor).
		Padding(m.paddingY, m.paddingX)

	// Create content with title inside border
	innerContent := titleStyle.Render(titleArt) + "\n\n" + m.list.View()

	// Add extended help inside border if toggled
	if m.showExtendedHelp {
		innerContent += "\n" + renderExtendedHelp()
	}

	content := border.Render(innerContent)

	// Center the content
	contentWidth := lipgloss.Width(content)
	leftPadding := (m.width - contentWidth) / 2
	if leftPadding < 0 {
		leftPadding = 0
	}

	centeredContent := lipgloss.NewStyle().
		PaddingLeft(leftPadding).
		Render(content)

	return centeredContent
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}