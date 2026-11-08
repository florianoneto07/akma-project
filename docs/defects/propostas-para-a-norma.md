# Propostas para a norma: opções de texto para os defeitos

**Criado em:** 8 de novembro de 2026 (MVP-14, ponto Q4) · **Base:** `docs/defects/registro-de-defeitos-3gpp.md`, que continua sendo a lista dos problemas da norma que o responsável vai enviar à 3GPP.

Pedido do responsável no Q4: para cada desvio que vem de um defeito da norma, as opções de texto que resolveriam o defeito, com a ideal indicada. Cada item traz o texto atual, o problema, as opções, a recomendação e o que o projeto faz hoje. Os textos propostos estão em inglês, como entrariam num *Change Request* (CR); a categoria é F (correção) em todos. A escolha entre as opções e o envio ficam com o responsável.

Os itens D-AL a D-AO entraram no registro nesta data (Q11) e também têm opções aqui.

| Item | Desvio | Norma | Grupo | Opção recomendada |
|---|---|---|---|---|
| D-F | DV-5 | TS 24.109 K.2 c) | CT1 | B: RFC 7616 e SHA-256 |
| D-C | DV-6 | TS 24.109 K.2 c) 3) | CT1 | A: `3GPP-bootstrapping-akma@<FQDN>` |
| D-G | DV-7 | TS 29.535 §4.2.2.3.2; TS 29.522 tab. 5.14.3.2.2-2 | CT3 | A: sem o 204 |
| D-U | DV-11 | TS 33.535 Anexo C | SA3 | B: o protocolo é da aplicação, com a regra das suítes |
| D-W | DV-11 | TS 33.535 C.2.1 | SA3 | A: as suítes de TS 33.210 prevalecem sobre as da RFC 7252 |
| D-B | DV-12 | `TS29503_Nudm_UECM.yaml` | CT4 | A: `post:` |
| D-AA | DV-13 | TS 29.222 tab. 8.5.4.2.8-1 | CT3 | A: `iss` com o CCF e `client_id` com o *invoker* |
| D-AB | DV-13 | TS 29.222 tab. 8.5.4.2.8-1 | CT3 | A: `exp` como instante |
| D-AE | DV-15 | TS 29.222 §7.3; TS 29.122 §5.2.2 | CT3 | A: ALPN |
| D-AL | DV-2 | TS 33.535 §6.2.1 passo 5, §6.3 passo 4 | SA3 | A: a subscrição só com a política de *roaming* |
| D-AM | — | TS 33.535 §6.8.2, §7.1.6; TS 29.522 tab. 5.14.5.3.4-1 | SA3, CT3 | B: `afId` opcional na notificação ao NEF |
| D-AN | — | TS 29.535 §4.2.2.2.2, tab. 5.1.4.3.2-2 | CT3 | A: GPSI só pela *feature*, e o registro só com SUPI |
| D-AO | — | TS 33.535, TS 29.535, TS 29.522 | SA3, CT3 | correções editoriais |

---

## D-F — A referência ao Digest no K.2 (DV-5)

- **Texto atual.** TS 24.109 V19.0.0, K.2 c): *"the UE shall generate the HTTP request and the AF shall authenticate the HTTP request using HTTP Digest. HTTP Digest authentication (see RFC 3310 [6]) shall be used with previously bootstrapped security association as follows:"*
- **Problema.** A RFC 3310 é o Digest AKA, em que a senha vem do RES da autenticação AKA. No AKMA a senha é a K_AF em base64 (K.2 c) 2)), e só a RFC 7616 funciona com ela. O K.2 b) remete à §5.2, cuja §5.2.2.1 já usa a RFC 7616 [36]; TS 33.222 §5.3.0, que TS 33.535 B.1.2.2 adapta, também, e desaconselha o MD5 (NOTA 1b).
- **Opções.**
  - **A.** Trocar só a referência: *"HTTP Digest authentication (see RFC 7616 [36]) shall be used with the AKMA Application Key as follows:"* Menor mudança possível; o algoritmo continua sem regra.
  - **B (ideal).** A troca de A e uma regra para o algoritmo: *"HTTP Digest authentication as specified in RFC 7616 [36] shall be used with the AKMA Application Key as follows:"* e um item novo: *"4) the "algorithm" parameter shall be "SHA-256"; the "MD5" algorithm shall not be used."* Dá interoperabilidade: os dois lados sabem o algoritmo.
  - **C.** Remeter ao modelo do GBA: *"HTTP Digest authentication shall be used as specified in clause 5.2.2.1, with the following changes:"*, e os itens 1) a 3) como estão. Mantém um só texto para os dois, mas herda do GBA o que o AKMA não usa (B-TID, GBA_U).
