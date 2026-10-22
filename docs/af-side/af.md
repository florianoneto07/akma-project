# AF de teste

**MVP:** MVP-3; perfis TLS no MVP-8 · **Data:** 18 de outubro de 2026; perfis em 21/10/2026 · **Verificado** contra o AAnF da VM do free5GC em 18/10/2026 (`teste-af.md`) e, com os perfis TLS, no `e2e.sh` (`e2e.md`)

## 1. O que o AF faz

`[3GPP]` TS 33.535 §6.2.1: o AF recebe do UE o A-KID na requisição de estabelecimento de sessão (passo 1). Se não tem contexto ativo para aquele A-KID, pede a K_AF ao AAnF com `Naanf_AKMA_ApplicationKey_Get`, enviando o A-KID e o próprio AF_ID (passo 2). Com a K_AF, responde ao UE (passo 9).

O AF serve os perfis de TS 24.109 Anexo K, cada um na sua porta (`ua-star-norma.md` §3):

| Perfil | Como o UE se autentica | AF_ID pedido ao AAnF |
|---|---|---|
| `digest-tls` (K.2) | HTTP Digest com a K_AF dentro de TLS 1.3 ou 1.2; o AF se autentica pelo certificado | `<FQDN>.010001<cipher suite>` da conexão |
| `psk-tls13` (K.3.3) | TLS 1.3 com a K_AF como PSK externa, `psk_dhe_ke` | `<FQDN>.0100000007` |
| `psk-tls12` (K.3.2) | TLS 1.2 `ECDHE-PSK` com a K_AF, *hint* `3GPP-AKMA` | `<FQDN>.010001D001` |
| `psk-dtls13` (TS 33.535 Anexo C) | DTLS 1.3 com a K_AF como PSK externa, `psk_dhe_ke`, *cookie* no `HelloRetryRequest`; CoAP `GET` dentro | `<FQDN>.0100000009` |
| `psk-dtls12` | DTLS 1.2 `ECDHE-PSK` com a K_AF, `HelloVerifyRequest`; CoAP dentro | `<FQDN>.010001D001` |
| `digest-http` (teste, desligado por padrão) | o Ua\* simplificado de `ua-star.md`, sem TLS | `<FQDN>.<uaStarId>` |

Nos perfis DTLS, o AF responde em CoAP (RFC 7252): `GET coaps://<FQDN>/ua/v1/session` recebe `2.05 Content` com o mesmo JSON dos perfis TLS (`ua-star-norma.md` §4.2). O AF só pede a K_AF ao AAnF depois do *cookie* (U16).

Os perfis PSK precisam do AF compilado com o wolfSSL (`af/scripts/build-wolfssl.sh` antes de `build-af.sh`); o binário fica sob GPLv3. Sem ele, o AF recusa partir com um perfil PSK na configuração.

Nos perfis Digest:

- **Desafio:** responde `401` com o `realm` `3GPP-bootstrapping-akma@af1.akma.example` a toda requisição sem credenciais. É também a iniciação do AKMA (§6.5).
- **K_AF:** pede-a ao AAnF configurado (§6.7) e a guarda até `expiry`.
- **Verificação do UE:** confere o `response` do Digest com a senha base64(K_AF).
- **Prova da chave:** responde com `rspauth` no `Authentication-Info`, provando ao UE que tem a mesma K_AF.
- **Expiração:** nunca usa uma K_AF vencida. Pela política configurada, pede-a de novo ou recusa com `KAF_EXPIRED` (§6.4.2).

Nos perfis PSK:

- **Identidade:** lê o A-KID depois de `3GPP-AKMA;`, pede a K_AF ao AAnF e a usa como PSK. Sem K_AF, recusa a identidade, e o *handshake* falha.
- **Expiração:** encerra a sessão TLS (`close_notify`) quando a K_AF vence. O *handshake* seguinte segue a política: `refetch` pede a K_AF de novo; `reject` recusa.
- **Retomada:** por *tickets*, válidos até a expiração da K_AF e no máximo 24 horas (TS 33.222 §5.4.1.4). No K.2, os *tickets* valem 24 horas (§5.3.1.5).

`[MVP]` O AF roda na VM do AF (192.168.16.129, ADR-0001-R2). O AAnF o autoriza pelo endereço de origem e pelo AF_ID (DV-1).

## 2. Arquivos

Módulo Go `af/` (`github.com/florianoneto07/akma-free5gc/af`). Depende de `shared/akma` (MVP-0) e do `free5gc/openapi` corrigido (D6), por `replace`; usa as mesmas versões de `free5gc/util` e `logrus` do AAnF.

