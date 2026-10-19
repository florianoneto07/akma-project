# Correções do free5GC para o AKMA

**Criado em:** 19 de outubro de 2026 · **Base:** free5GC v4.2.3 (NRF v1.4.5, UDM e AUSF v1.4.5, NEF v1.2.5, `openapi` v1.2.4)

Registro de tudo o que o projeto encontrou no free5GC e que precisa mudar para o AKMA funcionar conforme a norma. Pedido do responsável em 19/10/2026: este registro vai para o repositório público (MVP-7, plano seção 5). A regra de cada entrada:

- onde a norma diz o que deve ser feito (cláusula e trecho);
- como deve ser feito, segundo a norma;
- como o projeto trata hoje.

O AAnF não entra aqui: é uma NF nova, escrita pelo projeto (`aanf.md`). As NFs do free5GC só mudam pelos *patches* em `free5gc/patches/`.

| Nº | Onde | Assunto | Estado no projeto |
|---|---|---|---|
| F1 | NRF | descarta `aanfInfoList` | corrigido, *patch* `nrf/0001` |
| F2 | NRF | recusa descoberta de `AANF` por Routing Indicator | corrigido, *patch* `nrf/0002` |
| F3 | NRF | *token* com `aud` vazio | corrigido, *patch* `nrf/0003` |
| F4 | NRF | sem escopos por operação e por tipo de consumidor | aberto |
| F5 | `openapi` | `VerifyOAuth` não confere `exp` nem `aud` | contornado no AAnF; aberto nas NFs do free5GC |
| F6 | `openapi` | sem modelos nem cliente do Naanf_AKMA | corrigido, *patches* `openapi/0001` e `0002` |
| F7 | UDM | não envia `akmaInd` nem `routingId` | corrigido, *patch* `udm` |
| F8 | AUSF | não gera K_AKMA nem registra no AAnF | corrigido, *patch* `ausf` |
| F9 | *webconsole* | sem `akmaAllowed` e `routingId`; a edição os apaga | contornado por *script* |
| F10 | UDM | `id-translation-result` só traduz GPSI → SUPI | aberto (plano, seção 5) |
| F11 | NEF | sem a API `3gpp-akma` | em andamento (MVP-6) |
| F12 | NEF | 401 do *middleware* OAuth2 sem `ProblemDetails` nem `WWW-Authenticate` | aberto; mantido no MVP-6 por decisão do responsável |
| F13 | NEF | APIs do AF sem CAPIF e sem TLS mútuo | registrado como DV-9 |
| F14 | NEF | pede o *token* de `nnrf-nfm` para o tipo errado e não se descadastra | aberto |
| F15 | NRF | `complexQuery` sem `routing-indicator` para `AANF` | aberto |

---

## F1 — NRF descarta `aanfInfoList`

- **O que o free5GC faz.** `[FREE5GC]` `nnrfNFManagementOption` (`internal/context/management_data.go`) copia um a um os `*Info` que conhece (`udrInfo`, `udmInfo`, `ausfInfo`…) e descarta os demais. O perfil do AAnF fica sem os Routing Indicators que ele atende.
- **Onde a norma diz.** `[3GPP]` TS 29.510 tab. 6.1.6.2.2-1: `aanfInfoList` (mapa de `AanfInfo`) faz parte do `NFProfile`. §5.2.2.2.2: atributos que o NRF não conhece *"shall be stored as part of the NF's profile data in NRF"*.
- **Como deve ser.** O NRF guarda o `NFProfile` inteiro, inclusive `aanfInfoList`, e o devolve na descoberta.
- **No projeto.** *Patch* `free5gc/patches/nrf/0001-store-aanf-info-list.patch` (`nrf.md`, N1).

## F2 — NRF recusa a descoberta de `AANF` por Routing Indicator

