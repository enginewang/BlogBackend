package model

import "github.com/globalsign/mgo/bson"

type Comment struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"`
	Username  string        `bson:"username" json:"username"`
	Avatar    string        `bson:"avatar" json:"avatar"`
	Content   string        `bson:"content" json:"content"`
	Time      string        `bson:"time" json:"time"`
	ArticleId bson.ObjectId `bson:"articleId" json:"articleId"`
	ReplyId   bson.ObjectId `bson:"replyId,omitempty" json:"replyId,omitempty"`
}
