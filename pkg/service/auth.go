package service

import (
	"abir/models"
	"abir/pkg/repository"
	"abir/pkg/storage"
	"abir/pkg/utils"
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"math/rand"
	"strconv"
	"time"
)

const (
	tokenTTL = 30 * 24 * time.Hour
)

type AuthService struct {
	repo repository.Authorization
	redisClient *redis.Client
	fileStorage storage.Storage
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
	UserType string `json:"user_type"`
}

func NewAuthService(repo repository.Authorization, client *redis.Client, fileStorage storage.Storage) *AuthService {
	return &AuthService{repo: repo, redisClient: client, fileStorage: fileStorage}
}

func (s *AuthService) CreateClient(ctx context.Context, client models.Client, userId int) error {
	fileName, err := utils.GenerateFileName()
	if err != nil {
		return err
	}
	if client.Avatar != nil {
		uploadedFileName, err := s.fileStorage.Upload(ctx, storage.UploadInput{
			File:        *client.Avatar,
			Name:        fileName,
			Folder: 	"clients",
		})
		if err != nil {
			return err
		}
		client.Avatar = &uploadedFileName
	}
	return s.repo.CreateClient(client, userId)
}

func (s *AuthService) SendActivationCode(userId int, phone string) error {
	codeMin := 10000
	codeMax := 99999
	code := strconv.Itoa(rand.Intn(codeMax - codeMin) + codeMin)
	code = strconv.Itoa(11111)

	_, ok := s.redisClient.Get("update_phone"+ strconv.Itoa(userId)).Result()
	if ok == nil {
		return errors.New("try after a while")
	}
	err := s.redisClient.Set("update_phone"+strconv.Itoa(userId), code, 2 * time.Minute).Err()
	if err != nil {
		return err
	}
	err = s.repo.ClientCheckPhone(phone)
	if err != nil {
		return err
	}
	//err = utils.SendSms(login, "Your verification code - "+code)
	//if err != nil {
	//	return err
	//}
	return nil
}

func (s *AuthService) ClientUpdatePhone(userId int, phone, code string) error{
	activationCode, err := s.redisClient.Get("update_phone"+ strconv.Itoa(userId)).Result()
	if err != nil {
		return err
	}
	if activationCode != code {
		return errors.New("wrong code from sms")
	}
	return s.repo.ClientUpdatePhone(userId, phone)
}

func (s *AuthService) GetClient(userId int) (models.Client, error)  {
	return s.repo.GetClient(userId)
}

func (s *AuthService) GetDriver(userId int) (models.Driver, error)  {
	return s.repo.GetDriver(userId)
}

func (s *AuthService) GetDriverVerification(userId int) ([]models.DriverVerification, error){
	return s.repo.GetDriverVerification(userId)
}

func (s *AuthService) GetDriverCar(userId int) (models.DriverCar, error){
	return s.repo.GetDriverCar(userId)
}
func (s *AuthService) GetDriverCarInfo(langId, userId int) (models.DriverCarInfo, error){
	return s.repo.GetDriverCarInfo(langId, userId)
}
func (s *AuthService) GetDriverInfo(langId, driverId int) (models.Driver, models.DriverCar, models.DriverCarInfo, error){
	return s.repo.GetDriverInfo(langId, driverId)
}
func (s *AuthService) GenerateToken(login, password, userType string) (string, error)  {
	_, err := s.redisClient.Get(userType + "_login"+login).Result()
	if err != nil {
		return "", errors.New("code expired")
	}
	user, err := s.repo.GetUser(login, generatePasswordHash(password), userType)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt: time.Now().Unix(),
		},
		user.Id,
		userType,
	})

	return token.SignedString([]byte(viper.GetString("auth.signing_key")))
}

func (s *AuthService) ParseToken(accessToken string) (int, string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok{
			return nil, errors.New("invalid signing method")
		}
		return []byte(viper.GetString("auth.signing_key")), nil
	})
	if err != nil {
		return 0, "", err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, "", errors.New("token claims are not of type")
	}
	return claims.UserId, claims.UserType, nil
}
func (s *AuthService) ClientSendCode(login string) error {
	codeMin := 10000
	codeMax := 99999
	code := strconv.Itoa(rand.Intn(codeMax - codeMin) + codeMin)
	code = strconv.Itoa(11111)

	_, ok := s.redisClient.Get("client_login"+login).Result()
	if ok == nil {
		return errors.New("try after a while")
	}
	err := s.redisClient.Set("client_login"+login, code, 2 * time.Minute).Err()
	if err != nil {
		return err
	}
	_, err = s.repo.CreateOrUpdateClient(models.User{Login: login})
	if err != nil {
		return err
	}
	//err = utils.SendSms(login, "Your verification code - "+code)
	//if err != nil {
	//	return err
	//}

	return s.repo.ClientSendCode(login, generatePasswordHash(code))
}