- **O que o free5GC faz.** `[FREE5GC]` `validateQueryParameters` (`internal/sbi/processor/nf_discovery.go`) aceita só uma lista fixa de tipos de NF, sem `AANF`, e responde 400. O filtro de `routing-indicator` trata só AUSF e UDM.
- **Onde a norma diz.** `[3GPP]` TS 29.510 tab. 6.2.3.2.3.1-1: `routing-indicator` *"may be included if the target NF type is "AUSF", "AANF" or "UDM""*. Tab. 6.1.6.2.73-1 (`AanfInfo`): sem `routingIndicators`, *"the AAnf can serve any Routing Indicator"*. TS 33.535 §6.7: o AUSF, o NEF e os AF internos *"shall utilize the NRF to discover the AAnF instance(s) unless AAnF information is available by other means, e.g. locally configured"* e consideram o Routing Indicator do UE.
- **Como deve ser.** A descoberta aceita `target-nf-type=AANF`. Com `routing-indicator`, devolve os AAnF cujo `aanfInfoList` lista o RID ou não restringe RIDs.
- **No projeto.** *Patch* `free5gc/patches/nrf/0002-discover-aanf-by-routing-indicator.patch` (`nrf.md`, N2).

## F3 — NRF emite *token* com `aud` vazio

- **O que o free5GC faz.** `[FREE5GC]` `access_token.go` copia `targetNfInstanceId` para `aud`. Os clientes do free5GC não mandam esse campo, e o `aud` sai vazio.
- **Onde a norma diz.** `[3GPP]` TS 29.510 tab. 6.3.5.2.4-1: `aud` é obrigatório e *"shall contain the NF service producer's NF instance ID(s) (if the exact NF instance(s) of the NF service producer is known) or the NF type of NF service producers"*. TS 33.501 §13.4.1.1.2 passo 1: as *claims* incluem *"NF type of the NF Service Producer (audience)"*.
- **Como deve ser.** O NRF põe em `aud` a instância do produtor, quando o pedido a indica, ou o tipo de NF do produtor.
- **No projeto.** *Patch* `free5gc/patches/nrf/0003-access-token-audience.patch` (`nrf.md`, N3). O DV-8 foi encerrado com ele.

## F4 — NRF sem escopos por operação e por tipo de consumidor

- **O que o free5GC faz.** `[FREE5GC]` O NRF só emite escopos de serviço (por exemplo, `naanf-akma`). Ele não trata `allowedOperationsPerNfType` nem `allowedOperationsPerNfInstance` do `NFService`; nenhum código do NRF os menciona.
- **Onde a norma diz.** `[3GPP]`
  - TS 33.501 §13.4.1.1.1: o produtor pode registrar no seu perfil um *"additional scope"* com as operações permitidas *"per NF type of the NF Service Consumer or per NF instance ID"*.
  - TS 29.510 §6.1.6.2.3, NOTA 11: o NRF *"shall only grant such scope in the access token, if the scope is present in either "allowedOperationsPerNfType", for the specific NF type of the NF Service Consumer, or in "allowedOperationsPerNfInstance""*.
  - TS 29.535 tab. 5.1.9-1: escopos do Naanf_AKMA: `naanf-akma`, `naanf-akma:anchorkey`, `naanf-akma:applicationkeyget` e `naanf-akma:applicationkeyget:supi-access`. O último habilita a devolução do SUPI.
- **Como deve ser.**
  - O AAnF registra no `NFService` `naanf-akma`, por exemplo, `allowedOperationsPerNfType: {AUSF: [naanf-akma, naanf-akma:anchorkey], NEF: [naanf-akma, naanf-akma:applicationkeyget, naanf-akma:applicationkeyget:supi-access], AF: [...]}`.
  - O NRF só concede o escopo pedido se ele estiver na lista do tipo ou da instância do consumidor.
  - O AAnF confere o escopo de cada operação e só devolve o SUPI com `supi-access`.
- **No projeto.** Aberto. `[MVP]` O AAnF confere só `naanf-akma` (`aanf.md` §5).

## F5 — `oauth.VerifyOAuth` não confere `exp` nem `aud`

- **O que o free5GC faz.** `[FREE5GC]` `oauth.VerifyOAuth` (`openapi` v1.2.4, `oauth/oauth.go`) confere a assinatura e o escopo do *token*. Não confere a validade (`exp`) nem o destinatário (`aud`). Todas as NFs do free5GC que exigem *token* usam essa função, inclusive o *middleware* do NEF.
- **Onde a norma diz.** `[3GPP]` TS 33.501 §13.4.1.1.2 passo 2: o produtor *"checks that the audience claim in the access token matches its own identity or the NF type of NF Service Producer"* e *"checks that the access token has not expired by verifying the expiration time"*. TS 29.522 §7.2, para o NEF com CAPIF: o NEF confere *"whether the NEF identifier in the "token" matches its own published identifier"*.
- **Como deve ser.** `VerifyOAuth` recusa *token* vencido e *token* cujo `aud` não seja a instância nem o tipo do produtor. O 401 segue TS 29.500 §6.7.3 (F12).
- **No projeto.** O AAnF faz as duas verificações por conta própria (`aanf.md` §5). Nas NFs do free5GC, aberto.

