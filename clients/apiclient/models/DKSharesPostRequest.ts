/**
 * Wasp API
 * REST API for the Wasp node
 *
 * OpenAPI spec version: 0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { HttpFile } from '../http/http';

export class DKSharesPostRequest {
    /**
    * Names or hex encoded public keys of trusted peers to run DKG on.
    */
    'peerIdentities': Array<string>;
    /**
    * Should be =< len(PeerPublicIdentities)
    */
    'threshold': number;
    /**
    * Timeout in milliseconds.
    */
    'timeoutMS': number;

    static readonly discriminator: string | undefined = undefined;

    static readonly mapping: {[index: string]: string} | undefined = undefined;

    static readonly attributeTypeMap: Array<{name: string, baseName: string, type: string, format: string}> = [
        {
            "name": "peerIdentities",
            "baseName": "peerIdentities",
            "type": "Array<string>",
            "format": "string"
        },
        {
            "name": "threshold",
            "baseName": "threshold",
            "type": "number",
            "format": "int32"
        },
        {
            "name": "timeoutMS",
            "baseName": "timeoutMS",
            "type": "number",
            "format": "int32"
        }    ];

    static getAttributeTypeMap() {
        return DKSharesPostRequest.attributeTypeMap;
    }

    public constructor() {
    }
}
