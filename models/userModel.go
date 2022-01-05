package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"



)

type User struct {
	ID               primitive.ObjectID          `bson:"_id"`         //PRIMITIVE IS USED FOR USING BSON HERE
	Name             *string										`json:"name" validate:"required",min=2,max=100"`
	Password				 *string										`json:"password" validate:"required,min=6"`
	Email            *string										`json:"email" validate:"email,required"`
	Phone						 *string										`json:"phone" validate:"required"`
	Token            *string										`json:"token"`	
	User_type        *string										`json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	Refresh_token    *string										`json:"refresh_token"`
	Created_at       time.Time									`json:"created_at"`
	Updated_at       time.Time 									`json:"updated_at"`
	User_id          string										`json:"user_id"`

}