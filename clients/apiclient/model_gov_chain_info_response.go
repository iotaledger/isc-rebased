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

// checks if the GovChainInfoResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GovChainInfoResponse{}

// GovChainInfoResponse struct for GovChainInfoResponse
type GovChainInfoResponse struct {
	// ChainID (Bech32-encoded).
	ChainID string `json:"chainID"`
	// The chain owner address (Bech32-encoded).
	ChainOwnerId string `json:"chainOwnerId"`
	GasFeePolicy FeePolicy `json:"gasFeePolicy"`
	GasLimits Limits `json:"gasLimits"`
	// The EVM json rpc url
	MetadataEvmJsonRpcUrl string `json:"metadataEvmJsonRpcUrl"`
	// The EVM websocket url
	MetadataEvmWebSocketUrl string `json:"metadataEvmWebSocketUrl"`
	// The fully qualified public url leading to the chains metadata
	PublicUrl string `json:"publicUrl"`
}

// NewGovChainInfoResponse instantiates a new GovChainInfoResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGovChainInfoResponse(chainID string, chainOwnerId string, gasFeePolicy FeePolicy, gasLimits Limits, metadataEvmJsonRpcUrl string, metadataEvmWebSocketUrl string, publicUrl string) *GovChainInfoResponse {
	this := GovChainInfoResponse{}
	this.ChainID = chainID
	this.ChainOwnerId = chainOwnerId
	this.GasFeePolicy = gasFeePolicy
	this.GasLimits = gasLimits
	this.MetadataEvmJsonRpcUrl = metadataEvmJsonRpcUrl
	this.MetadataEvmWebSocketUrl = metadataEvmWebSocketUrl
	this.PublicUrl = publicUrl
	return &this
}

// NewGovChainInfoResponseWithDefaults instantiates a new GovChainInfoResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGovChainInfoResponseWithDefaults() *GovChainInfoResponse {
	this := GovChainInfoResponse{}
	return &this
}

// GetChainID returns the ChainID field value
func (o *GovChainInfoResponse) GetChainID() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ChainID
}

// GetChainIDOk returns a tuple with the ChainID field value
// and a boolean to check if the value has been set.
func (o *GovChainInfoResponse) GetChainIDOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ChainID, true
}

// SetChainID sets field value
func (o *GovChainInfoResponse) SetChainID(v string) {
	o.ChainID = v
}

// GetChainOwnerId returns the ChainOwnerId field value
func (o *GovChainInfoResponse) GetChainOwnerId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ChainOwnerId
}

// GetChainOwnerIdOk returns a tuple with the ChainOwnerId field value
// and a boolean to check if the value has been set.
func (o *GovChainInfoResponse) GetChainOwnerIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ChainOwnerId, true
}

// SetChainOwnerId sets field value
func (o *GovChainInfoResponse) SetChainOwnerId(v string) {
	o.ChainOwnerId = v
}

// GetGasFeePolicy returns the GasFeePolicy field value
func (o *GovChainInfoResponse) GetGasFeePolicy() FeePolicy {
	if o == nil {
		var ret FeePolicy
		return ret
	}

	return o.GasFeePolicy
}

// GetGasFeePolicyOk returns a tuple with the GasFeePolicy field value
// and a boolean to check if the value has been set.
func (o *GovChainInfoResponse) GetGasFeePolicyOk() (*FeePolicy, bool) {
	if o == nil {
		return nil, false
	}
	return &o.GasFeePolicy, true
}

// SetGasFeePolicy sets field value
func (o *GovChainInfoResponse) SetGasFeePolicy(v FeePolicy) {
	o.GasFeePolicy = v
}

// GetGasLimits returns the GasLimits field value
func (o *GovChainInfoResponse) GetGasLimits() Limits {
	if o == nil {
		var ret Limits
		return ret
	}

	return o.GasLimits
}

// GetGasLimitsOk returns a tuple with the GasLimits field value
// and a boolean to check if the value has been set.
func (o *GovChainInfoResponse) GetGasLimitsOk() (*Limits, bool) {
	if o == nil {
		return nil, false
	}
	return &o.GasLimits, true
}

// SetGasLimits sets field value
func (o *GovChainInfoResponse) SetGasLimits(v Limits) {
	o.GasLimits = v
}

// GetMetadataEvmJsonRpcUrl returns the MetadataEvmJsonRpcUrl field value
func (o *GovChainInfoResponse) GetMetadataEvmJsonRpcUrl() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.MetadataEvmJsonRpcUrl
}

// GetMetadataEvmJsonRpcUrlOk returns a tuple with the MetadataEvmJsonRpcUrl field value
// and a boolean to check if the value has been set.
func (o *GovChainInfoResponse) GetMetadataEvmJsonRpcUrlOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.MetadataEvmJsonRpcUrl, true
}

// SetMetadataEvmJsonRpcUrl sets field value
func (o *GovChainInfoResponse) SetMetadataEvmJsonRpcUrl(v string) {
	o.MetadataEvmJsonRpcUrl = v
}

// GetMetadataEvmWebSocketUrl returns the MetadataEvmWebSocketUrl field value
func (o *GovChainInfoResponse) GetMetadataEvmWebSocketUrl() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.MetadataEvmWebSocketUrl
}

// GetMetadataEvmWebSocketUrlOk returns a tuple with the MetadataEvmWebSocketUrl field value
// and a boolean to check if the value has been set.
func (o *GovChainInfoResponse) GetMetadataEvmWebSocketUrlOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.MetadataEvmWebSocketUrl, true
}

// SetMetadataEvmWebSocketUrl sets field value
func (o *GovChainInfoResponse) SetMetadataEvmWebSocketUrl(v string) {
	o.MetadataEvmWebSocketUrl = v
}

// GetPublicUrl returns the PublicUrl field value
func (o *GovChainInfoResponse) GetPublicUrl() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.PublicUrl
}

// GetPublicUrlOk returns a tuple with the PublicUrl field value
// and a boolean to check if the value has been set.
func (o *GovChainInfoResponse) GetPublicUrlOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PublicUrl, true
}

// SetPublicUrl sets field value
func (o *GovChainInfoResponse) SetPublicUrl(v string) {
	o.PublicUrl = v
}

func (o GovChainInfoResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GovChainInfoResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["chainID"] = o.ChainID
	toSerialize["chainOwnerId"] = o.ChainOwnerId
	toSerialize["gasFeePolicy"] = o.GasFeePolicy
	toSerialize["gasLimits"] = o.GasLimits
	toSerialize["metadataEvmJsonRpcUrl"] = o.MetadataEvmJsonRpcUrl
	toSerialize["metadataEvmWebSocketUrl"] = o.MetadataEvmWebSocketUrl
	toSerialize["publicUrl"] = o.PublicUrl
	return toSerialize, nil
}

type NullableGovChainInfoResponse struct {
	value *GovChainInfoResponse
	isSet bool
}

func (v NullableGovChainInfoResponse) Get() *GovChainInfoResponse {
	return v.value
}

func (v *NullableGovChainInfoResponse) Set(val *GovChainInfoResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableGovChainInfoResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableGovChainInfoResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGovChainInfoResponse(val *GovChainInfoResponse) *NullableGovChainInfoResponse {
	return &NullableGovChainInfoResponse{value: val, isSet: true}
}

func (v NullableGovChainInfoResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGovChainInfoResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


