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

func NewIndirectCfg(cfg *IndirectConfig) error {
	sql := `INSERT INTO indirect_cfg (protocol, "listen-addr", "listen-port",
			"dest-addr", "dest-port", acl, "deny-addr", "admit-addr", 
			"max-conns", memo) values(?,?,?,?,?,?,?,?,?,?)`

	stmt, err := db.Prepare(sql)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(cfg.Protocol, cfg.ListenAddr, cfg.ListenPort,
		cfg.DestAddr, cfg.DestPort, cfg.Acl, cfg.DenyAddr,
		cfg.AdmitAddr, cfg.MaxConns, cfg.Memo)
	if err != nil {
		return err
	}

	return nil
}

func RepeatConfig(cfg *IndirectConfig) (bool, error) {
	sql := `SELECT * FROM indirect_cfg WHERE protocol = ? AND "listen-port"=? LIMIT 1`

	// just need check protocol and listen-port.
	stmt, err := db.Prepare(sql)
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(cfg.Protocol, cfg.ListenPort)
	if err != nil {
		return false, err
	}

	defer rows.Close()

	if rows.Next() {
		return true, nil
	}

	return false, nil
}

func GetIndirectCfg(page, limit int) ([]*IndirectConfig, error) {
	cfgs := make([]*IndirectConfig, 0)
	sql := `SELECT * FROM indirect_cfg ORDER BY id DESC LIMIT ? OFFSET ?`

	stmt, err := db.Prepare(sql)
	if err != nil {
		return cfgs, err
	}

	defer stmt.Close()

	offset := (page - 1) * limit
	rows, err := stmt.Query(limit, offset)
	if err != nil {
		return cfgs, err
	}

	defer rows.Close()

	for rows.Next() {
		id := 0
		cfg := IndirectConfig{}
		err := rows.Scan(&id, &cfg.Protocol, &cfg.ListenAddr, &cfg.ListenPort, &cfg.DestAddr,
			&cfg.DestPort, &cfg.Acl, &cfg.DenyAddr, &cfg.AdmitAddr,
			&cfg.MaxConns, &cfg.Memo)
		if err != nil {
			return cfgs, err
		}
		cfgs = append(cfgs, &cfg)
	}

	return cfgs, nil
}

func GetIndirectTblRows() (int, error) {
	rows := 0
	sql := `SELECT count(*) FROM indirect_cfg`

	stmt, err := db.Prepare(sql)
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	err = stmt.QueryRow().Scan(&rows)
	if err != nil {
		return 0, err
	}

	return rows, nil
}

func DelIndirectCfg(protocol, listen_port *string) (bool, error) {
	sql := `DELETE FROM indirect_cfg WHERE protocol = ? AND "listen-port"=?`

	stmt, err := db.Prepare(sql)
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	res, err := stmt.Exec(*protocol, *listen_port)
	if err != nil {
		return false, err
	}

	r, err := res.RowsAffected()
	if err != nil || r <= 0 {
		return false, err
	}

	return true, nil
}

func UpdateIndirectCfg(field, val *string, cfg *IndirectConfig) (bool, error) {
	sql := `UPDATE indirect_cfg SET "` + *field + `"=? WHERE protocol=? AND "listen-port" = ?`

	stmt, err := db.Prepare(sql)
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	res, err := stmt.Exec(*val, cfg.Protocol, cfg.ListenPort)
	if err != nil {
		return false, err
	}

	r, err := res.RowsAffected()
	if r <= 0 && err != nil {
		return false, err
	}

	return true, nil
}
