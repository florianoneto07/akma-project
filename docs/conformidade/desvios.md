# Registro de desvios e escolhas entre textos normativos

**Criado em:** 16 de outubro de 2026
**Regra:** todo ponto em que o projeto não segue literalmente a norma — por limitação do free5GC, simplificação de MVP, ambiguidade ou contradição do texto — é apresentado com a cláusula, o motivo e o impacto, decidido pelo responsável pelo projeto e registrado aqui antes de ser implementado.

Estados: **decidido** (escolha explícita do responsável), **aprovado com o plano** (recomendação aprovada junto com o plano de desenvolvimento em 16/10/2026), **implementado** (com referência ao código e ao teste).

| Nº | Assunto | Estado | MVP |
|---|---|---|---|
| DV-1 | SBI sem TLS (h2c); CAPIF-3/4 sem TLS até o MVP-11 | **implementado** (MVP-1, 17/10/2026); estendido a CAPIF-3/4, **implementado** (MVP-9, 02/11/2026); **a encerrar** no MVP-11, aprovado em 02/11/2026 | MVP-1, MVP-9, MVP-11 |
| DV-2 | Subscrição de `ROAMING_STATUS` desabilitada | **implementado** (MVP-1, 17/10/2026) — parcial: a etapa não existe; o campo só aceita `false` | MVP-1 |
| DV-3 | Só SUPI do tipo IMSI | **implementado** (MVP-0, 16/10/2026) | MVP-0 |
| DV-4 | Sem NFKC no lado UE | **implementado** (MVP-4, 18/10/2026) | MVP-4 |
| DV-5 | Digest do Ua\* conforme RFC 7616 | **implementado** (MVP-3, 18/10/2026) | MVP-3 |
| DV-6 | `realm` do Digest conforme TS 33.535 | **implementado** (MVP-3, 18/10/2026) | MVP-3 |
| DV-7 | Nunca emitir 204 em `retrieve-applicationkey` | **implementado** (MVP-1, 17/10/2026) | MVP-1 |
| DV-8 | `aud` vazio tolerado no token OAuth2 | **encerrado** em 17/10/2026 (MVP-1b: NRF corrigido, AAnF estrito) | MVP-1 |
| DV-9 | NEF–AF externo sem TLS nem certificado de cliente | **implementado** (MVP-6, 19/10/2026); vale também para CAPIF-2e (MVP-9, 02/11/2026); **a encerrar** no MVP-11 (T14, opção b, 02/11/2026) | MVP-6, MVP-9, MVP-11 |
| DV-10 | Perfis TLS do Ua\* sem `TLS_ECDHE_PSK_WITH_AES_256_GCM_SHA384` e sem OCSP *stapling*, recomendados por TS 33.210 | **implementado** (MVP-8, 21/10/2026) | MVP-8 |
| DV-11 | CoAP sobre DTLS sem a suíte obrigatória de RFC 7252, `TLS_PSK_WITH_AES_128_CCM_8`; vale TS 33.210 | **implementado** (MVP-8, 21/10/2026) | MVP-8 |
| DV-12 | `trigger-auth` com POST, contra a precedência do YAML (TS 29.503 Anexo A.1), que declara `get` | **implementado** (MVP-12, 28/10/2026) | MVP-12 |
| DV-13 | *Claims* do *token* do CCF: `iss` com o identificador do *invoker*, ao lado de `client_id`, e `exp` como instante | **implementado** (MVP-9, 02/11/2026) | MVP-9 |
| DV-14 | CAPIF parcial: só o caminho do AKMA | **implementado** (MVP-9, 02/11/2026) | MVP-9 |
| DV-15 | HTTP/2 por ALPN no CCF, contra o *upgrade* do HTTP de TS 29.222 §7.3 | **implementado** (MVP-9, 02/11/2026) | MVP-9 |
| DV-16 | CAPIF-1e sem revogação de certificados (sem ponto de distribuição de CRL, CRL nem OCSP) e *token* do CCF sem `x5t#S256` | **implementado** (MVP-9, 02/11/2026); revisto na especificação do MVP-11 (04/11/2026): a revogação sai do desvio na fase (c), e ficam o OCSP e o `x5t#S256` | MVP-9, MVP-11 |
| DV-17 | SBI por IP, sem o FQDN no perfil nem o SNI, como modo opcional; o padrão segue a norma | **decidido** (MVP-11, 02/11/2026); a implementar | MVP-11 |
| DV-18 | CRL só por HTTP, sem o LDAPv3 que TS 33.310 §6.1a pede como método principal | **decidido** (MVP-11, 04/11/2026); a implementar | MVP-11 |

**Consolidação (MVP-5, 18/10/2026).** Nenhum desvio novo surgiu no MVP-5. Estado e evidência de cada um no roteiro ponta a ponta (`docs/af-side/e2e.md`, execução `20261018-171510`):

| Nº | Estado final | Evidência no laboratório |
|---|---|---|
| DV-1 | implementado | SBI em h2c e HTTP/1.1 sem TLS; o AF fora da política é recusado pelo endereço de origem [N03] |
| DV-2 | implementado, parcial | nenhuma subscrição de *roaming*; o fluxo completo funciona sem ela [R2], [R3] |
| DV-3 | implementado | só SUPI IMSI no AUSF, no AAnF e no ME |
| DV-4 | implementado | o ME recusa FQDN não ASCII na K_AF [S1] |
| DV-5 | implementado | Digest SHA-256 com `rspauth` em todo Ua\* [R2.ua] |
| DV-6 | implementado | `realm="3GPP-bootstrapping-akma@af1.akma.example"` no desafio do AF |
| DV-7 | implementado | nenhuma resposta 204 nos testes da API [N01] a [N16] |
| DV-8 | encerrado | `aud` exigido pelo AAnF com o NRF corrigido; sem *token*, 401 [N15], [N16] |

---

## DV-1 — SBI sem TLS

**Norma.** `[3GPP]` TS 33.535 §4.4.0: as interfaces AAnF↔AUSF, AAnF↔AF/NEF e AAnF↔UDM *"shall be confidentiality, integrity and replay protected"*. TS 33.535 §4.2.2 e TS 29.535 §4.2.1.3.2: o AF *"shall be authenticated and authorized by the operator network"* antes de receber K_AF. TS 33.501 §13.1.0: *"TLS shall be used for transport protection within a PLMN unless network security is provided by other means"*; NOTA 2: com interfaces confiáveis, a proteção criptográfica é decisão do operador.

