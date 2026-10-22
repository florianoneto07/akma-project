# Registro de desvios e escolhas entre textos normativos

**Criado em:** 16 de outubro de 2026
**Regra:** todo ponto em que o projeto não segue literalmente a norma — por limitação do free5GC, simplificação de MVP, ambiguidade ou contradição do texto — é apresentado com a cláusula, o motivo e o impacto, decidido pelo responsável pelo projeto e registrado aqui antes de ser implementado.

Estados: **decidido** (escolha explícita do responsável), **aprovado com o plano** (recomendação aprovada junto com o plano de desenvolvimento em 16/10/2026), **implementado** (com referência ao código e ao teste).

| Nº | Assunto | Estado | MVP |
|---|---|---|---|
| DV-1 | SBI sem TLS (h2c) | **implementado** (MVP-1, 17/10/2026) | MVP-1 |
| DV-2 | Subscrição de `ROAMING_STATUS` desabilitada | **implementado** (MVP-1, 17/10/2026) — parcial: a etapa não existe; o campo só aceita `false` | MVP-1 |
| DV-3 | Só SUPI do tipo IMSI | **implementado** (MVP-0, 16/10/2026) | MVP-0 |
| DV-4 | Sem NFKC no lado UE | **implementado** (MVP-4, 18/10/2026) | MVP-4 |
| DV-5 | Digest do Ua\* conforme RFC 7616 | **implementado** (MVP-3, 18/10/2026) | MVP-3 |
| DV-6 | `realm` do Digest conforme TS 33.535 | **implementado** (MVP-3, 18/10/2026) | MVP-3 |
| DV-7 | Nunca emitir 204 em `retrieve-applicationkey` | **implementado** (MVP-1, 17/10/2026) | MVP-1 |
| DV-8 | `aud` vazio tolerado no token OAuth2 | **encerrado** em 17/10/2026 (MVP-1b: NRF corrigido, AAnF estrito) | MVP-1 |
| DV-9 | NEF–AF externo sem TLS nem certificado de cliente | **implementado** (MVP-6, 19/10/2026) | MVP-6 |
| DV-10 | Perfis TLS do Ua\* sem `TLS_ECDHE_PSK_WITH_AES_256_GCM_SHA384` e sem OCSP *stapling*, recomendados por TS 33.210 | **implementado** (MVP-8, 21/10/2026) | MVP-8 |
| DV-11 | CoAP sobre DTLS sem a suíte obrigatória de RFC 7252, `TLS_PSK_WITH_AES_128_CCM_8`; vale TS 33.210 | **implementado** (MVP-8, 21/10/2026) | MVP-8 |

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
- **Origem do *token*.** `[MVP]` O AF externo obtém o *token* no NRF do operador, o que um AF fora da rede não faria. Fora do laboratório, quem emite o *token* é o CAPIF ou o servidor de autorização do NEF (§12.4, §12.5). Sem CAPIF, TS 29.522 §6 deixa o servidor de autorização para a configuração local. O CAPIF está planejado para o MVP-9 (plano, seção 5) e fecha esta parte do desvio.

**Implementação (19/10/2026).** `free5gc/patches/nef/0001-akma-api.patch` e `free5gc/patches/config/0002-nef-akma.patch` (`docs/free5gc-side/nef.md` §10). O NEF serve `3gpp-akma` em `http` em 192.168.16.130:8000. O AF de fora da rede obtém o *token* no NRF pelo túnel SSH. Verificação: o NEF responde 401 sem *token* (E04, `docs/af-side/teste-nef.md`) e 403 a um AF_ID fora da política (E06).

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