- **Recomendação.** B: corrige a referência e fecha a escolha do algoritmo, que hoje fica a cargo de cada implementação.
- **O projeto hoje.** Segue B: RFC 7616 com `algorithm=SHA-256`; o AF recusa `MD5` com 400.

## D-C — Os dois valores do `realm` (DV-6)

- **Texto atual.** TS 24.109 V19.0.0, K.2 b) 2): *"The realm attribute shall contain the constant string "3GPP-bootstrapping-akma"."* K.2 c) 3): *"the "realm" parameter shall contain two parts delimited by "@" sign where the first part is the constant string "3gpp-akma" and the latter part shall be the FQDN of the AF (e.g. "3gpp-akma@af1.operator.com")"*.
- **Problema.** No Digest, o cliente devolve o `realm` do desafio, e o valor entra no HA1: os dois itens têm de dar o mesmo `realm`. TS 33.535 B.1.2.2 passo 3 diz *"the AF shall include the "3GPP-bootstrapping-akma" within the WWW-Authenticate header field"*, e o modelo do GBA (§5.2.2.1 item 3) usa `<constante>@<FQDN>` com a constante do modo.
- **Opções.**
  - **A (ideal).** Em c) 3), *"the first part is the constant string "3GPP-bootstrapping-akma""*, com o exemplo *"3GPP-bootstrapping-akma@af1.operator.com"*. Uma palavra, numa norma, e as três fontes passam a dizer o mesmo.
  - **B.** O contrário: `3gpp-akma` em K.2 b) 2) e em TS 33.535 B.1.2.2. Muda duas normas, de dois grupos, e a constante do `realm` passa a ser igual ao *product token* do `User-Agent`, que tem outro papel.
  - **C.** Definir o `realm` inteiro uma vez, em b) 2) (*"The realm attribute shall contain two parts delimited by "@" sign: the constant string "3GPP-bootstrapping-akma" and the FQDN of the AF"*), e em c) 3) só remeter a ele. Resolve como A e tira a repetição.
- **Recomendação.** A, ou C se o grupo preferir tirar a repetição.
- **O projeto hoje.** Segue A: `realm="3GPP-bootstrapping-akma@<FQDN>"`, com a primeira parte configurável.

## D-G — O 204 sem condição em `retrieve-applicationkey` (DV-7)

- **Texto atual.** TS 29.535 V20.1.0, §4.2.2.3.2: *"If KAKMA is not present in the AAnF, the AAnF shall reply with an HTTP "403 Forbidden" status code … "K_AKMA_NOT_PRESENT""* e, na mesma cláusula, *"If the requested AKMA Application Key information for the UE does not exist, the AAnF shall respond with "204 No Content"."* O mesmo na tab. 5.1.4.3.2-2 e no YAML. TS 29.522 tab. 5.14.3.2.2-2 repete o 204 na API do NEF.
- **Problema.** Sem K_AKMA, a resposta é 403; com K_AKMA, o AAnF deriva a K_AF (TS 33.535 §6.2.1 passo 7). Não sobra caso para o 204.
- **Opções.**
  - **A (ideal).** Tirar o 204: a frase de §4.2.2.3.2, a linha da tab. 5.1.4.3.2-2 e a resposta `'204'` do `TS29535_Naanf_AKMA.yaml`; em TS 29.522, a linha da tab. 5.14.3.2.2-2 e a resposta do `TS29522_AKMA.yaml`.
  - **B.** Manter o 204 e dizer o caso, por exemplo: *"If the AAnF is not able to provide the AKMA Application Key information for the requested AF_ID for a reason other than the absence of KAKMA, …"*. Só serve se o grupo tiver um caso em mente; a recusa por política já é 403.
  - **C.** Trocar o 403 `K_AKMA_NOT_PRESENT` pelo 204. Tira a causa que o NEF tem de repassar ao AF (TS 29.522 §4.4.23.2) e quebra quem já segue a norma.
