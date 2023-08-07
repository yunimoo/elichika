package serverdb

import (
	"elichika/model"

	"time"
)

// card grade up trigger is responsible for showing the pop-up animation when openning a card after getting a new copy
// or right after performing a limit break using items

func (session *Session) AddTriggerCardGradeUp(id int64, trigger *model.TriggerCardGradeUp) {
	if id == 0 {
		id = time.Now().UnixNano()
	}
	if trigger != nil {
		trigger.TriggerID = id
	}
	session.TriggerCardGradeUps = append(session.TriggerCardGradeUps, id)
	session.TriggerCardGradeUps = append(session.TriggerCardGradeUps, trigger)
}

func (session *Session) AddTriggerBasic(id int64, trigger *model.TriggerBasic) {
	if id == 0 {
		id = time.Now().UnixNano()
	}
	if trigger != nil {
		trigger.TriggerID = id
	}
	session.TriggerBasics = append(session.TriggerBasics, id)
	session.TriggerBasics = append(session.TriggerBasics, trigger)
}
