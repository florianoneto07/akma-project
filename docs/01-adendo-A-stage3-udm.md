# Adendo A — Mapeamento Stage 3 das interfaces AUSF↔UDM e AAnF↔UDM

**Documento:** complemento ao *Baseline Normativo* v2.0
**Versão:** A.1
**Data:** 02 de setembro de 2026
**Errata:** 16 de outubro de 2026 — §A.5 e §A.7, marcadas como *(errata 16/10/2026)*
**Fonte incorporada:** arquivos OpenAPI do Anexo A de 3GPP TS 29.503 V19.7.0

---

## A.1 Arquivos incorporados e suas versões

| Arquivo | `info.version` | TS 29.503 declarada | Relevância para o AKMA |
|---|---|---|---|
| `TS29503_Nudm_UEAU.yaml` | 1.4.0 | 19.5.0 | **Alta** — `AuthenticationInfoResult` |
| `TS29503_Nudm_UECM.yaml` | 1.4.2 | 19.7.0 | **Alta** — `trigger-auth` |
| `TS29503_Nudm_SDM.yaml` | 2.4.2 | 19.7.0 | Média — `id-translation-result` |
| `TS29503_Nudm_EE.yaml` | 1.4.1 | 19.7.0 | Média — `ROAMING_STATUS` |
| `TS29503_Nudm_PP.yaml` | 1.4.1 | 19.6.0 | Não aplicável |
| `TS29503_Nudm_MT.yaml` | 1.3.1 | 19.7.0 | Não aplicável |
| `TS29503_Nudm_NIDDAU.yaml` | 1.3.0 | 19.5.0 | Não aplicável |
| `TS29503_Nudm_RSDS.yaml` | 1.3.0 | 19.5.0 | Não aplicável |
| `TS29503_Nudm_SSAU.yaml` | 1.2.0 | 19.5.0 | Não aplicável |
| `TS29503_Nudm_UEID.yaml` | 1.1.0 | 19.5.0 | Não aplicável |

`[3GPP]` A heterogeneidade das versões declaradas em `externalDocs` não indica inconsistência do acervo. Cada arquivo OpenAPI declara a versão da TS na qual foi alterado pela última vez; arquivos não modificados desde a V19.5.0 mantêm essa declaração. Todos pertencem ao Anexo A da TS 29.503 V19.7.0, que é o corpo normativo do acervo.

---

## A.2 Interface AUSF↔UDM: `AuthenticationInfoResult`

`[3GPP]` `TS29503_Nudm_UEAU.yaml`, esquema `AuthenticationInfoResult`:

```yaml
akmaInd:
  type: boolean
  default: false
routingId:
  type: string
  pattern: '^[0-9]{1,4}$'
```

Apenas `authType` consta em `required`. Isto confirma o registro da v2.0: `akmaInd` e `routingId` são opcionais no esquema, e a obrigatoriedade de emiti-los para assinantes com assinatura AKMA decorre de TS 33.501 V19.7.0, não do esquema JSON.

`[3GPP]` **`routingId` é do tipo `string`, com padrão de 1 a 4 dígitos decimais.** Isto é consistente com TS 23.003 V19.7.0 cláusula 2.2B, que estabelece que cada dígito do Routing Indicator é significativo e que `"012"` não equivale a `"12"`.

`[FREE5GC]` **Advertência de implementação.** O Routing Indicator deve ser armazenado e propagado como cadeia de caracteres em toda a cadeia UDM → AUSF → A-KID → AAnF. Convertê-lo para inteiro em qualquer ponto — em uma estrutura Go, em um documento MongoDB ou em uma configuração YAML — descarta zeros à esquerda e produz um A-KID distinto do gerado pelo UE. O sintoma seria uma falha de `K_AKMA_NOT_PRESENT` no AAnF apenas para assinantes cujo RID começa com zero, o que é difícil de diagnosticar por não afetar todos os assinantes.

### A.2.1 Corroboração da advertência sobre o prefixo `imsi-`

`[3GPP]` `TS29503_Nudm_UECM.yaml`, variável de URI `ueId`, e TS 29.503 Tabela 6.2.3.13.2-1:

