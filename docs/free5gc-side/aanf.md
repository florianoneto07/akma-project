# AAnF — AKMA Anchor Function

**MVP:** MVP-1 · **Data:** 17 de outubro de 2026 · **Verificado na VM do free5GC** (192.168.16.130) e a partir da VM do AF (192.168.16.129) em 17/10/2026

## 1. O que é e onde fica

`[3GPP]` TS 33.535 §4.2.1: o AAnF é a função âncora do AKMA na HPLMN. Ele guarda K_AKMA e o A-KID que o AUSF lhe entrega depois da autenticação primária. Deriva e entrega a K_AF a cada AF autorizado e apaga o contexto quando pedido. O serviço é o Naanf_AKMA (TS 29.535).

`[MVP]` O AAnF é uma NF nova, com a estrutura das NFs do free5GC v4.2.3 (o AUSF v1.4.5 serviu de modelo: `cmd`, `pkg/factory`, `pkg/service`, `internal/context`, `internal/sbi/{processor,consumer}`). Roda na máquina do free5GC (ADR-0001-R2). **Nenhum arquivo do free5GC é alterado.**

| No repositório | Conteúdo |
|---|---|
| `free5gc/aanf/cmd/main.go` | ponto de entrada (`-c` configuração, `-l` log) |
| `free5gc/aanf/pkg/factory/` | leitura e validação de `aanfcfg.yaml` |
| `free5gc/aanf/pkg/service/init.go` | ciclo de vida: contexto, servidor SBI, registro e cancelamento no NRF |
| `free5gc/aanf/internal/context/context.go` | contexto da NF, verificação de token, política de AFs |
| `free5gc/aanf/internal/context/akma_context.go` | armazenamento dos contextos AKMA |
| `free5gc/aanf/internal/sbi/{server,api_akma}.go` | servidor h2c, rotas, erros de protocolo |
| `free5gc/aanf/internal/sbi/processor/` | as três operações |
| `free5gc/aanf/internal/sbi/consumer/` | Nnrf_NFManagement |
| `free5gc/aanf/internal/util/problem.go` | `ProblemDetails` |
| `free5gc/aanf/tools/seed/` | ferramenta de teste: papel do AUSF e do AF. Desde o MVP-5, `register` aceita `--gpsi` e `--k-akma`, para os testes negativos do `register-anchorkey` (`docs/af-side/e2e.md` §4) |
| `free5gc/config/aanfcfg.yaml` | configuração de exemplo do laboratório |
| `free5gc/scripts/{prepare-openapi,build-aanf,run-aanf,deploy}.sh` | preparação, compilação, execução, implantação |

| Na instalação do free5GC | Origem |
|---|---|
| `~/free5gc/bin/aanf` | arquivo novo (binário estático) |
| `~/free5gc/config/aanfcfg.yaml` | arquivo novo |
| `~/free5gc/log/aanf/` | criado pelo `run-aanf.sh` (log, saída, pid). Com `./run.sh -aanf` (MVP-3a, `logs-akma.md` §5), o log vai para o `free5gc.log` das outras NFs |
| `~/free5gc/cert/aanf.pem`, `cert/aanf_<id>.pem` | criados pelo NRF (`registro-nrf.md` §4) |

Dependências: `free5gc/openapi` v1.2.4 com os *patches* do projeto (`openapi.md`), `free5gc/util` v1.3.2, `shared/akma` (MVP-0).

## 2. Configuração

Além dos campos comuns às NFs do free5GC (`sbi`, `serviceNameList`, `nrfUri`, `nrfCertPem`, `plmnSupportList`, `logger`):

| Campo | Exemplo | Regra | Marcação |
|---|---|---|---|
| `nfInstanceId` | UUID v4 fixo | opcional; sem ele, gerado a cada partida | `[MVP]` fixo para não acumular certificados no NRF |
| `serviceNameList` | `[naanf-akma]` | só esse valor é aceito | `[3GPP]` TS 29.535 §5.1.1 |
| `routingIndicators` | `["0000"]` | cadeias de 1 a 4 dígitos; vazio = atende qualquer RID. Mesmo sem aspas, `0000` é lido como `"0000"` (teste `TestRoutingIndicatorLeadingZerosKept`) | `[3GPP]` TS 29.510 `AanfInfo`; TS 23.003 §2.2B |
| `akma.kafLifetime` | `1h` | duração Go positiva; padrão `1h` | `[3GPP]` política do operador (TS 33.535 §5.2); `[MVP]` o valor |
| `akma.authorizedAfs[].afId` | `af1.akma.example.FF00000001` | forma de transporte `<FQDN>.<10 hex>` | `[3GPP]` TS 29.522 `AfId` |
| `akma.authorizedAfs[].sourceAddresses` | `[192.168.16.129]` | ao menos um endereço IP | `[MVP]` política local (seção 5) |
| `akma.roamingStatusSubscription` | `false` | `true` impede a partida | DV-2 |
| `akma.aKid` | `atidEncoding`, `atidLabel`, `acceptLabels` | padrão: base64, escreve `atid`, aceita `atid` e `akid` | ADR-0002 §4; defeito D-A |