func generatePasswordHash(password string) string  {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(viper.GetString("auth.hash_salt"))))
}

func (s *AuthService) DriverSendCode(login string) error {
	codeMin := 10000
	codeMax := 99999
	code := strconv.Itoa(rand.Intn(codeMax - codeMin) + codeMin)
	code = strconv.Itoa(11111)

	_, ok := s.redisClient.Get("driver_login"+login).Result()
	if ok == nil {
		return errors.New("try after a while")
	}
	err := s.redisClient.Set("driver_login"+login, code, 2 * time.Minute).Err()
	if err != nil {
		return err
	}
	_, err = s.repo.CreateOrUpdateDriver(models.User{Login: login})
	if err != nil {
		return err
	}
	//err = utils.SendSms(login, "Your verification code - "+code)
	//if err != nil {
	//	return err
	//}

	return s.repo.DriverSendCode(login, generatePasswordHash(code))
}

func (s *AuthService) CreateDriver(ctx context.Context, driver models.Driver, userId int) error {
	if driver.Photo != nil && *driver.Photo != "" {
		fileName, err := utils.GenerateFileName()
		if err != nil {
			return err
		}
		uploadedFileName, err := s.fileStorage.Upload(ctx, storage.UploadInput{
			File:        *driver.Photo,
			Name:        fileName,
			Folder: 	"drivers",
		})
		if err != nil {
			return err
		}
		driver.Photo = &uploadedFileName
	}
	return s.repo.CreateDriver(driver, userId)
}

func (s *AuthService) SendForModerating(userId int) error {
	return s.repo.SendForModerating(userId)
}
func (s *AuthService) UpdateDriver(ctx context.Context, driver models.Driver, userId int) error {
	if driver.Photo != nil && *driver.Photo != "" {
		fileName, err := utils.GenerateFileName()
		if err != nil {
			return err
		}
		uploadedFileName, err := s.fileStorage.Upload(ctx, storage.UploadInput{
			File:        *driver.Photo,
			Name:        fileName,
			Folder: 	"drivers",
		})
		if err != nil {
			return err
		}
		driver.Photo = &uploadedFileName
	}
	if driver.PassportCopy1 != nil && *driver.PassportCopy1 != "" {
		fileName, err := utils.GenerateFileName()
		if err != nil {
			return err
		}
		uploadedFileName, err := s.fileStorage.Upload(ctx, storage.UploadInput{
			File:        *driver.PassportCopy1,
			Name:        fileName,
			Folder: 	"drivers",
		})
		if err != nil {
			return err
		}
		driver.PassportCopy1 = &uploadedFileName
	}
	if driver.PassportCopy2 != nil && *driver.PassportCopy2 != "" {
		fileName, err := utils.GenerateFileName()
		if err != nil {
			return err
		}
		uploadedFileName, err := s.fileStorage.Upload(ctx, storage.UploadInput{
			File:        *driver.PassportCopy2,
			Name:        fileName,
			Folder: 	"drivers",
		})
		if err != nil {
			return err
		}
		driver.PassportCopy2 = &uploadedFileName
	}
	if driver.PassportCopy3 != nil && *driver.PassportCopy3 != "" {
		fileName, err := utils.GenerateFileName()
		if err != nil {
			return err
		}
		uploadedFileName, err := s.fileStorage.Upload(ctx, storage.UploadInput{
			File:        *driver.PassportCopy3,
			Name:        fileName,
			Folder: 	"drivers",
		})
		if err != nil {
			return err
		}
		driver.PassportCopy3 = &uploadedFileName
	}
	if driver.DriverLicensePhoto1 != nil && *driver.DriverLicensePhoto1 != "" {
		fileName, err := utils.GenerateFileName()
		if err != nil {
			return err
		}
		uploadedFileName, err := s.fileStorage.Upload(ctx, storage.UploadInput{
			File:        *driver.DriverLicensePhoto1,
			Name:        fileName,
			Folder: 	"drivers",
		})
		if err != nil {
			return err
		}
		driver.DriverLicensePhoto1 = &uploadedFileName
	}
	if driver.DriverLicensePhoto2 != nil && *driver.DriverLicensePhoto2 != "" {
		fileName, err := utils.GenerateFileName()
		if err != nil {
			return err
		}
		uploadedFileName, err := s.fileStorage.Upload(ctx, storage.UploadInput{
			File:        *driver.DriverLicensePhoto2,
			Name:        fileName,
			Folder: 	"drivers",
		})
		if err != nil {
			return err
		}
		driver.DriverLicensePhoto2 = &uploadedFileName
	}
	if driver.DriverLicensePhoto3 != nil && *driver.DriverLicensePhoto3 != "" {
		fileName, err := utils.GenerateFileName()
		if err != nil {
			return err
		}
		uploadedFileName, err := s.fileStorage.Upload(ctx, storage.UploadInput{
			File:        *driver.DriverLicensePhoto3,
			Name:        fileName,
			Folder: 	"drivers",
		})
		if err != nil {
			return err
		}
		driver.DriverLicensePhoto3 = &uploadedFileName
	}
	return s.repo.UpdateDriver(driver, userId)
}

