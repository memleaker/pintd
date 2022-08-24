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
