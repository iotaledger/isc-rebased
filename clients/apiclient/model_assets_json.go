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

// checks if the AssetsJSON type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &AssetsJSON{}

// AssetsJSON struct for AssetsJSON
type AssetsJSON struct {
	Coins []CoinJSON `json:"coins"`
	Objects [][]int32 `json:"objects"`
}

type _AssetsJSON AssetsJSON

// NewAssetsJSON instantiates a new AssetsJSON object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAssetsJSON(coins []CoinJSON, objects [][]int32) *AssetsJSON {
	this := AssetsJSON{}
	this.Coins = coins
	this.Objects = objects
	return &this
}

// NewAssetsJSONWithDefaults instantiates a new AssetsJSON object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAssetsJSONWithDefaults() *AssetsJSON {
	this := AssetsJSON{}
	return &this
}

// GetCoins returns the Coins field value
func (o *AssetsJSON) GetCoins() []CoinJSON {
	if o == nil {
		var ret []CoinJSON
		return ret
	}

	return o.Coins
}

// GetCoinsOk returns a tuple with the Coins field value
// and a boolean to check if the value has been set.
func (o *AssetsJSON) GetCoinsOk() ([]CoinJSON, bool) {
	if o == nil {
		return nil, false
	}
	return o.Coins, true
}

// SetCoins sets field value
func (o *AssetsJSON) SetCoins(v []CoinJSON) {
	o.Coins = v
}

// GetObjects returns the Objects field value
func (o *AssetsJSON) GetObjects() [][]int32 {
	if o == nil {
		var ret [][]int32
		return ret
	}

	return o.Objects
}

// GetObjectsOk returns a tuple with the Objects field value
// and a boolean to check if the value has been set.
func (o *AssetsJSON) GetObjectsOk() ([][]int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.Objects, true
}

// SetObjects sets field value
func (o *AssetsJSON) SetObjects(v [][]int32) {
	o.Objects = v
}

func (o AssetsJSON) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o AssetsJSON) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["coins"] = o.Coins
	toSerialize["objects"] = o.Objects
	return toSerialize, nil
}

func (o *AssetsJSON) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"coins",
		"objects",
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

	varAssetsJSON := _AssetsJSON{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varAssetsJSON)

	if err != nil {
		return err
	}

	*o = AssetsJSON(varAssetsJSON)

	return err
}

type NullableAssetsJSON struct {
	value *AssetsJSON
	isSet bool
}

func (v NullableAssetsJSON) Get() *AssetsJSON {
	return v.value
}

func (v *NullableAssetsJSON) Set(val *AssetsJSON) {
	v.value = val
	v.isSet = true
}

func (v NullableAssetsJSON) IsSet() bool {
	return v.isSet
}

func (v *NullableAssetsJSON) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAssetsJSON(val *AssetsJSON) *NullableAssetsJSON {
	return &NullableAssetsJSON{value: val, isSet: true}
}

func (v NullableAssetsJSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAssetsJSON) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