## 3. Operações

`[3GPP]` TS 29.535 Tab. 5.1.4.1-1: três operações personalizadas, todas `POST` em `{apiRoot}/naanf-akma/v1/`.

### 3.1 `register-anchorkey` — Naanf_AKMA_AnchorKey_Register

`[3GPP]` TS 29.535 §4.2.2.2.2; TS 33.535 §6.1 passos 4–5. Corpo `AkmaKeyInfo`; resposta 200 com o `AkmaKeyInfo` armazenado.

- `supi`, `aKId` e `kAkma` são obrigatórios. `[3GPP]` Tab. 5.1.6.2.2-1, NOTA: sem a *feature* `AKMA_GPSI_Support`, *"only the supi attribute shall be present"*. Faltando algum, a resposta é 400 `MANDATORY_IE_MISSING`, com `invalidParams`.
- `gpsi` presente → 400 `OPTIONAL_IE_INCORRECT`. `[INTERPRETAÇÃO]` A norma não indica a causa para um atributo que não deveria estar presente; esta é a de TS 29.500 Tab. 5.2.7.2-1 que mais se aproxima.
- `aKId` fora do formato NAI de TS 29.522, ou `kAkma` diferente de 64 caracteres hexadecimais → 400 `MANDATORY_IE_INCORRECT`. A codificação hexadecimal é decisão do ADR-0002.
- `supi` não é validado além de não ser vazio. `[3GPP]` O padrão de `Supi` em TS 29.571 admite qualquer cadeia não vazia (`…|.+`). O AAnF não usa o SUPI na KDF: a restrição a IMSI (DV-3) vale no AUSF e no UE.

### 3.2 `retrieve-applicationkey` — Naanf_AKMA_ApplicationKey_Get e _AnonUser_Get

`[3GPP]` TS 29.535 §4.2.2.3.2; TS 33.535 §6.2.1 e §6.2.2. Corpo `AkmaAfKeyRequest`. A ordem das verificações segue o §6.2.1:

1. formato: `afId` e `aKId` obrigatórios e bem formados (400);
2. passo 2, autorização do AF pela política local (seção 5) → 403 sem `cause`;
3. passo 2, presença de K_AKMA para o A-KID → 403 `K_AKMA_NOT_PRESENT`;
4. passo 7, K_AF: reaproveitada se o contexto já tem uma válida, derivada da K_AKMA corrente; senão derivada (TS 33.535 A.4) e guardada com o instante de expiração;
5. passo 8, resposta 200 com `kaf` (64 hexadecimais, minúsculas — ADR-0002) e `expiry` (RFC 3339, UTC, em segundos). Também vai `supi`, exceto quando `anonInd` é `true` (TS 29.535 Tab. 5.1.4.3.2-2 NOTA 2). `gpsi` nunca é enviado.

- **Nunca 204** (DV-7): sem K_AKMA a resposta é 403; com K_AKMA a K_AF é sempre derivável.
- **Nunca `ROAMING_AKMA_SERVICE_DENIED`** (DV-2): a subscrição de *roaming* está desligada e a *feature* `RoamingRestriction` não é suportada.
- **`notifUri` é ignorado.** `[3GPP]` TS 29.522 Tab. 5.14.5.3.2-1 o vincula à *feature* `RoamingRestriction`, não suportada. Nenhuma subscrição implícita é criada.

### 3.3 `remove-context` — Naanf_AKMA_ContextRemove

`[3GPP]` TS 29.535 §4.2.2.4.2; TS 33.535 §6.6.1. Corpo `CtxRemove` com `supi` (obrigatório nesta *release*, Tab. 5.1.6.2.3-1 NOTA). Apaga SUPI, A-KID, K_AKMA e as K_AF emitidas e responde 204. Sem contexto, responde 404 `AKMA_CONTEXT_NOT_FOUND`.

