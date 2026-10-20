# NEF — API `3gpp-akma`: especificação

**MVP:** MVP-6 · **Data:** 19 de outubro de 2026 · **Estado:** aprovada pelo responsável em 19/10/2026, com as decisões da seção 9. Implementada e verificada (seções 10 e 11).

Este documento define o que muda no NEF do free5GC, no AAnF e no AF para que um AF fora da rede do operador obtenha a K_AF pelo NEF (TS 33.535 §6.3). O *patch* no NEF, a mudança no AAnF e o modo externo do AF são implementados a partir dele. As decisões P11 (b) e DV-9 (b) já foram tomadas (plano, MVP-6); os pontos ainda abertos estão na seção 9.

## 1. Situação normativa e do free5GC

`[3GPP]` TS 33.535 §6.3, o pedido de K_AF via NEF:

1. o AF descobre a HPLMN pelo A-KID e manda ao NEF o A-KID e o AF_ID;
2. o NEF autoriza o AF, inclusive pela verificação do AF_ID, e seleciona o AAnF (§6.7);
3. o NEF manda `Naanf_AKMA_ApplicationKey_Get` ao AAnF com o A-KID;
4. o AAnF responde ao NEF com a K_AF, a validade e o SUPI;
5. o NEF repassa ao AF a K_AF e a validade e, por política local, o GPSI. *"The NEF shall not send the SUPI to the AF."*

`[3GPP]` TS 29.522 define a API do lado do AF:

- §4.4.23.2: `POST {apiRoot}/3gpp-akma/v1/retrieve`, com `AkmaAfKeyRequest`. Um `403` do AAnF com `K_AKMA_NOT_PRESENT` *"shall"* ser repassado ao AF;
- §5.14: a operação personalizada `Retrieve` (tab. 5.14.3.1-1), os tipos `AkmaAfKeyRequest` e `AkmaAfKeyData` (tab. 5.14.5.3.2-1 e 5.14.5.3.3-1), os erros (tab. 5.14.7.3-1) e a *feature* `RoamingRestriction` (tab. 5.14.6-1);
- §5.14.7.1: o tratamento de erros segue TS 29.122 §5.2.6.

`[FREE5GC]` O NEF v1.2.5 do free5GC v4.2.3:

- não tem nada de AKMA;
- serve as APIs do lado do AF (`3gpp-traffic-influence`, `3gpp-pfd-management`) no mesmo servidor das APIs SBI, montando só os serviços listados em `serviceList`;
- protege cada grupo de APIs com um *middleware* OAuth2, que confere o *token* do NRF para o nome do serviço quando o NRF tem `oauth: true` (`internal/sbi/server.go`, `internal/util/util.go`);
- usa o `openapi` v1.2.4, que não tem os modelos nem o cliente do Naanf_AKMA.

`3gpp-akma` já é um `ServiceName` do `openapi` v1.2.4 (TS 29.510).

## 2. Endereços e caminhos

```
                  VM do AF (.129)                         VM do free5GC (.130)
  AF externo ──(1) registro e token, pelo túnel SSH ──► NRF 127.0.0.10:8000
  AF externo ──(2) POST /3gpp-akma/v1/retrieve ───────► NEF 192.168.16.130:8000
                                                        NEF ──(3) descoberta do AAnF ──► NRF
                                                        NEF ──(4) retrieve-applicationkey ──► AAnF 192.168.16.130:8010
```

| Caminho | Hoje | No MVP-6 |
|---|---|---|
| AF → NEF | o NEF escuta em 127.0.0.5:8000, fora do alcance da VM do AF | `sbi.registerIPv4` e `bindingIPv4` em 192.168.16.130. A porta continua a 8000 do free5GC |
| AF → AAnF (AF interno) | 192.168.16.130:8000 | **192.168.16.130:8010** (decisão do ponto 2) |
| AF → NRF | o NRF escuta em 127.0.0.10:8000; o AF não o alcança | túnel SSH a partir da VM do AF (P12, seção 2.1) |
| NEF → AAnF | — | descoberta pelo NRF com o RID do A-KID (NRF corrigido no MVP-1b) e *token* `naanf-akma` |

