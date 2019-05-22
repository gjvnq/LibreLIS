package libLIS

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	hibp "github.com/mattevans/pwned-passwords"
	"github.com/pkg/errors"
	"github.com/raja/argon2pw"
)

const (
	USER_TYPE_NATURAL_PERSON = "N"
	USER_TYPE_LEGAL_PERSON   = "L"
	USER_TYPE_BOT            = "B"
	USER_TYPE_SYSTEM         = "S"
)

const (
	MIN_PASSWORD_LENGTH                = 12
	ERR_PASSWORD_TOO_SHORT             = "ERR_PASSWORD_TOO_SHORT"
	ERR_SUPER_USER_PASSWORD_IN_USE     = "ERR_SUPER_USER_PASSWORD_IN_USE"
	ERR_FAILED_TO_VERIFY_PASSWORD_HIBP = "ERR_FAILED_TO_VERIFY_PASSWORD_HIBP"
	ERR_FAILED_TO_HASH_PASSWORD        = "ERR_FAILED_TO_HASH_PASSWORD"
)

type User struct {
	id          uuid.UUID
	revision    int
	name        string
	email       string
	flair       string
	password    string
	isSuperUser bool
	type_       string
	creation    time.Time
	changed     time.Time
	changer     uuid.UUID
}

var hibpClient *hibp.Client

func init() {
	hibpClient = hibp.NewClient()
}

func (this User) Id() string {
	return this.id.String()
}

func (this *User) String() string {
	ans := fmt.Sprintf("[%s %s №%d %s <%s>", this.id.String(), this.type_, this.revision, this.name, this.email)
	if this.isSuperUser {
		ans += " (SUPER)"
	}
	ans += "]"
	return ans
}

func (this *User) VerifyPassword(password string) bool {
	ans, err := argon2pw.CompareHashWithPassword(this.password, password)
	if err != nil && err.Error() != "Password did not match" {
		TheLogger.Error(err)
		return false
	}
	TheLogger.DebugF("Password verification for %s: %v", this.String(), ans)
	return ans
}

func (this *User) SetPassword(password string) error {
	var err error

	if len(password) < MIN_PASSWORD_LENGTH {
		return errors.New(ERR_PASSWORD_TOO_SHORT)
	}
	if this.isSuperUser {
		pwned, err := hibpClient.Pwned.Compromised(password)
		if err != nil {
			TheLogger.ErrorF("Failed to verify password for super user %s: %v", this.String(), err.Error())
			return errors.Wrap(err, ERR_FAILED_TO_VERIFY_PASSWORD_HIBP)
		}
		if pwned {
			return errors.New(ERR_SUPER_USER_PASSWORD_IN_USE)
		}
	}

	this.password, err = argon2pw.GenerateSaltedHash(password)
	if err != nil {
		TheLogger.ErrorF("Failed to hash password for user %s: %v", this.String(), err.Error())
		return errors.Wrap(err, ERR_FAILED_TO_HASH_PASSWORD)
	}
	return nil
}

func LoadUserByEmail(email string) *User {
	return nil
}

func LoadUserById(id int) *User {
	return nil
}
