# Adendo C — Tipos comuns de TS 29.122 e reavaliação da LACUNA B

**Documento:** complemento ao *Baseline Normativo* v2.0 e aos Adendos A e B
**Versão:** C.1
**Data:** 02 de setembro de 2026
**Fonte incorporada:** 3GPP TS 29.122 V19.7.0 e `TS29122_CommonData.yaml`

---

## C.1 Escopo

TS 29.122 é referenciada pelo AKMA em dois pontos, ambos no arquivo `TS29522_AKMA.yaml`:

```yaml
notifUri:
  $ref: 'TS29122_CommonData.yaml#/components/schemas/Uri'
expiry:
  $ref: 'TS29122_CommonData.yaml#/components/schemas/DateTime'
```

Adicionalmente, TS 29.535 §5.1.6.1, Tabela 5.1.6.1-2, lista `Uri` de TS 29.122 entre os tipos reutilizados, com aplicabilidade `RoamingRestriction`.

O restante da especificação — APIs T8 de monitoramento, NIDD, provisionamento e afins — não tem relação com o AKMA.

---

## C.2 Encerramento da questão levantada no Adendo B §B.5.3

`[3GPP]` `TS29122_CommonData.yaml`:

```yaml
DateTime:
  format: date-time
  type: string
  description: string with format "date-time" as defined in OpenAPI.

Uri:
  type: string
  description: string providing an URI formatted according to IETF RFC 3986.
```

`[3GPP]` `TS29571_CommonData.yaml`:

```yaml
DateTime:
  format: date-time
  type: string
  description: string with format 'date-time' as defined in OpenAPI.
```

**As duas definições de `DateTime` são idênticas.** A questão registrada no Adendo B §B.5.3 fica encerrada: o atributo `expiry` de `AkmaAfKeyData` é uma cadeia com formato `date-time` da especificação OpenAPI, isto é, um instante absoluto conforme IETF RFC 3339. A escolha entre a definição de TS 29.122 e a de TS 29.571 não produz diferença de serialização.

`[3GPP]` A observação operacional registrada no Adendo B permanece válida e agora está confirmada: sendo `expiry` um instante absoluto, a validade de K_AF depende da concordância de relógio entre AAnF e AF. O testbed deve incluir sincronização de tempo, e o plano de testes negativos deve prever um caso de desvio deliberado de relógio.

`[3GPP]` TS 29.122 V19.7.0, §5.2.4.1: todas as URIs de API T8 têm a forma `{apiRoot}/<apiName>/<apiVersion>`, com `apiRoot` configurado por meios fora do escopo da especificação. Isto confirma a estrutura declarada no arquivo `TS29522_AKMA.yaml`, cujo servidor é `{apiRoot}/3gpp-akma/v1`. Sem efeito sobre o MVP, que não implementa o caminho via NEF.

---

## C.3 Um argumento novo sobre a LACUNA B

`[3GPP]` `TS29122_CommonData.yaml` define o tipo:

```yaml
Bytes:
  type: string
  description: >
    String with format "byte" as defined in OpenAPI Specification, i.e,
    base64-encoded characters.
```

Trata-se do tipo canônico da série 29 para transporte de dados binários em JSON, e é efetivamente usado — por exemplo, o atributo `sms` de `MsisdnLessMoSmsNotification` o referencia.

**A observação relevante é a seguinte.** O arquivo `TS29522_AKMA.yaml` **já importa** `TS29122_CommonData.yaml`, do qual obtém `Uri` e `DateTime`. O tipo `Bytes` estava, portanto, imediatamente disponível para o atributo `kaf`. Ainda assim, `kaf` é declarado como:

```yaml
kaf:
  type: string
```

sem `format`, sem `pattern` e sem `$ref` para `Bytes`. O mesmo vale para `kAkma` em `TS29535_Naanf_AKMA.yaml`.

`[INTERPRETAÇÃO]` Duas leituras são possíveis:

1. **Omissão deliberada.** Se a intenção fosse base64, o tipo `Bytes` teria sido referenciado, como é feito em outros pontos da mesma série. A ausência da referência sugere que a codificação não é base64 e que a especificação a deixou implicitamente a cargo de outra fonte ou da prática de implementação.
2. **Omissão editorial.** O atributo deveria referenciar `Bytes` e não o faz por lapso.

O argumento de disponibilidade favorece a primeira leitura e, por consequência, **reforça a decisão de projeto de adotar hexadecimal**. Trata-se, contudo, de argumento por ausência, e não de evidência positiva. A LACUNA B permanece formalmente aberta.

