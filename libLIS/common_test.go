package libLIS

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUUID2Str(t *testing.T) {
	my_uuid := uuid.UUID{}
	assert.Equal(t, "0000–0000–0000–0000–0000–0000–0000–0000–0010–0001", UUID2Str(my_uuid))
	my_uuid = uuid.MustParse("22792807-048a-4158-9971-ce73d67c4f22")
	assert.Equal(t, "2279–2807–048A–4158–9971–CE73–D67C–4F22–2C10–0600", UUID2Str(my_uuid))
}

func TestStr2UUID(t *testing.T) {
	right := uuid.UUID{}
	got, err := Str2UUID("0000_0000_0000_0000_0000_0000_0000_0000_0010_0001")
	require.Nil(t, err)
	assert.Equal(t, right, got)
	right = uuid.MustParse("22792807-048a-4158-9971-ce73d67c4f22")
	got, err = Str2UUID("  2279_2807_0-48!A_4158_99–7   1_CE73_D67C_4F22_2C10 0600	 ")
	require.Nil(t, err)
	assert.Equal(t, right, got)

	right = uuid.MustParse("22792807-048a-4158-9971-ce73d67c4f22")
	got, err = Str2UUID("2279–2807–048A–4158–9971–CE73–D67C–4F22–2C10–060")
	require.NotNil(t, err)
	require.Equal(t, ERR_INVALID_LISID_LENGTH, err.Error())
	require.Equal(t, uuid.UUID{}, got)

	right = uuid.MustParse("22792807-048a-4158-9971-ce73d67c4f22")
	got, err = Str2UUID("2279–2807–048A–4158–9971–CE73–D67C–4F22–2C10–0601")
	require.NotNil(t, err)
	require.Equal(t, ERR_INVALID_LISID_CHECKSUM, err.Error())
	require.Equal(t, uuid.UUID{}, got)
}
