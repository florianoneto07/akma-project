# Registro de desvios e escolhas entre textos normativos

**Criado em:** 16 de outubro de 2026
**Regra:** todo ponto em que o projeto não segue literalmente a norma — por limitação do free5GC, simplificação de MVP, ambiguidade ou contradição do texto — é apresentado com a cláusula, o motivo e o impacto, decidido pelo responsável pelo projeto e registrado aqui antes de ser implementado.

Estados: **decidido** (escolha explícita do responsável), **aprovado com o plano** (recomendação aprovada junto com o plano de desenvolvimento em 16/10/2026), **implementado** (com referência ao código e ao teste).

| Nº | Assunto | Estado | MVP |
|---|---|---|---|
| DV-1 | SBI sem TLS (h2c) | **implementado** (MVP-1, 17/10/2026) | MVP-1 |
| DV-2 | Subscrição de `ROAMING_STATUS` desabilitada | **implementado** (MVP-1, 17/10/2026) — parcial: a etapa não existe; o campo só aceita `false` | MVP-1 |
| DV-3 | Só SUPI do tipo IMSI | **implementado** (MVP-0, 16/10/2026) | MVP-0 |
| DV-4 | Sem NFKC no lado UE | aprovado com o plano | MVP-4 |
| DV-5 | Digest do Ua\* conforme RFC 7616 | aprovado com o plano | MVP-3 |
| DV-6 | `realm` do Digest conforme TS 33.535 | **decidido** em 16/10/2026 | MVP-3 |
| DV-7 | Nunca emitir 204 em `retrieve-applicationkey` | **implementado** (MVP-1, 17/10/2026) | MVP-1 |
| DV-8 | `aud` vazio tolerado no token OAuth2 | **encerrado** em 17/10/2026 (MVP-1b: NRF corrigido, AAnF estrito) | MVP-1 |

---

## DV-1 — SBI sem TLS

**Norma.** `[3GPP]` TS 33.535 §4.4.0: as interfaces AAnF↔AUSF, AAnF↔AF/NEF e AAnF↔UDM *"shall be confidentiality, integrity and replay protected"*. TS 33.535 §4.2.2 e TS 29.535 §4.2.1.3.2: o AF *"shall be authenticated and authorized by the operator network"* antes de receber K_AF. TS 33.501 §13.1.0: *"TLS shall be used for transport protection within a PLMN unless network security is provided by other means"*; NOTA 2: com interfaces confiáveis, a proteção criptográfica é decisão do operador.

**O que o projeto faz.** `[MVP]` `http` (h2c) em AUSF↔AAnF (*loopback*, mesma máquina) e em AF↔AAnF (rede virtual entre as duas VMs), declarando ambas interfaces confiáveis do laboratório. A política local do AAnF associa cada AF_ID autorizado ao endereço de origem do AF.

**Por quê.** As NFs do free5GC operam em `http` na configuração de laboratório; habilitar TLS exige certificados com SAN coerente com o endereço publicado no NRF (ADR-0001-R2 §5.5). Fica para o MVP-6+.

**Impacto.** O AF não é autenticado pelo transporte; um AF que conheça um A-KID pode apresentar o AF_ID de outro AF (ataque de falsificação de AF descrito por Akman et al., 2023 — `docs/literatura/analise-dos-artigos.md`). A associação AF_ID↔endereço reduz, mas não elimina, o risco. O suporte a TLS continua obrigatório para a NF (TS 33.501 §13.1.0).

**Implementação (17/10/2026).** O AAnF serve h2c (`sbi.scheme: http`) em `192.168.16.130:8000`, alcançável da VM do AF. A política local está em `akma.authorizedAfs`: cada AF_ID tem a lista de endereços de origem autorizados, e o endereço é o da conexão, não o de `X-Forwarded-For` (`free5gc/aanf/internal/context/context.go`, `AfPolicy`). Testes: `TestRetrieveAfNotAuthorized` e, na VM, AF_ID autorizado pedido de outra máquina → 403 (`docs/free5gc-side/aanf.md` §9). O suporte a `https` existe na configuração (`sbi.tls`), mas não foi exercitado.

---

## DV-2 — Subscrição de `ROAMING_STATUS` desabilitada

**Norma.** `[3GPP]` TS 33.535 §6.2.1 passo 5: *"the AAnF shall send a Nudm_EventExposure_Subscribe request to UDM … to request the RoamingStatusReport"*. `[3GPP]` TS 29.535 condiciona a recusa por *roaming* à *feature* opcional `RoamingRestriction` (Tab. 5.1.8-1, Tab. 5.1.7.3-1).

**O que o projeto faz.** `[MVP]` A subscrição existe como etapa configurável (`akma.roamingStatusSubscription`), desabilitada por padrão.

**Por quê.** Questão em aberto nº 1: o `shall` do Stage 2 convive com a *feature* opcional do Stage 3. Orientação do material de partida do projeto.