**Porta do AAnF.** Decisão do responsável (19/10/2026): o AAnF sai da 8000, que é a porta SBI de todas as NFs do free5GC, e vai para uma porta que o free5GC não usa. Nenhuma NF do free5GC muda de porta. As portas nas configurações do free5GC v4.2.3 são 2001, 2121, 2122, 3868, 3869, 4860, 5000, 8000, 8084, 8888, 9091, 20000, 27017 e 38412; mais as do UPF (2152 do GTP-U e 8805 do PFCP). A escolhida é a **8010**.

O AAnF muda em `aanfcfg.yaml` (`sbi.port`), e o AF em `afcfg.yaml` (`aanfUri`). O AUSF e o NEF não mudam, porque descobrem o AAnF pelo NRF. Os roteiros já leem o endereço das configurações. O tutorial e os documentos passam a dizer 8010.

### 2.1 De onde vem o *token* do AF externo (P12)

O que a norma diz:

- `[3GPP]` TS 33.501 §13.4.1.1.1: o NRF é o servidor de autorização OAuth 2.0 para as NFs consumidoras **dentro da PLMN**. O AF externo não é uma delas.
- `[3GPP]` TS 29.522 §7.2, TS 33.122 §6.5.2.3 e TS 33.501 §12.5: **com CAPIF**, o servidor de autorização é o CAPIF Core Function. O AF, como *API invoker*, obtém o *token* com `Obtain_Authorization` (TS 29.222 §5.6.2.3.2). O NEF confere no *token* o próprio identificador, o nome da API e o escopo.
- `[3GPP]` TS 29.522 §6: **sem CAPIF**, *"the access to the NEF Northbound APIs shall be authorized by means of OAuth2 protocol [...] based on local configuration, using the "Client Credentials" authorization grant as defined in clause 12.4 of 3GPP TS 33.501"*. TS 33.501 §12.4 exige a autorização OAuth.
- `[LACUNA]` Sem CAPIF, nenhuma das normas diz quem é o servidor de autorização: fica para a configuração local. Também não dizem como o AF alcança esse servidor na rede.

`[FREE5GC]` O free5GC não tem CAPIF. Nas APIs do AF que já existem, o NEF aceita só *tokens* do NRF.

`[INTERPRETAÇÃO]` Usar o NRF como servidor de autorização do AF externo cabe no *"based on local configuration"* de TS 29.522 §6, e é o que o free5GC já faz. Fica registrado em F13 de `correcoes-free5gc.md`: o caminho da norma para um AF externo é o CAPIF.

`[MVP]` O NRF escuta em *loopback*. Um túnel SSH dá ao AF da VM 192.168.16.129 o mesmo `nrfUri` das NFs, `http://127.0.0.10:8000`, sem mudar o NRF:

```bash
ssh -N -L 127.0.0.10:8000:127.0.0.10:8000 woca@192.168.16.130
```

## 3. Operação `retrieve` no NEF

### 3.1 Pedido

`POST http://192.168.16.130:8000/3gpp-akma/v1/retrieve`, com `Content-Type: application/json` e `Authorization: Bearer <token>`. O *token* é do NRF, com escopo `3gpp-akma`.

| Atributo de `AkmaAfKeyRequest` | P | No MVP-6 |
|---|---|---|
| `afId` | M | AF_ID na forma de transporte, `<FQDN>.<10 hex>` (ADR-0002) |
| `aKId` | M | A-KID |
| `anonInd` | O | repassado ao AAnF; `false` se omitido |
| `notifUri` | O | da *feature* `RoamingRestriction`, não suportada: ignorado, como no AAnF |
| `suppFeat` | C | sem *feature* suportada; a resposta não tem `suppFeat` |

### 3.2 Processamento, na ordem

