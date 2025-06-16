package models

import "github.com/golang-jwt/jwt/v5"

type MyClaims struct {
	Username              string `json:"username"`
	Id                    string `json:"id"`
	ProfilePictureAssetId string `json:"profile_picture_asset_id"`
	jwt.RegisteredClaims
}
