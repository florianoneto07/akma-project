# Adendo B — Camada SBA, registro no NRF e tipos comuns

**Documento:** complemento ao *Baseline Normativo* v2.0 e ao *Adendo A*
**Versão:** B.1
**Data:** 02 de setembro de 2026
**Fontes incorporadas:** 3GPP TS 29.571 V19.7.0, TS 29.510 V19.7.0, TS 29.500 V19.7.0, TS 29.501 V19.5.0, e os respectivos arquivos OpenAPI do Anexo A

---

## B.1 Encerramento do último item bloqueado da matriz

`[3GPP]` TS 29.510 V19.7.0 e `TS29510_Nnrf_NFManagement.yaml`:

- O tipo de NF do AAnF é **`"AANF"`**, presente na enumeração `NFType`.
- O nome de serviço **`naanf-akma`** consta da enumeração `ServiceName`.
- Existe o tipo `AanfInfo`, descrito como *"Represents the information relative to an AAnF NF Instance"*:

```yaml
AanfInfo:
  type: object
  properties:
    routingIndicators:
      type: array
      items:
        type: string
        pattern: '^[0-9]{1,4}$'
      minItems: 1
```

- O `NFProfile` possui o atributo `aanfInfoList`, um mapa cujas chaves são cadeias JSON únicas e cujos valores são objetos `AanfInfo`.

Isto encerra a última linha bloqueada da matriz de rastreabilidade do documento v2.0. O registro do AAnF no NRF está integralmente especificado.

### B.1.1 Observações sobre `AanfInfo`

`[3GPP]` `routingIndicators` **não consta de `required`**. Um AAnF pode, portanto, registrar-se sem declarar Routing Indicators. Isto é coerente com TS 33.535 §6.7, segundo a qual, quando o Routing Indicator do UE está no valor padrão, o consumidor pode selecionar qualquer instância de AAnF da rede doméstica.

`[3GPP]` `routingIndicators` é um **array**: uma instância de AAnF pode atender a mais de um Routing Indicator. E `aanfInfoList` é um **mapa por instância de serviço**, não um objeto único por NF. A modelagem interna do AAnF deve refletir essa cardinalidade caso venha a expor múltiplas instâncias de serviço.

`[3GPP]` O padrão `^[0-9]{1,4}$` é o mesmo de `routingId` em `AuthenticationInfoResult` (TS 29.503) e coerente com TS 23.003 §2.2B. A restrição de tipo cadeia, com zeros à esquerda significativos, registrada no Adendo A §A.2, aplica-se também aqui.

### B.1.2 Descoberta

`[3GPP]` `TS29510_Nnrf_NFDiscovery.yaml` define o parâmetro de consulta:

```yaml
- name: routing-indicator
  in: query
  description: routing indicator in SUCI
  schema:
    type: string
    pattern: '^[0-9]{1,4}$'
```

`[3GPP]` TS 29.510 V19.7.0 descreve o parâmetro como *"Routing Indicator information that allows to route network signalling with SUCI (see TS 23.003) to an AUSF, AAnF and UDM instance capable to serve the subscriber. May be included if the target NF type is 'AUSF', 'AANF' or 'UDM'."*

Com isso, a cláusula 6.7 de TS 33.535 fica integralmente mapeada no Stage 3:

| Requisito de TS 33.535 §6.7 | Realização Stage 3 |
|---|---|
| O AAnF é descoberto via NRF | `GET .../nnrf-disc/v1/nf-instances?target-nf-type=AANF&requester-nf-type=AUSF` |
| A seleção considera o Routing Indicator do UE | Parâmetro de consulta `routing-indicator` |
| AUSF, AF interno e NEF devem selecionar o mesmo conjunto de AAnFs | `AanfInfo.routingIndicators` publicado no perfil do AAnF |
| Com RID no valor padrão, qualquer instância pode ser selecionada | `routingIndicators` não é obrigatório em `AanfInfo` |
| Configuração local é alternativa admitida | Independente do NRF |

### B.1.3 Campos obrigatórios do registro

`[3GPP]` `NFProfile`: `nfInstanceId`, `nfType`, `nfStatus`, e ao menos um dentre `fqdn`, `ipv4Addresses` e `ipv6Addresses` (construção `anyOf`).

`[3GPP]` `NFService`: `serviceInstanceId`, `serviceName`, `versions`, `scheme`, `nfServiceStatus`.

`[FREE5GC]` O AUSF do free5GC já implementa registro no NRF com essa estrutura. O esqueleto do AAnF, derivado do AUSF conforme proposto na seção 6.3 do documento v2.0, herda o mecanismo; as alterações necessárias limitam-se a `nfType = "AANF"`, `serviceName = "naanf-akma"` e ao preenchimento de `aanfInfoList`.

