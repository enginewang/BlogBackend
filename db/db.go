package db

import (
	"github.com/globalsign/mgo"
)

type CollectionEntity func() (collection *mgo.Collection, clossConn func())

type Database struct {
	session  *mgo.Session
	database string
}

var GlobalDatabase *Database

const (
	CUser        = "user"
	CArticle     = "article"
	CAritcleTag  = "articleTag"
	CAritcleKind = "articleKind"
	CComment     = "comment"
)

func newDB(url string, database string) (*Database, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	db := &Database{
		session:  session,
		database: database,
	}

	cred := &mgo.Credential{
	Username: "root",
	Password: "Qwert@789",
	}
	err = session.Login(cred)
	if err != nil {
		panic(err)
	}
	return db, nil
}

func InitGlobalDB(url string, database string) error {
	d, err := newDB(url, database)
	if err != nil {
		return err
	}
	err = d.EnsureIndex()
	if err != nil {
		return err
	}
	GlobalDatabase = d
	return nil
}

func (d *Database) EnsureIndex() (err error) {
	return nil
}

func (d *Database) DB() (*mgo.Database, func()) {
	conn := d.session.Copy()
	return conn.DB(d.database), func() {
		conn.Close()
	}
}

func (d *Database) DropDatabase() (err error) {
	database, closeConn := d.DB()
	defer closeConn()
	err = database.DropDatabase()
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) User() (collection *mgo.Collection, clossConn func()) {
	database, closeConn := d.DB()
	c := database.C(CUser)
	return c, closeConn
}

func (d *Database) Article() (collection *mgo.Collection, clossConn func()) {
	database, closeConn := d.DB()
	c := database.C(CArticle)
	return c, closeConn
}

func (d *Database) ArticleTag() (collection *mgo.Collection, clossConn func()) {
	database, closeConn := d.DB()
	c := database.C(CAritcleTag)
	return c, closeConn
}

func (d *Database) ArticleKind() (collection *mgo.Collection, clossConn func()) {
	database, closeConn := d.DB()
	c := database.C(CAritcleKind)
	return c, closeConn
}

func (d *Database) Comment() (collection *mgo.Collection, clossConn func()) {
	database, closeConn := d.DB()
	c := database.C(CComment)
	return c, closeConn
}

func (d *Database) Collection(collectionName string) CollectionEntity {
	return func() (collection *mgo.Collection, closeConn func()) {
		database, closeConn := d.DB()
		c := database.C(collectionName)
		return c, closeConn
	}
}
