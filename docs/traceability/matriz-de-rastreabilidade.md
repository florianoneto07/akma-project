# Matriz de rastreabilidade

**Revisão:** 3.4 (MVP-2: UDM e AUSF; MVP-1b: NRF corrigido; D8, D10 e D13 de 17/10/2026; DV-8 encerrado)
**Data:** 02 de setembro de 2026; revisão 3.1 em 16 de outubro de 2026; revisões 3.2 a 3.4 em 17 de outubro de 2026

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
| `POST /naanf-akma/v1/register-anchorkey` | TS 29.535 V19.5.0 | 4.2.2.2.2, 5.1.4.2 | `free5gc/aanf/internal/sbi/processor/register_anchorkey.go` | `TestRegisterAnchorKey`, `TestRegisterBadRequests`; VM | **Implementado** (MVP-1) |
| `POST /naanf-akma/v1/retrieve-applicationkey` | TS 29.535 + TS 29.522 | 4.2.2.3.2, 5.14.5.3 | `free5gc/aanf/internal/sbi/processor/retrieve_applicationkey.go` | `TestRetrieveApplicationKey`, `TestRetrieveAnonymous`; VM (K_AF = vetor) | **Implementado** (MVP-1); `notifUri` ignorado (sem `RoamingRestriction`) |
| `POST /naanf-akma/v1/remove-context` | TS 29.535 V19.5.0 | 4.2.2.4.2, 5.1.4.4 | `free5gc/aanf/internal/sbi/processor/remove_context.go` | `TestRemoveContext`; VM | **Implementado** (MVP-1) |
| Erro `K_AKMA_NOT_PRESENT` (403) | TS 29.535 V19.5.0 | 5.1.7.3 | `free5gc/aanf/internal/sbi/processor/` | `TestRetrieveKAkmaNotPresent`; VM | **Implementado** (MVP-1) |
| Erro `AKMA_CONTEXT_NOT_FOUND` (404) | TS 29.535 V19.5.0 | 5.1.7.3 | `free5gc/aanf/internal/sbi/processor/` | `TestRemoveContext`; VM | **Implementado** (MVP-1) |
| `retrieve-applicationkey` nunca responde 204 | TS 29.535 V19.5.0 | 4.2.2.3.2 | `free5gc/aanf/internal/sbi/processor/` | conferido em todo teste de API | **Implementado** (MVP-1) — DV-7, defeito D-G |
| AF_ID autorizado associado ao AF que o pede | TS 33.535 V19.0.0 | 4.2.2, 6.2.1 passo 2 | `free5gc/aanf/internal/context/context.go` (`AfPolicy`) | `TestRetrieveAfNotAuthorized`; VM | **Implementado** (MVP-1) — risco residual sem TLS (DV-1) |
| Códigos HTTP obrigatórios para POST | TS 29.500 V19.7.0 | 5.2.7.1, 5.2.7.2 | `free5gc/aanf/internal/sbi/api_akma.go` | 8 códigos testados; 429, 502, 503 justificados (`docs/free5gc-side/aanf.md` §6) | **Implementado** (MVP-1) |
| Método não suportado → 501; versão de API desconhecida → 400 `INVALID_API` | TS 29.500 V19.7.0 | 5.2.7.2, Tab. 5.2.7.2-1 | `free5gc/aanf/internal/sbi/api_akma.go` | `TestRoutingErrors` | **Implementado** (MVP-1) |
| OAuth2 em `register-anchorkey` e `remove-context`: assinatura, escopo, `exp`, `aud` | TS 29.535 V19.5.0; TS 33.501 V19.7.0 | 5.1.9; 13.4.1.1.2 | `free5gc/aanf/internal/context/context.go` | `TestOAuth2`; VM (401 sem token; 200 com token do NRF corrigido) | **Implementado** (MVP-1; `aud` estrito desde o MVP-1b, DV-8 encerrado) |
| Negociação de *features* (nenhuma suportada) | TS 29.535 V19.5.0; TS 29.500 | 5.1.8; 6.6 | `free5gc/aanf/internal/sbi/processor/processor.go` | `TestRegisterFeatureNegotiation` | **Implementado** (MVP-1) |
| `oneOf` supi/gpsi em `AkmaKeyInfo` | `TS29535_Naanf_AKMA.yaml` v1.2.0; TS 29.535 | Anexo A.2; Tab. 5.1.6.2.2-1 | `free5gc/aanf/internal/sbi/processor/register_anchorkey.go` | `TestRegisterBadRequests` | **Implementado** (MVP-1) — sem `AKMA_GPSI_Support`, só `supi` |
| Modelos e cliente do Naanf_AKMA no `free5gc/openapi` | TS 29.535 Anexo A.2; TS 29.522 | 5.14.5.3 | `free5gc/patches/openapi/` | `TestOpenapiClient` | **Implementado** (MVP-1) — D6 |