`[HIPÓTESE]` A presença das constantes `AANF` e `naanf-akma` no repositório `free5gc/openapi` não foi verificada. A verificação deve preceder o início da fase F3.

---

## B.2 Códigos de estado HTTP obrigatórios

`[3GPP]` TS 29.500 V19.7.0, Tabela 5.2.7.1-1, para o método **POST** — único método usado pelas três operações do Naanf_AKMA:

| Classificação | Códigos |
|---|---|
| **Obrigatórios (M)** | 400, 401, 403, 404, 411, 413, 415, 429, 500, 502, 503 |
| Específicos do serviço (SS) | 200, 201, 202, 204, 303, 307, 308, 405, 408, 409, 410, 412, 501, 504 |
| Não aplicáveis (N/A) | 100, 300, 414, 422 |

`[3GPP]` Para os códigos 4xx e 5xx marcados como obrigatórios, a NF deve suportar o tratamento do objeto JSON `ProblemDetails` com `Content-Type` igual a `application/problem+json`, salvo se a definição da API especificar outro corpo de resposta.

**Validação cruzada.** O arquivo `TS29535_Naanf_AKMA.yaml` v1.2.0 declara, para as três operações, exatamente os códigos `400, 401, 403, 404, 411, 413, 415, 429, 500, 502, 503`, acrescidos de `307`, `308`, o código de sucesso específico e `default`. A lista do arquivo OpenAPI corresponde, portanto, precisamente ao conjunto obrigatório da Tabela 5.2.7.1-1 para POST. Isto confirma que a implementação do AAnF não precisa tratar códigos além dos já enumerados no arquivo OpenAPI, e fornece um critério objetivo de completude para os testes negativos da fase F7.

---

## B.3 Cabeçalhos customizados

`[3GPP]` TS 29.535 V19.5.0, §5.1.2.3, exige suporte aos cabeçalhos obrigatórios de TS 29.500 §5.2.3.2 e admite suporte aos opcionais de §5.2.3.3, sem definir cabeçalhos próprios.

`[3GPP]` TS 29.500 V19.7.0, Tabela 5.2.3.2.1-1 — obrigatórios:

| Cabeçalho | Aplicação ao AKMA |
|---|---|
| `3gpp-Sbi-Message-Priority` | Incluído quando é necessário transmitir prioridade da mensagem |
| `3gpp-Sbi-Callback` | **Relevante.** Deve ser incluído em POST de *callback* para consumidores em outra PLMN via SEPP e em comunicação indireta; **deveria** ser incluído em qualquer requisição de notificação de evento em comunicação direta |
| `3gpp-Sbi-Target-apiRoot` | Comunicação indireta via SCP |
| `3gpp-Sbi-Scp-apiRoot` | Comunicação indireta entre PLMNs |

`[3GPP]` **Consequência concreta:** a operação `Naanf_AKMA_Notify` (TS 29.535 §4.2.2.5) é uma notificação enviada pelo AAnF ao consumidor via POST no `notifUri`. Segundo TS 29.500, o cabeçalho `3gpp-Sbi-Callback` **deveria** ser incluído nessa requisição mesmo em comunicação direta. Este requisito não constava dos documentos anteriores e deve ser acrescentado à especificação da notificação, quando esta for implementada.

`[MVP]` Os quatro cabeçalhos dizem respeito a prioridade, *callbacks* e comunicação indireta. Nenhum é exercitado pelo fluxo básico do MVP em comunicação direta sem SCP. O suporte deve, ainda assim, constar do plano de conformidade da fase F7.

---

## B.4 Estrutura de URI e a DIVERGÊNCIA D

`[3GPP]` TS 29.501 V19.5.0, §4.4.2 — URI de operação customizada **não associada a recurso**:

```
{apiRoot}/<apiName>/<apiVersion>/<custOpName>
```

Esta é exatamente a forma das três operações do Naanf_AKMA, coerente com TS 29.535 §5.1.3, que declara não haver recursos definidos para a API nesta versão.

`[3GPP]` TS 29.501 V19.5.0, §4.6.1.2: as duas figuras normativas que ilustram operações customizadas — Figura 4.6.1.2-1, para operação sobre recurso, e Figura 4.6.1.2-2, para operação relacionada ao serviço — são intituladas *"Custom Operation … using HTTP POST"*. O texto refere-se ao corpo da resposta como *"POST response body"*.

`[3GPP]` TS 29.501 V19.5.0, Anexo C.4: *"The custom operation archetype can be used to model an **unsafe and non-idempotent** operation that is not a Create on a collection."*