1. **Token.** O *middleware* OAuth2 do free5GC, com `ServiceName_3GPP_AKMA`. Ausente ou inválido → 401.
   - `[FREE5GC]` O corpo é o do *middleware* existente, `{"error": "..."}` em `application/json`, sem `WWW-Authenticate`.
   - `[3GPP]` TS 29.122 tab. 5.2.6-1 pede `ProblemDetails` em `application/problem+json`; TS 29.500 §6.7.3 e RFC 6750 §3 pedem `WWW-Authenticate: Bearer`.
   - Decisão do responsável (19/10/2026): mantido no MVP-6 e documentado em F12 de `correcoes-free5gc.md`. A correção muda todas as APIs do NEF.
   - `[FREE5GC]` O *middleware* também não confere `exp` nem `aud` (F5).
2. **Forma do pedido** (TS 29.122 tab. 5.2.6-1, com `ProblemDetails` em `application/problem+json`):
   - `Content-Type` que não seja JSON → 415;
   - corpo que não é JSON → 400;
   - `afId` ou `aKId` ausente → 400, com `invalidParams`;
   - `afId` ou `aKId` malformado → 400, com `invalidParams`. A validação usa `shared/akma/afid` e `shared/akma/akid`, as mesmas do AAnF.
3. **Autorização do AF** (§6.3 passo 2; DV-9). O AF_ID tem de estar em `akma.authorizedAfs` do NEF, e o endereço de origem da conexão em `sourceAddresses` desse AF_ID. Senão → 403 com `ProblemDetails` sem `cause`. `[INTERPRETAÇÃO]` Nem TS 29.522 nem TS 29.122 definem causa para o caso; é o mesmo critério do AAnF (P3 do plano).
4. **Seleção do AAnF** (§6.7). O NEF descobre o AAnF no NRF com `target-nf-type=AANF`, `requester-nf-type=NEF`, `service-names=naanf-akma` e `routing-indicator` igual ao RID do A-KID. Nenhum AAnF → 503.
5. **Pedido ao AAnF** (§6.3 passo 3). `POST {apiRoot do AAnF}/naanf-akma/v1/retrieve-applicationkey` com `afId`, `aKId` e `anonInd`, e *token* `naanf-akma` obtido pelo NEF no NRF.
6. **Resposta ao AF**, pela tabela da seção 3.3.

### 3.3 Respostas

| Resposta do AAnF | Resposta do NEF ao AF | Base |
|---|---|---|
| 200 `AkmaAfKeyData` | 200 com `kaf` e `expiry`. **Nunca** `supi`; sem `gpsi` no MVP-6 | `[3GPP]` §6.3 passo 5; TS 29.522 §4.4.23.2 e tab. 5.14.5.3.3-1, NOTA 2 |
| 403 `K_AKMA_NOT_PRESENT` | 403 `K_AKMA_NOT_PRESENT` | `[3GPP]` TS 29.522 §4.4.23.2: *"shall relay"* |
| 403 sem `cause` (o NEF ou o AF_ID fora da política do AAnF, seção 4) | 403 sem `cause` | `[INTERPRETAÇÃO]` §4.4.23.2 pede só *"a proper error status code"* |
| outro 4xx | 500 | `[INTERPRETAÇÃO]` o pedido foi montado pelo NEF; um 4xx do AAnF é falha do NEF |
| 5xx, tempo esgotado, conexão recusada | 503 | `[INTERPRETAÇÃO]` |
| 204 | 204 | `[INTERPRETAÇÃO]` decisão do ponto 3, abaixo |

**O 204** (ponto 3, decisão do responsável em 19/10/2026: emitir e documentar).

