# Vetores de teste do projeto

**Lado:** comum
**Arquivo:** `shared/akma/testdata/vectors.json`
**Gerador:** `reference/kdf/gen_vectors.py`

> **Estes vetores não são do 3GPP.** TS 33.535 V19.0.0 não publica vetores de teste. Os valores esperados são calculados pelo oráculo Python do projeto, escrito a partir das normas, e conferidos contra a biblioteca Go e contra a KDF do free5GC.

## 1. Conteúdo

Cada vetor traz as entradas (K_AUSF, SUPI, MCC, MNC, RID, FQDN, identificador Ua\*) e as saídas: cadeias `S` e resultados de K_AKMA, A-TID e K_AF; `realm`; A-KID em base64, base64url e hexadecimal; as duas formas do AF_ID; K_AF em base64.

| Vetor | O que exercita |
|---|---|
| `worked-example` | o exemplo de `docs/00-baseline-normativo.md` §2.5 (chamado `baseline-2.5` até 20/10/2026) |
| `testbed-ue` | o assinante do testbed (`imsi-208930000000001`, RID `"0000"`) |
| `rid-012` | RID com zero à esquerda; MCC 345, MNC 12 (exemplo de TS 23.003 §28.2) |
| `mnc-3-digits` | MNC de três dígitos; RID de um dígito |
| `ua-tls13` | mesmo UE de `testbed-ue` com o identificador do TLS 1.3 |
| `nfkc-composed`, `nfkc-decomposed` | FQDN com `é` composto e decomposto |
| `nfkc-fullwidth` | FQDN em caracteres de largura total |

K_AUSF de cada vetor, exceto `worked-example`, é `SHA-256("akma-free5gc test vector: <nome>")` — determinístico e reproduzível.

O arquivo também registra:

- **igualdades esperadas:** `nfkc-composed` e `nfkc-decomposed` dão a mesma K_AF; `nfkc-fullwidth` dá a mesma K_AF que `testbed-ue`; `ua-tls13` tem a mesma K_AKMA de `testbed-ue` e outra K_AF;
- **contraexemplos:** o resultado errado de usar o SUPI com prefixo `imsi-` e a forma de transporte do AF_ID como parâmetros da KDF.

## 2. Como regenerar e conferir

```bash
python3 reference/kdf/gen_vectors.py > shared/akma/testdata/vectors.json
```

```bash
cd shared/akma && go test ./...
```

```bash
cd shared/crosscheck && go test ./...
```

O primeiro teste confere a biblioteca contra o oráculo; o segundo, contra `ueauth.GetKDFValue` do free5GC.

## 3. Resultado

16/10/2026: os 8 vetores reproduzidos octeto a octeto pela biblioteca Go; K_AKMA, A-TID e K_AF idênticos aos de `ueauth.GetKDFValue`; os três valores do baseline §2.5 conferem.
