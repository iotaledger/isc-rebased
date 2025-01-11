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

export class ChainRecord {
    'accessNodes': Array<string>;
    'isActive': boolean;

    static readonly discriminator: string | undefined = undefined;

    static readonly mapping: {[index: string]: string} | undefined = undefined;

    static readonly attributeTypeMap: Array<{name: string, baseName: string, type: string, format: string}> = [
        {
            "name": "accessNodes",
            "baseName": "accessNodes",
            "type": "Array<string>",
            "format": "string"
        },
        {
            "name": "isActive",
            "baseName": "isActive",
            "type": "boolean",
            "format": "boolean"
        }    ];

    static getAttributeTypeMap() {
        return ChainRecord.attributeTypeMap;
    }

    public constructor() {
    }
}
