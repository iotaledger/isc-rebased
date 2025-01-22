import { ResponseContext, RequestContext, HttpFile, HttpInfo } from '../http/http';
import { Configuration} from '../configuration'

import { AccountFoundriesResponse } from '../models/AccountFoundriesResponse';
import { AccountNFTsResponse } from '../models/AccountNFTsResponse';
import { AccountNonceResponse } from '../models/AccountNonceResponse';
import { AddUserRequest } from '../models/AddUserRequest';
import { AnchorMetricItem } from '../models/AnchorMetricItem';
import { AssetsJSON } from '../models/AssetsJSON';
import { AssetsResponse } from '../models/AssetsResponse';
import { AuthInfoModel } from '../models/AuthInfoModel';
import { BaseToken } from '../models/BaseToken';
import { BlockInfoResponse } from '../models/BlockInfoResponse';
import { BurnRecord } from '../models/BurnRecord';
import { CallTargetJSON } from '../models/CallTargetJSON';
import { ChainInfoResponse } from '../models/ChainInfoResponse';
import { ChainMessageMetrics } from '../models/ChainMessageMetrics';
import { ChainRecord } from '../models/ChainRecord';
import { CoinJSON } from '../models/CoinJSON';
import { CommitteeInfoResponse } from '../models/CommitteeInfoResponse';
import { CommitteeNode } from '../models/CommitteeNode';
import { ConsensusPipeMetrics } from '../models/ConsensusPipeMetrics';
import { ConsensusWorkflowMetrics } from '../models/ConsensusWorkflowMetrics';
import { ContractCallViewRequest } from '../models/ContractCallViewRequest';
import { ContractInfoResponse } from '../models/ContractInfoResponse';
import { ControlAddressesResponse } from '../models/ControlAddressesResponse';
import { DKSharesInfo } from '../models/DKSharesInfo';
import { DKSharesPostRequest } from '../models/DKSharesPostRequest';
import { ErrorMessageFormatResponse } from '../models/ErrorMessageFormatResponse';
import { EstimateGasRequestOffledger } from '../models/EstimateGasRequestOffledger';
import { EstimateGasRequestOnledger } from '../models/EstimateGasRequestOnledger';
import { EventJSON } from '../models/EventJSON';
import { EventsResponse } from '../models/EventsResponse';
import { FeePolicy } from '../models/FeePolicy';
import { FoundryOutputResponse } from '../models/FoundryOutputResponse';
import { GovAllowedStateControllerAddressesResponse } from '../models/GovAllowedStateControllerAddressesResponse';
import { GovChainInfoResponse } from '../models/GovChainInfoResponse';
import { GovChainOwnerResponse } from '../models/GovChainOwnerResponse';
import { GovPublicChainMetadata } from '../models/GovPublicChainMetadata';
import { InfoResponse } from '../models/InfoResponse';
import { L1Params } from '../models/L1Params';
import { Limits } from '../models/Limits';
import { LoginRequest } from '../models/LoginRequest';
import { LoginResponse } from '../models/LoginResponse';
import { NativeTokenIDRegistryResponse } from '../models/NativeTokenIDRegistryResponse';
import { NodeMessageMetrics } from '../models/NodeMessageMetrics';
import { NodeOwnerCertificateResponse } from '../models/NodeOwnerCertificateResponse';
import { OffLedgerRequest } from '../models/OffLedgerRequest';
import { OnLedgerRequest } from '../models/OnLedgerRequest';
import { OnLedgerRequestMetricItem } from '../models/OnLedgerRequestMetricItem';
import { PeeringNodeIdentityResponse } from '../models/PeeringNodeIdentityResponse';
import { PeeringNodeStatusResponse } from '../models/PeeringNodeStatusResponse';
import { PeeringTrustRequest } from '../models/PeeringTrustRequest';
import { ProtocolParameters } from '../models/ProtocolParameters';
import { PublicChainMetadata } from '../models/PublicChainMetadata';
import { PublisherStateTransactionItem } from '../models/PublisherStateTransactionItem';
import { Ratio32 } from '../models/Ratio32';
import { ReceiptResponse } from '../models/ReceiptResponse';
import { RequestIDsResponse } from '../models/RequestIDsResponse';
import { RequestJSON } from '../models/RequestJSON';
import { RequestProcessedResponse } from '../models/RequestProcessedResponse';
import { StateAnchor } from '../models/StateAnchor';
import { StateResponse } from '../models/StateResponse';
import { StateTransaction } from '../models/StateTransaction';
import { Type } from '../models/Type';
import { UnresolvedVMErrorJSON } from '../models/UnresolvedVMErrorJSON';
import { UpdateUserPasswordRequest } from '../models/UpdateUserPasswordRequest';
import { UpdateUserPermissionsRequest } from '../models/UpdateUserPermissionsRequest';
import { User } from '../models/User';
import { ValidationError } from '../models/ValidationError';
import { VersionResponse } from '../models/VersionResponse';
import { ObservableAuthApi } from './ObservableAPI';

import { AuthApiRequestFactory, AuthApiResponseProcessor} from "../apis/AuthApi";
export class PromiseAuthApi {
    private api: ObservableAuthApi

    public constructor(
        configuration: Configuration,
        requestFactory?: AuthApiRequestFactory,
        responseProcessor?: AuthApiResponseProcessor
    ) {
        this.api = new ObservableAuthApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * Get information about the current authentication mode
     */
    public authInfoWithHttpInfo(_options?: Configuration): Promise<HttpInfo<AuthInfoModel>> {
        const result = this.api.authInfoWithHttpInfo(_options);
        return result.toPromise();
    }

    /**
     * Get information about the current authentication mode
     */
    public authInfo(_options?: Configuration): Promise<AuthInfoModel> {
        const result = this.api.authInfo(_options);
        return result.toPromise();
    }

    /**
     * Authenticate towards the node
     * @param loginRequest The login request
     */
    public authenticateWithHttpInfo(loginRequest: LoginRequest, _options?: Configuration): Promise<HttpInfo<LoginResponse>> {
        const result = this.api.authenticateWithHttpInfo(loginRequest, _options);
        return result.toPromise();
    }

    /**
     * Authenticate towards the node
     * @param loginRequest The login request
     */
    public authenticate(loginRequest: LoginRequest, _options?: Configuration): Promise<LoginResponse> {
        const result = this.api.authenticate(loginRequest, _options);
        return result.toPromise();
    }


}



import { ObservableChainsApi } from './ObservableAPI';

import { ChainsApiRequestFactory, ChainsApiResponseProcessor} from "../apis/ChainsApi";
export class PromiseChainsApi {
    private api: ObservableChainsApi

