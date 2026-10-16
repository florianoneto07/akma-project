# Matriz de rastreabilidade

**Revisão:** 3.1 (consolida o baseline v2.0, os Adendos A a E, o ADR-0001-R2 e as decisões de 16/10/2026)
**Data:** 02 de setembro de 2026; revisão 3.1 em 16 de outubro de 2026

*Nota (16/10/2026):* os caminhos da coluna "Implementação" seguem a estrutura anterior ao plano. A estrutura aprovada (plano §3) é: `akma/` → `shared/akma/`; `NFs/aanf/` → `free5gc/aanf/`; `patches/` → `free5gc/patches/`; `af/` → `af/af/`; `ue/` → `af/ue/`. Os caminhos serão atualizados linha a linha à medida que cada MVP os criar. Os desvios citados como DV-n estão em `docs/conformidade/desvios.md`.

Legenda de estado:

- **Especificado** — parâmetros integralmente determinados por cláusula; pronto para implementação
- **Parcial** — parcialmente determinado; ver questão associada
- **Decisão** — determinado por decisão do projeto, registrada em ADR ou no registro de defeitos
- **Pendente** — fora do escopo do MVP ou aguardando decisão
- **Implementado** — código e testes no repositório, com o MVP indicado

## Derivação criptográfica

| Requisito | Norma | Cláusula | Implementação | Teste | Estado |
|---|---|---|---|---|---|
| KDF genérica HMAC-SHA-256 | TS 33.220 V19.1.0 | B.1, B.2.0, B.2.1 | `shared/akma/kdf/kdf.go` | `kdf_test.go`, `vectors_test.go`, `shared/crosscheck` | **Implementado** (MVP-0) |
| Codificação de cadeias UTF-8 + NFKC | TS 33.220 V19.1.0 | B.2.1.2 | `shared/akma/kdf/kdf.go` (`String`) | `kdf_test.go`; vetores `nfkc-*` | **Implementado** (MVP-0) |
| Codificação de comprimentos em 2 octetos | TS 33.220 V19.1.0 | B.2.1.3 | `shared/akma/kdf/kdf.go` | `kdf_test.go` | **Implementado** (MVP-0) |
| Derivação de K_AKMA — FC `0x80` | TS 33.535 V19.0.0 | A.2 | `shared/akma/derive/derive.go` | `vectors_test.go`, `shared/crosscheck` | **Implementado** (MVP-0) |
| Derivação de A-TID — FC `0x81` | TS 33.535 V19.0.0 | A.3 | `shared/akma/derive/derive.go` | `vectors_test.go`, `shared/crosscheck` | **Implementado** (MVP-0) |
| Derivação de K_AF — FC `0x82` | TS 33.535 V19.0.0 | A.4 | `shared/akma/derive/derive.go` | `vectors_test.go`, `shared/crosscheck` | **Implementado** (MVP-0) |
| SUPI como P0/P1, sem prefixo `imsi-` | TS 33.501 V19.7.0 | A.7.0 | `shared/akma/supi/supi.go` | `supi_test.go`, `derive_test.go`; contraexemplo nos vetores | **Implementado** (MVP-0); só IMSI (DV-3) |

## Identificadores

| Requisito | Norma | Cláusula | Implementação | Teste | Estado |
|---|---|---|---|---|---|
| Identificador Ua* — 5 octetos | TS 33.220 V19.1.0 | H.1–H.3 | `shared/akma/afid/uastar.go` | `afid_test.go`, `derive_test.go` | **Implementado** (MVP-0) |
| `AfId` de transporte — `<FQDN>.<10 hex>` | TS 29.522 V19.7.0 | Tab. 5.14.5.4.2-1 | `shared/akma/afid/afid.go` | `afid_test.go` | **Implementado** (MVP-0) — decisão P9: maiúsculas na escrita, as duas caixas na leitura |
| Conversão `AfId` → P0 da KDF | TS 33.535 A.4 + TS 29.522 | — | `shared/akma/afid/afid.go` | `afid_test.go`; contraexemplo nos vetores | **Implementado** (MVP-0) — decisão ADR-0002 |
| Home Network Domain | TS 23.003 V19.7.0 | 2.2, 28.2 | `shared/akma/akid/realm.go` | `akid_test.go` | **Implementado** (MVP-0) |
| Routing Indicator, cadeia, valor padrão `"0"` | TS 23.003 V19.7.0 | 2.2B | `shared/akma/akid/rid.go` | `akid_test.go` — casos `"0"`, `"012"`, `"12"`, `"0000"` | **Implementado** (MVP-0) |
| A-KID — formato NAI `rid<v>.atid<v>@<realm>` | TS 33.535 V19.0.0; TS 29.522 V19.7.0 | 6.1; Tab. 5.14.5.4.2-1 | `shared/akma/akid/akid.go`, `shared/akma/derive/derive.go` | `akid_test.go`, `vectors_test.go` | **Implementado** (MVP-0) — defeito D-A: `atid` na escrita, `atid`/`akid` na leitura |
| Codificação do valor do A-TID | — | — | `shared/akma/akid/akid.go` | `akid_test.go`, `vectors_test.go` | **Implementado** (MVP-0) — decisão ADR-0002 (base64; base64url e hexadecimal configuráveis) |
| Codificação de `kAkma` e `kaf` | TS 29.503 (`Kausf`) e TS 29.509 (precedentes) | — | `shared/akma/codec/key.go` | `key_test.go` | **Implementado** (MVP-0) — decisão ADR-0002 |

