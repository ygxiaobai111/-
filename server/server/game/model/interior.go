package model

type OpenCollectionRsp struct {
	Limit    int8  `json:"limit"`
	CurTimes int8  `json:"cur_times"`
	NextTime int64 `json:"next_time"`
}

type CollectionRsp struct {
	Gold     int   `json:"gold"`
	Limit    int8  `json:"limit"`
	CurTimes int8  `json:"cur_times"`
	NextTime int64 `json:"next_time"`
}
type TransformReq struct {
	From []int `json:"from"` //0 Wood 1 Iron 2 Stone 3 Grain
	To   []int `json:"to"`   //0 Wood 1 Iron 2 Stone 3 Grain
}

type TransformRsp struct {
}
