package afid

import (
	"encoding/hex"
	"strings"
)

// UaStarID is the five-octet Ua security protocol identifier
// (TS 33.220 Annex H.1).
type UaStarID [5]byte

var (
	// UaStarMVP identifies the project's simplified Ua* protocol. It lies in
	// the private range (first octet 0xFF, TS 33.220 H.2 NOTE 2; ADR-0003).
	UaStarMVP = UaStarID{0xFF, 0x00, 0x00, 0x00, 0x01}
	// UaStarTLS13 identifies shared-key TLS 1.3 (TS 33.220 H.3).
	UaStarTLS13 = UaStarID{0x01, 0x00, 0x00, 0x00, 0x07}
)

// String returns the identifier as ten uppercase hexadecimal characters, the
// form used inside AfId. TS 29.522 describes the identifier as "a string of 5
// octet"; the hexadecimal rendering follows the specification's example, and
// writing it in uppercase is a project decision (plan, P9).
func (u UaStarID) String() string {
	return strings.ToUpper(hex.EncodeToString(u[:]))
}
