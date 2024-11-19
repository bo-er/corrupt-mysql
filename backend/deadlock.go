package backend

import (
	"fmt"
	"sync"

	"github.com/bo-er/corrupt-mysql/pkg"
	"github.com/sirupsen/logrus"
)

const prepareDeadlockSQL = `
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

	// Get the original innodb_lock_wait_timeout
	var (
		variableName    string
		originalTimeout int
	)
	err = db.QueryRow("SHOW VARIABLES LIKE 'innodb_lock_wait_timeout'").Scan(&variableName, &originalTimeout)
	if err != nil {
		return fmt.Errorf("error getting original timeout: %w", err)
	}
	logrus.Infof("original timeout is %d", originalTimeout)
	newTimeout := 60
	err = pkg.Exec(db, "SET GLOBAL innodb_lock_wait_timeout = ?", newTimeout)
	if err != nil {
		return fmt.Errorf("error setting new timeout: %w", err)
	}
	defer func() {
		// Restore the original innodb_lock_wait_timeout
		err = pkg.Exec(db, "SET GLOBAL innodb_lock_wait_timeout = ?", originalTimeout)
		if err != nil {
			logrus.Errorf("error restoring original timeout: %v", err)
		}
	}()

	err = pkg.BatchExec(db, prepareDeadlockSQL)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(2)
	for _, sql := range []string{Thread1, Thread2} {
		sql := sql
		go func() {
			defer wg.Done()
			err = pkg.BatchExec(db, sql)
			if err != nil {
				// Handle the deadlock error (e.g., log it)
				logrus.Errorf("Error: %v", err)
			}
		}()
	}
	wg.Wait()

	return nil
}