func (s *AuthService) UpdateDriverCar(ctx context.Context,car models.DriverCar, userId int) error{
	if car.PhotoTexpasport1 != nil && *car.PhotoTexpasport1 != "" {
		fileName, err := utils.GenerateFileName()
		if err != nil {
			return err
		}
		uploadedFileName, err := s.fileStorage.Upload(ctx, storage.UploadInput{
			File:        *car.PhotoTexpasport1,
			Name:        fileName,
			Folder: 	"drivers",
		})
		if err != nil {
			return err
		}
		car.PhotoTexpasport1 = &uploadedFileName
	}
	if car.PhotoTexpasport2 != nil && *car.PhotoTexpasport2 != "" {
		fileName, err := utils.GenerateFileName()
		if err != nil {
			return err
		}
		uploadedFileName, err := s.fileStorage.Upload(ctx, storage.UploadInput{
			File:        *car.PhotoTexpasport2,
			Name:        fileName,
			Folder: 	"drivers",
		})
		if err != nil {
			return err
		}
		car.PhotoTexpasport2 = &uploadedFileName
	}
	if car.CarFront != nil && *car.CarFront != "" {
		fileName, err := utils.GenerateFileName()
		if err != nil {
			return err
		}
		uploadedFileName, err := s.fileStorage.Upload(ctx, storage.UploadInput{
			File:        *car.CarFront,
			Name:        fileName,
			Folder: 	"drivers",
		})
		if err != nil {
			return err
		}
		car.CarFront = &uploadedFileName
	}
	if car.CarBack != nil && *car.CarBack != "" {
		fileName, err := utils.GenerateFileName()
		if err != nil {
			return err
		}
		uploadedFileName, err := s.fileStorage.Upload(ctx, storage.UploadInput{
			File:        *car.CarBack,
			Name:        fileName,
			Folder: 	"drivers",
		})
		if err != nil {
			return err
		}
		car.CarBack = &uploadedFileName
	}
	if car.CarLeft != nil && *car.CarLeft != "" {
		fileName, err := utils.GenerateFileName()
		if err != nil {
			return err
		}
		uploadedFileName, err := s.fileStorage.Upload(ctx, storage.UploadInput{
			File:        *car.CarLeft,
			Name:        fileName,
			Folder: 	"drivers",
		})
		if err != nil {
			return err
		}
		car.CarLeft = &uploadedFileName
	}
	if car.CarRight != nil && *car.CarRight != "" {
		fileName, err := utils.GenerateFileName()
		if err != nil {
			return err
		}
		uploadedFileName, err := s.fileStorage.Upload(ctx, storage.UploadInput{
			File:        *car.CarRight,
			Name:        fileName,
			Folder: 	"drivers",
		})
		if err != nil {
			return err
		}
		car.CarRight = &uploadedFileName
	}
	if car.CarFrontRow != nil && *car.CarFrontRow != "" {
		fileName, err := utils.GenerateFileName()
		if err != nil {
			return err
		}
		uploadedFileName, err := s.fileStorage.Upload(ctx, storage.UploadInput{
			File:        *car.CarFrontRow,
			Name:        fileName,
			Folder: 	"drivers",
		})
		if err != nil {
			return err
		}
		car.CarFrontRow = &uploadedFileName
	}
	if car.CarFrontBack != nil && *car.CarFrontBack != "" {
		fileName, err := utils.GenerateFileName()
		if err != nil {
			return err
		}
		uploadedFileName, err := s.fileStorage.Upload(ctx, storage.UploadInput{
			File:        *car.CarFrontBack,
			Name:        fileName,
			Folder: 	"drivers",
		})
		if err != nil {
			return err
		}
		car.CarFrontBack = &uploadedFileName
	}
	if car.CarBaggage != nil && *car.CarBaggage != "" {
		fileName, err := utils.GenerateFileName()
		if err != nil {
			return err
		}
		uploadedFileName, err := s.fileStorage.Upload(ctx, storage.UploadInput{
			File:        *car.CarBaggage,
			Name:        fileName,
			Folder: 	"drivers",
		})
		if err != nil {
			return err
		}
		car.CarBaggage = &uploadedFileName
	}
	return s.repo.UpdateDriverCar(car, userId)
}