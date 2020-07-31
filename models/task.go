package models

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
)

// Task structure with id, name, and completion status
type Task struct {
	gorm.Model
	Name      string `gorm:"size:100;not null"        json:"name"`
	Completed bool   `gorm:"not null"                 json:"completed"`
}

func (t *Task) Prepare() {
	t.Name = strings.TrimSpace(t.Name)
	// v.CreatedBy = User{}
}

func (t *Task) Validate() error {
	if t.Name == "" {
		return errors.New("Name of task is required")
	}
	return nil
}

func (t *Task) Save(db *gorm.DB) (*Task, error) {
	var err error

	// Debug a single operation, show detailed log for this operation
	err = db.Debug().Create(&t).Error
	if err != nil {
		return &Task{}, err
	}
	return t, nil
}

func GetTasks(db *gorm.DB) (*[]Task, error) {
	tasks := []Task{}
	if err := db.Debug().Table("tasks").Find(&tasks).Error; err != nil {
		return &[]Task{}, err
	}
	return &tasks, nil
}

func GetTaskById(id int, db *gorm.DB) (*Task, error) {
	task := &Task{}
	if err := db.Debug().Table("tasks").Where("id = ?", id).First(task).Error; err != nil {
		return nil, err
	}
	return task, nil
}

// TODO: change to toggle task completion status
func (t *Task) UpdateTask(id int, db *gorm.DB) (*Task, error) {
	if err := db.Debug().Table("tasks").Where("id = ?", id).Updates(map[string]interface{}{
		"name":      t.Name,
		"completed": t.Completed}).Error; err != nil {
		return &Task{}, err
	}
	return t, nil
}

func DeleteTask(id int, db *gorm.DB) error {
	if err := db.Debug().Table("tasks").Where("id = ?", id).Delete(&Task{}).Error; err != nil {
		return err
	}
	return nil
}
