package Ticker

type Ticker struct {
	Name            string    `json:"ticker"`
	ResultCount     int64     `json:"results_count"`
	DatabaseLatency int64     `json:"db_latency"`
	Success         bool      `json:"success"`
	Results         []Result  `json:"results"`
	Map             ResultMap `json:"map"`
}

type keyMap struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type ResultMap struct {
	Map Maps `json:"map"`
}

type Maps struct {
	OrignalID   keyMap `json:"I"`
	Exchange    keyMap `json:"x"`
	Price       keyMap `json:"p"`
	ID          keyMap `json:"i"`
	Correction  keyMap `json:"e"`
	TRFID       keyMap `json:"r"`
	SIP         keyMap `json:"t"`
	Participant keyMap `json:"y"`
	Trf         keyMap `json:"f"`
	Sequence    keyMap `json:"q"`
	Conditions  keyMap `json:"c"`
	Size        keyMap `json:"s"`
	Tape        keyMap `json:"z"`
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
