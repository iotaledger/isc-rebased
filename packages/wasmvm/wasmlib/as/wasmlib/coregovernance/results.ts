// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the schema definition file instead

import * as wasmtypes from '../wasmtypes';
import * as sc from './index';

export class ArrayOfImmutableAddress extends wasmtypes.ScProxy {

    length(): u32 {
        return this.proxy.length();
    }

    getAddress(index: u32): wasmtypes.ScImmutableAddress {
        return new wasmtypes.ScImmutableAddress(this.proxy.index(index));
    }
}

export class ImmutableGetAllowedStateControllerAddressesResults extends wasmtypes.ScProxy {
    // native contract, so this is an Array16
    allowedStateControllerAddresses(): sc.ArrayOfImmutableAddress {
        return new sc.ArrayOfImmutableAddress(this.proxy.root(sc.ResultAllowedStateControllerAddresses));
    }
}

export class ArrayOfMutableAddress extends wasmtypes.ScProxy {

    appendAddress(): wasmtypes.ScMutableAddress {
        return new wasmtypes.ScMutableAddress(this.proxy.append());
    }

    clear(): void {
        this.proxy.clearArray();
    }

    length(): u32 {
        return this.proxy.length();
    }

    getAddress(index: u32): wasmtypes.ScMutableAddress {
        return new wasmtypes.ScMutableAddress(this.proxy.index(index));
    }
}

export class MutableGetAllowedStateControllerAddressesResults extends wasmtypes.ScProxy {
    // native contract, so this is an Array16
    allowedStateControllerAddresses(): sc.ArrayOfMutableAddress {
        return new sc.ArrayOfMutableAddress(this.proxy.root(sc.ResultAllowedStateControllerAddresses));
    }
}

export class ImmutableGetChainInfoResults extends wasmtypes.ScProxy {
    chainID(): wasmtypes.ScImmutableChainID {
        return new wasmtypes.ScImmutableChainID(this.proxy.root(sc.ResultChainID));
    }

    chainOwnerID(): wasmtypes.ScImmutableAgentID {
        return new wasmtypes.ScImmutableAgentID(this.proxy.root(sc.ResultChainOwnerID));
    }

    customMetadata(): wasmtypes.ScImmutableBytes {
        return new wasmtypes.ScImmutableBytes(this.proxy.root(sc.ResultCustomMetadata));
    }

    gasFeePolicyBytes(): wasmtypes.ScImmutableBytes {
        return new wasmtypes.ScImmutableBytes(this.proxy.root(sc.ResultGasFeePolicyBytes));
    }
}

export class MutableGetChainInfoResults extends wasmtypes.ScProxy {
    chainID(): wasmtypes.ScMutableChainID {
        return new wasmtypes.ScMutableChainID(this.proxy.root(sc.ResultChainID));
    }

    chainOwnerID(): wasmtypes.ScMutableAgentID {
        return new wasmtypes.ScMutableAgentID(this.proxy.root(sc.ResultChainOwnerID));
    }

    customMetadata(): wasmtypes.ScMutableBytes {
        return new wasmtypes.ScMutableBytes(this.proxy.root(sc.ResultCustomMetadata));
    }

    gasFeePolicyBytes(): wasmtypes.ScMutableBytes {
        return new wasmtypes.ScMutableBytes(this.proxy.root(sc.ResultGasFeePolicyBytes));
    }
}

export class MapBytesToImmutableBytes extends wasmtypes.ScProxy {

    getBytes(key: Uint8Array): wasmtypes.ScImmutableBytes {
        return new wasmtypes.ScImmutableBytes(this.proxy.key(wasmtypes.bytesToBytes(key)));
    }
}

export class ImmutableGetChainNodesResults extends wasmtypes.ScProxy {
    accessNodeCandidates(): sc.MapBytesToImmutableBytes {
        return new sc.MapBytesToImmutableBytes(this.proxy.root(sc.ResultAccessNodeCandidates));
    }

    accessNodes(): sc.MapBytesToImmutableBytes {
        return new sc.MapBytesToImmutableBytes(this.proxy.root(sc.ResultAccessNodes));
    }
}

export class MapBytesToMutableBytes extends wasmtypes.ScProxy {

    clear(): void {
        this.proxy.clearMap();
    }

    getBytes(key: Uint8Array): wasmtypes.ScMutableBytes {
        return new wasmtypes.ScMutableBytes(this.proxy.key(wasmtypes.bytesToBytes(key)));
    }
}

export class MutableGetChainNodesResults extends wasmtypes.ScProxy {
    accessNodeCandidates(): sc.MapBytesToMutableBytes {
        return new sc.MapBytesToMutableBytes(this.proxy.root(sc.ResultAccessNodeCandidates));
    }

    accessNodes(): sc.MapBytesToMutableBytes {
        return new sc.MapBytesToMutableBytes(this.proxy.root(sc.ResultAccessNodes));
    }
}

export class ImmutableGetChainOwnerResults extends wasmtypes.ScProxy {
    chainOwner(): wasmtypes.ScImmutableAgentID {
        return new wasmtypes.ScImmutableAgentID(this.proxy.root(sc.ResultChainOwner));
    }
}

export class MutableGetChainOwnerResults extends wasmtypes.ScProxy {
    chainOwner(): wasmtypes.ScMutableAgentID {
        return new wasmtypes.ScMutableAgentID(this.proxy.root(sc.ResultChainOwner));
    }
}

export class ImmutableGetFeePolicyResults extends wasmtypes.ScProxy {
    feePolicyBytes(): wasmtypes.ScImmutableBytes {
        return new wasmtypes.ScImmutableBytes(this.proxy.root(sc.ResultFeePolicyBytes));
    }
}

export class MutableGetFeePolicyResults extends wasmtypes.ScProxy {
    feePolicyBytes(): wasmtypes.ScMutableBytes {
        return new wasmtypes.ScMutableBytes(this.proxy.root(sc.ResultFeePolicyBytes));
    }
}
