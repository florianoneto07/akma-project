# Matriz de rastreabilidade

**Revisão:** 4 (MVP-5: etapas mínimas do fluxo ligadas a evidências do roteiro ponta a ponta)
**Data:** 02 de setembro de 2026; revisão 3.1 em 16 de outubro de 2026; revisões 3.2 a 3.4 em 17 de outubro de 2026; revisões 3.5 a 3.7 e 4 em 18 de outubro de 2026

*Nota (16/10/2026):* os caminhos da coluna "Implementação" seguem a estrutura anterior ao plano. A estrutura aprovada (plano §3) é: `akma/` → `shared/akma/`; `NFs/aanf/` → `free5gc/aanf/`; `patches/` → `free5gc/patches/`; `af/` → `af/af/`; `ue/` → `af/ue/`. Os caminhos serão atualizados linha a linha à medida que cada MVP os criar. Os desvios citados como DV-n estão em `docs/conformidade/desvios.md`.

Legenda de estado:

- **Especificado** — parâmetros integralmente determinados por cláusula; pronto para implementação
- **Parcial** — parcialmente determinado; ver questão associada
- **Decisão** — determinado por decisão do projeto, registrada em ADR ou no registro de defeitos
- **Pendente** — fora do escopo do MVP ou aguardando decisão
- **Implementado** — código e testes no repositório, com o MVP indicado

## Etapas mínimas do fluxo e evidências (revisão 4)

`[LACUNA]` O plano (MVP-5, critério 3) liga cada linha da "tabela de etapas mínimas" do material de partida a uma evidência. Esse material nunca entrou no repositório (plano §8, I17 e I23). `[INTERPRETAÇÃO]` A tabela abaixo o reconstrói a partir do texto de TS 33.535 V19.0.0, passo a passo: §6.1 (registro após a autenticação primária) e §6.2.1 (K_AF para um AF na rede do operador), com as cláusulas que eles invocam.

Evidências: os identificadores entre colchetes são verificações da execução `20261018-171510` de `af/scripts/e2e.sh` (`docs/af-side/e2e.md`): R na máquina do free5GC; S, N, B, C e L na do AF. Os demais são testes do repositório.

