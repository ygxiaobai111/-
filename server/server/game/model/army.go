package model

type ArmyListReq struct {
	CityId int `json:"cityId"`
}

type ArmyListRsp struct {
	CityId int    `json:"cityId"`
	Armys  []Army `json:"armys"`
}
