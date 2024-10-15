package util

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/coin"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/tools/wasp-cli/cli/wallet"
	"github.com/iotaledger/wasp/tools/wasp-cli/log"
)

//nolint:funlen,gocyclo
func ValueFromString(vtype, s string, chainID isc.ChainID) []byte {
	switch strings.ToLower(vtype) {
	case "address":
		addr, err := cryptolib.NewAddressFromHexString(s)
		log.Check(err)

		return addr.Bytes()
	case "agentid":
		return AgentIDFromString(s, chainID).Bytes()
	case "bigint":
		n, ok := new(big.Int).SetString(s, 10)
		if !ok {
			log.Fatal("error converting to bigint")
		}
		return n.Bytes()
	case "bool":
		b, err := strconv.ParseBool(s)
		log.Check(err)
		return codec.Encode[bool](b)
	case "bytes", "hex":
		b, err := cryptolib.DecodeHex(s)
		log.Check(err)
		return b
	case "chainid":
		chainid, err := isc.ChainIDFromString(s)
		log.Check(err)
		return chainid.Bytes()
	case "dict":
		d := dict.Dict{}
		err := d.UnmarshalJSON([]byte(s))
		log.Check(err)
		return codec.Encode[dict.Dict](d)
	case "file":
		return ReadFile(s)
	case "hash":
		hash, err := hashing.HashValueFromHex(s)
		log.Check(err)
		return hash.Bytes()
	case "hname":
		hn, err := isc.HnameFromString(s)
		log.Check(err)
		return hn.Bytes()
	case "int8":
		n, err := strconv.ParseInt(s, 10, 8)
		log.Check(err)
		return codec.Encode[int8](int8(n))
	case "int16":
		n, err := strconv.ParseInt(s, 10, 16)
		log.Check(err)
		return codec.Encode[int16](int16(n))
	case "int32":
		n, err := strconv.ParseInt(s, 10, 32)
		log.Check(err)
		return codec.Encode[int32](int32(n))
	case "int64", "int":
		n, err := strconv.ParseInt(s, 10, 64)
		log.Check(err)
		return codec.Encode[int64](n)
	case "nftid":
		nid, err := cryptolib.DecodeHex(s)
		log.Check(err)
		if len(nid) != iotago.NFTIDLength {
			log.Fatal("invalid nftid length")
		}
		return nid
	case "requestid":
		rid, err := isc.RequestIDFromString(s)
		log.Check(err)
		return rid.Bytes()
	case "string":
		return []byte(s)
	case "tokenid":
		tid, err := cryptolib.DecodeHex(s)
		log.Check(err)
		if len(tid) != iotago.FoundryIDLength {
			log.Fatal("invalid tokenid length")
		}
		return tid
	case "uint8":
		n, err := strconv.ParseUint(s, 10, 8)
		log.Check(err)
		return codec.Encode[uint8](uint8(n))
	case "uint16":
		n, err := strconv.ParseUint(s, 10, 16)
		log.Check(err)
		return codec.Encode[uint16](uint16(n))
	case "uint32":
		n, err := strconv.ParseUint(s, 10, 32)
		log.Check(err)
		return codec.Encode[uint32](uint32(n))
	case "uint64":
		n, err := strconv.ParseUint(s, 10, 64)
		log.Check(err)
		return codec.Encode[uint64](n)
	}
	log.Fatalf("ValueFromString: No handler for type %s", vtype)
	return nil
}