| Caminho | Conteúdo |
|---|---|
| `af/uastar/` | Digest de RFC 7616 com SHA-256, usado pelo AF e pelo cliente: cabeçalhos `WWW-Authenticate`, `Authorization` e `Authentication-Info`; `response` e `rspauth`; `realm`; *product token* |
| `af/uastar/testdata/gen_digest_vectors.py`, `digest_vectors.json` | implementação Python independente e os vetores dela: o exemplo de RFC 7616 §3.9.1 e o vetor `testbed-ue` com `qop=auth` e `auth-int` |
| `af/af/internal/config/` | `afcfg.yaml` e validação (FQDN como nome DNS, identificador Ua\*, `qop`, política de expiração) |
| `af/af/internal/aanf/` | `Naanf_AKMA_ApplicationKey_Get` pelo cliente do `openapi`; `K_AKMA_NOT_PRESENT` e demais falhas |
| `af/af/internal/server/` | servidor Ua\*: Digest nas verificações de `ua-star.md` §5.4, com o AF_ID do perfil e da conexão; K_AF da identidade PSK; recurso de sessão dos perfis PSK; K_AF por AF_ID e A-KID |
| `af/af/internal/profile/` | um servidor por perfil: certificado do K.2 e `crypto/tls`; PSK pelo `af/wolftls` (só com a *tag* `wolfssl`) |
| `af/akmatls/` | identificadores Ua\* de TS 33.220 H.3, identidade PSK, perfil TLS do K.2 com *tickets* de 24 horas |
| `af/wolftls/` (GPLv3) | ligação cgo ao wolfSSL: TLS e DTLS 1.2 e 1.3 com PSK, *tickets* presos à expiração da chave; no DTLS, um *socket* UDP separado por endereço de origem e *cookies* presos ao endereço |
| `af/coap/` | CoAP (RFC 7252): mensagem, opções, troca confirmável com resposta no ACK e retransmissão |
| `af/af/internal/logger/` | log no formato do free5GC, NF `AF`, categorias `Main`, `CFG`, `UA` e `Consumer` |
| `af/af/cmd/af/` | o programa `af` (`-c` configuração, `-l` arquivo de log) |
| `af/ue/uaclient/` | cliente Ua\* (biblioteca e comando `uaclient`) |
| `af/config/afcfg.yaml` | configuração do laboratório |
| `af/scripts/build-wolfssl.sh` (GPLv3) | baixa o wolfSSL 5.9.4, confere o SHA-256, compila a biblioteca estática e confere as macros |
| `af/scripts/make-test-ca.sh` | CA de teste e certificado do AF em `af/build/cert` |
| `af/scripts/build-af.sh` | compila `af/build/bin/af` e `af/build/bin/uaclient`, com o wolfSSL quando a biblioteca existe |

## 3. Configuração

Os campos e os padrões estão em `ua-star.md` §11; os perfis, em `ua-star-norma.md` §3.5. No laboratório:

| Campo | Valor |
|---|---|
| `fqdn` | `af1.akma.example` |
| `profiles` | `digest-tls` em `192.168.16.129:8443`, com `certificate` e `key` de `make-test-ca.sh`; `psk-tls13` em `:8444`; `psk-tls12` em `:8445`; `psk-dtls13` em UDP `:5684`; `psk-dtls12` em UDP `:5685`. O `digest-http` em `:8080` fica comentado: o `e2e.sh` o acrescenta |
| `uaStarId` | `FF00000001`, só do `digest-http` |
| `kafExpiry` | `refetch` |
| `aanfUri` | `http://192.168.16.130:8010` |

Sem `profiles`, o `listen` de antes vale como um perfil `digest-http`. Com `profiles`, `listen` não pode aparecer.

Cada AF_ID dos perfis precisa estar em `akma.authorizedAfs` do `aanfcfg.yaml`, com o endereço do AF (`docs/free5gc-side/aanf.md` §2): `af1.akma.example.` mais `0100011301`, `0100011302`, `0100011303`, `010001C02B` e `010001C02C` (K.2 com certificado ECDSA), `0100000007`, `010001D001`, `0100000009` e `FF00000001`. O mesmo vale para a política do NEF e para a do UE.

**AF fora da rede do operador (MVP-6).** Com `nefUri` no lugar de `aanfUri`, o AF pede a K_AF ao NEF (`POST {nefUri}/3gpp-akma/v1/retrieve`, TS 29.522 §4.4.23.2; TS 33.535 §6.3). Especificação: `docs/free5gc-side/nef.md`.

