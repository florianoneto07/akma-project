# Renovação da K_AKMA pelo AAnF (`Nudm_UECM_AuthTrigger`): proposta

**MVP:** MVP-12 · **Data:** 23 de outubro de 2026 · **Estado:** **aprovada** pelo responsável em 28/10/2026, com os pontos K1 a K12 da seção 8 como propostos; **implementada** em 28/10/2026 (seção 10). As mudanças no free5GC estão em `correcoes-free5gc.md` como F16 a F20.

Pedido do responsável em 22/10/2026: implementar a decisão do defeito D-B do registro de defeitos, a operação `trigger-auth` com o método POST, num MVP próprio, depois da fase (d) do MVP-8. Este documento define o que muda no AAnF, no UDM, no AMF e nos modelos do `openapi` para que o AAnF peça uma nova autenticação primária e, com ela, uma nova K_AKMA.

## 1. Situação normativa

`[3GPP]` TS 33.535:

- §6.4.1: *"KAKMA shall be re-keyed by running a successful primary authentication as described in clause 6.1."*
- §6.4.4: *"the AAnF may decide to refresh the KAKMA based on the operator's local authentication policy by sending the Nudm_UECM_AuthTrigger Request message to the UDM. The UDM may further decide whether to trigger the primary authentication as defined in clause 6.1.5 of TS 33.501"*. É opcional (*may*), e a política é local.
- §6.4.2: a K_AF não muda sozinha quando muda a K_AKMA; o AF segue com a K_AF que tem até ela vencer (o que o R3 do `e2e.sh` já verifica).

`[3GPP]` TS 33.501 §6.1.5, a autenticação primária disparada pela rede doméstica, opcional para a rede doméstica e a servidora:

- passo 0b: na `Nudm_UECM_Registration`, o AMF *"shall provide a callback URI within the AMF registration for the UDM to create an implicit subscription"* para a `Nudm_UECM_Re-AuthenticationNotification`;
- passo 1: o AAnF manda `Nudm_UECM_AuthTrigger` com o SUPI; NOTA A: o UDM pode recusar pela própria política. *"as AAnF sets the AF key expiry based on operator's local authentication policy, no frequent AF key expiry can happen and there is no risk of signalling overload"*;
- passo 2: o UDM manda a `Nudm_UECM_Re-AuthenticationNotification` ao AMF com o SUPI;
- passo 3: o AMF decide pela política e pelo estado do UE (*"if the UE is under handover, or if the UE is already under authentication"*); recusa com uma causa ou aceita. Se aceita e não alcança o UE, marca a autenticação como pendente e a faz quando o UE voltar;
- passo 4: o AMF *"starts the primary authentication procedure as defined in clause 6.1.2"*.

§14.2.6 e §14.2.7: as duas operações têm só o SUPI como entrada. §6.9.4.2: depois de uma autenticação primária, as chaves NAS novas vêm do novo K_AMF e entram em uso por um NAS SMC. §6.9.4.4: a troca das chaves AS depois de uma autenticação *"may be used"*.

`[3GPP]` TS 29.503:

- §5.3.2.16 e §6.2.3.13: `trigger-auth`, operação customizada em `{apiRoot}/nudm-uecm/v1/{ueId}/registrations/trigger-auth`, com `AuthTriggerInfo` (o SUPI) no corpo. Respostas: `204`; `403 REAUTHENTICATION_NOT_ALLOWED`; `404 USER_NOT_FOUND` ou `CONTEXT_NOT_FOUND`;
- §5.3.2.15 e §6.2.5.6: o UDM faz `POST` no `reauthNotifyCallbackUri` com `ReauthNotificationInfo` (o SUPI; opcionais `oldGuami` e `ausfRediscoveryInd`). O AMF responde `204`, ou `202` quando o UE não está alcançável e a autenticação fica pendente; recusas `403 REAUTHENTICATION_NOT_ALLOWED`, `404 CONTEXT_NOT_FOUND`, `409`;
- `reauthNotifyCallbackUri` é um atributo opcional de `Amf3GppAccessRegistration` e de `AmfNon3GppAccessRegistration`.

