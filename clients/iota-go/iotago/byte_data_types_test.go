package iotago_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/clients/iota-go/iotago"
)

func TestSerialization(t *testing.T) {
	hexStr := "0x12333aabcc"

	hexdata, err := iotago.NewHexData(hexStr)
	require.Nil(t, err)
	require.Equal(t, hexStr, hexdata.String())

	base64data := iotago.Bytes(hexdata.Data()).GetBase64Data()
	base64Str := base64data.String()

	t.Log(base64Str)
	t.Log(hexStr)
}

func TestHexdataJson(t *testing.T) {
	hexdata, err := iotago.NewHexData("0x12333aabcc")
	require.Nil(t, err)

	dataJson, err := json.Marshal(hexdata)
	require.Nil(t, err)

	hexdata2 := iotago.HexData{}
	err = json.Unmarshal(dataJson, &hexdata2)
	require.Nil(t, err)
	require.Equal(t, hexdata.Data(), hexdata2.Data())

	base64data := iotago.Bytes(hexdata.Data()).GetBase64Data()
	dataJsonb, err := json.Marshal(base64data)
	require.Nil(t, err)

	base64data2 := iotago.Base64Data{}
	err = json.Unmarshal(dataJsonb, &base64data2)
	require.Nil(t, err)
	require.Equal(t, base64data.Data(), base64data2.Data())
	require.Equal(t, hexdata.Data(), base64data2.Data())
}
