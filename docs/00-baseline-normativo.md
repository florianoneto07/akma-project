# Baseline Normativo para Implementação do 3GPP AKMA no free5GC

**Documento:** Parte 1 do *Roadmap de Implementação do 3GPP AKMA no free5GC*
**Versão:** 2.0 — substitui a v1.0
**Data:** 02 de setembro de 2026
**Alteração principal:** incorporação de TS 29.522, TS 33.220, TS 23.003, TS 33.501 e TS 29.503; encerramento de quatro das cinco lacunas bloqueantes e de três das cinco incertezas registradas na v1.0

---

## Convenções de marcação

| Marcação | Significado |
|---|---|
| `[3GPP]` | Requisito extraído de especificação 3GPP disponível, com indicação de cláusula |
| `[FREE5GC]` | Fato observado no código ou na documentação do free5GC |
| `[MVP]` | Simplificação proposta, com justificativa e impacto declarados |
| `[INTERPRETAÇÃO]` | Leitura adotada pelo projeto quando o texto normativo admite mais de uma |
| `[LACUNA]` | Informação necessária não determinável pelas fontes disponíveis |
| `[DIVERGÊNCIA]` | Inconsistência observada entre fontes ou dentro de uma mesma fonte |

---

## 1. Baseline normativo

### 1.1 Especificações do acervo

| Especificação | Título abreviado | Versão | Data |
|---|---|---|---|
| TS 33.535 | AKMA based on 3GPP credentials in the 5GS | **V19.0.0** | 2025-10 |
| TS 29.535 | AKMA Anchor Services; Stage 3 | **V19.5.0** | 2025-12 |
| TS 29.522 | NEF Northbound APIs | **V19.7.0** | 2026-06 |
| TS 33.501 | Security architecture and procedures for 5G System | **V19.7.0** | 2026-06 |
| TS 33.220 | Generic Authentication Architecture; GBA | **V19.1.0** | 2026-01 |
| TS 23.003 | Numbering, addressing and identification | **V19.7.0** | 2026-06 |
| TS 29.503 | Unified Data Management Services; Stage 3 | **V19.7.0** | 2026-06 |
| TS 23.502 | Procedures for the 5G System | V19.8.0 | 2026-06 |
| `TS29535_Naanf_AKMA.yaml` | OpenAPI do Naanf_AKMA | `1.2.0`, ref. TS 29.535 V19.5.0 | 2025 |

Baseline: **Release 19**.

### 1.2 Encerramento da INCERTEZA 4 (correspondência do arquivo OpenAPI)

`[3GPP]` O arquivo `TS29535_Naanf_AKMA.yaml` versão `1.2.0` declara correspondência com TS 29.535 V19.5.0 e é, portanto, o par correto do corpo normativo do acervo.

Comparação byte a byte com o arquivo anteriormente presente no projeto (`1.3.0-alpha.1`, referente a V20.0.0): **os dois arquivos diferem exclusivamente em três linhas de metadado** — `info.version`, o ano do aviso de copyright e `externalDocs.description`. Todos os caminhos, esquemas, escopos e códigos de resposta são idênticos.

Consequência para a observação registrada na v1.0: a construção `oneOf: [required: supi | required: gpsi]` no esquema `AkmaKeyInfo` **também está presente na versão publicada V19.5.0**, e não é artefato de um rascunho posterior. A aparente tensão com a Tabela 5.1.6.2.2-1, que classifica `supi` e `gpsi` como `C` com cardinalidade `0..1`, resolve-se por TS 29.535, Anexo A.1: *"This Annex takes precedence when being discrepant to other parts of the specification with respect to the encoding of information elements."* **O `oneOf` prevalece**: exatamente um dentre `supi` e `gpsi` deve estar presente. Combinado com a nota da tabela — sem a *feature* `AKMA_GPSI_Support`, apenas `supi` é admitido —, no MVP `supi` é obrigatório e `gpsi` é proibido.

**Decisão de baseline:** adotar TS 29.535 V19.5.0 + `TS29535_Naanf_AKMA.yaml` v1.2.0.

### 1.3 Estado das lacunas da v1.0

| Lacuna (v1.0) | Estado | Fonte que a encerra |
|---|---|---|
| Definição da KDF | **Encerrada** | TS 33.220 V19.1.0, Anexo B.2.0 |
| Codificação do SUPI como P0 | **Encerrada** | TS 33.501 V19.7.0, Anexo A.7.0 |
| Identificador de protocolo de segurança Ua* | **Encerrada** | TS 33.220 V19.1.0, Anexo H |
| Construção do A-KID | **Parcialmente encerrada** | TS 29.522 V19.7.0 §5.14.5.4.2 + TS 23.003 §28.2. Resta a codificação do valor do A-TID (seção 4.3) |
| Tipos `AKId`, `AkmaAfKeyRequest`, `AkmaAfKeyData`, `ServiceDisableNotif` | **Encerrada** | TS 29.522 V19.7.0 §§5.14.5.3, 5.14.5.4 e Anexo A.12 |

Especificações ainda ausentes, de criticidade média: **TS 29.571** (tipos comuns), **TS 29.500 / TS 29.501** (SBA, estrutura de URI, tratamento de erro), **TS 29.510** (NRF, `NFType` do AAnF, OAuth2), **TS 29.122** (tipos `Uri`, `DateTime` usados por TS 29.522), **TS 23.501**.

---

## 2. Funções de derivação de chave — especificação completa

### 2.1 A KDF genérica

`[3GPP]` TS 33.220 V19.1.0, Anexo B.2.0:

```
S = FC ‖ P0 ‖ L0 ‖ P1 ‖ L1 ‖ … ‖ Pn ‖ Ln

derived key = HMAC-SHA-256 ( Key , S )
```

- `FC` é um único octeto, ou dois octetos na forma `0xFF ‖ FC2`. Para AKMA, é sempre um único octeto (faixa `0x80`–`0x82`).
- `Li` é a codificação em **dois octetos** do número de octetos de `Pi`, no intervalo `[0, 65535]`, com o bit menos significativo do octeto string igual ao bit menos significativo de `k` (Anexo B.2.1.3) — isto é, ordem *big-endian*.
- Nenhum parâmetro isolado pode exceder 65535 octetos (Anexo B.1).
- Saída: **256 bits**, por construção do HMAC-SHA-256.

