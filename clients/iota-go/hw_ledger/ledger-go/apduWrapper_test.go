/*******************************************************************************
*   (c) Zondax AG
*
*  Licensed under the Apache License, Version 2.0 (the "License");
*  you may not use this file except in compliance with the License.
*  You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
*  Unless required by applicable law or agreed to in writing, software
*  distributed under the License is distributed on an "AS IS" BASIS,
*  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*  See the License for the specific language governing permissions and
*  limitations under the License.
********************************************************************************/

package ledger_go

import (
	"bytes"
	"fmt"
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func Test_SerializePacket_EmptyCommand(t *testing.T) {
	command := make([]byte, 1)

	_, _, err := SerializePacket(0x0101, command, 64, 0)
	assert.Nil(t, err, "Commands smaller than 3 bytes should return error")
}

func Test_SerializePacket_PacketSize(t *testing.T) {
	packetSize := 64
	type header struct {
		channel     uint16
		tag         uint8
		sequenceIdx uint16
		commandLen  uint16
	}

	h := header{channel: 0x0101, tag: 0x05, sequenceIdx: 0, commandLen: 32}

	command := make([]byte, h.commandLen)

	result, _, _ := SerializePacket(
		h.channel,
		command,
		packetSize,
		h.sequenceIdx)

	assert.Equal(t, len(result), packetSize, "Packet size is wrong")
}

func Test_SerializePacket_Header(t *testing.T) {
	packetSize := 64
	type header struct {
		channel     uint16
		tag         uint8
		sequenceIdx uint16
		commandLen  uint16
	}

	h := header{channel: 0x0101, tag: 0x05, sequenceIdx: 0, commandLen: 32}

	command := make([]byte, h.commandLen)

	result, _, _ := SerializePacket(
		h.channel,
		command,
		packetSize,
		h.sequenceIdx)

	assert.Equal(t, codec.Uint16(result), h.channel, "Channel not properly serialized")
	assert.Equal(t, result[2], h.tag, "Tag not properly serialized")
	assert.Equal(t, codec.Uint16(result[3:]), h.sequenceIdx, "SequenceIdx not properly serialized")
	assert.Equal(t, codec.Uint16(result[5:]), h.commandLen, "Command len not properly serialized")
}

func Test_SerializePacket_Offset(t *testing.T) {
	packetSize := 64
	type header struct {
		channel     uint16
		tag         uint8
		sequenceIdx uint16
		commandLen  uint16
	}

	h := header{channel: 0x0101, tag: 0x05, sequenceIdx: 0, commandLen: 100}

	command := make([]byte, h.commandLen)

	_, offset, _ := SerializePacket(
		h.channel,
		command,
		packetSize,
		h.sequenceIdx)

	assert.Equal(t, packetSize-int(unsafe.Sizeof(h))+1, offset, "Wrong offset returned. Offset must point to the next command byte that needs to be packetized.")
}

func Test_WrapCommandAPDU_NumberOfPackets(t *testing.T) {
	packetSize := 64
	type firstHeader struct {
		channel     uint16
		sequenceIdx uint16
		commandLen  uint16
		tag         uint8
	}

	h1 := firstHeader{channel: 0x0101, tag: 0x05, sequenceIdx: 0, commandLen: 100}

	command := make([]byte, h1.commandLen)

	result, _ := WrapCommandAPDU(
		h1.channel,
		command,
		packetSize)

	assert.Equal(t, packetSize*2, len(result), "Result buffer size is not correct")
}

func Test_WrapCommandAPDU_CheckHeaders(t *testing.T) {
	packetSize := 64
	type firstHeader struct {
		channel     uint16
		sequenceIdx uint16
		commandLen  uint16
		tag         uint8
	}

	h1 := firstHeader{channel: 0x0101, tag: 0x05, sequenceIdx: 0, commandLen: 100}

	command := make([]byte, h1.commandLen)

	result, _ := WrapCommandAPDU(
		h1.channel,
		command,
		packetSize)

	assert.Equal(t, h1.channel, codec.Uint16(result), "Channel not properly serialized")
	assert.Equal(t, h1.tag, result[2], "Tag not properly serialized")
	assert.Equal(t, 0, int(codec.Uint16(result[3:])), "SequenceIdx not properly serialized")
	assert.Equal(t, int(h1.commandLen), int(codec.Uint16(result[5:])), "Command len not properly serialized")

	offsetOfSecondPacket := packetSize
	assert.Equal(t, h1.channel, codec.Uint16(result[offsetOfSecondPacket:]), "Channel not properly serialized")
	assert.Equal(t, h1.tag, result[offsetOfSecondPacket+2], "Tag not properly serialized")
	assert.Equal(t, 1, int(codec.Uint16(result[offsetOfSecondPacket+3:])), "SequenceIdx not properly serialized")
}

func Test_WrapCommandAPDU_CheckData(t *testing.T) {
	packetSize := 64
	type firstHeader struct {
		channel     uint16
		sequenceIdx uint16
		commandLen  uint16
		tag         uint8
	}

	h1 := firstHeader{channel: 0x0101, tag: 0x05, sequenceIdx: 0, commandLen: 200}

	command := make([]byte, h1.commandLen)

	for i := range command {
		command[i] = byte(i % 256)
	}

	result, _ := WrapCommandAPDU(
		h1.channel,
		command,
		packetSize)

	// Check data in the first packet
	assert.True(t, bytes.Equal(command[0:64-7], result[7:64]))

	result = result[64:]
	command = command[64-7:]
	// Check data in the second packet
	assert.True(t, bytes.Equal(command[0:64-5], result[5:64]))

	result = result[64:]
	command = command[64-5:]
	// Check data in the third packet
	assert.True(t, bytes.Equal(command[0:64-5], result[5:64]))

	result = result[64:]
	command = command[64-5:]

	// Check data in the last packet
	assert.True(t, bytes.Equal(command[0:], result[5:5+len(command)]))

	// The remaining bytes in the result should be zeros
	result = result[5+len(command):]
	assert.True(t, bytes.Equal(result, make([]byte, len(result))))
}

func Test_DeserializePacket_FirstPacket(t *testing.T) {
	sampleCommand := []byte{'H', 'e', 'l', 'l', 'o', 0}

	packetSize := 64
	firstPacketHeaderSize := 7
	packet, _, _ := SerializePacket(0x0101, sampleCommand, packetSize, 0)

	output, totalSize, isSequenceZero, err := DeserializePacket(0x0101, packet, 0)

	assert.Nil(t, err, "Simple deserialize should not have errors")
	assert.Equal(t, len(sampleCommand), int(totalSize), "TotalSize is incorrect")
	assert.Equal(t, packetSize-firstPacketHeaderSize, len(output), "Size of the deserialized packet is wrong")
	assert.Equal(t, true, isSequenceZero, "Test Case Should Find Sequence == 0")
	assert.True(t, bytes.Equal(output[:len(sampleCommand)], sampleCommand), "Deserialized message does not match the original")
}

func Test_DeserializePacket_SecondMessage(t *testing.T) {
	sampleCommand := []byte{'H', 'e', 'l', 'l', 'o', 0}

	packetSize := 64
	firstPacketHeaderSize := 5 // second packet does not have responseLength (uint16) in the header
	packet, _, _ := SerializePacket(0x0101, sampleCommand, packetSize, 1)

	output, totalSize, isSequenceZero, err := DeserializePacket(0x0101, packet, 1)

	assert.Nil(t, err, "Simple deserialize should not have errors")
	assert.Equal(t, 0, int(totalSize), "TotalSize should not be returned from deserialization of non-first packet")
	assert.Equal(t, packetSize-firstPacketHeaderSize, len(output), "Size of the deserialized packet is wrong")
	assert.Equal(t, false, isSequenceZero, "Test Case Should Find Sequence == 1")
	assert.True(t, bytes.Equal(output[:len(sampleCommand)], sampleCommand), "Deserialized message does not match the original")
}

func Test_UnwrapApdu_SmokeTest(t *testing.T) {
	const channel uint16 = 0x8002

	inputSize := 200
	packetSize := 64

	// Initialize some dummy input
	input := make([]byte, inputSize)
	for i := range input {
		input[i] = byte(i % 256)
	}

	serialized, _ := WrapCommandAPDU(channel, input, packetSize)

	// Allocate enough packets
	pipe := make(chan []byte, int(math.Ceil(float64(inputSize)/float64(packetSize))))
	// Send all the packets to the pipe
	for len(serialized) > 0 {
		pipe <- serialized[:packetSize]
		serialized = serialized[packetSize:]
	}

	output, _ := UnwrapResponseAPDU(channel, pipe, packetSize)

	fmt.Printf("INPUT     : %x\n", input)
	fmt.Printf("SERIALIZED: %x\n", serialized)
	fmt.Printf("OUTPUT    : %x\n", output)

	assert.Equal(t, len(input), len(output), "Input and output messages have different size")
	assert.True(t,
		bytes.Equal(input, output),
		"Input message does not match message which was serialized and then deserialized")
}

func TestSerializePacketWithInvalidSize(t *testing.T) {
	_, _, err := SerializePacket(0x0101, []byte{1, 2}, 2, 0)
	assert.ErrorIs(t, err, ErrPacketSize)
}

func TestDeserializePacketWithInvalidChannel(t *testing.T) {
	packet := []byte{0x02, 0x02, 0x05, 0x00, 0x00, 0x00, 0x20}
	_, _, _, err := DeserializePacket(0x0101, packet, 0)
	assert.ErrorIs(t, err, ErrInvalidChannel)
}

func TestSerializeDeserialize(t *testing.T) {
	sampleCommand := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	channel := uint16(0x0101)
	packetSize := 64
	sequenceIdx := uint16(0)

	packet, _, err := SerializePacket(channel, sampleCommand, packetSize, sequenceIdx)
	assert.NoError(t, err)

	output, _, _, err := DeserializePacket(channel, packet, sequenceIdx)
	assert.NoError(t, err)

	assert.True(t, bytes.Equal(output[:len(sampleCommand)], sampleCommand), "Deserialized message does not match the original")
}
