package pkg

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
)

type Connect struct {
	User, Host, Password, DBName string
	Port                         int
}

func GetDB(c Connect) (*sql.DB, error) {
	fmt.Println(c)
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	cfg := mysql.Config{
		User:   c.User,
		Passwd: c.Password,
		Net:    "tcp",
		Addr:   addr,
		//DBName:               c.DBName,
		AllowNativePasswords: true,
	}

	var (
		db  *sql.DB
		err error
	)

	// Get a database handle.
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("open mysql(%s): %s", addr, err.Error())
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, fmt.Errorf("ping mysql(%s): %s", addr, pingErr.Error())
	}
	return db, nil
}

func BatchExec(db *sql.DB, sqls string) error {
	for _, sql := range strings.Split(strings.TrimSpace(sqls), ";") {
		sql := sql
		if strings.TrimSpace(sql) == "" {
			continue
		}
		_, err := db.Exec(sql)
		if err != nil {
			return fmt.Errorf("sql: %s,initializing databases: %s", sql, err.Error())
		}
		fmt.Println(sql)
	}
	return nil
}
