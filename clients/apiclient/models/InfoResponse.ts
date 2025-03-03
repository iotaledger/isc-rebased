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

import { L1Params } from '../models/L1Params';
import { HttpFile } from '../http/http';

export class InfoResponse {
    'l1Params': L1Params;
    /**
    * The net id of the node
    */
    'peeringURL': string;
    /**
    * The public key of the node (Hex)
    */
    'publicKey': string;
    /**
    * The version of the node
    */
    'version': string;

    static readonly discriminator: string | undefined = undefined;

    static readonly mapping: {[index: string]: string} | undefined = undefined;

    static readonly attributeTypeMap: Array<{name: string, baseName: string, type: string, format: string}> = [
        {
            "name": "l1Params",
            "baseName": "l1Params",
            "type": "L1Params",
            "format": ""
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
        },
        {
            "name": "version",
            "baseName": "version",
            "type": "string",
            "format": "string"
        }    ];

    static getAttributeTypeMap() {
        return InfoResponse.attributeTypeMap;
    }

    public constructor() {
    }
}
