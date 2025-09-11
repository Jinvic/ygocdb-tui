package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	"ygocdb-tui/internal/api"
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