### 3.4 *Features*

`[3GPP]` TS 29.535 §5.1.8; TS 29.500 §6.6. O AAnF não suporta nenhuma das duas *features* opcionais (`AKMA_GPSI_Support`, `RoamingRestriction`). Quando o consumidor envia `suppFeat`, a resposta traz `suppFeat: "0"`, a interseção vazia. Sem `suppFeat` no pedido, o atributo não vai na resposta.

## 4. Armazenamento

Em memória (`internal/context/akma_context.go`); `[MVP]` perdido ao reiniciar o AAnF. Índices por SUPI e por A-KID, protegidos por *mutex*.

| Regra | Origem |
|---|---|
| O AAnF guarda o par mais recente enviado pelo AUSF; o novo par apaga o antigo, e o A-KID antigo deixa de identificar K_AKMA | `[3GPP]` TS 33.535 §6.1 passo 4 e NOTA 1a |
| As K_AF já emitidas permanecem no contexto até expirar | `[3GPP]` §5.2 |
| K_AF só é derivada quando o contexto não tem uma válida para o AF, derivada da K_AKMA corrente; a expiração guardada é devolvida até ser atingida | `[3GPP]` §6.2.1 passo 7; §5.2 NOTA |
| O A-KID é comparado pelo conteúdo: RID como cadeia (`"0000"` ≠ `"000"`), A-TID em octetos, *realm* sem distinção de maiúsculas | `[3GPP]` TS 23.003 §2.2B; `[INTERPRETAÇÃO]` o *realm* é nome de domínio |
| Colisão de A-TID entre SUPIs: o índice fica com o registro mais recente | `[3GPP]` §6.1 NOTA 2 deixa detecção e tratamento em aberto; `[MVP]` a escolha |
| Chaves apagadas são sobrescritas com zeros (melhor esforço) | `[MVP]` |

## 5. Autorização

**`register-anchorkey` e `remove-context` — OAuth2.** `[3GPP]` TS 29.535 §5.1.9: o acesso *"may be authorized"* por OAuth2, conforme configuração local. `[FREE5GC]` O NRF informa se o OAuth2 está ligado (`registro-nrf.md` §3). O token é verificado assim (TS 33.501 §13.4.1.1.2 passo 2):

| Verificação | Quem faz |
|---|---|
| esquema `Bearer` | AAnF |
| assinatura RS512 com o certificado do NRF | `oauth.VerifyOAuth` do free5GC |
| `scope` contém `naanf-akma` | `oauth.VerifyOAuth`; os escopos finos de TS 29.535 Tab. 5.1.9-1 não são exigidos (`[MVP]`; TS 33.501 §13.4.1.0 torna os adicionais opcionais) |
| `exp` no futuro | AAnF (o `VerifyOAuth` do free5GC não verifica) |
| `aud` igual ao `nfInstanceId` do AAnF ou a `AANF`; ausente ou vazio → recusado | AAnF. Exige o NRF corrigido (`nrf.md`, N3); até o MVP-1b, `aud` vazio era tolerado (DV-8, encerrado) |

Token ausente ou inválido → 401 com `WWW-Authenticate: Bearer` (ou `Bearer error="invalid_token"`, RFC 6750 §3) e `ProblemDetails`.

**`retrieve-applicationkey` — política local.** `[3GPP]` TS 33.535 §6.2.1 passo 2: o AAnF verifica *"based on the configured local policy or based on the authorization information available in the signalling (i.e., Oauth2.0 token)"*. `[MVP]` O AF não alcança o NRF, então vale a política local. Cada AF_ID autorizado fica associado aos endereços de origem dos quais pode ser pedido. Isso atende à recomendação contra a falsificação de AF (`docs/literatura/`), com o risco residual do DV-1. O endereço vem da conexão, nunca de `X-Forwarded-For`. O FQDN é comparado com exatidão, porque entra na KDF. AF fora da política, ou AF_ID pedido de outro endereço → 403 sem `cause`, por falta de causa definida (P3 do plano).

## 6. Erros: os 11 códigos obrigatórios para POST

`[3GPP]` TS 29.535 §5.1.7.1; TS 29.500 Tab. 5.2.7.1-1 e §5.2.7.2. Toda resposta de erro leva `ProblemDetails` em `application/problem+json`. Testes em `free5gc/aanf/internal/sbi/api_akma_test.go`, sobre h2c.

