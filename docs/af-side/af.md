# AF de teste

**MVP:** MVP-3 · **Data:** 18 de outubro de 2026 · **Verificado** contra o AAnF da VM do free5GC em 18/10/2026 (`teste-af.md`)

## 1. O que o AF faz

`[3GPP]` TS 33.535 §6.2.1: o AF recebe do UE o A-KID na requisição de estabelecimento de sessão (passo 1). Se não tem contexto ativo para aquele A-KID, pede a K_AF ao AAnF com `Naanf_AKMA_ApplicationKey_Get`, enviando o A-KID e o próprio AF_ID (passo 2). Com a K_AF, responde ao UE (passo 9).

O AF do MVP faz isso pelo Ua\* simplificado de `ua-star.md`, que é HTTP Digest com os parâmetros AKMA, sem TLS:

- **Desafio:** responde `401` com o `realm` `3GPP-bootstrapping-akma@af1.akma.example` a toda requisição sem credenciais. É também a iniciação do AKMA (§6.5).
- **K_AF:** pede-a ao AAnF configurado (§6.7) e a guarda até `expiry`.
- **Verificação do UE:** confere o `response` do Digest com a senha base64(K_AF).
- **Prova da chave:** responde com `rspauth` no `Authentication-Info`, provando ao UE que tem a mesma K_AF.
- **Expiração:** nunca usa uma K_AF vencida. Pela política configurada, pede-a de novo ou recusa com `KAF_EXPIRED` (§6.4.2).

`[MVP]` O AF roda na VM do AF (192.168.16.129, ADR-0001-R2), sem TLS (ADR-0003; limitações em `ua-star.md` §10). O AAnF o autoriza pelo endereço de origem (DV-1).

## 2. Arquivos

Módulo Go `af/` (`github.com/florianoneto07/akma-free5gc/af`). Depende de `shared/akma` (MVP-0) e do `free5gc/openapi` corrigido (D6), por `replace`; usa as mesmas versões de `free5gc/util` e `logrus` do AAnF.

| Caminho | Conteúdo |
|---|---|
| `af/uastar/` | Digest de RFC 7616 com SHA-256, usado pelo AF e pelo cliente: cabeçalhos `WWW-Authenticate`, `Authorization` e `Authentication-Info`; `response` e `rspauth`; `realm`; *product token* |
| `af/uastar/testdata/gen_digest_vectors.py`, `digest_vectors.json` | implementação Python independente e os vetores dela: o exemplo de RFC 7616 §3.9.1 e o vetor `testbed-ue` com `qop=auth` e `auth-int` |
| `af/af/internal/config/` | `afcfg.yaml` e validação (FQDN como nome DNS, identificador Ua\*, `qop`, política de expiração) |
| `af/af/internal/aanf/` | `Naanf_AKMA_ApplicationKey_Get` pelo cliente do `openapi`; `K_AKMA_NOT_PRESENT` e demais falhas |
| `af/af/internal/server/` | servidor Ua\*: verificações na ordem de `ua-star.md` §5.4, desafios (nonce, opaque, `nc`), K_AF por A-KID |
| `af/af/internal/logger/` | log no formato do free5GC, NF `AF`, categorias `Main`, `CFG`, `UA` e `Consumer` |
| `af/af/cmd/af/` | o programa `af` (`-c` configuração, `-l` arquivo de log) |
| `af/ue/uaclient/` | cliente Ua\* (biblioteca e comando `uaclient`) |
| `af/config/afcfg.yaml` | configuração do laboratório |
| `af/scripts/build-af.sh` | compila `af/build/bin/af` e `af/build/bin/uaclient` |

## 3. Configuração

Os campos e os padrões estão em `ua-star.md` §11. No laboratório:

| Campo | Valor |
|---|---|
| `fqdn` | `af1.akma.example` |
| `uaStarId` | `FF00000001` |
| `listen` | `192.168.16.129:8080` |
| `kafExpiry` | `refetch` |
| `aanfUri` | `http://192.168.16.130:8010` |

O AF_ID resultante, `af1.akma.example.FF00000001`, precisa estar em `akma.authorizedAfs` do `aanfcfg.yaml`, com o endereço do AF (`docs/free5gc-side/aanf.md` §2).

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
| `info` | `AF_ID …; realm "…"; qop […]; nonce lifetime …; K_AF expiry policy …; AAnF …` (partida) |
| `info` | `AKMA initiation: challenge sent to <endereço> (product token 3gpp-akma: true\|false)` |
| `info` | `K_AF for A-KID <A-KID> obtained from the AAnF, expires <instante>` |
| `info` | `UE with A-KID <A-KID> authenticated (nc <nc>, qop <qop>)` |
| `info` | `No K_AKMA at the AAnF for A-KID <A-KID>: access rejected` |
| `info` | `K_AF for A-KID <A-KID> expired at <instante>: requesting it again from the AAnF (policy refetch)` ou `…: access rejected (policy reject)` |
| `warn` | `Digest response from <endereço> does not match the K_AF of A-KID <A-KID>: authentication failed` (RFC 7616 §3.4: registrar falhas de *login*) |
| `warn` | `Request from <endereço> for host "<Host>", not <FQDN>`; `Realm … is not …`; `Nonce count … already used …: possible replay` |
| `error` | `K_AF for A-KID <A-KID> not obtained: <motivo>` |
| `debug` | o SUPI recebido do AAnF |

A K_AF fica em `kdf.Key`, que se imprime como `[256-bit key]`; nem ela nem a senha do Digest aparecem em log.

## 5. Como reproduzir

Na máquina do AF, com Go 1.26.2 e o repositório clonado:

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

O teste de integração, com o AAnF real, e o resultado de cada critério estão em `teste-af.md`.