**O que o projeto faz.** `[MVP]` `http` (h2c) em AUSF↔AAnF (*loopback*, mesma máquina) e em AF↔AAnF (rede virtual entre as duas VMs), declarando ambas interfaces confiáveis do laboratório. A política local do AAnF associa cada AF_ID autorizado ao endereço de origem do AF.

**Por quê.** As NFs do free5GC operam em `http` na configuração de laboratório; habilitar TLS exige certificados com SAN coerente com o endereço publicado no NRF (ADR-0001-R2 §5.5). Fica para o MVP-6+.

**Impacto.** O AF não é autenticado pelo transporte; um AF que conheça um A-KID pode apresentar o AF_ID de outro AF (ataque de falsificação de AF descrito por Akman et al., 2023 — `docs/literatura/analise-dos-artigos.md`). A associação AF_ID↔endereço reduz, mas não elimina, o risco. O suporte a TLS continua obrigatório para a NF (TS 33.501 §13.1.0).

**Implementação (17/10/2026).** O AAnF serve h2c (`sbi.scheme: http`) em `192.168.16.130:8000` (8010 desde o MVP-6), alcançável da VM do AF. A política local está em `akma.authorizedAfs`: cada AF_ID tem a lista de endereços de origem autorizados, e o endereço é o da conexão, não o de `X-Forwarded-For` (`free5gc/aanf/internal/context/context.go`, `AfPolicy`). Testes: `TestRetrieveAfNotAuthorized` e, na VM, AF_ID autorizado pedido de outra máquina → 403 (`docs/free5gc-side/aanf.md` §9). O suporte a `https` existe na configuração (`sbi.tls`), mas não foi exercitado.

**Extensão (29/10/2026, MVP-9).** Decisão do responsável no ponto C3 do MVP-9 (plano, seção 4): CAPIF-3/4, entre o NEF e o CCF, fica em `http` num endereço de *loopback* do CCF até o MVP-11, que traz o TLS mútuo com os certificados das NFs. `[3GPP]` TS 33.122 V19.4.0 §6.6: *"TLS shall be used to provide integrity protection, replay protection and confidentiality protection. The support of TLS is mandatory."* e *"Certificate based mutual authentication shall be performed between the CAPIF entities using TLS."*; NOTA: *"It is up to the domain administrator's policy to protect interfaces within the trusted domain."* Como as NFs deste desvio, o CCF reconhece o NEF pela interface interna. CAPIF-1e, entre o AF e o CCF, tem TLS desde o MVP-9.

**Implementação da extensão (02/11/2026).** O NEF fala com o servidor interno do CCF em `http://127.0.0.1:8020` (`free5gc/patches/config/0003`). Nesse endereço, o CCF serve só a publicação e os dados do *invoker* (`docs/free5gc-side/capif.md` §3). Evidência: [R1.nef-published] e [CP04] na execução `20261102-093804` (`docs/af-side/e2e.md` §5.9). No [CP04], o `check-authentication` faz o NEF pedir ao CCF os dados do *invoker*.

**MVP-11 (02/11/2026).** Achado na proposta do MVP-11: `[3GPP]` TS 33.501 V19.7.0 §13.3.2.1: *"If the PLMN uses token-based authorization, the network shall use protection at the transport layer as described in clause 13.1."* O laboratório usa *tokens* do NRF e, desde o MVP-9, do CCF: a NOTA 2 de §13.1.0, base deste desvio, não cobre esse caso. O MVP-11, aprovado em 02/11/2026, encerra o desvio: TLS mútuo em toda a SBI (pontos T2 a T8), em CAPIF-3/4 (T15) e a política do AAnF pelo certificado, no lugar do endereço (T13).

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

**Implementação (18/10/2026).** `af/ue/ueransim/patches/0003-akma-kaf-nr-cli-selftest.patch`: `keys::IsValidAfFqdn` aceita só nome DNS ASCII. `CalculateKAf` e o comando `akma-kaf` do `nr-cli` recusam qualquer outro. O `nr-akma-selftest` confere a recusa nos três vetores de NFKC do MVP-0 (`docs/af-side/teste-ue.md` §1).

---

## DV-5 — Digest do Ua\* conforme RFC 7616

**Norma.** `[DIVERGÊNCIA]` TS 24.109 K.2 c): *"HTTP Digest authentication (see RFC 3310 [6]) shall be used"*. O mesmo K.2 b) remete à §5.2, cuja §5.2.2.1 usa RFC 7616; TS 33.535 B.1.2.2 segue TS 33.222 §5.3.0, que usa RFC 7235 e RFC 7616 e desaconselha MD5 (NOTA 1b).

**O que o projeto faz.** Digest conforme RFC 7616, `algorithm=SHA-256`, `qop=auth` (`auth-int` opcional, TS 24.109 §5.2.6), `rspauth` obrigatório.

**Por quê.** RFC 3310 é o Digest AKA, em que a senha vem do RES; com senha = base64(K_AF), só RFC 7616 é coerente.

**Impacto.** Diverge da referência literal de K.2 c). Candidato a defeito (CT1).

**Implementação (18/10/2026).** `af/uastar/` (`Input.Response`, `Input.RspAuth`): SHA-256, `qop=auth` e `auth-int`, conferidos contra o exemplo de RFC 7616 §3.9.1 e contra os vetores de uma implementação Python independente (`af/uastar/testdata/gen_digest_vectors.py`). O AF recusa `algorithm=MD5` com 400 (`TestRequestChecks`); o cliente recusa resposta sem `rspauth` válido (`TestAFWithoutKaf`). Especificação: `docs/af-side/ua-star.md`.

---

## DV-6 — `realm` do Digest conforme TS 33.535

**Norma.** `[DIVERGÊNCIA]` TS 33.535 B.1.2.2 (Stage 2, normativo), que adapta TS 33.222 §5.3.0: no passo 3, *"the AF shall include the "3GPP-bootstrapping-akma" within the WWW-Authenticate header field"*; no passo 4, o cliente verifica *"the FQDN in the realm attribute"*. TS 24.109 K.2 b)2): *"The realm attribute shall contain the constant string "3GPP-bootstrapping-akma""*. TS 24.109 §5.2.2.1 item 3 (o modelo GBA que K.2 adapta): `realm` em duas partes separadas por `@`, constante e FQDN. Em sentido contrário, TS 24.109 K.2 c)3): primeira parte `"3gpp-akma"`, exemplo `3gpp-akma@af1.operator.com`.