## Contexto do AAnF

| Requisito | Norma | Cláusula | Implementação | Teste | Estado |
|---|---|---|---|---|---|
| Contexto AKMA por SUPI | TS 33.535 V19.0.0 | 6.6.1, 6.2.1 | `free5gc/aanf/internal/context/akma_context.go` | `akma_context_test.go` | **Implementado** (MVP-1) — em memória |
| Substituição de (A-KID, K_AKMA) em reautenticação | TS 33.535 V19.0.0 | 6.1 NOTA 1a | `free5gc/aanf/internal/context/akma_context.go` | `TestReRegistrationReplacesPair`, `TestReRegistrationOverAPI` | **Implementado** (MVP-1) |
| K_AF reaproveitada até expirar; K_AF emitidas mantidas após novo K_AKMA | TS 33.535 V19.0.0 | 6.2.1 passo 7; 5.2 | `free5gc/aanf/internal/context/akma_context.go` | `TestKafReusedUntilExpiry` | **Implementado** (MVP-1) |
| A-KID transmitido por AF | TS 33.535 §6.8.1 NOTA 1, §6.2.1 passo 8 (exigido com subscrição de *roaming*); TS 29.522 `ServiceDisableNotif` | — | `free5gc/aanf/internal/context/akma_context.go` (`AfKeyRecord.AKIDSent`) | unitário | **Implementado** (MVP-1) — guardado; notificação fora do MVP |
| Estado de *roaming* por tipo de acesso | TS 33.535 §6.2.1 passo 7; `RoamingStatusReport` | — | `free5gc/aanf/internal/context/` | unitário | Pendente — DV-2 |
| `expiry` como instante absoluto | TS 29.522 + TS 29.122 | — | `free5gc/aanf/internal/context/akma_context.go` | `TestRetrieveApplicationKey` (RFC 3339, UTC); desvio de relógio no AF fica para o MVP-3 | **Implementado** no AAnF (MVP-1) |

## Registro e descoberta

| Requisito | Norma | Cláusula | Implementação | Teste | Estado |
|---|---|---|---|---|---|
| Registro do AAnF no NRF | TS 29.510 V19.7.0 | — | `free5gc/aanf/internal/sbi/consumer/nrf_service.go` | `TestRegisterNFInstance*`; VM | **Implementado** (MVP-1) |
| `nfType = "AANF"`, `serviceName = "naanf-akma"` | TS 29.510 V19.7.0 | — | `free5gc/aanf/internal/context/context.go` | `TestRegisterNFInstanceProfile`; VM (MongoDB) | **Implementado** (MVP-1) |
| `AanfInfo.routingIndicators` no perfil | TS 29.510 V19.7.0 | Tab. 6.1.6.2.73-1 | `free5gc/aanf/internal/sbi/consumer/nrf_service.go` | `TestRegisterNFInstanceProfile`; VM (MongoDB) | **Implementado** (MVP-1 no AAnF; guardado pelo NRF desde o MVP-1b) |
| NRF guarda `aanfInfoList` | TS 29.510 V19.7.0 | 5.2.2.2.2; Tab. 6.1.6.2.2-1 | `free5gc/patches/nrf/0001-store-aanf-info-list.patch` | `TestNnrfNFManagementDataModelKeepsAanfInfoList`; VM | **Implementado** (MVP-1b) |
| NRF descobre `AANF` por `routing-indicator` | TS 29.510 V19.7.0 | Tab. 6.2.3.2.3.1-1, Tab. 6.1.6.2.73-1 | `free5gc/patches/nrf/0002-discover-aanf-by-routing-indicator.patch` | `TestValidateQueryParametersNfTypes`, `TestBuildFilterAanfRoutingIndicator`; VM (MongoDB 5.0, quatro perfis sintéticos) | **Implementado** (MVP-1b); `complexQuery` fora do escopo |
| NRF preenche `aud` do token | TS 29.510 V19.7.0; TS 33.501 V19.7.0 | Tab. 6.3.5.2.4-1; 13.4.1.1.2 | `free5gc/patches/nrf/0003-access-token-audience.patch` | `TestAccessTokenAudience`; VM | **Implementado** (MVP-1b) |
| Descoberta do AAnF pelo AUSF, por `routing-indicator` | TS 29.510 V19.7.0; TS 33.535 §6.7 | — | `free5gc/patches/ausf/` (`aanfUri`) | `TestRegisterAkmaDiscoversAanf`; UERANSIM (`routing-indicator=0000`, 200) | **Implementado** (MVP-2) — `akma.aanfUri` como alternativa |
| Endereço do AAnF alcançável a partir da VM do AF | ADR-0001-R2 §5.2 | — | `free5gc/config/aanfcfg.yaml` | VM, 17/10/2026 | **Implementado** (MVP-1) |