`[DIVERGÊNCIA]` Defeito D-B: a tabela 6.2.3.13.4.1-1 mapeia `trigger-auth` para POST, e o `TS29503_Nudm_UECM.yaml` v1.4.2 declara `get:` com corpo obrigatório. O Anexo A.1 diz que o YAML *"takes precedence ... with respect to the encoding of information elements and methods"*. Seguir a tabela contraria essa regra. Ponto K2.

## 2. Situação no free5GC v4.2.3

`[FREE5GC]` Nenhuma parte da cadeia existe:

- **`openapi`** (v1.2.4 e v1.2.5): sem `AuthTriggerInfo`, sem `ReauthNotificationInfo`, sem o atributo `reauthNotifyCallbackUri` e sem a operação no cliente `Nudm_UEContextManagement`;
- **UDM**: sem a rota `trigger-auth`. Ele guarda na memória o registro do AMF (`UdmUeContext.Amf3GppAccessRegistration`) e já notifica o AMF pelo `deregCallbackUri` (`SendOnDeregistrationNotification`), o modelo a seguir. O perfil do UDM no NRF não tem `allowedNfTypes`: o NRF dá ao AAnF um *token* `nudm-uecm`;
- **AMF**: registra-se no UDM só com o `deregCallbackUri` e não tem a rota da notificação. A máquina de estados do GMM já tem a transição de `Registered` para `Authentication`, mas só a usa num `Registration Request` de atualização. `AuthenticationProcedure` pula a autenticação quando o contexto de segurança vale e, na falha, manda `Registration Reject`. Depois do SMC, o estado `ContextSetup` espera um `Registration Request` para responder;
- **AAnF**: só registra-se no NRF; não descobre outras NF, não pede *token* e não guarda quando cada contexto foi registrado.

**UE (UERANSIM).** Sem mudança prevista. O `nr-ue` aceita o `Authentication Request` em qualquer estado (`receiveAuthenticationRequest`). O `Security Mode Command` põe o contexto novo em uso, e o *patch* 0001 deriva ali a nova K_AKMA e o novo A-KID (`deriveAkmaKeys`). O AUSF, pelo *patch* do MVP-2, registra o novo contexto no AAnF, que substitui o anterior (TS 33.535 §6.1 NOTA 1a). Conferir na implementação.

## 3. Fluxo proposto

1. Um AF pede a K_AF de um A-KID (`Naanf_AKMA_ApplicationKey_Get`, direto ou pelo NEF). O AAnF responde como hoje. Se o contexto tem mais tempo que `kakmaRefresh` (K3), ele também pede ao UDM, em segundo plano, a renovação: descobre o UDM no NRF, obtém o *token* `nudm-uecm` e manda `POST .../{supi}/registrations/trigger-auth` com `{"supi": …}`.
2. O UDM aplica a sua política (K4), acha o registro do AMF do acesso 3GPP e faz `POST` no `reauthNotifyCallbackUri` com `ReauthNotificationInfo`. Responde ao AAnF com o resultado do AMF.
3. O AMF (K6), com o UE em `Registered` e conectado, responde `204` e começa a autenticação primária: `Authentication Request`, AUSF, `Authentication Response`, confirmação no AUSF. O AUSF deriva a nova K_AKMA e a registra no AAnF.
4. O AMF faz o NAS SMC com o novo K_AMF (TS 33.501 §6.9.4.2) e volta a `Registered`, sem o `ContextSetup` de um registro. A sessão PDU continua.
5. No `Security Mode Command`, o ME deriva a nova K_AKMA e o novo A-KID. Os pedidos de K_AF seguintes vêm com o A-KID novo; o AF segue aceitando a K_AF antiga até ela vencer (§6.4.2).

