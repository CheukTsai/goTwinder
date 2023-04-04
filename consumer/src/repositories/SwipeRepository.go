package repositories

import (
	"database/sql"
	"fmt"
	"goTwinderRMQConsumer/src/helpers"
	"goTwinderRMQConsumer/src/schemas"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func SwipesBatchInsertIntoTable(swipes []schemas.Swipe, db *sql.DB, tableName string) {
	if len(swipes) > 0 {
		tx, err := db.Begin()
		if (err != nil) {
			helpers.FailOnError(err, "Failed inserting records", "db.Begin()")
		}
		// defer tx.Rollback()
		stmtstring := fmt.Sprintf("INSERT INTO %s(userid, swipeeid) VALUES(?, ?)", tableName)
		stmt, err := tx.Prepare(stmtstring)
		for _, swipe := range swipes {
			stmt.Exec(strconv.Itoa(swipe.Swiper), strconv.Itoa(swipe.Swipee))
		}
		err = tx.Commit()
		if err != nil {
			helpers.FailOnError(err, "Failed inserting records", "commit()")
		}
		log.Printf("Successfully sent %d records into %s table", len(swipes), tableName)
	}
}