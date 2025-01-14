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

import { AssetsJSON } from '../models/AssetsJSON';
import { CallTargetJSON } from '../models/CallTargetJSON';
import { HttpFile } from '../http/http';

export class RequestJSON {
    'allowance': AssetsJSON;
    'assets': AssetsJSON;
    'callTarget': CallTargetJSON;
    /**
    * The gas budget (uint64 as string)
    */
    'gasBudget': string;
    'isEVM': boolean;
    'isOffLedger': boolean;
    'params': Array<Array<number>>;
    'requestId': string;
    'senderAccount': string;
    'targetAddress': string;

    static readonly discriminator: string | undefined = undefined;

    static readonly mapping: {[index: string]: string} | undefined = undefined;

    static readonly attributeTypeMap: Array<{name: string, baseName: string, type: string, format: string}> = [
        {
            "name": "allowance",
            "baseName": "allowance",
            "type": "AssetsJSON",
            "format": ""
        },
        {
            "name": "assets",
            "baseName": "assets",
            "type": "AssetsJSON",
            "format": ""
        },
        {
            "name": "callTarget",
            "baseName": "callTarget",
            "type": "CallTargetJSON",
            "format": ""
        },
        {
            "name": "gasBudget",
            "baseName": "gasBudget",
            "type": "string",
            "format": "string"
        },
        {
            "name": "isEVM",
            "baseName": "isEVM",
            "type": "boolean",
            "format": "boolean"
        },
        {
            "name": "isOffLedger",
            "baseName": "isOffLedger",
            "type": "boolean",
            "format": "boolean"
        },
        {
            "name": "params",
            "baseName": "params",
            "type": "Array<Array<number>>",
            "format": "int32"
        },
        {
            "name": "requestId",
            "baseName": "requestId",
            "type": "string",
            "format": "string"
        },
        {
            "name": "senderAccount",
            "baseName": "senderAccount",
            "type": "string",
            "format": "string"
        },
        {
            "name": "targetAddress",
            "baseName": "targetAddress",
            "type": "string",
            "format": "string"
        }    ];

    static getAttributeTypeMap() {
        return RequestJSON.attributeTypeMap;
    }

    public constructor() {
    }
}
