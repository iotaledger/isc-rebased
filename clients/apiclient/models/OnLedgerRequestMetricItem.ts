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

import { OnLedgerRequest } from '../models/OnLedgerRequest';
import { HttpFile } from '../http/http';

export class OnLedgerRequestMetricItem {
    'lastMessage': OnLedgerRequest;
    'messages': number;
    'timestamp': Date;

    static readonly discriminator: string | undefined = undefined;

    static readonly mapping: {[index: string]: string} | undefined = undefined;

    static readonly attributeTypeMap: Array<{name: string, baseName: string, type: string, format: string}> = [
        {
            "name": "lastMessage",
            "baseName": "lastMessage",
            "type": "OnLedgerRequest",
            "format": ""
        },
        {
            "name": "messages",
            "baseName": "messages",
            "type": "number",
            "format": "int32"
        },
        {
            "name": "timestamp",
            "baseName": "timestamp",
            "type": "Date",
            "format": "date-time"
        }    ];

    static getAttributeTypeMap() {
        return OnLedgerRequestMetricItem.attributeTypeMap;
    }

    public constructor() {
    }
}
