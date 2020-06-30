package controllers

import (
	"github.com/kataras/iris/v12"
	"mvc-login/datamodels"
	"mvc-login/services"
)

type UsersController struct {
	Ctx iris.Context
	Service services.UserService
}

// Get returns list of the users.
// Demo:
// curl -i -u admin:password http://localhost:8080/users
//
// The correct way if you have sensitive data:
// func (c *UsersController) Get() (results []viewmodels.User) {
// 	data := c.Service.GetAll()
//
// 	for _, user := range data {
// 		results = append(results, viewmodels.User{user})
// 	}
// 	return
// }
// otherwise just return the datamodels.
func (c *UsersController) Get() (results []datamodels.User) {
	return c.Service.GetAll()
}

// GetBy returns a user.
// Demo:
// curl -i -u admin:password http://localhost:8080/users/1
func (c *UsersController) GetBy(id int64) (user datamodels.User, found bool) {
	u, found := c.Service.GetByID(id)
	if !found {
		c.Ctx.Values().Set("message", "User could not be found!")
	}
	
	return u,found
}

// PutBy updates a user.
// Demo:
// curl -i -X PUT -u admin:password -F "username=kataras"
// -F "password=rawPasswordIsNotSafeIfOrNotHTTPs_You_Should_Use_A_client_side_lib_for_hash_as_well"
// http://localhost:8080/users/1
func (c *UsersController) PutBy(id int64) (datamodels.User, error) {
	u := datamodels.User{}

	if err := c.Ctx.ReadForm(&u); err != nil {
		return u, err
	}
	
	return c.Service.Update(id, u)
}

// DeleteBy deletes a user.
// Demo:
// curl -i -X DELETE -u admin:password http://localhost:8080/users/1
func (c *UsersController) DeleteBy(id int64) interface{} {
	wasDel := c.Service.DeleteByID(id)

	if wasDel {
		return map[string]interface{}{"deleted": id}
	}
	return iris.StatusBadRequest
}
