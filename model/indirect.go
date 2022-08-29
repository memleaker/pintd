package model

import (
	"pintd/plog"
)

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

func NewIndirectCfg(cfg *IndirectConfig) bool {
	stmt, err := db.Prepare(`INSERT INTO indirect_cfg (protocol, "listen-addr", "listen-port",
		"dest-addr", "dest-port", acl, "deny-addr", "admit-addr", 
		"max-conns", memo) values(?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		plog.Println("Insert into indirect_cfg ", err.Error())
		return false
	}

	defer stmt.Close()

	_, err = stmt.Exec(cfg.Protocol, cfg.ListenAddr, cfg.ListenPort,
		cfg.DestAddr, cfg.DestPort, cfg.Acl, cfg.DenyAddr,
		cfg.AdmitAddr, cfg.MaxConns, cfg.Memo)
	if err != nil {
		plog.Println("Insert into indirect_cfg ", err.Error())
		return false
	}

	return true
}

func RepeatConfig(cfg *IndirectConfig) bool {

	// just need check protocol and listen-port.
	stmt, err := db.Prepare(`SELECT * FROM indirect_cfg WHERE protocol = ? AND "listen-port"=? LIMIT 1`)
	if err != nil {
		plog.Println("Select from indirect_cfg ", err.Error())
		return false
	}

	defer stmt.Close()

	rows, err := stmt.Query(cfg.Protocol, cfg.ListenPort)
	if err != nil {
		plog.Println("Select from indirect_cfg ", err.Error())
		return false
	}

	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}

func GetIndirectCfg(page, limit int) []*IndirectConfig {
	cfgs := make([]*IndirectConfig, 0)

	stmt, err := db.Prepare(`SELECT * FROM indirect_cfg ORDER BY id DESC LIMIT ? OFFSET ?`)
	if err != nil {
		plog.Println("Select from indirect_cfg ", err.Error())
		return cfgs
	}

	defer stmt.Close()

	offset := (page - 1) * limit
	rows, err := stmt.Query(limit, offset)
	if err != nil {
		plog.Println("Select from indirect_cfg ", err.Error())
		return cfgs
	}

	for rows.Next() {
		id := 0
		cfg := IndirectConfig{}
		err := rows.Scan(&id, &cfg.Protocol, &cfg.ListenAddr, &cfg.ListenPort, &cfg.DestAddr,
			&cfg.DestPort, &cfg.Acl, &cfg.DenyAddr, &cfg.AdmitAddr,
			&cfg.MaxConns, &cfg.Memo)
		if err != nil {
			plog.Println(err.Error())
		}
		cfgs = append(cfgs, &cfg)
	}

	rows.Close()
	return cfgs
}
