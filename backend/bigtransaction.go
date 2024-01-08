package backend

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/bo-er/corrupt-mysql/pkg"
	"github.com/inhies/go-bytesize"
)

const prepareBigTransactionSQL = `
DROP DATABASE IF EXISTS corrupt_mysql_bigtransaction_test;
CREATE DATABASE corrupt_mysql_bigtransaction_test;
USE corrupt_mysql_bigtransaction_test;
DROP TABLE IF EXISTS t;
CREATE TABLE t(x int primary key auto_increment);
INSERT INTO t() values(),(),(),();
`

// calling this 10 times we get a transaction of size 10KB
const quadraticGrowthSQL = `
USE corrupt_mysql_bigtransaction_test;
insert into t(x) select x+(select count(1) from t) from t;
`

func getBytes(size string) (float64, error) {
	b, err := bytesize.Parse(size)
	if err != nil {
		return 0, fmt.Errorf("parse %s to bytes: %s", size, err.Error())
	}
	parts := strings.Fields(b.Format("%f ", "kilobyte", true))
	if len(parts) != 2 {
		return 0, fmt.Errorf("parse %s to bytes: %s, got %s", size, "failure in kilobyte conversion", parts)
	}
	kb, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, fmt.Errorf("parse %s to uint64: %s", parts[0], err.Error())
	}
	return kb, nil
}

func getPower(size string) (int, error) {
	bytes, err := getBytes(size)
	if err != nil {
		return 0, fmt.Errorf("parse %s to bytes: %s", size, err.Error())
	}
	if bytes < 10 {
		return 0, errors.New("please enter a number that's bigger than 10 kilobytes")
	}
	ratio := bytes / 10
	return int(math.Log2(math.Ceil(ratio))), nil
}

func CreatesBigTransactions(c pkg.Connect, maxSize string) error {
	power, err := getPower(maxSize)
	if err != nil {
		return fmt.Errorf("getPower: %s", err.Error())
	}
	db, err := pkg.GetDB(c)
	if err != nil {
		return err
	}
	err = pkg.BatchExec(db, prepareBigTransactionSQL)
	if err != nil {
		return err
	}
	for i := 1; i <= 10+power; i++ {
		err = pkg.BatchExec(db, quadraticGrowthSQL)
		if err != nil {
			return fmt.Errorf("calling(%s) for the %dth time: %s", quadraticGrowthSQL, i, err.Error())
		}
	}
	return nil
}
