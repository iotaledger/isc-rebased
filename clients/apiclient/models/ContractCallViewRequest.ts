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

export class ContractCallViewRequest {
    /**
    * Encoded arguments to be passed to the function
    */
    'arguments': Array<string>;
    'block'?: string;
    /**
    * The contract name as HName (Hex)
    */
    'contractHName': string;
    /**
    * The contract name
    */
    'contractName': string;
    /**
    * The function name as HName (Hex)
    */
    'functionHName': string;
    /**
    * The function name
    */
    'functionName': string;

    static readonly discriminator: string | undefined = undefined;

    static readonly attributeTypeMap: Array<{name: string, baseName: string, type: string, format: string}> = [
        {
            "name": "arguments",
            "baseName": "arguments",
            "type": "Array<string>",
            "format": "string"
        },
        {
            "name": "block",
            "baseName": "block",
            "type": "string",
            "format": "string"
        },
        {
            "name": "contractHName",
            "baseName": "contractHName",
            "type": "string",
            "format": "string"
        },
        {
            "name": "contractName",
            "baseName": "contractName",
            "type": "string",
            "format": "string"
        },
        {
            "name": "functionHName",
            "baseName": "functionHName",
            "type": "string",
            "format": "string"
        },
        {
            "name": "functionName",
            "baseName": "functionName",
            "type": "string",
            "format": "string"
        }    ];

    static getAttributeTypeMap() {
        return ContractCallViewRequest.attributeTypeMap;
    }

    public constructor() {
    }
}

