# Teste do AF e do Ua\* (MVP-3)

**Data:** 18 de outubro de 2026 · AF e cliente na VM do AF (192.168.16.129); AAnF e NRF na VM do free5GC (192.168.16.130)

## 1. Testes automáticos

```bash
cd af && go test -race -count=1 ./...
```

| Pacote | Testes | O que confirmam |
|---|---|---|
| `uastar` | `TestDigestVectors`, `TestRFC7616Example` | `response` e `rspauth` iguais aos do Python e ao exemplo de RFC 7616 §3.9.1 |
| `uastar` | `TestChallengeRoundTrip`, `TestCredentialsRoundTrip`, `TestParseCredentialsRejects`, `TestQuotedStringEscapes`, `TestAuthInfoRoundTrip`, `TestSplitRealm`, `TestHasProductToken`, `TestNC` | sintaxe dos cabeçalhos (aspas conforme RFC 7616 §3.3 a §3.5), `username*` e `userhash` recusados |
| `af/internal/aanf` | `TestApplicationKey`, `TestApplicationKeyErrors`, `TestApplicationKeyUnreachable` | `afId` na forma de transporte e `aKId` sem alteração; `K_AKMA_NOT_PRESENT` separado das demais falhas, com um AAnF falso em h2c |
| `af/internal/config` | `TestDefaults`, `TestInvalid`, `TestSampleConfig`, `TestUnknownField` | validação e o `afcfg.yaml` do laboratório |
| `af/internal/server` | `TestInitiation` | desafio a pedido sem credenciais, com e sem *product token*; o AAnF não é chamado |
| `af/internal/server` | `TestSessionEstablished`, `TestAuthInt`, `TestWrongKafRejected` | sessão com a K_AF certa, três requisições e um único pedido ao AAnF; `auth-int`; K_AF errada recusada |
| `af/internal/server` | `TestKeyNotAvailable`, `TestAanfUnavailable` | `403 AKMA_KEY_NOT_AVAILABLE` e `503 AANF_UNAVAILABLE`, sem desafio |
| `af/internal/server` | `TestRequestChecks`, `TestReplay`, `TestHostPathMethod`, `TestStaleNonce` | as verificações de `ua-star.md` §5.4: `uri`, MD5, `qop`, A-KID, `realm`, nonce, opaque, `nc`, `Host`, caminho, método, `stale=true` |
| `af/internal/server` | `TestKafExpiryRefetch`, `TestKafExpiryReject`, `TestExpiredOnArrival` | K_AF nunca usada a partir do instante de `expiry`, nas duas políticas, com relógio controlado |
| `af/internal/server` | `TestNoKeyInLogs` | K_AF em hexadecimal e em base64 ausentes do log em nível `trace` |
| `ue/uaclient` | `TestHonestAF`, `TestAFWithoutKaf`, `TestChallengeChecks`, `TestHostAndProductToken` | o cliente recusa AF sem `rspauth` válido e desafio de outro AF ou sem AKMA, sem chegar a enviar o A-KID |

Resultado em 18/10/2026: todos passam, também com `-race`.

## 2. Integração com o AAnF

Pré-requisitos: AAnF e NRF na VM do free5GC e o contexto do vetor `testbed-ue` semeado, como no MVP-1 (`docs/free5gc-side/aanf.md` §9):

```bash
~/akma-tools/seed register --nrf http://127.0.0.10:8000 --supi imsi-208930000000001 --k-ausf 141ef548f86e3a2a1ad949331caa9e6bb12401afa321bde4f9f23b84870dd654 --rid 0000 --mcc 208 --mnc 93
```

Na VM do AF, o AF numa aba (`af/build/bin/af -c af/config/afcfg.yaml`) e o cliente em outra:

```bash
af/build/bin/uaclient -vector testbed-ue -n 3
```

O cliente lê o A-KID e a K_AF do vetor `testbed-ue` de `shared/akma`, a K_AF que o oráculo deriva para `af1.akma.example.FF00000001`. Só autentica se a K_AF que o AF obteve do AAnF for a mesma.

| # | Caso | Comando | Resultado em 18/10/2026 |
|---|---|---|---|
| 1 | pedido sem parâmetros AKMA | `GET /ua/v1/session` sem `Authorization`, `User-Agent: curl/8.0` | `401`, `WWW-Authenticate: Digest realm="3GPP-bootstrapping-akma@af1.akma.example", qop="auth", algorithm=SHA-256, nonce="…", opaque="…", charset=UTF-8` |
| 2 | K_AF certa, três requisições | `uaclient -vector testbed-ue -n 3` | três sessões, `nc` 1 a 3, cada uma com `rspauth` válido; K_AF até `2026-10-18T17:28:25Z`. O AF pediu a K_AF ao AAnF uma vez |
| 3 | K_AF errada | `uaclient -akid <A-KID do vetor> -kaf 00…01` | `401`; o AF registrou `Digest response … does not match` |
| 4 | A-KID sem K_AKMA no AAnF | `uaclient -vector rid-012` | `403 AKMA_KEY_NOT_AVAILABLE`; o AAnF registrou `K_AKMA not present for A-KID rid012.atidCoZy…` |
| 5 | `Host` de outro AF | `uaclient -vector testbed-ue -fqdn af2.akma.example` | `421 MISDIRECTED_REQUEST`, sem pedido ao AAnF |
| 6 | AF_ID fora da política do AAnF | segundo AF com `fqdn: af2.akma.example`, na porta 8081 | `503 AANF_UNAVAILABLE`; o AAnF registrou `AF at 192.168.16.129 not authorized for AF_ID af2.akma.example.FF00000001` (403) |

No log do AAnF:

```
[INFO][AANF][AKMA] K_AF of imsi-208930000000001 provided to AF_ID af1.akma.example.FF00000001 for A-KID rid0000.atidiz1FN8DkMGyxwtfTVYsKUsWnzOJdBOLlQEiQlK2gHNA=@5gc.mnc093.mcc208.3gppnetwork.org, expires 2026-10-18T17:28:25Z
```

Os relógios das duas VMs estavam de acordo (13:28 −03:00 e 09:28 −07:00).

Depois do teste: `seed remove`, AF, AAnF e NRF parados.

## 3. Critérios de aceitação do MVP-3

| Critério | Verificação | Resultado |
|---|---|---|
| 1. K_AF do AF igual à do oráculo | integração, caso 2: o `response` do cliente, feito com a K_AF do oráculo, confere com a K_AF que o AF recebeu do AAnF | atendido |
| 2. `afId` na forma de transporte e `aKId` sem alteração | `TestApplicationKey`; linha do AAnF acima | atendido |
| 3. Ua\* aceito com a K_AF certa e recusado com qualquer outra | casos 2 e 3; `TestWrongKafRejected` | atendido |
| 4. Confirmação de chave nos dois sentidos | `rspauth` conferido pelo cliente no caso 2; `TestAFWithoutKaf` | atendido |
| 5. K_AF nunca usada após `expiry` | `TestKafExpiryRefetch`, `TestKafExpiryReject`, `TestExpiredOnArrival`, com relógio controlado. Na integração, o AAnF usa tempo de vida de 1 h, e o teste não esperou a expiração | atendido nos testes automáticos |
| 6. Pedido sem parâmetros AKMA recebe a iniciação | caso 1; `TestInitiation` | atendido |
| 7. Nenhuma chave em log | `TestNoKeyInLogs`; nenhuma cadeia de 64 hexadecimais em `aanf.log` e `aanf.out`, nem a K_AF do vetor, em hexadecimal ou base64; nenhuma na tela do AF | atendido |
