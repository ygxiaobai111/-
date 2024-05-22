package model

type ConfigReq struct {
}

type ConfigRsp struct {
	Confs []Conf
}

type Conf struct {
	Type     int8   `json:"type"`
	Level    int8   `json:"level"`
	Name     string `json:"name"`
	Wood     int    `json:"Wood"`
	Iron     int    `json:"iron"`
	Stone    int    `json:"stone"`
	Grain    int    `json:"grain"`
	Durable  int    `json:"durable"`  //耐久
	Defender int    `json:"defender"` //防御等级
}