```
pattern: "(imsi-[0-9]{5,15}|nai-.+|.+)"
```

Isto confirma, no Stage 3, que o SUPI trafega nas interfaces SBI **com** o prefixo `imsi-` ou `nai-`. O parâmetro P0/P1 da KDF, por TS 33.501 Anexo A.7.0, é o IMSI conforme TS 23.003 cláusula 2.2 — **sem** prefixo. A conversão entre as duas representações é obrigatória e ocorre exatamente na fronteira entre a camada SBI e a camada criptográfica, tanto no AUSF quanto no AAnF.

Esta é a corroboração independente da primeira advertência de alto impacto registrada na seção 2.3 do documento v2.0.

---

## A.3 `Nudm_UECM_AuthTrigger`: nova divergência

### A.3.1 Método HTTP

`[DIVERGÊNCIA] — DIVERGÊNCIA D.` TS 29.503 V19.7.0, Tabela 6.2.3.13.4.1-1, declara para a operação customizada `trigger-auth` o **mapeamento para o método POST**. A tabela 6.2.3.13.4.2.2-1 é intitulada *"Data structures supported by the POST Request Body on this resource"*.

O arquivo `TS29503_Nudm_UECM.yaml` v1.4.2 declara a mesma operação como:

```yaml
/{ueId}/registrations/trigger-auth:
  get:
    summary: trigger the primary (re-)authentication
    operationId: authTrigger
    requestBody:
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/AuthTriggerInfo'
      required: true
```

Duas observações:

1. A declaração é `get`, e não `post`.
2. A operação declara `requestBody` com `required: true`. Um `GET` com corpo de requisição obrigatório contraria a semântica de HTTP e não é suportado de forma previsível por servidores e clientes HTTP/2 correntes.

`[DIVERGÊNCIA]` **A regra formal de precedência agrava o problema.** TS 29.503 V19.7.0, Anexo A.1, contém a mesma cláusula presente em TS 29.535 e TS 29.522:

> *"This Annex takes precedence when being discrepant to other parts of the specification with respect to the encoding of information elements **and methods** within the API(s)."*

Aplicada literalmente, a cláusula faria prevalecer o `get` do arquivo OpenAPI. O resultado seria uma operação não implementável de forma interoperável e em contradição com TS 29.501, segundo a qual operações customizadas sem recurso associado são mapeadas para POST.

`[INTERPRETAÇÃO]` Adotar **POST**, seguindo o corpo da especificação e TS 29.501, e registrar a divergência. Fundamento: a cláusula de precedência do Anexo A destina-se a resolver discrepâncias de codificação, não a validar construções que violem a especificação de protocolo subjacente. Trata-se, aparentemente, de defeito editorial do arquivo OpenAPI, análogo ao já registrado na DIVERGÊNCIA A (rótulo `atid`/`akid` em TS 29.522).

Como a operação está fora do escopo do MVP — TS 33.535 §6.4.4 emprega linguagem permissiva —, a decisão não bloqueia o cronograma. Recomenda-se, contudo, documentar as duas divergências editoriais (A e D) em um registro único do repositório, por serem candidatas a submissão ao 3GPP CT4.

### A.3.2 Obrigatoriedade do SUPI

`[DIVERGÊNCIA]` O esquema `AuthTriggerInfo` não declara `required`:

```yaml
AuthTriggerInfo:
  description: This data type contains the SUPI for authorization information.
  type: object
  properties:
    supi:
      $ref: 'TS29571_CommonData.yaml#/components/schemas/Supi'
```

Já TS 33.501 V19.7.0 cláusula 14.2.6 declara *"Input, Required: SUPI"*, e a Tabela 6.2.3.13.4.2.2-1 de TS 29.503 classifica `AuthTriggerInfo` como `M`, com a descrição *"Contains the SUPI"*. A implementação deve tratar `supi` como obrigatório.

### A.3.3 Erro de aplicação adicional

`[3GPP]` TS 29.503 V19.7.0, Tabela 6.2.3.13.4.2.2-2, define para `trigger-auth` o erro de aplicação **`REAUTHENTICATION_NOT_ALLOWED`** com estado HTTP **403 Forbidden**, além de `404 Not Found` para usuário inexistente. Este erro não constava do documento v2.0 e deve ser acrescentado à lista de erros a tratar quando a operação for implementada.

