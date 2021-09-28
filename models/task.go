package models

import (
	"context"
	"go-swan/common/constants"
	"go-swan/database"
	"go-swan/logs"
	"strconv"
	"time"
)

type Task struct {
	Id             int      `json:"id"`
	TaskName       string   `json:"task_name"`
	Description    string   `json:"description"`
	TaskFileName   string   `json:"task_file_name"`
	CreatedOn      string   `json:"created_on"`
	UserId         int      `json:"user_id"`
	Status         string   `json:"status"`
	Tags           string   `json:"tags"`
	MinerId        *int     `json:"miner_id"`
	Type           *string  `json:"type"`
	IsPublic       int      `json:"is_public"`
	MinPrice       *float64 `json:"min_price"`
	MaxPrice       *float64 `json:"max_price"`
	ExpireDays     *int     `json:"expire_days"`
	Uuid           string   `json:"uuid"`
	CuratedDataset string   `json:"curated_dataset"`
	UpdatedOn      string   `json:"updated_on"`
	BidMode        *int     `json:"bid_mode"`
	FastRetrieval  *int     `json:"fast_retrieval"`
}

func GetTasks(pageNum int, pageSize int, status string) ([]*Task, error) {
	var tasks []*Task
	err := database.GetDB().Where("Status=?", status).Offset(pageNum).Limit(pageSize).Find(&tasks).Error

	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return tasks, nil
}

func GetAutoBidTasks(pageNum int, pageSize int, status string) ([]*Task, error) {
	var tasks []*Task
	err := database.GetDB().Where("is_public=1 and bid_mode=1 and miner_id is null and status=?", status).Offset(pageNum).Limit(pageSize).Find(&tasks).Error

	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return tasks, nil
}

func GetTaskById(id int) (*Task, error) {
	var task Task
	err := database.GetDB().Where("id=?", id).First(&task).Error

	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return &task, nil
}

func AddTask(task *Task) error {
	err := database.GetDB().Create(task).Error

	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	return err
}

func TaskAssignMiner(taskId, minerId, autoBidTaskCnt int, lastAutoBidAt int64) error {
	taskInfo := make(map[string]interface{})
	taskInfo["miner_id"] = minerId
	taskInfo["status"] = constants.TASK_STATUS_ASSIGNED
	taskInfo["updated_on"] = strconv.FormatInt(time.Now().UnixNano()/1e9, 10)

	lastAutoBidInfo := make(map[string]interface{})
	lastAutoBidInfo["auto_bid_task_cnt"] = autoBidTaskCnt
	lastAutoBidInfo["last_auto_bid_at"] = lastAutoBidAt

	ctx := context.Background()
	db := database.GetDB().BeginTx(ctx, nil)

	err := db.Model(&Task{}).Where("id=?", taskId).Update(taskInfo).Error
	if err != nil {
		db.Rollback()
		logs.GetLogger().Error(err)
		return err
	}

	err = db.Model(Miner{}).Where("id=?", minerId).Update(lastAutoBidInfo).Error
	if err != nil {
		db.Rollback()
		logs.GetLogger().Error(err)
		return err
	}

	err = db.Commit().Error
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	logs.GetLogger().Info("updated")

	return nil
}

func TaskUpdateStatus(taskId int, status string) error {
	taskInfo := make(map[string]interface{})
	taskInfo["status"] = status
	taskInfo["updated_on"] = time.Now().UnixNano() / 1e9

	err := database.GetDB().Model(&Task{}).Where("id=?", taskId).Update(taskInfo).Error

	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	return err
}

func EditTask(task Task) error {
	err := database.GetDB().Model(&Task{}).Where("id=?", task.Id).Update(task).Error

	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	return err
}

func DeleteTask(id int) error {
	err := database.GetDB().Where("id=?", id).Delete(Task{}).Error

	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	return err
}