## 4. Componentes

| Componente | Mudança |
|---|---|
| `free5gc/patches/openapi` | modelos `AuthTriggerInfo` e `ReauthNotificationInfo`; `reauthNotifyCallbackUri` em `Amf3GppAccessRegistration` e `AmfNon3GppAccessRegistration`; operação `TriggerAuth` (POST) no cliente `Nudm_UEContextManagement` |
| `free5gc/patches/udm` (novo *patch*) | rota `POST /{ueId}/registrations/trigger-auth`; política (K4); notificação ao AMF |
| `free5gc/patches/amf` (novo diretório) | `reauthNotifyCallbackUri` no registro no UDM; rota `POST /namf-callback/v1/reauth/{supi}`; autenticação disparada pela rede: estado do UE, *flag* de pendência, NAS SMC e volta a `Registered` |
| `free5gc/aanf` | instante de registro de cada contexto; `kakmaRefresh` na configuração; descoberta do UDM e *token* no NRF; cliente `trigger-auth` |
| `free5gc/aanf/tools/seed` | comando `trigger-auth --supi`, no papel do AAnF, para os testes do UDM |
| `free5gc/scripts/patch-nf.sh` | alvo `amf` (já cobre qualquer NF em `NFs/`) |
| `af/scripts/e2e.sh`, `e2e-remote.sh` | passo novo (K11) |

## 5. Logs

Na categoria AKMA (D16), no nível *info*, sem chaves:

- AAnF: `K_AKMA of <SUPI> registered <t> ago: refresh requested from the UDM (TS 33.535 cl. 6.4.4)` e o resultado;
- UDM: `AuthTrigger for <SUPI> from <tipo de NF>: re-authentication notified to the AMF <id>` ou a recusa e a causa;
- AMF: `Re-authentication of <SUPI> requested by the UDM: started` ou `pending (UE not connected)` ou a recusa;
- AUSF e UE: os logs de hoje, com o A-KID novo.

## 6. Verificação prevista

| Caso | Esperado |
|---|---|
| K_AKMA mais velha que `kakmaRefresh`, pedido de K_AF | o AF recebe a K_AF; o AAnF pede o `trigger-auth`; o UDM notifica o AMF; o UE se autentica de novo e tem um A-KID novo, igual ao do AUSF e ao do AAnF |
| sessão PDU depois da renovação | continua: o Ua\* do UE com o A-KID novo funciona pelo plano de usuário |
| AF com a K_AF anterior | aceita até vencer (§6.4.2) |
| segundo pedido do AAnF para o mesmo contexto | nenhum `trigger-auth` novo |
| `trigger-auth` repetido dentro do intervalo do UDM | `403 REAUTHENTICATION_NOT_ALLOWED` |
| SUPI desconhecido; UE sem registro de AMF; corpo sem `supi` | `404 USER_NOT_FOUND`; `404 CONTEXT_NOT_FOUND`; `400` |
| `GET` em `trigger-auth` | `405` (K2) |
| UE não conectado | `202` e autenticação no retorno do UE; nos testes de unidade do AMF, porque o UERANSIM não fica em CM-IDLE no laboratório |
| `kakmaRefresh` desligado (padrão) | nada muda: o `e2e.sh` atual passa igual |
| logs | nenhuma chave |

## 7. Limitações

- Só acesso 3GPP: o laboratório não tem não-3GPP.
- Sem a troca das chaves AS (TS 33.501 §6.9.4.4, opcional): o UERANSIM não tem segurança AS. O gNB segue com o K_gNB do K_AMF anterior até o próximo estabelecimento do contexto AS.
- `oldGuami` e `ausfRediscoveryInd` não são usados: o laboratório tem um AMF e um AUSF.

## 8. Pontos para aprovação

