/*
Wasp API

REST API for the Wasp node

API version: 0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package apiclient

import (
	"encoding/json"
	"time"
	"bytes"
	"fmt"
)

// checks if the BlockInfoResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &BlockInfoResponse{}

// BlockInfoResponse struct for BlockInfoResponse
type BlockInfoResponse struct {
	BlockIndex uint32 `json:"blockIndex"`
	// The burned gas (uint64 as string)
	GasBurned string `json:"gasBurned"`
	// The charged gas fee (uint64 as string)
	GasFeeCharged string `json:"gasFeeCharged"`
	NumOffLedgerRequests uint32 `json:"numOffLedgerRequests"`
	NumSuccessfulRequests uint32 `json:"numSuccessfulRequests"`
	Timestamp time.Time `json:"timestamp"`
	TotalRequests uint32 `json:"totalRequests"`
}

type _BlockInfoResponse BlockInfoResponse

// NewBlockInfoResponse instantiates a new BlockInfoResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBlockInfoResponse(blockIndex uint32, gasBurned string, gasFeeCharged string, numOffLedgerRequests uint32, numSuccessfulRequests uint32, timestamp time.Time, totalRequests uint32) *BlockInfoResponse {
	this := BlockInfoResponse{}
	this.BlockIndex = blockIndex
	this.GasBurned = gasBurned
	this.GasFeeCharged = gasFeeCharged
	this.NumOffLedgerRequests = numOffLedgerRequests
	this.NumSuccessfulRequests = numSuccessfulRequests
	this.Timestamp = timestamp
	this.TotalRequests = totalRequests
	return &this
}

// NewBlockInfoResponseWithDefaults instantiates a new BlockInfoResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBlockInfoResponseWithDefaults() *BlockInfoResponse {
	this := BlockInfoResponse{}
	return &this
}

// GetBlockIndex returns the BlockIndex field value
func (o *BlockInfoResponse) GetBlockIndex() uint32 {
	if o == nil {
		var ret uint32
		return ret
	}

	return o.BlockIndex
}

// GetBlockIndexOk returns a tuple with the BlockIndex field value
// and a boolean to check if the value has been set.
func (o *BlockInfoResponse) GetBlockIndexOk() (*uint32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.BlockIndex, true
}

// SetBlockIndex sets field value
func (o *BlockInfoResponse) SetBlockIndex(v uint32) {
	o.BlockIndex = v
}

// GetGasBurned returns the GasBurned field value
func (o *BlockInfoResponse) GetGasBurned() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.GasBurned
}

// GetGasBurnedOk returns a tuple with the GasBurned field value
// and a boolean to check if the value has been set.
func (o *BlockInfoResponse) GetGasBurnedOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.GasBurned, true
}

// SetGasBurned sets field value
func (o *BlockInfoResponse) SetGasBurned(v string) {
	o.GasBurned = v
}

// GetGasFeeCharged returns the GasFeeCharged field value
func (o *BlockInfoResponse) GetGasFeeCharged() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.GasFeeCharged
}

// GetGasFeeChargedOk returns a tuple with the GasFeeCharged field value
// and a boolean to check if the value has been set.
func (o *BlockInfoResponse) GetGasFeeChargedOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.GasFeeCharged, true
}

// SetGasFeeCharged sets field value
func (o *BlockInfoResponse) SetGasFeeCharged(v string) {
	o.GasFeeCharged = v
}

// GetNumOffLedgerRequests returns the NumOffLedgerRequests field value
func (o *BlockInfoResponse) GetNumOffLedgerRequests() uint32 {
	if o == nil {
		var ret uint32
		return ret
	}

	return o.NumOffLedgerRequests
}

// GetNumOffLedgerRequestsOk returns a tuple with the NumOffLedgerRequests field value
// and a boolean to check if the value has been set.
func (o *BlockInfoResponse) GetNumOffLedgerRequestsOk() (*uint32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.NumOffLedgerRequests, true
}

// SetNumOffLedgerRequests sets field value
func (o *BlockInfoResponse) SetNumOffLedgerRequests(v uint32) {
	o.NumOffLedgerRequests = v
}

// GetNumSuccessfulRequests returns the NumSuccessfulRequests field value
func (o *BlockInfoResponse) GetNumSuccessfulRequests() uint32 {
	if o == nil {
		var ret uint32
		return ret
	}

	return o.NumSuccessfulRequests
}

// GetNumSuccessfulRequestsOk returns a tuple with the NumSuccessfulRequests field value
// and a boolean to check if the value has been set.
func (o *BlockInfoResponse) GetNumSuccessfulRequestsOk() (*uint32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.NumSuccessfulRequests, true
}

// SetNumSuccessfulRequests sets field value
func (o *BlockInfoResponse) SetNumSuccessfulRequests(v uint32) {
	o.NumSuccessfulRequests = v
}

// GetTimestamp returns the Timestamp field value
func (o *BlockInfoResponse) GetTimestamp() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.Timestamp
}

// GetTimestampOk returns a tuple with the Timestamp field value
// and a boolean to check if the value has been set.
func (o *BlockInfoResponse) GetTimestampOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Timestamp, true
}

// SetTimestamp sets field value
func (o *BlockInfoResponse) SetTimestamp(v time.Time) {
	o.Timestamp = v
}

// GetTotalRequests returns the TotalRequests field value
func (o *BlockInfoResponse) GetTotalRequests() uint32 {
	if o == nil {
		var ret uint32
		return ret
	}

	return o.TotalRequests
}

// GetTotalRequestsOk returns a tuple with the TotalRequests field value
// and a boolean to check if the value has been set.
func (o *BlockInfoResponse) GetTotalRequestsOk() (*uint32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TotalRequests, true
}

// SetTotalRequests sets field value
func (o *BlockInfoResponse) SetTotalRequests(v uint32) {
	o.TotalRequests = v
}

func (o BlockInfoResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o BlockInfoResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["blockIndex"] = o.BlockIndex
	toSerialize["gasBurned"] = o.GasBurned
	toSerialize["gasFeeCharged"] = o.GasFeeCharged
	toSerialize["numOffLedgerRequests"] = o.NumOffLedgerRequests
	toSerialize["numSuccessfulRequests"] = o.NumSuccessfulRequests
	toSerialize["timestamp"] = o.Timestamp
	toSerialize["totalRequests"] = o.TotalRequests
	return toSerialize, nil
}

func (o *BlockInfoResponse) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"blockIndex",
		"gasBurned",
		"gasFeeCharged",
		"numOffLedgerRequests",
		"numSuccessfulRequests",
		"timestamp",
		"totalRequests",
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

	varBlockInfoResponse := _BlockInfoResponse{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varBlockInfoResponse)

	if err != nil {
		return err
	}

	*o = BlockInfoResponse(varBlockInfoResponse)

	return err
}

type NullableBlockInfoResponse struct {
	value *BlockInfoResponse
	isSet bool
}

func (v NullableBlockInfoResponse) Get() *BlockInfoResponse {
	return v.value
}

func (v *NullableBlockInfoResponse) Set(val *BlockInfoResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableBlockInfoResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableBlockInfoResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBlockInfoResponse(val *BlockInfoResponse) *NullableBlockInfoResponse {
	return &NullableBlockInfoResponse{value: val, isSet: true}
}

func (v NullableBlockInfoResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBlockInfoResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


