package model

import (
	"src/abs"
	"src/table"
)

type bbsSpecialPopulation struct {
	abs.Model
}

func NewBbsSpecialPopulation() *bbsSpecialPopulation {
	return &bbsSpecialPopulation{}
}

func (sp *bbsSpecialPopulation) FindSpecialPopulation(ids []string, fields string) []table.BbsSpecialPopulation {
	var specialPopulationList []table.BbsSpecialPopulation
	sp.GetBbs().Select(fields).Where("user_id in (?)", ids).Find(&specialPopulationList)
	return specialPopulationList
}
