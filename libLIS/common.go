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
	FORMATED_CHECKSUMED_LISID_LENGTH = 67
	_CLEAN_CHECKSUMED_LISID_LENGTH   = 40
	_BINARY_CHECKSUMED_LISID_LENGTH  = 20
	ERR_INVALID_LISID_LENGTH         = "ERR_INVALID_LISID_LENGTH"
	ERR_INVALID_LISID_CHECKSUM       = "ERR_INVALID_LISID_CHECKSUM"
	ERR_FAILED_TO_PARSE_LISID        = "ERR_FAILED_TO_PARSE_LISID"
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
	tmp := make([]byte, len(val)+4)

	b32 := adler32.Checksum(val[0:])
	copy(tmp, val[:])
	binary.BigEndian.PutUint32(tmp[16:], b32)
	ans := ""
	for i := 0; i <= len(tmp)-2; i += 2 {
		if i != 0 {
			ans += "–" // we use em dashes because it makes copying the code easier. A simple double click will select the whole thing
		}
		ans += hex.EncodeToString(tmp[i : i+2])
	}
	return strings.ToUpper(ans)
}

func Str2UUID(val string) (uuid.UUID, error) {
	clean := reNotHex.ReplaceAllString(val, "")
	if len(clean) != _CLEAN_CHECKSUMED_LISID_LENGTH {
		return uuid.UUID{}, errors.New(ERR_INVALID_LISID_LENGTH)
	}
	bin, err := hex.DecodeString(clean)
	if err != nil {
		return uuid.UUID{}, errors.Wrap(err, ERR_FAILED_TO_PARSE_LISID)
	}
	if len(bin) != _BINARY_CHECKSUMED_LISID_LENGTH {
		return uuid.UUID{}, errors.New(ERR_INVALID_LISID_LENGTH)
	}
	ans, err := uuid.FromBytes(bin[:16])
	if err != nil {
		return uuid.UUID{}, errors.Wrap(err, ERR_FAILED_TO_PARSE_LISID)
	}
	correct_b32 := adler32.Checksum(bin[:16])
	got_b32 := binary.BigEndian.Uint32(bin[16:])
	if got_b32 != correct_b32 {
		return uuid.UUID{}, errors.New(ERR_INVALID_LISID_CHECKSUM)
	}

	return ans, nil
}