- `[3GPP]` TS 29.522 tab. 5.14.3.2.2-2: o NEF responde 204 *"if the requested data does not exist"*.
- `[3GPP]` Quem sabe se o dado existe é o AAnF, que tem a mesma regra em TS 29.535 §4.2.2.3.2.
- `[INTERPRETAÇÃO]` O caso em que o NEF emite 204: quando o AAnF responde 204. O NEF repassa a resposta, sem corpo.
- Com o AAnF deste projeto, o caso não ocorre: ele nunca emite 204 (DV-7) e responde `403 K_AKMA_NOT_PRESENT`, que o NEF repassa. Com outro AAnF que emita 204, o NEF segue a tabela de TS 29.522. Por isso o NEF não tem desvio, e o DV-7 continua só no AAnF.

O NEF não emite 307 nem 308: há um NEF só.

### 3.4 Acesso anônimo

`[3GPP]` TS 33.535 §6.2.2 e TS 29.535 §4.2.2.3.2: com `anonInd: true`, o AAnF responde sem `supi`. TS 29.522 tab. 5.14.5.3.3-1, NOTA 3: com `anonInd: true` no pedido do AF, a resposta do NEF não tem `gpsi`.

O NEF repassa o `anonInd` como veio. A resposta ao AF é a mesma da seção 3.3, porque o MVP-6 nunca devolve `gpsi`.

## 4. Autorização do NEF no AAnF (P11, opção b)

`[3GPP]` TS 33.535 §6.2.1 passo 2: o AAnF verifica o pedido *"based on the configured local policy or based on the authorization information available in the signalling (i.e., Oauth2.0 token)"*. §6.3 passo 3: pelo NEF, o AAnF processa o pedido como no §6.2.1.

Nova lista na configuração do AAnF:

```yaml
configuration:
  akma:
    authorizedNefs: # NEFs that may request K_AF for external AFs (TS 33.535 cl. 6.3)
      - sourceAddresses:
          - 192.168.16.130
        afIds: # external AF_IDs accepted from this NEF
          - af1.akma.example.FF00000001
```

Regra no `retrieve-applicationkey`:

1. Se houver `Authorization`, o *token* é verificado como no `register-anchorkey` (`aanf.md` §5): escopo `naanf-akma`, `exp`, `aud`. Inválido → 401.
2. Com *token* válido e o endereço de origem numa entrada de `authorizedNefs`, o pedido é tratado como vindo do NEF. O `afId` tem de estar nos `afIds` dessa entrada; senão → 403 sem `cause`.
3. Nos demais casos vale `authorizedAfs`, como hoje: AF_ID ↔ endereço de origem.

Sem `authorizedNefs`, o comportamento é o do MVP-5.

### 4.1 O que a norma diz sobre reconhecer o NEF (ponto 5)

A norma não tem um passo próprio para o AAnF reconhecer o NEF. §6.3 passo 3 manda o AAnF tratar o pedido do NEF como o do AF (§6.2.1). O reconhecimento vem dos mecanismos gerais entre NFs:

| Mecanismo | Norma | O que dá ao AAnF |
|---|---|---|
| Autenticação entre NFs | `[3GPP]` TS 33.501 §13.3.2.1: com TLS, a autenticação é a do TLS, com os certificados das NFs. Sem TLS, *"authentication between NFs within one PLMN may be implicit by NDS/IP or physical security"* | quem está do outro lado da conexão |
| Autorização estática | `[3GPP]` TS 33.501 §13.3.0: sem *token*, o produtor *"shall check authorization of the NF Service Consumer based on its local policy"* | a política local |
| Autorização por *token* | `[3GPP]` TS 33.501 §13.4.1.1.2: o *token* traz a instância do consumidor (`sub`), o tipo do produtor (`aud`) e os escopos. O NRF só o emite depois de conferir a instância e o tipo do consumidor | que o NRF autorizou o consumidor para `naanf-akma` |
| Escopos por tipo de consumidor | `[3GPP]` TS 33.501 §13.4.1.1.1, TS 29.510 §6.1.6.2.3 NOTA 11, TS 29.535 tab. 5.1.9-1: o AAnF pode limitar `naanf-akma:applicationkeyget` e `naanf-akma:applicationkeyget:supi-access` a tipos de NF | quais operações o tipo do consumidor pode usar |