## F6 — `openapi` sem os modelos e o cliente do Naanf_AKMA

- **O que o free5GC faz.** `[FREE5GC]` O `openapi` v1.2.4 tem o tipo de NF `AANF`, o serviço `naanf-akma` e `AanfInfo`. Não tem os tipos de dados do Naanf_AKMA nem um cliente.
- **Onde a norma diz.** `[3GPP]`
  - TS 29.535 §5.1.6.2: `AkmaKeyInfo` e `CtxRemove`;
  - TS 29.522 §5.14.5.3: `AkmaAfKeyRequest`, `AkmaAfKeyData` e `ServiceDisableNotif`;
  - TS 29.535 Anexo A: a API em OpenAPI, que prevalece na codificação (A.1).
- **Como deve ser.** Os modelos e o cliente são gerados de `TS29535_Naanf_AKMA.yaml` e `TS29522_AKMA.yaml` pela mesma ferramenta que o free5GC usa para as outras APIs.
- **No projeto.** *Patches* `free5gc/patches/openapi/0001-naanf-akma-models.patch` e `0002-naanf-akma-client.patch`, escritos à mão no estilo dos gerados (`openapi.md`, D6). Gerar pela ferramenta fica para o *upstream*.

## F7 — UDM não envia `akmaInd` nem `routingId`