**O que o projeto faz.** Segue TS 33.535: `realm = "3GPP-bootstrapping-akma@<FQDN do AF>"`. O conteúdo — a constante e o FQDN — vem de TS 33.535 B.1.2.2; a forma `<constante>@<FQDN>` vem do Stage 3 (TS 24.109 §5.2.2.1 e K.2 b)2)), porque o Stage 2 diz o que o `realm` contém, não como as partes se separam. A primeira parte fica configurável.

**Por quê.** Três fontes (TS 33.535 B.1.2.2, TS 24.109 K.2 b)2) e a estrutura do GBA) contra uma (K.2 c)3)). A decisão anterior do projeto (registro de defeitos D-C, Adendo D §D.4.1) seguia c)3) com argumento que as duas leituras satisfazem (`docs/baseline/conferencia-normativa.md` §3.10).

**Impacto.** Diverge apenas de TS 24.109 K.2 c)3). Um UE que siga c)3) e verifique a primeira parte do `realm` não interopera — no MVP o cliente é do próprio projeto. O registro de defeitos (D-C) e o Adendo D precisam ser atualizados.

**Implementação (18/10/2026).** O AF desafia com `realm="3GPP-bootstrapping-akma@af1.akma.example"` (`realmPrefix` configurável em `af/config/afcfg.yaml`) e recusa credenciais com outro `realm`. O cliente só aceita desafio cuja primeira parte seja a esperada e cuja segunda parte seja o FQDN que ele usou (`TestChallengeChecks`). Verificado com o AAnF da VM do free5GC (`docs/af-side/teste-af.md`).

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

---

## DV-9 — NEF–AF externo sem TLS nem certificado de cliente

**Norma.** `[3GPP]` TS 33.501:

- §12.2: entre o NEF e um AF fora do domínio do operador, *"mutual authentication based on client and server certificates shall be performed between the NEF and AF using TLS"*;
- §12.3: *"TLS shall be used to provide integrity protection, replay protection and confidentiality protection for the interface between the NEF and the AF. The support of TLS is mandatory"*;
- §12.4: o NEF *"shall authorize the requests from AF using OAuth-based authorization mechanism"* (RFC 6749).

TS 33.535 §6.3 passo 2: o NEF autoriza o AF, inclusive pela verificação do AF_ID.

**O que o free5GC faz.** `[FREE5GC]` O NEF v1.2.5 serve as APIs do lado do AF (`3gpp-traffic-influence`, `3gpp-pfd-management`) no mesmo servidor SBI das outras NFs:

- em `http` ou `https` (`sbi.scheme`), só com certificado de servidor;
- cada grupo de APIs protegido por um *middleware* OAuth2 que confere o *token* emitido pelo NRF para o nome do serviço, quando o NRF tem `oauth: true`.

Não há verificação de certificado de cliente.

**O que o projeto faz.** `[MVP]` A API `3gpp-akma` do NEF segue o padrão do free5GC:

- `http`, como o resto do laboratório;
- *token* OAuth2 do NRF para `3gpp-akma`, conferido pelo mesmo *middleware*: atende a §12.4, com o NRF no papel de servidor de autorização;
- política local do NEF que associa cada AF_ID externo aos endereços de origem do AF, como a do AAnF no DV-1.

Não há TLS nem autenticação mútua por certificados (§12.2, §12.3).

**Por quê.** Decidido pelo responsável em 19/10/2026 (MVP-6, opção b), entre duas opções:

- (a) TLS com certificado de cliente, só para `3gpp-akma`;
- (b) esta.

O certificado de cliente exige mudar o servidor SBI comum do NEF do free5GC. O TLS das interfaces SBI e NEF–AF fica para um MVP próprio, junto com o DV-1 (plano, seção 5).

**Impacto.**

- **Transporte.** A K_AF trafega em claro entre o NEF e o AF externo. Um atacante na rede entre os dois a lê.
- **Autenticação do AF.** Ela se apoia no *token* e no endereço de origem, e não num certificado. O risco é o mesmo do DV-1, agora do lado do NEF: um AF que obtenha *token* para `3gpp-akma` e esteja num endereço autorizado pode apresentar o AF_ID de outro.
- **Origem do *token*.** `[MVP]` O AF externo obtém o *token* no NRF do operador, o que um AF fora da rede não faria. Fora do laboratório, quem emite o *token* é o CAPIF ou o servidor de autorização do NEF (§12.4, §12.5). Sem CAPIF, TS 29.522 §6 deixa o servidor de autorização para a configuração local. O CAPIF do MVP-9, aprovado em 29/10/2026 (plano, seção 4), fecha esta parte do desvio quando for implementado. Com o CAPIF, o método 3 também pede TLS entre o AF e o NEF (TS 33.122 §6.5.2.3 passo 5), que fica com este desvio até o MVP-11.

**Implementação (19/10/2026).** `free5gc/patches/nef/0001-akma-api.patch` e `free5gc/patches/config/0002-nef-akma.patch` (`docs/free5gc-side/nef.md` §10). O NEF serve `3gpp-akma` em `http` em 192.168.16.130:8000. O AF de fora da rede obtém o *token* no NRF pelo túnel SSH. Verificação: o NEF responde 401 sem *token* (E04, `docs/af-side/teste-nef.md`) e 403 a um AF_ID fora da política (E06).

**CAPIF (02/11/2026, MVP-9).** O AF externo também pode obter o *token* do CCF, que é o caminho da norma (TS 29.522 §7.2), e o NEF aceita os dois modos (`allowNrfTokens: true`). CAPIF-2e, entre o AF e o NEF, continua em `http` até o MVP-11, com a mesma política de endereço. O `check-authentication` também passa por essa política. Evidência: [CP04], [CP06], [CP16], [CP17] e [R2c] na execução `20261102-093804` (`docs/af-side/e2e.md` §5.9).

