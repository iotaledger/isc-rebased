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

// checks if the FeePolicy type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &FeePolicy{}

// FeePolicy struct for FeePolicy
type FeePolicy struct {
	EvmGasRatio Ratio32 `json:"evmGasRatio"`
	GasPerToken Ratio32 `json:"gasPerToken"`
	// The validator fee share.
	ValidatorFeeShare int32 `json:"validatorFeeShare"`
}

type _FeePolicy FeePolicy

// NewFeePolicy instantiates a new FeePolicy object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewFeePolicy(evmGasRatio Ratio32, gasPerToken Ratio32, validatorFeeShare int32) *FeePolicy {
	this := FeePolicy{}
	this.EvmGasRatio = evmGasRatio
	this.GasPerToken = gasPerToken
	this.ValidatorFeeShare = validatorFeeShare
	return &this
}

// NewFeePolicyWithDefaults instantiates a new FeePolicy object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewFeePolicyWithDefaults() *FeePolicy {
	this := FeePolicy{}
	return &this
}

// GetEvmGasRatio returns the EvmGasRatio field value
func (o *FeePolicy) GetEvmGasRatio() Ratio32 {
	if o == nil {
		var ret Ratio32
		return ret
	}

	return o.EvmGasRatio
}

// GetEvmGasRatioOk returns a tuple with the EvmGasRatio field value
// and a boolean to check if the value has been set.
func (o *FeePolicy) GetEvmGasRatioOk() (*Ratio32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.EvmGasRatio, true
}

// SetEvmGasRatio sets field value
func (o *FeePolicy) SetEvmGasRatio(v Ratio32) {
	o.EvmGasRatio = v
}

// GetGasPerToken returns the GasPerToken field value
func (o *FeePolicy) GetGasPerToken() Ratio32 {
	if o == nil {
		var ret Ratio32
		return ret
	}

	return o.GasPerToken
}

// GetGasPerTokenOk returns a tuple with the GasPerToken field value
// and a boolean to check if the value has been set.
func (o *FeePolicy) GetGasPerTokenOk() (*Ratio32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.GasPerToken, true
}

// SetGasPerToken sets field value
func (o *FeePolicy) SetGasPerToken(v Ratio32) {
	o.GasPerToken = v
}

// GetValidatorFeeShare returns the ValidatorFeeShare field value
func (o *FeePolicy) GetValidatorFeeShare() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.ValidatorFeeShare
}

// GetValidatorFeeShareOk returns a tuple with the ValidatorFeeShare field value
// and a boolean to check if the value has been set.
func (o *FeePolicy) GetValidatorFeeShareOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ValidatorFeeShare, true
}

// SetValidatorFeeShare sets field value
func (o *FeePolicy) SetValidatorFeeShare(v int32) {
	o.ValidatorFeeShare = v
}

func (o FeePolicy) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o FeePolicy) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["evmGasRatio"] = o.EvmGasRatio
	toSerialize["gasPerToken"] = o.GasPerToken
	toSerialize["validatorFeeShare"] = o.ValidatorFeeShare
	return toSerialize, nil
}

func (o *FeePolicy) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"evmGasRatio",
		"gasPerToken",
		"validatorFeeShare",
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

	varFeePolicy := _FeePolicy{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varFeePolicy)

	if err != nil {
		return err
	}

	*o = FeePolicy(varFeePolicy)

	return err
}

type NullableFeePolicy struct {
	value *FeePolicy
	isSet bool
}

func (v NullableFeePolicy) Get() *FeePolicy {
	return v.value
}

func (v *NullableFeePolicy) Set(val *FeePolicy) {
	v.value = val
	v.isSet = true
}

func (v NullableFeePolicy) IsSet() bool {
	return v.isSet
}

func (v *NullableFeePolicy) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableFeePolicy(val *FeePolicy) *NullableFeePolicy {
	return &NullableFeePolicy{value: val, isSet: true}
}

func (v NullableFeePolicy) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableFeePolicy) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


