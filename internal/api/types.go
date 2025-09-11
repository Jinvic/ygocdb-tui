package api

// Card represents a Yu-Gi-Oh! card
type Card struct {
	CID    int    `json:"cid"`
	ID     int    `json:"id"`
	CnName string `json:"cn_name"`
	ScName string `json:"sc_name"`
	MdName string `json:"md_name"`
	NwbbsN string `json:"nwbbs_n"`
	CnocgN string `json:"cnocg_n"`
	JpRuby string `json:"jp_ruby"`
	JpName string `json:"jp_name"`
	EnName string `json:"en_name"`
	Text   Text   `json:"text"`
	Data   Data   `json:"data"`
}

// Text represents card text information
type Text struct {
	Name  string `json:"name"`
	Types string `json:"types"`
	PDesc string `json:"pdesc"`
	Desc  string `json:"desc"`
}

// Data represents card data information
type Data struct {
	OT      int `json:"ot"`
	Setcode int `json:"setcode"`
	Type    int `json:"type"`
	Atk     int `json:"atk"`
	Def     int `json:"def"`
	Level   int `json:"level"`
	Race    int `json:"race"`
	Attrib  int `json:"attribute"`
}

// SearchResponse represents the response from search API
type SearchResponse struct {
	Result []Card `json:"result"`
	Next   int    `json:"next"`
}

// GetCardResponse represents the response from get card API
type GetCardResponse struct {
	ID   int  `json:"id"`
	Data Data `json:"data"`
	Text Text `json:"text"`
}