| Campo | Uso |
|---|---|
| `nefUri` | API do NEF, `http://192.168.16.130:8000`. Exclusivo com `aanfUri` |
| `nrfUri` | NRF em que o AF se registra como `nfType: AF` e de que obtém o *token* `3gpp-akma` (`nef.md` §2.1). Obrigatório com `nefUri`. No laboratório, `http://127.0.0.10:8000` pelo túnel SSH |
| `nfInstanceId` | UUID do AF no NRF; sem ele, um UUID novo a cada partida |
| `anonymous` | `true` manda `anonInd: true`, direto ao AAnF ou via NEF: a K_AF vem sem o SUPI (TS 33.535 §6.2.2). Vale nos dois modos |

No modo externo, o AF se registra no NRF na partida e se descadastra ao encerrar. A resposta do NEF nunca traz o SUPI (§6.3 passo 5).

`[MVP]` A opção `clockOffset` (duração Go, como `2h` ou `-2h`) desloca o relógio do AF. É só de teste: simula um AF cujo relógio discorda do relógio do AAnF (MVP-5, `e2e.md` §4). Com ela, o AF registra um `warn` na partida.

## 4. Log

Mesmo formato das NFs do free5GC (MVP-3a), na categoria `UA` para o Ua\* e `Consumer` para o pedido ao AAnF:

| Nível | Mensagem |
|---|---|
| `info` | `FQDN …; realm "…"; qop […]; nonce lifetime …; K_AF expiry policy …; anonymous …`, `Built with wolfSSL: PSK profiles available`, `AF inside the operator network: K_AF from the AAnF at …` ou `AF outside the operator network: K_AF through the NEF at …; registered at the NRF … as AF … (OAuth2 …)`, e `Ua* server <perfil> listening on <endereço> for <URL>, AF_ID …` (partida) |
| `info` | `AKMA initiation: challenge sent to <endereço> (product token 3gpp-akma: true\|false)` |
| `info` | `K_AF for A-KID <A-KID> and AF_ID <AF_ID> obtained from the AAnF, expires <instante>`; no modo externo, `obtained from the NEF` |
| `info` | `UE with A-KID <A-KID> authenticated (AF_ID <AF_ID>, nc <nc>, qop <qop>)` |
| `info` | `TLS 1.3 PSK handshake with <endereço>: A-KID <A-KID>, AF_ID <AF_ID>, cipher suite 1301, resumed true\|false`; `TLS 1.2 …` |
| `info` | `UE with A-KID <A-KID> served over PSK TLS (AF_ID <AF_ID>)`; nos perfis DTLS, `served over PSK DTLS`; `DTLS 1.3 PSK handshake with …` |
| `info` | `K_AF of A-KID <A-KID>, AF_ID <AF_ID> expired: TLS session with <endereço> closed` |
| `info` | `No K_AKMA at the AAnF for A-KID <A-KID>: access rejected`; nos perfis PSK, `… (AF_ID <AF_ID>): handshake refused` |
| `warn` | `TLS 1.3 PSK handshake with <endereço> failed: <erro do wolfSSL>: <causa>`; `PSK identity from <endereço> is not 3GPP-AKMA;<A-KID>: handshake refused`; erros de *handshake* do K.2 |
| `info` | `K_AF for A-KID <A-KID> expired at <instante>: requesting it again from the AAnF (policy refetch)` ou `…: access rejected (policy reject)` |
| `warn` | `Digest response from <endereço> does not match the K_AF of A-KID <A-KID>: authentication failed` (RFC 7616 §3.4: registrar falhas de *login*) |
| `warn` | `Request from <endereço> for host "<Host>", not <FQDN>`; `Realm … is not …`; `Nonce count … already used …: possible replay` |
| `error` | `K_AF for A-KID <A-KID> not obtained: <motivo>` |
| `debug` | o SUPI recebido do AAnF |

A K_AF fica em `kdf.Key`, que se imprime como `[256-bit key]`; nem ela nem a senha do Digest aparecem em log.

## 5. Como reproduzir

Na máquina do AF, com Go 1.26.2 e o repositório clonado:

```bash
af/scripts/build-wolfssl.sh
```

```bash
af/scripts/make-test-ca.sh
```

```bash
af/scripts/build-af.sh
```

```bash
af/build/bin/af -c af/config/afcfg.yaml
```

O AAnF tem de estar no ar, com o AF autorizado na política. Para testar sem UE, semeie um contexto no AAnF com o `seed` e use o `uaclient` com o mesmo vetor (`teste-af.md`).

## 6. Como verificar

```bash
cd af && go test -race ./...
```

```bash
cd af && go test -race -tags wolfssl ./...
```

O teste de integração, com o AAnF real, e o resultado de cada critério estão em `teste-af.md`.
