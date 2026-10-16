# Baseline de especificações

**Release alvo:** 3GPP Release 19
**Data de fixação:** 02 de setembro de 2026

Todas as decisões de implementação deste projeto referem-se às versões abaixo. Qualquer alteração de versão exige revisão da matriz de rastreabilidade.

## Especificações normativas

| Especificação | Título abreviado | Versão | Papel no AKMA |
|---|---|---|---|
| TS 33.535 | AKMA based on 3GPP credentials in the 5GS | V19.0.0 | Stage 2 do AKMA |
| TS 33.501 | Security architecture and procedures for 5G System | V19.7.0 | K_AUSF, SUPI como P0 (Anexo A.7.0), `Nudm_UECM_AuthTrigger` (§14.2.6) |
| TS 33.220 | Generic Authentication Architecture; GBA | V19.1.0 | KDF (Anexo B.2), identificador Ua* (Anexo H), B-TID |
| TS 33.222 | Access to network application functions using HTTPS | V19.0.0 | Contexto dos perfis TLS; não menciona AKMA |
| TS 23.003 | Numbering, addressing and identification | V19.7.0 | Routing Indicator (§2.2B), Home Network Domain (§28.2) |
| TS 23.502 | Procedures for the 5G System | V19.8.0 | Contexto de procedimentos 5GS |
| TS 24.109 | Bootstrapping interface and network application function interface; Stage 3 | V19.0.0 | **Stage 3 do Ua*** (Anexo K) |
| TS 29.535 | AKMA Anchor Services; Stage 3 | V19.5.0 | Stage 3 do Naanf_AKMA |
| TS 29.522 | Network Exposure Function Northbound APIs | V19.7.0 | `AKId`, `AfId`, `AkmaAfKeyRequest`, `AkmaAfKeyData`, `ServiceDisableNotif` |
| TS 29.503 | Unified Data Management Services; Stage 3 | V19.7.0 | `akmaInd`, `routingId`, `id-translation-result`, `trigger-auth`, `ROAMING_STATUS` |
| TS 29.509 | Authentication Server Services; Stage 3 | V19.6.0 | Precedente de codificação de material de chave |
| TS 29.510 | Network Function Repository Services; Stage 3 | V19.7.0 | `NFType` `AANF`, `AanfInfo`, descoberta por `routing-indicator` |
| TS 29.571 | Common Data Types for Service Based Interfaces | V19.7.0 | `Supi`, `Gpsi`, `SupportedFeatures`, `VarUeId` |
| TS 29.500 | Technical Realization of Service Based Architecture | V19.7.0 | HTTP/2, códigos de estado, cabeçalhos customizados |
| TS 29.501 | Principles and Guidelines for Services Definition | V19.5.0 | Estrutura de URI, operações customizadas |
| TS 29.122 | T8 reference point for Northbound APIs | V19.7.0 | `Uri`, `DateTime` |
| TS 29.505 | Usage of the Unified Data Repository services for Subscription Data; Stage 3 | **a fixar** — norma a obter antes do MVP-2 (plano, D8) | Semântica de `akmaAllowed` e `routingId` no `AuthenticationSubscription` |

**Nota (16/10/2026):** o baseline de TS 23.502 permanece V19.8.0; a cópia no acervo é a V19.9.0, usada para consulta. Ver [Adendo E](../05-adendo-E-ts23502-v19-9-0.md).

## Arquivos OpenAPI

| Arquivo | Versão | Corresponde a |
|---|---|---|
| `TS29535_Naanf_AKMA.yaml` | 1.2.0 | TS 29.535 V19.5.0 |
| `TS29522_AKMA.yaml` | — | TS 29.522 V19.7.0, Anexo A.12 |
| `TS29503_Nudm_UEAU.yaml` | 1.4.0 | TS 29.503 |
| `TS29503_Nudm_UECM.yaml` | 1.4.2 | TS 29.503 |
| `TS29503_Nudm_SDM.yaml` | 2.4.2 | TS 29.503 |
| `TS29503_Nudm_EE.yaml` | 1.4.1 | TS 29.503 |
| `TS29509_Nausf_UEAuthentication.yaml` | — | TS 29.509 |
| `TS29510_Nnrf_NFManagement.yaml` | — | TS 29.510 |
| `TS29510_Nnrf_NFDiscovery.yaml` | — | TS 29.510 |
| `TS29571_CommonData.yaml` | — | TS 29.571 |
| `TS29122_CommonData.yaml` | — | TS 29.122 |

Os arquivos OpenAPI declaram, em `externalDocs`, a versão da TS em que foram alterados pela última vez. Versões heterogêneas entre arquivos da mesma TS são normais e não indicam inconsistência do acervo.

**Nota sobre precedência.** TS 29.535 Anexo A.1, TS 29.522 e TS 29.503 Anexo A.1 contêm cláusula segundo a qual o Anexo OpenAPI prevalece em caso de discrepância quanto à codificação de elementos de informação e métodos. O projeto aplica essa regra, com a exceção registrada no defeito D-B.

## Base de código

| Item | Valor | Estado |
|---|---|---|
| free5GC — versão | v4.2.3 | Fixado em 15/10/2026 (plano, D4) |
| free5GC — *commit hash* | `3b34a08e93a9b334f0f4005d3a3a9f79b66d59b9` | Fixado |
| Submódulos das NFs — *commit hash* | tabela em `ambiente-verificado.md` §2.1 | Fixado |
| `free5gc/openapi` / `free5gc/util` | v1.2.4 / v1.3.2 (AMF: `openapi` v1.2.5-0.20260527003827-02dc71b4d94f) | Fixado |
| Versão de Go do ambiente | `go1.26.2` | Fixado |
| Versão de kernel da máquina do free5GC | `5.4.0-150-generic` | Fixado; relevante apenas para o UPF |

Preenchido em 15/10/2026 a partir da inspeção registrada em [`ambiente-verificado.md`](ambiente-verificado.md).