**Impacto.** Sem subscrição, o AAnF não detecta UE em *roaming* nem recusa por política de *roaming*; o testbed é sem *roaming*. `ROAMING_AKMA_SERVICE_DENIED` nunca é emitido.

**Implementação (17/10/2026).** `akma.roamingStatusSubscription` existe no `aanfcfg.yaml`, mas o valor `true` impede a partida do AAnF com mensagem que cita este desvio (`free5gc/aanf/pkg/factory/config.go`; teste `TestInvalidConfigs`). A etapa configurável prevista acima fica para quando o AAnF consumir o UDM. A *feature* `RoamingRestriction` não é anunciada, e o `notifUri` dos pedidos é ignorado. Todo teste de API confere que `ROAMING_AKMA_SERVICE_DENIED` não aparece em nenhuma resposta (`api_akma_test.go`, função `do`).

---

## DV-3 — Só SUPI do tipo IMSI

**Norma.** `[3GPP]` TS 33.501 A.7.0 define o P0 para SUPI do tipo IMSI, NAI, GLI e GCI; TS 33.535 A.2 e A.3 remetem a esse P0.

**O que o projeto faz.** `[MVP]` Aceita apenas SUPI `imsi-…`; os demais tipos são recusados com erro explícito na camada da KDF.

**Por quê.** O free5GC e o UERANSIM do testbed usam IMSI.

**Impacto.** Assinantes com SUPI NAI/GLI/GCI não obtêm AKMA.

**Implementação (16/10/2026).** `shared/akma/supi/supi.go`: `FromSBI` devolve `ErrUnsupportedType` para `nai-`, `gci-` e `gli-`; `Param.Validate` recusa qualquer valor fora de `^[0-9]{5,15}$`, inclusive com o prefixo `imsi-`. Testes: `supi_test.go` (`TestUnsupportedTypes`, `TestPrefixRejectedAtKDFLayer`) e `derive_test.go` (`TestSBIPrefixRejected`).

---

## DV-4 — Sem NFKC no lado UE

**Norma.** `[3GPP]` TS 33.220 B.2.1.2: cadeias em UTF-8 com NFKC.

**O que o projeto faz.** `[MVP]` O ME (*patch* no UERANSIM, em C++) não aplica NFKC e recusa FQDN não ASCII. AAnF e biblioteca Go aplicam NFKC.

**Por quê.** O UERANSIM não tem biblioteca de normalização Unicode.

**Impacto.** Nenhum para FQDN ASCII (invariante sob NFKC); AF com FQDN internacionalizado não é suportado no lado UE.

---

## DV-5 — Digest do Ua\* conforme RFC 7616

**Norma.** `[DIVERGÊNCIA]` TS 24.109 K.2 c): *"HTTP Digest authentication (see RFC 3310 [6]) shall be used"*. O mesmo K.2 b) remete à §5.2, cuja §5.2.2.1 usa RFC 7616; TS 33.535 B.1.2.2 segue TS 33.222 §5.3.0, que usa RFC 7235 e RFC 7616 e desaconselha MD5 (NOTA 1b).

**O que o projeto faz.** Digest conforme RFC 7616, `algorithm=SHA-256`, `qop=auth` (`auth-int` opcional, TS 24.109 §5.2.6), `rspauth` obrigatório.

**Por quê.** RFC 3310 é o Digest AKA, em que a senha vem do RES; com senha = base64(K_AF), só RFC 7616 é coerente.

**Impacto.** Diverge da referência literal de K.2 c). Candidato a defeito (CT1).

---

## DV-6 — `realm` do Digest conforme TS 33.535

**Norma.** `[DIVERGÊNCIA]` TS 33.535 B.1.2.2 (Stage 2, normativo), que adapta TS 33.222 §5.3.0: no passo 3, *"the AF shall include the "3GPP-bootstrapping-akma" within the WWW-Authenticate header field"*; no passo 4, o cliente verifica *"the FQDN in the realm attribute"*. TS 24.109 K.2 b)2): *"The realm attribute shall contain the constant string "3GPP-bootstrapping-akma""*. TS 24.109 §5.2.2.1 item 3 (o modelo GBA que K.2 adapta): `realm` em duas partes separadas por `@`, constante e FQDN. Em sentido contrário, TS 24.109 K.2 c)3): primeira parte `"3gpp-akma"`, exemplo `3gpp-akma@af1.operator.com`.

**O que o projeto faz.** Segue TS 33.535: `realm = "3GPP-bootstrapping-akma@<FQDN do AF>"`. O conteúdo — a constante e o FQDN — vem de TS 33.535 B.1.2.2; a forma `<constante>@<FQDN>` vem do Stage 3 (TS 24.109 §5.2.2.1 e K.2 b)2)), porque o Stage 2 diz o que o `realm` contém, não como as partes se separam. A primeira parte fica configurável.