**Reforço da DIVERGÊNCIA D.** O método GET é, por definição de IETF RFC 9110, seguro e idempotente. A declaração de `trigger-auth` como `get:` no arquivo `TS29503_Nudm_UECM.yaml` contradiz, portanto, não apenas a Tabela 6.2.3.13.4.1-1 do corpo da TS 29.503, mas também o arquétipo de operação customizada definido em TS 29.501 Anexo C.4 e o padrão de método das figuras de §4.6.1.2.

A interpretação registrada no Adendo A — adotar POST — fica sustentada por três fontes independentes. A divergência permanece como defeito editorial do arquivo OpenAPI, candidato a submissão ao 3GPP CT4 juntamente com a DIVERGÊNCIA A.

---

## B.5 Tipos comuns: confirmações e uma lacuna que persiste

### B.5.1 Tipos confirmados

`[3GPP]` `TS29571_CommonData.yaml`:

| Tipo | Definição | Observação |
|---|---|---|
| `Supi` | `string`, `pattern: '^(imsi-[0-9]{5,15}\|nai-.+\|gci-.+\|gli-.+\|.+)$'` | Formatado como `"imsi-<imsi>"`, com `<imsi>` conforme TS 23.003 §2.2 |
| `Gpsi` | `string`, `pattern: '^(msisdn-[0-9]{5,15}\|extid-[^@]+@[^@]+\|.+)$'` | |
| `VarUeId` | `string`, admite `imsi-`, `nai-`, `msisdn-`, `extid-`, `gci-`, `gli-` | Usado em `trigger-auth` e `id-translation-result` |
| `SupportedFeatures` | `string`, `pattern: '^[A-Fa-f0-9]*$'` | Máscara de bits em hexadecimal; cada caractere representa quatro *features* |
| `DateTime` | `string`, `format: date-time` | |

`[3GPP]` **Terceira corroboração independente do prefixo.** A descrição do tipo `Supi` em TS 29.571 declara que a cadeia é formatada como `"imsi-<imsi>"`, sendo `<imsi>` o valor conforme TS 23.003 §2.2. Confirma-se, agora pela definição do próprio tipo comum, que a representação SBI do SUPI carrega o prefixo, enquanto o parâmetro da KDF, por TS 33.501 Anexo A.7.0, é o `<imsi>` isolado. A advertência da seção 2.3 do documento v2.0 está triplamente fundamentada — por TS 33.501, por TS 29.503 e por TS 29.571.

### B.5.2 A LACUNA B não é encerrada

`[LACUNA]` `TS29571_CommonData.yaml` **não define nenhum tipo comum para material de chave criptográfica**. Não há tipo `Key`, `KeyMaterial` ou equivalente, nem padrão hexadecimal genérico aplicável. Os atributos `kAkma` (TS 29.535) e `kaf` (TS 29.522) permanecem declarados apenas como `type: string`, sem `pattern` nem `format`.

`[INTERPRETAÇÃO]` A decisão de projeto registrada na v2.0 — **hexadecimal minúsculo de 64 caracteres** — mantém-se, agora com um argumento adicional: a convenção de codificação hexadecimal para valores binários é sistemática em TS 29.571, que a emprega em `SupportedFeatures`, `Nid`, `PlmnId` estendido e diversos identificadores de rede. A escolha é, portanto, consistente com o estilo da série 29, embora continue sendo decisão do projeto e não requisito 3GPP.

### B.5.3 Um ponto a verificar quando TS 29.122 chegar

`[3GPP]` Observa-se, no arquivo `TS29522_AKMA.yaml`, que os atributos `notifUri` de `AkmaAfKeyRequest` e `expiry` de `AkmaAfKeyData` referenciam **`TS29122_CommonData.yaml`**, e não `TS29571_CommonData.yaml`:

```yaml
notifUri:
  $ref: 'TS29122_CommonData.yaml#/components/schemas/Uri'
expiry:
  $ref: 'TS29122_CommonData.yaml#/components/schemas/DateTime'
```

Isto decorre de TS 29.522 pertencer ao domínio das APIs de exposição *northbound*, que reutilizam os tipos comuns de TS 29.122. **A definição de `DateTime` aplicável a `expiry` é, portanto, a de TS 29.122, e não a de TS 29.571.** As duas são presumivelmente equivalentes, mas a equivalência deve ser verificada quando TS 29.122 estiver disponível, pois `expiry` determina o formato de serialização do tempo de vida de K_AF.

`[3GPP]` **Implicação operacional do tipo `DateTime`.** Sendo `expiry` um instante absoluto, e não uma duração, a validade de K_AF depende da concordância de relógio entre AAnF, AF e, indiretamente, UE. Um desvio de relógio no AF produz aceitação de chaves expiradas ou rejeição de chaves válidas. O ambiente experimental deve, portanto, incluir sincronização de tempo entre os contêineres, e o plano de testes negativos da fase F7 deve prever um caso de desvio de relógio deliberado.

---