**MVP-11 (02/11/2026).** Decisão do responsável no ponto T14 do MVP-11, opção (b): o AF de fora usa só o modo CAPIF, com TLS mútuo em CAPIF-2e pelo método 3 (TS 33.122 §6.5.2.3 passo 5) e o certificado que o CCF emitiu na entrada. O NEF deixa de aceitar *token* do NRF na `3gpp-akma` (`allowNrfTokens: false`), e o AF de fora deixa de falar com o NRF, que é o servidor de autorização das NFs do operador (TS 33.501 §13.4.1.1.1). O desvio se encerra inteiro, sem atalho: o transporte (§12.2 e §12.3), a autenticação mútua (§12.5) e a origem do *token* (TS 29.522 §7.2). O modo sem CAPIF (TS 29.522 §6) fica no código do NEF, desligado; um servidor de autorização próprio para ele pode entrar no MVP-14.

**Especificação do MVP-11 (04/11/2026).** Aprovada com os pontos T19 e T34 de `docs/free5gc-side/tls-sbi.md`. O transporte de CAPIF-2e entra na fase (a), com o grupo do NEF do MVP-6 no modo CAPIF. A ligação do certificado ao *token* (o CN do certificado igual ao `client_id`) e `allowNrfTokens: false` entram na fase (c). A política de endereço do NEF sai do modo CAPIF: `sourceAddresses` fica opcional, só para o modo sem CAPIF. `[LACUNA]` TS 29.222 não tem atributo para a indicação do TLS mútuo que TS 33.122 §6.5.2.3 passo 5 atribui ao CCF (defeito D-AJ): a política do CCF é sempre o TLS mútuo.

---

## DV-10 — Perfis TLS do Ua\* sem `D0 02` e sem OCSP *stapling*

**Norma.** `[3GPP]` TS 33.210 V19.3.0 §6.2, o perfil TLS para o qual TS 33.222 §5.3.1 e §5.4.1 e TS 24.109 K.2 e K.3.1 remetem, por TS 33.310 Anexo E. Duas recomendações (*should*), não obrigações (*shall*):

- §6.2.3, PSK no TLS 1.2: *"the following cipher suites are mandatory to support and recommended to use: TLS_ECDHE_PSK_WITH_AES_128_GCM_SHA256"* e *"Support of the following cipher suite is recommended: TLS_ECDHE_PSK_WITH_AES_256_GCM_SHA384"* (RFC 8442, `D0 02`);
- §6.2.2, TLS 1.3: *"The TLS Certificate Status Request extension (i.e., "OCSP stapling") as defined in RFC 6066 and RFC 8446 should be supported"*; §6.2.3, TLS 1.2: *"The OCSP Status (a.k.a. certificate status request) extension, defined in RFC 6066 should be supported"*.

**O que o projeto faz.** `[MVP]`

- **`D0 02`.** O perfil `psk-tls12` só tem `TLS_ECDHE_PSK_WITH_AES_128_GCM_SHA256` (`D0 01`), a suíte obrigatória: o AF só a aceita, e o `ClientHello` do UE só oferece essa com PSK. O identificador Ua\* fica `01 00 01 D0 01`.
- **OCSP *stapling*.** O AF do K.2 (`digest-tls`) não grampeia resposta OCSP ao certificado. O cliente Go manda a extensão `status_request`, e o AF a ignora. Os perfis PSK não têm certificado, e para eles a extensão não se aplica.

**Por quê.** Decidido pelo responsável em 21/10/2026, ao fechar a fase (b) do MVP-8:

- o wolfSSL 5.9.4 não tem a suíte `D0 02`: o código só define `ECDHE-PSK-AES128-GCM-SHA256` entre as suítes `ECDHE_PSK` com AES-GCM (`wolfssl/internal.h`; lista de suítes da biblioteca compilada, conferida em 21/10/2026). O `crypto/tls` do Go não tem PSK. Fazer a suíte exigiria mexer no wolfSSL;
- a CA de teste do laboratório não tem respondedor OCSP. Sem ele, não há o que grampear. TS 33.222 §5.3.1.0 também deixa a revogação de certificados opcional: *"Support of certificate revocation and of the related fields in certificates is optional"*.

**Impacto.**

- **`D0 02`.** O PSK TLS 1.2 fica só com AES-128-GCM e SHA-256, que é o nível da suíte obrigatória. Um UE ou AF de outra implementação que oferecesse só `D0 02` não conectaria; como `D0 01` é obrigatória, toda implementação conforme a tem. Nenhum impacto no TLS 1.3, cujas suítes vêm do RFC 8446.
- **OCSP.** O UE não recebe, no *handshake*, prova de que o certificado do AF não foi revogado. Com a CA de teste, não há revogação nenhuma. Num ambiente com PKI de verdade, o `crypto/tls` grampeia a resposta que a aplicação lhe der (`tls.Certificate.OCSPStaple`); o AF precisaria buscá-la no respondedor da CA.

**Implementação (21/10/2026).** `af/wolftls/wolftls.go` (listas de suítes do servidor e do cliente); `af/akmatls/akmatls.go` (`ServerConfig`, sem `OCSPStaple`). Evidência: `TestTLS12Hello` (`ClientHello` com `D0 01` como única suíte PSK; `ServerHello` com `D0 01`); [R2t.psk12] e [T04] na execução `20261021-164401` (`docs/af-side/e2e.md` §5.4). Especificação: `docs/af-side/ua-star-norma.md` §3.8, itens 2 e 7.

**Como fechar.** `D0 02`: quando o wolfSSL a tiver, ou com outra pilha de PSK TLS; basta acrescentá-la às listas e o AF_ID `<FQDN>.010001D002` às políticas. OCSP: com uma CA que tenha respondedor, o AF busca a resposta e a põe em `OCSPStaple`.

---

## DV-11 — CoAP sobre DTLS sem `TLS_PSK_WITH_AES_128_CCM_8`

**Norma.** Dois textos em conflito:

- `[3GPP]` TS 33.535 C.2.1: *"The TLS profile specified in TS 33.210 clause 6.2 apply to DTLS 1.3"*. TS 33.210 V19.3.0 §6.2.3: *"Only cipher suites with AEAD (e.g. GCM) and PFS (i.e., ECDHE) shall be supported"*; com PSK, `TLS_ECDHE_PSK_WITH_AES_128_GCM_SHA256` é obrigatória. A NOTA 2 de §6.2.1 estende o perfil ao DTLS;
- RFC 7252 §9.1.3.1, sobre o CoAP no modo PSK: *"Implementations in these modes MUST support the mandatory-to-implement cipher suite TLS_PSK_WITH_AES_128_CCM_8"*. É uma suíte de DTLS 1.2 sem ECDHE, sem sigilo futuro.