//nolint:funlen,gocyclo
func ValueToString(vtype string, v []byte) string {
	switch strings.ToLower(vtype) {
	case "address":
		addr, err := codec.Decode[*cryptolib.Address](v)
		log.Check(err)
		return addr.String()
	case "agentid":
		aid, err := codec.Decode[isc.AgentID](v)
		log.Check(err)
		return aid.String()
	case "bigint":
		n := new(big.Int).SetBytes(v)
		return n.String()
	case "bool":
		b, err := codec.Decode[bool](v)
		log.Check(err)
		if b {
			return "true"
		}
		return "false"
	case "bytes", "hex":
		return iotago.EncodeHex(v)
	case "chainid":
		cid, err := codec.Decode[isc.ChainID](v)
		log.Check(err)
		return cid.String()
	case "dict":
		d, err := codec.Decode[dict.Dict](v)
		log.Check(err)
		s, err := d.MarshalJSON()
		log.Check(err)
		return string(s)
	case "hash":
		hash, err := codec.Decode[hashing.HashValue](v)
		log.Check(err)
		return hash.String()
	case "hname":
		hn, err := codec.Decode[isc.Hname](v)
		log.Check(err)
		return hn.String()
	case "int8":
		n, err := codec.Decode[int8](v)
		log.Check(err)
		return fmt.Sprintf("%d", n)
	case "int16":
		n, err := codec.Decode[int16](v)
		log.Check(err)
		return fmt.Sprintf("%d", n)
	case "int32":
		n, err := codec.Decode[int32](v)
		log.Check(err)
		return fmt.Sprintf("%d", n)
	case "int64", "int":
		n, err := codec.Decode[int64](v)
		log.Check(err)
		return fmt.Sprintf("%d", n)
	case "nftid":
		nid, err := codec.Decode[NFTID](v)
		log.Check(err)
		return nid.String()
	case "requestid":
		rid, err := codec.Decode[isc.RequestID](v)
		log.Check(err)
		return rid.String()
	case "string":
		return fmt.Sprintf("%q", string(v))
	case "tokenid":
		tid, err := codec.Decode[coin.Type](v)
		log.Check(err)
		return tid.String()
	case "uint8":
		n, err := codec.Decode[uint8](v)
		log.Check(err)
		return fmt.Sprintf("%d", n)
	case "uint16":
		n, err := codec.Decode[uint16](v)
		log.Check(err)
		return fmt.Sprintf("%d", n)
	case "uint32":
		n, err := codec.Decode[uint32](v)
		log.Check(err)
		return fmt.Sprintf("%d", n)
	case "uint64":
		n, err := codec.Decode[uint64](v)
		log.Check(err)
		return fmt.Sprintf("%d", n)
	}

	log.Fatalf("ValueToString: No handler for type %s", vtype)
	return ""
}

func EncodeParams(params []string, chainID isc.ChainID) dict.Dict {
	d := dict.New()
	if len(params)%4 != 0 {
		log.Fatal("Params format: <type> <key> <type> <value> ...")
	}
	for i := 0; i < len(params)/4; i++ {
		ktype := params[i*4]
		k := params[i*4+1]
		vtype := params[i*4+2]
		v := params[i*4+3]

		key := kv.Key(ValueFromString(ktype, k, chainID))
		val := ValueFromString(vtype, v, chainID)
		d.Set(key, val)
	}
	return d
}

func PrintDictAsJSON(d dict.Dict) {
	log.Check(json.NewEncoder(os.Stdout).Encode(d))
}

func UnmarshalDict() dict.Dict {
	var d dict.Dict
	log.Check(json.NewDecoder(os.Stdin).Decode(&d))
	return d
}

func AgentIDFromArgs(args []string, chainID isc.ChainID) isc.AgentID {
	if len(args) == 0 {
		return isc.NewAddressAgentID(wallet.Load().Address())
	}
	return AgentIDFromString(args[0], chainID)
}

func AgentIDFromString(s string, chainID isc.ChainID) isc.AgentID {
	if s == "common" {
		return accounts.CommonAccount()
	}
	// allow EVM addresses as AgentIDs without the chain specified
	if strings.HasPrefix(s, "0x") && !strings.Contains(s, isc.AgentIDStringSeparator) {
		s = s + isc.AgentIDStringSeparator + chainID.String()
	}
	agentID, err := isc.AgentIDFromString(s)
	log.Check(err, "cannot parse AgentID")
	return agentID
}
