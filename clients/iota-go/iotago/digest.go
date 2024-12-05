package iotago

type Digest = Base58
type ObjectDigest = Digest
type TransactionDigest = Digest
type TransactionEffectsDigest = Digest
type TransactionEventsDigest = Digest
type CheckpointDigest = Digest
type CertificateDigest = Digest
type CheckpointContentsDigest = Digest

func NewDigest(str string) (*Digest, error) {
	return NewBase58(str)
}

func MustNewDigest(str string) *Digest {
	digest, err := NewBase58(str)
	if err != nil {
		panic(err)
	}
	return digest
}

func (d Digest) Bytes() []byte {
	return d.Data()
}

func DigestFromBytes(b []byte) *Digest {
	ret := Digest(b)
	return &ret
}