    public constructor(
        configuration: Configuration,
        requestFactory?: ChainsApiRequestFactory,
        responseProcessor?: ChainsApiResponseProcessor
    ) {
        this.api = new ObservableChainsApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * Activate a chain
     * @param chainID ChainID (Hex Address)
     */
    public activateChainWithHttpInfo(chainID: string, _options?: Configuration): Promise<HttpInfo<void>> {
        const result = this.api.activateChainWithHttpInfo(chainID, _options);
        return result.toPromise();
    }

    /**
     * Activate a chain
     * @param chainID ChainID (Hex Address)
     */
    public activateChain(chainID: string, _options?: Configuration): Promise<void> {
        const result = this.api.activateChain(chainID, _options);
        return result.toPromise();
    }

    /**
     * Configure a trusted node to be an access node.
     * @param chainID ChainID (Hex Address)
     * @param peer Name or PubKey (hex) of the trusted peer
     */
    public addAccessNodeWithHttpInfo(chainID: string, peer: string, _options?: Configuration): Promise<HttpInfo<void>> {
        const result = this.api.addAccessNodeWithHttpInfo(chainID, peer, _options);
        return result.toPromise();
    }

    /**
     * Configure a trusted node to be an access node.
     * @param chainID ChainID (Hex Address)
     * @param peer Name or PubKey (hex) of the trusted peer
     */
    public addAccessNode(chainID: string, peer: string, _options?: Configuration): Promise<void> {
        const result = this.api.addAccessNode(chainID, peer, _options);
        return result.toPromise();
    }

    /**
     * Execute a view call. Either use HName or Name properties. If both are supplied, HName are used.
     * Call a view function on a contract by Hname
     * @param chainID ChainID (Hex Address)
     * @param contractCallViewRequest Parameters
     */
    public callViewWithHttpInfo(chainID: string, contractCallViewRequest: ContractCallViewRequest, _options?: Configuration): Promise<HttpInfo<Array<string>>> {
        const result = this.api.callViewWithHttpInfo(chainID, contractCallViewRequest, _options);
        return result.toPromise();
    }

    /**
     * Execute a view call. Either use HName or Name properties. If both are supplied, HName are used.
     * Call a view function on a contract by Hname
     * @param chainID ChainID (Hex Address)
     * @param contractCallViewRequest Parameters
     */
    public callView(chainID: string, contractCallViewRequest: ContractCallViewRequest, _options?: Configuration): Promise<Array<string>> {
        const result = this.api.callView(chainID, contractCallViewRequest, _options);
        return result.toPromise();
    }

    /**
     * Deactivate a chain
     * @param chainID ChainID (Hex Address)
     */
    public deactivateChainWithHttpInfo(chainID: string, _options?: Configuration): Promise<HttpInfo<void>> {
        const result = this.api.deactivateChainWithHttpInfo(chainID, _options);
        return result.toPromise();
    }

    /**
     * Deactivate a chain
     * @param chainID ChainID (Hex Address)
     */
    public deactivateChain(chainID: string, _options?: Configuration): Promise<void> {
        const result = this.api.deactivateChain(chainID, _options);
        return result.toPromise();
    }

    /**
     * dump accounts information into a humanly-readable format
     * @param chainID ChainID (Hex Address)
     */
    public dumpAccountsWithHttpInfo(chainID: string, _options?: Configuration): Promise<HttpInfo<void>> {
        const result = this.api.dumpAccountsWithHttpInfo(chainID, _options);
        return result.toPromise();
    }

    /**
     * dump accounts information into a humanly-readable format
     * @param chainID ChainID (Hex Address)
     */
    public dumpAccounts(chainID: string, _options?: Configuration): Promise<void> {
        const result = this.api.dumpAccounts(chainID, _options);
        return result.toPromise();
    }

    /**
     * Estimates gas for a given off-ledger ISC request
     * @param chainID ChainID (Hex Address)
     * @param request Request
     */
    public estimateGasOffledgerWithHttpInfo(chainID: string, request: EstimateGasRequestOffledger, _options?: Configuration): Promise<HttpInfo<ReceiptResponse>> {
        const result = this.api.estimateGasOffledgerWithHttpInfo(chainID, request, _options);
        return result.toPromise();
    }

    /**
     * Estimates gas for a given off-ledger ISC request
     * @param chainID ChainID (Hex Address)
     * @param request Request
     */
    public estimateGasOffledger(chainID: string, request: EstimateGasRequestOffledger, _options?: Configuration): Promise<ReceiptResponse> {
        const result = this.api.estimateGasOffledger(chainID, request, _options);
        return result.toPromise();
    }

    /**
     * Estimates gas for a given on-ledger ISC request
     * @param chainID ChainID (Hex Address)
     * @param request Request
     */
    public estimateGasOnledgerWithHttpInfo(chainID: string, request: EstimateGasRequestOnledger, _options?: Configuration): Promise<HttpInfo<ReceiptResponse>> {
        const result = this.api.estimateGasOnledgerWithHttpInfo(chainID, request, _options);
        return result.toPromise();
    }

    /**
     * Estimates gas for a given on-ledger ISC request
     * @param chainID ChainID (Hex Address)
     * @param request Request
     */
    public estimateGasOnledger(chainID: string, request: EstimateGasRequestOnledger, _options?: Configuration): Promise<ReceiptResponse> {
        const result = this.api.estimateGasOnledger(chainID, request, _options);
        return result.toPromise();
    }

    /**
     * Get information about a specific chain
     * @param chainID ChainID (Hex Address)
     * @param [block] Block index or trie root
     */
    public getChainInfoWithHttpInfo(chainID: string, block?: string, _options?: Configuration): Promise<HttpInfo<ChainInfoResponse>> {
        const result = this.api.getChainInfoWithHttpInfo(chainID, block, _options);
        return result.toPromise();
    }

    /**
     * Get information about a specific chain
     * @param chainID ChainID (Hex Address)
     * @param [block] Block index or trie root
     */
    public getChainInfo(chainID: string, block?: string, _options?: Configuration): Promise<ChainInfoResponse> {
        const result = this.api.getChainInfo(chainID, block, _options);
        return result.toPromise();
    }

    /**
     * Get a list of all chains
     */
    public getChainsWithHttpInfo(_options?: Configuration): Promise<HttpInfo<Array<ChainInfoResponse>>> {
        const result = this.api.getChainsWithHttpInfo(_options);
        return result.toPromise();
    }

    /**
     * Get a list of all chains
     */
    public getChains(_options?: Configuration): Promise<Array<ChainInfoResponse>> {
        const result = this.api.getChains(_options);
        return result.toPromise();
    }

    /**
     * Get information about the deployed committee
     * @param chainID ChainID (Hex Address)
     * @param [block] Block index or trie root
     */
    public getCommitteeInfoWithHttpInfo(chainID: string, block?: string, _options?: Configuration): Promise<HttpInfo<CommitteeInfoResponse>> {
        const result = this.api.getCommitteeInfoWithHttpInfo(chainID, block, _options);
        return result.toPromise();
    }

    /**
     * Get information about the deployed committee
     * @param chainID ChainID (Hex Address)
     * @param [block] Block index or trie root
     */
    public getCommitteeInfo(chainID: string, block?: string, _options?: Configuration): Promise<CommitteeInfoResponse> {
        const result = this.api.getCommitteeInfo(chainID, block, _options);
        return result.toPromise();
    }

    /**
     * Get all available chain contracts
     * @param chainID ChainID (Hex Address)
     * @param [block] Block index or trie root
     */
    public getContractsWithHttpInfo(chainID: string, block?: string, _options?: Configuration): Promise<HttpInfo<Array<ContractInfoResponse>>> {
        const result = this.api.getContractsWithHttpInfo(chainID, block, _options);
        return result.toPromise();
    }

    /**
     * Get all available chain contracts
     * @param chainID ChainID (Hex Address)
     * @param [block] Block index or trie root
     */
    public getContracts(chainID: string, block?: string, _options?: Configuration): Promise<Array<ContractInfoResponse>> {
        const result = this.api.getContracts(chainID, block, _options);
        return result.toPromise();
    }

    /**
     * Get the contents of the mempool.
     * @param chainID ChainID (Hex Address)
     */
    public getMempoolContentsWithHttpInfo(chainID: string, _options?: Configuration): Promise<HttpInfo<Array<number>>> {
        const result = this.api.getMempoolContentsWithHttpInfo(chainID, _options);
        return result.toPromise();
    }

    /**
     * Get the contents of the mempool.
     * @param chainID ChainID (Hex Address)
     */
    public getMempoolContents(chainID: string, _options?: Configuration): Promise<Array<number>> {
        const result = this.api.getMempoolContents(chainID, _options);
        return result.toPromise();
    }

    /**
     * Get a receipt from a request ID
     * @param chainID ChainID (Hex Address)
     * @param requestID RequestID (Hex)
     */
    public getReceiptWithHttpInfo(chainID: string, requestID: string, _options?: Configuration): Promise<HttpInfo<ReceiptResponse>> {
        const result = this.api.getReceiptWithHttpInfo(chainID, requestID, _options);
        return result.toPromise();
    }

    /**
     * Get a receipt from a request ID
     * @param chainID ChainID (Hex Address)
     * @param requestID RequestID (Hex)
     */
    public getReceipt(chainID: string, requestID: string, _options?: Configuration): Promise<ReceiptResponse> {
        const result = this.api.getReceipt(chainID, requestID, _options);
        return result.toPromise();
    }

    /**
     * Fetch the raw value associated with the given key in the chain state
     * @param chainID ChainID (Hex Address)
     * @param stateKey State Key (Hex)
     */
    public getStateValueWithHttpInfo(chainID: string, stateKey: string, _options?: Configuration): Promise<HttpInfo<StateResponse>> {
        const result = this.api.getStateValueWithHttpInfo(chainID, stateKey, _options);
        return result.toPromise();
    }

    /**
     * Fetch the raw value associated with the given key in the chain state
     * @param chainID ChainID (Hex Address)
     * @param stateKey State Key (Hex)
     */
    public getStateValue(chainID: string, stateKey: string, _options?: Configuration): Promise<StateResponse> {
        const result = this.api.getStateValue(chainID, stateKey, _options);
        return result.toPromise();
    }

    /**
     * Remove an access node.
     * @param chainID ChainID (Hex Address)
     * @param peer Name or PubKey (hex) of the trusted peer
     */
    public removeAccessNodeWithHttpInfo(chainID: string, peer: string, _options?: Configuration): Promise<HttpInfo<void>> {
        const result = this.api.removeAccessNodeWithHttpInfo(chainID, peer, _options);
        return result.toPromise();
    }

    /**
     * Remove an access node.
     * @param chainID ChainID (Hex Address)
     * @param peer Name or PubKey (hex) of the trusted peer
     */
    public removeAccessNode(chainID: string, peer: string, _options?: Configuration): Promise<void> {
        const result = this.api.removeAccessNode(chainID, peer, _options);
        return result.toPromise();
    }

    /**
     * Sets the chain record.
     * @param chainID ChainID (Hex Address)
     * @param chainRecord Chain Record
     */
    public setChainRecordWithHttpInfo(chainID: string, chainRecord: ChainRecord, _options?: Configuration): Promise<HttpInfo<void>> {
        const result = this.api.setChainRecordWithHttpInfo(chainID, chainRecord, _options);
        return result.toPromise();
    }

    /**
     * Sets the chain record.
     * @param chainID ChainID (Hex Address)
     * @param chainRecord Chain Record
     */
    public setChainRecord(chainID: string, chainRecord: ChainRecord, _options?: Configuration): Promise<void> {
        const result = this.api.setChainRecord(chainID, chainRecord, _options);
        return result.toPromise();
    }

    /**
     * Ethereum JSON-RPC
     * @param chainID ChainID (Hex Address)
     */
    public v1ChainsChainIDEvmPostWithHttpInfo(chainID: string, _options?: Configuration): Promise<HttpInfo<void>> {
        const result = this.api.v1ChainsChainIDEvmPostWithHttpInfo(chainID, _options);
        return result.toPromise();
    }

    /**
     * Ethereum JSON-RPC
     * @param chainID ChainID (Hex Address)
     */
    public v1ChainsChainIDEvmPost(chainID: string, _options?: Configuration): Promise<void> {
        const result = this.api.v1ChainsChainIDEvmPost(chainID, _options);
        return result.toPromise();
    }

    /**
     * Ethereum JSON-RPC (Websocket transport)
     * @param chainID ChainID (Hex Address)
     */
    public v1ChainsChainIDEvmWsGetWithHttpInfo(chainID: string, _options?: Configuration): Promise<HttpInfo<void>> {
        const result = this.api.v1ChainsChainIDEvmWsGetWithHttpInfo(chainID, _options);
        return result.toPromise();
    }

    /**
     * Ethereum JSON-RPC (Websocket transport)
     * @param chainID ChainID (Hex Address)
     */
    public v1ChainsChainIDEvmWsGet(chainID: string, _options?: Configuration): Promise<void> {
        const result = this.api.v1ChainsChainIDEvmWsGet(chainID, _options);
        return result.toPromise();
    }

    /**
     * Wait until the given request has been processed by the node
     * @param chainID ChainID (Hex Address)
     * @param requestID RequestID (Hex)
     * @param [timeoutSeconds] The timeout in seconds, maximum 60s
     * @param [waitForL1Confirmation] Wait for the block to be confirmed on L1
     */
    public waitForRequestWithHttpInfo(chainID: string, requestID: string, timeoutSeconds?: number, waitForL1Confirmation?: boolean, _options?: Configuration): Promise<HttpInfo<ReceiptResponse>> {
        const result = this.api.waitForRequestWithHttpInfo(chainID, requestID, timeoutSeconds, waitForL1Confirmation, _options);
        return result.toPromise();
    }

    /**
     * Wait until the given request has been processed by the node
     * @param chainID ChainID (Hex Address)
     * @param requestID RequestID (Hex)
     * @param [timeoutSeconds] The timeout in seconds, maximum 60s
     * @param [waitForL1Confirmation] Wait for the block to be confirmed on L1
     */
    public waitForRequest(chainID: string, requestID: string, timeoutSeconds?: number, waitForL1Confirmation?: boolean, _options?: Configuration): Promise<ReceiptResponse> {
        const result = this.api.waitForRequest(chainID, requestID, timeoutSeconds, waitForL1Confirmation, _options);
        return result.toPromise();
    }


}



import { ObservableCorecontractsApi } from './ObservableAPI';

import { CorecontractsApiRequestFactory, CorecontractsApiResponseProcessor} from "../apis/CorecontractsApi";
export class PromiseCorecontractsApi {
    private api: ObservableCorecontractsApi

