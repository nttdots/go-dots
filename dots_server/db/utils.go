package db

import (
  "database/sql"
  "time"

  "github.com/go-xorm/xorm"
  log "github.com/sirupsen/logrus"

  "github.com/nttdots/go-dots/dots_server/models"
)

type Tx struct {
  Engine  *xorm.Engine
  Session *xorm.Session
}

func WithTransaction(f func(*Tx) (interface{}, error)) (interface{}, error) {
  engine, err := models.ConnectDB()
  if err != nil {
    log.WithError(err).Error("Failed connect to database.")
    return nil, err
  }

  session := engine.NewSession()
  defer session.Close()

  err = session.Begin()
  if err != nil {
    log.WithError(err).Error("Failed begin transaction.")
    return nil, err
  }

  ret, err := f(&Tx{ engine, session })
  if err != nil {
    session.Rollback()
    return nil, err
  }

  err = session.Commit()
  if err != nil {
    log.WithError(err).Error("Failed to commit transaction.")
    session.Rollback()
    return nil, err
  }
  return ret, nil
}

func AsNullString(p *string) sql.NullString {
  if p == nil {
    return sql.NullString { Valid: false }
  } else {
    return sql.NullString { Valid: true, String: *p }
  }
}

func AsStringPointer(s sql.NullString) *string {
  if s.Valid {
    return &s.String
  } else {
    return nil
  }
}

func AsDateTime(t time.Time) string {
  return t.Format("2006-01-02 15:04:05")
}
