package user

import (
	"BlogBackend/controller/auth"
	"BlogBackend/db"
	"BlogBackend/model"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
	"net/http"
	"strings"
)

const BaseURL = "/api/user"

type Controller struct{}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Initialize(e echo.Echo) (err error) {
	e.GET(BaseURL+"/all", getAllUsers, auth.IsLoggedIn, needAdminAuth)
	e.GET(BaseURL+"/allUserInfo", getAllUsersInfo)
	e.GET(BaseURL+"/:id", getUser, auth.IsLoggedIn, needUserAuth)
	e.GET(BaseURL+"/username/:username/addToLikeList/:articleId", addToLikeList)
	e.GET(BaseURL+"/username/:username/removeFromLikeList/:articleId", removeFromLikeList)
	e.PUT(BaseURL+"/username/:username", updateUser, auth.IsLoggedIn, needUserAuth)
	e.GET(BaseURL+"/username/:username", getUserByUserName, auth.IsLoggedIn, needUserAuth)
	e.POST(BaseURL, newUser, auth.IsLoggedIn, needAdminAuth)
	e.DELETE(BaseURL+"/:id", deleteUser, auth.IsLoggedIn, needUserAuth)
	return nil
}

func getAllUsers(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.User()
	defer closeConn()
	var results []model.User
	err = collection.Find(nil).All(&results)
	if err != nil {
		return err
	}
	fmt.Println(results)
	return c.JSON(http.StatusOK, results)
}


func getAllUsersInfo(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.User()
	defer closeConn()
	var results []model.UserInfo
	err = collection.Find(nil).All(&results)
	if err != nil {
		return err
	}
	fmt.Println(results)
	return c.JSON(http.StatusOK, results)
}


func getUser(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.User()
	defer closeConn()
	id := c.Param("id")
	var result model.User
	err = collection.FindId(bson.ObjectIdHex(id)).One(&result)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func updateUser(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.User()
	defer closeConn()
	var updatedUser model.User
	err = c.Bind(&updatedUser)
	if err != nil {
		return err
	}
	username := c.Param("username")
	err = collection.Update(bson.M{"username": username}, updatedUser)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, updatedUser)
}

func getUserByUserName(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.User()
	defer closeConn()
	username := c.Param("username")
	var result model.User
	err = collection.Find(bson.M{"username": username}).One(&result)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func addToLikeList(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.User()
	defer closeConn()
	username := c.Param("username")
	var user model.User
	err = collection.Find(bson.M{"username": username}).One(&user)
	if err != nil {
		return err
	}
	articleId := c.Param("articleId")

	if user.LikeList == "" {
		user.LikeList = articleId
	} else {
		userLikeList := strings.Split(user.LikeList, ",")
		for _, aid := range userLikeList {
			if aid == articleId {
				return c.JSON(http.StatusNoContent, "已经点过赞了！")
			}
		}
		userLikeList = append(userLikeList, articleId)
		newUserLikeList := strings.Join(userLikeList, ",")
		user.LikeList = newUserLikeList
	}
	err = collection.Update(bson.M{"username": username}, user)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, user)
}

func removeFromLikeList(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.User()
	defer closeConn()
	username := c.Param("username")
	var user model.User
	err = collection.Find(bson.M{"username": username}).One(&user)
	if err != nil {
		return err
	}
	articleId := c.Param("articleId")
	userLikeList := strings.Split(user.LikeList, ",")
	for i, aid := range userLikeList {
		if aid == articleId {
			userLikeList = append(userLikeList[:i], userLikeList[i+1:]...)
			newUserLikeList := strings.Join(userLikeList, ",")
			user.LikeList = newUserLikeList
			err = collection.Update(bson.M{"username": username}, user)
			if err != nil {
				return err
			}
			return c.JSON(http.StatusOK, user)
		}
	}
	return c.JSON(http.StatusNoContent, "收藏列表里没有！")
}

func newUser(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.User()
	defer closeConn()
	var newUser model.User
	err = c.Bind(&newUser)
	// 通过用户注册界面注册的用户只能是可以评论的读者
	newUser.Role = 2
	if err != nil {
		return err
	}
	newUser.LikeList = ""
	newUser.Coin = 0
	newUser.Exp = 0
	newUser.Level = 1
	newUser.Extra = ""
	newUser.Id = bson.NewObjectId()
	err = collection.Insert(&newUser)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, newUser)
}

func deleteUser(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.User()
	defer closeConn()
	userId := c.Param("id")
	err = collection.RemoveId(bson.ObjectIdHex(userId))
	return c.NoContent(http.StatusNoContent)
}

func needAdminAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
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
	return func(c echo.Context) error {
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
