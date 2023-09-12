package sql

import (
	"fmt"
	"github.com/ansible-semaphore/semaphore/db"
	"github.com/masterminds/squirrel"
	"time"
)

func (d *SqlDb) getEvents(q squirrel.SelectBuilder, params db.RetrieveQueryParams) (events []db.Event, err error) {

	if params.Count > 0 {
		q = q.Limit(uint64(params.Count))
	}

	query, args, err := q.ToSql()

	if err != nil {
		return
	}

	_, err = d.selectAll(&events, query, args...)

	if err != nil {
		return
	}

	err = db.FillEvents(d, events)

	return
}

func (d *SqlDb) CreateEvent(evt db.Event) (newEvent db.Event, err error) {
	// 调整时间为东八区
	// todo 调整失败，还无法显示正确东八区时间，寄
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println("获取时区对象失败：", err)
		return
	}
	created := time.Now().In(loc)
	_, err = d.exec(
		"insert into event(user_id, project_id, object_id, object_type, description, created) values (?, ?, ?, ?, ?, ?)",
		evt.UserID,
		evt.ProjectID,
		evt.ObjectID,
		evt.ObjectType,
		evt.Description,
		created)

	if err != nil {
		return
	}

	newEvent = evt
	newEvent.Created = created
	return
}

func (d *SqlDb) GetUserEvents(userID int, params db.RetrieveQueryParams) ([]db.Event, error) {
	q := squirrel.Select("event.*, p.name as project_name").
		From("event").
		LeftJoin("project as p on event.project_id=p.id").
		OrderBy("created desc").
		LeftJoin("project__user as pu on pu.project_id=p.id").
		Where("p.id IS NULL or pu.user_id=?", userID)

	return d.getEvents(q, params)
}

func (d *SqlDb) GetEvents(projectID int, params db.RetrieveQueryParams) ([]db.Event, error) {
	q := squirrel.Select("event.*, p.name as project_name").
		From("event").
		LeftJoin("project as p on event.project_id=p.id").
		OrderBy("created desc").
		Where("event.project_id=?", projectID)

	return d.getEvents(q, params)
}