    public constructor(
        configuration: Configuration,
        requestFactory?: CorecontractsApiRequestFactory,
        responseProcessor?: CorecontractsApiResponseProcessor
    ) {
        this.api = new ObservableCorecontractsApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * Get all assets belonging to an account
     * @param chainID ChainID (Hex Address)
     * @param agentID AgentID (Hex Address for L1 accounts | Hex for EVM)
     * @param [block] Block index or trie root
     */
    public accountsGetAccountBalanceWithHttpInfo(chainID: string, agentID: string, block?: string, _options?: Configuration): Promise<HttpInfo<AssetsResponse>> {
        const result = this.api.accountsGetAccountBalanceWithHttpInfo(chainID, agentID, block, _options);
        return result.toPromise();
    }

    /**
     * Get all assets belonging to an account
     * @param chainID ChainID (Hex Address)
     * @param agentID AgentID (Hex Address for L1 accounts | Hex for EVM)
     * @param [block] Block index or trie root
     */
    public accountsGetAccountBalance(chainID: string, agentID: string, block?: string, _options?: Configuration): Promise<AssetsResponse> {
        const result = this.api.accountsGetAccountBalance(chainID, agentID, block, _options);
        return result.toPromise();
    }

    /**
     * Get all foundries owned by an account
     * @param chainID ChainID (Hex Address)
     * @param agentID AgentID (Hex Address for L1 accounts, Hex for EVM)
     * @param [block] Block index or trie root
     */
    public accountsGetAccountFoundriesWithHttpInfo(chainID: string, agentID: string, block?: string, _options?: Configuration): Promise<HttpInfo<AccountFoundriesResponse>> {
        const result = this.api.accountsGetAccountFoundriesWithHttpInfo(chainID, agentID, block, _options);
        return result.toPromise();
    }

    /**
     * Get all foundries owned by an account
     * @param chainID ChainID (Hex Address)
     * @param agentID AgentID (Hex Address for L1 accounts, Hex for EVM)
     * @param [block] Block index or trie root
     */
    public accountsGetAccountFoundries(chainID: string, agentID: string, block?: string, _options?: Configuration): Promise<AccountFoundriesResponse> {
        const result = this.api.accountsGetAccountFoundries(chainID, agentID, block, _options);
        return result.toPromise();
    }

    /**
     * Get all NFT ids belonging to an account
     * @param chainID ChainID (Hex Address)
     * @param agentID AgentID (Hex Address for L1 accounts | Hex for EVM)
     * @param [block] Block index or trie root
     */
    public accountsGetAccountNFTIDsWithHttpInfo(chainID: string, agentID: string, block?: string, _options?: Configuration): Promise<HttpInfo<AccountNFTsResponse>> {
        const result = this.api.accountsGetAccountNFTIDsWithHttpInfo(chainID, agentID, block, _options);
        return result.toPromise();
    }

    /**
     * Get all NFT ids belonging to an account
     * @param chainID ChainID (Hex Address)
     * @param agentID AgentID (Hex Address for L1 accounts | Hex for EVM)
     * @param [block] Block index or trie root
     */
    public accountsGetAccountNFTIDs(chainID: string, agentID: string, block?: string, _options?: Configuration): Promise<AccountNFTsResponse> {
        const result = this.api.accountsGetAccountNFTIDs(chainID, agentID, block, _options);
        return result.toPromise();
    }

    /**
     * Get the current nonce of an account
     * @param chainID ChainID (Hex Address)
     * @param agentID AgentID (Hex Address for L1 accounts | Hex for EVM)
     * @param [block] Block index or trie root
     */
    public accountsGetAccountNonceWithHttpInfo(chainID: string, agentID: string, block?: string, _options?: Configuration): Promise<HttpInfo<AccountNonceResponse>> {
        const result = this.api.accountsGetAccountNonceWithHttpInfo(chainID, agentID, block, _options);
        return result.toPromise();
    }

    /**
     * Get the current nonce of an account
     * @param chainID ChainID (Hex Address)
     * @param agentID AgentID (Hex Address for L1 accounts | Hex for EVM)
     * @param [block] Block index or trie root
     */
    public accountsGetAccountNonce(chainID: string, agentID: string, block?: string, _options?: Configuration): Promise<AccountNonceResponse> {
        const result = this.api.accountsGetAccountNonce(chainID, agentID, block, _options);
        return result.toPromise();
    }

    /**
     * Get the foundry output
     * @param chainID ChainID (Hex Address)
     * @param serialNumber Serial Number (uint32)
     * @param [block] Block index or trie root
     */
    public accountsGetFoundryOutputWithHttpInfo(chainID: string, serialNumber: number, block?: string, _options?: Configuration): Promise<HttpInfo<FoundryOutputResponse>> {
        const result = this.api.accountsGetFoundryOutputWithHttpInfo(chainID, serialNumber, block, _options);
        return result.toPromise();
    }

    /**
     * Get the foundry output
     * @param chainID ChainID (Hex Address)
     * @param serialNumber Serial Number (uint32)
     * @param [block] Block index or trie root
     */
    public accountsGetFoundryOutput(chainID: string, serialNumber: number, block?: string, _options?: Configuration): Promise<FoundryOutputResponse> {
        const result = this.api.accountsGetFoundryOutput(chainID, serialNumber, block, _options);
        return result.toPromise();
    }

    /**
     * Get the NFT data by an ID
     * @param chainID ChainID (Hex Address)
     * @param nftID NFT ID (Hex)
     * @param [block] Block index or trie root
     */
    public accountsGetNFTDataWithHttpInfo(chainID: string, nftID: string, block?: string, _options?: Configuration): Promise<HttpInfo<void>> {
        const result = this.api.accountsGetNFTDataWithHttpInfo(chainID, nftID, block, _options);
        return result.toPromise();
    }

    /**
     * Get the NFT data by an ID
     * @param chainID ChainID (Hex Address)
     * @param nftID NFT ID (Hex)
     * @param [block] Block index or trie root
     */
    public accountsGetNFTData(chainID: string, nftID: string, block?: string, _options?: Configuration): Promise<void> {
        const result = this.api.accountsGetNFTData(chainID, nftID, block, _options);
        return result.toPromise();
    }

    /**
     * Get a list of all registries
     * @param chainID ChainID (Hex Address)
     * @param [block] Block index or trie root
     */
    public accountsGetNativeTokenIDRegistryWithHttpInfo(chainID: string, block?: string, _options?: Configuration): Promise<HttpInfo<NativeTokenIDRegistryResponse>> {
        const result = this.api.accountsGetNativeTokenIDRegistryWithHttpInfo(chainID, block, _options);
        return result.toPromise();
    }

    /**
     * Get a list of all registries
     * @param chainID ChainID (Hex Address)
     * @param [block] Block index or trie root
     */
    public accountsGetNativeTokenIDRegistry(chainID: string, block?: string, _options?: Configuration): Promise<NativeTokenIDRegistryResponse> {
        const result = this.api.accountsGetNativeTokenIDRegistry(chainID, block, _options);
        return result.toPromise();
    }

    /**
     * Get all stored assets
     * @param chainID ChainID (Hex Address)
     * @param [block] Block index or trie root
     */
    public accountsGetTotalAssetsWithHttpInfo(chainID: string, block?: string, _options?: Configuration): Promise<HttpInfo<AssetsResponse>> {
        const result = this.api.accountsGetTotalAssetsWithHttpInfo(chainID, block, _options);
        return result.toPromise();
    }

    /**
     * Get all stored assets
     * @param chainID ChainID (Hex Address)
     * @param [block] Block index or trie root
     */
    public accountsGetTotalAssets(chainID: string, block?: string, _options?: Configuration): Promise<AssetsResponse> {
        const result = this.api.accountsGetTotalAssets(chainID, block, _options);
        return result.toPromise();
    }

    /**
     * Get the block info of a certain block index
     * @param chainID ChainID (Hex Address)
     * @param blockIndex BlockIndex (uint32)
     * @param [block] Block index or trie root
     */
    public blocklogGetBlockInfoWithHttpInfo(chainID: string, blockIndex: number, block?: string, _options?: Configuration): Promise<HttpInfo<BlockInfoResponse>> {
        const result = this.api.blocklogGetBlockInfoWithHttpInfo(chainID, blockIndex, block, _options);
        return result.toPromise();
    }

    /**
     * Get the block info of a certain block index
     * @param chainID ChainID (Hex Address)
     * @param blockIndex BlockIndex (uint32)
     * @param [block] Block index or trie root
     */
    public blocklogGetBlockInfo(chainID: string, blockIndex: number, block?: string, _options?: Configuration): Promise<BlockInfoResponse> {
        const result = this.api.blocklogGetBlockInfo(chainID, blockIndex, block, _options);
        return result.toPromise();
    }

    /**
     * Get the control addresses
     * @param chainID ChainID (Hex Address)
     * @param [block] Block index or trie root
     */
    public blocklogGetControlAddressesWithHttpInfo(chainID: string, block?: string, _options?: Configuration): Promise<HttpInfo<ControlAddressesResponse>> {
        const result = this.api.blocklogGetControlAddressesWithHttpInfo(chainID, block, _options);
        return result.toPromise();
    }

    /**
     * Get the control addresses
     * @param chainID ChainID (Hex Address)
     * @param [block] Block index or trie root
     */
    public blocklogGetControlAddresses(chainID: string, block?: string, _options?: Configuration): Promise<ControlAddressesResponse> {
        const result = this.api.blocklogGetControlAddresses(chainID, block, _options);
        return result.toPromise();
    }

    /**
     * Get events of a block
     * @param chainID ChainID (Hex Address)
     * @param blockIndex BlockIndex (uint32)
     * @param [block] Block index or trie root
     */
    public blocklogGetEventsOfBlockWithHttpInfo(chainID: string, blockIndex: number, block?: string, _options?: Configuration): Promise<HttpInfo<EventsResponse>> {
        const result = this.api.blocklogGetEventsOfBlockWithHttpInfo(chainID, blockIndex, block, _options);
        return result.toPromise();
    }

    /**
     * Get events of a block
     * @param chainID ChainID (Hex Address)
     * @param blockIndex BlockIndex (uint32)
     * @param [block] Block index or trie root
     */
    public blocklogGetEventsOfBlock(chainID: string, blockIndex: number, block?: string, _options?: Configuration): Promise<EventsResponse> {
        const result = this.api.blocklogGetEventsOfBlock(chainID, blockIndex, block, _options);
        return result.toPromise();
    }

    /**
     * Get events of the latest block
     * @param chainID ChainID (Hex Address)
     * @param [block] Block index or trie root
     */
    public blocklogGetEventsOfLatestBlockWithHttpInfo(chainID: string, block?: string, _options?: Configuration): Promise<HttpInfo<EventsResponse>> {
        const result = this.api.blocklogGetEventsOfLatestBlockWithHttpInfo(chainID, block, _options);
        return result.toPromise();
    }

    /**
     * Get events of the latest block
     * @param chainID ChainID (Hex Address)
     * @param [block] Block index or trie root
     */
    public blocklogGetEventsOfLatestBlock(chainID: string, block?: string, _options?: Configuration): Promise<EventsResponse> {
        const result = this.api.blocklogGetEventsOfLatestBlock(chainID, block, _options);
        return result.toPromise();
    }

    /**
     * Get events of a request
     * @param chainID ChainID (Hex Address)
     * @param requestID RequestID (Hex)
     * @param [block] Block index or trie root
     */
    public blocklogGetEventsOfRequestWithHttpInfo(chainID: string, requestID: string, block?: string, _options?: Configuration): Promise<HttpInfo<EventsResponse>> {
        const result = this.api.blocklogGetEventsOfRequestWithHttpInfo(chainID, requestID, block, _options);
        return result.toPromise();
    }

    /**
     * Get events of a request
     * @param chainID ChainID (Hex Address)
     * @param requestID RequestID (Hex)
     * @param [block] Block index or trie root
     */
    public blocklogGetEventsOfRequest(chainID: string, requestID: string, block?: string, _options?: Configuration): Promise<EventsResponse> {
        const result = this.api.blocklogGetEventsOfRequest(chainID, requestID, block, _options);
        return result.toPromise();
    }

    /**
     * Get the block info of the latest block
     * @param chainID ChainID (Hex Address)
     * @param [block] Block index or trie root
     */
    public blocklogGetLatestBlockInfoWithHttpInfo(chainID: string, block?: string, _options?: Configuration): Promise<HttpInfo<BlockInfoResponse>> {
        const result = this.api.blocklogGetLatestBlockInfoWithHttpInfo(chainID, block, _options);
        return result.toPromise();
    }

    /**
     * Get the block info of the latest block
     * @param chainID ChainID (Hex Address)
     * @param [block] Block index or trie root
     */
    public blocklogGetLatestBlockInfo(chainID: string, block?: string, _options?: Configuration): Promise<BlockInfoResponse> {
        const result = this.api.blocklogGetLatestBlockInfo(chainID, block, _options);
        return result.toPromise();
    }

    /**
     * Get the request ids for a certain block index
     * @param chainID ChainID (Hex Address)
     * @param blockIndex BlockIndex (uint32)
     * @param [block] Block index or trie root
     */
    public blocklogGetRequestIDsForBlockWithHttpInfo(chainID: string, blockIndex: number, block?: string, _options?: Configuration): Promise<HttpInfo<RequestIDsResponse>> {
        const result = this.api.blocklogGetRequestIDsForBlockWithHttpInfo(chainID, blockIndex, block, _options);
        return result.toPromise();
    }

    /**
     * Get the request ids for a certain block index
     * @param chainID ChainID (Hex Address)
     * @param blockIndex BlockIndex (uint32)
     * @param [block] Block index or trie root
     */
    public blocklogGetRequestIDsForBlock(chainID: string, blockIndex: number, block?: string, _options?: Configuration): Promise<RequestIDsResponse> {
        const result = this.api.blocklogGetRequestIDsForBlock(chainID, blockIndex, block, _options);
        return result.toPromise();
    }

    /**
     * Get the request ids for the latest block
     * @param chainID ChainID (Hex Address)
     * @param [block] Block index or trie root
     */
    public blocklogGetRequestIDsForLatestBlockWithHttpInfo(chainID: string, block?: string, _options?: Configuration): Promise<HttpInfo<RequestIDsResponse>> {
        const result = this.api.blocklogGetRequestIDsForLatestBlockWithHttpInfo(chainID, block, _options);
        return result.toPromise();
    }

    /**
     * Get the request ids for the latest block
     * @param chainID ChainID (Hex Address)
     * @param [block] Block index or trie root
     */
    public blocklogGetRequestIDsForLatestBlock(chainID: string, block?: string, _options?: Configuration): Promise<RequestIDsResponse> {
        const result = this.api.blocklogGetRequestIDsForLatestBlock(chainID, block, _options);
        return result.toPromise();
    }

    /**
     * Get the request processing status
     * @param chainID ChainID (Hex Address)
     * @param requestID RequestID (Hex)
     * @param [block] Block index or trie root
     */
    public blocklogGetRequestIsProcessedWithHttpInfo(chainID: string, requestID: string, block?: string, _options?: Configuration): Promise<HttpInfo<RequestProcessedResponse>> {
        const result = this.api.blocklogGetRequestIsProcessedWithHttpInfo(chainID, requestID, block, _options);
        return result.toPromise();
    }

    /**
     * Get the request processing status
     * @param chainID ChainID (Hex Address)
     * @param requestID RequestID (Hex)
     * @param [block] Block index or trie root
     */
    public blocklogGetRequestIsProcessed(chainID: string, requestID: string, block?: string, _options?: Configuration): Promise<RequestProcessedResponse> {
        const result = this.api.blocklogGetRequestIsProcessed(chainID, requestID, block, _options);
        return result.toPromise();
    }

    /**
     * Get the receipt of a certain request id
     * @param chainID ChainID (Hex Address)
     * @param requestID RequestID (Hex)
     * @param [block] Block index or trie root
     */
    public blocklogGetRequestReceiptWithHttpInfo(chainID: string, requestID: string, block?: string, _options?: Configuration): Promise<HttpInfo<ReceiptResponse>> {
        const result = this.api.blocklogGetRequestReceiptWithHttpInfo(chainID, requestID, block, _options);
        return result.toPromise();
    }

    /**
     * Get the receipt of a certain request id
     * @param chainID ChainID (Hex Address)
     * @param requestID RequestID (Hex)
     * @param [block] Block index or trie root
     */
    public blocklogGetRequestReceipt(chainID: string, requestID: string, block?: string, _options?: Configuration): Promise<ReceiptResponse> {
        const result = this.api.blocklogGetRequestReceipt(chainID, requestID, block, _options);
        return result.toPromise();
    }

    /**
     * Get all receipts of a certain block
     * @param chainID ChainID (Hex Address)
     * @param blockIndex BlockIndex (uint32)
     * @param [block] Block index or trie root
     */
    public blocklogGetRequestReceiptsOfBlockWithHttpInfo(chainID: string, blockIndex: number, block?: string, _options?: Configuration): Promise<HttpInfo<Array<ReceiptResponse>>> {
        const result = this.api.blocklogGetRequestReceiptsOfBlockWithHttpInfo(chainID, blockIndex, block, _options);
        return result.toPromise();
    }

    /**
     * Get all receipts of a certain block
     * @param chainID ChainID (Hex Address)
     * @param blockIndex BlockIndex (uint32)
     * @param [block] Block index or trie root
     */
    public blocklogGetRequestReceiptsOfBlock(chainID: string, blockIndex: number, block?: string, _options?: Configuration): Promise<Array<ReceiptResponse>> {
        const result = this.api.blocklogGetRequestReceiptsOfBlock(chainID, blockIndex, block, _options);
        return result.toPromise();
    }

    /**
     * Get all receipts of the latest block
     * @param chainID ChainID (Hex Address)
     * @param [block] Block index or trie root
     */
    public blocklogGetRequestReceiptsOfLatestBlockWithHttpInfo(chainID: string, block?: string, _options?: Configuration): Promise<HttpInfo<Array<ReceiptResponse>>> {
        const result = this.api.blocklogGetRequestReceiptsOfLatestBlockWithHttpInfo(chainID, block, _options);
        return result.toPromise();
    }

    /**
     * Get all receipts of the latest block
     * @param chainID ChainID (Hex Address)
     * @param [block] Block index or trie root
     */
    public blocklogGetRequestReceiptsOfLatestBlock(chainID: string, block?: string, _options?: Configuration): Promise<Array<ReceiptResponse>> {
        const result = this.api.blocklogGetRequestReceiptsOfLatestBlock(chainID, block, _options);
        return result.toPromise();
    }

    /**
     * Get the error message format of a specific error id
     * @param chainID ChainID (Hex Address)
     * @param contractHname Contract (Hname as Hex)
     * @param errorID Error Id (uint16)
     * @param [block] Block index or trie root
     */
    public errorsGetErrorMessageFormatWithHttpInfo(chainID: string, contractHname: string, errorID: number, block?: string, _options?: Configuration): Promise<HttpInfo<ErrorMessageFormatResponse>> {
        const result = this.api.errorsGetErrorMessageFormatWithHttpInfo(chainID, contractHname, errorID, block, _options);
        return result.toPromise();
    }

    /**
     * Get the error message format of a specific error id
     * @param chainID ChainID (Hex Address)
     * @param contractHname Contract (Hname as Hex)
     * @param errorID Error Id (uint16)
     * @param [block] Block index or trie root
     */
    public errorsGetErrorMessageFormat(chainID: string, contractHname: string, errorID: number, block?: string, _options?: Configuration): Promise<ErrorMessageFormatResponse> {
        const result = this.api.errorsGetErrorMessageFormat(chainID, contractHname, errorID, block, _options);
        return result.toPromise();
    }

    /**
     * Returns the allowed state controller addresses
     * Get the allowed state controller addresses
     * @param chainID ChainID (Hex Address)
     * @param [block] Block index or trie root
     */
    public governanceGetAllowedStateControllerAddressesWithHttpInfo(chainID: string, block?: string, _options?: Configuration): Promise<HttpInfo<GovAllowedStateControllerAddressesResponse>> {
        const result = this.api.governanceGetAllowedStateControllerAddressesWithHttpInfo(chainID, block, _options);
        return result.toPromise();
    }

    /**
     * Returns the allowed state controller addresses
     * Get the allowed state controller addresses
     * @param chainID ChainID (Hex Address)
     * @param [block] Block index or trie root
     */
    public governanceGetAllowedStateControllerAddresses(chainID: string, block?: string, _options?: Configuration): Promise<GovAllowedStateControllerAddressesResponse> {
        const result = this.api.governanceGetAllowedStateControllerAddresses(chainID, block, _options);
        return result.toPromise();
    }

    /**
     * If you are using the common API functions, you most likely rather want to use \'/v1/chains/:chainID\' to get information about a chain.
     * Get the chain info
     * @param chainID ChainID (Hex Address)
     * @param [block] Block index or trie root
     */
    public governanceGetChainInfoWithHttpInfo(chainID: string, block?: string, _options?: Configuration): Promise<HttpInfo<GovChainInfoResponse>> {
        const result = this.api.governanceGetChainInfoWithHttpInfo(chainID, block, _options);
        return result.toPromise();
    }

    /**
     * If you are using the common API functions, you most likely rather want to use \'/v1/chains/:chainID\' to get information about a chain.
     * Get the chain info
     * @param chainID ChainID (Hex Address)
     * @param [block] Block index or trie root
     */
    public governanceGetChainInfo(chainID: string, block?: string, _options?: Configuration): Promise<GovChainInfoResponse> {
        const result = this.api.governanceGetChainInfo(chainID, block, _options);
        return result.toPromise();
    }

    /**
     * Returns the chain owner
     * Get the chain owner
     * @param chainID ChainID (Hex Address)
     * @param [block] Block index or trie root
     */
    public governanceGetChainOwnerWithHttpInfo(chainID: string, block?: string, _options?: Configuration): Promise<HttpInfo<GovChainOwnerResponse>> {
        const result = this.api.governanceGetChainOwnerWithHttpInfo(chainID, block, _options);
        return result.toPromise();
    }

    /**
     * Returns the chain owner
     * Get the chain owner
     * @param chainID ChainID (Hex Address)
     * @param [block] Block index or trie root
     */
    public governanceGetChainOwner(chainID: string, block?: string, _options?: Configuration): Promise<GovChainOwnerResponse> {
        const result = this.api.governanceGetChainOwner(chainID, block, _options);
        return result.toPromise();
    }


}



import { ObservableDefaultApi } from './ObservableAPI';

import { DefaultApiRequestFactory, DefaultApiResponseProcessor} from "../apis/DefaultApi";
export class PromiseDefaultApi {
    private api: ObservableDefaultApi

