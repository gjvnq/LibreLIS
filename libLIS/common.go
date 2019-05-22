package libLIS

import (
	"encoding/binary"
	"encoding/hex"
	"hash/adler32"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var SysUUID uuid.UUID

var reNotHex *regexp.Regexp

const (
	ERR_INVALID_UUID_LENGTH  = "ERR_INVALID_UUID_LENGTH"
	ERR_FAILED_TO_PARSE_UUID = "ERR_FAILED_TO_PARSE_UUID"
)

func init() {
	reNotHex = regexp.MustCompile(`[^A-F0-9a-f]*`)
}

func Start() {
	var err error

	err = getSysUUID()
	if err != nil {
		panic(err)
	}
}

func getSysUUID() error {
	err := DB.QueryRow("SELECT `user_id` FROM `user` WHERE `user_type` = ?", "S").Scan(&SysUUID)
	if err == nil {
		sys_user := User{}
		sys_user.type_ = USER_TYPE_SYSTEM
	}
	return err
}

func UUID2Str(val uuid.UUID) string {
	tmp := make([]byte, 0)

	b32 := adler32.Checksum(val[0:])

	tmp = append(tmp, val[0:]...)
	tmp = append(tmp, 0, 0, 0, 0)
	binary.BigEndian.PutUint32(tmp[len(val):], b32)
	ans := ""
	for i := 0; i < len(tmp)-2; i += 2 {
		if i != 0 {
			ans += "_" // we use underscores because it makes copying the code easier. A simple double click will select the whole thing
		}
		ans += hex.EncodeToString(tmp[i : i+2])
	}
	return strings.ToUpper(ans)
}

func Str2UUID(val string) (uuid.UUID, error) {
	// remove all character expect 0123456789ABCDEF
	clean := reNotHex.ReplaceAllString(val, "")
	if len(clean) != 36 {
		return uuid.UUID{}, errors.New(ERR_INVALID_UUID_LENGTH)
	}
	bin, err := hex.DecodeString(clean)
	if err != nil {
		return uuid.UUID{}, errors.Wrap(err, ERR_FAILED_TO_PARSE_UUID)
	}
	ans, err := uuid.FromBytes(bin[:16])
	if err != nil {
		return uuid.UUID{}, errors.Wrap(err, ERR_FAILED_TO_PARSE_UUID)
	}
	// check adler32

	return ans, nil
}