- **Recomendação.** A.
- **O projeto hoje.** O AAnF nunca responde 204; o AF e o NEF tratam o 204, por segurança.

## D-U — O protocolo de aplicação dentro do DTLS (DV-11)

- **Texto atual.** TS 33.535 V19.0.0, Anexo C: só o DTLS, sem o protocolo que vai dentro. TS 33.220 Anexo O, que ele adapta: *"It does not have any dependency on the application layer protocol."*
- **Problema.** O UE e o AF precisam do mesmo protocolo para trocar dados depois do *handshake*. Com o CoAP, o mais comum sobre DTLS, vem a suíte obrigatória da RFC 7252, que TS 33.210 proíbe (D-W).
- **Opções.**
  - **A.** Fixar o CoAP (RFC 7252) sobre DTLS, com o recurso e o conteúdo, como o Anexo D faz para o OSCORE. Interoperável, mas tira da aplicação a escolha que §4.4.1 NOTA 1 lhe dá.
  - **B (ideal).** Dizer que o protocolo é da aplicação e acrescentar a regra das suítes de D-W: *"The application layer protocol carried over DTLS is application specific. When CoAP (IETF RFC 7252) is used, the cipher suites of clause C.2.1 apply."* Respeita §4.4.1 e fecha o conflito.
  - **C.** Deixar como está.
- **Recomendação.** B.
- **O projeto hoje.** CoAP, esquema `coaps`, por escolha da aplicação (U14).

## D-W — A suíte obrigatória do CoAP contra TS 33.210 (DV-11)

- **Texto atual.** TS 33.535 V19.0.0, C.2.1: *"The TLS profile specified in TS 33.210 [13] clause 6.2 apply to DTLS 1.3"*. TS 33.210 V20.0.0 §6.2.3: *"Only cipher suites with AEAD (e.g. GCM) and PFS (i.e., ECDHE) shall be supported."* RFC 7252 §9.1.3.1: *"Implementations in these modes MUST support the mandatory-to-implement cipher suite TLS_PSK_WITH_AES_128_CCM_8"*, sem ECDHE.
- **Problema.** Um AF com CoAP sobre DTLS não consegue cumprir as duas regras.
- **Opções.**
  - **A (ideal).** Em C.2.1: *"When CoAP (IETF RFC 7252) is used over DTLS, the cipher suites of TS 33.210 clause 6.2 apply instead of the mandatory-to-implement cipher suites of clause 9.1.3 of IETF RFC 7252."* Diz qual regra vale.
  - **B.** Admitir `TLS_PSK_WITH_AES_128_CCM_8` para UEs com pouca capacidade. Volta a uma suíte sem sigilo futuro, contra §6.2.3.
  - **C.** A de A e, em TS 33.210 §6.2.3, recomendar `TLS_ECDHE_PSK_WITH_AES_128_CCM_8_SHA256` (RFC 8442, `D0 03`) para DTLS em dispositivos restritos: mantém o ECDHE e o CCM_8.
- **Recomendação.** A; C, se o grupo quiser cuidar dos dispositivos restritos.
- **O projeto hoje.** Segue A: só suítes com ECDHE e AEAD.

## D-B — O método do `trigger-auth` (DV-12)

- **Texto atual.** TS 29.503 V20.1.0, tab. 6.2.3.13.4.1-1 e §5.3.2.16.2: POST. `TS29503_Nudm_UECM.yaml` v1.5.0-alpha.2: `get:`, com `requestBody` obrigatório. O Anexo A.1 dá precedência ao YAML nos métodos.
- **Problema.** A regra de precedência leva a um GET com corpo obrigatório, que contraria a RFC 9110 §9.3.1 e o modelo das operações customizadas de TS 29.501 §4.6.1.2, sempre com POST.
- **Opções.**
  - **A (ideal).** No YAML, `post:` em `/{ueId}/registrations/trigger-auth`; e, no esquema `AuthTriggerInfo`, `required: [supi]`, como a tab. 6.2.3.13.4.2.2-1 (M) e TS 33.501 §14.2.6 (*"Input, Required: SUPI"*).
  - **B.** Trocar o texto para GET sem corpo, com o SUPI só no caminho. Muda o modelo das operações customizadas.
  - **C.** Uma exceção na regra do Anexo A.1. Deixa o erro no YAML.
