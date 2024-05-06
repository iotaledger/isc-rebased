package types

type SuiX_GetLatestSuiSystemState struct {
	Jsonrpc string `json:"jsonrpc,omitempty"`
	Result  struct {
		Epoch                                 string `json:"epoch,omitempty"`
		ProtocolVersion                       string `json:"protocolVersion,omitempty"`
		SystemStateVersion                    string `json:"systemStateVersion,omitempty"`
		StorageFundTotalObjectStorageRebates  string `json:"storageFundTotalObjectStorageRebates,omitempty"`
		StorageFundNonRefundableBalance       string `json:"storageFundNonRefundableBalance,omitempty"`
		ReferenceGasPrice                     string `json:"referenceGasPrice,omitempty"`
		SafeMode                              bool   `json:"safeMode,omitempty"`
		SafeModeStorageRewards                string `json:"safeModeStorageRewards,omitempty"`
		SafeModeComputationRewards            string `json:"safeModeComputationRewards,omitempty"`
		SafeModeStorageRebates                string `json:"safeModeStorageRebates,omitempty"`
		SafeModeNonRefundableStorageFee       string `json:"safeModeNonRefundableStorageFee,omitempty"`
		EpochStartTimestampMs                 string `json:"epochStartTimestampMs,omitempty"`
		EpochDurationMs                       string `json:"epochDurationMs,omitempty"`
		StakeSubsidyStartEpoch                string `json:"stakeSubsidyStartEpoch,omitempty"`
		MaxValidatorCount                     string `json:"maxValidatorCount,omitempty"`
		MinValidatorJoiningStake              string `json:"minValidatorJoiningStake,omitempty"`
		ValidatorLowStakeThreshold            string `json:"validatorLowStakeThreshold,omitempty"`
		ValidatorVeryLowStakeThreshold        string `json:"validatorVeryLowStakeThreshold,omitempty"`
		ValidatorLowStakeGracePeriod          string `json:"validatorLowStakeGracePeriod,omitempty"`
		StakeSubsidyBalance                   string `json:"stakeSubsidyBalance,omitempty"`
		StakeSubsidyDistributionCounter       string `json:"stakeSubsidyDistributionCounter,omitempty"`
		StakeSubsidyCurrentDistributionAmount string `json:"stakeSubsidyCurrentDistributionAmount,omitempty"`
		StakeSubsidyPeriodLength              string `json:"stakeSubsidyPeriodLength,omitempty"`
		StakeSubsidyDecreaseRate              int    `json:"stakeSubsidyDecreaseRate,omitempty"`
		TotalStake                            string `json:"totalStake,omitempty"`
		ActiveValidators                      []struct {
			SuiAddress                   string `json:"suiAddress,omitempty"`
			ProtocolPubkeyBytes          string `json:"protocolPubkeyBytes,omitempty"`
			NetworkPubkeyBytes           string `json:"networkPubkeyBytes,omitempty"`
			WorkerPubkeyBytes            string `json:"workerPubkeyBytes,omitempty"`
			ProofOfPossessionBytes       string `json:"proofOfPossessionBytes,omitempty"`
			Name                         string `json:"name,omitempty"`
			Description                  string `json:"description,omitempty"`
			ImageURL                     string `json:"imageUrl,omitempty"`
			ProjectURL                   string `json:"projectUrl,omitempty"`
			NetAddress                   string `json:"netAddress,omitempty"`
			P2PAddress                   string `json:"p2pAddress,omitempty"`
			PrimaryAddress               string `json:"primaryAddress,omitempty"`
			WorkerAddress                string `json:"workerAddress,omitempty"`
			NextEpochProtocolPubkeyBytes any    `json:"nextEpochProtocolPubkeyBytes,omitempty"`
			NextEpochProofOfPossession   any    `json:"nextEpochProofOfPossession,omitempty"`
			NextEpochNetworkPubkeyBytes  any    `json:"nextEpochNetworkPubkeyBytes,omitempty"`
			NextEpochWorkerPubkeyBytes   any    `json:"nextEpochWorkerPubkeyBytes,omitempty"`
			NextEpochNetAddress          any    `json:"nextEpochNetAddress,omitempty"`
			NextEpochP2PAddress          any    `json:"nextEpochP2pAddress,omitempty"`
			NextEpochPrimaryAddress      any    `json:"nextEpochPrimaryAddress,omitempty"`
			NextEpochWorkerAddress       any    `json:"nextEpochWorkerAddress,omitempty"`
			VotingPower                  string `json:"votingPower,omitempty"`
			OperationCapID               string `json:"operationCapId,omitempty"`
			GasPrice                     string `json:"gasPrice,omitempty"`
			CommissionRate               string `json:"commissionRate,omitempty"`
			NextEpochStake               string `json:"nextEpochStake,omitempty"`
			NextEpochGasPrice            string `json:"nextEpochGasPrice,omitempty"`
			NextEpochCommissionRate      string `json:"nextEpochCommissionRate,omitempty"`
			StakingPoolID                string `json:"stakingPoolId,omitempty"`
			StakingPoolActivationEpoch   string `json:"stakingPoolActivationEpoch,omitempty"`
			StakingPoolDeactivationEpoch any    `json:"stakingPoolDeactivationEpoch,omitempty"`
			StakingPoolSuiBalance        string `json:"stakingPoolSuiBalance,omitempty"`
			RewardsPool                  string `json:"rewardsPool,omitempty"`
			PoolTokenBalance             string `json:"poolTokenBalance,omitempty"`
			PendingStake                 string `json:"pendingStake,omitempty"`
			PendingTotalSuiWithdraw      string `json:"pendingTotalSuiWithdraw,omitempty"`
			PendingPoolTokenWithdraw     string `json:"pendingPoolTokenWithdraw,omitempty"`
			ExchangeRatesID              string `json:"exchangeRatesId,omitempty"`
			ExchangeRatesSize            string `json:"exchangeRatesSize,omitempty"`
		} `json:"activeValidators,omitempty"`
		PendingActiveValidatorsID   string `json:"pendingActiveValidatorsId,omitempty"`
		PendingActiveValidatorsSize string `json:"pendingActiveValidatorsSize,omitempty"`
		PendingRemovals             []any  `json:"pendingRemovals,omitempty"`
		StakingPoolMappingsID       string `json:"stakingPoolMappingsId,omitempty"`
		StakingPoolMappingsSize     string `json:"stakingPoolMappingsSize,omitempty"`
		InactivePoolsID             string `json:"inactivePoolsId,omitempty"`
		InactivePoolsSize           string `json:"inactivePoolsSize,omitempty"`
		ValidatorCandidatesID       string `json:"validatorCandidatesId,omitempty"`
		ValidatorCandidatesSize     string `json:"validatorCandidatesSize,omitempty"`
		AtRiskValidators            []any  `json:"atRiskValidators,omitempty"`
		ValidatorReportRecords      []any  `json:"validatorReportRecords,omitempty"`
	} `json:"result,omitempty"`
	ID int `json:"id,omitempty"`
}
