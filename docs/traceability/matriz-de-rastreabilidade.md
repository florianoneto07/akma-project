# Matriz de rastreabilidade

**Revisão:** 3 (consolida o baseline v2.0 e os Adendos A a D)
**Data:** 02 de setembro de 2026

Legenda de estado:

- **Especificado** — parâmetros integralmente determinados por cláusula; pronto para implementação
- **Parcial** — parcialmente determinado; ver questão associada
- **Decisão** — determinado por decisão do projeto, registrada em ADR ou no registro de defeitos
- **Pendente** — fora do escopo do MVP ou aguardando decisão

## Derivação criptográfica

| Requisito | Norma | Cláusula | Implementação | Teste | Estado |
|---|---|---|---|---|---|
| KDF genérica HMAC-SHA-256 | TS 33.220 V19.1.0 | B.2.0, B.2.1 | `akma/kdf/kdf.go` | `kdf_test.go` | Especificado |
| Codificação de cadeias UTF-8 + NFKC | TS 33.220 V19.1.0 | B.2.1.2 | `akma/kdf/encoding.go` | `encoding_test.go` | Especificado |
| Codificação de comprimentos em 2 octetos | TS 33.220 V19.1.0 | B.2.1.3 | `akma/kdf/kdf.go` | `kdf_test.go` | Especificado |
| Derivação de K_AKMA — FC `0x80` | TS 33.535 V19.0.0 | A.2 | `akma/kdf/kakma.go` | `kakma_test.go` | Especificado |
| Derivação de A-TID — FC `0x81` | TS 33.535 V19.0.0 | A.3 | `akma/kdf/atid.go` | `atid_test.go` | Especificado |
| Derivação de K_AF — FC `0x82` | TS 33.535 V19.0.0 | A.4 | `akma/kdf/kaf.go` | `kaf_test.go` | Especificado |
| SUPI como P0/P1, sem prefixo `imsi-` | TS 33.501 V19.7.0 | A.7.0 | `akma/kdf/supi.go` | `supi_test.go` | Especificado |

## Identificadores

| Requisito | Norma | Cláusula | Implementação | Teste | Estado |
|---|---|---|---|---|---|
| Identificador Ua* — 5 octetos | TS 33.220 V19.1.0 | H.1–H.3 | `akma/afid/uastar.go` | `uastar_test.go` | Especificado |
| `AfId` de transporte — `<FQDN>.<10 hex>` | TS 29.522 V19.7.0 | 5.14.5.4.2 | `akma/afid/afid.go` | `afid_test.go` | Especificado |
| Conversão `AfId` → P0 da KDF | TS 33.535 A.4 + TS 29.522 | — | `akma/afid/afid.go` | `afid_kdf_test.go` | Decisão — ADR-0002 |
| Home Network Domain | TS 23.003 V19.7.0 | 28.2 | `akma/akid/realm.go` | `realm_test.go` | Especificado |
| Routing Indicator, cadeia, valor padrão `"0"` | TS 23.003 V19.7.0 | 2.2B | `akma/akid/rid.go` | `rid_test.go` — casos `"0"`, `"012"`, `"0000"` | Especificado |
| A-KID — formato NAI `rid<v>.atid<v>@<realm>` | TS 29.522 V19.7.0 | 5.14.5.4.2 | `akma/akid/akid.go` | `akid_test.go` | Parcial — defeito D-A |
| Codificação do valor do A-TID | — | — | `akma/akid/akid.go` | `akid_test.go` | Decisão — ADR-0002 |
| Codificação de `kAkma` e `kaf` | TS 29.509 V19.6.0 (precedente) | — | `akma/codec/` | `codec_test.go` | Decisão — ADR-0002 |

## API Naanf_AKMA

| Requisito | Norma | Cláusula | Implementação | Teste | Estado |
|---|---|---|---|---|---|
| `POST /naanf-akma/v1/register-anchorkey` | TS 29.535 V19.5.0 | 4.2.2.2.2, 5.1.4.2 | `NFs/aanf/internal/sbi/` | API | Especificado |
| `POST /naanf-akma/v1/retrieve-applicationkey` | TS 29.535 + TS 29.522 | 4.2.2.3.2, 5.14.5.3 | `NFs/aanf/internal/sbi/` | API | Especificado |
| `POST /naanf-akma/v1/remove-context` | TS 29.535 V19.5.0 | 4.2.2.4.2, 5.1.4.4 | `NFs/aanf/internal/sbi/` | API | Especificado |
| Erro `K_AKMA_NOT_PRESENT` (403) | TS 29.535 V19.5.0 | 5.1.7.3 | `NFs/aanf/internal/sbi/` | negativo | Especificado |
| Erro `AKMA_CONTEXT_NOT_FOUND` (404) | TS 29.535 V19.5.0 | 5.1.7.3 | `NFs/aanf/internal/sbi/` | negativo | Especificado |
| Códigos HTTP obrigatórios para POST | TS 29.500 V19.7.0 | 5.2.7.1 | `NFs/aanf/internal/sbi/` | negativo | Especificado — 11 códigos |
| `oneOf` supi/gpsi em `AkmaKeyInfo` | `TS29535_Naanf_AKMA.yaml` v1.2.0 | Anexo A.1 | `NFs/aanf/internal/sbi/` | API | Especificado |

