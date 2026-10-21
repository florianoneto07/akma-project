# Ua\* profiles of the test AF

The test AF and the Ua\* client implement the TLS profiles of TS 24.109 Annex K, which TS 33.535 Annex B adapts to AKMA from the GBA profiles of TS 33.222. Each profile runs on its own port of the AF.

| Profile | Specification | How the UE authenticates | Ua\* identifier (TS 33.220 H.3) |
|---|---|---|---|
| `digest-tls` | TS 24.109 K.2; TS 33.222 cl. 5.3 | HTTP Digest keyed with K_AF inside TLS 1.3 or 1.2; the AF authenticates with its certificate | `01 00 01 yy zz`, `yy zz` the cipher suite of the connection |
| `psk-tls13` | TS 24.109 K.3.3; TS 33.222 cl. 5.4.0.2 | TLS 1.3 with K_AF as external pre-shared key | `01 00 00 00 07` |
| `psk-tls12` | TS 24.109 K.3.2; TS 33.222 cl. 5.4.0.1 | TLS 1.2 `ECDHE_PSK` with K_AF | `01 00 01 D0 01` |
| `digest-http` | none: test only | the simplified Ua\* of [ua-star.md](ua-star.md), without TLS | `FF 00 00 00 01` (private range) |

The AF_ID is the FQDN of the AF followed by the identifier, `<FQDN>.<10 hex digits>` in the APIs. Each profile, and in K.2 and TLS 1.2 each cipher suite, has its own AF_ID and so its own K_AF. The AAnF, the NEF and the key service of the UE must authorize every AF_ID that the AF can use.

## 1. K.2: Digest inside TLS

1. The UE opens TLS with the AF, checks that the certificate of the AF matches the FQDN of the AF (RFC 2818 cl. 3.1) and sends no certificate.
2. The UE sends the request with the product token `3gpp-akma` in `User-Agent`. The AF answers `401` with a Digest challenge whose realm is `3GPP-bootstrapping-akma@<FQDN>`.
3. The UE checks that the realm names the AF it opened TLS with. Only then does it ask the ME for K_AF, for the AF_ID with the cipher suite of the connection, and send the A-KID as `username` with base64(K_AF) as password.
4. The AF asks the AAnF for K_AF with the same AF_ID, checks the Digest and proves K_AF back with `rspauth`.

The Digest is the one of [ua-star.md](ua-star.md) (RFC 7616 with SHA-256). With the ECDSA certificate of the example, the AF_IDs are `0100011301`, `0100011302` and `0100011303` in TLS 1.3 and `010001C02B` and `010001C02C` in TLS 1.2.

## 2. PSK profiles

- **TLS 1.3.** The UE asks the ME for the K_AF of `<FQDN>.0100000007` before the `ClientHello`, and offers the identity `3GPP-AKMA;<A-KID>` with `psk_key_exchange_modes` set to `psk_dhe_ke` only: TS 33.210 cl. 6.2.2 forbids `psk_ke`. The AF reads the A-KID from the identity, asks the AAnF for the K_AF of the same AF_ID and uses it as the external PSK, with SHA-256 (RFC 8446 cl. 4.2.11) and `TLS_AES_128_GCM_SHA256`. No early data.
- **TLS 1.2.** The AF selects `TLS_ECDHE_PSK_WITH_AES_128_GCM_SHA256` (RFC 8442, mandatory in TS 33.210 cl. 6.2.3) and sends the hint `3GPP-AKMA`. The UE then asks the ME for the K_AF of `<FQDN>.010001D001`, the identifier of that cipher suite, and sends the identity `3GPP-AKMA;<A-KID>`. K_AF is the PSK of RFC 5489. Extended Master Secret and `renegotiation_info` are on; renegotiation itself is never accepted.

K_AF, the full 32 octets, is the PSK: the specifications only say that the PSK is derived from K_AF. If the AF cannot get K_AF, it refuses the identity and the handshake fails: TLS has no way to carry a cause. Once the handshake is done, the UE sends `GET /ua/v1/session` without HTTP authentication and gets the A-KID and the expiry of K_AF.

The AF requires the `server_name` of the client to be its FQDN. In both versions the client also offers other authentication methods, as TS 33.222 cl. 5.4.0 asks, but it trusts no certificate in these profiles: if the AF picked one, the handshake would fail. K.2 is the profile with a certificate.

