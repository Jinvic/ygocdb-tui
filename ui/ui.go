package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"ygocdb-tui/api"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)
	
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)
			
	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(0, 1)
			
	resultStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2)
			
	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2)
			
	paginationStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Padding(1, 0)
			
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).Render
)

type model struct {
	textInput   textinput.Model
	results     []api.Card
	card        *api.GetCardResponse
	selected    int
	err         error
	mode        mode
	loading     bool
	apiClient   *api.Client
	query       string
	start       int
	next        int
	pageHistory []int // 记录页面历史，用于正确返回上一页
}

type mode int

const (
	searchMode mode = iota
	resultMode
	cardMode
)

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "输入卡片名称或ID"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40

	return model{
		textInput:   ti,
		results:     []api.Card{},
		card:        nil,
		selected:    -1,
		err:         nil,
		mode:        searchMode,
		loading:     false,
		apiClient:   api.NewClient(),
		query:       "",
		start:       0,
		next:        0,
		pageHistory: make([]int, 0),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// Search cards command
func searchCards(query string, start int) tea.Cmd {
	return func() tea.Msg {
		// Check if query is a number (card ID) and we're on the first page
		if cardID, err := strconv.Atoi(query); err == nil && start == 0 {
			// Query by card ID (only for first page)
			client := api.NewClient()
			card, err := client.GetCardByID(cardID)
			if err != nil {
				return searchErrorMsg{err}
			}
			return searchByIDResultMsg{card}
		}
		
		// Search by name with pagination
		client := api.NewClient()
		results, err := client.SearchCards(query, start)
		if err != nil {
			return searchErrorMsg{err}
		}
		return searchResultMsg{results, query, start}
	}
}

// Get card by ID command
func getCardByID(id int) tea.Cmd {
	return func() tea.Msg {
		client := api.NewClient()
		card, err := client.GetCardByID(id)
		if err != nil {
			return searchErrorMsg{err}
		}
		return cardResultMsg{card}
	}
}

// Messages
type searchResultMsg struct {
	results *api.SearchResponse
	query   string
	start   int
}

type searchByIDResultMsg struct {
	card *api.GetCardResponse
}

type cardResultMsg struct {
	card *api.GetCardResponse
}

type searchErrorMsg struct {
	err error
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			if m.mode == searchMode {
				return m, tea.Quit
			} else if m.mode == resultMode || m.mode == cardMode {
				// Go back to search mode
				m.mode = searchMode
				m.results = []api.Card{}
				m.card = nil
				m.selected = -1
				m.textInput.Focus()
				m.start = 0
				m.next = 0
				m.pageHistory = make([]int, 0) // 清空页面历史
				return m, nil
			}

		case tea.KeyEnter:
			if m.mode == searchMode && !m.loading {
				// Search for cards
				query := m.textInput.Value()
				if query != "" {
					m.query = query
					m.start = 0
					m.pageHistory = make([]int, 0) // 开始新搜索时清空页面历史
					m.loading = true
					m.textInput.Blur()
					return m, searchCards(query, 0)
				}
			} else if m.mode == resultMode && len(m.results) > 0 {
				// View selected card
				if m.selected >= 0 && m.selected < len(m.results) {
					m.loading = true
					return m, getCardByID(m.results[m.selected].ID)
				}
			} else if m.mode == cardMode {
				// Back to results
				m.mode = resultMode
				m.card = nil
				m.loading = false
				return m, nil
			}

		case tea.KeyUp:
			if m.mode == resultMode && len(m.results) > 0 {
				m.selected--
				if m.selected < 0 {
					m.selected = len(m.results) - 1
				}
			}
			return m, nil

		case tea.KeyDown:
			if m.mode == resultMode && len(m.results) > 0 {
				m.selected++
				if m.selected >= len(m.results) {
					m.selected = 0
				}
			}
			return m, nil

		case tea.KeyRight:
			// Next page
			if m.mode == resultMode && m.next > 0 && !m.loading {
				m.pageHistory = append(m.pageHistory, m.start) // 记录当前页到历史
				m.loading = true
				m.start = m.next
				return m, searchCards(m.query, m.start)
			}
			return m, nil

		case tea.KeyLeft:
			// Previous page
			if m.mode == resultMode && len(m.pageHistory) > 0 && !m.loading {
				m.loading = true
				// 从历史记录中获取上一页的start位置
				prevStart := m.pageHistory[len(m.pageHistory)-1]
				m.pageHistory = m.pageHistory[:len(m.pageHistory)-1] // 移除最后一条记录
				m.start = prevStart
				return m, searchCards(m.query, m.start)
			}
			return m, nil
		}

	case searchResultMsg:
		m.loading = false
		m.mode = resultMode
		m.results = msg.results.Result
		m.selected = 0
		m.next = msg.results.Next
		if len(m.results) == 0 {
			m.err = fmt.Errorf("未找到相关卡片")
		}
		return m, nil

	case searchByIDResultMsg:
		m.loading = false
		m.mode = cardMode
		m.card = msg.card
		return m, nil

	case cardResultMsg:
		m.loading = false
		m.mode = cardMode
		m.card = msg.card
		return m, nil

	case searchErrorMsg:
		m.loading = false
		m.err = msg.err
		m.textInput.Focus()
		return m, nil

	// Handle input changes
	case tea.WindowSizeMsg:
		// Handle window resizing if needed
	}

	// Update text input
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var b strings.Builder
	
	switch m.mode {
	case searchMode:
		b.WriteString(titleStyle.Render("游戏王卡片查询工具 (百鸽API)"))
		b.WriteString("\n\n")
		b.WriteString(inputStyle.Render(m.textInput.View()))
		b.WriteString("\n\n")
		
		if m.loading {
			b.WriteString("搜索中...")
		} else if m.err != nil {
			b.WriteString(fmt.Sprintf("错误: %v\n\n", m.err))
			m.err = nil // Reset error after displaying
		}
		
		b.WriteString(helpStyle("按 Enter 搜索，按 Esc 退出"))
		
	case resultMode:
		b.WriteString(titleStyle.Render("搜索结果"))
		b.WriteString("\n\n")
		
		if m.loading {
			b.WriteString("加载中...")
		} else if len(m.results) == 0 {
			b.WriteString("未找到相关卡片")
		} else {
			for i, result := range m.results {
				if i == m.selected {
					b.WriteString("> " + resultStyle.Render(formatCardSummary(result)) + "\n\n")
				} else {
					b.WriteString("  " + formatCardSummary(result) + "\n\n")
				}
			}
			
			// Pagination info
			b.WriteString("\n")
			b.WriteString(paginationStyle.Render(formatPagination(m.start, m.next, len(m.pageHistory))))
		}
		
		b.WriteString("\n")
		b.WriteString(helpStyle("使用 ↑/↓ 选择卡片，←/→ 翻页，按 Enter 查看详情，按 Esc 返回"))

	case cardMode:
		b.WriteString(titleStyle.Render("卡片详情"))
		b.WriteString("\n\n")
		
		if m.loading {
			b.WriteString("加载中...")
		} else if m.card != nil {
			b.WriteString(cardStyle.Render(formatCardDetails(*m.card)))
		}
		
		b.WriteString("\n\n")
		b.WriteString(helpStyle("按 Enter 或 Esc 返回搜索结果"))
	}

	return appStyle.Render(b.String())
}

