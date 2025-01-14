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

export class EstimateGasRequestOffledger {
    /**
    * The address to estimate gas for(Hex)
    */
    'fromAddress': string;
    /**
    * Offledger Request (Hex)
    */
    'requestBytes': string;

    static readonly discriminator: string | undefined = undefined;

    static readonly mapping: {[index: string]: string} | undefined = undefined;

    static readonly attributeTypeMap: Array<{name: string, baseName: string, type: string, format: string}> = [
        {
            "name": "fromAddress",
            "baseName": "fromAddress",
            "type": "string",
            "format": "string"
        },
        {
            "name": "requestBytes",
            "baseName": "requestBytes",
            "type": "string",
            "format": "string"
        }    ];

    static getAttributeTypeMap() {
        return EstimateGasRequestOffledger.attributeTypeMap;
    }

    public constructor() {
    }
}
