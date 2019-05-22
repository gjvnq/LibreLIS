package libLIS

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUUID2Str(t *testing.T) {
	my_uuid := uuid.UUID{}
	assert.Equal(t, "0000_0000_0000_0000_0000_0000_0000_0000_0010", UUID2Str(my_uuid))
	my_uuid = uuid.MustParse("22792807-048a-4158-9971-ce73d67c4f22")
	assert.Equal(t, "2279_2807_048A_4158_9971_CE73_D67C_4F22_2C10", UUID2Str(my_uuid))
}

func TestStr2UUID(t *testing.T) {
	right := uuid.UUID{}
	got, err := Str2UUID("0000_0000_0000_0000_0000_0000_0000_0000_0010")
	require.Nil(t, err)
	assert.Equal(t, right, got)
	right = uuid.MustParse("22792807-048a-4158-9971-ce73d67c4f22")
	got, err = Str2UUID("  2279_2807_0-48!A_4158_997   1_CE73_D67C_4F22_2C10  ")
	require.Nil(t, err)
	assert.Equal(t, right, got)

	right = uuid.MustParse("22792807-048a-4158-9971-ce73d67c4f22")
	got, err = Str2UUID("2279_2807_048A_4158_9971_CE73_D67C_4F22_2C1")
	require.NotNil(t, err)
	got, err = Str2UUID("2279_2807_048A_4158_9971_CE73_D67C_4F22_2C11")
	require.NotNil(t, err)
}