| Nº | Ponto | Proposta | Base |
|---|---|---|---|
| K1 | Escopo | a cadeia inteira: AAnF, UDM, AMF e modelos do `openapi`, só no acesso 3GPP; o UE sem mudança | TS 33.535 §6.4.4; TS 33.501 §6.1.5 |
| K2 | Método do `trigger-auth` | POST, como a tabela 6.2.3.13.4.1-1 e o defeito D-B decidiram; `GET` recebe `405`. Como contraria a precedência do Anexo A.1, registrar como DV-12 em `desvios.md` | `[DIVERGÊNCIA]` TS 29.503 tabela 6.2.3.13.4.1-1 e §5.3.2.16.2 (*"POST request"*) contra o YAML e o Anexo A.1; TS 29.501 §4.6.1.2; RFC 9110 |
| K3 | Política do AAnF | `akma.kakmaRefresh` (duração; `0`, o padrão, desliga). Num pedido de K_AF de um contexto mais velho que isso, o AAnF responde e pede o `trigger-auth` em segundo plano, uma vez por contexto; se o contexto não mudar em 5 minutos, pede de novo. Alternativas: (b) um temporizador por contexto, que renova também quem não usa AKMA; (c) só por comando de operação | `[LACUNA]` §6.4.4 deixa a política ao operador. A NOTA A de TS 33.501 §6.1.5 liga a renovação ao uso das chaves dos AF, como em (a) |
| K4 | Política do UDM | aceita e notifica o AMF; `403 REAUTHENTICATION_NOT_ALLOWED` para outro pedido do mesmo SUPI em menos de 60 s; `404 USER_NOT_FOUND` sem assinatura; `404 CONTEXT_NOT_FOUND` sem AMF registrado; `400` sem `supi`. Responde ao AAnF depois da resposta do AMF: `204` se o AMF aceitou (`204` ou `202`), a recusa do AMF repassada | `[LACUNA]` a norma não diz se o UDM espera o AMF; NOTA A (recusa pela política); D-B, observação secundária (`AuthTriggerInfo` sem `required`) |
| K5 | *Callback* do AMF | `reauthNotifyCallbackUri` = `<AMF>/namf-callback/v1/reauth/<SUPI>`, no registro do acesso 3GPP, como o `deregCallbackUri` | TS 33.501 §6.1.5 passo 0b; TS 29.503 §6.2.6.2.2 |
| K6 | AMF ao receber a notificação | UE em `Registered` e conectado: `204` e autenticação já; UE registrado e não conectado: `202`, autenticação pendente, feita depois do próximo `Service Request` ou registro aceito; UE em outro procedimento (autenticação, SMC, registro, desregistro): `403 REAUTHENTICATION_NOT_ALLOWED`; UE desconhecido: `404 CONTEXT_NOT_FOUND`. A falha da autenticação segue o tratamento atual do AMF | TS 33.501 §6.1.5 passo 3; TS 29.503 tabela 6.2.5.6-3 |
| K7 | Depois da autenticação | NAS SMC com o novo K_AMF e volta direta a `Registered`; sem troca das chaves AS | TS 33.501 §6.9.4.2 (*shall*) e §6.9.4.4 (*may*); `[MVP]` o UERANSIM não tem segurança AS |
| K8 | Modelos | no *patch* do `openapi` do projeto, como os do Naanf_AKMA; `oldGuami` e `ausfRediscoveryInd` lidos e ignorados | TS 29.503 §6.2.6.2.35, §6.2.6.2.36 |
| K9 | AAnF como consumidor | descobre o UDM no NRF e pede o *token* `nudm-uecm` quando o OAuth2 está ligado, como o AUSF faz | TS 29.510; TS 33.501 §13.4.1 |
| K10 | Ferramenta de teste | `seed trigger-auth --supi`, no papel do AAnF, para os casos do UDM sem esperar a política | como o `seed register` do MVP-1 |
| K11 | Verificação | passo novo do `e2e.sh`, depois do R6: o AAnF reiniciado com `kakmaRefresh: 20s`, um Ua\* do UE depois disso e as conferências da seção 6; a configuração do AAnF volta ao fim | — |
| K12 | Defeitos novos | registrar como D-O: TS 29.503 §5.3.2.16.2 passo 2a diz *"the NF Service Consumer responds with 204"*, quando quem responde é o UDM; e *"manadatory"* na nota da tabela 6.2.3.13.4.2.2-2, como em TS 29.535 (D-I) | editorial |

