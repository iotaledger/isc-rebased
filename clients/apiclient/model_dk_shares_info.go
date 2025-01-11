/*
Wasp API

REST API for the Wasp node

API version: 0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package apiclient

import (
	"encoding/json"
	"bytes"
	"fmt"
)

// checks if the DKSharesInfo type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &DKSharesInfo{}

// DKSharesInfo struct for DKSharesInfo
type DKSharesInfo struct {
	// New generated shared address.
	Address string `json:"address"`
	// Identities of the nodes sharing the key. (Hex)
	PeerIdentities []string `json:"peerIdentities"`
	PeerIndex uint32 `json:"peerIndex"`
	// Used public key. (Hex)
	PublicKey string `json:"publicKey"`
	// Public key shares for all the peers. (Hex)
	PublicKeyShares []string `json:"publicKeyShares"`
	Threshold uint32 `json:"threshold"`
}

type _DKSharesInfo DKSharesInfo

// NewDKSharesInfo instantiates a new DKSharesInfo object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDKSharesInfo(address string, peerIdentities []string, peerIndex uint32, publicKey string, publicKeyShares []string, threshold uint32) *DKSharesInfo {
	this := DKSharesInfo{}
	this.Address = address
	this.PeerIdentities = peerIdentities
	this.PeerIndex = peerIndex
	this.PublicKey = publicKey
	this.PublicKeyShares = publicKeyShares
	this.Threshold = threshold
	return &this
}

// NewDKSharesInfoWithDefaults instantiates a new DKSharesInfo object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDKSharesInfoWithDefaults() *DKSharesInfo {
	this := DKSharesInfo{}
	return &this
}

// GetAddress returns the Address field value
func (o *DKSharesInfo) GetAddress() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Address
}

// GetAddressOk returns a tuple with the Address field value
// and a boolean to check if the value has been set.
func (o *DKSharesInfo) GetAddressOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Address, true
}

// SetAddress sets field value
func (o *DKSharesInfo) SetAddress(v string) {
	o.Address = v
}

// GetPeerIdentities returns the PeerIdentities field value
func (o *DKSharesInfo) GetPeerIdentities() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.PeerIdentities
}

// GetPeerIdentitiesOk returns a tuple with the PeerIdentities field value
// and a boolean to check if the value has been set.
func (o *DKSharesInfo) GetPeerIdentitiesOk() ([]string, bool) {
	if o == nil {
		return nil, false
	}
	return o.PeerIdentities, true
}

// SetPeerIdentities sets field value
func (o *DKSharesInfo) SetPeerIdentities(v []string) {
	o.PeerIdentities = v
}

// GetPeerIndex returns the PeerIndex field value
func (o *DKSharesInfo) GetPeerIndex() uint32 {
	if o == nil {
		var ret uint32
		return ret
	}

	return o.PeerIndex
}

// GetPeerIndexOk returns a tuple with the PeerIndex field value
// and a boolean to check if the value has been set.
func (o *DKSharesInfo) GetPeerIndexOk() (*uint32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PeerIndex, true
}

// SetPeerIndex sets field value
func (o *DKSharesInfo) SetPeerIndex(v uint32) {
	o.PeerIndex = v
}

// GetPublicKey returns the PublicKey field value
func (o *DKSharesInfo) GetPublicKey() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.PublicKey
}

// GetPublicKeyOk returns a tuple with the PublicKey field value
// and a boolean to check if the value has been set.
func (o *DKSharesInfo) GetPublicKeyOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PublicKey, true
}

// SetPublicKey sets field value
func (o *DKSharesInfo) SetPublicKey(v string) {
	o.PublicKey = v
}

// GetPublicKeyShares returns the PublicKeyShares field value
func (o *DKSharesInfo) GetPublicKeyShares() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.PublicKeyShares
}

// GetPublicKeySharesOk returns a tuple with the PublicKeyShares field value
// and a boolean to check if the value has been set.
func (o *DKSharesInfo) GetPublicKeySharesOk() ([]string, bool) {
	if o == nil {
		return nil, false
	}
	return o.PublicKeyShares, true
}

// SetPublicKeyShares sets field value
func (o *DKSharesInfo) SetPublicKeyShares(v []string) {
	o.PublicKeyShares = v
}

// GetThreshold returns the Threshold field value
func (o *DKSharesInfo) GetThreshold() uint32 {
	if o == nil {
		var ret uint32
		return ret
	}

	return o.Threshold
}

// GetThresholdOk returns a tuple with the Threshold field value
// and a boolean to check if the value has been set.
func (o *DKSharesInfo) GetThresholdOk() (*uint32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Threshold, true
}

// SetThreshold sets field value
func (o *DKSharesInfo) SetThreshold(v uint32) {
	o.Threshold = v
}

func (o DKSharesInfo) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o DKSharesInfo) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["address"] = o.Address
	toSerialize["peerIdentities"] = o.PeerIdentities
	toSerialize["peerIndex"] = o.PeerIndex
	toSerialize["publicKey"] = o.PublicKey
	toSerialize["publicKeyShares"] = o.PublicKeyShares
	toSerialize["threshold"] = o.Threshold
	return toSerialize, nil
}

func (o *DKSharesInfo) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"address",
		"peerIdentities",
		"peerIndex",
		"publicKey",
		"publicKeyShares",
		"threshold",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err;
	}

	for _, requiredProperty := range(requiredProperties) {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varDKSharesInfo := _DKSharesInfo{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varDKSharesInfo)

	if err != nil {
		return err
	}

	*o = DKSharesInfo(varDKSharesInfo)

	return err
}

type NullableDKSharesInfo struct {
	value *DKSharesInfo
	isSet bool
}

func (v NullableDKSharesInfo) Get() *DKSharesInfo {
	return v.value
}

func (v *NullableDKSharesInfo) Set(val *DKSharesInfo) {
	v.value = val
	v.isSet = true
}

func (v NullableDKSharesInfo) IsSet() bool {
	return v.isSet
}

func (v *NullableDKSharesInfo) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDKSharesInfo(val *DKSharesInfo) *NullableDKSharesInfo {
	return &NullableDKSharesInfo{value: val, isSet: true}
}

func (v NullableDKSharesInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDKSharesInfo) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


