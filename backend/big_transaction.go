package backend

import (
	"errors"
	"fmt"
	"log"
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
flush logs;
`

// calling this 10 times we get a transaction of size 10KB
const quadraticGrowthSQL = `insert into t(x) select x+(select count(1) from t) from t;`

func getKiloBytes(size string) (float64, error) {
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

// we do have a pattern
// 20 10MB
// 21 20MB
// 22 40MB
// 23 80MB
// 24 160MB
func getPower(size string) (int, error) {
	kb, err := getKiloBytes(size)
	if err != nil {
		return 0, fmt.Errorf("parse %s to bytes: %s", size, err.Error())
	}
	if kb < 10 {
		return 0, errors.New("please enter a number that's bigger than 10 kilobytes")
	}
	fmt.Println("bytes: ", kb)
	ratio := (kb / 10240) * 20
	return int(math.Ceil(ratio)), nil
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
	// maybe there is a better way to do this.
	_, err = db.Exec(`USE corrupt_mysql_bigtransaction_test;`)
	if err != nil {
		return fmt.Errorf("use corrupt_mysql_bigtransaction_test: %s", err.Error())
	}
	for i := 1; i <= 10+power; i++ {
		log.Printf("executing SQL(%s), times: %d\n", quadraticGrowthSQL, i)
		_, err = db.Exec(quadraticGrowthSQL)
		if err != nil {
			return fmt.Errorf("calling(%s) for the %dth time: %s", quadraticGrowthSQL, i, err.Error())
		}
	}
	return nil
}
