package service

import (
	"abir/models"
	"abir/pkg/repository"
	"abir/pkg/utils"
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type CreditCardsService struct {
	repo        repository.CreditCards
	redisClient *redis.Client
	fcmClient   *utils.FCMClient
}

func NewCreditCardsService(repo repository.CreditCards, redisClient *redis.Client, fcmClient *utils.FCMClient) *CreditCardsService {
	return &CreditCardsService{repo: repo, redisClient: redisClient, fcmClient: fcmClient}
}

func (s *CreditCardsService) Get(userId int) ([]models.CreditCards, error) {
	return s.repo.Get(userId)
}
func (s *CreditCardsService) Store(creditCard models.CreditCards, userId int) (int, error) {
	return s.repo.Store(creditCard, userId)
}
func (s *CreditCardsService) SendActivationCode(creditCardId, userId int) (string, error) {
	codeMin := 10000
	codeMax := 99999
	code := strconv.Itoa(rand.Intn(codeMax-codeMin) + codeMin)
	code = strconv.Itoa(11111)

	_, ok := s.redisClient.Get("card_activation" + strconv.Itoa(creditCardId)).Result()
	if ok == nil {
		return "", errors.New("try after a while")
	}
	err := s.redisClient.Set("card_activation"+strconv.Itoa(creditCardId), code, 2*time.Minute).Err()
	if err != nil {
		return "", err
	}
	creditCard, err := s.repo.GetSingleCard(creditCardId, userId)
	if err != nil {
		return "", err
	}
	if creditCard.CardNumber == nil || creditCard.CardExpiration == nil {
		return "", errors.New("card not found")
	}
	card, err := utils.GetCardToken(*creditCard.CardNumber, *creditCard.CardExpiration)
	if err != nil {
		return "", err
	}
	phone := card.Card.Phone
	if len(phone) != 12 {
		phone = "998" + phone
	}
	phone = utils.HidePhone(phone)
	//err = utils.SendSms(login, "Your verification code - "+code)
	//if err != nil {
	//	return err
	//}
	return phone, nil
}
func (s *CreditCardsService) Activate(creditCardId, userId int, code string) error {
	activationCode, err := s.redisClient.Get("card_activation" + strconv.Itoa(creditCardId)).Result()
	if err != nil {
		return err
	}
	if activationCode != code {
		return errors.New("wrong code from sms")
	}
	return s.repo.Activate(creditCardId, userId)
}
func (s *CreditCardsService) Delete(creditCardId, userId int) error {
	return s.repo.Delete(creditCardId, userId)
}
