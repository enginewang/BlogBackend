package model

import (
	"github.com/globalsign/mgo/bson"
)

type Article struct {
	Id             bson.ObjectId   `bson:"_id,omitempty" json:"_id,omitempty"`
	Title          string          `bson:"title" json:"title"`
	Desc           string          `bson:"desc" json:"desc"`
	Cover          string          `bson:"cover" json:"cover"`
	PubTime        string          `bson:"pubTime" json:"pubTime"`
	Author         string          `bson:"author" json:"author"`
	Content        string          `bson:"content" json:"content"`
	Type           string          `bson:"type" json:"type"`
	Tags           string          `bson:"tags" json:"tags"`
	Kind           string          `bson:"kind" json:"kind"`
	ReadCount      int             `bson:"readCount" json:"readCount"`
	LikeCount      int             `bson:"likeCount" json:"likeCount"`
}

type ArticleSimple struct {
	Id             bson.ObjectId   `bson:"_id,omitempty" json:"_id,omitempty"`
	Title          string          `bson:"title" json:"title"`
	Desc           string          `bson:"desc" json:"desc"`
	Cover          string          `bson:"cover" json:"cover"`
	PubTime        string          `bson:"pubTime" json:"pubTime"`
	Type           string          `bson:"type" json:"type"`
	Tags           string          `bson:"tags" json:"tags"`
	Kind           string          `bson:"kind" json:"kind"`
}

type ArticleTag struct {
	Id         bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"`
	ArticleTag string        `bson:"articleTag" json:"articleTag"`
}

type ArticleKind struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"`
	ArticleKind string        `bson:"articleKind" json:"articleKind"`
}
