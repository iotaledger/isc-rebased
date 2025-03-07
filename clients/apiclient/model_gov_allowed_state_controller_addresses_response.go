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

// checks if the GovAllowedStateControllerAddressesResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GovAllowedStateControllerAddressesResponse{}

// GovAllowedStateControllerAddressesResponse struct for GovAllowedStateControllerAddressesResponse
type GovAllowedStateControllerAddressesResponse struct {
	// The allowed state controller addresses (Hex Address)
	Addresses []string `json:"addresses,omitempty"`
}

// NewGovAllowedStateControllerAddressesResponse instantiates a new GovAllowedStateControllerAddressesResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGovAllowedStateControllerAddressesResponse() *GovAllowedStateControllerAddressesResponse {
	this := GovAllowedStateControllerAddressesResponse{}
	return &this
}

// NewGovAllowedStateControllerAddressesResponseWithDefaults instantiates a new GovAllowedStateControllerAddressesResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGovAllowedStateControllerAddressesResponseWithDefaults() *GovAllowedStateControllerAddressesResponse {
	this := GovAllowedStateControllerAddressesResponse{}
	return &this
}

// GetAddresses returns the Addresses field value if set, zero value otherwise.
func (o *GovAllowedStateControllerAddressesResponse) GetAddresses() []string {
	if o == nil || IsNil(o.Addresses) {
		var ret []string
		return ret
	}
	return o.Addresses
}

// GetAddressesOk returns a tuple with the Addresses field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GovAllowedStateControllerAddressesResponse) GetAddressesOk() ([]string, bool) {
	if o == nil || IsNil(o.Addresses) {
		return nil, false
	}
	return o.Addresses, true
}

// HasAddresses returns a boolean if a field has been set.
func (o *GovAllowedStateControllerAddressesResponse) HasAddresses() bool {
	if o != nil && !IsNil(o.Addresses) {
		return true
	}

	return false
}

// SetAddresses gets a reference to the given []string and assigns it to the Addresses field.
func (o *GovAllowedStateControllerAddressesResponse) SetAddresses(v []string) {
	o.Addresses = v
}

func (o GovAllowedStateControllerAddressesResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GovAllowedStateControllerAddressesResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Addresses) {
		toSerialize["addresses"] = o.Addresses
	}
	return toSerialize, nil
}

type NullableGovAllowedStateControllerAddressesResponse struct {
	value *GovAllowedStateControllerAddressesResponse
	isSet bool
}

func (v NullableGovAllowedStateControllerAddressesResponse) Get() *GovAllowedStateControllerAddressesResponse {
	return v.value
}

func (v *NullableGovAllowedStateControllerAddressesResponse) Set(val *GovAllowedStateControllerAddressesResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableGovAllowedStateControllerAddressesResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableGovAllowedStateControllerAddressesResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGovAllowedStateControllerAddressesResponse(val *GovAllowedStateControllerAddressesResponse) *NullableGovAllowedStateControllerAddressesResponse {
	return &NullableGovAllowedStateControllerAddressesResponse{value: val, isSet: true}
}

func (v NullableGovAllowedStateControllerAddressesResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGovAllowedStateControllerAddressesResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