- **Recomendação.** A.
- **O projeto hoje.** Segue A: POST; GET recebe 405.

## D-AA — O `iss` do *token* do CAPIF (DV-13)

- **Texto atual.** TS 29.222 V20.1.0, tab. 8.5.4.2.8-1: `iss` *"shall contain the API invoker Identifier as defined in Annex C.2.2 of 3GPP TS 33.122 [16] and in clause 4.1.1 of IETF RFC 7519 [24]"*. TS 33.122 V19.4.0 C.2.2 não tem `iss` e põe o *invoker* em `client_id`; a RFC 7519 §4.1.1 diz que `iss` é quem emite o JWT.
- **Opções.**
  - **A (ideal).** Na tabela e no YAML: `iss`, *"This attribute shall contain the identifier of the CAPIF core function that issued the access token, as defined in clause 4.1.1 of IETF RFC 7519 [24]"*; e um atributo novo, `client_id` (M), *"This attribute shall contain the API invoker Identifier as defined in Annex C.2.2 of 3GPP TS 33.122 [16]"*. As duas normas e a RFC passam a concordar.
  - **B.** Tirar `iss` e pôr `client_id`, exatamente como TS 33.122 C.2.2.
  - **C.** Pôr `iss` com o *invoker* em TS 33.122. Contraria a RFC 7519.
- **Recomendação.** A.
- **O projeto hoje.** `iss` e `client_id`, os dois com o identificador do *invoker*. Com A aprovada, o `iss` do CCF passaria a levar o identificador do CCF.

## D-AB — O `exp` do *token* do CAPIF como duração (DV-13)

- **Texto atual.** TS 29.222 V20.1.0, tab. 8.5.4.2.8-1: `exp`, do tipo `DurationSec`, *"the number of seconds after which the access_token is considered to be expired as defined in clause 4.1.4 of IETF RFC 7519 [24]"*. A RFC 7519 §4.1.4 e TS 33.122 C.6 tratam `exp` como instante.
- **Opções.**
  - **A (ideal).** Tipo `Uinteger` e a descrição *"This attribute shall contain the expiration time of the access token, as a NumericDate, as defined in clause 4.1.4 of IETF RFC 7519 [24]"*. A duração continua em `expires_in` da resposta.
  - **B.** Manter a duração e dizer que não é a `exp` da RFC 7519. Um JWT com `exp` fora da RFC confunde qualquer biblioteca.
- **Recomendação.** A.
- **O projeto hoje.** Segue A.

## D-AE — O HTTP/2 pelo *upgrade* do HTTP (DV-15)

- **Texto atual.** TS 29.222 V20.1.0 §7.3 e TS 29.122 V19.7.0 §5.2.2.1, que TS 29.522 §5.2 aplica ao NEF: *"A functional entity desiring to use HTTP/2 shall use the HTTP upgrade mechanism to negotiate applicable HTTP version as described in IETF RFC 9113 [10]."* A RFC 9113 §3.1 deprecou esse uso, e §3.2 manda o ALPN sobre TLS.
- **Opções.**
  - **A (ideal).** Nas duas normas: *"HTTP/2 over TLS shall be negotiated with the Application-Layer Protocol Negotiation (ALPN) extension, as described in clause 3.2 of IETF RFC 9113."*
  - **B.** Tirar a frase, e a RFC 9113 resolve.
  - **C.** Manter o *upgrade* só para conexões sem TLS, que as duas normas não admitem.
- **Recomendação.** A.
- **O projeto hoje.** ALPN, no CCF e no NEF.

## D-AL — A subscrição de *roaming* obrigatória no Stage 2 e opcional no Stage 3 (DV-2)

