package session

import "GemmaORM/log"

// Begin 开启事务
func (s *Session) Begin() error {
	log.Info("transaction begin")
	var err error
	s.tx, err = s.db.Begin()
	if err != nil {
		log.Error(err)
	}
	return err
}

// Commit 提交事务
func (s *Session) Commit() error {
	log.Info("transaction commit")
	if err := s.tx.Commit(); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// RollBack 回滚事务
func (s *Session) RollBack() error {
	log.Info("transaction rollback")
	if err := s.tx.Rollback(); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