| Etapa | Cláusula | Implementação | Evidência | Estado |
|---|---|---|---|---|
| AUSF obtém do UDM os dados de autenticação (`Nudm_UEAuthentication_Get`) | 6.1 passo 1 | free5GC (sem alteração) | [R1.ue], [R2.registration] | **Verificado** |
| UDM inclui a indicação AKMA e o RID | 6.1 passo 2; TS 33.501 6.1.3.x passo 2; TS 29.503 Tab. 6.3.6.2.3-1 | `free5gc/patches/udm/` | log `[UDM][AKMA] AKMA indication sent … Routing Indicator 0000` em [R2]; ausente sem assinatura [R5.no-indication] | **Verificado** |
| AUSF gera K_AKMA e A-KID após a autenticação bem-sucedida, com o RID do UDM | 6.1 passo 3; A.2, A.3 | `free5gc/patches/ausf/` | [R2.akid-ausf]; `TestAkma5GAKA`, `TestAkmaEAPAKAPrime` | **Verificado** |
| UE gera K_AKMA e A-KID antes de comunicar com um AF | 6.1 passo 3; 6.2.1 passo 1 | `af/ue/ueransim/patches/0001` | A-KID do UE igual ao do AUSF e do AAnF [R2.akid-ausf], [R2.akid-aanf]; [S1] | **Verificado** |
| AUSF seleciona o AAnF | 6.1 passo 4; 6.7 | `free5gc/patches/ausf/`, `free5gc/patches/nrf/0002` | log `AAnF for Routing Indicator 0000 … (discovered at the NRF)` em [R2]; sem AAnF registrado, `no uri for naanf-akma` [R4.ausf-error] | **Verificado** |
| AUSF envia A-KID, K_AKMA e SUPI (`Naanf_AKMA_AnchorKey_Register`); AAnF guarda o mais recente | 6.1 passo 4 e NOTA 1a | `free5gc/patches/ausf/`, `free5gc/aanf/` | [R2.akid-aanf]; substituição [R3.replaced]; A-KID anterior recusado [B1] | **Verificado** |
| AAnF responde ao AUSF | 6.1 passo 5 | `free5gc/aanf/` | log `AKMA key material … registered at the AAnF` em [R2] | **Verificado** |
| Autenticação primária não depende do AAnF | 6.1 passo 3 (D9) | `free5gc/patches/ausf/` (`startAkma`) | [R4.registration], [R4.ausf-error] | **Verificado** — D9 |
| UE inclui o A-KID no pedido de sessão; K_AF derivada pelo ME | 6.2.1 passo 1; A.4 | `0003-akma-kaf-nr-cli-selftest.patch`; `af/ue/uaclient/` | [R2.ua], [R3.ua-new], [R5.ua-restored] | **Verificado** |
| AF sem contexto ativo pede a K_AF ao AAnF com A-KID e AF_ID | 6.2.1 passo 2; 6.7 (configuração local) | `af/af/internal/aanf/` | log do AF `K_AF … obtained from the AAnF`; `TestApplicationKey` | **Verificado** |
| AAnF verifica a política local para o AF | 6.2.1 passo 2 | `free5gc/aanf/` (`AfPolicy`) | [N03] (403); o AF fora da política recebe 503 (MVP-3) | **Verificado** — DV-1 |
| AAnF verifica a presença da K_AKMA | 6.2.1 passo 2 | `free5gc/aanf/` | [N06], [B1], [C1], [N21] (403 `K_AKMA_NOT_PRESENT`) | **Verificado** — DV-7 |
| GPSI pelo UDM | 6.2.1 passos 3 e 4 | — | — | Não aplicável: a política local do MVP não pede GPSI |
| Subscrição de `ROAMING_STATUS` | 6.2.1 passos 5 e 6 | — | — | **Desvio DV-2**: etapa desabilitada |
| AAnF deriva a K_AF se ainda não a tem e guarda o `expiry` | 6.2.1 passo 7; A.4; 5.2 | `free5gc/aanf/` | [N01b]; a mesma `expiry` em pedidos seguintes (log do AAnF em [R2] e [N01]) | **Verificado** |
| AAnF responde com SUPI, K_AF e `expiry` | 6.2.1 passo 8 | `free5gc/aanf/` | [N01b] | **Verificado** |
| AF responde ao UE; com falha, recusa com causa | 6.2.1 passo 9 | `af/af/internal/server/` | 200 [R2.ua]; `403 AKMA_KEY_NOT_AVAILABLE` [R5.ua]; `503 AANF_UNAVAILABLE` [R4.ua] | **Verificado** |
| Iniciação do AKMA pelo AF | 6.5 | `af/af/internal/server/` | `AKMA initiation: challenge sent` no log do AF em cada Ua\*; `TestInitiation` | **Verificado** |
| K_AF vencida no AF; K_AF anterior válida até o `expiry` | 6.4.2; 5.2; 4.4.0 | `af/af/internal/server/keys.go` | [R3.ua-old]; [N22], [N23]; `TestKafExpiry*` | **Verificado** |
| Remoção do contexto AKMA | 6.6 | `free5gc/aanf/` | [R6.removed], [R6.not-found], [C1] | **Verificado** |
| Aplicações do UE sem acesso à K_AKMA | 4.4.3 | `0003-…` (`akma-kaf`) | `docs/af-side/teste-ue.md` §4 | **Verificado** |
| Nenhuma chave em log | plano, todos os MVPs | todos | [R7.no-key-*], [L1]; `TestNoKeyInLogs` | **Verificado** |

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

## Lado UE (ME)