    public constructor(
        configuration: Configuration,
        requestFactory?: DefaultApiRequestFactory,
        responseProcessor?: DefaultApiResponseProcessor
    ) {
        this.api = new ObservableDefaultApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * Returns 200 if the node is healthy.
     */
    public getHealthWithHttpInfo(_options?: Configuration): Promise<HttpInfo<void>> {
        const result = this.api.getHealthWithHttpInfo(_options);
        return result.toPromise();
    }

    /**
     * Returns 200 if the node is healthy.
     */
    public getHealth(_options?: Configuration): Promise<void> {
        const result = this.api.getHealth(_options);
        return result.toPromise();
    }

    /**
     * The websocket connection service
     */
    public v1WsGetWithHttpInfo(_options?: Configuration): Promise<HttpInfo<void>> {
        const result = this.api.v1WsGetWithHttpInfo(_options);
        return result.toPromise();
    }

    /**
     * The websocket connection service
     */
    public v1WsGet(_options?: Configuration): Promise<void> {
        const result = this.api.v1WsGet(_options);
        return result.toPromise();
    }


}



import { ObservableMetricsApi } from './ObservableAPI';

import { MetricsApiRequestFactory, MetricsApiResponseProcessor} from "../apis/MetricsApi";
export class PromiseMetricsApi {
    private api: ObservableMetricsApi

