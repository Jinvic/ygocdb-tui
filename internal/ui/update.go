package ui

import (
	"fmt"
	"ygocdb-tui/internal/api"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
)

func (m model) Init() tea.Cmd {
	return textinput.Blink
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