## API Naanf_AKMA

| Requisito | Norma | Cláusula | Implementação | Teste | Estado |
|---|---|---|---|---|---|
| `POST /naanf-akma/v1/register-anchorkey` | TS 29.535 V19.5.0 | 4.2.2.2.2, 5.1.4.2 | `NFs/aanf/internal/sbi/` | API | Especificado |
| `POST /naanf-akma/v1/retrieve-applicationkey` | TS 29.535 + TS 29.522 | 4.2.2.3.2, 5.14.5.3 | `NFs/aanf/internal/sbi/` | API | Especificado |
| `POST /naanf-akma/v1/remove-context` | TS 29.535 V19.5.0 | 4.2.2.4.2, 5.1.4.4 | `NFs/aanf/internal/sbi/` | API | Especificado |
| Erro `K_AKMA_NOT_PRESENT` (403) | TS 29.535 V19.5.0 | 5.1.7.3 | `NFs/aanf/internal/sbi/` | negativo | Especificado |
| Erro `AKMA_CONTEXT_NOT_FOUND` (404) | TS 29.535 V19.5.0 | 5.1.7.3 | `NFs/aanf/internal/sbi/` | negativo | Especificado |
| `retrieve-applicationkey` nunca responde 204 | TS 29.535 V19.5.0 | 4.2.2.3.2 | `NFs/aanf/internal/sbi/` | negativo | Decisão — DV-7, defeito D-G |
| AF_ID autorizado associado ao AF que o pede | TS 33.535 V19.0.0 | 4.2.2, 6.2.1 passo 2 | `NFs/aanf/internal/sbi/` | negativo | Especificado — risco residual sem TLS (DV-1) |
| Códigos HTTP obrigatórios para POST | TS 29.500 V19.7.0 | 5.2.7.1 | `NFs/aanf/internal/sbi/` | negativo | Especificado — 11 códigos |
| `oneOf` supi/gpsi em `AkmaKeyInfo` | `TS29535_Naanf_AKMA.yaml` v1.2.0 | Anexo A.1 | `NFs/aanf/internal/sbi/` | API | Especificado |

## Contexto do AAnF

| Requisito | Norma | Cláusula | Implementação | Teste | Estado |
|---|---|---|---|---|---|
| Contexto AKMA por SUPI | TS 33.535 V19.0.0 | 6.6.1, 6.2.1 | `NFs/aanf/internal/context/` | unitário | Especificado |
| Substituição de (A-KID, K_AKMA) em reautenticação | TS 33.535 V19.0.0 | 6.1 NOTA 1a | `NFs/aanf/internal/context/` | integração | Especificado |
| A-KID transmitido por AF | TS 33.535 §6.8.1 NOTA 1, §6.2.1 passo 8 (exigido com subscrição de *roaming*); TS 29.522 `ServiceDisableNotif` | — | `NFs/aanf/internal/context/` | unitário | Especificado |
| Estado de *roaming* por tipo de acesso | TS 33.535 §6.2.1 passo 7; `RoamingStatusReport` | — | `NFs/aanf/internal/context/` | unitário | Pendente |
| `expiry` como instante absoluto | TS 29.522 + TS 29.122 | — | `NFs/aanf/internal/context/` | negativo — desvio de relógio | Especificado |

## Registro e descoberta

