package models

type Secret struct {
	ID    int32  `json:"id,omitempty"`
	Value string `json:"value,required"`
	Key   string `json:"key,required"`
	Url   string `json:"url,required"`
}

type User struct {
	ID       int32  `json:"id,omitempty"`
	Username string `json:"username,required"`
}