## Contexto do AAnF

| Requisito | Norma | Cláusula | Implementação | Teste | Estado |
|---|---|---|---|---|---|
| Contexto AKMA por SUPI | TS 33.535 V19.0.0 | 6.6.1, 6.2.1 | `NFs/aanf/internal/context/` | unitário | Especificado |
| Substituição de (A-KID, K_AKMA) em reautenticação | TS 33.535 V19.0.0 | 6.1 NOTA 1a | `NFs/aanf/internal/context/` | integração | Especificado |
| A-KID transmitido por AF | TS 33.535 §6.8.1 NOTA 1; TS 29.522 `ServiceDisableNotif` | — | `NFs/aanf/internal/context/` | unitário | Especificado |
| Estado de *roaming* por tipo de acesso | TS 33.535 §6.2.1 passo 7; `RoamingStatusReport` | — | `NFs/aanf/internal/context/` | unitário | Pendente |
| `expiry` como instante absoluto | TS 29.522 + TS 29.122 | — | `NFs/aanf/internal/context/` | negativo — desvio de relógio | Especificado |

## Registro e descoberta

| Requisito | Norma | Cláusula | Implementação | Teste | Estado |
|---|---|---|---|---|---|
| Registro do AAnF no NRF | TS 29.510 V19.7.0 | — | `NFs/aanf/internal/sbi/consumer/` | integração | Especificado |
| `nfType = "AANF"`, `serviceName = "naanf-akma"` | TS 29.510 V19.7.0 | — | `NFs/aanf/pkg/factory/`, `patches/openapi/` | integração | Especificado |
| `AanfInfo.routingIndicators` no perfil | TS 29.510 V19.7.0 | — | `NFs/aanf/internal/sbi/consumer/` | integração | Especificado |
| Descoberta por `routing-indicator` | TS 29.510 V19.7.0; TS 33.535 §6.7 | — | `patches/ausf/internal/sbi/consumer/` | integração | Especificado |
| Endereço publicado alcançável entre VMs | ADR-0001 §6.2 | — | configuração | integração | Especificado |

## Integração AUSF e UDM

| Requisito | Norma | Cláusula | Implementação | Teste | Estado |
|---|---|---|---|---|---|
| `akmaInd` e `routingId` na resposta do UDM | TS 33.501 V19.7.0 + TS 29.503 §6.3.6.2.3 | — | `patches/udm/`, `patches/openapi/` | integração | Especificado |
| Assinatura AKMA no perfil do assinante | TS 33.535 V19.0.0 | 4.2.5 | `patches/udm/`, *webconsole* | integração | Especificado |
| Acionamento do AKMA no AUSF — 5G-AKA | TS 33.535 V19.0.0 | 6.1 passo 3 | `patches/ausf/` — `Auth5gAkaComfirmRequestProcedure` | integração | Especificado |
| Acionamento do AKMA no AUSF — EAP-AKA' | TS 33.535 V19.0.0 | 6.1 passo 3 | `patches/ausf/` — `EapAuthComfirmRequestProcedure` | integração | Especificado |
| Consumidor `Naanf_AKMA` no AUSF | TS 29.535 V19.5.0 | 4.2.2.2 | `patches/ausf/internal/sbi/consumer/` | integração | Especificado |
| GPSI via `GET /{ueId}/id-translation-result` | TS 29.503 V19.7.0 | 5.2.2.2.10 | — | — | Pendente — política local |
| Subscrição `ROAMING_STATUS` com `immediateFlag` | TS 29.503 V19.7.0 | 6.4.6.2.8 | — | — | Pendente — questão 1 |
| `POST /{ueId}/registrations/trigger-auth` | TS 33.501 §14.2.6; TS 29.503 §5.3.2.16 | — | — | — | Pendente — defeito D-B |

## Ponto de referência Ua*

| Requisito | Norma | Cláusula | Implementação | Teste | Estado |
|---|---|---|---|---|---|
| Ua* simplificado com id de faixa privada | TS 33.220 H.2 NOTA 2 | — | `af/`, `ue/` | ponta a ponta | Decisão — ADR-0003 |
| Perfil HTTP Digest — username, senha, realm | TS 24.109 V19.0.0 | K.2 | `af/` | ponta a ponta | Especificado — defeito D-C |
| Perfil TLS 1.2 PSK | TS 24.109 V19.0.0 | K.3.2 | `af/` | ponta a ponta | Pendente |
| Perfil TLS 1.3 PSK — `"3GPP-AKMA" ";" A-KID` | TS 24.109 V19.0.0 | K.3.3 | `af/` | ponta a ponta | Pendente — defeito D-E |
| Ausência de renovação de K_AF | TS 33.535 §6.4.3; TS 24.109 K.2, K.3 | — | — | — | Especificado |

## Requisitos de ambiente

| Requisito | Origem | Verificação | Estado |
|---|---|---|---|
| Sincronização NTP entre VMs | ADR-0001 §6.3 | `chronyc tracking` | Especificado |
| Alcançabilidade bidirecional VM1 ↔ VM2 | ADR-0001 §6.1 | `curl` ao NRF e ao AAnF | Especificado |
| Proteção das interfaces SBI | TS 33.535 §4.4.0 | TLS habilitado | Pendente — fase F11 |
