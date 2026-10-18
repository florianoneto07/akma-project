// Command uaclient establishes an application session with the AF over the
// Ua* protocol of the MVP (docs/af-side/ua-star.md). In MVP-3 the A-KID and
// K_AF are given on the command line or taken from a shared/akma vector; in
// MVP-4 they come from the ME of UERANSIM. K_AF is never printed.
package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/florianoneto07/akma-free5gc/af/ue/uaclient"
	"github.com/florianoneto07/akma-free5gc/shared/akma/codec"
	"github.com/florianoneto07/akma-free5gc/shared/akma/kdf"
)

func main() {
	var (
		fqdn      = flag.String("fqdn", "af1.akma.example", "FQDN of the AF (sent in Host, checked in the realm)")
		connect   = flag.String("connect", "192.168.16.129:8080", "address of the AF")
		path      = flag.String("path", "/ua/v1/session", "protected resource")
		aKid      = flag.String("akid", "", "A-KID")
		kafHex    = flag.String("kaf", "", "K_AF, 64 hexadecimal characters")
		kafB64    = flag.String("kaf-base64", "", "K_AF in base64")
		vector    = flag.String("vector", "", "take A-KID and K_AF from this vector of -vectors")
		vectors   = flag.String("vectors", "shared/akma/testdata/vectors.json", "vectors file of shared/akma")
		qop       = flag.String("qop", "auth", "qop: auth or auth-int")
		count     = flag.Int("n", 1, "number of requests in the session")
		userAgent = flag.String("user-agent", uaclient.DefaultUserAgent, "User-Agent")
	)
	flag.Parse()
	if err := run(*fqdn, *connect, *path, *aKid, *kafHex, *kafB64, *vector, *vectors, *qop, *userAgent, *count); err != nil {
		msg := err.Error()
		if !strings.HasPrefix(msg, "uaclient:") {
			msg = "uaclient: " + msg
		}
		fmt.Fprintln(os.Stderr, msg)
		os.Exit(1)
	}
}

func run(fqdn, connect, path, aKid, kafHex, kafB64, vector, vectors, qop, userAgent string, count int) error {
	var kaf kdf.Key
	var err error
	switch {
	case vector != "":
		aKid, kaf, err = fromVector(vectors, vector)
	case kafHex != "":
		kaf, err = codec.ParseHexKey(kafHex)
	case kafB64 != "":
		kaf, err = parseBase64Key(kafB64)
	default:
		err = errors.New("give -vector, -kaf or -kaf-base64")
	}
	if err != nil {
		return err
	}
	if aKid == "" {
		return errors.New("-akid is required")
	}
	s, err := uaclient.NewSession(uaclient.Config{
		FQDN: fqdn, Connect: connect, Path: path, AKID: aKid, KAf: kaf, Qop: qop, UserAgent: userAgent,
	})
	if err != nil {
		return err
	}
	fmt.Printf("AF %s at %s, A-KID %s\n", fqdn, connect, aKid)
	for i := 0; i < count; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		r, err := s.Establish(ctx)
		cancel()
		if err != nil {
			return err
		}
		fmt.Printf("session established (nc %s): AF authenticated by rspauth; K_AF valid until %s\n",
			r.NC, r.KafExpiry.Format(time.RFC3339))
	}
	return nil
}

func fromVector(path, name string) (string, kdf.Key, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", kdf.Key{}, err
	}
	var doc struct {
		Vectors []struct {
			Name string `json:"name"`
			AKid struct {
				Base64 string `json:"base64"`
			} `json:"a_kid"`
			KAf string `json:"k_af"`
		} `json:"vectors"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return "", kdf.Key{}, err
	}
	for _, v := range doc.Vectors {
		if v.Name == name {
			k, err := codec.ParseHexKey(v.KAf)
			return v.AKid.Base64, k, err
		}
	}
	return "", kdf.Key{}, fmt.Errorf("vector %q not in %s", name, path)
}

func parseBase64Key(s string) (kdf.Key, error) {
	var k kdf.Key
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil || len(b) != len(k) {
		return k, errors.New("-kaf-base64: not 32 octets in base64")
	}
	copy(k[:], b)
	return k, nil
}
