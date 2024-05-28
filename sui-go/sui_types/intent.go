package sui_types

// import (
// 	"github.com/iotaledger/wasp/sui-go/lib"
// )

// type IntentScope struct {
// 	TransactionData         *serialization.EmptyEnum // Used for a user signature on a transaction data.
// 	TransactionEffects      *serialization.EmptyEnum // Used for an authority signature on transaction effects.
// 	CheckpointSummary       *serialization.EmptyEnum // Used for an authority signature on a checkpoint summary.
// 	PersonalMessage         *serialization.EmptyEnum // Used for a user signature on a personal message.
// 	SenderSignedTransaction *serialization.EmptyEnum // Used for an authority signature on a user signed transaction.
// 	ProofOfPossession       *serialization.EmptyEnum // Used as a signature representing an authority's proof of possession of its authority protocol key.
// 	HeaderDigest            *serialization.EmptyEnum // Used for narwhal authority signature on header digest.
// }

// func (i IntentScope) IsBcsEnum() {}

// type IntentVersion struct {
// 	V0 *serialization.EmptyEnum
// }

// func (i IntentVersion) IsBcsEnum() {}

// type AppId struct {
// 	Sui     *serialization.EmptyEnum
// 	Narwhal *serialization.EmptyEnum
// }

// func (a AppId) IsBcsEnum() {}

// type Intent struct {
// 	Scope   IntentScope
// 	Version IntentVersion
// 	AppId   AppId
// }

// func DefaultIntent() Intent {
// 	return Intent{
// 		Scope: IntentScope{
// 			TransactionData: &serialization.EmptyEnum{},
// 		},
// 		Version: IntentVersion{
// 			V0: &serialization.EmptyEnum{},
// 		},
// 		AppId: AppId{
// 			Sui: &serialization.EmptyEnum{},
// 		},
// 	}
// }

// type IntentValue interface {
// 	TransactionData | ~[]byte
// }

// type IntentMessage[T IntentValue] struct {
// 	Intent Intent
// 	Value  T
// }

// func MessageWithIntent[T IntentValue](intent Intent, value T) IntentMessage[T] {
// 	return IntentMessage[T]{
// 		Intent: intent,
// 		Value:  value,
// 	}
// }
