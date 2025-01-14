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

export class PeeringTrustRequest {
    'name': string;
    /**
    * The peering URL of the peer
    */
    'peeringURL': string;
    /**
    * The peers public key encoded in Hex
    */
    'publicKey': string;

    static readonly discriminator: string | undefined = undefined;

    static readonly mapping: {[index: string]: string} | undefined = undefined;

    static readonly attributeTypeMap: Array<{name: string, baseName: string, type: string, format: string}> = [
        {
            "name": "name",
            "baseName": "name",
            "type": "string",
            "format": "string"
        },
        {
            "name": "peeringURL",
            "baseName": "peeringURL",
            "type": "string",
            "format": "string"
        },
        {
            "name": "publicKey",
            "baseName": "publicKey",
            "type": "string",
            "format": "string"
        }    ];

    static getAttributeTypeMap() {
        return PeeringTrustRequest.attributeTypeMap;
    }

    public constructor() {
    }
}