## Integração AUSF e UDM

| Requisito | Norma | Cláusula | Implementação | Teste | Estado |
|---|---|---|---|---|---|
| `akmaInd` e `routingId` na resposta do UDM | TS 33.501 V19.7.0; TS 29.503 V19.7.0 | 6.1.3.1 passo 2, 6.1.3.2.0 passo 2; Tab. 6.3.6.2.3-1 | `free5gc/patches/udm/` (`internal/sbi/processor/akma.go`) | `TestSetAkmaIndication`, `TestGenerateAuthDataAkma`, `TestGenerateAuthDataWithoutAkma`; banco de testes | **Implementado** (MVP-2) |
| Fonte do RID no UDM: `routingId` do `AuthenticationSubscription` no UDR; comparação com o RID do SUCI, quando houver, com divergência em log | TS 29.505 V19.6.0; TS 29.503 V19.7.0; TS 33.535 §6.1 | Tab. 5.4.2.2-1; Tab. 6.3.6.2.3-1 | `free5gc/patches/udm/` (`suci.RoutingIndicator`, `setAkmaIndication`) | `TestRoutingIndicator`, `TestSetAkmaIndication`; UERANSIM | **Implementado** (MVP-2) — D10 |
| Assinatura AKMA no perfil do assinante (`akmaAllowed`, `routingId`) | TS 33.535 V19.0.0; TS 29.505 V19.6.0 | 3.1, 4.2.5; Tab. 5.4.2.2-1 | `free5gc/scripts/provision-akma-subscriber.sh` (a *webconsole* não tem os campos) | VM (MongoDB) | **Implementado** (MVP-2) |
| Acionamento do AKMA no AUSF — 5G-AKA | TS 33.535 V19.0.0; TS 33.501 V19.7.0 | 6.1 passo 3; 6.1.3.2.0 passo 11 | `free5gc/patches/ausf/` — `Auth5gAkaComfirmRequestProcedure` | `TestAkma5GAKA`; UERANSIM | **Implementado** (MVP-2) |
| Acionamento do AKMA no AUSF — EAP-AKA' | TS 33.535 V19.0.0; TS 33.501 V19.7.0 | 6.1 passo 3; 6.1.3.1 passo 10 | `free5gc/patches/ausf/` — `EapAuthComfirmRequestProcedure` | `TestAkmaEAPAKAPrime` | **Implementado** (MVP-2) |
| Domínio de origem do A-KID: MCC e MNC do SUCI; com SUPI, da PLMN única de `plmnSupportList` | TS 33.535 V19.0.0; TS 23.003 V19.7.0 | 6.1; 2.2, 28.2 | `free5gc/patches/ausf/` (`homeNetwork`) | `TestHomeNetworkFromSuci`, `TestBuildAkmaKeyInfoMatchesOracle` | **Implementado** (MVP-2) — `[INTERPRETAÇÃO]` |
| Registro no AAnF sem bloquear a resposta ao AMF; falha no AAnF não altera a autenticação | TS 33.535 V19.0.0 | 6.1 passo 3 | `free5gc/patches/ausf/` (`startAkma`) | `TestStartAkmaDoesNotBlock`; UERANSIM com AAnF parado | **Implementado** (MVP-2) — D9 |
| Consumidor `Naanf_AKMA` no AUSF | TS 29.535 V19.5.0 | 4.2.2.2 | `free5gc/patches/ausf/` (`internal/sbi/consumer/aanf_service.go`) | `TestRegisterAkmaDiscoversAanf`, `TestRegisterAkmaConfiguredUri`; banco de testes; UERANSIM | **Implementado** (MVP-2) |
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
| AF alcança a porta do AAnF | ADR-0001-R2 §5.1, §5.2 | requisição a partir de `192.168.16.129` | **Verificado** em 17/10/2026 (MVP-1) |
| Proteção das interfaces SBI | TS 33.535 §4.4.0; TS 33.501 §13.1.0 | TLS habilitado | Decisão — DV-1; TLS no MVP-6+ |