**Por quê.** Três fontes (TS 33.535 B.1.2.2, TS 24.109 K.2 b)2) e a estrutura do GBA) contra uma (K.2 c)3)). A decisão anterior do projeto (registro de defeitos D-C, Adendo D §D.4.1) seguia c)3) com argumento que as duas leituras satisfazem (`docs/baseline/conferencia-normativa.md` §3.10).

**Impacto.** Diverge apenas de TS 24.109 K.2 c)3). Um UE que siga c)3) e verifique a primeira parte do `realm` não interopera — no MVP o cliente é do próprio projeto. O registro de defeitos (D-C) e o Adendo D precisam ser atualizados.

---

## DV-7 — Nunca emitir 204 em `retrieve-applicationkey`

**Norma.** `[DIVERGÊNCIA]` TS 29.535 §4.2.2.3.2 e Tab. 5.1.4.3.2-2: *"If the requested AKMA Application Key information for the UE does not exist, the AAnF shall respond with "204 No Content""*. A mesma cláusula: sem K_AKMA, 403 `K_AKMA_NOT_PRESENT`. TS 33.535 §6.2.1 passo 7: com K_AKMA presente, o AAnF deriva K_AF.

**O que o projeto faz.** O AAnF nunca responde 204; sem K_AKMA, responde 403 `K_AKMA_NOT_PRESENT`.

**Por quê.** Não há condição distinta que leve ao 204: sem K_AKMA a resposta é 403, e com K_AKMA a K_AF é sempre derivável.

**Impacto.** Nenhum consumidor conforme deixa de funcionar (o 204 continua previsto no cliente, que o trata). Candidato a defeito (CT3).

**Implementação (17/10/2026).** `free5gc/aanf/internal/sbi/processor/retrieve_applicationkey.go`: sem K_AKMA, 403 `K_AKMA_NOT_PRESENT`; com K_AKMA, 200. Todo pedido a `retrieve-applicationkey` nos testes de API confere que a resposta não é 204. O cliente do *patch* do `openapi` trata o 204 (`NoContent`).

---

## DV-8 — `aud` vazio tolerado no token OAuth2

**Norma.** `[3GPP]` TS 33.501 §13.4.1.1.2 passo 2: verificada a integridade, o produtor *"checks that the audience claim in the access token matches its own identity or the NF type of NF Service Producer"* e *"checks that the access token has not expired by verifying the expiration time"*. Na mesma cláusula, as *claims* do token emitido pelo NRF incluem *"NF type of the NF Service Producer (audience)"*.

**O que o free5GC faz.** `[FREE5GC]` `oauth.VerifyOAuth` (`free5gc/openapi` v1.2.4) verifica assinatura e escopo, mas não verifica `exp` nem `aud`: o campo `exp` é próprio do modelo e não chega à validação do `jwt`. O NRF v1.4.5 preenche `aud` com o `targetNfInstanceId` do pedido de token, que os clientes do free5GC não enviam; o `aud` sai vazio.

**O que o projeto faz.** O AAnF verifica `exp` (conforme) e aceita `aud` igual ao seu `nfInstanceId` ou a `AANF`. **Tolera `aud` vazio**; um `aud` não vazio diferente desses é recusado.

**Por quê.** Exigir `aud` recusaria todo token do NRF atual, e `register-anchorkey` e `remove-context` só funcionariam depois de um *patch* no NRF. Decidido pelo responsável em 17/10/2026, entre três opções: esta, a estrita e a do free5GC (sem `exp` nem `aud`).

**Impacto.** Um token emitido para outro produtor, com `aud` vazio e escopo `naanf-akma`, é aceito. Como o NRF só emite escopo `naanf-akma` para os tipos em `allowedNfTypes` do AAnF (AUSF, NEF, AF), o risco fica nesses consumidores. Corrigido quando o *patch* futuro no NRF preencher `aud` (`docs/free5gc-side/registro-nrf.md` §5, N3); então a tolerância sai.

**Implementação (17/10/2026).** `free5gc/aanf/internal/context/context.go` (`VerifyAccessToken`, `audienceAccepted`). Testes: `TestOAuth2` (expirado, sem `exp`, `aud` de outra instância e de outro tipo recusados; `aud` vazio, igual ao ID, igual a `AANF` e em lista aceitos).

**Encerramento (17/10/2026).** A decisão D13 foi revista no mesmo dia: o NRF foi corrigido (MVP-1b, *patch* `free5gc/patches/nrf/0003-access-token-audience.patch`, `docs/free5gc-side/nrf.md`) e passou a preencher `aud` com o tipo de NF do produtor (TS 29.510 Tab. 6.3.5.2.4-1). O AAnF deixou de tolerar `aud` vazio ou ausente. Testes: `TestOAuth2` (`aud` vazio, ausente e lista vazia recusados) e, na VM, `register-anchorkey` aceito com token do NRF corrigido. O desvio deixa de existir; o registro fica como histórico.
