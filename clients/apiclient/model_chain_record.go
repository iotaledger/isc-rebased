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

// checks if the ChainRecord type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ChainRecord{}

// ChainRecord struct for ChainRecord
type ChainRecord struct {
	AccessNodes []string `json:"accessNodes"`
	IsActive bool `json:"isActive"`
}

type _ChainRecord ChainRecord

// NewChainRecord instantiates a new ChainRecord object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewChainRecord(accessNodes []string, isActive bool) *ChainRecord {
	this := ChainRecord{}
	this.AccessNodes = accessNodes
	this.IsActive = isActive
	return &this
}

// NewChainRecordWithDefaults instantiates a new ChainRecord object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewChainRecordWithDefaults() *ChainRecord {
	this := ChainRecord{}
	return &this
}

// GetAccessNodes returns the AccessNodes field value
func (o *ChainRecord) GetAccessNodes() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.AccessNodes
}

// GetAccessNodesOk returns a tuple with the AccessNodes field value
// and a boolean to check if the value has been set.
func (o *ChainRecord) GetAccessNodesOk() ([]string, bool) {
	if o == nil {
		return nil, false
	}
	return o.AccessNodes, true
}

// SetAccessNodes sets field value
func (o *ChainRecord) SetAccessNodes(v []string) {
	o.AccessNodes = v
}

// GetIsActive returns the IsActive field value
func (o *ChainRecord) GetIsActive() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.IsActive
}

// GetIsActiveOk returns a tuple with the IsActive field value
// and a boolean to check if the value has been set.
func (o *ChainRecord) GetIsActiveOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.IsActive, true
}

// SetIsActive sets field value
func (o *ChainRecord) SetIsActive(v bool) {
	o.IsActive = v
}

func (o ChainRecord) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ChainRecord) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["accessNodes"] = o.AccessNodes
	toSerialize["isActive"] = o.IsActive
	return toSerialize, nil
}

func (o *ChainRecord) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"accessNodes",
		"isActive",
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

	varChainRecord := _ChainRecord{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varChainRecord)

	if err != nil {
		return err
	}

	*o = ChainRecord(varChainRecord)

	return err
}

type NullableChainRecord struct {
	value *ChainRecord
	isSet bool
}

func (v NullableChainRecord) Get() *ChainRecord {
	return v.value
}

func (v *NullableChainRecord) Set(val *ChainRecord) {
	v.value = val
	v.isSet = true
}

func (v NullableChainRecord) IsSet() bool {
	return v.isSet
}

func (v *NullableChainRecord) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableChainRecord(val *ChainRecord) *NullableChainRecord {
	return &NullableChainRecord{value: val, isSet: true}
}

func (v NullableChainRecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableChainRecord) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


