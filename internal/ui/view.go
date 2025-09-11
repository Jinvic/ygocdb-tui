package ui

import (
	"fmt"
	"strings"
)

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