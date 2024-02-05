package example

import (
	"fmt"
	"snowflake/core"
	"snowflake/log"

	_ "github.com/go-sql-driver/mysql"
)

// 插入空投
func (s *ExampleServer) SaveAirdrop(data *Airdrop) bool {
	db := core.OpenMysql(s.ctx, s.module)
	if db == nil {
		return false
	}
	defer db.Close()

	query := `INSERT INTO airdrop(symbol, address, share, amount, hashcode, timestamp)VALUES(?,?,?,?,?,?)
			  ON DUPLICATE KEY UPDATE symbol=VALUES(symbol),address=VALUES(address),share=VALUES(share),amount=VALUES(amount),timestamp=VALUES(timestamp);`
	_, err := db.Exec(query, data.Symbol, data.Address, data.Share.String(), data.Amount.String(), data.Hashcode, data.Timestamp)
	if err != nil {
		s.Error("failed to save airdrop: %s -> %s", err.Error(), query)
		return false
	}

	return true
}

// 清除数据
func (s *ExampleServer) Truncate(table string) {
	db := core.OpenMysql(s.ctx, s.module)
	if db == nil {
		return
	}
	defer db.Close()

	query := fmt.Sprintf("Truncate %s", table)
	_, err := db.Exec(query)
	if err != nil {
		s.Error("failed to truncate categories: %s -> %s", err.Error(), query)
		return
	}
}

// 获取空投
func (s *ExampleServer) GetAllAirdrops() []*Airdrop {
	db := core.OpenMysql(s.ctx, s.module)
	if db == nil {
		return nil
	}
	defer db.Close()

	query := "SELECT * FROM airdrop where amount>0;"
	rows, err := db.Query(query)
	if err != nil {
		s.Error("failed to query airdrop: %s -> %s", err.Error(), query)
	}

	airdrops := make([]*Airdrop, 0, 0)
	for rows.Next() {
		var item Airdrop
		err := rows.Scan(&item.Symbol, &item.Address, &item.Share, &item.Amount, &item.Hashcode, &item.Timestamp)
		if err != nil {
			log.Errorf("failed to execute airdrop method: %s", err.Error())
			continue
		}
		airdrops = append(airdrops, &item)
	}

	return airdrops
}
