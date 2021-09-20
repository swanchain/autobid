package test

import (
	"fmt"
	"go-swan/common/constants"
	"go-swan/common/utils"
	"go-swan/logs"
	"go-swan/models"
	"math/rand"
)

func Test() {
	//TestTask_GetTasks()
	//TestTask_GetAutoBidTasks()
	//TestMiner_GetAllMiners()
	models.TaskAssignMiner(1389, 14)
}

func TestTask_GetTasks() {
	tasks, err := models.GetTasks(0, 10, constants.TASK_STATUS_CREATED)
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}

	for _, task := range tasks {
		logs.GetLogger().Info(utils.ToJson(task))
	}
}

func TestTask_GetAutoBidTasks() {
	tasks, err := models.GetAutoBidTasks(0, 10, constants.TASK_STATUS_CREATED)
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}

	for _, task := range tasks {
		logs.GetLogger().Info(utils.ToJson(task))
	}
}

func TestTask_GetTaskById() {
	task, err := models.GetTaskById(1)
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}

	logs.GetLogger().Info(utils.ToJson(task))
}

func TestMiner_GetMiners() {
	miners, err := models.GetMiners(0, 10, constants.MINER_STATUS_ACTIVE)
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}

	for _, miner := range miners {
		logs.GetLogger().Info(utils.ToJson(miner))
	}
}

func TestMiner_GetAllMiners() {
	miners, err := models.GetAllMinersOrderByScore(constants.MINER_STATUS_ACTIVE)
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}

	totalScore := 0

	for _, miner := range miners {
		totalScore = totalScore + miner.Score
		//logs.GetLogger().Info(utils.ToJson(miner))
		//logs.GetLogger().Info("Score:", miner.Score, rand.Float64())
	}

	var minScorePercent float64 = 1
	var maxScorePercent float64 = 0

	for _, miner := range miners {
		miner.ScorePercent = float64(miner.Score) / float64(totalScore)
		if miner.ScorePercent < minScorePercent {
			minScorePercent = miner.ScorePercent
		}

		if miner.ScorePercent > maxScorePercent {
			maxScorePercent = miner.ScorePercent
		}
		//fmt.Println(miner.ScorePercent)
		//logs.GetLogger().Info(utils.ToJson(miner))
		//logs.GetLogger().Info("Score:", miner.Score, rand.Float64())
	}

	multiple := 1.0 / maxScorePercent

	for _, miner := range miners {
		oldScorePercent := miner.ScorePercent
		miner.ScorePercent = miner.ScorePercent * multiple
		fmt.Println(oldScorePercent, miner.Score, miner.ScorePercent)
	}

	fmt.Println(minScorePercent, maxScorePercent, multiple)

	for j := 0; j < 100; j++ {
		ratio := rand.ExpFloat64()
		for i, miner := range miners {
			//score := miner.ScorePercent * multiple
			//logs.GetLogger().Info("Score:", ratio, randNum)
			if ratio < miner.ScorePercent {
				fmt.Println(i, " ScorePercent=", miner.ScorePercent, " MinerFid=", miner.MinerFid, " ratio=", ratio, " is selected")
				break
			}
		}
	}
}
