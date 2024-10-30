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

import { Type } from '../models/Type';
import { HttpFile } from '../http/http';

export class CoinJSON {
    /**
    * The base tokens (uint64 as string)
    */
    'balance': string;
    'coinType': Type;

    static readonly discriminator: string | undefined = undefined;

    static readonly attributeTypeMap: Array<{name: string, baseName: string, type: string, format: string}> = [
        {
            "name": "balance",
            "baseName": "balance",
            "type": "string",
            "format": "string"
        },
        {
            "name": "coinType",
            "baseName": "coinType",
            "type": "Type",
            "format": ""
        }    ];

    static getAttributeTypeMap() {
        return CoinJSON.attributeTypeMap;
    }

    public constructor() {
    }
}

