package Ticker

type Ticker struct {
	Results         []Result `json:"results"`
	Success         bool     `json:"success"`
	Map             Maps     `json:"map"`
	Name            string   `json:"ticker"`
	ResultCount     int64    `json:"results_count"`
	DatabaseLatency int64    `json:"db_latency"`
}

type keyMap struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Maps struct {
	Tape        keyMap `json:"z"`
	Conditions  keyMap `json:"c"`
	Exchange    keyMap `json:"x"`
	Trf         keyMap `json:"f"`
	Sequence    keyMap `json:"q"`
	ID          keyMap `json:"i"`
	OrignalID   keyMap `json:"I"`
	Correction  keyMap `json:"e"`
	TRFID       keyMap `json:"r"`
	SIP         keyMap `json:"t"`
	Participant keyMap `json:"y"`
	Size        keyMap `json:"s"`
	Price       keyMap `json:"p"`
}

type Result struct {
	SIP         int64   `json:"t"`
	Participant int64   `json:"y"`
	Sequence    int     `json:"q"`
	ID          string  `json:"i"`
	Exchange    int     `json:"x"`
	Size        int     `json:"s"`
	Conditions  []int   `json:"c"`
	Price       float64 `json:"p"`
	Tape        int     `json:"z"`
}
