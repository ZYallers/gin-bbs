package model

import (
	"src/abs"
	"src/table"
)

type bbsCircle struct {
	abs.Model
}

func NewBbsCircle() *bbsCircle {
	return &bbsCircle{}
}

func (c *bbsCircle) FindCircleIcon() []table.BbsCircle {
	var circleList []table.BbsCircle
	c.GetBbs().Where("state=?", "1").Find(&circleList)
	return circleList
}

func (c *bbsCircle) FindCircle(ids []string, fields string) []table.BbsCircle {
	var circleList []table.BbsCircle
	c.GetBbs().Where(ids).Select(fields).Find(&circleList)
	return circleList
}