func formatCardSummary(card api.Card) string {
	return fmt.Sprintf("%s (%d)", card.CnName, card.ID)
}

func formatPagination(start, next int, historyLen int) string {
	page := start/10 + 1
	if start == 0 && next == 0 {
		return fmt.Sprintf("第 %d 页", page)
	} else if start == 0 {
		return fmt.Sprintf("第 %d 页 | 按 → 查看下一页", page)
	} else if next == 0 {
		return fmt.Sprintf("第 %d 页 | 按 ← 查看上一页", page)
	} else {
		return fmt.Sprintf("第 %d 页 | 按 ← 查看上一页 | 按 → 查看下一页", page)
	}
}

func formatCardDetails(card api.GetCardResponse) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("卡片密码: %d\n", card.ID))
	b.WriteString(fmt.Sprintf("名称: %s\n", card.Text.Name))
	b.WriteString(fmt.Sprintf("类型: %s\n", getCardType(card.Data.Type)))
	b.WriteString(fmt.Sprintf("种族: %s\n", getCardRace(card.Data.Race)))
	b.WriteString(fmt.Sprintf("属性: %s\n", getCardAttribute(card.Data.Attrib)))
	
	if card.Data.Level > 0 {
		b.WriteString(fmt.Sprintf("星级: %d\n", card.Data.Level))
	}
	
	if card.Data.Atk >= 0 {
		b.WriteString(fmt.Sprintf("攻击力: %d\n", card.Data.Atk))
	}
	
	if card.Data.Def >= 0 {
		b.WriteString(fmt.Sprintf("守备力: %d\n", card.Data.Def))
	}
	
	b.WriteString(fmt.Sprintf("\n效果:\n%s", card.Text.Desc))
	
	return b.String()
}

func getCardType(typ int) string {
	switch typ {
	case 17:
		return "怪兽|通常"
	case 33:
		return "怪兽|效果"
	case 65:
		return "魔法|通常"
	case 129:
		return "陷阱|通常"
	default:
		return fmt.Sprintf("未知类型(%d)", typ)
	}
}

func getCardRace(race int) string {
	switch race {
	case 1:
		return "战士"
	case 2:
		return "魔法师"
	case 4:
		return "天使"
	case 8:
		return "恶魔"
	case 16:
		return "不死"
	case 32:
		return "机械"
	case 64:
		return "水"
	case 128:
		return "炎"
	case 256:
		return "岩石"
	case 512:
		return "鸟兽"
	case 1024:
		return "植物"
	case 2048:
		return "昆虫"
	case 4096:
		return "雷"
	case 8192:
		return "龙"
	case 16384:
		return "兽"
	case 32768:
		return "兽战士"
	case 65536:
		return "恐龙"
	case 131072:
		return "鱼"
	case 262144:
		return "海龙"
	case 524288:
		return "爬虫类"
	case 1048576:
		return "念动力"
	case 2097152:
		return "幻神兽"
	case 4194304:
		return "创造神"
	case 8388608:
		return "幻龙"
	default:
		return fmt.Sprintf("未知种族(%d)", race)
	}
}

func getCardAttribute(attr int) string {
	switch attr {
	case 1:
		return "地"
	case 2:
		return "水"
	case 4:
		return "炎"
	case 8:
		return "风"
	case 16:
		return "光"
	case 32:
		return "暗"
	case 64:
		return "神"
	default:
		return fmt.Sprintf("未知属性(%d)", attr)
	}
}

// Start initializes and starts the TUI application
func Start() error {
	p := tea.NewProgram(initialModel())
	_, err := p.Run()
	return err
}