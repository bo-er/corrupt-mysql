package backend

import (
	"sync"

	"github.com/bo-er/corrupt-mysql/pkg"
)

const initSQL = `
DROP DATABASE IF EXISTS corrupt_mysql_deadlock_test;
CREATE DATABASE corrupt_mysql_deadlock_test;
USE corrupt_mysql_deadlock_test;
DROP TABLE IF EXISTS people;
DROP TABLE IF EXISTS students;
CREATE TABLE people (
    id             serial PRIMARY KEY,
    age integer
  );
CREATE TABLE students (
    id      serial PRIMARY KEY,
    age     integer
  );
  INSERT INTO people (id, age) VALUES (1, 8);
  INSERT INTO students (id, age) VALUES (1, 8);
`

const Thread1 = `
BEGIN;
USE corrupt_mysql_deadlock_test;
UPDATE people SET age =  18 WHERE id = 1;
UPDATE students SET age =  18 WHERE id = 1;
`

const Thread2 = `
BEGIN;
USE corrupt_mysql_deadlock_test;
UPDATE students SET age =  18 WHERE id = 1;
UPDATE people SET age =  18 WHERE id = 1;
`

func CreateDeadlock(c pkg.Connect) error {
	db, err := pkg.GetDB(c)
	if err != nil {
		return err
	}
	err = pkg.BatchExec(db, initSQL)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(2)
	for _, sql := range []string{Thread1, Thread2} {
		sql := sql
		go func() {
			err = pkg.BatchExec(db, sql)
			if err != nil {
				panic(err.Error())
			}
		}()
	}
	wg.Wait()
	return nil
}
