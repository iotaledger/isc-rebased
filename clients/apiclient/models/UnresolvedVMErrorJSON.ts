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

export class UnresolvedVMErrorJSON {
    'code'?: string;
    'params'?: Array<string>;

    static readonly discriminator: string | undefined = undefined;

    static readonly mapping: {[index: string]: string} | undefined = undefined;

    static readonly attributeTypeMap: Array<{name: string, baseName: string, type: string, format: string}> = [
        {
            "name": "code",
            "baseName": "code",
            "type": "string",
            "format": "string"
        },
        {
            "name": "params",
            "baseName": "params",
            "type": "Array<string>",
            "format": "string"
        }    ];

    static getAttributeTypeMap() {
        return UnresolvedVMErrorJSON.attributeTypeMap;
    }

    public constructor() {
    }
}
