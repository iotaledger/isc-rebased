// TODO: better import syntax?
import {BaseAPIRequestFactory, RequiredError, COLLECTION_FORMATS} from './baseapi';
import {Configuration} from '../configuration';
import {RequestContext, HttpMethod, ResponseContext, HttpFile, HttpInfo} from '../http/http';
import {ObjectSerializer} from '../models/ObjectSerializer';
import {ApiException} from './exception';
import {canConsumeForm, isCodeInRange} from '../util';
import {SecurityAuthentication} from '../auth/auth';


import { ChainInfoResponse } from '../models/ChainInfoResponse';
import { ChainRecord } from '../models/ChainRecord';
import { CommitteeInfoResponse } from '../models/CommitteeInfoResponse';
import { ContractCallViewRequest } from '../models/ContractCallViewRequest';
import { ContractInfoResponse } from '../models/ContractInfoResponse';
import { EstimateGasRequestOffledger } from '../models/EstimateGasRequestOffledger';
import { EstimateGasRequestOnledger } from '../models/EstimateGasRequestOnledger';
import { ReceiptResponse } from '../models/ReceiptResponse';
import { RotateChainRequest } from '../models/RotateChainRequest';
import { StateResponse } from '../models/StateResponse';
import { ValidationError } from '../models/ValidationError';

/**
 * no description
 */
export class ChainsApiRequestFactory extends BaseAPIRequestFactory {

