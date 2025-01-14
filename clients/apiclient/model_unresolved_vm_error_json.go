/*
Wasp API

REST API for the Wasp node

API version: 0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package apiclient

import (
	"encoding/json"
)

// checks if the UnresolvedVMErrorJSON type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &UnresolvedVMErrorJSON{}

// UnresolvedVMErrorJSON struct for UnresolvedVMErrorJSON
type UnresolvedVMErrorJSON struct {
	Code *string `json:"code,omitempty"`
	Params []string `json:"params,omitempty"`
}

// NewUnresolvedVMErrorJSON instantiates a new UnresolvedVMErrorJSON object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUnresolvedVMErrorJSON() *UnresolvedVMErrorJSON {
	this := UnresolvedVMErrorJSON{}
	return &this
}

// NewUnresolvedVMErrorJSONWithDefaults instantiates a new UnresolvedVMErrorJSON object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUnresolvedVMErrorJSONWithDefaults() *UnresolvedVMErrorJSON {
	this := UnresolvedVMErrorJSON{}
	return &this
}

// GetCode returns the Code field value if set, zero value otherwise.
func (o *UnresolvedVMErrorJSON) GetCode() string {
	if o == nil || IsNil(o.Code) {
		var ret string
		return ret
	}
	return *o.Code
}

// GetCodeOk returns a tuple with the Code field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UnresolvedVMErrorJSON) GetCodeOk() (*string, bool) {
	if o == nil || IsNil(o.Code) {
		return nil, false
	}
	return o.Code, true
}

// HasCode returns a boolean if a field has been set.
func (o *UnresolvedVMErrorJSON) HasCode() bool {
	if o != nil && !IsNil(o.Code) {
		return true
	}

	return false
}

// SetCode gets a reference to the given string and assigns it to the Code field.
func (o *UnresolvedVMErrorJSON) SetCode(v string) {
	o.Code = &v
}

// GetParams returns the Params field value if set, zero value otherwise.
func (o *UnresolvedVMErrorJSON) GetParams() []string {
	if o == nil || IsNil(o.Params) {
		var ret []string
		return ret
	}
	return o.Params
}

// GetParamsOk returns a tuple with the Params field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UnresolvedVMErrorJSON) GetParamsOk() ([]string, bool) {
	if o == nil || IsNil(o.Params) {
		return nil, false
	}
	return o.Params, true
}

// HasParams returns a boolean if a field has been set.
func (o *UnresolvedVMErrorJSON) HasParams() bool {
	if o != nil && !IsNil(o.Params) {
		return true
	}

	return false
}

// SetParams gets a reference to the given []string and assigns it to the Params field.
func (o *UnresolvedVMErrorJSON) SetParams(v []string) {
	o.Params = v
}

func (o UnresolvedVMErrorJSON) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o UnresolvedVMErrorJSON) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Code) {
		toSerialize["code"] = o.Code
	}
	if !IsNil(o.Params) {
		toSerialize["params"] = o.Params
	}
	return toSerialize, nil
}

type NullableUnresolvedVMErrorJSON struct {
	value *UnresolvedVMErrorJSON
	isSet bool
}

func (v NullableUnresolvedVMErrorJSON) Get() *UnresolvedVMErrorJSON {
	return v.value
}

func (v *NullableUnresolvedVMErrorJSON) Set(val *UnresolvedVMErrorJSON) {
	v.value = val
	v.isSet = true
}

func (v NullableUnresolvedVMErrorJSON) IsSet() bool {
	return v.isSet
}

func (v *NullableUnresolvedVMErrorJSON) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUnresolvedVMErrorJSON(val *UnresolvedVMErrorJSON) *NullableUnresolvedVMErrorJSON {
	return &NullableUnresolvedVMErrorJSON{value: val, isSet: true}
}

func (v NullableUnresolvedVMErrorJSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUnresolvedVMErrorJSON) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