## 9. Arquivos previstos

- `free5gc/patches/openapi`, `free5gc/patches/udm`, `free5gc/patches/amf` (novo);
- `free5gc/aanf` (contexto, configuração, consumidor) e `free5gc/aanf/tools/seed`;
- `free5gc/config/aanfcfg.yaml` (`kakmaRefresh: 0`, comentado);
- `af/scripts/e2e.sh`, `e2e-remote.sh`;
- `docs/free5gc-side/` (este documento, `aanf.md`, `udm.md`, um `amf.md` novo), `desvios.md` (DV-12), registro de defeitos (D-O), plano, matriz, tutorial e guia público.

## 10. Implementação

Feita em 28/10/2026, nos arquivos da seção 9:

| Componente | Onde | Documento |
|---|---|---|
| modelos e clientes | `free5gc/patches/openapi/0003-home-network-triggered-authentication.patch` | `openapi.md` |
| UDM | `free5gc/patches/udm/0002-auth-trigger.patch` | `udm.md` §6 |
| AMF | `free5gc/patches/amf/0001-home-network-triggered-authentication.patch` | `amf.md` |
| AAnF | `free5gc/aanf`: `akma.kakmaRefresh`, `Registered` do contexto, consumidor do UDM | `aanf.md` |
| ferramenta | `seed trigger-auth` | `aanf.md` |
| roteiro | passo R6r de `af/scripts/e2e-remote.sh` | `docs/af-side/e2e.md` |

**Decisões do responsável na implementação (28/10/2026).**

1. *Token* na notificação ao AMF. `[3GPP]` TS 29.500 §6.7.3 (*"For the notify operation ... OAuth 2.0 access token is not used"*) e TS 33.501 §13.4.1.0. `[FREE5GC]` O AMF exige um *token* `namf-callback` em todas as rotas de *callback*. Decidido seguir a norma: o UDM manda a notificação sem *token*, e a rota nova do AMF fica fora da verificação. As outras notificações ficam como estão, registradas como F20.
2. Recusas do AMF. O K6 previa `403 REAUTHENTICATION_NOT_ALLOWED` para todo procedimento em curso. A tabela 6.2.5.6-3 de TS 29.503 tem causas próprias. Decidido segui-la: `409 TEMPORARY_REJECT_REGISTRATION_ONGOING` para registro em curso, `409 TEMPORARY_REJECT_HANDOVER_ONGOING` para *handover*, e `403` para autenticação ou desregistro em curso. Como o `trigger-auth` não lista o `409`, o UDM responde `403` ao AAnF nos dois casos, com a causa do AMF no `detail`.

**Outras diferenças em relação às seções 3 a 7.**

