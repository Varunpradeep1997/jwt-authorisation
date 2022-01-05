package helper

import (
	"errors"

	"fmt"

	"github.com/gin-gonic/gin"
)


func CheckUserType(c *gin.Context,role string) (err error) {

	userType:=c.GetString("user_type")
	userId := c.GetString("uid")
	err=nil
	if userType !=role{
		err=errors.New("Unauthorized access for the resource - user type error")
		err=fmt.Errorf("%q %q  %q ",role,userType,userId)
	//	err=fmt.Errorf("%q",userType)

		return err
	}
	return err
}

 func MatchUserTypeToUid(c *gin.Context,userId string)(err error) {
userType:=c.GetString("user_type")
uid:=c.GetString("uid")
err=nil

if userType=="USER" && uid != userId {
	err=errors.New("Unauthorised access to this resource - USER and uid mismatch")
	return err
}
err=CheckUserType(c,userType)
return err
 }