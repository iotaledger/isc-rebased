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

// checks if the Limits type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Limits{}

// Limits struct for Limits
type Limits struct {
	// The maximum gas per external view call
	MaxGasExternalViewCall int64 `json:"maxGasExternalViewCall"`
	// The maximum gas per block
	MaxGasPerBlock int64 `json:"maxGasPerBlock"`
	// The maximum gas per request
	MaxGasPerRequest int64 `json:"maxGasPerRequest"`
	// The minimum gas per request
	MinGasPerRequest int64 `json:"minGasPerRequest"`
}

type _Limits Limits

// NewLimits instantiates a new Limits object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewLimits(maxGasExternalViewCall int64, maxGasPerBlock int64, maxGasPerRequest int64, minGasPerRequest int64) *Limits {
	this := Limits{}
	this.MaxGasExternalViewCall = maxGasExternalViewCall
	this.MaxGasPerBlock = maxGasPerBlock
	this.MaxGasPerRequest = maxGasPerRequest
	this.MinGasPerRequest = minGasPerRequest
	return &this
}

// NewLimitsWithDefaults instantiates a new Limits object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewLimitsWithDefaults() *Limits {
	this := Limits{}
	return &this
}

// GetMaxGasExternalViewCall returns the MaxGasExternalViewCall field value
func (o *Limits) GetMaxGasExternalViewCall() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.MaxGasExternalViewCall
}

// GetMaxGasExternalViewCallOk returns a tuple with the MaxGasExternalViewCall field value
// and a boolean to check if the value has been set.
func (o *Limits) GetMaxGasExternalViewCallOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.MaxGasExternalViewCall, true
}

// SetMaxGasExternalViewCall sets field value
func (o *Limits) SetMaxGasExternalViewCall(v int64) {
	o.MaxGasExternalViewCall = v
}

// GetMaxGasPerBlock returns the MaxGasPerBlock field value
func (o *Limits) GetMaxGasPerBlock() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.MaxGasPerBlock
}

// GetMaxGasPerBlockOk returns a tuple with the MaxGasPerBlock field value
// and a boolean to check if the value has been set.
func (o *Limits) GetMaxGasPerBlockOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.MaxGasPerBlock, true
}

// SetMaxGasPerBlock sets field value
func (o *Limits) SetMaxGasPerBlock(v int64) {
	o.MaxGasPerBlock = v
}

// GetMaxGasPerRequest returns the MaxGasPerRequest field value
func (o *Limits) GetMaxGasPerRequest() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.MaxGasPerRequest
}

// GetMaxGasPerRequestOk returns a tuple with the MaxGasPerRequest field value
// and a boolean to check if the value has been set.
func (o *Limits) GetMaxGasPerRequestOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.MaxGasPerRequest, true
}

// SetMaxGasPerRequest sets field value
func (o *Limits) SetMaxGasPerRequest(v int64) {
	o.MaxGasPerRequest = v
}

// GetMinGasPerRequest returns the MinGasPerRequest field value
func (o *Limits) GetMinGasPerRequest() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.MinGasPerRequest
}

// GetMinGasPerRequestOk returns a tuple with the MinGasPerRequest field value
// and a boolean to check if the value has been set.
func (o *Limits) GetMinGasPerRequestOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.MinGasPerRequest, true
}

// SetMinGasPerRequest sets field value
func (o *Limits) SetMinGasPerRequest(v int64) {
	o.MinGasPerRequest = v
}

func (o Limits) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Limits) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["maxGasExternalViewCall"] = o.MaxGasExternalViewCall
	toSerialize["maxGasPerBlock"] = o.MaxGasPerBlock
	toSerialize["maxGasPerRequest"] = o.MaxGasPerRequest
	toSerialize["minGasPerRequest"] = o.MinGasPerRequest
	return toSerialize, nil
}

func (o *Limits) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"maxGasExternalViewCall",
		"maxGasPerBlock",
		"maxGasPerRequest",
		"minGasPerRequest",
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

	varLimits := _Limits{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varLimits)

	if err != nil {
		return err
	}

	*o = Limits(varLimits)

	return err
}

type NullableLimits struct {
	value *Limits
	isSet bool
}

func (v NullableLimits) Get() *Limits {
	return v.value
}

func (v *NullableLimits) Set(val *Limits) {
	v.value = val
	v.isSet = true
}

func (v NullableLimits) IsSet() bool {
	return v.isSet
}

func (v *NullableLimits) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableLimits(val *Limits) *NullableLimits {
	return &NullableLimits{value: val, isSet: true}
}

func (v NullableLimits) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableLimits) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


