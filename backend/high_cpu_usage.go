package backend

import (
	"fmt"

	"github.com/bo-er/corrupt-mysql/pkg"
)

const prepareHighCPUUsageSQL = `
DROP DATABASE IF EXISTS corrupt_mysql_highcpuusage_test;
CREATE DATABASE corrupt_mysql_highcpuusage_test;
USE corrupt_mysql_highcpuusage_test;
DROP TABLE IF EXISTS t;
CREATE TABLE t(x int primary key auto_increment);
INSERT INTO t() values(),(),(),();
`

func MakeCPUUsageHigh(c pkg.Connect) error {
	db, err := pkg.GetDB(c)
	if err != nil {
		return err
	}
	err = pkg.BatchExec(db, prepareBigTransactionSQL)
	if err != nil {
		return err
	}
	// maybe there is a better way to do this.
	_, err = db.Exec(`USE corrupt_mysql_highcpuusage_test;`)
	if err != nil {
		return fmt.Errorf("use corrupt_mysql_highcpuusage_test: %s", err.Error())
	}
	// preparing some data
	for i := 1; i <= 20; i++ {
		_, err = db.Exec(quadraticGrowthSQL)
		if err != nil {
			return fmt.Errorf("calling(%s) for the %dth time: %s", quadraticGrowthSQL, i, err.Error())
		}
	}
	// calling rand() on each entry.
	_, err = db.Exec(`SELECT * FROM t order by rand() limit 1;`)
	if err != nil {
		return fmt.Errorf("calling rand() on each entry: %s", err.Error())
	}
	return nil
}