- **O que o free5GC faz.** `[FREE5GC]` O UDM v1.4.5 não preenche os dois campos do `AuthenticationInfoResult`, embora o `openapi` os tenha.
- **Onde a norma diz.** `[3GPP]` TS 33.501 §6.1.3.1 passo 2 (EAP-AKA') e §6.1.3.2.0 passo 2 (5G AKA): *"If a subscriber has an AKMA subscription, the UDM shall include the AKMA indication and Routing indicator in the Nudm_UEAuthentication_Get Response."* TS 29.503 tab. 6.3.6.2.3-1: `routingId` *"shall be present if the akmaInd is set to "true""*.
- **Como deve ser.** Com `akmaAllowed: true` na assinatura (TS 29.505 tab. 5.4.2.2-1), o UDM responde com `akmaInd: true` e `routingId` igual ao Routing Indicator do SUCI.
- **No projeto.** *Patch* `free5gc/patches/udm/` (`udm.md`). O RID vem do `routingId` do UDR (decisão D10).

## F8 — AUSF não gera o material AKMA nem registra no AAnF

- **O que o free5GC faz.** `[FREE5GC]` O AUSF v1.4.5 ignora `akmaInd` e `routingId` da resposta do UDM.
- **Onde a norma diz.** `[3GPP]` TS 33.535 §6.1:
  - passo 3: com a indicação do UDM, o AUSF *"shall store the KAUSF and generate the AKMA Anchor Key (KAKMA) and the A-KID from KAUSF after the primary authentication procedure is successfully completed"*;
  - passo 4: seleciona o AAnF (§6.7) e manda SUPI, A-KID e K_AKMA com `Naanf_AKMA_AnchorKey_Register`;
  - *"The AUSF shall use the RID received from the UDM"*.

  Derivações em TS 33.535 A.2 e A.3.
- **Como deve ser.** O AUSF deriva K_AKMA e A-TID ao fim da autenticação primária bem-sucedida, monta o A-KID com o RID do UDM e registra no AAnF descoberto pelo NRF.
- **No projeto.** *Patch* `free5gc/patches/ausf/` (`ausf.md`). O registro não bloqueia a resposta ao AMF (decisão D9).

## F9 — *webconsole* sem os dados de assinatura AKMA

- **O que o free5GC faz.** `[FREE5GC]` A *webconsole* não tem campos para `akmaAllowed` e `routingId`. Ao editar o assinante, ela regrava o documento e apaga os dois.
- **Onde a norma diz.** `[3GPP]` TS 33.535 §3.1 e §4.2.5: os *"AKMA subscription data"* ficam no UDM. TS 29.505 tab. 5.4.2.2-1: `akmaAllowed` e `routingId` do `AuthenticationSubscription`; `routingId` *"shall be present if the akmaAllowed is set to "true""*.
- **Como deve ser.** A *webconsole* edita os dois atributos e os preserva ao editar o resto do assinante. Com `akmaAllowed: true`, exige `routingId` (`^[0-9]{1,4}$`).
- **No projeto.** Contornado por `free5gc/scripts/provision-akma-subscriber.sh` (`assinatura-akma.md`).

## F10 — UDM traduz só GPSI → SUPI

- **O que o free5GC faz.** `[FREE5GC]` `GetIdTranslationResultProcedure` (`internal/sbi/processor/subscriber_data_management.go`) trata o `ueId` sempre como GPSI e devolve o SUPI.
- **Onde a norma diz.** `[3GPP]` TS 29.503 §5.2.2.2.10: o consumidor (*"e.g. NEF, AAnF"*) recebe *"the SUPI/GPSI that corresponds to the provided GPSI/SUPI"*; o `ueId` *"shall be a SUPI or GPSI"*. TS 33.535 §6.3 passo 5: por política local, o NEF traduz o SUPI em GPSI pelo `Nudm_SubscriberDataManagement` para mandá-lo ao AF externo.
- **Como deve ser.** Com SUPI no `ueId`, o UDM devolve o GPSI do assinante (MSISDN ou identificador externo), levando em conta o AF quando os parâmetros do AF vierem (`af-id`, entre outros).
- **No projeto.** Aberto. O NEF do MVP-6 não devolve GPSI, o que a norma permite (*"optionally"*).

## F11 — NEF sem a API `3gpp-akma`

- **O que o free5GC faz.** `[FREE5GC]` O NEF v1.2.5 não tem nada de AKMA.
- **Onde a norma diz.** `[3GPP]` TS 33.535 §6.3 e §7.3.2; TS 29.522 §4.4.23.2 e §5.14: `POST {apiRoot}/3gpp-akma/v1/retrieve`. O NEF autoriza o AF, pede a K_AF ao AAnF com `Naanf_AKMA_ApplicationKey_Get` e nunca repassa o SUPI ao AF.
- **Como deve ser.** `nef.md`.
- **No projeto.** Em andamento, MVP-6.

## F12 — 401 do *middleware* OAuth2 do NEF sem `ProblemDetails` nem `WWW-Authenticate`

- **O que o free5GC faz.** `[FREE5GC]` `RouterAuthorizationCheck.Check` (`internal/util/util.go` do NEF) responde 401 com o corpo `{"error": "<mensagem>"}` em `application/json` e sem o cabeçalho `WWW-Authenticate`. Isso vale para todos os grupos de APIs do NEF: os do AF (`3gpp-traffic-influence`, `3gpp-pfd-management`, `3gpp-akma`) e os SBI (`nnef-pfdmanagement`, `nnef-callback`, `nnef-oam`).
- **Onde a norma diz.** `[3GPP]`
  - **APIs do AF.** TS 29.522 §5.14.7.1 remete a TS 29.122 §5.2.6. Tab. 5.2.6-1 de TS 29.122: o 401 leva corpo `ProblemDetails`, e a NOTA 2 manda o *Content-Type* `application/problem+json`.
  - **APIs SBI.** TS 29.500 §6.7.3: o produtor que recusa pedido sem *token* ou com *token* inválido *"shall return the HTTP "401 Unauthorized" status code together with the "WWW-Authenticate" header"*, com o esquema `Bearer`, `realm` igual ao URI da API e `error="invalid_token"` quando havia *token*.
  - **RFC 6750 §3:** o cabeçalho `WWW-Authenticate` na resposta a pedido sem credencial válida.
- **Como deve ser.** O 401 leva:
  - `WWW-Authenticate: Bearer realm="<URI da API>"`, com `error="invalid_token"` quando o pedido trazia *token*;
  - corpo `ProblemDetails` (`status: 401` e `detail`) em `application/problem+json`.
- **No projeto.** Aberto. Mantido no MVP-6 por decisão do responsável (19/10/2026), porque corrigir muda todas as APIs do NEF.

## F13 — APIs do AF no NEF sem CAPIF e sem TLS mútuo

- **O que o free5GC faz.** `[FREE5GC]` O free5GC não tem CAPIF. O NEF autoriza o AF só com *token* do NRF e serve as APIs do AF no servidor SBI, com TLS só de servidor (`sbi.scheme: https`) ou sem TLS.
- **Onde a norma diz.** `[3GPP]`
  - TS 33.501 §12.2: autenticação mútua NEF–AF *"based on client and server certificates"*, com TLS;
  - TS 33.501 §12.3: TLS *"shall"* proteger a interface; o suporte é obrigatório;
  - TS 33.501 §12.4: autorização OAuth *"shall"*;
  - TS 29.522 §6: sem CAPIF, autorização OAuth2 com *"Client Credentials"*, *"based on local configuration"*. A norma não diz quem é o servidor de autorização nesse caso (`[LACUNA]`);
  - TS 29.522 §7.2 e TS 33.122 §6.5.2.3: com CAPIF, o servidor de autorização é o CAPIF Core Function. O NEF confere no *token* o próprio identificador, o nome da API e o escopo.
- **Como deve ser.** TLS com certificados de cliente e servidor entre o NEF e o AF. Autorização por OAuth2: com CAPIF, *token* do CAPIF Core Function, obtido por `Obtain_Authorization` (TS 29.222 §5.6.2.3.2); sem CAPIF, de um servidor de autorização definido por configuração local.
- **No projeto.** DV-9: `http`, *token* do NRF e política local AF_ID ↔ endereço de origem no NEF.

## F14 — NEF pede o *token* de `nnrf-nfm` para o tipo errado

- **O que o free5GC faz.** `[FREE5GC]` `DeregisterNFInstance` (`internal/sbi/consumer/nrf_service.go` do NEF) pede o *token* de `nnrf-nfm` com `targetNfType=NEF`. O NRF recusa (`invalid_scope`), e o NEF não se descadastra ao encerrar. A falha já aparece com o NRF original.
- **Onde a norma diz.** `[3GPP]` TS 33.501 §13.4.1.1.2 passo 1: o pedido de *token* *"shall include the NF type of the expected NF Service Producer instance"*. O produtor de `nnrf-nfm` é o NRF.
- **Como deve ser.** `GetTokenCtx(ServiceName_NNRF_NFM, NrfNfManagementNfType_NRF)`, como nas outras NFs do free5GC.
- **No projeto.** Aberto. Não é do AKMA, mas afeta o NEF que o MVP-6 usa: depois de cada parada, o perfil antigo fica no NRF.

## F15 — `complexQuery` sem `routing-indicator` para `AANF`

- **O que o free5GC faz.** `[FREE5GC]` A descoberta por `complexQuery` (trecho `[Query-27]` de `nf_discovery.go`) não trata `routing-indicator` para `AANF`, como já não trata para os tipos que não são AUSF nem UDM.
- **Onde a norma diz.** `[3GPP]` TS 29.510 §6.2.3.2.3.1: o NRF *"may support"* a expressão `complexQuery` (TS 29.501), que combina os parâmetros de consulta da tab. 6.2.3.2.3.1-1, entre eles `routing-indicator`. Quem não a suporta *"shall reject"* o pedido com `INVALID_QUERY_PARAM`.
- **Como deve ser.** O filtro de F2 também vale dentro de `complexQuery`; ou, se o NRF não tratar o parâmetro, recusa o pedido com `INVALID_QUERY_PARAM`.
- **No projeto.** Aberto. Nenhum consumidor do projeto usa `complexQuery`.

---

## Integração sem cláusula normativa

Mudanças que o AKMA exige no free5GC sem que uma norma as determine:

| Onde | O quê | No projeto |
|---|---|---|
| `run.sh`, `force_kill.sh` | iniciar e parar o AAnF com as outras NFs | *patch* `free5gc/patches/run/` (D15, `logs-akma.md`) |
| `config/` | `aanfcfg.yaml`; no MVP-6, `3gpp-akma` no `nefcfg.yaml` e o N2 e o N3 fora do *loopback* | `free5gc/patches/config/` |
| `test/` | testes AKMA no banco de testes Go | `free5gc/patches/test/` (`verificacao.md`) |