`[3GPP]` TS 33.220 Anexo B.2.1.2 — codificação de cadeias de caracteres: UTF-8 conforme IETF RFC 3629, aplicando **Normalization Form KC (NFKC)**.

`[DIVERGÊNCIA]` TS 33.535 V19.0.0, Anexo A.1.1, remete à "KDF especificada no Anexo **B.2.2** de TS 33.220". Em TS 33.220 V19.1.0, o Anexo B.2.2 intitula-se *"FC value allocations"*; a definição da KDF está no Anexo **B.2.0** (*General*). Trata-se de deriva de referência decorrente de renumeração do Anexo B. A remissão pretendida é inequívoca (Anexo B.2, *Generic key derivation function*), e a implementação deve seguir B.2.0. Registrar na matriz de rastreabilidade.

`[3GPP]` TS 33.535 Anexo A.1.2: os valores de FC alocados a TS 33.535 estão na faixa `0x80 – 0x82`.

### 2.2 Parâmetros das três derivações

`[3GPP]` TS 33.535 V19.0.0, Anexos A.2, A.3 e A.4:

| | **K_AKMA** (A.2) | **A-TID** (A.3) | **K_AF** (A.4) |
|---|---|---|---|
| `FC` | `0x80` | `0x81` | `0x82` |
| `P0` | `"AKMA"` | `"A-TID"` | `AF_ID` |
| `L0` | `0x00 0x04` | `0x00 0x05` | comprimento de `AF_ID` |
| `P1` | `SUPI` | `SUPI` | — |
| `L1` | comprimento do `SUPI` | comprimento do `SUPI` | — |
| `KEY` | `K_AUSF` | `K_AUSF` | `K_AKMA` |
| Saída | 256 bits | 256 bits | 256 bits |

### 2.3 Codificação do SUPI como parâmetro de entrada

`[3GPP]` TS 33.535 Anexos A.2 e A.3 remetem ao parâmetro P0 do Anexo A.7.0 de TS 33.501. TS 33.501 V19.7.0, Anexo A.7.0, estabelece:

- SUPI do tipo **IMSI** → P0 é o **IMSI conforme TS 23.003, cláusula 2.2**;
- SUPI do tipo **network specific identifier** → P0 é o NAI conforme TS 23.003, cláusula 28.7.2;
- SUPI do tipo **GLI** → NAI conforme TS 23.003, cláusula 28.16.2;
- SUPI do tipo **GCI** → NAI conforme TS 23.003, cláusula 28.15.2;
- em todos os casos, P0 é representado como **cadeia de caracteres conforme TS 33.220 B.2.1.2** (UTF-8 + NFKC).

`[FREE5GC]` **Advertência de implementação de alto impacto.** O free5GC representa o SUPI internamente na forma `imsi-208930000000003`. O parâmetro de entrada da KDF é o IMSI conforme TS 23.003 cláusula 2.2, isto é, **apenas os dígitos decimais**, sem o prefixo `imsi-`. Utilizar a cadeia interna do free5GC como P1 produz K_AKMA e A-TID incorretos e, por consequência, um K_AF divergente entre UE e rede. A conversão deve ser explícita, testada e documentada. A mesma observação aplica-se ao prefixo `nai-` nos demais tipos de SUPI.

### 2.4 Identificador de protocolo de segurança Ua* e construção do AF_ID

`[3GPP]` TS 33.535 Anexo A.4: `AF_ID = FQDN do AF ‖ Ua* security protocol identifier`, com o identificador especificado no Anexo H de TS 33.220.

`[3GPP]` TS 33.220 V19.1.0, Anexo H.1: o identificador é uma **cadeia de cinco octetos**. O primeiro octeto identifica a organização que especifica o protocolo; os quatro restantes identificam o protocolo dentro da responsabilidade daquela organização.

`[3GPP]` Anexo H.2 — octeto de organização:

| Valor | Organização |
|---|---|
| `0x00` | Valor padrão para protocolos não especificados de outro modo. Quando usado, apenas o identificador `(0x00,0x00,0x00,0x00,0x00)` deve ser empregado. NOTA 1 registra que protocolos com este identificador não podem ser separados entre si |
| `0x01` | 3GPP |
| `0x02` | 3GPP2 |
| `0x03` | Open Mobile Alliance |
| `0x04` | GSMA |
| `0x01`–`0xFE` | Organizações registradas |
| `0xFF` | **Faixa privada** |

`[3GPP]` NOTA 2 do Anexo H.2: identificadores com `0xFF` no primeiro octeto podem ser usados para definir protocolos locais ou experimentais **sem necessidade de registro**, com a ressalva de que uma falha de segurança em um protocolo sobre Ua* pode ser explorada contra outro protocolo sobre Ua*.

`[3GPP]` Anexo H.3 — identificadores relevantes para os perfis Ua* de TS 33.535:

| Identificador | Protocolo |
|---|---|
| `(0x01,0x00,0x00,0x00,0x07)` | TLS 1.3 com chave compartilhada, TS 33.222 cláusula 5.4.0.2 |
| `(0x01,0x00,0x01,yy,zz)` | TLS 1.2 com chave compartilhada (TS 33.222 §§5.3 e 5.4.0.1), onde `yy,zz` é o código de *CipherSuite* do registro IANA TLS |
| `(0x01,0x00,0x00,0x00,0x08)` | OSCORE, Anexo P de TS 33.220 |
| `(0x01,0x00,0x00,0x00,0x09)` | DTLS, Anexo O de TS 33.220 |

`[MVP]` **Consequência favorável ao MVP.** Um protocolo Ua* simplificado, criado para demonstração, pode utilizar legitimamente um identificador da faixa privada `0xFF`, por exemplo `(0xFF,0x00,0x00,0x00,0x01)`. Isto **não constitui desvio da norma**: o Anexo H.2, NOTA 2, prevê expressamente esse uso. A alternativa de usar o valor padrão `(0x00,0x00,0x00,0x00,0x00)` também é normativa, mas elimina a separação criptográfica entre protocolos e é, portanto, menos adequada.

`[INTERPRETAÇÃO]` **Duas representações distintas de AF_ID.** É preciso distinguir:

1. **AF_ID como parâmetro P0 da KDF** (TS 33.535 A.4): concatenação de FQDN com os cinco octetos do identificador. Por analogia direta com o `NAF_Id` do GBA — TS 33.220 estabelece que "os parâmetros de entrada IMPI e a parte FQDN do NAF_Id devem ser codificados conforme o Anexo B.2.1" —, adota-se: `P0 = UTF-8_NFKC(FQDN) ‖ 5 octetos brutos`, e `L0 = len(UTF-8(FQDN)) + 5`. TS 33.535 não afirma isto de forma explícita; a leitura é do projeto e deve constar como decisão registrada.
2. **`afId` como atributo transportado na API** (TS 29.522 V19.7.0, Tabela 5.14.5.4-1, tipo `AfId`): cadeia formatada como `"<FQDN>.<Ua* security protocol id>"`, com o identificador representado por **dez dígitos hexadecimais**. Exemplo da especificação: FQDN `www.app1.com`, identificador `0100BC0001`, resultando em `www.app1.com.0100BC0001`.

O AAnF recebe a forma (2) e deve reconstruir a forma (1) antes de invocar a KDF. A regra de análise sintática é determinada: separar os dez últimos caracteres hexadecimais precedidos por um ponto. **Nem TS 33.535 nem TS 29.522 declaram explicitamente que a forma (2) não é o P0**; assumir o contrário produziria um K_AF distinto no UE e no AAnF. Esta é a segunda advertência de implementação de alto impacto do documento, e deve ser objeto de teste dedicado.

### 2.5 Exemplo trabalhado — vetor gerado pelo projeto

`[MVP]` Os valores abaixo **não provêm do 3GPP**. TS 33.535 V19.0.0 não fornece vetores de teste. O exemplo serve para validar a implementação da cadeia `S` e é reproduzível a partir dos parâmetros declarados.

Entradas:

```
K_AUSF     = 5c7c9be7bf6a1a1a2c0b0f6d4c9c2f819c1a3b4d5e6f708192a3b4c5d6e7f809
SUPI       = 208930000000003          (IMSI, sem prefixo "imsi-")
FQDN do AF = af.akma.free5gc.org
Ua* id     = FF 00 00 00 01           (faixa privada, TS 33.220 H.2)
```

Cadeias `S` e saídas:

```
K_AKMA   S   = 80 414b4d41 0004 323038393330303030303030303033 000f
         out = d6aa856e127fb8de3eb1a1d0cebf52a7df6b69b4a6c68c21c18c1c27d466dbba

A-TID    S   = 81 412d54494 4 0005 323038393330303030303030303033 000f
         out = 1914b806ff737fe54dd5cd9d3e19cc184a8d2e22cc445ee3fbd9fc3eca329152

K_AF     S   = 82 61662e616b6d612e667265653567632e6f7267 ff00000001 0018
         out = 91df18f308a46e13f032722b1c9448886c59cbae352f1e17e5a1e9d635965daa
```

Leitura da terceira cadeia: `FC = 0x82`; `P0` = os 19 octetos UTF-8 de `af.akma.free5gc.org` seguidos dos 5 octetos `ff00000001`, totalizando 24 octetos; `L0 = 0x0018 = 24`.

Origem dos valores esperados: implementação de referência em Python escrita diretamente a partir de TS 33.220 B.2.0/B.2.1. O procedimento de validação previsto é a comparação byte a byte com a implementação em Go, mais a verificação cruzada de derivações não modificadas contra a implementação em C do artefato AKMA+ (Zenodo, DOI 10.5281/zenodo.14649847).

---

## 3. Construção e formato do A-KID

### 3.1 O que TS 33.535 estabelece

`[3GPP]` TS 33.535 V19.0.0, cláusula 6.1: o A-KID deve estar em formato NAI conforme a cláusula 2.2 da IETF RFC 7542, isto é, `username@realm`. A parte `username` deve incluir o RID e o A-TID; a parte `realm` deve incluir o Home Network Identifier. O A-TID é derivado de K_AUSF conforme o Anexo A.3. O AUSF deve usar o RID recebido do UDM.

### 3.2 TS 23.003 não trata de AKMA

`[3GPP]` **Constatação verificada:** TS 23.003 V19.7.0 não contém nenhuma ocorrência do termo "AKMA", "A-KID" ou "A-TID". A remissão de TS 33.535 §6.7 a TS 23.003 diz respeito exclusivamente ao valor padrão do Routing Indicator.

`[3GPP]` TS 23.003 V19.7.0, cláusula 2.2B — **Routing Indicator**: 1 a 4 dígitos decimais atribuídos pelo operador da rede doméstica e provisionados no USIM. Cada dígito é significativo — o valor `"012"` não equivale a `"12"`. **Se nenhum Routing Indicator estiver configurado no USIM ou no ME, o campo é definido como o valor `0`**, isto é, um único dígito decimal `"0"`. Este é o valor padrão referido em TS 33.535 §6.7.

`[3GPP]` TS 23.003 V19.7.0, cláusula 28.2 — **Home Network Domain** para o 5GC:

```
"5gc.mnc<MNC>.mcc<MCC>.3gppnetwork.org"
```

com `<MNC>` e `<MCC>` de **três dígitos**; se o MNC tiver apenas dois dígitos significativos, insere-se um `"0"` à esquerda. Exemplo da especificação: MCC 345, MNC 12 → `5gc.mnc012.mcc345.3gppnetwork.org`. Para SNPN, a forma é `"5gc.nid<NID>.mnc<MNC>.mcc<MCC>.3gppnetwork.org"`.

### 3.3 O formato concreto está em TS 29.522

`[3GPP]` TS 29.522 V19.7.0, Tabela 5.14.5.4-1, tipo **`AKId`**:

> O AKMA Key Identifier deve estar em formato NAI conforme a cláusula 2.2 da IETF RFC 7542, formatado como `"<username>@<realm>"`, onde `<username>` deve incluir o Routing Indicator e o A-TID no formato `"rid<value>.atid<value>"`, sendo `"rid"` e `"atid"` rótulos que indicam o Routing Indicator e o A-TID, e `<realm>` deve incluir o Home Network Domain.

Exemplo constante da mesma tabela: Routing Indicator `012`, A-TID `019345346`, Home Network Domain `5gc.mnc012.mcc345.3gppnetwork.org`, resultando em

```
rid012.akid019345346@5gc.mnc012.mcc345.3gppnetwork.org
```

### 3.4 Duas questões que permanecem em aberto