    public constructor(
        configuration: Configuration,
        requestFactory?: MetricsApiRequestFactory,
        responseProcessor?: MetricsApiResponseProcessor
    ) {
        this.api = new ObservableMetricsApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * Get chain specific message metrics.
     * @param chainID ChainID (Hex Address)
     */
    public getChainMessageMetricsWithHttpInfo(chainID: string, _options?: Configuration): Promise<HttpInfo<ChainMessageMetrics>> {
        const result = this.api.getChainMessageMetricsWithHttpInfo(chainID, _options);
        return result.toPromise();
    }

    /**
     * Get chain specific message metrics.
     * @param chainID ChainID (Hex Address)
     */
    public getChainMessageMetrics(chainID: string, _options?: Configuration): Promise<ChainMessageMetrics> {
        const result = this.api.getChainMessageMetrics(chainID, _options);
        return result.toPromise();
    }

    /**
     * Get chain pipe event metrics.
     * @param chainID ChainID (Hex Address)
     */
    public getChainPipeMetricsWithHttpInfo(chainID: string, _options?: Configuration): Promise<HttpInfo<ConsensusPipeMetrics>> {
        const result = this.api.getChainPipeMetricsWithHttpInfo(chainID, _options);
        return result.toPromise();
    }

    /**
     * Get chain pipe event metrics.
     * @param chainID ChainID (Hex Address)
     */
    public getChainPipeMetrics(chainID: string, _options?: Configuration): Promise<ConsensusPipeMetrics> {
        const result = this.api.getChainPipeMetrics(chainID, _options);
        return result.toPromise();
    }

    /**
     * Get chain workflow metrics.
     * @param chainID ChainID (Hex Address)
     */
    public getChainWorkflowMetricsWithHttpInfo(chainID: string, _options?: Configuration): Promise<HttpInfo<ConsensusWorkflowMetrics>> {
        const result = this.api.getChainWorkflowMetricsWithHttpInfo(chainID, _options);
        return result.toPromise();
    }

    /**
     * Get chain workflow metrics.
     * @param chainID ChainID (Hex Address)
     */
    public getChainWorkflowMetrics(chainID: string, _options?: Configuration): Promise<ConsensusWorkflowMetrics> {
        const result = this.api.getChainWorkflowMetrics(chainID, _options);
        return result.toPromise();
    }

    /**
     * Get accumulated message metrics.
     */
    public getNodeMessageMetricsWithHttpInfo(_options?: Configuration): Promise<HttpInfo<NodeMessageMetrics>> {
        const result = this.api.getNodeMessageMetricsWithHttpInfo(_options);
        return result.toPromise();
    }

    /**
     * Get accumulated message metrics.
     */
    public getNodeMessageMetrics(_options?: Configuration): Promise<NodeMessageMetrics> {
        const result = this.api.getNodeMessageMetrics(_options);
        return result.toPromise();
    }


}



import { ObservableNodeApi } from './ObservableAPI';

import { NodeApiRequestFactory, NodeApiResponseProcessor} from "../apis/NodeApi";
export class PromiseNodeApi {
    private api: ObservableNodeApi