| Requisito | Norma | Cláusula | Implementação | Teste | Estado |
|---|---|---|---|---|---|
| ME deriva K_AKMA e A-KID da K_AUSF antes de comunicar com um AF | TS 33.535 V19.0.0 | 6.1; 6.2.1 passo 1 | `af/ue/ueransim/patches/0001-akma-key-derivation.patch` (`NasMm::deriveAkmaKeys`) | VM: A-KID do UE igual ao do AUSF e do AAnF | **Implementado** (MVP-3a, D14); instante da derivação `[INTERPRETAÇÃO]` (`docs/af-side/ue-ueransim.md` §2) |
| K_AKMA (FC `0x80`) e A-TID (FC `0x81`) no ME; SUPI = IMSI | TS 33.535 V19.0.0; TS 33.501 V19.7.0 | A.2, A.3; A.7.0 | `af/ue/ueransim/patches/0001-akma-key-derivation.patch` (`keys::CalculateKAkma`, `keys::CalculateATid`) | `nr-akma-selftest` (8 vetores do MVP-0); VM (A-KID igual ao da rede) | **Implementado** (MVP-3a; conferido no MVP-4) |
| A-KID no ME: RID do SUCI, A-TID em base64, domínio com MNC de três dígitos | TS 33.535 V19.0.0; TS 29.522 V19.7.0; TS 23.003 V19.7.0 | 6.1; Tab. 5.14.5.4.2-1; 2.2B, 28.2 | idem (`keys::ConstructAKid`) | VM, RID `"0000"` | **Implementado** (MVP-3a) |
| Opção `akma` na configuração do UE | — | — | idem; `0002-free5gc-ue-akma.patch` | VM: sem a opção, nenhuma derivação | **Implementado** (MVP-3a) — decisão do plano |
| K_AF no ME, sob demanda, para um AF_ID | TS 33.535 V19.0.0; TS 33.220 V19.1.0 | 6.2.1 passo 1; A.4; Anexo H | `0003-akma-kaf-nr-cli-selftest.patch` (`keys::CalculateKAf`) | `nr-akma-selftest`; Ua\* com a K_AF do UE (`docs/af-side/teste-ue.md`) | **Implementado** (MVP-4) |
| Aplicações do UE sem acesso à K_AKMA | TS 33.535 V19.0.0 | 4.4.3 | comando `akma-kaf` do `nr-cli`: só A-KID e K_AF | `teste-ue.md` §4 | **Implementado** (MVP-4) — `nr-cli` é interface de teste: entrega K_AF a qualquer processo local |
| FQDN do AF em UTF-8 com NFKC na K_AF | TS 33.220 V19.1.0 | B.2.1.2 | ME recusa FQDN não ASCII | `nr-akma-selftest` (vetores de NFKC) | **Implementado** (MVP-4) — DV-4 |
| Cliente Ua\* do UE pelo plano de usuário | TS 33.535 V19.0.0 | 6.2.1 passo 1 | `af/ue/uaclient/` (`-nr-cli`, `-bind`), `af/scripts/ue-netns.sh` | `teste-ue.md` §2 e §3 | **Implementado** (MVP-4) — D5 |

## Logs

| Requisito | Origem | Implementação | Teste | Estado |
|---|---|---|---|---|
| Nenhuma chave em log, em nível algum | plano, critério de todos os MVPs | UDM, AUSF, AAnF, UERANSIM | `TestNoKeyInLogs` (AAnF); VM: nenhuma cadeia de 64 hexadecimais nos logs e na tela do UE | **Verificado** em 18/10/2026 |
| Eventos AKMA em `info`, categoria `AKMA`, com SUPI, A-KID e AF_ID | D16 | `free5gc/patches/{udm,ausf}/`, `free5gc/aanf/` | VM (`docs/free5gc-side/verificacao.md` §4) | **Implementado** (MVP-3a) |
| AAnF na tela do `run.sh` (`-aanf`) | D15 | `free5gc/patches/run/0001-aanf-option.patch` | VM | **Implementado** (MVP-3a) |

## Ponto de referência Ua*

