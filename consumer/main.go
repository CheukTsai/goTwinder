package main

import (
	"bytes"
	"encoding/json"
	"goTwinderRMQConsumer/src/tools"
	"goTwinderRMQConsumer/src/helpers"
	"goTwinderRMQConsumer/src/managers"
	"goTwinderRMQConsumer/src/repositories"
	"goTwinderRMQConsumer/src/schemas"
	"log"
)

var (
	numChannel = 10
	batchSize = 2
	likesTable = "likes"
	dislikesTable = "dislikes"
)

func main() {
	db := managers.MySqlConnectDatabase()
	cp := tools.NewRMQChannelPool(numChannel)
	
	for i := 0; i < numChannel; i++ {
		threadNum, cnt, likes, dislikes := i, 0, []schemas.Swipe{}, []schemas.Swipe{}
		go func() {
			for {
				msgs := managers.RMQConsumeWithQName(cp, "swipeQueue", threadNum)
				for d := range msgs {
					log.Printf("Received a message: %s from thread %d", d.Body, threadNum)
					var swipe schemas.Swipe
					e := json.NewDecoder(bytes.NewReader(d.Body)).Decode(&swipe)
					helpers.FailOnError(e, "Failed decoding json", "json.Newcoder().Decode()")
					if swipe.Like {
						likes = append(likes, swipe)
					} else {
						dislikes = append(dislikes, swipe)
					}
					cnt++
					if (cnt == batchSize) {
						repositories.SwipesBatchInsertIntoTable(likes, db, likesTable)
						repositories.SwipesBatchInsertIntoTable(dislikes, db, dislikesTable)
						cnt, likes, dislikes = 0, likes[:0], dislikes[:0]
					}
				}
			}
		}()
	}
	select{}
}