    /**
     * Activate a chain
     * @param chainID ChainID (Hex Address)
     */
    public async activateChain(chainID: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'chainID' is not null or undefined
        if (chainID === null || chainID === undefined) {
            throw new RequiredError("ChainsApi", "activateChain", "chainID");
        }


        // Path Params
        const localVarPath = '/v1/chains/{chainID}/activate'
            .replace('{' + 'chainID' + '}', encodeURIComponent(String(chainID)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.POST);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        let authMethod: SecurityAuthentication | undefined;
        // Apply auth methods
        authMethod = _config.authMethods["Authorization"]
        if (authMethod?.applySecurityAuthentication) {
            await authMethod?.applySecurityAuthentication(requestContext);
        }
        
        const defaultAuth: SecurityAuthentication | undefined = _options?.authMethods?.default || this.configuration?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Configure a trusted node to be an access node.
     * @param chainID ChainID (Hex Address)
     * @param peer Name or PubKey (hex) of the trusted peer
     */
    public async addAccessNode(chainID: string, peer: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'chainID' is not null or undefined
        if (chainID === null || chainID === undefined) {
            throw new RequiredError("ChainsApi", "addAccessNode", "chainID");
        }


        // verify required parameter 'peer' is not null or undefined
        if (peer === null || peer === undefined) {
            throw new RequiredError("ChainsApi", "addAccessNode", "peer");
        }


        // Path Params
        const localVarPath = '/v1/chains/{chainID}/access-node/{peer}'
            .replace('{' + 'chainID' + '}', encodeURIComponent(String(chainID)))
            .replace('{' + 'peer' + '}', encodeURIComponent(String(peer)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.PUT);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        let authMethod: SecurityAuthentication | undefined;
        // Apply auth methods
        authMethod = _config.authMethods["Authorization"]
        if (authMethod?.applySecurityAuthentication) {
            await authMethod?.applySecurityAuthentication(requestContext);
        }
        
        const defaultAuth: SecurityAuthentication | undefined = _options?.authMethods?.default || this.configuration?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Execute a view call. Either use HName or Name properties. If both are supplied, HName are used.
     * Call a view function on a contract by Hname
     * @param chainID ChainID (Hex Address)
     * @param contractCallViewRequest Parameters
     */
    public async callView(chainID: string, contractCallViewRequest: ContractCallViewRequest, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'chainID' is not null or undefined
        if (chainID === null || chainID === undefined) {
            throw new RequiredError("ChainsApi", "callView", "chainID");
        }


        // verify required parameter 'contractCallViewRequest' is not null or undefined
        if (contractCallViewRequest === null || contractCallViewRequest === undefined) {
            throw new RequiredError("ChainsApi", "callView", "contractCallViewRequest");
        }


        // Path Params
        const localVarPath = '/v1/chains/{chainID}/callview'
            .replace('{' + 'chainID' + '}', encodeURIComponent(String(chainID)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.POST);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        // Body Params
        const contentType = ObjectSerializer.getPreferredMediaType([
            "application/json"
        ]);
        requestContext.setHeaderParam("Content-Type", contentType);
        const serializedBody = ObjectSerializer.stringify(
            ObjectSerializer.serialize(contractCallViewRequest, "ContractCallViewRequest", ""),
            contentType
        );
        requestContext.setBody(serializedBody);

        
        const defaultAuth: SecurityAuthentication | undefined = _options?.authMethods?.default || this.configuration?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Deactivate a chain
     * @param chainID ChainID (Hex Address)
     */
    public async deactivateChain(chainID: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'chainID' is not null or undefined
        if (chainID === null || chainID === undefined) {
            throw new RequiredError("ChainsApi", "deactivateChain", "chainID");
        }


        // Path Params
        const localVarPath = '/v1/chains/{chainID}/deactivate'
            .replace('{' + 'chainID' + '}', encodeURIComponent(String(chainID)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.POST);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        let authMethod: SecurityAuthentication | undefined;
        // Apply auth methods
        authMethod = _config.authMethods["Authorization"]
        if (authMethod?.applySecurityAuthentication) {
            await authMethod?.applySecurityAuthentication(requestContext);
        }
        
        const defaultAuth: SecurityAuthentication | undefined = _options?.authMethods?.default || this.configuration?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * dump accounts information into a humanly-readable format
     * @param chainID ChainID (Hex Address)
     */
    public async dumpAccounts(chainID: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'chainID' is not null or undefined
        if (chainID === null || chainID === undefined) {
            throw new RequiredError("ChainsApi", "dumpAccounts", "chainID");
        }


        // Path Params
        const localVarPath = '/v1/chains/{chainID}/dump-accounts'
            .replace('{' + 'chainID' + '}', encodeURIComponent(String(chainID)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.POST);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        let authMethod: SecurityAuthentication | undefined;
        // Apply auth methods
        authMethod = _config.authMethods["Authorization"]
        if (authMethod?.applySecurityAuthentication) {
            await authMethod?.applySecurityAuthentication(requestContext);
        }
        
        const defaultAuth: SecurityAuthentication | undefined = _options?.authMethods?.default || this.configuration?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Estimates gas for a given off-ledger ISC request
     * @param chainID ChainID (Hex Address)
     * @param request Request
     */
    public async estimateGasOffledger(chainID: string, request: EstimateGasRequestOffledger, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'chainID' is not null or undefined
        if (chainID === null || chainID === undefined) {
            throw new RequiredError("ChainsApi", "estimateGasOffledger", "chainID");
        }


        // verify required parameter 'request' is not null or undefined
        if (request === null || request === undefined) {
            throw new RequiredError("ChainsApi", "estimateGasOffledger", "request");
        }


        // Path Params
        const localVarPath = '/v1/chains/{chainID}/estimategas-offledger'
            .replace('{' + 'chainID' + '}', encodeURIComponent(String(chainID)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.POST);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        // Body Params
        const contentType = ObjectSerializer.getPreferredMediaType([
            "application/json"
        ]);
        requestContext.setHeaderParam("Content-Type", contentType);
        const serializedBody = ObjectSerializer.stringify(
            ObjectSerializer.serialize(request, "EstimateGasRequestOffledger", ""),
            contentType
        );
        requestContext.setBody(serializedBody);

        
        const defaultAuth: SecurityAuthentication | undefined = _options?.authMethods?.default || this.configuration?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Estimates gas for a given on-ledger ISC request
     * @param chainID ChainID (Hex Address)
     * @param request Request
     */
    public async estimateGasOnledger(chainID: string, request: EstimateGasRequestOnledger, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'chainID' is not null or undefined
        if (chainID === null || chainID === undefined) {
            throw new RequiredError("ChainsApi", "estimateGasOnledger", "chainID");
        }


        // verify required parameter 'request' is not null or undefined
        if (request === null || request === undefined) {
            throw new RequiredError("ChainsApi", "estimateGasOnledger", "request");
        }


        // Path Params
        const localVarPath = '/v1/chains/{chainID}/estimategas-onledger'
            .replace('{' + 'chainID' + '}', encodeURIComponent(String(chainID)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.POST);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        // Body Params
        const contentType = ObjectSerializer.getPreferredMediaType([
            "application/json"
        ]);
        requestContext.setHeaderParam("Content-Type", contentType);
        const serializedBody = ObjectSerializer.stringify(
            ObjectSerializer.serialize(request, "EstimateGasRequestOnledger", ""),
            contentType
        );
        requestContext.setBody(serializedBody);

        
        const defaultAuth: SecurityAuthentication | undefined = _options?.authMethods?.default || this.configuration?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Get information about a specific chain
     * @param chainID ChainID (Hex Address)
     * @param block Block index or trie root
     */
    public async getChainInfo(chainID: string, block?: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'chainID' is not null or undefined
        if (chainID === null || chainID === undefined) {
            throw new RequiredError("ChainsApi", "getChainInfo", "chainID");
        }



        // Path Params
        const localVarPath = '/v1/chains/{chainID}'
            .replace('{' + 'chainID' + '}', encodeURIComponent(String(chainID)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.GET);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")

        // Query Params
        if (block !== undefined) {
            requestContext.setQueryParam("block", ObjectSerializer.serialize(block, "string", "string"));
        }


        
        const defaultAuth: SecurityAuthentication | undefined = _options?.authMethods?.default || this.configuration?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Get a list of all chains
     */
    public async getChains(_options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // Path Params
        const localVarPath = '/v1/chains';

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.GET);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        let authMethod: SecurityAuthentication | undefined;
        // Apply auth methods
        authMethod = _config.authMethods["Authorization"]
        if (authMethod?.applySecurityAuthentication) {
            await authMethod?.applySecurityAuthentication(requestContext);
        }
        
        const defaultAuth: SecurityAuthentication | undefined = _options?.authMethods?.default || this.configuration?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Get information about the deployed committee
     * @param chainID ChainID (Hex Address)
     * @param block Block index or trie root
     */
    public async getCommitteeInfo(chainID: string, block?: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'chainID' is not null or undefined
        if (chainID === null || chainID === undefined) {
            throw new RequiredError("ChainsApi", "getCommitteeInfo", "chainID");
        }



        // Path Params
        const localVarPath = '/v1/chains/{chainID}/committee'
            .replace('{' + 'chainID' + '}', encodeURIComponent(String(chainID)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.GET);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")

        // Query Params
        if (block !== undefined) {
            requestContext.setQueryParam("block", ObjectSerializer.serialize(block, "string", "string"));
        }


        let authMethod: SecurityAuthentication | undefined;
        // Apply auth methods
        authMethod = _config.authMethods["Authorization"]
        if (authMethod?.applySecurityAuthentication) {
            await authMethod?.applySecurityAuthentication(requestContext);
        }
        
        const defaultAuth: SecurityAuthentication | undefined = _options?.authMethods?.default || this.configuration?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Get all available chain contracts
     * @param chainID ChainID (Hex Address)
     * @param block Block index or trie root
     */
    public async getContracts(chainID: string, block?: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'chainID' is not null or undefined
        if (chainID === null || chainID === undefined) {
            throw new RequiredError("ChainsApi", "getContracts", "chainID");
        }



        // Path Params
        const localVarPath = '/v1/chains/{chainID}/contracts'
            .replace('{' + 'chainID' + '}', encodeURIComponent(String(chainID)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.GET);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")

        // Query Params
        if (block !== undefined) {
            requestContext.setQueryParam("block", ObjectSerializer.serialize(block, "string", "string"));
        }


        let authMethod: SecurityAuthentication | undefined;
        // Apply auth methods
        authMethod = _config.authMethods["Authorization"]
        if (authMethod?.applySecurityAuthentication) {
            await authMethod?.applySecurityAuthentication(requestContext);
        }
        
        const defaultAuth: SecurityAuthentication | undefined = _options?.authMethods?.default || this.configuration?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Get the contents of the mempool.
     * @param chainID ChainID (Hex Address)
     */
    public async getMempoolContents(chainID: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'chainID' is not null or undefined
        if (chainID === null || chainID === undefined) {
            throw new RequiredError("ChainsApi", "getMempoolContents", "chainID");
        }


        // Path Params
        const localVarPath = '/v1/chains/{chainID}/mempool'
            .replace('{' + 'chainID' + '}', encodeURIComponent(String(chainID)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.GET);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        let authMethod: SecurityAuthentication | undefined;
        // Apply auth methods
        authMethod = _config.authMethods["Authorization"]
        if (authMethod?.applySecurityAuthentication) {
            await authMethod?.applySecurityAuthentication(requestContext);
        }
        
        const defaultAuth: SecurityAuthentication | undefined = _options?.authMethods?.default || this.configuration?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Get a receipt from a request ID
     * @param chainID ChainID (Hex Address)
     * @param requestID RequestID (Hex)
     */
    public async getReceipt(chainID: string, requestID: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'chainID' is not null or undefined
        if (chainID === null || chainID === undefined) {
            throw new RequiredError("ChainsApi", "getReceipt", "chainID");
        }


        // verify required parameter 'requestID' is not null or undefined
        if (requestID === null || requestID === undefined) {
            throw new RequiredError("ChainsApi", "getReceipt", "requestID");
        }


        // Path Params
        const localVarPath = '/v1/chains/{chainID}/receipts/{requestID}'
            .replace('{' + 'chainID' + '}', encodeURIComponent(String(chainID)))
            .replace('{' + 'requestID' + '}', encodeURIComponent(String(requestID)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.GET);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        
        const defaultAuth: SecurityAuthentication | undefined = _options?.authMethods?.default || this.configuration?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Fetch the raw value associated with the given key in the chain state
     * @param chainID ChainID (Hex Address)
     * @param stateKey State Key (Hex)
     */
    public async getStateValue(chainID: string, stateKey: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'chainID' is not null or undefined
        if (chainID === null || chainID === undefined) {
            throw new RequiredError("ChainsApi", "getStateValue", "chainID");
        }


        // verify required parameter 'stateKey' is not null or undefined
        if (stateKey === null || stateKey === undefined) {
            throw new RequiredError("ChainsApi", "getStateValue", "stateKey");
        }


        // Path Params
        const localVarPath = '/v1/chains/{chainID}/state/{stateKey}'
            .replace('{' + 'chainID' + '}', encodeURIComponent(String(chainID)))
            .replace('{' + 'stateKey' + '}', encodeURIComponent(String(stateKey)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.GET);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        
        const defaultAuth: SecurityAuthentication | undefined = _options?.authMethods?.default || this.configuration?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Remove an access node.
     * @param chainID ChainID (Hex Address)
     * @param peer Name or PubKey (hex) of the trusted peer
     */
    public async removeAccessNode(chainID: string, peer: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'chainID' is not null or undefined
        if (chainID === null || chainID === undefined) {
            throw new RequiredError("ChainsApi", "removeAccessNode", "chainID");
        }


        // verify required parameter 'peer' is not null or undefined
        if (peer === null || peer === undefined) {
            throw new RequiredError("ChainsApi", "removeAccessNode", "peer");
        }


        // Path Params
        const localVarPath = '/v1/chains/{chainID}/access-node/{peer}'
            .replace('{' + 'chainID' + '}', encodeURIComponent(String(chainID)))
            .replace('{' + 'peer' + '}', encodeURIComponent(String(peer)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.DELETE);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        let authMethod: SecurityAuthentication | undefined;
        // Apply auth methods
        authMethod = _config.authMethods["Authorization"]
        if (authMethod?.applySecurityAuthentication) {
            await authMethod?.applySecurityAuthentication(requestContext);
        }
        
        const defaultAuth: SecurityAuthentication | undefined = _options?.authMethods?.default || this.configuration?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Rotate a chain
     * @param chainID ChainID (Hex Address)
     * @param rotateRequest RotateRequest
     */
    public async rotateChain(chainID: string, rotateRequest?: RotateChainRequest, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'chainID' is not null or undefined
        if (chainID === null || chainID === undefined) {
            throw new RequiredError("ChainsApi", "rotateChain", "chainID");
        }



        // Path Params
        const localVarPath = '/v1/chains/{chainID}/rotate'
            .replace('{' + 'chainID' + '}', encodeURIComponent(String(chainID)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.POST);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        // Body Params
        const contentType = ObjectSerializer.getPreferredMediaType([
            "application/json"
        ]);
        requestContext.setHeaderParam("Content-Type", contentType);
        const serializedBody = ObjectSerializer.stringify(
            ObjectSerializer.serialize(rotateRequest, "RotateChainRequest", ""),
            contentType
        );
        requestContext.setBody(serializedBody);

        let authMethod: SecurityAuthentication | undefined;
        // Apply auth methods
        authMethod = _config.authMethods["Authorization"]
        if (authMethod?.applySecurityAuthentication) {
            await authMethod?.applySecurityAuthentication(requestContext);
        }
        
        const defaultAuth: SecurityAuthentication | undefined = _options?.authMethods?.default || this.configuration?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Sets the chain record.
     * @param chainID ChainID (Hex Address)
     * @param chainRecord Chain Record
     */
    public async setChainRecord(chainID: string, chainRecord: ChainRecord, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'chainID' is not null or undefined
        if (chainID === null || chainID === undefined) {
            throw new RequiredError("ChainsApi", "setChainRecord", "chainID");
        }


        // verify required parameter 'chainRecord' is not null or undefined
        if (chainRecord === null || chainRecord === undefined) {
            throw new RequiredError("ChainsApi", "setChainRecord", "chainRecord");
        }


        // Path Params
        const localVarPath = '/v1/chains/{chainID}/chainrecord'
            .replace('{' + 'chainID' + '}', encodeURIComponent(String(chainID)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.POST);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        // Body Params
        const contentType = ObjectSerializer.getPreferredMediaType([
            "application/json"
        ]);
        requestContext.setHeaderParam("Content-Type", contentType);
        const serializedBody = ObjectSerializer.stringify(
            ObjectSerializer.serialize(chainRecord, "ChainRecord", ""),
            contentType
        );
        requestContext.setBody(serializedBody);

        let authMethod: SecurityAuthentication | undefined;
        // Apply auth methods
        authMethod = _config.authMethods["Authorization"]
        if (authMethod?.applySecurityAuthentication) {
            await authMethod?.applySecurityAuthentication(requestContext);
        }
        
        const defaultAuth: SecurityAuthentication | undefined = _options?.authMethods?.default || this.configuration?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Ethereum JSON-RPC
     * @param chainID ChainID (Hex Address)
     */
    public async v1ChainsChainIDEvmPost(chainID: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'chainID' is not null or undefined
        if (chainID === null || chainID === undefined) {
            throw new RequiredError("ChainsApi", "v1ChainsChainIDEvmPost", "chainID");
        }


        // Path Params
        const localVarPath = '/v1/chains/{chainID}/evm'
            .replace('{' + 'chainID' + '}', encodeURIComponent(String(chainID)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.POST);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        
        const defaultAuth: SecurityAuthentication | undefined = _options?.authMethods?.default || this.configuration?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Ethereum JSON-RPC (Websocket transport)
     * @param chainID ChainID (Hex Address)
     */
    public async v1ChainsChainIDEvmWsGet(chainID: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'chainID' is not null or undefined
        if (chainID === null || chainID === undefined) {
            throw new RequiredError("ChainsApi", "v1ChainsChainIDEvmWsGet", "chainID");
        }


        // Path Params
        const localVarPath = '/v1/chains/{chainID}/evm/ws'
            .replace('{' + 'chainID' + '}', encodeURIComponent(String(chainID)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.GET);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        
        const defaultAuth: SecurityAuthentication | undefined = _options?.authMethods?.default || this.configuration?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Wait until the given request has been processed by the node
     * @param chainID ChainID (Hex Address)
     * @param requestID RequestID (Hex)
     * @param timeoutSeconds The timeout in seconds, maximum 60s
     * @param waitForL1Confirmation Wait for the block to be confirmed on L1
     */
    public async waitForRequest(chainID: string, requestID: string, timeoutSeconds?: number, waitForL1Confirmation?: boolean, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'chainID' is not null or undefined
        if (chainID === null || chainID === undefined) {
            throw new RequiredError("ChainsApi", "waitForRequest", "chainID");
        }


        // verify required parameter 'requestID' is not null or undefined
        if (requestID === null || requestID === undefined) {
            throw new RequiredError("ChainsApi", "waitForRequest", "requestID");
        }




        // Path Params
        const localVarPath = '/v1/chains/{chainID}/requests/{requestID}/wait'
            .replace('{' + 'chainID' + '}', encodeURIComponent(String(chainID)))
            .replace('{' + 'requestID' + '}', encodeURIComponent(String(requestID)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.GET);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")

        // Query Params
        if (timeoutSeconds !== undefined) {
            requestContext.setQueryParam("timeoutSeconds", ObjectSerializer.serialize(timeoutSeconds, "number", "int32"));
        }

        // Query Params
        if (waitForL1Confirmation !== undefined) {
            requestContext.setQueryParam("waitForL1Confirmation", ObjectSerializer.serialize(waitForL1Confirmation, "boolean", "boolean"));
        }


        
        const defaultAuth: SecurityAuthentication | undefined = _options?.authMethods?.default || this.configuration?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

}

export class ChainsApiResponseProcessor {

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to activateChain
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async activateChainWithHttpInfo(response: ResponseContext): Promise<HttpInfo<void >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, undefined);
        }
        if (isCodeInRange("304", response.httpStatusCode)) {
            throw new ApiException<undefined>(response.httpStatusCode, "Chain was not activated", undefined, response.headers);
        }
        if (isCodeInRange("401", response.httpStatusCode)) {
            const body: ValidationError = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ValidationError", ""
            ) as ValidationError;
            throw new ApiException<ValidationError>(response.httpStatusCode, "Unauthorized (Wrong permissions, missing token)", body, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: void = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "void", ""
            ) as void;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to addAccessNode
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async addAccessNodeWithHttpInfo(response: ResponseContext): Promise<HttpInfo<void >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("201", response.httpStatusCode)) {
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, undefined);
        }
        if (isCodeInRange("401", response.httpStatusCode)) {
            const body: ValidationError = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ValidationError", ""
            ) as ValidationError;
            throw new ApiException<ValidationError>(response.httpStatusCode, "Unauthorized (Wrong permissions, missing token)", body, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: void = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "void", ""
            ) as void;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to callView
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async callViewWithHttpInfo(response: ResponseContext): Promise<HttpInfo<Array<string> >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: Array<string> = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "Array<string>", "string"
            ) as Array<string>;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: Array<string> = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "Array<string>", "string"
            ) as Array<string>;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to deactivateChain
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async deactivateChainWithHttpInfo(response: ResponseContext): Promise<HttpInfo<void >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, undefined);
        }
        if (isCodeInRange("304", response.httpStatusCode)) {
            throw new ApiException<undefined>(response.httpStatusCode, "Chain was not deactivated", undefined, response.headers);
        }
        if (isCodeInRange("401", response.httpStatusCode)) {
            const body: ValidationError = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ValidationError", ""
            ) as ValidationError;
            throw new ApiException<ValidationError>(response.httpStatusCode, "Unauthorized (Wrong permissions, missing token)", body, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: void = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "void", ""
            ) as void;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to dumpAccounts
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async dumpAccountsWithHttpInfo(response: ResponseContext): Promise<HttpInfo<void >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, undefined);
        }
        if (isCodeInRange("401", response.httpStatusCode)) {
            const body: ValidationError = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ValidationError", ""
            ) as ValidationError;
            throw new ApiException<ValidationError>(response.httpStatusCode, "Unauthorized (Wrong permissions, missing token)", body, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: void = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "void", ""
            ) as void;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to estimateGasOffledger
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async estimateGasOffledgerWithHttpInfo(response: ResponseContext): Promise<HttpInfo<ReceiptResponse >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: ReceiptResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ReceiptResponse", ""
            ) as ReceiptResponse;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: ReceiptResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ReceiptResponse", ""
            ) as ReceiptResponse;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to estimateGasOnledger
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async estimateGasOnledgerWithHttpInfo(response: ResponseContext): Promise<HttpInfo<ReceiptResponse >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: ReceiptResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ReceiptResponse", ""
            ) as ReceiptResponse;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: ReceiptResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ReceiptResponse", ""
            ) as ReceiptResponse;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to getChainInfo
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async getChainInfoWithHttpInfo(response: ResponseContext): Promise<HttpInfo<ChainInfoResponse >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: ChainInfoResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ChainInfoResponse", ""
            ) as ChainInfoResponse;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: ChainInfoResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ChainInfoResponse", ""
            ) as ChainInfoResponse;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to getChains
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async getChainsWithHttpInfo(response: ResponseContext): Promise<HttpInfo<Array<ChainInfoResponse> >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: Array<ChainInfoResponse> = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "Array<ChainInfoResponse>", ""
            ) as Array<ChainInfoResponse>;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("401", response.httpStatusCode)) {
            const body: ValidationError = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ValidationError", ""
            ) as ValidationError;
            throw new ApiException<ValidationError>(response.httpStatusCode, "Unauthorized (Wrong permissions, missing token)", body, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: Array<ChainInfoResponse> = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "Array<ChainInfoResponse>", ""
            ) as Array<ChainInfoResponse>;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to getCommitteeInfo
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async getCommitteeInfoWithHttpInfo(response: ResponseContext): Promise<HttpInfo<CommitteeInfoResponse >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: CommitteeInfoResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "CommitteeInfoResponse", ""
            ) as CommitteeInfoResponse;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("401", response.httpStatusCode)) {
            const body: ValidationError = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ValidationError", ""
            ) as ValidationError;
            throw new ApiException<ValidationError>(response.httpStatusCode, "Unauthorized (Wrong permissions, missing token)", body, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: CommitteeInfoResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "CommitteeInfoResponse", ""
            ) as CommitteeInfoResponse;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to getContracts
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async getContractsWithHttpInfo(response: ResponseContext): Promise<HttpInfo<Array<ContractInfoResponse> >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: Array<ContractInfoResponse> = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "Array<ContractInfoResponse>", ""
            ) as Array<ContractInfoResponse>;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("401", response.httpStatusCode)) {
            const body: ValidationError = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ValidationError", ""
            ) as ValidationError;
            throw new ApiException<ValidationError>(response.httpStatusCode, "Unauthorized (Wrong permissions, missing token)", body, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: Array<ContractInfoResponse> = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "Array<ContractInfoResponse>", ""
            ) as Array<ContractInfoResponse>;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to getMempoolContents
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async getMempoolContentsWithHttpInfo(response: ResponseContext): Promise<HttpInfo<Array<number> >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: Array<number> = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "Array<number>", "int32"
            ) as Array<number>;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("401", response.httpStatusCode)) {
            const body: ValidationError = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ValidationError", "int32"
            ) as ValidationError;
            throw new ApiException<ValidationError>(response.httpStatusCode, "Unauthorized (Wrong permissions, missing token)", body, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: Array<number> = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "Array<number>", "int32"
            ) as Array<number>;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to getReceipt
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async getReceiptWithHttpInfo(response: ResponseContext): Promise<HttpInfo<ReceiptResponse >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: ReceiptResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ReceiptResponse", ""
            ) as ReceiptResponse;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("404", response.httpStatusCode)) {
            throw new ApiException<undefined>(response.httpStatusCode, "Chain or request id not found", undefined, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: ReceiptResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ReceiptResponse", ""
            ) as ReceiptResponse;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to getStateValue
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async getStateValueWithHttpInfo(response: ResponseContext): Promise<HttpInfo<StateResponse >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: StateResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "StateResponse", ""
            ) as StateResponse;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: StateResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "StateResponse", ""
            ) as StateResponse;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to removeAccessNode
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async removeAccessNodeWithHttpInfo(response: ResponseContext): Promise<HttpInfo<void >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, undefined);
        }
        if (isCodeInRange("401", response.httpStatusCode)) {
            const body: ValidationError = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ValidationError", ""
            ) as ValidationError;
            throw new ApiException<ValidationError>(response.httpStatusCode, "Unauthorized (Wrong permissions, missing token)", body, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: void = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "void", ""
            ) as void;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to rotateChain
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async rotateChainWithHttpInfo(response: ResponseContext): Promise<HttpInfo<void >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, undefined);
        }
        if (isCodeInRange("401", response.httpStatusCode)) {
            const body: ValidationError = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ValidationError", ""
            ) as ValidationError;
            throw new ApiException<ValidationError>(response.httpStatusCode, "Unauthorized (Wrong permissions, missing token)", body, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: void = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "void", ""
            ) as void;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to setChainRecord
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async setChainRecordWithHttpInfo(response: ResponseContext): Promise<HttpInfo<void >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("201", response.httpStatusCode)) {
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, undefined);
        }
        if (isCodeInRange("401", response.httpStatusCode)) {
            const body: ValidationError = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ValidationError", ""
            ) as ValidationError;
            throw new ApiException<ValidationError>(response.httpStatusCode, "Unauthorized (Wrong permissions, missing token)", body, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: void = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "void", ""
            ) as void;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to v1ChainsChainIDEvmPost
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async v1ChainsChainIDEvmPostWithHttpInfo(response: ResponseContext): Promise<HttpInfo< void>> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("0", response.httpStatusCode)) {
            throw new ApiException<undefined>(response.httpStatusCode, "successful operation", undefined, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, undefined);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to v1ChainsChainIDEvmWsGet
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async v1ChainsChainIDEvmWsGetWithHttpInfo(response: ResponseContext): Promise<HttpInfo< void>> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("0", response.httpStatusCode)) {
            throw new ApiException<undefined>(response.httpStatusCode, "successful operation", undefined, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, undefined);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to waitForRequest
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async waitForRequestWithHttpInfo(response: ResponseContext): Promise<HttpInfo<ReceiptResponse >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: ReceiptResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ReceiptResponse", ""
            ) as ReceiptResponse;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("404", response.httpStatusCode)) {
            throw new ApiException<undefined>(response.httpStatusCode, "The chain or request id not found", undefined, response.headers);
        }
        if (isCodeInRange("408", response.httpStatusCode)) {
            throw new ApiException<undefined>(response.httpStatusCode, "The waiting time has reached the defined limit", undefined, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: ReceiptResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ReceiptResponse", ""
            ) as ReceiptResponse;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

}
