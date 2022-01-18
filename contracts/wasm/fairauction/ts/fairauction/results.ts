// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the json schema instead

import * as wasmlib from "wasmlib";
import * as sc from "./index";

export class ImmutableGetInfoResults extends wasmlib.ScMapID {
    bidders(): wasmlib.ScImmutableInt32 {
		return new wasmlib.ScImmutableInt32(this.mapID, wasmlib.Key32.fromString(sc.ResultBidders));
	}

    color(): wasmlib.ScImmutableColor {
		return new wasmlib.ScImmutableColor(this.mapID, wasmlib.Key32.fromString(sc.ResultColor));
	}

    creator(): wasmlib.ScImmutableAgentID {
		return new wasmlib.ScImmutableAgentID(this.mapID, wasmlib.Key32.fromString(sc.ResultCreator));
	}

    deposit(): wasmlib.ScImmutableInt64 {
		return new wasmlib.ScImmutableInt64(this.mapID, wasmlib.Key32.fromString(sc.ResultDeposit));
	}

    description(): wasmlib.ScImmutableString {
		return new wasmlib.ScImmutableString(this.mapID, wasmlib.Key32.fromString(sc.ResultDescription));
	}

    duration(): wasmlib.ScImmutableInt32 {
		return new wasmlib.ScImmutableInt32(this.mapID, wasmlib.Key32.fromString(sc.ResultDuration));
	}

    highestBid(): wasmlib.ScImmutableInt64 {
		return new wasmlib.ScImmutableInt64(this.mapID, wasmlib.Key32.fromString(sc.ResultHighestBid));
	}

    highestBidder(): wasmlib.ScImmutableAgentID {
		return new wasmlib.ScImmutableAgentID(this.mapID, wasmlib.Key32.fromString(sc.ResultHighestBidder));
	}

    minimumBid(): wasmlib.ScImmutableInt64 {
		return new wasmlib.ScImmutableInt64(this.mapID, wasmlib.Key32.fromString(sc.ResultMinimumBid));
	}

    numTokens(): wasmlib.ScImmutableInt64 {
		return new wasmlib.ScImmutableInt64(this.mapID, wasmlib.Key32.fromString(sc.ResultNumTokens));
	}

    ownerMargin(): wasmlib.ScImmutableInt64 {
		return new wasmlib.ScImmutableInt64(this.mapID, wasmlib.Key32.fromString(sc.ResultOwnerMargin));
	}

    whenStarted(): wasmlib.ScImmutableInt64 {
		return new wasmlib.ScImmutableInt64(this.mapID, wasmlib.Key32.fromString(sc.ResultWhenStarted));
	}
}

export class MutableGetInfoResults extends wasmlib.ScMapID {
    bidders(): wasmlib.ScMutableInt32 {
		return new wasmlib.ScMutableInt32(this.mapID, wasmlib.Key32.fromString(sc.ResultBidders));
	}

    color(): wasmlib.ScMutableColor {
		return new wasmlib.ScMutableColor(this.mapID, wasmlib.Key32.fromString(sc.ResultColor));
	}

    creator(): wasmlib.ScMutableAgentID {
		return new wasmlib.ScMutableAgentID(this.mapID, wasmlib.Key32.fromString(sc.ResultCreator));
	}

    deposit(): wasmlib.ScMutableInt64 {
		return new wasmlib.ScMutableInt64(this.mapID, wasmlib.Key32.fromString(sc.ResultDeposit));
	}

    description(): wasmlib.ScMutableString {
		return new wasmlib.ScMutableString(this.mapID, wasmlib.Key32.fromString(sc.ResultDescription));
	}

    duration(): wasmlib.ScMutableInt32 {
		return new wasmlib.ScMutableInt32(this.mapID, wasmlib.Key32.fromString(sc.ResultDuration));
	}

    highestBid(): wasmlib.ScMutableInt64 {
		return new wasmlib.ScMutableInt64(this.mapID, wasmlib.Key32.fromString(sc.ResultHighestBid));
	}

    highestBidder(): wasmlib.ScMutableAgentID {
		return new wasmlib.ScMutableAgentID(this.mapID, wasmlib.Key32.fromString(sc.ResultHighestBidder));
	}

    minimumBid(): wasmlib.ScMutableInt64 {
		return new wasmlib.ScMutableInt64(this.mapID, wasmlib.Key32.fromString(sc.ResultMinimumBid));
	}

    numTokens(): wasmlib.ScMutableInt64 {
		return new wasmlib.ScMutableInt64(this.mapID, wasmlib.Key32.fromString(sc.ResultNumTokens));
	}

    ownerMargin(): wasmlib.ScMutableInt64 {
		return new wasmlib.ScMutableInt64(this.mapID, wasmlib.Key32.fromString(sc.ResultOwnerMargin));
	}

    whenStarted(): wasmlib.ScMutableInt64 {
		return new wasmlib.ScMutableInt64(this.mapID, wasmlib.Key32.fromString(sc.ResultWhenStarted));
	}
}
