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

// checks if the CoinJSON type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CoinJSON{}

// CoinJSON struct for CoinJSON
type CoinJSON struct {
	// The balance (uint64 as string)
	Balance string `json:"balance"`
	CoinType Type `json:"coinType"`
}

type _CoinJSON CoinJSON

// NewCoinJSON instantiates a new CoinJSON object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCoinJSON(balance string, coinType Type) *CoinJSON {
	this := CoinJSON{}
	this.Balance = balance
	this.CoinType = coinType
	return &this
}

// NewCoinJSONWithDefaults instantiates a new CoinJSON object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCoinJSONWithDefaults() *CoinJSON {
	this := CoinJSON{}
	return &this
}

// GetBalance returns the Balance field value
func (o *CoinJSON) GetBalance() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Balance
}

// GetBalanceOk returns a tuple with the Balance field value
// and a boolean to check if the value has been set.
func (o *CoinJSON) GetBalanceOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Balance, true
}

// SetBalance sets field value
func (o *CoinJSON) SetBalance(v string) {
	o.Balance = v
}

// GetCoinType returns the CoinType field value
func (o *CoinJSON) GetCoinType() Type {
	if o == nil {
		var ret Type
		return ret
	}

	return o.CoinType
}

// GetCoinTypeOk returns a tuple with the CoinType field value
// and a boolean to check if the value has been set.
func (o *CoinJSON) GetCoinTypeOk() (*Type, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CoinType, true
}

// SetCoinType sets field value
func (o *CoinJSON) SetCoinType(v Type) {
	o.CoinType = v
}

func (o CoinJSON) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CoinJSON) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["balance"] = o.Balance
	toSerialize["coinType"] = o.CoinType
	return toSerialize, nil
}

func (o *CoinJSON) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"balance",
		"coinType",
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

	varCoinJSON := _CoinJSON{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varCoinJSON)

	if err != nil {
		return err
	}

	*o = CoinJSON(varCoinJSON)

	return err
}

type NullableCoinJSON struct {
	value *CoinJSON
	isSet bool
}

func (v NullableCoinJSON) Get() *CoinJSON {
	return v.value
}

func (v *NullableCoinJSON) Set(val *CoinJSON) {
	v.value = val
	v.isSet = true
}

func (v NullableCoinJSON) IsSet() bool {
	return v.isSet
}

func (v *NullableCoinJSON) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCoinJSON(val *CoinJSON) *NullableCoinJSON {
	return &NullableCoinJSON{value: val, isSet: true}
}

func (v NullableCoinJSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCoinJSON) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


