package taskresult

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/lk153/quizgame-ai-serving/lib/strings"
)

type TaskResultEntity struct {
	ID      string  `bson:"id" json:"id" example:"35f1b935-58b1-42ed-8eea-10062906b84f"`
	Name    string  `bson:"name" json:"name"`
	Score   float64 `bson:"score" json:"score"`
	Comment string  `bson:"comment" json:"comment"`
}

func init() {
	uuid.EnableRandPool()
}

func NewAI() *TaskResultEntity {
	uuidStr := uuid.NewString()

	return &TaskResultEntity{
		ID: uuidStr,
	}
}

func (u *TaskResultEntity) SetName(name string) {
	u.Name = name
}

func (u *TaskResultEntity) Validate() (isValid bool, err error) {
	if strings.IsEmpty(u.ID) {
		isValid = false
		err = fmt.Errorf("task result's id is empty")
		return
	}

	if strings.IsEmpty(u.Name) {
		isValid = false
		err = fmt.Errorf("task result's name is empty")
		return
	}

	isValid = true
	return
}
