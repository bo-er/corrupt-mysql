package backend

import (
	"fmt"

	"github.com/bo-er/corrupt-mysql/pkg"
	"github.com/sirupsen/logrus"
)

const enablePerformanceSchemaSQL = `
UPDATE performance_schema.setup_instruments
SET ENABLED='YES', TIMED='YES'
WHERE NAME LIKE 'stage/innodb/alter%';

UPDATE performance_schema.setup_consumers
SET ENABLED='YES'
WHERE NAME LIKE '%stage%';
`

const createBigDDLSQL = `
CREATE DATABASE IF NOT EXISTS TESTDB;
USE TESTDB;

CREATE TABLE IF NOT EXISTS  t1(x int primary key auto_increment);
INSERT INTO t1() VALUES (),(),(),();
INSERT INTO t1(x) SELECT x+(SELECT COUNT(*) FROM  t1) FROM t1;
`

// why using this COPY algorithm?
// see https://dev.mysql.com/blog-archive/mysql-8-0-innodb-now-supports-instant-add-column/
const triggerBigDDLSQL = `
ALTER TABLE testdb.t1 ADD COLUMN col1 INT NOT NULL, ALGORITHM=COPY;
`

func CreateBigDDL(c pkg.Connect) error {
	db, err := pkg.GetDB(c)
	if err != nil {
		return err
	}

	// Execute preparatory SQL to enable performance schema monitoring
	err = pkg.BatchExec(db, enablePerformanceSchemaSQL)
	if err != nil {
		return fmt.Errorf("error enabling performance schema: %w", err)
	}
	logrus.Info("Performance schema monitoring enabled")

	err = pkg.BatchExec(db, createBigDDLSQL)
	if err != nil {
		return err
	}
	defer func() {
		// Drop the TESTDB database
		_, err = db.Exec("DROP DATABASE TESTDB")
		if err != nil {
			logrus.Errorf("error dropping database: %v", err)
		}
	}()

	// Execute the doubling logic in a loop
	for i := 0; i < 23; i++ {
		logrus.Infof("inserting data, round %d/%d", i, 23)
		err = pkg.Exec(db, "INSERT INTO t1(x) SELECT x + (SELECT COUNT(*) FROM t1) FROM t1")
		if err != nil {
			return fmt.Errorf("error doubling data: %w", err)
		}
	}

	// Trigger the big DDL
	fmt.Println("Triggering big DDL...")
	_, err = db.Exec(triggerBigDDLSQL)
	if err != nil {
		return fmt.Errorf("error triggering big DDL: %w", err)
	}

	return nil
}