`[DIVERGÊNCIA] — DIVERGÊNCIA A: rótulo `atid` versus `akid`.` No **mesmo parágrafo** de TS 29.522 V19.7.0, o texto normativo prescreve o rótulo `"atid"` e o exemplo imediatamente subsequente utiliza `"akid"`. As duas formas são mutuamente incompatíveis para fins de interoperabilidade: o AAnF que analisa o A-KID recebido do AF precisa reconhecer um rótulo determinado.

O histórico de alterações da própria TS 29.522 registra duas correções sobre este tipo — CR 0269 *"Correction to AKId in AKMA API"* (17.1.0), CR 0559 *"AKID encoding clarification"* (17.6.0) — e uma terceira em CR 1784 *"AKId correction"* (19.6.0, CT#111, CP-260067), imediatamente anterior à versão do acervo. A persistência da inconsistência sugere que o exemplo não foi atualizado junto com o texto.

`[INTERPRETAÇÃO]` Adotar o rótulo **`atid`**, por ser o prescrito no texto normativo, sendo o exemplo material ilustrativo. Implementar o analisador sintático do AAnF de forma **tolerante na leitura** — aceitando `atid` e `akid` — e **estrita na escrita** — emitindo apenas `atid`. Registrar a decisão e submeter a questão à discussão na *issue* free5GC #1146 ou diretamente ao 3GPP CT4, dado que se trata de defeito editorial verificável.

`[LACUNA] — LACUNA A: codificação do valor do A-TID na cadeia.` Nenhuma das especificações do acervo declara como o A-TID — saída de 256 bits do Anexo A.3 de TS 33.535 — é representado como texto dentro da parte `username`. As alternativas plausíveis são hexadecimal (64 caracteres), base64url sem preenchimento (43 caracteres) ou outra codificação.

Observe-se que o exemplo de TS 29.522 emprega `019345346`, isto é, **nove dígitos decimais**, o que é incompatível com uma saída de 256 bits sob qualquer das codificações acima. O exemplo é, portanto, meramente sintático e não pode ser usado para inferir a codificação.

Consequência prática: **o A-KID não pode ser construído de forma interoperável enquanto esta lacuna não for fechada.** Duas ações são recomendadas, em paralelo:

1. Consultar TS 24.501 e as demais especificações da série 24 que definam o transporte do A-KID sobre protocolos Ua*, e o Anexo B de TS 33.535 (perfis TLS), onde o A-KID é usado como identidade PSK — a codificação pode estar fixada nesse contexto.
2. Enquanto pendente, parametrizar a codificação no código (`akid.Encoding = hex | base64url`), com **hexadecimal minúsculo como padrão do projeto**, por consistência com a representação de material de chave já usada no free5GC, e declarar explicitamente que se trata de escolha do projeto e não de requisito 3GPP.

### 3.5 Implicação para a seleção do AAnF

`[3GPP]` O formato `rid<value>.atid<value>@<Home Network Domain>` satisfaz o requisito da cláusula 4.4.2 de TS 33.535 — o AF deve poder identificar, a partir do A-KID, o AAnF que atende o UE — porque tanto o Routing Indicator quanto o domínio da rede doméstica são recuperáveis por análise sintática direta da cadeia.

---

## 4. Modelo de dados completo

### 4.1 Tipos definidos por TS 29.535

**`AkmaKeyInfo`** (§5.1.6.2.2; OpenAPI v1.2.0):

| Atributo | Tipo | P | Card. | Observação |
|---|---|---|---|---|
| `supi` | `Supi` | C | 0..1 | `oneOf` com `gpsi`; único admitido sem `AKMA_GPSI_Support` |
| `gpsi` | `Gpsi` | C | 0..1 | `AKMA_GPSI_Support` |
| `aKId` | `AKId` | M | 1 | |
| `kAkma` | `string` | M | 1 | |
| `suppFeat` | `SupportedFeatures` | C | 0..1 | Apenas quando há negociação |

**`CtxRemove`** (§5.1.6.2.3): `supi` (`Supi`, C, 0..1); nesta versão apenas `supi` pode indicar o contexto, e deve ser incluído.

### 4.2 Tipos reutilizados de TS 29.522

`[3GPP]` TS 29.522 V19.7.0, §5.14.5.3 e Anexo A.12.

**`AkmaAfKeyRequest`** (Tabela 5.14.5.3.2-1):

| Atributo | Tipo | P | Card. | Aplicabilidade |
|---|---|---|---|---|
| `afId` | `AfId` | **M** | 1 | |
| `aKId` | `AKId` | **M** | 1 | |
| `anonInd` | `boolean` | O | 0..1 | Padrão `false` se omitido |
| `notifUri` | `Uri` | O | 0..1 | `RoamingRestriction` |
| `suppFeat` | `SupportedFeatures` | C | 0..1 | Apenas quando há negociação |

**`AkmaAfKeyData`** (Tabela 5.14.5.3.3-1):

| Atributo | Tipo | P | Card. | Observação |
|---|---|---|---|---|
| `kaf` | `string` | **M** | 1 | K_AF |
| `expiry` | `DateTime` | **M** | 1 | Tempo de expiração de K_AF |
| `gpsi` | `Gpsi` | O | 0..1 | Ver nota de aplicabilidade abaixo |
| `supi` | `Supi` | C | 0..1 | Ver nota de aplicabilidade abaixo |
| `suppFeat` | `SupportedFeatures` | O | 0..1 | |

`[3GPP]` **`expiry` é do tipo `DateTime`**, não uma duração. O AAnF fornece um instante absoluto de expiração, e não um número de segundos. Esta observação afeta diretamente a implementação da política de tempo de vida de K_AF e o formato de armazenamento no contexto AKMA.

`[3GPP]` **Aplicabilidade dependente da interface** — não se trata de contradição, mas de duas notas distintas sobre o mesmo tipo em interfaces diferentes:

- TS 29.522, Tabela 5.14.5.3.3-1, NOTA 2 (interface **NEF → AF externo**, N33): `gpsi` pode ser incluído e `supi` **não é aplicável**. NOTA 3: se `anonInd` estava definido como `true` na requisição, `gpsi` não deve ser incluído.
- TS 29.535, Tabela 5.1.4.3.2-2, NOTA 2 (interface **AAnF → AF interno/NEF**, N62/N63): `gpsi` **não é aplicável** e `supi` deve ser incluído, exceto no procedimento de acesso anônimo.

Isto é consistente com TS 33.535 §6.3, passo 5: o NEF não deve enviar o SUPI ao AF externo.

**`ServiceDisableNotif`** (Tabela 5.14.5.3.4-1): atributo único `aKId` (`AKId`, M, 1), descrito como *"contém o A-KID previamente transmitido ao AF"* — confirmando, no Stage 3, o requisito de modelo de dados derivado da NOTA 1 da cláusula 6.8.1 de TS 33.535.

**`AfId`** e **`AKId`**: ver seções 2.4 e 3.3.

### 4.3 Lacuna residual: codificação de `kAkma` e `kaf`

`[LACUNA] — LACUNA B.` Tanto `kAkma` (TS 29.535) quanto `kaf` (TS 29.522) são declarados apenas como `type: string`, sem `pattern`, `format` ou nota sobre codificação. Sendo ambos saídas de 256 bits, a representação textual — hexadecimal ou base64 — não está fixada por nenhuma das especificações do acervo.

`[MVP]` Adotar **hexadecimal minúsculo de 64 caracteres**, por consistência com a representação de `Kausf` e `Kseaf` já empregada no free5GC, e declarar a escolha como decisão de implementação. Verificar em TS 29.571, quando disponível, se existe um tipo comum aplicável.

---

## 5. Integração com o UDM: encerramento de duas incertezas

### 5.1 INCERTEZA 5 encerrada — `akmaInd` e `routingId`

`[3GPP]` TS 29.503 V19.7.0, Tabela 6.3.6.2.3-1, tipo `AuthenticationInfoResult`:

| Atributo | Tipo | P | Card. | Descrição |
|---|---|---|---|---|
| `akmaInd` | `boolean` | O | 0..1 | `true`: as chaves AKMA precisam ser geradas para o UE; `false` (padrão): não são necessárias |
| `routingId` | `string` | **C** | 0..1 | **Deve estar presente se `akmaInd` for `true`.** Quando presente, contém o Routing Indicator contido no SUCI (TS 23.003, cláusula 2.2B) |

`[3GPP]` **A obrigatoriedade é mais forte do que TS 33.535 sugere.** TS 33.535 §6.1, passo 2, emprega linguagem permissiva: o UDM *"pode"* indicar ao AUSF se a chave âncora precisa ser gerada. TS 33.501 V19.7.0, nas cláusulas de autenticação primária, estabelece de forma obrigatória:

> Se um assinante possui uma assinatura AKMA, o UDM **deve** incluir a indicação AKMA e o Routing indicator na resposta `Nudm_UEAuthentication_Get`.

A mesma sentença aparece duas vezes em TS 33.501, nos fluxos de 5G-AKA e de EAP-AKA'. TS 33.501 §14.2.2 confirma, entre as saídas opcionais de `Nudm_UEAuthentication_Get`: *"AKMA Indication and Routing indicator, if the subscriber has an AKMA subscription"*.

Conclusão para a implementação: `akmaInd` é opcional no esquema JSON, mas sua omissão para um assinante com assinatura AKMA constitui não conformidade com TS 33.501. O UDM do free5GC deve emiti-lo sempre que o perfil do assinante contiver assinatura AKMA.

`[3GPP]` TS 33.501 V19.7.0, cláusula 6.2.2: K_AUSF pode ser armazenada com segurança no AUSF conforme a política do operador, por exemplo quando SoR, UPU **ou AKMA** são suportados pela HPLMN. Isto fundamenta normativamente a persistência de K_AUSF no AUSF do free5GC, já presente no `AusfUeContext`.

`[3GPP]` TS 29.503 V19.7.0, cláusula 5.4.2 (descrição de `Nudm_UEAuthentication_Get`): *"Se o AUSF estiver configurado para armazenar Kausf (por exemplo, com base em seu suporte a SoRProtection / UPUProtection / derivação de chave AKMA após autenticação primária), o AUSF deve preservar Kausf e informação relacionada (por exemplo, SUPI) após a conclusão da autenticação primária."*

### 5.2 INCERTEZA 2 encerrada — `Nudm_UECM_AuthTrigger` existe e está especificada

`[3GPP]` A operação está definida em **TS 33.501 V19.7.0, cláusula 14.2.6**:

- Nome: `Nudm_UECM_AuthTrigger`
- Descrição: permite que a NF solicite ao UDM que acione uma (re-)autenticação primária conforme a cláusula 6.1.5
- Entrada obrigatória: **SUPI**. Entrada opcional: nenhuma. Saídas: nenhuma

`[3GPP]` O Stage 3 está em **TS 29.503 V19.7.0, cláusulas 5.3.2.16 e 6.2.3.13**:

- `POST {apiRoot}/nudm-uecm/<apiVersion>/{ueId}/registrations/trigger-auth` (operação customizada `trigger-auth`)
- Corpo da requisição: `AuthTriggerInfo` (M, 1), contendo o SUPI
- Sucesso: **204 No Content**
- Usuário inexistente: **404 Not Found** com `ProblemDetails`
- A cláusula nomeia explicitamente o AAnF como consumidor exemplar

`[3GPP]` TS 33.501 V19.7.0, cláusula 6.1.5.2, passo 1a-c: o UDM decide com base em eventos, incluindo *"requisições de NF tais como requisições do AAnF conforme definido em TS 33.535"*. A NOTA A da mesma cláusula registra que o UDM pode decidir não prosseguir com o acionamento, com base em sua política local, e observa que, como o AAnF define a expiração da chave de aplicação segundo a política de autenticação do operador, não há expiração frequente nem risco de sobrecarga de sinalização.

`[DIVERGÊNCIA]` A lista de operações do serviço `Nudm_UECM` em TS 23.502 V19.8.0, cláusula 5.2.3.2, **não inclui** `AuthTrigger`. O histórico de TS 29.503 registra a introdução da operação em CR 1133 (CT#101, 18.3.0). Conclui-se que a enumeração de TS 23.502 não está atualizada. **A divergência apontada na v1.0 deste documento resolve-se em favor de TS 33.501 e TS 29.503**, que são as fontes específicas.

### 5.3 INCERTEZA 3 encerrada — o AAnF é consumidor reconhecido do UDM

`[3GPP]` TS 29.503 V19.7.0, cláusula 5.2.2.1, enumera os consumidores do serviço `Nudm_SubscriberDataManagement`: *"AMF, **AAnF**, SMF, SMSF, GMLC, NEF, 5G DDNMF, NWDAF, DCCF, trusted AF, LMF"*. O histórico registra a inclusão em CR 1180, *"Addition of 'AaNF' in Nudm Subscriber management service"* (CT#102, 18.4.0).

`[3GPP]` A obtenção do GPSI pelo AAnF (TS 33.535 §6.2.1, passos 3–4) mapeia-se, no Stage 3, para a cláusula **5.2.2.2.10 de TS 29.503**: o consumidor — a cláusula cita explicitamente NEF e AAnF — envia um **`GET` ao recurso `IdTranslationResult`**, na forma `/{ueId}/id-translation-result`, onde `{ueId}` é um SUPI ou GPSI, obtendo o identificador correspondente.

Isto refina o registrado na v1.0: a operação relevante não é uma leitura genérica de dados de assinatura, mas a **tradução de identificadores**.

`[DIVERGÊNCIA]` A lista de consumidores de `Nudm_EventExposure_Subscribe` em TS 23.502 V19.8.0, cláusula 5.2.3.5.2, permanece restrita a NEF e SMS-GMSC. O histórico de TS 29.503 registra CR 1253, *"Update the Nudm_EventExposure_Subscribe for AKMA service"* (CT#104, 18.6.0). Conclui-se, novamente, que a enumeração de TS 23.502 não é exaustiva.

`[3GPP]` No Stage 3, o evento é identificado pelo valor de enumeração **`"ROAMING_STATUS"`** e o relatório correspondente é o tipo **`RoamingStatusReport`** (TS 29.503 §6.4.6.2.8). O atributo `immediateFlag` definido como `true` permite que o UDM inclua o estado corrente do evento já na resposta à subscrição, o que corresponde diretamente ao passo 6 da cláusula 6.2.1 de TS 33.535.

### 5.4 INCERTEZA 1 permanece aberta

`[DIVERGÊNCIA]` A questão registrada na v1.0 — se a subscrição de `ROAMING_STATUS` pelo AAnF é incondicional (TS 33.535 §6.2.1 passo 5 usa `shall`) ou condicionada ao suporte da *feature* `RoamingRestriction` — **não foi resolvida** pelas especificações agora disponíveis.

Elemento novo que reforça a segunda leitura: TS 29.522 V19.7.0 marca o atributo `notifUri` de `AkmaAfKeyRequest` como aplicável a `RoamingRestriction`, e o erro `ROAMING_AKMA_SERVICE_DENIED` igualmente. Se toda a reação ao estado de *roaming* está condicionada à *feature*, uma subscrição incondicional produziria informação sem consumidor.

`[MVP]` Manter a recomendação da v1.0: implementar a subscrição como etapa configurável, desabilitada por padrão no MVP, com registro explícito na matriz de rastreabilidade de que se trata de desvio consciente em relação à leitura literal do passo 5.

---

## 6. Alterações revisadas no free5GC

### 6.1 AUSF

`[FREE5GC]` Pontos de integração confirmados por inspeção de `internal/sbi/processor/ue_authentication.go` no repositório `free5gc/ausf`, ramo `main` (mantidos da v1.0):

| Elemento | Papel no AKMA |
|---|---|
| `UeAuthPostRequestProcedure` | Leitura de `akmaInd` e `routingId` da resposta do UDM |
| `ausfUeContext.Kausf` | K_AUSF em representação hexadecimal; insumo dos Anexos A.2 e A.3 |
| `Auth5gAkaComfirmRequestProcedure` | Acionamento do AKMA após sucesso de 5G-AKA |
| `EapAuthComfirmRequestProcedure`, ramo `AKA_CHALLENGE_SUBTYPE` | Acionamento do AKMA após sucesso de EAP-AKA' |
| `ueauth.GetKDFValue` / `ueauth.KDFLen` (`github.com/free5gc/util/ueauth`) | Infraestrutura de KDF |
| `p.Consumer()` em `internal/sbi/consumer/` | Padrão para o novo consumidor `Naanf_AKMA` |

`[FREE5GC]` **Verificação pendente e agora especificável.** A função `ueauth.GetKDFValue(key, FC, P0, L0, …)` deve ser confrontada com TS 33.220 B.2.0 em três pontos: se a saída é integral de 256 bits ou truncada; se `KDFLen` produz a codificação de dois octetos *big-endian* exigida por B.2.1.3; e se as cadeias de caracteres passam por normalização NFKC conforme B.2.1.2. A ausência de NFKC não afeta cadeias ASCII — `"AKMA"`, `"A-TID"`, dígitos de IMSI e FQDN em ASCII —, mas afetaria FQDNs internacionalizados. Registrar o resultado da verificação.

### 6.2 UDM

| Alteração | Fundamento normativo | Categoria |
|---|---|---|
| Emitir `akmaInd` e `routingId` em `Nudm_UEAuthentication_Get Response` para assinantes com assinatura AKMA | TS 33.501 V19.7.0 (linguagem obrigatória) + TS 29.503 §6.3.6.2.3 | `[3GPP]` |
| Adicionar assinatura AKMA ao perfil do assinante (UDR/MongoDB, *webconsole*) | TS 33.535 §4.2.5 | `[3GPP]` |
| Adicionar `AkmaInd` e `RoutingId` ao modelo `AuthenticationInfoResult` em `free5gc/openapi` | Consequência de implementação | `[FREE5GC]` |
| Implementar `GET /{ueId}/id-translation-result` para consumo pelo AAnF | TS 29.503 §5.2.2.2.10; TS 33.535 §6.2.1 passos 3–4 | `[3GPP]`, fora do MVP por política local |
| Implementar `Nudm_EventExposure_Subscribe` com evento `ROAMING_STATUS` | TS 29.503 §6.4.6.2.8; TS 33.535 §6.2.1 passos 5–6 | `[3GPP]`, decisão pendente (INCERTEZA 1) |
| Implementar `POST /{ueId}/registrations/trigger-auth` | TS 33.501 §14.2.6; TS 29.503 §5.3.2.16 | `[3GPP]`, fora do MVP |

`[FREE5GC]` **Origem do Routing Indicator — questão encerrada.** TS 29.503 §6.3.6.2.3 é explícita: `routingId` *"deve incluir o Routing Indicator contido no SUCI"*. O UDM, portanto, extrai o RID do SUCI recebido na requisição de autenticação, e não de dados de assinatura. Quando o SUCI não traz RID configurado, aplica-se o valor padrão de TS 23.003 §2.2B, o dígito único `"0"`. Isto simplifica a alteração no free5GC: o SUCI já é processado no caminho de autenticação, e o RID é um de seus campos.

### 6.3 AAnF

`[3GPP]` Conteúdo do contexto AKMA por UE, agora com os tipos determinados:

```
SUPI              (chave primária; identificador usado por CtxRemove)
A-KID             corrente          — string NAI  "rid<v>.atid<v>@<HN domain>"
K_AKMA            corrente          — 256 bits
GPSI              quando obtido via id-translation-result
por AF atendido:
    AF_ID         forma de transporte ("<FQDN>.<10 hex>") e forma de KDF (FQDN‖5 octetos)
    K_AF          256 bits
    expiry        instante absoluto (DateTime)
    A-KID transmitido                — TS 33.535 §6.8.1 NOTA 1; TS 29.522 ServiceDisableNotif
    notifUri + expiração da subscrição implícita   — TS 33.535 §6.2.1 passo 2
```

`[3GPP]` Comportamento de substituição (TS 33.535 §6.1, NOTA 1a): ao receber novo par (A-KID, K_AKMA) para o mesmo SUPI, o AAnF apaga o par anterior e armazena o novo. Os registros de K_AF por AF permanecem válidos até a expiração, conforme §5.2.

`[LACUNA]` Permanece por determinar o valor de `NFType` do AAnF no perfil registrado no NRF (TS 29.510) e sua presença entre as constantes de `free5gc/openapi`.

### 6.4 AF e Ua*

`[MVP]` Com o Anexo H de TS 33.220 disponível, a especificação do Ua* simplificado do MVP fica completa e normativamente fundamentada:

- Identificador de protocolo de segurança Ua*: faixa privada, `(0xFF,0x00,0x00,0x00,0x01)`, admitida pelo Anexo H.2 NOTA 2
- `afId` transportado: `"<FQDN do AF>.ff00000001"`
- `AF_ID` usado como P0: `UTF-8(FQDN) ‖ 0xFF 0x00 0x00 0x00 0x01`

A evolução para o perfil TLS 1.3 PSK do Anexo B.1.3.2.2 de TS 33.535 implica apenas trocar o identificador por `(0x01,0x00,0x00,0x00,0x07)`, sem alterar a arquitetura.

---

## 7. Matriz de rastreabilidade (revisão 2)

| Requisito | Norma | Cláusula | Implementação prevista | Teste | Estado |
|---|---|---|---|---|---|
| KDF genérica HMAC-SHA-256 | TS 33.220 V19.1.0 | B.2.0, B.2.1 | `akma/kdf/kdf.go` | `kdf_test.go` | **Especificado** |
| Codificação de cadeias (UTF-8 + NFKC) | TS 33.220 V19.1.0 | B.2.1.2 | `akma/kdf/encoding.go` | `encoding_test.go` | **Especificado** |
| Derivação de K_AKMA (FC 0x80) | TS 33.535 V19.0.0 | A.2 | `akma/kdf/kakma.go` | `kakma_test.go` | **Especificado** |
| Derivação de A-TID (FC 0x81) | TS 33.535 V19.0.0 | A.3 | `akma/kdf/atid.go` | `atid_test.go` | **Especificado** |
| Derivação de K_AF (FC 0x82) | TS 33.535 V19.0.0 | A.4 | `akma/kdf/kaf.go` | `kaf_test.go` | **Especificado** |
| SUPI como P0/P1 (IMSI sem prefixo) | TS 33.501 V19.7.0 | A.7.0 | `akma/kdf/supi.go` | `supi_test.go` | **Especificado** |
| Ua* security protocol identifier (5 octetos) | TS 33.220 V19.1.0 | H.1–H.3 | `akma/afid/uastar.go` | `uastar_test.go` | **Especificado** |
| Construção do `AfId` de transporte | TS 29.522 V19.7.0 | 5.14.5.4.2 | `akma/afid/afid.go` | `afid_test.go` | **Especificado** |
| Conversão `AfId` → P0 da KDF | TS 33.535 A.4 + TS 29.522 | — | `akma/afid/afid.go` | `afid_kdf_test.go` | `[INTERPRETAÇÃO]` registrada |
| Home Network Domain | TS 23.003 V19.7.0 | 28.2 | `akma/akid/realm.go` | `realm_test.go` | **Especificado** |
| Routing Indicator e valor padrão | TS 23.003 V19.7.0 | 2.2B | `akma/akid/rid.go` | `rid_test.go` | **Especificado** |
| Construção do A-KID (formato NAI) | TS 29.522 V19.7.0 | 5.14.5.4.2 | `akma/akid/akid.go` | `akid_test.go` | **Parcial** — DIVERGÊNCIA A e LACUNA A |
| Codificação de `kAkma` / `kaf` | — | — | `akma/codec/` | — | **LACUNA B** — decisão do projeto |
| Contexto AKMA do AAnF | TS 33.535 V19.0.0 | 6.6.1, 6.2.1, 6.8.1 | `NFs/aanf/internal/context/` | unitário | **Especificado** |
| `POST /register-anchorkey` | TS 29.535 V19.5.0 | 4.2.2.2.2, 5.1.4.2 | `NFs/aanf/internal/sbi/` | API | **Especificado** |
| `POST /retrieve-applicationkey` | TS 29.535 V19.5.0 + TS 29.522 | 4.2.2.3.2, 5.14.5.3 | `NFs/aanf/internal/sbi/` | API | **Especificado** |
| `POST /remove-context` | TS 29.535 V19.5.0 | 4.2.2.4.2, 5.1.4.4 | `NFs/aanf/internal/sbi/` | API | **Especificado** |
| Erro `K_AKMA_NOT_PRESENT` | TS 29.535 V19.5.0 | 5.1.7.3 | `NFs/aanf/internal/sbi/` | negativo | **Especificado** |
| Erro `AKMA_CONTEXT_NOT_FOUND` | TS 29.535 V19.5.0 | 5.1.7.3 | `NFs/aanf/internal/sbi/` | negativo | **Especificado** |
| Registro do AAnF no NRF | TS 29.510 | — | `NFs/aanf/internal/sbi/consumer/` | integração | **Bloqueado** — norma ausente |
| `akmaInd` / `routingId` no UDM | TS 33.501 V19.7.0 + TS 29.503 §6.3.6.2.3 | — | `patches/udm/`, `patches/openapi/` | integração | **Especificado** |
| Acionamento do AKMA no AUSF (5G-AKA) | TS 33.535 V19.0.0 | 6.1 passo 3 | `patches/ausf/` | integração | **Especificado** |
| Acionamento do AKMA no AUSF (EAP-AKA') | TS 33.535 V19.0.0 | 6.1 passo 3 | `patches/ausf/` | integração | **Especificado** |
| Consumidor `Naanf_AKMA` no AUSF | TS 29.535 V19.5.0 | 4.2.2.2 | `patches/ausf/internal/sbi/consumer/` | integração | **Especificado** |
| Substituição de (A-KID, K_AKMA) | TS 33.535 V19.0.0 | 6.1 NOTA 1a | `NFs/aanf/internal/context/` | integração | **Especificado** |
| GPSI via `id-translation-result` | TS 29.503 V19.7.0 | 5.2.2.2.10 | — | — | Fora do MVP (política local) |
| Subscrição `ROAMING_STATUS` | TS 29.503 V19.7.0 | 6.4.6.2.8 | — | — | **INCERTEZA 1** |
| `Nudm_UECM_AuthTrigger` | TS 33.501 §14.2.6 + TS 29.503 §5.3.2.16 | — | — | — | Fora do MVP |

---

## 8. Registro de incertezas — revisão 2

| Nº | Assunto | Estado |
|---|---|---|
| 1 | Obrigatoriedade da subscrição de `ROAMING_STATUS` | **Aberta.** Reforçada a leitura de condicionamento à *feature* `RoamingRestriction` |
| 2 | Definição de `Nudm_UECM_AuthTrigger` | **Encerrada.** TS 33.501 §14.2.6; TS 29.503 §§5.3.2.16, 6.2.3.13 |
| 3 | AAnF como consumidor do UDM | **Encerrada.** TS 29.503 §5.2.2.1 e §5.2.2.2.10 |
| 4 | Correspondência corpo/OpenAPI da TS 29.535 | **Encerrada.** YAML v1.2.0; diferença de metadado apenas |
| 5 | `AkmaInd`/`RoutingId` no `free5gc/openapi` | **Aberta.** Requer inspeção direta do repositório |
| **A** | Rótulo `atid` versus `akid` em TS 29.522 | **Nova, aberta.** Defeito editorial aparente |
| **B** | Codificação textual de `kAkma` e `kaf` | **Nova, aberta.** Não especificada no acervo |
| **C** | Codificação do valor do A-TID na cadeia do A-KID | **Nova, aberta.** Bloqueia interoperabilidade do A-KID |

---

## 9. Sequência de trabalho revisada

A Fase F0 da v1.0 está substancialmente cumprida. A sequência revisada é:

| Fase | Conteúdo | Critério de conclusão |
|---|---|---|
| **F0′ — Aquisição residual** | TS 29.571, TS 29.500, TS 29.501, TS 29.510, TS 29.122. Fixar *commit hash* do free5GC e dos submódulos | Baseline completo registrado no repositório |
| **F1 — Resolução das três lacunas novas** | LACUNA A (codificação do A-TID): consultar TS 24.501 e o Anexo B de TS 33.535. DIVERGÊNCIA A: decidir e registrar. LACUNA B: decidir e registrar | Documento de decisões de codificação, revisado |
| **F2 — KDF e vetores** | Três derivações, construção do A-KID e do AF_ID, validação cruzada Go × Python | Suíte verde; vetores rotulados quanto à origem |
| **F3 — AAnF autônomo** | Esqueleto, configuração, registro no NRF, três *endpoints*, contexto AKMA, três erros | AAnF registra no NRF; testes cobrem *endpoints* e erros |
| **F4 — Integração AUSF e UDM** | `akmaInd`/`routingId`; derivação no AUSF; consumidor `Naanf_AKMA`; substituição em reautenticação | Contexto AKMA criado no AAnF após registro do UE |
| **F5 — AF e Ua\* simplificado** | AF mínimo com Ua* id de faixa privada | K_AF obtida e verificada contra derivação independente |
| **F6 — UE** | Derivação no lado UE; transporte do A-KID | K_AF do UE idêntica à do AAnF |
| **F7 — End-to-end e testes negativos** | Cenário completo; testes negativos correspondentes a comportamento normativo | Critérios de aceitação satisfeitos |
| **F8 em diante** | Acesso anônimo, remoção de contexto, caminho NEF, *roaming*, `AuthTrigger`, desempenho, verificação formal, *upstream* | Conforme escopo então definido |

A F1 é curta e pode ser conduzida em paralelo com o início da F2, uma vez que as três questões pendentes afetam apenas a serialização, não a derivação criptográfica.

---

## 10. Síntese das alterações em relação à v1.0

**Encerrado:** a definição da KDF; a codificação do SUPI como parâmetro de entrada; o identificador de protocolo de segurança Ua*; os quatro tipos de dados de TS 29.522; a existência e o Stage 3 de `Nudm_UECM_AuthTrigger`; o reconhecimento do AAnF como consumidor do UDM; o mapeamento da obtenção de GPSI para `id-translation-result`; o mapeamento do estado de *roaming* para o evento `ROAMING_STATUS`; a correspondência entre o corpo da TS 29.535 e seu arquivo OpenAPI; e a origem do Routing Indicator.

**Reforçado:** TS 33.501 emprega linguagem obrigatória onde TS 33.535 §6.1 passo 2 é permissiva, o que torna a emissão de `akmaInd` e `routingId` um requisito, e não uma opção, para assinantes com assinatura AKMA.

**Aberto e novo:** a inconsistência entre texto normativo e exemplo no rótulo do A-TID dentro do A-KID (DIVERGÊNCIA A); a ausência de especificação da codificação textual do A-TID (LACUNA C) e do material de chave (LACUNA B). As três dizem respeito exclusivamente à serialização; nenhuma afeta as derivações criptográficas, que estão agora integralmente determinadas.

**Duas advertências de implementação de alto impacto**, ambas capazes de produzir divergência silenciosa de K_AF entre UE e rede: o prefixo `imsi-` na representação interna de SUPI do free5GC (seção 2.3) e a distinção entre a forma de transporte e a forma de KDF do AF_ID (seção 2.4).