- **Texto atual.** TS 33.535 V19.0.0 §6.2.1 passo 5: *"Once receiving the request from the AF, the AAnF shall send a Nudm_EventExposure_Subscribe request to UDM with SUPI/GPSI to request the RoamingStatusReport from the UDM."* §6.3 passo 4 repete. TS 29.535 V20.1.0 tab. 5.1.8-1: `RoamingRestriction`, *feature* opcional.
- **Opções.**
  - **A (ideal).** Condicionar o passo: *"If the AAnF applies a roaming policy (see the RoamingRestriction feature in TS 29.535), the AAnF shall send a Nudm_EventExposure_Subscribe request …"*, e o mesmo em §6.3 passo 4. Um AAnF sem política de *roaming* não precisa do estado.
  - **B.** Tornar `RoamingRestriction` obrigatória em TS 29.535.
  - **C.** Uma NOTA em TS 29.535 sobre o AAnF sem a *feature*.
- **Recomendação.** A.
- **O projeto hoje.** Sem a subscrição (DV-2, mantido até haver uma VPLMN, Q5).

## D-AM — O AF_ID na notificação de serviço desligado

- **Texto atual.** TS 33.535 V19.0.0 §6.8.2 passo 5: *"The notification includes the A-KID and the AF_ID."* §7.1.6: *"Input, Required: A-KID"*, sem opcionais. TS 29.522 tab. 5.14.5.3.4-1: `ServiceDisableNotif` só com `aKId`. No mesmo passo 5: *"If AF(s) are determined at step 5"*, que é o próprio passo; o certo é o passo 4.
- **Opções.**
  - **A.** Alinhar o Stage 2 ao Stage 3: tirar o AF_ID do passo 5; o NEF acha o AF pelo URI de notificação que ele mesmo deu ao AAnF.
  - **B (ideal).** Alinhar o Stage 3 ao Stage 2: `afId` opcional na notificação do AAnF ao NEF, e AF_ID em *"Input, Optional"* de §7.1.6. O NEF consegue encaminhar mesmo com um URI só para vários AF.
- Nas duas opções, *"determined at step 4"*.
- **Recomendação.** B.
- **O projeto hoje.** Sem §6.8 (depende do DV-2).

## D-AN — GPSI no AKMA: a *feature* contra a nota

- **Texto atual.** TS 29.535 V20.1.0 §4.2.2.3.2: com `AKMA_GPSI_Support`, *"either the SUPI within the "supi" attribute or the GPSI within the "gpsi" attribute"*. Tab. 5.1.4.3.2-2, NOTA 2: *"For the "AkmaAfKeyData" data structure used in the current release of this specification, the "gpsi" attribute is not applicable and the "supi" attribute shall be included"*. §4.2.2.2.2: com a *feature*, o `AkmaKeyInfo` do registro leva *"SUPI as "supi" attribute or GPSI as "gpsi" attribute"*; TS 33.535 §6.1 passo 4 e §7.1.2 pedem o SUPI no registro, e o AUSF não tem o GPSI.
- **Opções.**
  - **A (ideal).** NOTA 2: *"When the "AKMA_GPSI_Support" feature is not supported, the "gpsi" attribute is not applicable and the "supi" attribute shall be included."* Em §4.2.2.2.2, só o SUPI; o GPSI vem do UDM (TS 33.535 §6.2.1 passos 3 e 4).
  - **B.** Admitir o GPSI no registro também em TS 33.535. O AUSF não tem de onde tirá-lo.
- **Recomendação.** A.
- **O projeto hoje.** Sem `AKMA_GPSI_Support`; o MVP-10 decide com este item.

## D-AO — Editoriais

- **TS 33.535 V19.0.0.** O nome da operação: `Naanf_AKMA_ServiceDisableNotification` (§6.8, §7.1.6) e `Naanf_AKMA_Context_Remove` (§6.6, §7.1.4), contra `Naanf_AKMA_Notify` e `Naanf_AKMA_ContextRemove` em TS 29.535. Tab. 7.3.1-1 dá *"Request/Response"* para `ServiceDisableNotification`, e a tab. 7.1.1-1, *"Notification/Response"*. A tab. 7.1.1-1 não lista o NEF entre os consumidores de `ServiceDisableNotification`, que §6.8.2 usa.
- **TS 29.522 V20.1.0.** "RoamingRestiction" em §4.4.23.2.
- **Proposta única.** Os nomes de TS 29.535 em TS 33.535; *"Notification/Response"* nas duas tabelas; *"AF, NEF"* na tab. 7.1.1-1; "RoamingRestriction".