| Código | Quando | Causa | Teste |
|---|---|---|---|
| 400 | JSON inválido ou de tipo errado; atributo obrigatório ausente ou incorreto; versão da API diferente de `v1` | `INVALID_MSG_FORMAT`, `MANDATORY_IE_MISSING`, `MANDATORY_IE_INCORRECT`, `OPTIONAL_IE_INCORRECT`, `INVALID_API` | `TestRegisterBadRequests`, `TestRetrieveBadRequests`, `TestRoutingErrors` |
| 401 | token ausente ou inválido (seção 5) | — | `TestOAuth2` |
| 403 | K_AKMA ausente; AF não autorizado | `K_AKMA_NOT_PRESENT`; sem causa | `TestRetrieveKAkmaNotPresent`, `TestRetrieveAfNotAuthorized` |
| 404 | contexto a remover inexistente; operação desconhecida; caminho fora da API | `AKMA_CONTEXT_NOT_FOUND`; `RESOURCE_URI_STRUCTURE_NOT_FOUND`; sem causa | `TestRemoveContext`, `TestRoutingErrors` |
| 411 | `Content-Length` com valor inválido, ou diferente do corpo recebido | `INCORRECT_LENGTH` | `TestInvalidContentLengthValue`; `TestIncorrectContentLength` (quadros HTTP/2 montados à mão) |
| 413 | corpo acima de 16 000 000 octetos (TS 29.501 §6.2), declarado ou não | `MAX_JSON_SIZE_EXCEEDED` | `TestContentTooLarge` |
| 415 | `Content-Type` diferente de `application/json` (parâmetros como `charset` aceitos) | — | `TestUnsupportedMediaType` |
| 429 | **não emitido**: sem controle de sobrecarga no MVP. As causas `NF_CONGESTION_RISK` pressupõem esse controle (TS 29.500 §6.4) | — | justificado |
| 500 | falha interna (pânico num *handler*, erro da KDF) | `SYSTEM_FAILURE` | `TestInternalError` |
| 502 | **não emitido**: `INBOUND_SERVER_ERROR` exige que o AAnF chame outro servidor para atender ao pedido. No MVP ele não chama nenhum: não busca GPSI e não subscreve *roaming* (DV-2) | — | justificado; volta a valer com a subscrição ao UDM |
| 503 | **não emitido**: `NF_CONGESTION` também pressupõe controle de sobrecarga | — | justificado |

Também cobertos: **501** para métodos diferentes de `POST`, porque nenhum recurso da API os suporta (TS 29.500 §5.2.7.2); e barra final no caminho → 404, sem redirecionamento. `[INTERPRETAÇÃO]` Para 411, TS 29.500 só dá a causa `INCORRECT_LENGTH` (*"incorrect value of a Content-length header field"*). O AAnF a aplica também quando o valor não bate com o corpo. A ausência de `Content-Length` é aceita, porque o HTTP/2 delimita o corpo por quadros.

## 7. Logs

Nenhuma chave aparece em log, em nível algum. Desde o MVP-3a (decisão D16), SUPI, A-KID e AF_ID aparecem em `info`, na categoria `AKMA`; até o MVP-2 só apareciam em `debug`. O catálogo das mensagens está em `logs-akma.md`. Teste: `TestNoKeyInLogs`, em nível `trace`, procura K_AKMA e K_AF (hexadecimal e base64) no log. Na VM, em 17/10/2026: zero ocorrências de K_AKMA e K_AF em `aanf.log` e `aanf.out` depois do fluxo completo.

## 8. Como reproduzir numa instalação limpa

Pré-requisitos: free5GC v4.2.3 instalado e funcionando; Go 1.26.2; este repositório clonado. Os comandos supõem a topologia do projeto: repositório na máquina do AF, free5GC em `woca@192.168.16.130`. Numa máquina só, rode `build-aanf.sh` nela e copie os arquivos à mão.

1. Preparar o `openapi` e compilar (gera `free5gc/build/bin/aanf` e `seed`):

```bash
free5gc/scripts/build-aanf.sh
```

2. Ajustar `free5gc/config/aanfcfg.yaml`: `registerIPv4` e `bindingIPv4` com o endereço da máquina do free5GC alcançável pelo AF; `plmnSupportList`; `routingIndicators`; e, em `authorizedAfs`, o AF_ID e o endereço de cada AF.

3. Copiar para a máquina do free5GC (só arquivos novos; recusa se o AAnF estiver rodando):

