package application

import (
	"errors"

	"apiserver/pkg/util/jsonx"
)

type AppStatus int32
type UpdateStatus int32

const (
	AppBuilding  AppStatus = 0
	AppSuccessed AppStatus = 1
	AppFailed    AppStatus = 2
	AppRunning   AppStatus = 3

	StartFailed    UpdateStatus = 10
	StartSuccessed UpdateStatus = 11

	StopFailed    UpdateStatus = 20
	StopSuccessed UpdateStatus = 21

	ScaleFailed    UpdateStatus = 30
	ScaleSuccessed UpdateStatus = 31

	UpdateConfigFailed    UpdateStatus = 40
	UpdateConfigSuccessed UpdateStatus = 41

	RedeploymentFailed    UpdateStatus = 50
	RedeploymentSuccessed UpdateStatus = 51
)

//App is struct of application
type App struct {
	Id                       int          `json:"id" xorm:"pk not null autoincr int(11)"`
	Name                     string       `json:"name" xorm:"varchar(256)"`
	ReplicationConrollerName string       `json:"replicationConrollerName" xorm:"varchar(256)"`
	Region                   string       `json:"region" xorm:"varchar(256)"`
	Memory                   string       `json:"memory" xorm:"varchar(11)"`
	Cpu                      string       `json:"cpu" xorm:"varchar(11)"`
	InstanceCount            int32        `json:"instanceCount" xorm:"int(11)"`
	Envs                     string       `json:"envs" xorm:"varchar(256)"`
	Ports                    string       `json:"ports" xorm:"varchar(256)"`
	Image                    string       `json:"image" xorm:"varchar(256)"`
	Status                   AppStatus    `json:"status" xorm:"int(1)"` //构建中 0 成功 1 失败 2 运行中 3
	UpdateStatus             UpdateStatus `json:"status" xorm:"int(2)"`
	UserName                 string       `json:"userName" xorm:"varchar(256)"`
	Remark                   string       `json:"remark" xorm:"varchar(1024)"`
}

func (app *App) String() string {
	appStr, err := jsonx.ToJson(app)
	if err != nil {
		log.Errorf("node to string err :%s", err.Error())
		return ""
	}
	return appStr
}

func (app *App) Insert() error {
	_, err := engine.Insert(app)

	if err != nil {
		return err
	}

	return nil
}

func (app *App) Delete() error {
	_, err := engine.Id(app.Id).Delete(app)

	if err != nil {
		return err
	}

	return nil
}

func (app *App) Update() error {
	_, err := engine.Id(app.Id).Update(app)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) QueryOne() (*App, error) {
	has, err := engine.Id(app.Id).Get(app)

	if err != nil {
		return nil, err
	}

	if !has {
		return nil, errors.New("the query data not exist")
	}

	return app, nil
}

func (app *App) QuerySet() ([]*App, error) {
	appSet := []*App{}
	err := engine.Where("1 and 1 order by id desc").Find(&appSet)

	if err != nil {
		return nil, err
	}

	return appSet, nil
}
