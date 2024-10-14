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

export class GovAllowedStateControllerAddressesResponse {
    /**
    * The allowed state controller addresses (Hex Address)
    */
    'addresses'?: Array<string>;

    static readonly discriminator: string | undefined = undefined;

    static readonly attributeTypeMap: Array<{name: string, baseName: string, type: string, format: string}> = [
        {
            "name": "addresses",
            "baseName": "addresses",
            "type": "Array<string>",
            "format": "string"
        }    ];

    static getAttributeTypeMap() {
        return GovAllowedStateControllerAddressesResponse.attributeTypeMap;
    }

    public constructor() {
    }
}