## 3. Session resumption and K_AF expiry

- **PSK profiles.** Resumption by session tickets only, in TLS 1.2 (RFC 5077) and in TLS 1.3. A ticket is valid until the earlier of the expiry of the K_AF of its full handshake and 24 hours (TS 33.222 cl. 5.4.1.4); the AF keeps the A-KID of each ticket. The AF closes a TLS session (`close_notify`) when its K_AF expires; the next handshake follows `kafExpiry`: `refetch` asks the AAnF again, `reject` refuses.
- **K.2.** Tickets are accepted for 24 hours from their issue (TS 33.222 cl. 5.3.1.5). K_AF does not enter TLS here: each Digest request is checked against K_AF and its expiry, as in [ua-star.md](ua-star.md).

## 4. Configuration

In `afcfg.yaml`, one entry per profile:

```yaml
configuration:
  fqdn: af1.akma.example
  profiles:
    - type: digest-tls
      listen: 192.168.16.129:8443
      certificate: ../build/cert/af1.akma.example.pem   # relative to this file
      key: ../build/cert/af1.akma.example.key
    - type: psk-tls13
      listen: 192.168.16.129:8444
    - type: psk-tls12
      listen: 192.168.16.129:8445
    # - type: digest-http          # test only, off by default
    #   listen: 192.168.16.129:8080
  uaStarId: FF00000001             # identifier of digest-http
```

Without `profiles`, a top-level `listen` is a single `digest-http` server, as in earlier versions.

The client picks the profile with `-profile`. Without a port in `-connect`, it uses the one of the example: 8443, 8444, 8445 or 8080.

| Option | |
|---|---|
| `-profile` | `digest-tls`, `psk-tls13`, `psk-tls12` or `digest-http` (default) |
| `-ca` | PEM certificate of the CA of the AF, required with `digest-tls` |
| `-tls12` | `digest-tls` in TLS 1.2 only |
| `-reconnect` | a new TLS connection for each request, resuming the previous session |
| `-psk-ke-only` | `psk-tls13` offering `psk_ke` only, which the AF refuses: a negative test |

In TLS 1.3, wolfSSL tells the client that it offered a ticket, not whether the AF took it; the log of the AF says whether the session was resumed.

## 5. Build and licenses

Go's `crypto/tls` does K.2. It has neither the PSK cipher suites of TS 33.210 nor external PSKs in TLS 1.3, so the PSK profiles use wolfSSL through `af/wolftls`, a cgo binding of this repository:

```bash
af/scripts/build-wolfssl.sh
```

```bash
af/scripts/make-test-ca.sh
```

```bash
af/scripts/build-af.sh
```

`build-wolfssl.sh` downloads wolfSSL 5.9.4, checks the SHA-256 of the archive, builds a static library in `af/build/wolfssl` and checks its options: TLS 1.2 and 1.3, DTLS 1.3, PSK, AES-GCM and AES-CCM, SNI, Extended Master Secret, session tickets; no finite field DH, no TLS below 1.2, no early data. `build-af.sh` builds static binaries with the build tag `wolfssl` when the library is there; `AF_WOLFSSL=no` builds without it, and the PSK profiles are then refused.

wolfSSL is licensed under the GNU General Public License v3.0. `af/wolftls` and `af/scripts/build-wolfssl.sh` are under the same license, and so are the AF and client binaries built with the tag. The rest of `af/` is under the Apache License 2.0.

## 6. Limitations

- `TLS_ECDHE_PSK_WITH_AES_256_GCM_SHA384`, recommended by TS 33.210 cl. 6.2.3, is not in wolfSSL and is not offered.
- The certificate status request extension (OCSP stapling), recommended by TS 33.210 cl. 6.2, is not supported; the test CA has no OCSP responder.
- HTTP/1.1 only.
- One profile per port: a port that negotiates K.2 or K.3 in the same `ClientHello` (TS 33.222 cl. 5.4.0.1 NOTE 3) is not supported.
- The tickets of an AF are sealed with a key of the process: they do not survive a restart.
- The DTLS and OSCORE profiles of TS 33.535 Annexes C and D are not implemented.