- O AMF escolhe um ngKSI novo para a autenticação (TS 33.501 §6.1.3.2 passo 6). O UE recusa um ngKSI em uso com `NGKSI_ALREADY_IN_USE`.
- O NAS SMC dessa autenticação vai com RINMR e HDP em zero. O free5GC punha HDP = 1 depois de um registro periódico ou de mobilidade, e o UE derivaria o K_AMF horizontalmente.
- O AMF ignora um `NAS message container` no `Security Mode Complete` dessa autenticação. O UERANSIM põe ali o último `Registration Request` quando esse registro foi mandado sem contexto de segurança, também fora de um registro; o AMF do free5GC o tratava como um registro novo. O UE continua sem mudança (K1).
- Se o AUSF não responde, o AMF não manda `Registration Reject`, que só responde a um registro. Responde `500 SYSTEM_FAILURE` ao UDM, que repassa `500` ao AAnF; o UE segue registrado com o contexto corrente, e o AAnF volta a pedir no primeiro pedido de K_AF depois de 5 minutos.
- A autenticação pendente começa logo depois do `Service Request` aceito ou da volta a `Registered` depois de um registro. Não entra no meio do procedimento.
- O registro do AMF sem `reauthNotifyCallbackUri`, de um AMF sem o *patch*, recebe `403 REAUTHENTICATION_NOT_ALLOWED` no UDM.
- O UDM só libera o SUPI do intervalo de 60 s quando o AMF recusa.
- O AAnF manda `User-Agent: AANF-<nfInstanceId>` (TS 29.500 tabela 5.2.2.2-1); o UDM tira dali o tipo de NF dos seus logs.
- O `seed trigger-auth` registra-se no NRF como um AF temporário, como os outros comandos do `seed`. Registrar-se como AANF poria uma segunda instância de AAnF na descoberta do AUSF.

**Limitações, além das da seção 7.**

- O UDR sem o *patch* do `openapi` não guarda o `reauthNotifyCallbackUri`. O UDM usa o registro que tem na memória; depois de reiniciar o UDM, o `trigger-auth` responde `404 CONTEXT_NOT_FOUND` até o UE se registrar de novo.
- O AMF não serializa os eventos do GMM que chegam por NAS e por SBI, como nas outras operações SBI do free5GC.

## 11. Verificação

Em 28/10/2026, na VM do free5GC (192.168.16.130), com o UE do UERANSIM no *namespace* e o AF na VM do AF. Execução `20261028-162316` do `e2e.sh`: **161 verificações, nenhuma falha, nenhuma pulada**. São as 141 da fase (d) do MVP-8 e as 20 do passo novo R6r (`docs/af-side/e2e.md` §5.7).

| Caso da seção 6 | Resultado |
|---|---|
| K_AKMA mais velha que `kakmaRefresh`, pedido de K_AF | atendido: o AF recebeu a K_AF; o AAnF pediu o `trigger-auth`, o UDM notificou o AMF, e o AMF autenticou o UE registrado e fez o NAS SMC [R6r.ua, R6r.aanf, R6r.udm, R6r.amf, R6r.smc]. O UE derivou um A-KID novo, igual ao do AUSF e ao do AAnF [R6r.ue, R6r.new-akid, R6r.akid-ausf, R6r.akid-aanf] |
| sessão PDU depois da renovação | atendido: nenhuma interface TUN nova, e o Ua\* seguinte passou pela mesma sessão com o A-KID novo [R6r.pdu, R6r.ua-new] |
| AF com a K_AF anterior | atendido: o AF segue com a K_AF que tem até ela vencer; foi por isso que o R6r pede, depois de 20 s, a K_AF de outro AF_ID. O caso já era coberto pelo R3.ua-old |
| segundo pedido do AAnF para o mesmo contexto | atendido: uma renovação só, nenhuma para a K_AKMA nova [R6r.once]; o intervalo de 5 minutos, nos testes de unidade (`TestRefreshDue`, `TestRefreshKAkma`) |
| `trigger-auth` repetido dentro do intervalo do UDM | atendido: `403 REAUTHENTICATION_NOT_ALLOWED` [R6r.udm-interval] |
| SUPI desconhecido; UE sem registro de AMF; corpo sem `supi` | atendido: `404 USER_NOT_FOUND` [R6r.udm-unknown]; `400 MANDATORY_IE_MISSING` [R6r.udm-no-supi]; o `CONTEXT_NOT_FOUND`, nos testes de unidade do UDM (`TestAuthTrigger`) |
| `GET` em `trigger-auth` | atendido: `405` com `Allow: POST` [R6r.udm-get] |
| UE não conectado | nos testes de unidade do AMF (`TestReauthNotificationIdle`), como previsto |
| `kakmaRefresh` desligado (padrão) | atendido: as 141 verificações anteriores passam com a configuração instalada, que tem `kakmaRefresh: 0` |
| logs | atendido: nenhuma cadeia de 64 hexadecimais nos logs [R7, L1] |

