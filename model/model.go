package model

import (
	"BlogBackend/db"
	"github.com/globalsign/mgo/bson"
)

type Entity interface {
	GenerateID()
	CollectionName() string
}

type Helper struct{
	t Entity
}

func NewHelper(entity Entity) *Helper {
	return &Helper{
		t: entity,
	}
}

func (h *Helper) GetCollection() db.CollectionEntity {
	return db.GlobalDatabase.Collection(h.t.CollectionName())
}

func (h *Helper) All(result interface{}) (err error) {
	collection, closeConn := h.GetCollection()()
	defer closeConn()
	return collection.Find(nil).All(result)
}

func (h *Helper) Query(query bson.M, result interface{}) (err error)  {
	collection, closeConn := h.GetCollection()()
	defer closeConn()
	return collection.Find(query).All(result)
}

func (h *Helper) QueryOne(query bson.M, result Entity) (err error)  {
	collection, closeConn := h.GetCollection()()
	defer closeConn()
	return collection.Find(query).One(result)
}

func (h *Helper) CountNum(query interface{}) (num int, err error)  {
	collection, closeConn := h.GetCollection()()
	defer closeConn()
	return collection.Find(query).Count()
}