`[FREE5GC]` O NRF do free5GC não trata os escopos por tipo de consumidor (F4 de `correcoes-free5gc.md`). O laboratório não tem TLS (DV-1).

`[MVP]` A proposta da seção 4 usa o que a norma permite sem TLS:

- o endereço de origem como a autenticação implícita de §13.3.2.1;
- o *token* do NRF (§13.4.1.1.2);
- a lista de AF_ID como política local (§13.3.0; §6.2.1 passo 2).

Na VM 192.168.16.130, AUSF e NEF saem do mesmo endereço. Uma NF local com *token* `naanf-akma` poderia se passar pelo NEF, mas só para os AF_ID da lista. A forma da norma para fechar isso é o TLS com certificados de NF (DV-1) e os escopos por tipo (F4).

## 5. AF no modo externo

Mudanças em `af/config/afcfg.yaml`; o detalhe vai para `af.md` na implementação:

| Campo | Uso |
|---|---|
| `nefUri` | API do NEF, `http://192.168.16.130:8000`. Exclusivo com `aanfUri`: com ele, o AF é externo |
| `nrfUri` | NRF para o registro e o *token*, `http://127.0.0.10:8000` pelo túnel |
| `anonymous` | `true` manda `anonInd: true`, direto ao AAnF ou via NEF. Padrão `false` |

No modo externo, o AF:

- se registra no NRF como `nfType: AF` na partida, como o `seed` faz;
- pede *token* `3gpp-akma` para o NEF e o renova antes do `exp`;
- se descadastra ao encerrar.

O Ua\* não muda: o mesmo AF_ID dá a mesma K_AF, direto ou via NEF.

## 6. Mudanças no NEF do free5GC

*Patch* em `free5gc/patches/nef/`, aplicado por `patch-nf.sh nef`, no molde dos *patches* do AUSF:

| Arquivo do NEF | O quê |
|---|---|
| `go.mod` | `replace` do `openapi` para a cópia corrigida (D6) e de `shared/akma`, como no AUSF |
| `pkg/factory/config.go` | serviço `3gpp-akma`, prefixo `/3gpp-akma/v1`, seção `akma.authorizedAfs` |
| `internal/sbi/server.go` | grupo `3gpp-akma` com o *middleware* OAuth2 |
| `internal/sbi/api_akma.go` | rota `POST /retrieve` |
| `internal/sbi/processor/akma.go` | seção 3.2 |
| `internal/sbi/consumer/aanf_service.go` | descoberta do AAnF e cliente Naanf_AKMA |
| `internal/context/` | política de AF_ID |
| `internal/logger/` | categoria `AKMA` |

Configuração (`nefcfg.yaml`, por *patch* em `free5gc/patches/config/`):

- `3gpp-akma` em `serviceList`;
- `akma.authorizedAfs`;
- o endereço 192.168.16.130 do NEF, na porta 8000 (seção 2);
- no `aanfcfg.yaml`, a porta 8010 do AAnF.

## 7. Logs

No formato do free5GC, categoria `AKMA` (D16):

| Nível | Mensagem |
|---|---|
| info | `K_AF request from AF_ID <afId> for A-KID <aKId>` (com `anonymous` quando for o caso) |
| info | `K_AF for A-KID <aKId> relayed to AF_ID <afId>, expires <expiry>` |
| warn | AF fora da política, AAnF não encontrado, erro do AAnF |

O NEF nunca registra a K_AF. `[MVP]` O SUPI também não: o NEF não precisa dele e não pode repassá-lo.

## 8. Verificação prevista

| Teste | Onde |
|---|---|
| processamento da seção 3.2 contra um AAnF falso: os códigos da seção 3.3, sem `supi`, `anonInd` repassado | testes Go do *patch* do NEF |
| regra da seção 4: NEF autorizado, AF_ID fora da lista, *token* inválido, sem `authorizedNefs` | testes Go do AAnF |
| critérios 1 a 6 do MVP-6 | `e2e.sh`, grupos novos |