A cadeia levou 140 ms, do pedido do AAnF ao novo contexto no AAnF:

```
12:26:41.930 [AANF][AKMA] K_AKMA of imsi-208930000000001 registered 23s ago: refresh requested from the UDM (TS 33.535 cl. 6.4.4)
12:26:41.942 [AANF][AKMA] UDM for Nudm_UECM_AuthTrigger: http://127.0.0.3:8000 (discovered at the NRF)
12:26:42.013 [AMF][AKMA] Re-authentication of imsi-208930000000001 requested by the UDM: started (ngKSI 1)
12:26:42.014 [UDM][AKMA] AuthTrigger for imsi-208930000000001 from AANF at 127.0.0.1: re-authentication notified to the AMF bcb61676-…
12:26:42.015 [AANF][AKMA] Refresh of K_AKMA of imsi-208930000000001 accepted by the UDM: new primary authentication requested
12:26:42.054 [AUSF][AKMA] K_AKMA and A-KID derived for imsi-208930000000001: A-KID rid0000.atidfpSD…@5gc.mnc093.mcc208.3gppnetwork.org
12:26:42.056 [UE] AKMA key material derived, A-KID[rid0000.atidfpSD…@5gc.mnc093.mcc208.3gppnetwork.org]
12:26:42.057 [AMF][Gmm] NAS message container in the Security Mode Complete of a re-authentication: ignored
12:26:42.057 [AMF][AKMA] Re-authentication of imsi-208930000000001 requested by the UDM: completed, new KAMF in use (ngKSI 1)
12:26:42.070 [AANF][AKMA] AKMA context of imsi-208930000000001 registered, A-KID rid0000.atidfpSD…@5gc.mnc093.mcc208.3gppnetwork.org (previous context replaced)
```

**Duas execuções antes dessa**, com falhas que levaram a correções:

1. `20261028-151849`, 151 verificações e 10 falhas, todas no R6r. O roteiro pedia, depois de 20 s, a K_AF do mesmo AF_ID; o AF já a tinha guardada e não consultou o AAnF, e a renovação não começou. O R6r passou a pedir a K_AF de outro AF_ID (PSK TLS 1.3). O `seed trigger-auth` repetia o prefixo `/nudm-uecm/v1` no caminho e recebia 404; corrigido. O UDM mostrava `GO` como tipo de NF do `User-Agent` do Go; passou a aceitar só um tipo de NF.
2. `20261028-154345`, 160 verificações e 1 falha, o R6r.smc. O UERANSIM pôs o último `Registration Request` no `Security Mode Complete` (seção 10), e o AMF entrou em `ContextSetup`, onde ficou. O AMF passou a ignorar esse conteúdo durante a reautenticação (teste `TestSecurityModeCompleteOfReauth`).

**Testes de unidade**, todos com `-race`: UDM (`TestAuthTrigger`, `TestAuthTriggerInterval`, `TestAuthTriggerRequests`, `TestRequester`), AMF (`TestReauthNotificationRefusals`, `TestReauthNotificationIdle`, `TestReauthNotificationStart`, `TestReauthNotificationAbort`, `TestSecurityModeCommandOfReauth`, `TestSecurityModeCompleteOfReauth`) e AAnF (`TestRefreshDue`, `TestKakmaRefresh`, `TestTriggerAuth`, `TestRefreshKAkma`); o `openapi` com o `0003` compila nas duas versões. O banco de testes do free5GC, com o AMF e o UDM corrigidos, compila, e o `TestAkma` passa.