## B.6 Estado consolidado

### B.6.1 Matriz de rastreabilidade — linhas revisadas

| Requisito | Norma | Cláusula | Implementação prevista | Teste | Estado |
|---|---|---|---|---|---|
| Registro do AAnF no NRF | TS 29.510 V19.7.0; `TS29510_Nnrf_NFManagement.yaml` | — | `NFs/aanf/internal/sbi/consumer/` | integração | **Especificado** (antes: bloqueado) |
| `nfType = "AANF"`, `serviceName = "naanf-akma"` | TS 29.510 V19.7.0 | — | `NFs/aanf/pkg/factory/`, `patches/openapi/` | integração | **Especificado** |
| `AanfInfo.routingIndicators` no perfil | TS 29.510 V19.7.0 | — | `NFs/aanf/internal/sbi/consumer/` | integração | **Especificado** |
| Descoberta por `routing-indicator` | TS 29.510 V19.7.0; TS 33.535 §6.7 | — | `patches/ausf/internal/sbi/consumer/` | integração | **Especificado** |
| Códigos HTTP obrigatórios para POST | TS 29.500 V19.7.0 | 5.2.7.1 | `NFs/aanf/internal/sbi/` | negativo | **Especificado** — lista fechada de 11 códigos |
| `3gpp-Sbi-Callback` na notificação | TS 29.500 V19.7.0 | 5.2.3.2 | — | — | Fora do MVP; acrescentar ao plano de conformidade |
| Codificação de `kAkma` / `kaf` | — | — | `akma/codec/` | — | **LACUNA B** — permanece decisão do projeto |
| `expiry` e sincronização de relógio | TS 29.522; TS 29.122 | — | `NFs/aanf/internal/context/` | negativo (desvio de relógio) | **Parcial** — aguarda TS 29.122 |

### B.6.2 Questões em aberto

| Nº | Assunto | Estado |
|---|---|---|
| 1 | Obrigatoriedade da subscrição de `ROAMING_STATUS` | Aberta |
| 5 | `AkmaInd`, `RoutingId`, `AANF` e `naanf-akma` no `free5gc/openapi` | Aberta — inspeção direta do repositório |
| A | Rótulo `atid` versus `akid` em TS 29.522 | Aberta — defeito editorial |
| B | Codificação textual de `kAkma` e `kaf` | **Aberta e provavelmente definitiva.** TS 29.571 não define tipo aplicável |
| C | Codificação do valor do A-TID | Aberta — consultar TS 33.222 e TS 24.109 |
| D | Método HTTP de `trigger-auth` | Aberta — **três fontes** sustentam a adoção de POST |

### B.6.3 Especificações ainda ausentes

| Especificação | Finalidade | Criticidade |
|---|---|---|
| **TS 29.122** | Tipos `Uri` e `DateTime` usados por TS 29.522 | Média — afeta a serialização de `expiry` e `notifUri` |
| **TS 33.222** | Perfis TLS de que os Anexos B de TS 33.535 são adaptação | Média — pode encerrar a LACUNA C |
| **TS 24.109** | Protocolos Ua sobre HTTP; construção do B-TID no GBA | Média — precedente para a LACUNA C |

Nenhuma delas bloqueia as fases F2 a F5.

---

## B.7 Síntese

A incorporação de TS 29.571, TS 29.510, TS 29.500 e TS 29.501 produz quatro resultados:

1. **O registro e a descoberta do AAnF estão integralmente especificados.** `NFType` `AANF`, `ServiceName` `naanf-akma`, `AanfInfo.routingIndicators` e o parâmetro de consulta `routing-indicator` mapeiam a cláusula 6.7 de TS 33.535 no Stage 3 sem lacunas. Esta era a última linha bloqueada da matriz.
2. **O conjunto de códigos HTTP a suportar está fechado** em onze códigos obrigatórios para POST, e o arquivo OpenAPI do Naanf_AKMA declara exatamente esse conjunto — validação cruzada que fornece critério objetivo de completude para os testes negativos.
3. **A DIVERGÊNCIA D está sustentada por três fontes independentes**: o corpo da TS 29.503, o arquétipo de operação customizada de TS 29.501 Anexo C.4, e o padrão de método das figuras de §4.6.1.2. A adoção de POST deixa de ser interpretação isolada.
4. **A LACUNA B provavelmente não será encerrada.** TS 29.571 é a especificação de tipos comuns da série 29 e não define tipo para material de chave. A codificação de `kAkma` e `kaf` permanece decisão de implementação, e assim deve ser documentada no trabalho.

Acrescenta-se um requisito operacional não registrado anteriormente: sendo `expiry` um instante absoluto, o ambiente experimental exige sincronização de relógio entre os componentes, e o plano de testes deve prever um caso de desvio deliberado.
