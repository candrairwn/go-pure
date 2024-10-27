package utils

import (
	"crypto/rsa"
	"time"

	"github.com/candrairwn/go-pure/api/delivery/http/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	ErrTokenExpired = jwt.ErrTokenExpired
	privateKey      *rsa.PrivateKey
	publicKey       *rsa.PublicKey
)

type JWTUtil struct {
	ViperCustom *viper.Viper
	Log         *zap.SugaredLogger
}

func NewJWTUtil(viper *viper.Viper, log *zap.SugaredLogger) *JWTUtil {
	return &JWTUtil{ViperCustom: viper, Log: log}
}

func (J *JWTUtil) LoadFileKeys() error {
	privateData, err := ReadFileReturnByte(J.ViperCustom.GetString("JWT_SECRET_KEY"), J.Log)
	if err != nil {
		J.Log.Errorw("Error reading private key", "error", err)
		return err
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateData)
	if err != nil {
		J.Log.Errorw("Error parsing private key", "error", err)
		return err
	}

	publicData, err := ReadFileReturnByte(J.ViperCustom.GetString("JWT_PUBLIC_KEY"), J.Log)
	if err != nil {
		J.Log.Errorw("Error reading public key", "error", err)
		return err
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(publicData)
	if err != nil {
		J.Log.Errorw("Error parsing public key", "error", err)
		return err
	}

	privateKey = privKey
	publicKey = pubKey

	J.Log.Info("Success load keys")

	return nil

}

func CreateAccessToken(user model.UserJWT, expiry int64, log *zap.SugaredLogger) (string, error) {

	if expiry == 0 {
		expiry = time.Now().Add(time.Hour * 24).Unix()
	}

	claims := jwt.MapClaims{
		"iss":    viper.GetString("APP_NAME"),
		"expiry": expiry,
		"data": model.UserJWT{
			Version:    user.Version,
			Id:         user.Id,
			Username:   user.Username,
			IdTipeUser: user.IdTipeUser,
			IdProdi:    user.IdProdi,
			NamaProdi:  user.NamaProdi,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		log.Errorw("Error creating access token", "error", err)
		return "", err
	}

	return tokenString, nil
}

func VerifyAccessToken(tokenString string, log *zap.SugaredLogger) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, jwt.ErrInvalidKeyType
		}
		return publicKey, nil
	})
}