---

## A.4 Obtenção do GPSI pelo AAnF

`[3GPP]` `TS29503_Nudm_SDM.yaml` v2.4.2:

```
GET {apiRoot}/nudm-sdm/v2/{ueId}/id-translation-result
```

- `ueId` (path, obrigatório): tipo `VarUeId`
- Parâmetros de consulta opcionais: `supported-features`, `af-id`, `app-port-id`, `mtc-provider-info`, `requested-gpsi-type`
- Cabeçalhos condicionais: `If-None-Match`, `If-Modified-Since`
- Escopo OAuth2 específico: **`nudm-sdm:id-translation-result:read`**, além de `nudm-sdm`

Resposta, esquema `IdTranslationResult`:

| Atributo | Obrigatório | Observação |
|---|---|---|
| `supi` | **sim** | Único atributo em `required` |
| `gpsi` | não | |
| `additionalSupis` | não | array, `minItems: 1` |
| `additionalGpsis` | não | array, `minItems: 1` |
| `supportedFeatures` | não | |

`[3GPP]` **Consequência para o AAnF:** `gpsi` é opcional na resposta. O AAnF deve tratar a ausência do GPSI como resultado normal e não como erro. Isto se articula com TS 33.535 §6.2.1 passo 8, segundo o qual a escolha entre enviar SUPI ou GPSI ao AF é do AAnF, por política local: se o GPSI não estiver disponível no UDM, a única resposta possível é com SUPI.

`[3GPP]` O parâmetro de consulta `af-id` permite que o AAnF informe ao UDM a identidade do AF para o qual a tradução é solicitada, o que é coerente com a redação de TS 33.535 §6.2.1 passo 3 — *"se o AAnF determina que este AF específico necessita de GPSI"*.

---

## A.5 Estado de *roaming*: `Nudm_EventExposure`

`[3GPP]` `TS29503_Nudm_EE.yaml` v1.4.1:

```
POST {apiRoot}/nudm-ee/v1/{ueIdentity}/ee-subscriptions
```

- `ueIdentity` admite, entre outros, `imsi-`, `nai-`, `msisdn-`, `extid-`, `gci-`, `gli-`, `extgroupid-` e `anyUE`
- Corpo: `EeSubscription`, com `callbackReference` (`Uri`) e `monitoringConfigurations` (mapa com ao menos uma entrada) **obrigatórios**
- Sucesso: **201 Created**, com cabeçalho `Location` obrigatório apontando para `{apiRoot}/nudm-ee/v1/{ueIdentity}/ee-subscriptions/{subscriptionId}`, e corpo `CreatedEeSubscription`
- Escopo OAuth2 específico: `nudm-ee:subscription:create`

Evento e relatório:

- Valor de enumeração do tipo de evento: **`ROAMING_STATUS`**
- Esquema `RoamingStatusReport`:

```yaml
RoamingStatusReport:
  required: [roaming, newServingPlmn]
  properties:
    roaming:        { type: boolean }
    newServingPlmn: { $ref: '.../PlmnId' }
    accessType:     { $ref: '.../AccessType' }
    purged:         { type: boolean, enum: [true] }
```

`[3GPP]` **Articulação com TS 33.535 §6.2.1 passo 7.** A cláusula estabelece que, em registro duplo, o UE é tratado como em *roaming* se ao menos uma das PLMNs servidoras assim indicar. O esquema `RoamingStatusReport` fornece exatamente os elementos necessários a essa avaliação: o indicador booleano `roaming` e o `accessType` que discrimina o acesso ao qual o relatório se refere. O AAnF deve, portanto, manter o estado de *roaming* **por tipo de acesso** e aplicar disjunção lógica, e não armazenar um único booleano por UE.