| Requisito | Norma | Cláusula | Implementação | Teste | Estado |
|---|---|---|---|---|---|
| Ua* simplificado com id de faixa privada | TS 33.220 H.2 NOTA 2 | — | `af/uastar/`, `af/af/`, `af/ue/uaclient/`; `docs/af-side/ua-star.md` | `docs/af-side/teste-af.md` | **Implementado** (MVP-3) — ADR-0003 |
| Perfil HTTP Digest — username, senha, realm | TS 24.109 V19.0.0; TS 33.535 V19.0.0; RFC 7616 | K.2; B.1.2.2; 3.3–3.5 | `af/uastar/` | `TestDigestVectors` (RFC 7616 §3.9.1 e Python independente); integração | **Implementado** (MVP-3) — DV-5 (RFC 7616, SHA-256), DV-6 (`realm = "3GPP-bootstrapping-akma@<FQDN>"`) |
| Confirmação de chave nos dois sentidos (`rspauth`) | TS 33.535 V19.0.0; RFC 7616 | 4.4.1; 3.5 | `af/af/internal/server/`, `af/ue/uaclient/` | `TestAFWithoutKaf`, `TestWrongKafRejected`; integração | **Implementado** (MVP-3) |
| *Product token* `3gpp-akma`; o UE verifica o FQDN no `realm` | TS 24.109 V19.0.0; TS 33.535 V19.0.0 | K.2 b) 1), d); B.1.2.2 passos 2 e 4 | `af/ue/uaclient/` | `TestChallengeChecks`, `TestHostAndProductToken` | **Implementado** (MVP-3) |
| FQDN do AF levado no `Host` e conferido pelo AF | TS 33.220 V19.1.0; TS 33.222 V19.0.0 | 4.5.2 pré-requisito (3); 5.3.0 passo 2 | `af/af/internal/server/` | `TestHostPathMethod`; integração (421) | **Implementado** (MVP-3) — `[INTERPRETAÇÃO]` 421 para outro nome |
| AF pede K_AF ao AAnF sem contexto ativo; `afId` e `aKId` | TS 33.535 V19.0.0; TS 29.535 V19.5.0 | 6.2.1 passos 2 e 9; 4.2.2.3 | `af/af/internal/aanf/` | `TestApplicationKey`; integração | **Implementado** (MVP-3) |
| Erro sem o `realm` do desafio quando o AF não obtém K_AF | TS 33.535 V19.0.0 | B.1.2.2 passo 6 | `af/af/internal/server/` | `TestKeyNotAvailable`, `TestAanfUnavailable` | **Implementado** (MVP-3) |
| Iniciação do AKMA pelo AF | TS 33.535 V19.0.0 | 6.5 | `af/af/internal/server/` (o `401` com desafio) | `TestInitiation`; integração | **Implementado** (MVP-3) |
| Expiração da K_AF no AF: nunca usada após `expiry`; renovar no AAnF ou recusar | TS 33.535 V19.0.0 | 4.4.0, 4.4.1, 6.4.2 | `af/af/internal/server/keys.go` (`kafExpiry`) | `TestKafExpiryRefetch`, `TestKafExpiryReject`, `TestExpiredOnArrival` | **Implementado** (MVP-3) |
| Perfil TLS 1.2 PSK | TS 24.109 V19.0.0 | K.3.2 | `af/` | ponta a ponta | Pendente |
| Perfil TLS 1.3 PSK — `"3GPP-AKMA" ";" A-KID` | TS 24.109 V19.0.0 | K.3.3 | `af/` | ponta a ponta | Pendente — defeito D-E |
| Ausência de renovação de K_AF | TS 33.535 §6.4.3; TS 24.109 K.2, K.3 | — | — | — | Especificado |

## Requisitos de ambiente

| Requisito | Origem | Verificação | Estado |
|---|---|---|---|
| Sincronização NTP entre a VM do AF e a do free5GC (par AAnF↔AF) | ADR-0001-R2 §5.3 | `timedatectl` nas duas VMs | Especificado |
| AF alcança a porta do AAnF | ADR-0001-R2 §5.1, §5.2 | requisição a partir de `192.168.16.129` | **Verificado** em 17/10/2026 (MVP-1) |
| UE alcança o AF pelo plano de usuário (`uesimtun0` → UPF → N6) | D5; `docs/free5gc-side/plano-de-usuario.md` | `ping` e Ua\* a partir do `uesimtun0` no *namespace* do UE; contadores da `upfgtp` | **Verificado** em 18/10/2026 (MVP-4): UE no próprio *namespace*, NAT no N6 |
| Proteção das interfaces SBI | TS 33.535 §4.4.0; TS 33.501 §13.1.0 | TLS habilitado | Decisão — DV-1; TLS no MVP-6+ |