| Requisito | Norma | Cláusula | Implementação | Teste | Estado |
|---|---|---|---|---|---|
| Registro do AAnF no NRF | TS 29.510 V19.7.0 | — | `NFs/aanf/internal/sbi/consumer/` | integração | Especificado |
| `nfType = "AANF"`, `serviceName = "naanf-akma"` | TS 29.510 V19.7.0 | — | `NFs/aanf/pkg/factory/`, `patches/openapi/` | integração | Especificado |
| `AanfInfo.routingIndicators` no perfil | TS 29.510 V19.7.0 | — | `NFs/aanf/internal/sbi/consumer/` | integração | Especificado |
| Descoberta por `routing-indicator` | TS 29.510 V19.7.0; TS 33.535 §6.7 | — | `patches/ausf/internal/sbi/consumer/` | integração | Especificado |
| Endereço do AAnF alcançável a partir da VM do AF | ADR-0001-R2 §5.2 | — | configuração | integração | Especificado |

## Integração AUSF e UDM

| Requisito | Norma | Cláusula | Implementação | Teste | Estado |
|---|---|---|---|---|---|
| `akmaInd` e `routingId` na resposta do UDM | TS 33.501 V19.7.0 + TS 29.503 §6.3.6.2.3 | — | `patches/udm/`, `patches/openapi/` | integração | Especificado |
| Assinatura AKMA no perfil do assinante | TS 33.535 V19.0.0 | 4.2.5 | `patches/udm/`, *webconsole* | integração | Especificado |
| Acionamento do AKMA no AUSF — 5G-AKA | TS 33.535 V19.0.0 | 6.1 passo 3 | `patches/ausf/` — `Auth5gAkaComfirmRequestProcedure` | integração | Especificado |
| Acionamento do AKMA no AUSF — EAP-AKA' | TS 33.535 V19.0.0 | 6.1 passo 3 | `patches/ausf/` — `EapAuthComfirmRequestProcedure` | integração | Especificado |
| Registro no AAnF sem bloquear a resposta ao AMF; falha no AAnF não altera a autenticação | TS 33.535 V19.0.0 | 6.1 passo 3 | `patches/ausf/` | integração | Decisão — D9 do plano |
| Consumidor `Naanf_AKMA` no AUSF | TS 29.535 V19.5.0 | 4.2.2.2 | `patches/ausf/internal/sbi/consumer/` | integração | Especificado |
| GPSI via `GET /{ueId}/id-translation-result` | TS 29.503 V19.7.0 | 5.2.2.2.10 | — | — | Pendente — política local |
| Subscrição `ROAMING_STATUS` com `immediateFlag` | TS 29.503 V19.7.0 | 5.5.2.2.2, 6.4.6.2.8 | — | — | Decisão — DV-2 (configurável, desabilitada) |
| `POST /{ueId}/registrations/trigger-auth` | TS 33.501 §14.2.6; TS 29.503 §5.3.2.16 | — | — | — | Pendente — defeito D-B |

## Ponto de referência Ua*

| Requisito | Norma | Cláusula | Implementação | Teste | Estado |
|---|---|---|---|---|---|
| Ua* simplificado com id de faixa privada | TS 33.220 H.2 NOTA 2 | — | `af/`, `ue/` | ponta a ponta | Decisão — ADR-0003 |
| Perfil HTTP Digest — username, senha, realm | TS 24.109 V19.0.0; TS 33.535 V19.0.0 | K.2; B.1.2.2 | `af/` | ponta a ponta | Decisão — DV-5 (RFC 7616), DV-6 (`realm = "3GPP-bootstrapping-akma@<FQDN>"`); defeitos D-C, D-F |
| Confirmação de chave nos dois sentidos (`rspauth`) | TS 33.535 V19.0.0; RFC 7616 | 4.4.1 | `af/`, `ue/` | ponta a ponta | Especificado — critério do MVP-3 (`docs/literatura/`) |
| Perfil TLS 1.2 PSK | TS 24.109 V19.0.0 | K.3.2 | `af/` | ponta a ponta | Pendente |
| Perfil TLS 1.3 PSK — `"3GPP-AKMA" ";" A-KID` | TS 24.109 V19.0.0 | K.3.3 | `af/` | ponta a ponta | Pendente — defeito D-E |
| Ausência de renovação de K_AF | TS 33.535 §6.4.3; TS 24.109 K.2, K.3 | — | — | — | Especificado |

## Requisitos de ambiente

| Requisito | Origem | Verificação | Estado |
|---|---|---|---|
| Sincronização NTP entre a VM do AF e a do free5GC (par AAnF↔AF) | ADR-0001-R2 §5.3 | `timedatectl` nas duas VMs | Especificado |
| AF alcança a porta do AAnF | ADR-0001-R2 §5.1, §5.2 | requisição a partir de `192.168.16.129` | Especificado |
| Proteção das interfaces SBI | TS 33.535 §4.4.0; TS 33.501 §13.1.0 | TLS habilitado | Decisão — DV-1; TLS no MVP-6+ |