`[3GPP]` A cláusula 5.5.2.2.2 de TS 29.503 *(errata 16/10/2026: antes citada como 5.2.3.5.2)* registra que, com `immediateFlag` definido como `true`, o UDM *"may include the current status of the event if available in the response"*. Isto corresponde diretamente ao passo 6 de TS 33.535 §6.2.1, no qual o AAnF recebe a informação de *roaming* na própria resposta à subscrição, sem aguardar uma notificação posterior. **Sem `immediateFlag`, o passo 6 não é satisfeito** e o AAnF ficaria bloqueado à espera de uma mudança de estado que pode nunca ocorrer. *(errata 16/10/2026)* Mesmo com `immediateFlag`, o relatório imediato é permissivo (*"may"*, *"if available"*); TS 23.502 V19.9.0 §5.2.3.5.2 diz o mesmo (*"First corresponding event report is included, if corresponding information is available"*). O AAnF precisa tratar a resposta sem relatório.

Esta constatação é relevante para a INCERTEZA 1: ela demonstra que os passos 5 e 6 da cláusula 6.2.1 podem ser implementados de forma síncrona quando o UDM tem o estado disponível *(errata 16/10/2026: a síncronia não é garantida)*, o que torna a leitura literal — subscrição incondicional — tecnicamente viável, embora ainda não resolva a questão da condicionalidade em relação à *feature* `RoamingRestriction`.

---

## A.6 A LACUNA C à luz do Anexo B de TS 33.535

Examinando os perfis Ua* de TS 33.535 V19.0.0 com a hipótese de que a codificação do A-TID pudesse estar fixada no contexto de uso do A-KID:

`[3GPP]` Anexo B.1.2.2 (autenticação do UE por chave compartilhada com autenticação do AF por certificado): no passo 5, o cliente envia um cabeçalho `Authorization` com Digest, **usando o A-KID como *username*** e K_AF como senha no cálculo do Digest.

`[3GPP]` Anexo B.1.3.2.1 (TLS 1.2): o UE envia uma mensagem `ClientKeyExchange` com uma **identidade PSK composta pela cadeia `"3GPP-AKMA"` e o A-KID**. O AF utiliza o prefixo `"3GPP-AKMA"` para reconhecer a modalidade e o A-KID para obter K_AF junto ao AAnF.

`[3GPP]` Anexo B.1.3.2.2 (TLS 1.3): as identidades PSK no `ClientHello` incluem o prefixo de espaço de nomes `"3GPP-AKMA"` e o A-KID.

`[3GPP]` Anexo C.1.1 e C.2.2 (DTLS 1.3): a identidade PSK especificada em B.1 aplica-se também ao DTLS.

**Resultado: a LACUNA C não é encerrada.** Nenhum dos anexos declara a codificação textual do valor do A-TID.

`[INTERPRETAÇÃO]` Os anexos, contudo, impõem uma **restrição que estreita o espaço de soluções**: o A-KID é utilizado como *username* de HTTP Digest e como identidade PSK de TLS/DTLS. Em ambos os contextos ele deve ser uma cadeia de caracteres imprimível, o que exclui a representação binária direta do valor de 256 bits e mantém como candidatas apenas codificações textuais — hexadecimal ou base64url. Isto reforça, sem confirmar, a decisão de projeto registrada na v2.0 de adotar **hexadecimal minúsculo de 64 caracteres** como padrão parametrizável.

Observa-se ainda que o Anexo B.1.3.2.2 remete a TS 33.222 cláusula 5.4.0.2, que é exatamente o protocolo associado ao identificador Ua* `(0x01,0x00,0x00,0x00,0x07)` do Anexo H de TS 33.220. A cadeia de referências é, nesse ponto, internamente consistente.

`[LACUNA]` Fontes ainda não consultadas que podem encerrar a LACUNA C: **TS 33.222** (perfis TLS de que os anexos de TS 33.535 são adaptação) e **TS 24.109** (protocolos Ua sobre HTTP, incluindo a forma do B-TID no GBA, cuja construção é análoga à do A-KID).

---

## A.7 Atualização da matriz de rastreabilidade

Acrescentar ou revisar as seguintes linhas:

| Requisito | Norma | Cláusula | Implementação prevista | Teste | Estado |
|---|---|---|---|---|---|
| `routingId` como cadeia com zeros à esquerda preservados | TS 29.503 §6.3.6.2.3; `TS29503_Nudm_UEAU.yaml` | — | `patches/udm/`, `akma/akid/rid.go` | `rid_test.go` com RID `"0"`, `"012"`, `"0000"` | **Especificado** |
| Conversão SUPI SBI ↔ SUPI KDF | TS 33.501 A.7.0; TS 29.503 Tab. 6.2.3.13.2-1 | — | `akma/kdf/supi.go` | `supi_test.go` | **Especificado** |
| GPSI via `id-translation-result` | TS 29.503; `TS29503_Nudm_SDM.yaml` | 5.2.2.2.10 | `NFs/aanf/internal/sbi/consumer/` | integração | Fora do MVP; **tratamento de `gpsi` ausente é obrigatório** |
| Subscrição `ROAMING_STATUS` com `immediateFlag` | TS 29.503 §5.5.2.2.2 *(errata)*; `TS29503_Nudm_EE.yaml` | 6.4.6.2.8 | — | — | **INCERTEZA 1** |
| Estado de *roaming* por tipo de acesso | TS 33.535 §6.2.1 passo 7; `RoamingStatusReport` | — | `NFs/aanf/internal/context/` | unitário | **Especificado** |
| `trigger-auth` — método HTTP | TS 29.503 Tab. 6.2.3.13.4.1-1 × `TS29503_Nudm_UECM.yaml` | — | — | — | **DIVERGÊNCIA D** — adotar POST |
| Erro `REAUTHENTICATION_NOT_ALLOWED` | TS 29.503 Tab. 6.2.3.13.4.2.2-2 | — | — | — | Fora do MVP |

---

## A.8 Estado consolidado das questões em aberto

| Nº | Assunto | Estado |
|---|---|---|
| 1 | Obrigatoriedade da subscrição de `ROAMING_STATUS` | **Aberta.** Demonstrada a viabilidade técnica da leitura literal via `immediateFlag` |
| 5 | `AkmaInd`/`RoutingId` no `free5gc/openapi` | **Aberta.** Requer inspeção direta do repositório |
| A | Rótulo `atid` versus `akid` em TS 29.522 | **Aberta.** Defeito editorial aparente |
| B | Codificação textual de `kAkma` e `kaf` | **Aberta.** Verificar em TS 29.571 |
| C | Codificação do valor do A-TID | **Aberta, com espaço de soluções estreitado.** Consultar TS 33.222 e TS 24.109 |
| **D** | Método HTTP de `trigger-auth` | **Nova, aberta.** Adotar POST; fora do MVP |

Especificações ainda ausentes do acervo: **TS 29.571** (pode encerrar a LACUNA B), **TS 29.510** (único item ainda bloqueado na matriz — `NFType` do AAnF), **TS 29.500**, **TS 29.501**, **TS 29.122**, **TS 33.222** e **TS 24.109** (podem encerrar a LACUNA C).

---

## A.9 Síntese

A incorporação dos arquivos OpenAPI de TS 29.503 confirma integralmente o mapeamento Stage 3 registrado no documento v2.0 para as três interfaces do AAnF e do AUSF com o UDM, e acrescenta quatro elementos:

1. `routingId` é uma cadeia de 1 a 4 dígitos decimais com zeros à esquerda significativos, o que impõe uma restrição de tipo em toda a cadeia de propagação.
2. O padrão de `ueId` em TS 29.503 corrobora, de forma independente, que o SUPI trafega com prefixo nas interfaces SBI, reforçando a advertência sobre a conversão exigida antes da KDF.
3. A operação `trigger-auth` apresenta divergência de método HTTP entre o corpo da especificação e o arquivo OpenAPI, agravada por uma cláusula de precedência que, aplicada literalmente, produziria uma construção não implementável. A questão está fora do escopo do MVP.
4. O esquema `RoamingStatusReport` fornece `accessType`, o que torna implementável a regra de registro duplo de TS 33.535 §6.2.1 passo 7 e impõe que o estado de *roaming* seja mantido por tipo de acesso no contexto do AAnF.

Nenhum dos elementos acima altera a especificação das derivações criptográficas, que permanece completa desde a v2.0, nem o caminho crítico do MVP.