    public constructor(
        configuration: Configuration,
        requestFactory?: NodeApiRequestFactory,
        responseProcessor?: NodeApiResponseProcessor
    ) {
        this.api = new ObservableNodeApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * Distrust a peering node
     * @param peer Name or PubKey (hex) of the trusted peer
     */
    public distrustPeerWithHttpInfo(peer: string, _options?: Configuration): Promise<HttpInfo<void>> {
        const result = this.api.distrustPeerWithHttpInfo(peer, _options);
        return result.toPromise();
    }

    /**
     * Distrust a peering node
     * @param peer Name or PubKey (hex) of the trusted peer
     */
    public distrustPeer(peer: string, _options?: Configuration): Promise<void> {
        const result = this.api.distrustPeer(peer, _options);
        return result.toPromise();
    }

    /**
     * Generate a new distributed key
     * @param dKSharesPostRequest Request parameters
     */
    public generateDKSWithHttpInfo(dKSharesPostRequest: DKSharesPostRequest, _options?: Configuration): Promise<HttpInfo<DKSharesInfo>> {
        const result = this.api.generateDKSWithHttpInfo(dKSharesPostRequest, _options);
        return result.toPromise();
    }

    /**
     * Generate a new distributed key
     * @param dKSharesPostRequest Request parameters
     */
    public generateDKS(dKSharesPostRequest: DKSharesPostRequest, _options?: Configuration): Promise<DKSharesInfo> {
        const result = this.api.generateDKS(dKSharesPostRequest, _options);
        return result.toPromise();
    }

    /**
     * Get basic information about all configured peers
     */
    public getAllPeersWithHttpInfo(_options?: Configuration): Promise<HttpInfo<Array<PeeringNodeStatusResponse>>> {
        const result = this.api.getAllPeersWithHttpInfo(_options);
        return result.toPromise();
    }

    /**
     * Get basic information about all configured peers
     */
    public getAllPeers(_options?: Configuration): Promise<Array<PeeringNodeStatusResponse>> {
        const result = this.api.getAllPeers(_options);
        return result.toPromise();
    }

    /**
     * Return the Wasp configuration
     */
    public getConfigurationWithHttpInfo(_options?: Configuration): Promise<HttpInfo<{ [key: string]: string; }>> {
        const result = this.api.getConfigurationWithHttpInfo(_options);
        return result.toPromise();
    }

    /**
     * Return the Wasp configuration
     */
    public getConfiguration(_options?: Configuration): Promise<{ [key: string]: string; }> {
        const result = this.api.getConfiguration(_options);
        return result.toPromise();
    }

    /**
     * Get information about the shared address DKS configuration
     * @param sharedAddress SharedAddress (Hex Address)
     */
    public getDKSInfoWithHttpInfo(sharedAddress: string, _options?: Configuration): Promise<HttpInfo<DKSharesInfo>> {
        const result = this.api.getDKSInfoWithHttpInfo(sharedAddress, _options);
        return result.toPromise();
    }

    /**
     * Get information about the shared address DKS configuration
     * @param sharedAddress SharedAddress (Hex Address)
     */
    public getDKSInfo(sharedAddress: string, _options?: Configuration): Promise<DKSharesInfo> {
        const result = this.api.getDKSInfo(sharedAddress, _options);
        return result.toPromise();
    }

    /**
     * Returns private information about this node.
     */
    public getInfoWithHttpInfo(_options?: Configuration): Promise<HttpInfo<InfoResponse>> {
        const result = this.api.getInfoWithHttpInfo(_options);
        return result.toPromise();
    }

    /**
     * Returns private information about this node.
     */
    public getInfo(_options?: Configuration): Promise<InfoResponse> {
        const result = this.api.getInfo(_options);
        return result.toPromise();
    }

    /**
     * Get basic peer info of the current node
     */
    public getPeeringIdentityWithHttpInfo(_options?: Configuration): Promise<HttpInfo<PeeringNodeIdentityResponse>> {
        const result = this.api.getPeeringIdentityWithHttpInfo(_options);
        return result.toPromise();
    }

    /**
     * Get basic peer info of the current node
     */
    public getPeeringIdentity(_options?: Configuration): Promise<PeeringNodeIdentityResponse> {
        const result = this.api.getPeeringIdentity(_options);
        return result.toPromise();
    }

    /**
     * Get trusted peers
     */
    public getTrustedPeersWithHttpInfo(_options?: Configuration): Promise<HttpInfo<Array<PeeringNodeIdentityResponse>>> {
        const result = this.api.getTrustedPeersWithHttpInfo(_options);
        return result.toPromise();
    }

    /**
     * Get trusted peers
     */
    public getTrustedPeers(_options?: Configuration): Promise<Array<PeeringNodeIdentityResponse>> {
        const result = this.api.getTrustedPeers(_options);
        return result.toPromise();
    }

    /**
     * Returns the node version.
     */
    public getVersionWithHttpInfo(_options?: Configuration): Promise<HttpInfo<VersionResponse>> {
        const result = this.api.getVersionWithHttpInfo(_options);
        return result.toPromise();
    }

    /**
     * Returns the node version.
     */
    public getVersion(_options?: Configuration): Promise<VersionResponse> {
        const result = this.api.getVersion(_options);
        return result.toPromise();
    }

    /**
     * Gets the node owner
     */
    public ownerCertificateWithHttpInfo(_options?: Configuration): Promise<HttpInfo<NodeOwnerCertificateResponse>> {
        const result = this.api.ownerCertificateWithHttpInfo(_options);
        return result.toPromise();
    }

    /**
     * Gets the node owner
     */
    public ownerCertificate(_options?: Configuration): Promise<NodeOwnerCertificateResponse> {
        const result = this.api.ownerCertificate(_options);
        return result.toPromise();
    }

    /**
     * Shut down the node
     */
    public shutdownNodeWithHttpInfo(_options?: Configuration): Promise<HttpInfo<void>> {
        const result = this.api.shutdownNodeWithHttpInfo(_options);
        return result.toPromise();
    }

    /**
     * Shut down the node
     */
    public shutdownNode(_options?: Configuration): Promise<void> {
        const result = this.api.shutdownNode(_options);
        return result.toPromise();
    }

    /**
     * Trust a peering node
     * @param peeringTrustRequest Info of the peer to trust
     */
    public trustPeerWithHttpInfo(peeringTrustRequest: PeeringTrustRequest, _options?: Configuration): Promise<HttpInfo<void>> {
        const result = this.api.trustPeerWithHttpInfo(peeringTrustRequest, _options);
        return result.toPromise();
    }

    /**
     * Trust a peering node
     * @param peeringTrustRequest Info of the peer to trust
     */
    public trustPeer(peeringTrustRequest: PeeringTrustRequest, _options?: Configuration): Promise<void> {
        const result = this.api.trustPeer(peeringTrustRequest, _options);
        return result.toPromise();
    }


}



import { ObservableRequestsApi } from './ObservableAPI';

import { RequestsApiRequestFactory, RequestsApiResponseProcessor} from "../apis/RequestsApi";
export class PromiseRequestsApi {
    private api: ObservableRequestsApi

