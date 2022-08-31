package model

type Logging struct {
	Id      string `json:"id"`
	Time    string `json:"time"`
	Content string `json:"content"`
}

func NewLog(log string) error {
	sql := `INSERT INTO log (content) values(?)`

	stmt, err := db.Prepare(sql)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(log)
	if err != nil {
		return err
	}

	return nil
}

func GetLogTblRows() (int, error) {
	rows := 0
	sql := `SELECT count(*) FROM log`

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

func GetLog(page, limit int) ([]*Logging, error) {
	logs := make([]*Logging, 0)
	sql := `SELECT * FROM log ORDER BY id DESC LIMIT ? OFFSET ?`

	stmt, err := db.Prepare(sql)
	if err != nil {
		return logs, err
	}

	defer stmt.Close()

	offset := (page - 1) * limit
	rows, err := stmt.Query(limit, offset)
	if err != nil {
		return logs, err
	}

	defer rows.Close()

	for rows.Next() {
		log := Logging{}
		err := rows.Scan(&log.Id, &log.Time, &log.Content)
		if err != nil {
			return logs, err
		}
		logs = append(logs, &log)
	}

	return logs, nil
}

func DelLog(id string) (bool, error) {
	sql := `DELETE FROM log WHERE id=?`

	stmt, err := db.Prepare(sql)
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return false, err
	}

	r, err := res.RowsAffected()
	if err != nil || r <= 0 {
		return false, err
	}

	return true, nil
}

func DelMoreLog(logs []Logging) (bool, error) {
	for i := 0; i < len(logs); i++ {
		ok, err := DelLog(logs[i].Id)
		if !ok || err != nil {
			return ok, err
		}
	}

	return true, nil
}