`[DIVERGÊNCIA]` O Anexo C de TS 33.535 não diz o protocolo de aplicação. O projeto pôs CoAP dentro do DTLS (ponto U14 de `docs/af-side/ua-star-norma.md`), e com ele veio a regra da RFC 7252.

**O que o projeto faz.** Segue TS 33.210 (ponto U15): os perfis `psk-dtls13` e `psk-dtls12` só têm suítes com ECDHE e AEAD. No DTLS 1.3, `TLS_AES_128_GCM_SHA256` com `psk_dhe_ke`; no DTLS 1.2, `TLS_ECDHE_PSK_WITH_AES_128_GCM_SHA256`. `TLS_PSK_WITH_AES_128_CCM_8` não é oferecida nem aceita.

**Por quê.** Decidido pelo responsável em 21/10/2026, com os pontos U13 a U21. O Ua\* do AKMA é um ponto de referência 3GPP, e TS 33.535 C.2.1 aplica a ele o perfil de TS 33.210, que proíbe suítes sem ECDHE. A RFC 7252 pensa em dispositivos restritos, sem capacidade para ECDHE, o que não é o caso de um UE.

**Impacto.** Um cliente CoAP que só tenha `TLS_PSK_WITH_AES_128_CCM_8` não conecta ao AF. É o que TS 33.210 quer: sem sigilo futuro, quem obtivesse a K_AF leria as sessões passadas.

**Implementação (21/10/2026).** `af/wolftls/wolftls.go` (listas de suítes, as mesmas do TLS); `af/af/internal/profile/psk_wolfssl.go` (perfis DTLS); `af/coap` (CoAP). Verificação na execução do `e2e.sh` da fase (c) (`docs/af-side/e2e.md`).

## DV-12 — `trigger-auth` com POST, contra a precedência do YAML

**Norma.** Dois textos de TS 29.503 V19.7.0 em conflito (defeito D-B):

- `[3GPP]` tabela 6.2.3.13.4.1-1: a operação customizada `trigger-auth` em `/{ueId}/registrations/trigger-auth` tem o método mapeado POST. §5.3.2.16.2 passo 1: *"The NF Service Consumer sends a POST request (custom method: trigger-auth) to the UDM"*. A tabela 6.2.3.13.4.2.2-1 tem o título *"Data structures supported by the POST Request Body"*;
- `[3GPP]` `TS29503_Nudm_UECM.yaml` v1.4.2 (Anexo A): `get:` com `requestBody` obrigatório. O Anexo A.1 diz que o YAML *"takes precedence ... with respect to the encoding of information elements and methods"*.

`[DIVERGÊNCIA]` Seguir a tabela contraria a regra de precedência do Anexo A.1. Seguir o YAML exige um GET com corpo obrigatório, que contraria a semântica do GET (RFC 9110 §9.3.1) e o modelo das operações customizadas de TS 29.501 §4.6.1.2, sempre com POST.

**O que o projeto faz.** POST, como a tabela e o procedimento. O UDM responde `405` a um GET em `trigger-auth`, com `Allow: POST`. O cliente do *patch* `openapi/0003` (`AuthTrigger`) e o AAnF usam POST.

**Por quê.** Decidido pelo responsável no defeito D-B (22/10/2026) e confirmado no ponto K2 da especificação do MVP-12, aprovada em 28/10/2026 (`docs/free5gc-side/renovacao-kakma.md`). Três textos de TS 29.503 e o modelo de TS 29.501 apontam POST; só o YAML aponta GET, e um GET com corpo não é tratado de forma previsível por servidores e clientes HTTP.

**Impacto.** Um consumidor que siga o YAML ao pé da letra manda GET e recebe 405. Se a 3GPP corrigir o YAML para POST, o desvio se encerra sem mudança no código.

**Implementação (28/10/2026).** `free5gc/patches/udm/0002-auth-trigger.patch` (rota POST, GET com 405); `free5gc/patches/openapi/0003-home-network-triggered-authentication.patch` (cliente); `free5gc/aanf/internal/sbi/consumer/udm_service.go`. Verificação: teste de unidade `TestAuthTriggerRequests` do UDM e o passo R6r do `e2e.sh` (`docs/af-side/e2e.md`).

---

## DV-13 — *Claims* do *token* do CCF: `iss` com o identificador do *invoker* e `exp` como instante

**Norma.** Três textos em conflito (defeitos D-AA e D-AB):

- `[3GPP]` TS 29.222 V19.8.0, tabela 8.5.4.2.8-1 (`AccessTokenClaims`), e `TS29222_CAPIF_Security_API.yaml`: `iss`, `scope` e `exp` obrigatórios.
  - `iss`: *"This attribute shall contain the API invoker Identifier as defined in Annex C.2.2 of 3GPP TS 33.122 [16] and in clause 4.1.1 of IETF RFC 7519 [24]."*
  - `exp`, do tipo `DurationSec`: *"This attribute shall contain the number of seconds after which the access_token is considered to be expired as defined in clause 4.1.4 of IETF RFC 7519 [24]."*
- `[3GPP]` TS 33.122 V19.4.0, C.2.2: o *token* leva `exp` (*"The expiration time of the access token"*), `client_id` (*"The identifier of the API Invoker making the API request"*) e `scope`, todos *"REQUIRED"*; não há `iss`. C.6: *"If the current time is beyond the time specified by the "exp" parameter, the associated token shall no longer be considered valid"*.
- RFC 7519: §4.1.1, *"The "iss" (issuer) claim identifies the principal that issued the JWT. The processing of this claim is generally application specific."*; §4.1.4, `exp` *"identifies the expiration time on or after which the JWT MUST NOT be accepted for processing"*, com valor *NumericDate*.

`[DIVERGÊNCIA]` Os dois pontos estão iguais em TS 29.222 V20.1.0. TS 33.122 não tem versão mais nova.

**O que o projeto faz.** O *token* do CCF leva:

- `iss` e `client_id`, os dois com o identificador do *invoker*;
- `exp` como instante (*NumericDate*), um inteiro que cabe no tipo `DurationSec` do YAML. A duração vai em `expires_in` da resposta (tabela 8.5.4.2.7-1).

O NEF confere `exp` como instante, com tolerância de até 30 s (TS 33.122 C.2.2).

**Por quê.** Decidido pelo responsável em 29/10/2026, com o ponto C9 do MVP-9 (plano, seção 4). Com as duas *claims*, TS 29.222 e TS 33.122 ficam atendidas, e a RFC 7519 deixa o tratamento de `iss` à aplicação. Um `exp` relativo só teria sentido com o instante da emissão (`iat`), que nenhuma das duas normas pede, e contrariaria a RFC 7519, que TS 29.222 cita, e TS 33.122 C.6.

**Impacto.** Um AEF de outra implementação que espere `iss` com o identificador do CCF recusa o *token*. Um CCF de outra implementação que mande `exp` relativo tem os *tokens* recusados pelo NEF, como vencidos. O AF não é afetado: o *token* é opaco para ele, que usa `expires_in` (TS 33.122 C.5).

**Implementação (02/11/2026).** O CCF emite o *token* com `iss` e `client_id` iguais ao identificador do *invoker* e `exp` como instante (`free5gc/ccf/internal/context/context.go`, `IssueToken`). O NEF confere os dois e `exp`, com 30 s de tolerância (`free5gc/patches/nef/0002`). Evidência: [CP05] e [CP06] na execução `20261102-093804` (`docs/af-side/e2e.md` §5.9); `TestCapifFlow` no CCF e `TestCapifAccessToken` no NEF.

---

## DV-14 — CAPIF parcial: só o caminho do AKMA

**Norma.**

- `[3GPP]` TS 29.522 V19.7.0 §7.1: *"When CAPIF is used with an NEF that is used for external exposure, the NEF shall support the following as defined in 3GPP TS 29.222 [12]:"* o AEF e as APIs de CAPIF-2/2e e CAPIF-3/3e; o APF e as de CAPIF-4/4e; a função de gestão de API (*API management function*) e as de CAPIF-5/5e; pelo menos um dos métodos de segurança.
- `[3GPP]` TS 23.222 V19.10.0 §6.3.3 e §6.3.4: as capacidades do CCF (entre outras, o *logging* e a auditoria das chamadas, o monitoramento, a tarifação, as políticas e a interconexão) e as do AEF (entre outras, a revogação, o *logging* no CCF e a ocultação da topologia).
- `[3GPP]` TS 33.122 V19.4.0 §6.8, passos 7 a 10: na saída do *invoker*, o CCF manda ao AEF uma notificação de evento, e o AEF apaga os dados de segurança do *invoker*.

**O que o projeto faz.** O CCF e o NEF implementam do CAPIF só o caminho do AKMA (plano, MVP-9, ponto C2):

- CCF: `Onboard_API_Invoker`, só com a entrada automática (201), e `Offboard_API_Invoker`; `Discover_Service_API`; `Publish_Service_API`, `Unpublish_Service_API` e `Get_Service_API`; `Obtain_Security_Method`, `Obtain_Authorization` e `Obtain_API_Invoker_Info`;
- NEF: o AEF (`check-authentication` e a verificação do *token*) e o APF (a publicação de `3gpp-akma`). O registro do provedor fica na configuração (ponto C6).

Ficam de fora: a atualização da entrada e a entrada com aprovação; os eventos, inclusive o aviso ao AEF na saída do *invoker*; a revogação; o *logging* das chamadas; a auditoria; a política de acesso; o roteamento; o registro do provedor pela API; a descoberta aberta; a interconexão de CCF; o RNAA; a ocultação de topologia; a tarifação. Os métodos 1 e 2 de CAPIF-2e pedem TLS no NEF e ficam para o MVP-11 (DV-9).

**Por quê.** Decidido pelo responsável em 29/10/2026, com o ponto C16 do MVP-9. O MVP-9 leva o *token* do AF externo ao caminho da norma (TS 29.522 §7.2); as outras funções não entram no caminho do AKMA.

**Impacto.**

- Sem a revogação e sem o aviso na saída, um *token* já emitido vale até `exp`, mesmo depois da saída do *invoker*. TS 33.122 C.6 usa `exp` como forma de revogação; a vida curta do *token* limita a janela, e a especificação do MVP-9 fixa por quanto tempo o NEF guarda os dados do *invoker*.
- Sem *logging*, auditoria e tarifação no CCF, as chamadas ficam registradas só no log do NEF (categoria `AKMA`).
- Um *invoker* de outra implementação que dependa de uma API de fora (os eventos, por exemplo) não funciona com este CCF.

**Implementação (02/11/2026).** O CCF serve só as operações listadas acima (`free5gc/ccf/`). O NEF faz só a publicação, o `check-authentication` e a conferência do *token* (`free5gc/patches/nef/0002`). Evidência: o grupo CP e o R2c da execução `20261102-093804` (`docs/af-side/e2e.md` §5.9).

---

## DV-15 — HTTP/2 por ALPN no CCF, contra o *upgrade* do HTTP de TS 29.222 §7.3

**Norma.** Dois textos em conflito (defeito D-AE):

- `[3GPP]` TS 29.222 V19.8.0 §7.3: *"For CAPIF APIs, the support of HTTP/1.1 (IETF RFC 9112 [4], IETF RFC 9110 [5], and IETF RFC 9111 [8]) over TLS is mandatory and the support of HTTP/2 (IETF RFC 9113 [10]) over TLS is recommended. TLS shall be used as specified in 3GPP TS 33.122 [16]."* e *"A functional entity desiring to use HTTP/2 shall use the HTTP upgrade mechanism to negotiate applicable HTTP version as described in IETF RFC 9113 [10]."* Igual na V20.1.0.
- RFC 9113 §3.1, sobre o `h2c` no *Upgrade* do HTTP: *"This usage was never widely deployed and is deprecated by this document."* §3.2: *"A client that makes a request to an "https" URI uses TLS [TLS13] with the ALPN extension [TLS-ALPN]."* e *"The "h2c" protocol identifier MUST NOT be sent by a client or selected by a server; the "h2c" protocol identifier describes a protocol that does not use TLS."*