    public constructor(
        configuration: Configuration,
        requestFactory?: RequestsApiRequestFactory,
        responseProcessor?: RequestsApiResponseProcessor
    ) {
        this.api = new ObservableRequestsApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * Post an off-ledger request
     * @param offLedgerRequest Offledger request as JSON. Request encoded in Hex
     */
    public offLedgerWithHttpInfo(offLedgerRequest: OffLedgerRequest, _options?: Configuration): Promise<HttpInfo<void>> {
        const result = this.api.offLedgerWithHttpInfo(offLedgerRequest, _options);
        return result.toPromise();
    }

    /**
     * Post an off-ledger request
     * @param offLedgerRequest Offledger request as JSON. Request encoded in Hex
     */
    public offLedger(offLedgerRequest: OffLedgerRequest, _options?: Configuration): Promise<void> {
        const result = this.api.offLedger(offLedgerRequest, _options);
        return result.toPromise();
    }


}



import { ObservableUsersApi } from './ObservableAPI';

import { UsersApiRequestFactory, UsersApiResponseProcessor} from "../apis/UsersApi";
export class PromiseUsersApi {
    private api: ObservableUsersApi

    public constructor(
        configuration: Configuration,
        requestFactory?: UsersApiRequestFactory,
        responseProcessor?: UsersApiResponseProcessor
    ) {
        this.api = new ObservableUsersApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * Add a user
     * @param addUserRequest The user data
     */
    public addUserWithHttpInfo(addUserRequest: AddUserRequest, _options?: Configuration): Promise<HttpInfo<void>> {
        const result = this.api.addUserWithHttpInfo(addUserRequest, _options);
        return result.toPromise();
    }

    /**
     * Add a user
     * @param addUserRequest The user data
     */
    public addUser(addUserRequest: AddUserRequest, _options?: Configuration): Promise<void> {
        const result = this.api.addUser(addUserRequest, _options);
        return result.toPromise();
    }

    /**
     * Change user password
     * @param username The username
     * @param updateUserPasswordRequest The users new password
     */
    public changeUserPasswordWithHttpInfo(username: string, updateUserPasswordRequest: UpdateUserPasswordRequest, _options?: Configuration): Promise<HttpInfo<void>> {
        const result = this.api.changeUserPasswordWithHttpInfo(username, updateUserPasswordRequest, _options);
        return result.toPromise();
    }

    /**
     * Change user password
     * @param username The username
     * @param updateUserPasswordRequest The users new password
     */
    public changeUserPassword(username: string, updateUserPasswordRequest: UpdateUserPasswordRequest, _options?: Configuration): Promise<void> {
        const result = this.api.changeUserPassword(username, updateUserPasswordRequest, _options);
        return result.toPromise();
    }

    /**
     * Change user permissions
     * @param username The username
     * @param updateUserPermissionsRequest The users new permissions
     */
    public changeUserPermissionsWithHttpInfo(username: string, updateUserPermissionsRequest: UpdateUserPermissionsRequest, _options?: Configuration): Promise<HttpInfo<void>> {
        const result = this.api.changeUserPermissionsWithHttpInfo(username, updateUserPermissionsRequest, _options);
        return result.toPromise();
    }

    /**
     * Change user permissions
     * @param username The username
     * @param updateUserPermissionsRequest The users new permissions
     */
    public changeUserPermissions(username: string, updateUserPermissionsRequest: UpdateUserPermissionsRequest, _options?: Configuration): Promise<void> {
        const result = this.api.changeUserPermissions(username, updateUserPermissionsRequest, _options);
        return result.toPromise();
    }

    /**
     * Deletes a user
     * @param username The username
     */
    public deleteUserWithHttpInfo(username: string, _options?: Configuration): Promise<HttpInfo<void>> {
        const result = this.api.deleteUserWithHttpInfo(username, _options);
        return result.toPromise();
    }

    /**
     * Deletes a user
     * @param username The username
     */
    public deleteUser(username: string, _options?: Configuration): Promise<void> {
        const result = this.api.deleteUser(username, _options);
        return result.toPromise();
    }

    /**
     * Get a user
     * @param username The username
     */
    public getUserWithHttpInfo(username: string, _options?: Configuration): Promise<HttpInfo<User>> {
        const result = this.api.getUserWithHttpInfo(username, _options);
        return result.toPromise();
    }

    /**
     * Get a user
     * @param username The username
     */
    public getUser(username: string, _options?: Configuration): Promise<User> {
        const result = this.api.getUser(username, _options);
        return result.toPromise();
    }

    /**
     * Get a list of all users
     */
    public getUsersWithHttpInfo(_options?: Configuration): Promise<HttpInfo<Array<User>>> {
        const result = this.api.getUsersWithHttpInfo(_options);
        return result.toPromise();
    }

    /**
     * Get a list of all users
     */
    public getUsers(_options?: Configuration): Promise<Array<User>> {
        const result = this.api.getUsers(_options);
        return result.toPromise();
    }


}



