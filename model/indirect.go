package model

// Indirect.

const (
	AclWhiteList = "whitelist"
	AclBlackList = "blacklist"
	AclDisable   = "disabled"
)

type IndirectConfig struct {
	Protocol   string `json:"protocol"`
	ListenAddr string `json:"listen-addr"`
	ListenPort string `json:"listen-port"`
	DestAddr   string `json:"dest-addr"`
	DestPort   string `json:"dest-port"`
	Acl        string `json:"acl"`
	DenyAddr   string `json:"deny-addr"`
	AdmitAddr  string `json:"admit-addr"`
	MaxConns   string `json:"max-conns"`
	Memo       string `json:"memo"`
}

type IndirectState struct {
	Protocol     string `json:"protocol"`
	SrcAddr      string `json:"src-addr"`
	SrcPort      string `json:"src-port"`
	DestAddr     string `json:"dest-addr"`
	DestPort     string `json:"dest-port"`
	RunningTime  string `json:"running-time"`
	ForwardFlow  string `json:"forward-flow"`
	RealTimeFlow string `json:"realtime-flow"`
}