`[DIVERGÊNCIA]` A RFC que §7.3 cita para o mecanismo o deprecou e, sobre TLS, que o próprio §7.3 exige, não o admite.

**O que o projeto faz.** O CCF serve CAPIF-1e com HTTP/1.1 e HTTP/2, negociados por ALPN (`h2` e `http/1.1`), e não faz o *upgrade* do HTTP. O servidor interno (CAPIF-3/4, sem TLS até o MVP-11, DV-1) serve só HTTP/1.1.

**Por quê.** Decidido pelo responsável em 29/10/2026, com o ponto C29 da especificação do MVP-9 (`docs/free5gc-side/capif.md`). Sobre TLS, o HTTP/2 só se negocia por ALPN, como fazem o `crypto/tls` e o `net/http` do Go.

**Impacto.** Um cliente que siga §7.3 ao pé da letra e mande `Upgrade: h2c` numa conexão TLS continua em HTTP/1.1, que §7.3 torna obrigatório: não fica sem serviço.

**Implementação (02/11/2026).** `free5gc/ccf/internal/sbi/server.go` (`TLSConfig`, com `h2` e `http/1.1` no ALPN). Evidência: `TestCapifTLS` no CCF, com HTTP/2 por ALPN e HTTP/1.1 no TLS 1.2. Na execução `20261102-093804` (`docs/af-side/e2e.md` §5.9), usaram o CCF o AF, cujo cliente Go oferece `h2` no ALPN, e o `capif_api.py`, em HTTP/1.1.

---

## DV-16 — CAPIF-1e sem revogação de certificados, e *token* do CCF sem `x5t#S256`

**Norma.**

- `[3GPP]` TS 33.122 V19.4.0 §6.3.1.1, em CAPIF-1e: *"Certificate based authentication shall follow the profiles given in 3GPP TS 33.310 [2], subclauses 6.1.3a and 6.1.4a. The structure of the PKI used for the certificate is out of scope of the present document."*
- `[3GPP]` TS 33.310 V19.5.0:
  - §6.1.3a, no certificado de entidade TLS: *"Mandatory non-critical Distribution points: CRL distribution point."*;
  - §6.1.1, em todo certificado: *"CRLs as specified in subclause 6.1a shall be supported for certificate revocation verification."* e *"OCSP as specified in subclause 6.1b should be supported for certificate revocation verification."*
- `[3GPP]` TS 33.210 V19.3.0, o perfil TLS de CAPIF-1e (TS 33.122 §6.1 e §6.3.1.1, por TS 33.310 Anexo E): o OCSP *stapling* *"should be supported"*, no TLS 1.3 (§6.2.2) e no TLS 1.2 (§6.2.3). São as frases do DV-10.
- `[3GPP]` TS 33.210 V19.3.0 §6.3.3.4, no JWS do *token* de acesso: *"The "x5t#S256" header parameter should be supported."*

**O que o projeto faz.** `[MVP]`

- Os certificados do CCF e os dos *invokers* não têm o ponto de distribuição de CRL. O CCF não publica CRL, não responde OCSP e não confere revogação.
- O CCF não grampeia resposta OCSP ao seu certificado em CAPIF-1e.
- O *token* do CCF não leva `x5t#S256`. A chave que o assina não tem certificado: o NEF a recebe por configuração (`ccf-token.pub.pem`), e o *token* leva `kid`, a impressão digital da chave (TS 33.210 §6.3.3.1; RFC 7638).

O resto dos perfis é seguido: TS 33.310 §6.1.1, §6.1.3a e §6.1.4a; TS 33.210 §6.2 e §6.3.3.1. A conferência está em `docs/free5gc-side/capif.md` §4.2, §5.8 e §15.1.

**Por quê.** Decidido pelo responsável em 30/10/2026, na implementação do CCF:

- TS 33.122 §6.3.1.1 deixa a estrutura da PKI fora da norma, e o laboratório não tem serviço de revogação. Uma CRL pediria mais uma porta HTTP no CCF, o endereço dela nos certificados e a consulta em cada ponta. O OCSP fica de fora pelo mesmo motivo do DV-10.
- O efeito da revogação já existe em CAPIF-1e. O CCF confere o identificador do *invoker*, que é o CN do certificado, em toda operação. Depois da saída (TS 33.122 §6.8 passo 3), o CCF não conhece mais esse identificador e recusa o pedido: `403` ou `404`, ou `invalid_client` no *token*. O NEF pede os dados do *invoker* ao CCF e, depois da saída, recebe `404`.
- O `x5t#S256` é uma recomendação (*should*). Com uma única chave de *token*, configurada no NEF, o `kid` basta para identificá-la.

**Impacto.**

- O certificado de um *invoker* que saiu continua válido até a data de fim, um ano por padrão, mas o CCF não o aceita mais em nenhuma operação.
- No MVP-11, quando o NEF conferir o certificado do *invoker* em CAPIF-2e (TS 33.122 §6.5.2.2), a revogação passa a importar também no NEF. O ponto volta lá, com os certificados da SBI: TS 33.310 §6.1.3c também exige o ponto de distribuição de CRL.
- Um cliente que exija o ponto de distribuição de CRL ou o OCSP *stapling* recusa o certificado do CCF. O cliente do Go e o módulo `ssl` do Python não exigem nenhum dos dois.
- Um AEF de outra implementação que só ache a chave do *token* por `x5t#S256` não consegue conferir o *token*. Ele pode usar o `kid`, que TS 33.210 §6.3.3.1 manda suportar.

**Implementação (30/10/2026).**

- `free5gc/scripts/make-capif-certs.sh`: certificados sem `crlDistributionPoints`;
- `free5gc/ccf/internal/context/context.go`: `Onboard`, o certificado do *invoker*;
- `free5gc/ccf/internal/sbi/server.go`: `TLSConfig`, sem `OCSPStaple`;
- `free5gc/ccf/internal/jws/jws.go`: `kid`, sem `x5t#S256`.

Testes, em `free5gc/ccf/internal/sbi/capif_test.go`: `checkInvokerCertificate`, para certificados sem ponto de distribuição de CRL e sem OCSP; `TestCapifFlow`, para o cabeçalho do *token* só com `alg`, `kid` e `typ`; `TestCapifTLS`, para a sessão sem resposta OCSP.

