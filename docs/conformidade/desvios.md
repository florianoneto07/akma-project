# Registro de desvios e escolhas entre textos normativos

**Criado em:** 16 de outubro de 2026
**Regra:** todo ponto em que o projeto não segue literalmente a norma — por limitação do free5GC, simplificação de MVP, ambiguidade ou contradição do texto — é apresentado com a cláusula, o motivo e o impacto, decidido pelo responsável pelo projeto e registrado aqui antes de ser implementado.

Estados: **decidido** (escolha explícita do responsável), **aprovado com o plano** (recomendação aprovada junto com o plano de desenvolvimento em 16/10/2026), **implementado** (com referência ao código e ao teste).

| Nº | Assunto | Estado | MVP |
|---|---|---|---|
| DV-1 | SBI sem TLS (h2c) | aprovado com o plano | MVP-1 |
| DV-2 | Subscrição de `ROAMING_STATUS` desabilitada | aprovado com o plano (orientação do material de partida) | MVP-1 |
| DV-3 | Só SUPI do tipo IMSI | **implementado** (MVP-0, 16/10/2026) | MVP-0 |
| DV-4 | Sem NFKC no lado UE | aprovado com o plano | MVP-4 |
| DV-5 | Digest do Ua\* conforme RFC 7616 | aprovado com o plano | MVP-3 |
| DV-6 | `realm` do Digest conforme TS 33.535 | **decidido** em 16/10/2026 | MVP-3 |
| DV-7 | Nunca emitir 204 em `retrieve-applicationkey` | **decidido** em 16/10/2026 | MVP-1 |

---

## DV-1 — SBI sem TLS

**Norma.** `[3GPP]` TS 33.535 §4.4.0: as interfaces AAnF↔AUSF, AAnF↔AF/NEF e AAnF↔UDM *"shall be confidentiality, integrity and replay protected"*. TS 33.535 §4.2.2 e TS 29.535 §4.2.1.3.2: o AF *"shall be authenticated and authorized by the operator network"* antes de receber K_AF. TS 33.501 §13.1.0: *"TLS shall be used for transport protection within a PLMN unless network security is provided by other means"*; NOTA 2: com interfaces confiáveis, a proteção criptográfica é decisão do operador.

**O que o projeto faz.** `[MVP]` `http` (h2c) em AUSF↔AAnF (*loopback*, mesma máquina) e em AF↔AAnF (rede virtual entre as duas VMs), declarando ambas interfaces confiáveis do laboratório. A política local do AAnF associa cada AF_ID autorizado ao endereço de origem do AF.

**Por quê.** As NFs do free5GC operam em `http` na configuração de laboratório; habilitar TLS exige certificados com SAN coerente com o endereço publicado no NRF (ADR-0001-R2 §5.5). Fica para o MVP-6+.

**Impacto.** O AF não é autenticado pelo transporte; um AF que conheça um A-KID pode apresentar o AF_ID de outro AF (ataque de falsificação de AF descrito por Akman et al., 2023 — `docs/literatura/analise-dos-artigos.md`). A associação AF_ID↔endereço reduz, mas não elimina, o risco. O suporte a TLS continua obrigatório para a NF (TS 33.501 §13.1.0).

---

## DV-2 — Subscrição de `ROAMING_STATUS` desabilitada

**Norma.** `[3GPP]` TS 33.535 §6.2.1 passo 5: *"the AAnF shall send a Nudm_EventExposure_Subscribe request to UDM … to request the RoamingStatusReport"*. `[3GPP]` TS 29.535 condiciona a recusa por *roaming* à *feature* opcional `RoamingRestriction` (Tab. 5.1.8-1, Tab. 5.1.7.3-1).

**O que o projeto faz.** `[MVP]` A subscrição existe como etapa configurável (`akma.roamingStatusSubscription`), desabilitada por padrão.

**Por quê.** Questão em aberto nº 1: o `shall` do Stage 2 convive com a *feature* opcional do Stage 3. Orientação do material de partida do projeto.

**Impacto.** Sem subscrição, o AAnF não detecta UE em *roaming* nem recusa por política de *roaming*; o testbed é sem *roaming*. `ROAMING_AKMA_SERVICE_DENIED` nunca é emitido.

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
