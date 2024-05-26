package query

import (
	"cwgo_db/biz/dao/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

var db *gorm.DB
var userQuery user

func setup() {
	dsn := "root:root@tcp(127.0.0.1:3306)/testdb?charset=utf8&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	//err = db.AutoMigrate(&model.User{})
	//if err != nil {
	//	panic("failed to migrate database")
	//}
	db.Exec("USE testdb")
	db.Exec("CREATE TABLE IF NOT EXISTS users (" + "id INT AUTO_INCREMENT PRIMARY KEY, " + "name VARCHAR(100) NOT NULL, " + "email VARCHAR(100) UNIQUE NOT NULL)")
	userQuery = newUser(db)
}

func teardown() {
	db.Exec("DROP TABLE IF EXISTS users")
}

func TestCreateUser(t *testing.T) {
	setup()
	defer teardown()

	newUser := &model.User{
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}

	err := userQuery.userDo.Create(newUser)
	assert.NoError(t, err)
	assert.NotEqual(t, 1, newUser.ID)
}

func TestGetUserByID(t *testing.T) {
	setup()
	defer teardown()

	newUser := &model.User{
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}

	userQuery.userDo.Create(newUser)

	user, err := userQuery.userDo.First()
	assert.NoError(t, err)
	assert.Equal(t, newUser.Name, user.Name)
	assert.Equal(t, newUser.Email, user.Email)
}

func TestFindUsers(t *testing.T) {
	setup()
	defer teardown()

	newUser1 := &model.User{
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}
	newUser2 := &model.User{
		Name:  "Jane Doe",
		Email: "jane.doe@example.com",
	}

	userQuery.userDo.Create(newUser1)
	userQuery.userDo.Create(newUser2)

	users, err := userQuery.userDo.Find()
	assert.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestDeleteUser(t *testing.T) {
	setup()
	defer teardown()

	newUser := &model.User{
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}

	userQuery.userDo.Create(newUser)

	_, err := userQuery.userDo.Delete(newUser)
	assert.NoError(t, err)

	_, err = userQuery.userDo.First()
	assert.Error(t, err)
}
