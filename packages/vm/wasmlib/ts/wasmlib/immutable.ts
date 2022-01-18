// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// immutable proxies to host objects

import { base58Encode } from "./context";
import {Convert} from "./convert";
import {ScAddress,ScAgentID,ScChainID,ScColor,ScHash,ScHname,ScRequestID} from "./hashtypes";
import * as host from "./host";
import {Key32,MapKey} from "./keys";

// value proxy for immutable ScAddress in host container
export class ScImmutableAddress {
    objID: i32;
    keyID: Key32;

    constructor(objID: i32, keyID: Key32) {
        this.objID = objID;
        this.keyID = keyID;
    }

    // check if value exists in host container
    exists(): boolean {
        return host.exists(this.objID, this.keyID, host.TYPE_ADDRESS);
    }

    // human-readable string representation
    toString(): string {
        return this.value().toString();
    }

    // get value from host container
    value(): ScAddress {
        return ScAddress.fromBytes(host.getBytes(this.objID, this.keyID, host.TYPE_ADDRESS));
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// array proxy for immutable array of ScAddress
export class ScImmutableAddressArray {
    objID: i32;

    constructor(id: i32) {
        this.objID = id;
    }

    // get value proxy for item at index, index can be 0..length()-1
    getAddress(index: i32): ScImmutableAddress {
        return new ScImmutableAddress(this.objID, new Key32(index));
    }

    // number of items in array
    length(): i32 {
        return host.getLength(this.objID);
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// value proxy for immutable ScAgentID in host container
export class ScImmutableAgentID {
    objID: i32;
    keyID: Key32;

    constructor(objID: i32, keyID: Key32) {
        this.objID = objID;
        this.keyID = keyID;
    }

    // check if value exists in host container
    exists(): boolean {
        return host.exists(this.objID, this.keyID, host.TYPE_AGENT_ID);
    }

    // human-readable string representation
    toString(): string {
        return this.value().toString();
    }

    // get value from host container
    value(): ScAgentID {
        return ScAgentID.fromBytes(host.getBytes(this.objID, this.keyID, host.TYPE_AGENT_ID));
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// array proxy for immutable array of ScAgentID
export class ScImmutableAgentIDArray {
    objID: i32;

    constructor(id: i32) {
        this.objID = id;
    }

    // get value proxy for item at index, index can be 0..length()-1
    getAgentID(index: i32): ScImmutableAgentID {
        return new ScImmutableAgentID(this.objID, new Key32(index));
    }

    // number of items in array
    length(): i32 {
        return host.getLength(this.objID);
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// value proxy for immutable Bool in host container
export class ScImmutableBool {
    objID: i32;
    keyID: Key32;

    constructor(objID: i32, keyID: Key32) {
        this.objID = objID;
        this.keyID = keyID;
    }

    // check if value exists in host container
    exists(): boolean {
        return host.exists(this.objID, this.keyID, host.TYPE_BOOL);
    }

    // human-readable string representation
    toString(): string {
        return this.value().toString();
    }

    // get value from host container
    value(): boolean {
        let bytes = host.getBytes(this.objID, this.keyID, host.TYPE_BOOL);
        return bytes[0] != 0;
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// array proxy for immutable array of Bool
export class ScImmutableBoolArray {
    objID: i32;

    constructor(id: i32) {
        this.objID = id;
    }

    // get value proxy for item at index, index can be 0..length()-1
    getBool(index: i32): ScImmutableBool {
        return new ScImmutableBool(this.objID, new Key32(index));
    }

    // number of items in array
    length(): i32 {
        return host.getLength(this.objID);
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// value proxy for immutable bytes array in host container
export class ScImmutableBytes {
    objID: i32;
    keyID: Key32;

    constructor(objID: i32, keyID: Key32) {
        this.objID = objID;
        this.keyID = keyID;
    }

    // check if value exists in host container
    exists(): boolean {
        return host.exists(this.objID, this.keyID, host.TYPE_BYTES);
    }

    // human-readable string representation
    toString(): string {
        return base58Encode(this.value());
    }

    // get value from host container
    value(): u8[] {
        return host.getBytes(this.objID, this.keyID, host.TYPE_BYTES);
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// array proxy for immutable array of byte array
export class ScImmutableBytesArray {
    objID: i32;

    constructor(id: i32) {
        this.objID = id;
    }

    // get value proxy for item at index, index can be 0..length()-1
    getBytes(index: i32): ScImmutableBytes {
        return new ScImmutableBytes(this.objID, new Key32(index));
    }

    // number of items in array
    length(): i32 {
        return host.getLength(this.objID);
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// value proxy for immutable ScChainID in host container
export class ScImmutableChainID {
    objID: i32;
    keyID: Key32;

    constructor(objID: i32, keyID: Key32) {
        this.objID = objID;
        this.keyID = keyID;
    }

    // check if value exists in host container
    exists(): boolean {
        return host.exists(this.objID, this.keyID, host.TYPE_CHAIN_ID);
    }

    // human-readable string representation
    toString(): string {
        return this.value().toString();
    }

    // get value from host container
    value(): ScChainID {
        return ScChainID.fromBytes(host.getBytes(this.objID, this.keyID, host.TYPE_CHAIN_ID));
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// array proxy for immutable array of ScChainID
export class ScImmutableChainIDArray {
    objID: i32;

    constructor(id: i32) {
        this.objID = id;
    }

    // get value proxy for item at index, index can be 0..length()-1
    getChainID(index: i32): ScImmutableChainID {
        return new ScImmutableChainID(this.objID, new Key32(index));
    }

    // number of items in array
    length(): i32 {
        return host.getLength(this.objID);
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// value proxy for immutable ScColor in host container
export class ScImmutableColor {
    objID: i32;
    keyID: Key32;

    constructor(objID: i32, keyID: Key32) {
        this.objID = objID;
        this.keyID = keyID;
    }

    // check if value exists in host container
    exists(): boolean {
        return host.exists(this.objID, this.keyID, host.TYPE_COLOR);
    }

    // human-readable string representation
    toString(): string {
        return this.value().toString();
    }

    // get value from host container
    value(): ScColor {
        return ScColor.fromBytes(host.getBytes(this.objID, this.keyID, host.TYPE_COLOR));
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// array proxy for immutable array of ScColor
export class ScImmutableColorArray {
    objID: i32;

    constructor(id: i32) {
        this.objID = id;
    }

    // get value proxy for item at index, index can be 0..length()-1
    getColor(index: i32): ScImmutableColor {
        return new ScImmutableColor(this.objID, new Key32(index));
    }

    // number of items in array
    length(): i32 {
        return host.getLength(this.objID);
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// value proxy for immutable ScHash in host container
export class ScImmutableHash {
    objID: i32;
    keyID: Key32;

    constructor(objID: i32, keyID: Key32) {
        this.objID = objID;
        this.keyID = keyID;
    }

    // check if value exists in host container
    exists(): boolean {
        return host.exists(this.objID, this.keyID, host.TYPE_HASH);
    }

    // human-readable string representation
    toString(): string {
        return this.value().toString();
    }

    // get value from host container
    value(): ScHash {
        return ScHash.fromBytes(host.getBytes(this.objID, this.keyID, host.TYPE_HASH));
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// array proxy for immutable array of ScHash
export class ScImmutableHashArray {
    objID: i32;

    constructor(id: i32) {
        this.objID = id;
    }

    // get value proxy for item at index, index can be 0..length()-1
    getHash(index: i32): ScImmutableHash {
        return new ScImmutableHash(this.objID, new Key32(index));
    }

    // number of items in array
    length(): i32 {
        return host.getLength(this.objID);
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// value proxy for immutable ScHname in host container
export class ScImmutableHname {
    objID: i32;
    keyID: Key32;

    constructor(objID: i32, keyID: Key32) {
        this.objID = objID;
        this.keyID = keyID;
    }

    // check if value exists in host container
    exists(): boolean {
        return host.exists(this.objID, this.keyID, host.TYPE_HNAME);
    }

    // human-readable string representation
    toString(): string {
        return this.value().toString();
    }

    // get value from host container
    value(): ScHname {
        return ScHname.fromBytes(host.getBytes(this.objID, this.keyID, host.TYPE_HNAME));
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// array proxy for immutable array of ScHname
export class ScImmutableHnameArray {
    objID: i32;

    constructor(id: i32) {
        this.objID = id;
    }

    // get value proxy for item at index, index can be 0..length()-1
    getHname(index: i32): ScImmutableHname {
        return new ScImmutableHname(this.objID, new Key32(index));
    }

    // number of items in array
    length(): i32 {
        return host.getLength(this.objID);
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// value proxy for immutable Int8 in host container
export class ScImmutableInt8 {
    objID: i32;
    keyID: Key32;

    constructor(objID: i32, keyID: Key32) {
        this.objID = objID;
        this.keyID = keyID;
    }

    // check if value exists in host container
    exists(): boolean {
        return host.exists(this.objID, this.keyID, host.TYPE_INT8);
    }

    // human-readable string representation
    toString(): string {
        return this.value().toString();
    }

    // get value from host container
    value(): i8 {
        let bytes = host.getBytes(this.objID, this.keyID, host.TYPE_INT8);
        return bytes[0] as i8;
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// array proxy for immutable array of Int8
export class ScImmutableInt8Array {
    objID: i32;

    constructor(id: i32) {
        this.objID = id;
    }

    // get value proxy for item at index, index can be 0..length()-1
    getInt8(index: i32): ScImmutableInt8 {
        return new ScImmutableInt8(this.objID, new Key32(index));
    }

    // number of items in array
    length(): i32 {
        return host.getLength(this.objID);
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// value proxy for immutable Int16 in host container
export class ScImmutableInt16 {
    objID: i32;
    keyID: Key32;

    constructor(objID: i32, keyID: Key32) {
        this.objID = objID;
        this.keyID = keyID;
    }

    // check if value exists in host container
    exists(): boolean {
        return host.exists(this.objID, this.keyID, host.TYPE_INT16);
    }

    // human-readable string representation
    toString(): string {
        return this.value().toString();
    }

    // get value from host container
    value(): i16 {
        let bytes = host.getBytes(this.objID, this.keyID, host.TYPE_INT16);
        return Convert.toI16(bytes);
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// array proxy for immutable array of Int16
export class ScImmutableInt16Array {
    objID: i32;

    constructor(id: i32) {
        this.objID = id;
    }

    // get value proxy for item at index, index can be 0..length()-1
    getInt16(index: i32): ScImmutableInt16 {
        return new ScImmutableInt16(this.objID, new Key32(index));
    }

    // number of items in array
    length(): i32 {
        return host.getLength(this.objID);
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// value proxy for immutable Int32 in host container
export class ScImmutableInt32 {
    objID: i32;
    keyID: Key32;

    constructor(objID: i32, keyID: Key32) {
        this.objID = objID;
        this.keyID = keyID;
    }

    // check if value exists in host container
    exists(): boolean {
        return host.exists(this.objID, this.keyID, host.TYPE_INT32);
    }

    // human-readable string representation
    toString(): string {
        return this.value().toString();
    }

    // get value from host container
    value(): i32 {
        let bytes = host.getBytes(this.objID, this.keyID, host.TYPE_INT32);
        return Convert.toI32(bytes);
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// array proxy for immutable array of Int32
export class ScImmutableInt32Array {
    objID: i32;

    constructor(id: i32) {
        this.objID = id;
    }

    // get value proxy for item at index, index can be 0..length()-1
    getInt32(index: i32): ScImmutableInt32 {
        return new ScImmutableInt32(this.objID, new Key32(index));
    }

    // number of items in array
    length(): i32 {
        return host.getLength(this.objID);
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// value proxy for immutable Int64 in host container
export class ScImmutableInt64 {
    objID: i32;
    keyID: Key32;

    constructor(objID: i32, keyID: Key32) {
        this.objID = objID;
        this.keyID = keyID;
    }

    // check if value exists in host container
    exists(): boolean {
        return host.exists(this.objID, this.keyID, host.TYPE_INT64);
    }

    // human-readable string representation
    toString(): string {
        return this.value().toString();
    }

    // get value from host container
    value(): i64 {
        let bytes = host.getBytes(this.objID, this.keyID, host.TYPE_INT64);
        return Convert.toI64(bytes);
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// array proxy for immutable array of Int64
export class ScImmutableInt64Array {
    objID: i32;

    constructor(id: i32) {
        this.objID = id;
    }

    // get value proxy for item at index, index can be 0..length()-1
    getInt64(index: i32): ScImmutableInt64 {
        return new ScImmutableInt64(this.objID, new Key32(index));
    }

    // number of items in array
    length(): i32 {
        return host.getLength(this.objID);
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// map proxy for immutable map
export class ScImmutableMap {
    objID: i32;

    constructor(id: i32) {
        this.objID = id;
    }

    callFunc(keyID: Key32, params: u8[]): u8[] {
        return host.callFunc(this.objID, keyID, params);
    }

    // get value proxy for immutable ScAddress field specified by key
    getAddress(key: MapKey): ScImmutableAddress {
        return new ScImmutableAddress(this.objID, key.getKeyID());
    }

    // get array proxy for ScImmutableAddressArray specified by key
    getAddressArray(key: MapKey): ScImmutableAddressArray {
        let arrID = host.getObjectID(this.objID, key.getKeyID(), host.TYPE_ADDRESS | host.TYPE_ARRAY);
        return new ScImmutableAddressArray(arrID);
    }

    // get value proxy for immutable ScAgentID field specified by key
    getAgentID(key: MapKey): ScImmutableAgentID {
        return new ScImmutableAgentID(this.objID, key.getKeyID());
    }

    // get array proxy for ScImmutableAgentIDArray specified by key
    getAgentIDArray(key: MapKey): ScImmutableAgentIDArray {
        let arrID = host.getObjectID(this.objID, key.getKeyID(), host.TYPE_AGENT_ID | host.TYPE_ARRAY);
        return new ScImmutableAgentIDArray(arrID);
    }

    // get value proxy for immutable Bool field specified by key
    getBool(key: MapKey): ScImmutableBool {
        return new ScImmutableBool(this.objID, key.getKeyID());
    }

    // get array proxy for ScImmutableBoolArray specified by key
    getBoolArray(key: MapKey): ScImmutableBoolArray {
        let arrID = host.getObjectID(this.objID, key.getKeyID(), host.TYPE_BOOL | host.TYPE_ARRAY);
        return new ScImmutableBoolArray(arrID);
    }

    // get value proxy for immutable bytes array field specified by key
    getBytes(key: MapKey): ScImmutableBytes {
        return new ScImmutableBytes(this.objID, key.getKeyID());
    }

    // get array proxy for ScImmutableBytesArray specified by key
    getBytesArray(key: MapKey): ScImmutableBytesArray {
        let arrID = host.getObjectID(this.objID, key.getKeyID(), host.TYPE_BYTES | host.TYPE_ARRAY);
        return new ScImmutableBytesArray(arrID);
    }

    // get value proxy for immutable ScChainID field specified by key
    getChainID(key: MapKey): ScImmutableChainID {
        return new ScImmutableChainID(this.objID, key.getKeyID());
    }

    // get array proxy for ScImmutableChainIDArray specified by key
    getChainIDArray(key: MapKey): ScImmutableChainIDArray {
        let arrID = host.getObjectID(this.objID, key.getKeyID(), host.TYPE_CHAIN_ID | host.TYPE_ARRAY);
        return new ScImmutableChainIDArray(arrID);
    }

    // get value proxy for immutable ScColor field specified by key
    getColor(key: MapKey): ScImmutableColor {
        return new ScImmutableColor(this.objID, key.getKeyID());
    }

    // get array proxy for ScImmutableColorArray specified by key
    getColorArray(key: MapKey): ScImmutableColorArray {
        let arrID = host.getObjectID(this.objID, key.getKeyID(), host.TYPE_COLOR | host.TYPE_ARRAY);
        return new ScImmutableColorArray(arrID);
    }

    // get value proxy for immutable ScHash field specified by key
    getHash(key: MapKey): ScImmutableHash {
        return new ScImmutableHash(this.objID, key.getKeyID());
    }

    // get array proxy for ScImmutableHashArray specified by key
    getHashArray(key: MapKey): ScImmutableHashArray {
        let arrID = host.getObjectID(this.objID, key.getKeyID(), host.TYPE_HASH | host.TYPE_ARRAY);
        return new ScImmutableHashArray(arrID);
    }

    // get value proxy for immutable ScHname field specified by key
    getHname(key: MapKey): ScImmutableHname {
        return new ScImmutableHname(this.objID, key.getKeyID());
    }

    // get array proxy for ScImmutableHnameArray specified by key
    getHnameArray(key: MapKey): ScImmutableHnameArray {
        let arrID = host.getObjectID(this.objID, key.getKeyID(), host.TYPE_HNAME | host.TYPE_ARRAY);
        return new ScImmutableHnameArray(arrID);
    }

    // get value proxy for immutable Int8 field specified by key
    getInt8(key: MapKey): ScImmutableInt8 {
        return new ScImmutableInt8(this.objID, key.getKeyID());
    }

    // get array proxy for ScImmutableInt8Array specified by key
    getInt8Array(key: MapKey): ScImmutableInt8Array {
        let arrID = host.getObjectID(this.objID, key.getKeyID(), host.TYPE_INT8 | host.TYPE_ARRAY);
        return new ScImmutableInt8Array(arrID);
    }

    // get value proxy for immutable Int16 field specified by key
    getInt16(key: MapKey): ScImmutableInt16 {
        return new ScImmutableInt16(this.objID, key.getKeyID());
    }

    // get array proxy for ScImmutableInt16Array specified by key
    getInt16Array(key: MapKey): ScImmutableInt16Array {
        let arrID = host.getObjectID(this.objID, key.getKeyID(), host.TYPE_INT16 | host.TYPE_ARRAY);
        return new ScImmutableInt16Array(arrID);
    }

    // get value proxy for immutable Int32 field specified by key
    getInt32(key: MapKey): ScImmutableInt32 {
        return new ScImmutableInt32(this.objID, key.getKeyID());
    }

    // get array proxy for ScImmutableInt32Array specified by key
    getInt32Array(key: MapKey): ScImmutableInt32Array {
        let arrID = host.getObjectID(this.objID, key.getKeyID(), host.TYPE_INT32 | host.TYPE_ARRAY);
        return new ScImmutableInt32Array(arrID);
    }

    // get value proxy for immutable Int64 field specified by key
    getInt64(key: MapKey): ScImmutableInt64 {
        return new ScImmutableInt64(this.objID, key.getKeyID());
    }

    // get array proxy for ScImmutableInt64Array specified by key
    getInt64Array(key: MapKey): ScImmutableInt64Array {
        let arrID = host.getObjectID(this.objID, key.getKeyID(), host.TYPE_INT64 | host.TYPE_ARRAY);
        return new ScImmutableInt64Array(arrID);
    }

    // get map proxy for ScImmutableMap specified by key
    getMap(key: MapKey): ScImmutableMap {
        let mapID = host.getObjectID(this.objID, key.getKeyID(), host.TYPE_MAP);
        return new ScImmutableMap(mapID);
    }

    // get array proxy for ScImmutableMapArray specified by key
    getMapArray(key: MapKey): ScImmutableMapArray {
        let arrID = host.getObjectID(this.objID, key.getKeyID(), host.TYPE_MAP | host.TYPE_ARRAY);
        return new ScImmutableMapArray(arrID);
    }

    // get value proxy for immutable ScRequestID field specified by key
    getRequestID(key: MapKey): ScImmutableRequestID {
        return new ScImmutableRequestID(this.objID, key.getKeyID());
    }

    // get array proxy for ScImmutableRequestIDArray specified by key
    getRequestIDArray(key: MapKey): ScImmutableRequestIDArray {
        let arrID = host.getObjectID(this.objID, key.getKeyID(), host.TYPE_REQUEST_ID | host.TYPE_ARRAY);
        return new ScImmutableRequestIDArray(arrID);
    }

    // get value proxy for immutable UTF-8 text string field specified by key
    getString(key: MapKey): ScImmutableString {
        return new ScImmutableString(this.objID, key.getKeyID());
    }

    // get array proxy for ScImmutableStringArray specified by key
    getStringArray(key: MapKey): ScImmutableStringArray {
        let arrID = host.getObjectID(this.objID, key.getKeyID(), host.TYPE_STRING | host.TYPE_ARRAY);
        return new ScImmutableStringArray(arrID);
    }

    // get value proxy for immutable Uint8 field specified by key
    getUint8(key: MapKey): ScImmutableUint8 {
        return new ScImmutableUint8(this.objID, key.getKeyID());
    }

    // get array proxy for ScImmutableUint8Array specified by key
    getUint8Array(key: MapKey): ScImmutableUint8Array {
        let arrID = host.getObjectID(this.objID, key.getKeyID(), host.TYPE_INT8 | host.TYPE_ARRAY);
        return new ScImmutableUint8Array(arrID);
    }

    // get value proxy for immutable Uint16 field specified by key
    getUint16(key: MapKey): ScImmutableUint16 {
        return new ScImmutableUint16(this.objID, key.getKeyID());
    }

    // get array proxy for ScImmutableUint16Array specified by key
    getUint16Array(key: MapKey): ScImmutableUint16Array {
        let arrID = host.getObjectID(this.objID, key.getKeyID(), host.TYPE_INT16 | host.TYPE_ARRAY);
        return new ScImmutableUint16Array(arrID);
    }

    // get value proxy for immutable Uint32 field specified by key
    getUint32(key: MapKey): ScImmutableUint32 {
        return new ScImmutableUint32(this.objID, key.getKeyID());
    }

    // get array proxy for ScImmutableUint32Array specified by key
    getUint32Array(key: MapKey): ScImmutableUint32Array {
        let arrID = host.getObjectID(this.objID, key.getKeyID(), host.TYPE_INT32 | host.TYPE_ARRAY);
        return new ScImmutableUint32Array(arrID);
    }

    // get value proxy for immutable Uint64 field specified by key
    getUint64(key: MapKey): ScImmutableUint64 {
        return new ScImmutableUint64(this.objID, key.getKeyID());
    }

    // get array proxy for ScImmutableUint64Array specified by key
    getUint64Array(key: MapKey): ScImmutableUint64Array {
        let arrID = host.getObjectID(this.objID, key.getKeyID(), host.TYPE_INT64 | host.TYPE_ARRAY);
        return new ScImmutableUint64Array(arrID);
    }

    mapID(): i32 {
        return this.objID;
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// array proxy for immutable array of maps
export class ScImmutableMapArray {
    objID: i32;

    constructor(id: i32) {
        this.objID = id;
    }

    // get value proxy for item at index, index can be 0..length()-1
    getMap(index: i32): ScImmutableMap {
        let mapID = host.getObjectID(this.objID, new Key32(index), host.TYPE_MAP);
        return new ScImmutableMap(mapID);
    }

    // number of items in array
    length(): i32 {
        return host.getLength(this.objID);
    }
}

// value proxy for immutable ScRequestID in host container
export class ScImmutableRequestID {
    objID: i32;
    keyID: Key32;

    constructor(objID: i32, keyID: Key32) {
        this.objID = objID;
        this.keyID = keyID;
    }

    // check if value exists in host container
    exists(): boolean {
        return host.exists(this.objID, this.keyID, host.TYPE_REQUEST_ID);
    }

    // human-readable string representation
    toString(): string {
        return this.value().toString();
    }

    // get value from host container
    value(): ScRequestID {
        return ScRequestID.fromBytes(host.getBytes(this.objID, this.keyID, host.TYPE_REQUEST_ID));
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// array proxy for immutable array of ScRequestID
export class ScImmutableRequestIDArray {
    objID: i32;

    constructor(id: i32) {
        this.objID = id;
    }

    // get value proxy for item at index, index can be 0..length()-1
    getRequestID(index: i32): ScImmutableRequestID {
        return new ScImmutableRequestID(this.objID, new Key32(index));
    }

    // number of items in array
    length(): i32 {
        return host.getLength(this.objID);
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// value proxy for immutable UTF-8 text string in host container
export class ScImmutableString {
    objID: i32;
    keyID: Key32;

    constructor(objID: i32, keyID: Key32) {
        this.objID = objID;
        this.keyID = keyID;
    }

    // check if value exists in host container
    exists(): boolean {
        return host.exists(this.objID, this.keyID, host.TYPE_STRING);
    }

    // human-readable string representation
    toString(): string {
        return this.value();
    }

    // get value from host container
    value(): string {
        let bytes = host.getBytes(this.objID, this.keyID, host.TYPE_STRING);
        return Convert.toString(bytes);
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// array proxy for immutable array of UTF-8 text string
export class ScImmutableStringArray {
    objID: i32;

    constructor(id: i32) {
        this.objID = id;
    }

    // get value proxy for item at index, index can be 0..length()-1
    getString(index: i32): ScImmutableString {
        return new ScImmutableString(this.objID, new Key32(index));
    }

    // number of items in array
    length(): i32 {
        return host.getLength(this.objID);
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// value proxy for immutable Uint8 in host container
export class ScImmutableUint8 {
    objID: i32;
    keyID: Key32;

    constructor(objID: i32, keyID: Key32) {
        this.objID = objID;
        this.keyID = keyID;
    }

    // check if value exists in host container
    exists(): boolean {
        return host.exists(this.objID, this.keyID, host.TYPE_INT8);
    }

    // human-readable string representation
    toString(): string {
        return this.value().toString();
    }

    // get value from host container
    value(): u8 {
        let bytes = host.getBytes(this.objID, this.keyID, host.TYPE_INT8);
        return bytes[0] as u8;
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// array proxy for immutable array of Uint8
export class ScImmutableUint8Array {
    objID: i32;

    constructor(id: i32) {
        this.objID = id;
    }

    // get value proxy for item at index, index can be 0..length()-1
    getUint8(index: i32): ScImmutableUint8 {
        return new ScImmutableUint8(this.objID, new Key32(index));
    }

    // number of items in array
    length(): u32 {
        return host.getLength(this.objID);
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// value proxy for immutable Uint16 in host container
export class ScImmutableUint16 {
    objID: i32;
    keyID: Key32;

    constructor(objID: i32, keyID: Key32) {
        this.objID = objID;
        this.keyID = keyID;
    }

    // check if value exists in host container
    exists(): boolean {
        return host.exists(this.objID, this.keyID, host.TYPE_INT16);
    }

    // human-readable string representation
    toString(): string {
        return this.value().toString();
    }

    // get value from host container
    value(): u16 {
        let bytes = host.getBytes(this.objID, this.keyID, host.TYPE_INT16);
        return Convert.toI16(bytes) as u16;
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// array proxy for immutable array of Uint16
export class ScImmutableUint16Array {
    objID: i32;

    constructor(id: i32) {
        this.objID = id;
    }

    // get value proxy for item at index, index can be 0..length()-1
    getUint16(index: i32): ScImmutableUint16 {
        return new ScImmutableUint16(this.objID, new Key32(index));
    }

    // number of items in array
    length(): i32 {
        return host.getLength(this.objID);
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// value proxy for immutable Uint32 in host container
export class ScImmutableUint32 {
    objID: i32;
    keyID: Key32;

    constructor(objID: i32, keyID: Key32) {
        this.objID = objID;
        this.keyID = keyID;
    }

    // check if value exists in host container
    exists(): boolean {
        return host.exists(this.objID, this.keyID, host.TYPE_INT32);
    }

    // human-readable string representation
    toString(): string {
        return this.value().toString();
    }

    // get value from host container
    value(): u32 {
        let bytes = host.getBytes(this.objID, this.keyID, host.TYPE_INT32);
        return Convert.toI32(bytes) as u32;
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// array proxy for immutable array of Uint32
export class ScImmutableUint32Array {
    objID: i32;

    constructor(id: i32) {
        this.objID = id;
    }

    // get value proxy for item at index, index can be 0..length()-1
    getUint32(index: i32): ScImmutableUint32 {
        return new ScImmutableUint32(this.objID, new Key32(index));
    }

    // number of items in array
    length(): i32 {
        return host.getLength(this.objID);
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// value proxy for immutable Uint64 in host container
export class ScImmutableUint64 {
    objID: i32;
    keyID: Key32;

    constructor(objID: i32, keyID: Key32) {
        this.objID = objID;
        this.keyID = keyID;
    }

    // check if value exists in host container
    exists(): boolean {
        return host.exists(this.objID, this.keyID, host.TYPE_INT64);
    }

    // human-readable string representation
    toString(): string {
        return this.value().toString();
    }

    // get value from host container
    value(): u64 {
        let bytes = host.getBytes(this.objID, this.keyID, host.TYPE_INT64);
        return Convert.toI64(bytes) as u64;
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// array proxy for immutable array of Uint64
export class ScImmutableUint64Array {
    objID: i32;

    constructor(id: i32) {
        this.objID = id;
    }

    // get value proxy for item at index, index can be 0..length()-1
    getUint64(index: i32): ScImmutableUint64 {
        return new ScImmutableUint64(this.objID, new Key32(index));
    }

    // number of items in array
    length(): i32 {
        return host.getLength(this.objID);
    }
}
