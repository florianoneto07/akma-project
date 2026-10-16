# Identificadores do AKMA

**Lado:** comum
**Código:** `shared/akma/akid`, `shared/akma/afid`, `shared/akma/supi`
**MVP:** MVP-0

---

## 1. Routing Indicator (RID)

`[3GPP]` TS 23.003 V19.7.0 §2.2B, item 3: 1 a 4 dígitos decimais, **cada um significativo** — `"012"` não é `"12"`. Sem RID configurado no USIM ou no ME, o valor é o dígito único `"0"`. `[3GPP]` `routingId` (TS 29.503) e `AanfInfo.routingIndicators` (TS 29.510) usam o padrão `^[0-9]{1,4}$`.

O RID é sempre cadeia (`akid.RID`). Convertê-lo para inteiro em qualquer ponto da cadeia UDM → AUSF → A-KID → AAnF descarta zeros à esquerda e quebra só os assinantes cujo RID começa com zero. O UE do testbed usa `"0000"`.

`[3GPP]` O AUSF usa o RID recebido do UDM (TS 33.535 §6.1); o AF obtém o RID do próprio A-KID (§6.7 NOTA 1).

## 2. Home Network Domain

`[3GPP]` TS 23.003 §28.2: `5gc.mnc<MNC>.mcc<MCC>.3gppnetwork.org`, com MNC e MCC em três dígitos; MNC de dois dígitos recebe um `"0"` à esquerda. Exemplo da norma: MCC 345, MNC 12 → `5gc.mnc012.mcc345.3gppnetwork.org`.

`[3GPP]` O comprimento do MNC depende do MCC (TS 23.003 §2.2) e não pode ser deduzido do IMSI. A biblioteca exige MCC e MNC explícitos (`akid.NewHomeNetwork`). Domínios de SNPN (`5gc.nid<NID>…`) não são tratados no MVP.

## 3. A-KID

`[3GPP]` TS 33.535 §6.1: NAI `username@realm` (RFC 7542 §2.2); o `username` contém o RID e o A-TID; o `realm`, o identificador da rede doméstica. `[3GPP]` TS 29.522 Tab. 5.14.5.4.2-1, tipo `AKId`: `"rid<value>.atid<value>@<realm>"`, com o Home Network Domain no `realm`.

| Elemento | Escolha | Fundamento |
|---|---|---|
| Rótulo do A-TID | `atid` na escrita; `atid` e `akid` aceitos na leitura | `[DIVERGÊNCIA]` o texto de TS 29.522 diz `atid`, o exemplo usa `akid` (defeito D-A) |
| Codificação do A-TID | base64 RFC 4648 §4, com preenchimento — 44 caracteres | `[INTERPRETAÇÃO]` nenhuma norma fixa (lacuna C); analogia com o B-TID do GBA (ADR-0002). `base64url` e hexadecimal configuráveis |
| Comprimento do A-TID | 32 octetos | `[3GPP]` saída da KDF (TS 33.535 A.3) |

Exemplo (UE do testbed, vetor `testbed-ue`):

```
rid0000.atidiz1FN8DkMGyxwtfTVYsKUsWnzOJdBOLlQEiQlK2gHNA=@5gc.mnc093.mcc208.3gppnetwork.org
```

O exemplo de TS 29.522 (`rid012.akid019345346@…`) é só sintático: `019345346` não pode ser um A-TID de 256 bits, e a biblioteca o recusa.

`[3GPP]` No perfil OSCORE (TS 33.535 D.3.5), o A-KID vai em CBOR como `[RID : tstr, A-TID : bstr, HPLMN-ID : tstr]` — o A-TID é tratado como cadeia de octetos. Isso não fixa a forma textual na NAI.

## 4. Identificador do AF (AF_ID)

| Forma | Definição | Uso |
|---|---|---|
| KDF | `UTF-8/NFKC(FQDN) ‖ 5 octetos` (TS 33.535 A.4; `[INTERPRETAÇÃO]` codificação do FQDN por analogia com o `NAF_Id`) | P0 de K_AF |
| Transporte | `"<FQDN>.<identificador Ua*>"` (TS 29.522, tipo `AfId`) | campo `afId` do `AkmaAfKeyRequest` |

`[3GPP]` TS 29.522 descreve o identificador no `AfId` como *"a string of 5 octet"*. `[INTERPRETAÇÃO]` A forma em dez caracteres hexadecimais vem do exemplo (`www.app1.com.0100BC0001`). **Decisão do projeto (P9, 16/10/2026):** escrita em maiúsculas, como no exemplo; leitura aceita as duas caixas.

O AAnF recebe a forma de transporte e reconstrói a forma da KDF (`afid.Parse` → `KDFParam`). Usar a cadeia de transporte como P0 dá outra K_AF, sem erro.

## 5. Identificador de protocolo de segurança Ua*

`[3GPP]` TS 33.220 Anexo H: cinco octetos; o primeiro identifica a organização, `0xFF` é a faixa privada (H.2 NOTA 2).

| Constante | Valor | Uso |
|---|---|---|
| `afid.UaStarMVP` | `FF 00 00 00 01` | Ua\* simplificado do MVP (ADR-0003) |
| `afid.UaStarTLS13` | `01 00 00 00 07` | TLS 1.3 com chave compartilhada (H.3) — evolução futura |

Identificadores diferentes dão K_AF diferentes para o mesmo AF (vetores `testbed-ue` e `ua-tls13`).
