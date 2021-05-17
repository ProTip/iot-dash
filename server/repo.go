package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
)

type AppRepo struct {
	*sql.DB
}

const (
	PlanStandard   = "STD"
	PlanEnterprise = "ENT"
)

type Account struct {
	Id            string
	Plan          string
	IotUserLimit  int
	AdminUsername string
	AdminPassword string
	count         int
}

type BearerToken struct {
	AccountId string
	Token     string
}

func (r AppRepo) AccountGetById(id string) {
	stmt, err := r.DB.Prepare("SELECT * from account WHERE id=?")
	if err != nil {
		log.Fatal(err)
	}

	acct, err := scanAccount(stmt.QueryRow(id))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(acct)
}

func (r AppRepo) AccountGetByUsername(name string) (*Account, error) {
	stmt, err := r.DB.Prepare("SELECT * from account  WHERE admin_username=?")
	if err != nil {
		return nil, err
	}

	acct, err := scanAccount(stmt.QueryRow(name))
	if err != nil {
		return nil, err
	}

	fmt.Println(acct)
	return acct, err
}

func (r AppRepo) AccountGetByBearerToken(token string) (*Account, error) {
	sum := sha256.Sum256([]byte(token))

	hashed := hex.EncodeToString(sum[:])

	stmt, err := r.DB.Prepare("SELECT account.* from account JOIN bearer_token token ON token.account_id = account.id  WHERE token.token=?")
	if err != nil {
		return nil, err
	}

	acct, err := scanAccount(stmt.QueryRow(hashed))
	if err != nil {
		return nil, err
	} else {
		return acct, err
	}
}

func (r AppRepo) AccountGetIotUserCount(account string) (int, error) {
	count := 0

	stmt, err := r.DB.Prepare("SELECT count from account JOIN iot_user user ON user.account_id = account.id  WHERE account.id=?")
	if err != nil {
		return count, err
	}

	err = stmt.QueryRow(account).Scan(&count)
	return count, err
}

// Registers IOT user login and increments user count if new. Returns "limit reached" error if plan limit prevents insert.
func (r AppRepo) AccountRegisterIotUser(account string, user string) (err error) {
	tx, _ := r.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Checks limit and insets in a single statement
	stmt, _ := tx.Prepare(`
INSERT OR ABORT INTO iot_user
SELECT ? acct, ? user
WHERE EXISTS(SELECT 1 FROM account WHERE id = acct AND count < iot_user_limit);
	`)

	var res sql.Result
	if res, err = stmt.Exec(account, user); err != nil {
		// TODO handle duplicates more gracefully
		panic(err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		// TODO Custom errors or better response
		err = errors.New("limit reached")
		panic(err)
	}

	stmt, _ = tx.Prepare("UPDATE account SET count = count + 1 WHERE id = ?")
	if _, err = stmt.Exec(account); err != nil {
		panic(err)
	}

	tx.Commit()

	return nil
}

// Upgrades account plan and return any SQL errors. This is "idempotent".
func (r AppRepo) AccountUpgrade(account string) error {
	stmt, err := r.DB.Prepare("UPDATE account SET plan = ?, iot_user_limit = ? WHERE id = ?")
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(PlanEnterprise, 100, account); err != nil {
		return err
	} else {
		return nil
	}
}

func scanAccount(row *sql.Row) (*Account, error) {
	acct := &Account{}
	err := row.Scan(&acct.Id, &acct.Plan, &acct.IotUserLimit, &acct.AdminUsername, &acct.AdminPassword, &acct.count)
	return acct, err
}