```bash
free5gc/scripts/deploy.sh woca@192.168.16.130
```

4. Corrigir o NRF (`nrf.md` §4); sem isso, os tokens chegam sem `aud` e o AAnF os recusa. Com o NRF em execução (`run.sh` do free5GC, ou só o NRF), iniciar o AAnF:

```bash
ssh woca@192.168.16.130 bash -s -- start < free5gc/scripts/run-aanf.sh
```

O mesmo *script* aceita `stop`, `status` e `foreground`; na própria máquina do free5GC, rode-o diretamente. O AAnF roda a partir de `~/free5gc`, como as outras NFs, para que `cert/nrf.pem` resolva igual.

## 9. Como verificar

Registro no NRF: `registro-nrf.md` §7. Fluxo das três operações com a ferramenta `seed` e o vetor `testbed-ue` (`shared/akma/testdata/vectors.json`):

`[MVP]` O `seed` fala com o AAnF no endereço de `--aanf`, cujo padrão é o do laboratório (`http://192.168.16.130:8000`). Em outra máquina, passe `--aanf http://<sbi.registerIPv4>:8000` em `register`, `retrieve` e `remove`; sem isso, a resposta é `connection refused` ou, pior, vem de outro AAnF.

1. Na máquina do free5GC, no papel do AUSF, com token obtido do NRF:

```bash
~/akma-tools/seed register --nrf http://127.0.0.10:8000 --supi imsi-208930000000001 --k-ausf 141ef548f86e3a2a1ad949331caa9e6bb12401afa321bde4f9f23b84870dd654 --rid 0000 --mcc 208 --mnc 93
```

A saída deve mostrar `a-kid: rid0000.atidiz1FN8DkMGyxwtfTVYsKUsWnzOJdBOLlQEiQlK2gHNA=@5gc.mnc093.mcc208.3gppnetwork.org`.

2. Na máquina do AF, no papel do AF:

```bash
free5gc/build/bin/seed retrieve --af-id af1.akma.example.FF00000001 --akid 'rid0000.atidiz1FN8DkMGyxwtfTVYsKUsWnzOJdBOLlQEiQlK2gHNA=@5gc.mnc093.mcc208.3gppnetwork.org'
```

A saída deve mostrar `kaf: 00c16845c2e4e1256755b84bb06d94e2e659d8ba53c59e874fc149bffd5ea04b`, que é a K_AF do vetor, com `expiry` uma hora à frente, em UTC, e o `supi`.

3. Na máquina do free5GC: `~/akma-tools/seed remove --nrf http://127.0.0.10:8000 --supi imsi-208930000000001` responde `removed`. Repetido, responde `404 (AKMA_CONTEXT_NOT_FOUND)`. O passo 2 passa a responder `403 (K_AKMA_NOT_PRESENT)`.

Resultado em 17/10/2026, todos conforme o esperado:

| Caso | Resposta |
|---|---|
| registro com token, na máquina do free5GC | 200; A-KID igual ao do vetor |
| registro sem token, da máquina do AF | 401 |
| K_AF pedida da máquina do AF | 200; K_AF igual à do vetor; `expiry` em UTC |
| pedido anônimo (`--anon`) | 200, sem `supi` |
| AF_ID fora da política | 403 sem causa |
| AF_ID autorizado pedido da máquina do free5GC (outro endereço) | 403 sem causa |
| A-KID desconhecido | 403 `K_AKMA_NOT_PRESENT` |
| remoção, remoção repetida, novo pedido | 204, 404 `AKMA_CONTEXT_NOT_FOUND`, 403 `K_AKMA_NOT_PRESENT` |

Os erros de protocolo (411, 413, 415, 501…) não foram repetidos pela rede: a máquina do AF não tem `curl`. Eles estão cobertos pelos testes h2c da seção 6, que exercitam o mesmo servidor.

## 10. Limitações

- Sem TLS (DV-1); sem subscrição de *roaming* (DV-2).
- Contexto em memória, perdido ao reiniciar.
- Exige o NRF corrigido (`nrf.md`): sem o *patch* 0003, os tokens chegam sem `aud` e `register-anchorkey` e `remove-context` respondem 401. Sem os *patches* 0001 e 0002, o perfil fica sem `aanfInfoList` e o AAnF não pode ser descoberto.
- O AAnF não entra no `run.sh` do free5GC: alterar esse arquivo fica para quando o AAnF for proposto *upstream*.
