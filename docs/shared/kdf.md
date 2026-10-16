# Derivação de chaves do AKMA

**Lado:** comum (AUSF, AAnF e UE derivam com as mesmas regras)
**Código:** `shared/akma/kdf`, `shared/akma/derive`, `shared/akma/supi`
**MVP:** MVP-0

---

## 1. A KDF genérica

`[3GPP]` TS 33.220 V19.1.0, Anexo B.2.0:

```
S = FC ‖ P0 ‖ L0 ‖ P1 ‖ L1 ‖ … ‖ Pn ‖ Ln
chave derivada = HMAC-SHA-256(Key, S)
```

- `FC`: um octeto (B.2.0 admite também `0xFF ‖ FC2`, não usado pelo AKMA).
- `Li`: comprimento de `Pi` em octetos, em dois octetos, mais significativo primeiro (B.2.1.3). Nenhum parâmetro passa de 65535 octetos (B.1).
- Cadeias de caracteres: UTF-8 com normalização NFKC (B.2.1.2).
- Saída: 256 bits, sem truncamento.

`[DIVERGÊNCIA]` TS 33.535 A.1.1 remete à "KDF do Anexo B.2.2 de TS 33.220"; na V19.1.0, B.2.2 é a alocação de valores de FC e a KDF está em B.2.0. Deriva de referência por renumeração; o projeto segue B.2.0.

`[FREE5GC]` `ueauth.GetKDFValue` (`free5gc/util` v1.3.2) implementa a mesma função — saída integral, comprimentos em dois octetos — mas não aplica NFKC: o chamador entrega os octetos já codificados. A validação cruzada (seção 5) confirma igualdade octeto a octeto.

## 2. As três derivações

`[3GPP]` TS 33.535 V19.0.0, Anexos A.2, A.3 e A.4; FC alocados em A.1.2 e em TS 33.220 B.2.2.

| | K_AKMA (A.2) | A-TID (A.3) | K_AF (A.4) |
|---|---|---|---|
| FC | `0x80` | `0x81` | `0x82` |
| P0 | `"AKMA"` | `"A-TID"` | AF_ID |
| L0 | `0x00 0x04` | `0x00 0x05` | comprimento do AF_ID |
| P1 | SUPI | SUPI | — |
| Chave | K_AUSF | K_AUSF | K_AKMA |

O A-TID deriva de K_AUSF, não de K_AKMA. O AAnF nunca recebe K_AUSF.

## 3. O SUPI como parâmetro

`[3GPP]` TS 33.535 A.2 e A.3: o SUPI é o mesmo valor do P0 de TS 33.501 A.7.0 — para SUPI do tipo IMSI, o IMSI de TS 23.003 §2.2, isto é, só os dígitos. `[3GPP]` Nas interfaces SBI, o SUPI é `"imsi-<imsi>"` (TS 29.571, tipo `Supi`).

**Falha silenciosa.** Passar `"imsi-208930000000001"` à KDF produz K_AKMA, A-TID e K_AF errados, sem erro. O vetor `counterexamples/supi-with-imsi-prefix` de `shared/akma/testdata/vectors.json` mostra o resultado errado.

Na biblioteca, a conversão é feita por `supi.FromSBI`, que devolve o tipo `supi.Param`; `derive.KAkma` e `derive.ATid` só aceitam `supi.Param` e recusam um valor que ainda carregue o prefixo. `[MVP]` Só SUPI do tipo IMSI (desvio DV-3).

## 4. O AF_ID como parâmetro

`[3GPP]` TS 33.535 A.4: `AF_ID = FQDN do AF ‖ identificador de protocolo de segurança Ua*` (TS 33.220 Anexo H, cinco octetos).

`[INTERPRETAÇÃO]` P0 = `UTF-8/NFKC(FQDN) ‖ 5 octetos brutos`, `L0 = len(FQDN) + 5`, por analogia com o `NAF_Id` do GBA (TS 33.220 B.3). A forma de transporte do `AfId` (`"<FQDN>.<10 hexadecimais>"`, TS 29.522) **não** é o P0 — ver `identificadores.md`.

Na biblioteca, `derive.KAf` recebe um `afid.AfID`, não uma cadeia: não há como passar a forma de transporte como P0 por engano. O vetor `counterexamples/af-id-transport-as-p0` mostra o resultado errado que ela produziria.

## 5. Como a correção é verificada

`[MVP]` TS 33.535 V19.0.0 não publica vetores de teste. A correção é verificada por implementações independentes:

| Implementação | Onde | Escopo |
|---|---|---|
| Oráculo Python, escrito a partir das normas | `reference/kdf/` | KDF, derivações, A-KID, AF_ID |
| Biblioteca Go do projeto | `shared/akma` | idem |
| `ueauth.GetKDFValue` do free5GC (terceiros) | `shared/crosscheck` | KDF genérica, nas três derivações |
| ME no UERANSIM, C++ (terceiros + projeto) | MVP-4 | derivações no UE |

Resultado em 16/10/2026: as três primeiras concordam octeto a octeto nos 8 vetores (`vetores.md`).
