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

import { AssetsResponse } from '../models/AssetsResponse';
import { HttpFile } from '../http/http';

export class FoundryOutputResponse {
    'assets': AssetsResponse;
    'foundryId': string;

    static readonly discriminator: string | undefined = undefined;

    static readonly mapping: {[index: string]: string} | undefined = undefined;

    static readonly attributeTypeMap: Array<{name: string, baseName: string, type: string, format: string}> = [
        {
            "name": "assets",
            "baseName": "assets",
            "type": "AssetsResponse",
            "format": ""
        },
        {
            "name": "foundryId",
            "baseName": "foundryId",
            "type": "string",
            "format": "string"
        }    ];

    static getAttributeTypeMap() {
        return FoundryOutputResponse.attributeTypeMap;
    }

    public constructor() {
    }
}
