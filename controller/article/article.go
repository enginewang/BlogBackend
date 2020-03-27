package article

import (
	"BlogBackend/controller/auth"
	"BlogBackend/db"
	"BlogBackend/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
	"net/http"
)

const BaseURL = "/api/article"

type Controller struct {}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Initialize(e echo.Echo) (err error) {
	e.GET(BaseURL+"/all", getAllArticles)
	e.GET(BaseURL+"/indexArticleList", getIndexArticle)
	e.GET(BaseURL+"/:id/:visit", getArticleDetail)
	e.GET(BaseURL+"/title/:title", getArticleByTitle)
	e.GET(BaseURL+"/clickLove/:id", loveArticle)
	e.GET(BaseURL+"/clickCancelLove/:id", cancelLoveArticle)
	e.POST(BaseURL, newArticle, auth.IsLoggedIn, needAdminAuth)
	e.DELETE(BaseURL+"/:id", deleteArticle, auth.IsLoggedIn, needAdminAuth)
	e.PUT(BaseURL+"/:id", updateArticle, auth.IsLoggedIn, needAdminAuth)

	e.GET(BaseURL+"/tag/all", getAllArticleTags)
	e.GET(BaseURL+"/tag/:tagName", getArticleTag)
	e.POST(BaseURL+"/tag", newArticleTag, auth.IsLoggedIn, needAdminAuth)
	e.DELETE(BaseURL+"/tag/tagName/:tagName", deleteArticleTagByName, auth.IsLoggedIn, needAdminAuth)

	e.GET(BaseURL+"/kind/all", getAllArticleKinds)
	e.GET(BaseURL+"/kind/:kindName", getArticleKind)
	e.POST(BaseURL+"/kind", newArticleKind, auth.IsLoggedIn, needAdminAuth)
	e.DELETE(BaseURL+"/kind/kindName/:kindName", deleteArticleKindByName, auth.IsLoggedIn, needAdminAuth)
	e.DELETE(BaseURL+"/kind/:id", deleteArticleKind, auth.IsLoggedIn, needAdminAuth)
	return nil
}

func getAllArticles(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.Article()
	defer closeConn()
	var results []model.Article
	err = collection.Find(nil).Sort("-pubTime").All(&results)
	if err!=nil {
		return err
	}
	return c.JSON(http.StatusOK, results)
}

func getIndexArticle(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.Article()
	defer closeConn()
	var results []model.ArticleSimple
	err = collection.Find(nil).Sort("-pubTime").All(&results)
	if err!=nil {
		return err
	}
	return c.JSON(http.StatusOK, results)
}

func getArticleDetail(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.Article()
	defer closeConn()
	id := c.Param("id")
	visit := c.Param("visit")
	var result model.Article
	err = collection.FindId(bson.ObjectIdHex(id)).One(&result)
	if visit == "1" {
		result.ReadCount++
		err = collection.UpdateId(bson.ObjectIdHex(id), result)
	}
	if err!=nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func getArticleByTitle(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.Article()
	defer closeConn()
	title := c.Param("title")
	var result model.Article
	err = collection.Find(bson.M{"title": title}).One(&result)
	if err!=nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func newArticle(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.Article()
	defer closeConn()
	var newArticle model.Article
	err = c.Bind(&newArticle)
	if err != nil{
		return err
	}
	newArticle.Id = bson.NewObjectId()
	err = collection.Insert(&newArticle)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, newArticle)
}

func deleteArticle(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.Article()
	defer closeConn()
	kindId := c.Param("id")
	err = collection.RemoveId(bson.ObjectIdHex(kindId))
	return c.NoContent(http.StatusNoContent)
}

func updateArticle(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.Article()
	defer closeConn()
	var updatedArticle model.Article
	err = c.Bind(&updatedArticle)
	if err != nil{
		return err
	}
	id := c.Param("id")
	err = collection.UpdateId(bson.ObjectIdHex(id), updatedArticle)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, updatedArticle)
}

func loveArticle(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.Article()
	defer closeConn()
	id := c.Param("id")
	var article model.Article
	err = collection.FindId(bson.ObjectIdHex(id)).One(&article)
	if err != nil{
		return err
	}
	article.LikeCount+=1
	err = collection.UpdateId(bson.ObjectIdHex(id), article)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, article)
}

func cancelLoveArticle(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.Article()
	defer closeConn()
	id := c.Param("id")
	var article model.Article
	err = collection.FindId(bson.ObjectIdHex(id)).One(&article)
	if err != nil{
		return err
	}
	article.LikeCount-=1
	err = collection.UpdateId(bson.ObjectIdHex(id), article)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, article)
}


func getAllArticleTags(c echo.Context) (err error)  {
	collection, closeConn := db.GlobalDatabase.ArticleTag()
	defer closeConn()
	var results []model.ArticleTag
	err = collection.Find(nil).All(&results)
	if err!=nil {
		return err
	}
	return c.JSON(http.StatusOK, results)
}

func getArticleTag(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.ArticleTag()
	defer closeConn()
	tagName := c.Param("tagName")
	var result model.ArticleTag
	err = collection.Find(bson.M{"articleTag": tagName}).One(&result)
	if err!=nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func newArticleTag(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.ArticleTag()
	defer closeConn()
	var newArticleTag model.ArticleTag
	err = c.Bind(&newArticleTag)
	if err != nil{
		return err
	}
	newArticleTag.Id = bson.NewObjectId()
	err = collection.Insert(&newArticleTag)
	return c.JSON(http.StatusOK, newArticleTag)
}

func deleteArticleTagByName(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.ArticleTag()
	defer closeConn()
	tagName := c.Param("tagName")
	err = collection.Remove(bson.M{"articleTag": tagName})
	return c.NoContent(http.StatusNoContent)
}

func getAllArticleKinds(c echo.Context) (err error)  {
	collection, closeConn := db.GlobalDatabase.ArticleKind()
	defer closeConn()
	var results []model.ArticleKind
	err = collection.Find(nil).All(&results)
	if err!=nil {
		return err
	}
	return c.JSON(http.StatusOK, results)
}

func getArticleKind(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.ArticleKind()
	defer closeConn()
	kindName := c.Param("kindName")
	var result model.ArticleKind
	err = collection.Find(bson.M{"articleKind": kindName}).One(&result)
	if err!=nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func newArticleKind(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.ArticleKind()
	defer closeConn()
	var newArticleKind model.ArticleKind
	err = c.Bind(&newArticleKind)
	if err != nil{
		return err
	}
	newArticleKind.Id = bson.NewObjectId()
	err = collection.Insert(&newArticleKind)
	return c.JSON(http.StatusOK, newArticleKind)
}

func deleteArticleKindByName(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.ArticleKind()
	defer closeConn()
	kindName := c.Param("kindName")
	err = collection.Remove(bson.M{"articleKind": kindName})
	return c.NoContent(http.StatusNoContent)
}

func deleteArticleKind(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.ArticleKind()
	defer closeConn()
	kindId := c.Param("id")
	err = collection.RemoveId(bson.ObjectIdHex(kindId))
	return c.NoContent(http.StatusNoContent)
}

func needAdminAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func (c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		name := claims["name"].(string)
		if name == "admin" {
			return next(c)
		} else {
			return c.String(http.StatusUnauthorized, "您没有管理员权限，不能请求！")
		}
	}
}

func needUserAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func (c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		name := claims["name"].(string)
		if name != "" {
			return next(c)
		} else {
			return c.String(http.StatusUnauthorized, "您没有登陆，不能请求！")
		}
	}
}