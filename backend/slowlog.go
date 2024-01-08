package backend

import (
	"github.com/bo-er/corrupt-mysql/pkg"
)

const slowlogSQL = `
DROP DATABASE IF EXISTS corrupt_mysql_deadlock_test;
CREATE DATABASE corrupt_mysql_deadlock_test;
USE corrupt_mysql_deadlock_test;
DROP TABLE IF EXISTS people;
CREATE TABLE people (
    id             serial PRIMARY KEY,
    age integer
  );
  INSERT INTO people (id, age) VALUES (1, 8);
SET @@SESSION.long_query_time = 1;
SELECT sleep(2),id FROM people;
`

func CreateSlowlog(c pkg.Connect) error {
	db, err := pkg.GetDB(c)
	if err != nil {
		return err
	}
	return pkg.BatchExec(db, slowlogSQL)
}
