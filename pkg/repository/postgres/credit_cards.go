package postgres

import (
	"abir/models"
	"abir/pkg/utils"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type CreditCardsPostgres struct {
	db *sqlx.DB
}

func NewCreditCardsPostgres(db *sqlx.DB) *CreditCardsPostgres {
	return &CreditCardsPostgres{db: db}
}


func (r *CreditCardsPostgres) Get(userId int) ([]models.CreditCards, error) {
	var lists []models.CreditCards
	query := fmt.Sprintf("SELECT id, card_info, is_active, is_main FROM %s WHERE user_id = $1 ORDER BY id ASC", creditCardsTable)
	err := r.db.Select(&lists, query, userId)
	for i, list := range lists {
		cardInfo, err := utils.DecryptMessage(*list.CardInfo)
		if err != nil {
			return []models.CreditCards{}, err
		}
		cardInfoArr := strings.Split(cardInfo, ":")
		if len(cardInfoArr) != 2 {
			return []models.CreditCards{}, errors.New("invalid cardInfo")
		}
		list.CardNumber = &cardInfoArr[0]
		list.CardExpiration = &cardInfoArr[1]
		lists[i] = list
	}
	return lists, err
}
func (r *CreditCardsPostgres) GetSingleCard(creditCardId, userId int) (models.CreditCards, error){
	var creditCard models.CreditCards
	usrQuery := fmt.Sprintf("SELECT card_info FROM %s WHERE id=$1 AND user_id=$2", creditCardsTable)
	err := r.db.Get(&creditCard, usrQuery, creditCardId, userId)
	if err != nil {
		return models.CreditCards{}, err
	}
	cardInfo, err := utils.DecryptMessage(*creditCard.CardInfo)
	if err != nil {
		return models.CreditCards{}, err
	}
	cardInfoArr := strings.Split(cardInfo, ":")
	if len(cardInfoArr) != 2 {
		return models.CreditCards{}, errors.New("invalid cardInfo")
	}
	creditCard.CardNumber = &cardInfoArr[0]
	creditCard.CardExpiration = &cardInfoArr[1]
	return creditCard, nil
}

func (r *CreditCardsPostgres) Store(creditCard models.CreditCards, userId int) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (card_info, user_id) SELECT $1,$2 WHERE NOT EXISTS (SELECT id FROM %s WHERE card_info = $3 AND user_id = $2) RETURNING id", creditCardsTable, creditCardsTable)
	cardInfo, err := utils.EncryptMessage(fmt.Sprintf("%s:%s", *creditCard.CardNumber, *creditCard.CardExpiration))
	if err != nil {
		return 0, err
	}
	var id int
	row := r.db.QueryRow(query, cardInfo, userId, cardInfo)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows{
			return 0, errors.New("already exists")
		}
		return 0, err
	}
	return id, err
}
func (r *CreditCardsPostgres) Activate(creditCardId, userId int) error {
	updateQuery := fmt.Sprintf(`UPDATE %s SET is_active = $1 WHERE id = $2 AND user_id = $3`,
		creditCardsTable)
	_, err := r.db.Exec(updateQuery, true, creditCardId, userId)
	return err
}
func (r *CreditCardsPostgres) Delete(creditCardId, userId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1 AND user_id = $2`,
		creditCardsTable)
	_, err := r.db.Exec(query, creditCardId, userId)
	return err
}