`[LACUNA]` **Fonte que pode encerrá-la por precedente direto: TS 29.509** (*5G System; Authentication Server Services; Stage 3*). É a especificação que transporta K_SEAF e K_AUSF sobre SBI, no serviço `Nausf_UEAuthentication`, e constitui o análogo mais próximo do transporte de K_AKMA e K_AF. Se TS 29.509 fixar a codificação desses atributos — por exemplo, com um `pattern` hexadecimal —, o precedente resolve a questão para o AKMA por analogia direta.

`[FREE5GC]` Observa-se, adicionalmente, que o AUSF do free5GC já armazena e transporta `Kausf` e `Kseaf` em representação hexadecimal, o que é consistente com a decisão do projeto. Isto é evidência de prática de implementação, não de requisito normativo, e deve ser assim registrado.

**Recomendação:** obter TS 29.509 antes de fixar definitivamente a codificação. Sendo esta a única questão pendente que afeta diretamente a serialização das operações do MVP, o esforço se justifica.

---

## C.4 Estado consolidado das questões em aberto

| Nº | Assunto | Estado | Fonte que pode encerrar |
|---|---|---|---|
| 1 | Obrigatoriedade da subscrição de `ROAMING_STATUS` | Aberta | Histórico de CRs de TS 33.535 |
| 5 | `AkmaInd`, `RoutingId`, `AANF`, `naanf-akma` no `free5gc/openapi` | Aberta | Inspeção do repositório |
| A | Rótulo `atid` versus `akid` em TS 29.522 | Aberta — defeito editorial | 3GPP CT3/CT4 |
| B | Codificação textual de `kAkma` e `kaf` | Aberta, com argumento adicional a favor de hexadecimal | **TS 29.509** |
| C | Codificação do valor do A-TID | Aberta | **TS 33.222**, **TS 24.109** |
| D | Método HTTP de `trigger-auth` | Aberta — três fontes sustentam POST | 3GPP CT4 |
| — | Equivalência de `DateTime` entre TS 29.122 e TS 29.571 | **Encerrada** — definições idênticas | — |

### C.4.1 Especificações ainda ausentes do acervo

| Especificação | Finalidade | Criticidade |
|---|---|---|
| **TS 29.509** | Precedente para a codificação de material de chave sobre SBI (LACUNA B) | Média-alta — é a única pendência que afeta a serialização do MVP |
| **TS 33.222** | Perfis TLS de que os Anexos B de TS 33.535 são adaptação (LACUNA C) | Média |
| **TS 24.109** | Protocolos Ua sobre HTTP; construção do B-TID no GBA (LACUNA C) | Média |

Nenhuma delas bloqueia as fases F2 a F5. As LACUNAS B e C afetam apenas a serialização e podem ser tratadas de forma parametrizada no código, conforme já previsto.

---

## C.5 Situação do baseline documental

Com TS 29.122, o conjunto de especificações **diretamente referenciadas** por TS 33.535 V19.0.0 e TS 29.535 V19.5.0 está completo no acervo:

| Especificação | Versão | Papel no AKMA |
|---|---|---|
| TS 33.535 | V19.0.0 | Stage 2 do AKMA |
| TS 29.535 | V19.5.0 | Stage 3 do Naanf_AKMA |
| TS 29.522 | V19.7.0 | Tipos `AKId`, `AfId`, `AkmaAfKeyRequest`, `AkmaAfKeyData`, `ServiceDisableNotif`; Nnef_AKMA |
| TS 33.501 | V19.7.0 | K_AUSF, SUPI como P0, `Nudm_UECM_AuthTrigger`, autenticação acionada pela rede doméstica |
| TS 33.220 | V19.1.0 | KDF (Anexo B.2), identificador de protocolo Ua* (Anexo H) |
| TS 23.003 | V19.7.0 | Routing Indicator, Home Network Domain |
| TS 29.503 | V19.7.0 | `akmaInd`, `routingId`, `id-translation-result`, `trigger-auth`, `ROAMING_STATUS` |
| TS 29.510 | V19.7.0 | `NFType` `AANF`, `AanfInfo`, descoberta por `routing-indicator` |
| TS 29.571 | V19.7.0 | `Supi`, `Gpsi`, `SupportedFeatures`, `VarUeId` |
| TS 29.500 | V19.7.0 | HTTP/2, códigos de estado, cabeçalhos customizados |
| TS 29.501 | V19.5.0 | Estrutura de URI, operações customizadas, tratamento de erro |
| TS 29.122 | V19.7.0 | `Uri`, `DateTime` |
| TS 23.502 | V19.8.0 | Contexto de procedimentos 5GS |

**A fase F0′ está concluída.** Restam três especificações de consulta pontual — TS 29.509, TS 33.222 e TS 24.109 —, todas destinadas a encerrar questões de serialização que não bloqueiam a implementação.

O caminho crítico do projeto passa agora da aquisição documental para a execução: fase F1, curta, para registrar as decisões de codificação pendentes, e fase F2, implementação das derivações e do A-KID com validação cruzada.