## 9. Pontos para aprovação

Respostas do responsável em 19/10/2026:

| Nº | Ponto | Decisão |
|---|---|---|
| 1 | P12: como o AF externo obtém o *token* | **decidido:** no MVP-6, o NRF é o servidor de autorização (*"based on local configuration"*, TS 29.522 §6), com túnel SSH no laboratório (seção 2.1). O caminho da norma, o CAPIF, entra no MVP-8 (plano, seção 5; F13 de `correcoes-free5gc.md`) |
| 2 | Endereço do NEF | **decidido:** o AAnF vai para a porta 8010; o NEF fica na porta 8000 do free5GC e só troca o endereço de *loopback* por 192.168.16.130 (seção 2) |
| 3 | 204 de TS 29.522 tab. 5.14.3.2.2-2 | **decidido:** o NEF emite 204 quando o AAnF responde 204 (seção 3.3) |
| 4 | Códigos de erro `[INTERPRETAÇÃO]` da seção 3.3 | **decidido:** a tabela da seção 3.3 |
| 5 | Como o AAnF reconhece o NEF | **decidido:** endereço de origem, *token* e lista de AF_ID (seções 4 e 4.1) |
| 6 | 401 do *middleware* sem `ProblemDetails` | **decidido:** mantido no MVP-6 e documentado (seção 3.2; F12 de `correcoes-free5gc.md`) |

## 10. Implementação

| Onde | O quê |
|---|---|
| `free5gc/patches/nef/0001-akma-api.patch` | seções 3 e 6. Arquivos novos: `internal/context/akma.go`, `internal/sbi/api_akma.go`, `internal/sbi/consumer/aanf_service.go`, `internal/sbi/processor/akma.go` e `akma_test.go`. Alterados: `go.mod`, `go.sum`, `nef_context.go`, `logger.go`, `consumer.go`, `server.go`, `config.go` |
| `free5gc/patches/config/0002-nef-akma.patch` | `nefcfg.yaml`: endereço 192.168.16.130, `3gpp-akma`, `akma.authorizedAfs` com `af1` e a entrada de teste `af3` (aceita pelo NEF e não pelo AAnF, para o negativo E07) |
| `free5gc/aanf/` | `akma.authorizedNefs` e porta 8010 (seção 4) |
| `af/af/internal/nef`, `af/af/internal/nrf` | cliente da API do NEF e registro do AF no NRF (seção 5) |

`[FREE5GC]` Como os outros clientes do NEF, o cliente do AAnF usa `http.DefaultClient` e o *hook* de métricas do free5GC. É também o que permite simular o NRF e o AAnF com `gock` nos testes do NEF.

Como reproduzir, na máquina do free5GC, com o free5GC parado: `patch-nf.sh nef apply`, depois `patch-nf.sh config apply` (tutorial, parte 9.1).

## 11. Verificação

| O quê | Onde | Resultado em 19/10/2026 |
|---|---|---|
| seção 3.2, as respostas da seção 3.3, `anonInd` repassado, sem `supi` e sem chave no log | `akma_test.go` do NEF (12 casos); a suíte original do NEF continua passando | passa, também com `-race` |
| regra da seção 4 | `TestRetrieveViaNef`, `TestNefPolicy` e os testes de configuração do AAnF | passa, também com `-race` |
| cliente do NEF, registro no NRF, `anonymous` | testes de `af/af/internal/{nef,nrf,aanf,config}` | passa, também com `-race` |
| Ua\* do UE em outra máquina com o AF interno, o externo e o externo anônimo | `docs/af-side/teste-nef.md` | atendido |
| negativos da API do NEF (E00 a E12) e o Ua\* do UE com os AF extras (R2n) | `af/scripts/e2e.sh` | 83 verificações, nenhuma falha (execução `20261019-213720`; `docs/af-side/e2e.md` §4.1 e §5.2) |
