package ui

import (
	"fmt"
	"strings"
	"ygocdb-tui/internal/api"
)

// formatCardSummary formats a card summary for display
func formatCardSummary(card api.Card) string {
	return fmt.Sprintf("%s (%d)", card.CnName, card.ID)
}

// formatPagination formats pagination information for display
func formatPagination(currentPage, totalPages int) string {
	page := currentPage + 1 // Convert to 1-based indexing for display
	if totalPages <= 1 {
		return fmt.Sprintf("第 %d 页", page)
	} else if currentPage == 0 {
		return fmt.Sprintf("第 %d 页 | 按 → 查看下一页", page)
	} else if currentPage == totalPages-1 {
		return fmt.Sprintf("第 %d 页 | 按 ← 查看上一页", page)
	} else {
		return fmt.Sprintf("第 %d 页 | 按 ← 查看上一页 | 按 → 查看下一页", page)
	}
}

// formatCardDetails formats card details for display
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

// getCardType returns the Chinese name for a card type
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

// getCardRace returns the Chinese name for a card race
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

// getCardAttribute returns the Chinese name for a card attribute
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