No laboratório, na execução `20261102-093804` (`docs/af-side/e2e.md` §5.9), o AF e o `capif_api.py` usaram o CCF sem CRL e sem OCSP.

**Como fechar.**

- Uma CRL assinada pela CA do CCF, servida por HTTP num endereço posto nos certificados (`crlDistributionPoints`), com o certificado do *invoker* revogado na saída.
- O OCSP *stapling*, com um respondedor.
- O `x5t#S256`, com um certificado para a chave do *token*.

Rever no MVP-11.

**MVP-11 (02/11/2026).** O MVP-11, aprovado nesta data, traz dois pontos que mudam este desvio: o T7, que manda conferir a CRL em toda ponta TLS (TS 33.310 §6.3.1), e o T14, opção (b), em que o NEF passa a conferir o certificado do *invoker* em CAPIF-2e. A especificação do MVP-11 revê o desvio com os dois.

**Revisão (04/11/2026, MVP-11).** Decisão do responsável no ponto T36 de `docs/free5gc-side/tls-sbi.md` (§8.4). Na fase (c):

- os certificados que o CCF emite e o `ccf.pem` ganham o ponto de distribuição `http://192.168.16.130:8030/ccf-ca.crl`;
- o CCF assina a CRL da sua CA e revoga o certificado do *invoker* na saída (TS 33.122 §6.8);
- o CCF, o NEF e o AF conferem essa CRL nas suas pontas.

A parte da revogação sai deste desvio com a implementação. Ficam o OCSP *stapling*, como no DV-10, e o `x5t#S256` do *token* do CCF, cuja chave continua sem certificado. A CRL vai só por HTTP (DV-18).

---

## DV-17 — SBI por IP, como modo opcional

**Norma.**

- `[3GPP]` TS 29.510 V19.7.0 §6.1.6.2.2, NOTA 1: *"If the NF supports the NF services with "https" URI scheme (i.e use of TLS is mandatory), then the FQDN shall be provided in the NF Profile or the NF Service profile (see clause 6.1.6.2.3) and it shall be used to construct the target URI"*. §6.1.6.2.3, NOTA 8, diz o mesmo.
- `[3GPP]` TS 33.501 V19.7.0 §13.1.0: *"TLS clients shall include the SNI extension as specified in RFC 9113"*. A RFC 6066 §3 não aceita endereço IP no SNI.
- `[3GPP]` TS 33.310 V19.5.0 §6.1.3c.3: *"subjectAltName shall not contain only IP address in TLS server certificates"*; o certificado do servidor tem o DNS-ID com o FQDN da NF.

**O que o projeto faz.** `[MVP]` Dois modos, escolhidos na configuração de cada NF, do AAnF, do CCF e do AF (ponto T4 do MVP-11):

- por nome, o padrão e o da norma: o FQDN no perfil da NF, nos URIs e no SNI;
- por IP, opcional: os URIs com o endereço, como antes do MVP-11; o cliente confere o IP do certificado e não manda o SNI.

Os certificados são os mesmos nos dois modos, com o FQDN e o IP no `subjectAltName`; nunca só o IP.

**Por quê.** Decidido pelo responsável em 02/11/2026, no ponto T4 do MVP-11: o modo por IP serve a ambientes em que os nomes não resolvem. O modo padrão segue a norma; o registro existe porque o modo por IP, quando escolhido, não segue os três textos acima.

**Impacto.** No modo por IP, o perfil da NF não tem o FQDN que TS 29.510 pede com `https`, o cliente não manda o SNI, e a identidade do servidor é conferida pelo IP do certificado, não pelo nome. O TLS mútuo, a CRL e a ligação do *token* ao certificado (T7, T9 e T10) não mudam.

**Implementação.** A fazer no MVP-11, pela especificação `docs/free5gc-side/tls-sbi.md`, aprovada em 04/11/2026 (T21; §3.2). A escolha é a chave `sbi.fqdn` de cada NF e do AAnF (`internal.fqdn` no CCF): com ela, por nome; sem ela, por IP, com um aviso no log quando há `https`. Quem só é cliente, como o AF, segue a forma do URI configurado. Nas configurações do projeto, o modo é por nome.

---

## DV-18 — CRL só por HTTP, sem LDAPv3

**Norma.** `[3GPP]` TS 33.310 V19.5.0:

- §6.1a: *"CRL retrieval with LDAPv3 (RFC 4510 [67], RFC 4517 [68], RFC 4523 [69], RFC 4512 [70]) shall be supported as the primary method. HTTP may be used for checking the revocation status of TLS and NE certificates."*
- §6.1.1: *"CRLs as specified in subclause 6.1a shall be supported for certificate revocation verification."*
- tab. 6.1.3c.3-1, `cRLDistributionPoint`: *"this indicates if the CRL is available for retrieval using access protocol and location with LDAP or HTTP URI."*

**O que o projeto faz.** `[MVP]` As CRL da CA do operador e da CA do CCF ficam num serviço HTTP, em `http://192.168.16.130:8030/<ca>.crl`. Cada ponta TLS as baixa pelo ponto de distribuição do certificado do outro lado. O LDAPv3 não é suportado: nem o serviço nem as pontas falam LDAP, e um certificado cujo ponto de distribuição não tenha URI `http` é recusado.

**Por quê.** Decidido pelo responsável em 04/11/2026, no ponto T27 do MVP-11, opção (a), entre três:

- (a) esta;
- (b) ler o *"HTTP may be used"* como suficiente para entidades TLS, sem registrar desvio;
- (c) o LDAPv3 também, com um servidor de diretório na máquina do free5GC e um cliente LDAP em Go.

O próprio §6.1a permite o HTTP para os certificados TLS, que são os deste MVP. O laboratório não tem diretório. O registro existe porque o *"shall be supported"* do LDAPv3 fica sem atender.

**Impacto.** Um operador com a CRL só num diretório LDAP não serve a este código: o certificado com o ponto de distribuição `ldap://` é recusado. A conferência da CRL em cada ponta, a falha fechada e o resto do T7 não mudam.

**Implementação.** A fazer no MVP-11, pela especificação `docs/free5gc-side/tls-sbi.md` (§5).

**Como fechar.** Um leitor de CRL por LDAPv3 no pacote `sbitls`, com um diretório no ponto de distribuição.
