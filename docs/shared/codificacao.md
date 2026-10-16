# Codificação de chaves e identificadores por camada

**Lado:** comum
**Código:** `shared/akma/codec`, `shared/akma/akid`
**MVP:** MVP-0
**Decisões:** ADR-0002

Internamente, K_AUSF, K_AKMA, A-TID e K_AF são sempre 32 octetos (`kdf.Key`). A forma textual depende da camada e é aplicada só na fronteira.

| Camada | Objeto | Forma | Fundamento |
|---|---|---|---|
| SBI (JSON) | `kAkma` (TS 29.535), `kaf` (TS 29.522) | 64 caracteres hexadecimais; minúsculas na escrita; padrão ancorado `^[A-Fa-f0-9]{64}$` na leitura | `[INTERPRETAÇÃO]` as normas declaram só `type: string` (defeito D-D); precedente de `Kausf` em TS 29.503 (`^[A-Fa-f0-9]{64}$`) e de `Kseaf` em TS 29.509 |
| NAI | valor do A-TID dentro do A-KID | base64 RFC 4648 §4, 44 caracteres | `[INTERPRETAÇÃO]` lacuna C; analogia com o B-TID (TS 33.220 §4.5.2) |
| Ua\* (HTTP Digest) | K_AF como senha | base64 RFC 4648 | `[3GPP]` TS 24.109 K.2 c)2) |
| `AfId` | identificador Ua\* | 10 caracteres hexadecimais, maiúsculas na escrita | `[INTERPRETAÇÃO]` exemplo de TS 29.522; decisão P9 |

Confundir as camadas produz falha silenciosa: K_AF certa no SBI e autenticação recusada no Ua\*, ou o contrário.

`kdf.Key` não revela o valor ao ser formatada (`fmt` imprime um marcador fixo). Para obter o texto de uma chave é preciso chamar `codec` explicitamente — isso evita chaves em log.
