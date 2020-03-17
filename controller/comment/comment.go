package comment

import (
	"BlogBackend/db"
	"BlogBackend/model"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
	"net/http"
)

const BaseURL = "/api/comment"

type Controller struct {}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Initialize(e echo.Echo) (err error) {
	e.GET(BaseURL+"/all", getAllComments)
	e.GET(BaseURL+"/:id", getComment)
	e.GET(BaseURL + "/articleId/:articleId", getCommentsByArticleId)
	e.GET(BaseURL + "/replyId/:replyId", getCommentsByReplyId)
	e.POST(BaseURL, newComment)
	e.DELETE(BaseURL+"/:id", deleteComment)
	return nil
}

func getAllComments(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.Comment()
	defer closeConn()
	var results []model.Comment
	err = collection.Find(nil).All(&results)
	if err!=nil {
		return err
	}
	fmt.Println(results)
	return c.JSON(http.StatusOK, results)
}

func getComment(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.Comment()
	defer closeConn()
	id := c.Param("id")
	fmt.Println(id)
	var result model.Comment
	err = collection.FindId(bson.ObjectIdHex(id)).One(&result)
	if err!=nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func getCommentsByArticleId(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.Comment()
	defer closeConn()
	articleId := c.Param("articleId")
	//print(articleId)
	//print(bson.ObjectIdHex(articleId))
	var result []model.Comment
	err = collection.Find(bson.M{"articleId": bson.ObjectIdHex(articleId)}).Sort("-time").All(&result)
	if err!=nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func getCommentsByReplyId(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.Comment()
	defer closeConn()
	articleId := c.Param("replyId")
	var result []model.Comment
	_ = collection.Find(bson.M{"replyId": bson.ObjectIdHex(articleId)}).All(&result)
	return c.JSON(http.StatusOK, result)
}

func newComment(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.Comment()
	defer closeConn()
	var newComment model.Comment
	err = c.Bind(&newComment)
	if err != nil{
		return err
	}
	newComment.Id = bson.NewObjectId()
	err = collection.Insert(&newComment)
	if err != nil {
		return err
	}
	fmt.Println(newComment)
	return c.JSON(http.StatusOK, newComment)
}

func deleteComment(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.Comment()
	defer closeConn()
	id := c.Param("id")
	err = collection.RemoveId(bson.ObjectIdHex(id))
	return c.NoContent(http.StatusNoContent)
}
