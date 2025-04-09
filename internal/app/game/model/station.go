package model

type StationType int

const (
	StationTypeNormal StationType = iota
	StationTypeTransfer
)

type LineType int

const (
	Line1 LineType = iota + 1
	Line2
	Line3
	Line4
	Line5
	Line19
)

type CurrentStation struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Station struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Type           StationType    `json:"type"`
	Price          int            `json:"price"`
	CurrentStation CurrentStation `json:"CurrentStation"`
	Owner          *Team          `json:"owner"`
	ConnectedTo    []*Station     `json:"connectedTo"`
	Line           LineType       `json:"line"`
	Next           *Station       `json:"-"` // 同线路的下一个站
	Prev           *Station       `json:"-"` // 同线路的上一个站
}

type MetroLine struct {
	Type     LineType   `json:"type"`
	Name     string     `json:"name"`
	Stations []*Station `json:"stations"`
}
