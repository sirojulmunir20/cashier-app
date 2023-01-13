package repository

import (
	"a21hc3NpZ25tZW50/db"
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"fmt"
	"time"
)

type SessionsRepository struct {
	db db.DB
}

func NewSessionsRepository(db db.DB) SessionsRepository {
	return SessionsRepository{db}
}

func (u *SessionsRepository) ReadSessions() ([]model.Session, error) {
	records, err := u.db.Load("sessions")
	if err != nil {
		return nil, err
	}

	var listSessions []model.Session
	err = json.Unmarshal([]byte(records), &listSessions)
	if err != nil {
		return nil, err
	}

	return listSessions, nil
}

func (u *SessionsRepository) DeleteSessions(tokenTarget string) error {
	listSessions, err := u.ReadSessions()
	if err != nil {
		return err
	}

	dataSessions := []model.Session{}
	// Select target token and delete from listSessions
	for _, v := range listSessions {
		if v.Token == tokenTarget {
			continue
		} else {
			dataSessions = append(dataSessions, v)
		}
	}
	
	// for i := 0; i < len(listSessions); i++ {
	// 	if listSessions[i].Token == tokenTarget {
	// 		listSessions[i] = model.Session{}
	// 	}
	// }
	jsonData, err := json.Marshal(dataSessions)
	if err != nil {
		return err
	}

	err = u.db.Save("sessions", jsonData)
	if err != nil {
		return err
	}

	return nil
}

func (u *SessionsRepository) AddSessions(session model.Session) error {
	listSession, err := u.ReadSessions()
	if err != nil {
		return err
	}
	listSession = append(listSession, session)
	dataSession, err := json.Marshal(listSession)
	if err != nil {
		return err
	}
	_ = u.db.Save("sessions", dataSession)
	return nil // TODO: replace this
}

func (u *SessionsRepository) CheckExpireToken(token string) (model.Session, error) {
	expire, err := u.TokenExist(token)
	if err != nil {
		return model.Session{}, err
	}
	if !u.TokenExpired(expire) { // jika tdk expired
		return expire, nil
	}
	if u.TokenExpired(expire) {
		u.DeleteSessions(expire.Token)
		return model.Session{}, fmt.Errorf("Token is Expired!")
	}
	return model.Session{}, nil // TODO: replace this
}

func (u *SessionsRepository) ResetSessions() error {
	err := u.db.Reset("sessions", []byte("[]"))
	if err != nil {
		return err
	}

	return nil
}

func (u *SessionsRepository) TokenExist(req string) (model.Session, error) {
	listSessions, err := u.ReadSessions()
	if err != nil {
		return model.Session{}, err
	}
	for _, element := range listSessions {
		if element.Token == req {
			return element, nil
		}
	}
	return model.Session{}, fmt.Errorf("Token Not Found!")
}

func (u *SessionsRepository) TokenExpired(s model.Session) bool {
	return s.Expiry.Before(time.Now())
}
