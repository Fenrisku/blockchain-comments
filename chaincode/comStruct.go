package main

type Comment struct {
	ObjectType   string `json:"docType"`
	Cid          int    `json:"cid"`
	UsrName      string `json:"nick_name"`
	CommentScore int    `json:"first_comment_score"`

	CommentDetail string `json:"first_comment"`
	CommentTime   int    `json:"first_comment_time"`
	StudyTime     int    `json:first_comment_study_time`

	ClassName  string
	AgencyName string
	ID         string
	TxId       string
	Timestamp  string
	Signature  []byte
	Bind       []byte

	Historys []HistoryItemN // 当前com的历史记录
}

type HistoryItemN struct {
	TxId    string
	Comment Comment
}

type BlockCount struct {
	Name  string
	Score ScoreCount
}

type ScoreCount struct {
	Five  int
	Four  int
	Three int
	Two   int
	One   int
}
