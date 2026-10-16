# Conferência normativa da documentação

**Data:** 15 de outubro de 2026; **revisão 2:** 16 de outubro de 2026, com TS 33.535, TS 29.535 (e YAML), TS 23.502 e TR 121 917
**Escopo:** conferir, contra o texto das normas do acervo (`specs/`), as afirmações dos documentos `00` a `04`, dos ADRs, do registro de defeitos, da matriz, do material de partida do projeto e dos documentos de 15/10/2026 (`ambiente-verificado.md`, ADR-0001-R2, plano).

---

## 1. Método e cobertura

Os `.docx` foram convertidos em texto preservando títulos de cláusula e linhas de tabela; os OpenAPI foram lidos diretamente. Leitura integral das partes específicas de AKMA e do Ua; nas normas extensas, leitura de todas as ocorrências de AKMA/AAnF e de todas as cláusulas citadas pela documentação.

| Norma | O que foi lido |
|---|---|
| TS 33.220 V19.1.0 | §4.5.2; Anexo B inteiro; Anexo H inteiro |
| TS 33.222 V19.0.0 | §5.2 a §5.4 inteiros; confirmada a ausência de AKMA |
| TS 24.109 V19.0.0 | §5.1 a §5.3; Anexo K inteiro |
| TS 29.522 V19.7.0 | §5.14 inteira; Anexo A.12 inteiro; histórico (CRs de `AKId`) |
| TS 23.003 V19.7.0 | §2.2, §2.2B, §28.2; confirmada a ausência de AKMA, A-KID e A-TID |
| TS 33.501 V19.7.0 | todas as ocorrências de AKMA/AAnF; §6.1.1.1, §6.1.3.1, §6.1.3.2.0, §6.1.5, §6.2.2.1, §13.1.0, §13.4.1.0, §14.2.2, §14.2.6, A.7.0 |
| TS 29.503 V19.7.0 | todas as ocorrências de AKMA/AAnF/`routingId`; §5.2.2.1, §5.2.2.2.10, §5.3.2.16, §5.4.2.2.2, §5.5.2.2.2, Tab. 6.2.3.13.4-*, 6.3.6.2.2-1, 6.3.6.2.3-1, 6.4.6.2.8-1; histórico; YAML de UEAU, UECM, SDM e EE |
| TS 29.509 V19.6.0 | Tab. 6.1.6.2.2-1; YAML de UEAuthentication, SoRProtection e UPUProtection |
| TS 29.510 V19.7.0 | todas as ocorrências de AAnF/AKMA; Tab. 6.1.6.2.73-1; parâmetro `routing-indicator`; YAML de NFManagement e NFDiscovery |
| TS 29.571 V19.7.0 | YAML: `Supi`, `Gpsi`, `DateTime`, `SupportedFeatures`, `VarUeId`, `Bytes` |
| TS 29.500 V19.7.0 | §5.2.3.2.1, §5.2.7.1, §5.2.7.2, §5.2.7.4 |
| TS 29.501 V19.5.0 | §4.4.2, §4.6.1.2, §6.2, Anexo C.4 |
| TS 29.122 V19.7.0 | `DateTime`, `Uri`, `Bytes` (YAML); §5.2.6 |
| TS 33.535 V19.0.0 | **inteira** (cláusulas 1 a 7, Anexos A a E) |
| TS 29.535 V19.5.0 | **inteira**; `TS29535_Naanf_AKMA.yaml` v1.2.0 inteiro |
| TS 23.502 **V19.9.0** | busca de AKMA/AAnF no texto inteiro; §5.2.3.2, §5.2.3.3.2, §5.2.3.5.2 |
| ETSI TR 121 917 V17.0.1 | §16.2 e §16.3 (itens de trabalho AKMA e AKMA_TLS) |
| TS 29.500 V19.7.0 (complemento) | Tab. 5.2.7.2-1 (erros de protocolo e de aplicação) |

**Fora do acervo:** TS 29.505 (semântica de `akmaAllowed` no UDR). **Versão de TS 23.502:** o baseline permanece V19.8.0; a cópia do acervo é a V19.9.0, usada para consulta (Adendo E). As duas observações registradas para TS 23.502 valem na V19.9.0.

---

## 2. Confirmado

| Afirmação | Fonte conferida |
|---|---|
| KDF `S = FC ‖ P0 ‖ L0 ‖ …`, HMAC-SHA-256, `Li` em dois octetos, parâmetro ≤ 65535 octetos | TS 33.220 B.1, B.2.0, B.2.1.3 |
| Cadeias em UTF-8 com NFKC | TS 33.220 B.2.1.2 |
| FC `0x80`–`0x82` alocados a TS 33.535; a KDF está em B.2.0 e B.2.2 é "FC value allocations" | TS 33.220 B.2.2 |
| Identificador Ua\* de 5 octetos; `0x00…00` não separa protocolos; `0xFF` faixa privada sem registro | TS 33.220 H.1, H.2 NOTAS 1 e 2 |
| `01 00 00 00 07` (TLS 1.3), `01 00 01 yy zz` (TLS 1.2), `…08` (OSCORE), `…09` (DTLS) | TS 33.220 H.3 |
| B-TID = `base64encode(RAND)@BSF_servers_domain_name`, RFC 4648; TMPI em base64 | TS 33.220 §4.5.2 passo 7, B.4 |
| Identificador fixo no TLS 1.3 para derivar a chave antes do ClientHello | TS 33.222 §5.4.0.2 NOTA 0 |
| Formato preciso da PSK identity em TS 24.109 | TS 33.222 §5.4.0.1 passo 3, §5.4.0.2 passo 1 |
| TS 33.222 não menciona AKMA | busca no texto |
| K.2: `3gpp-akma` no `User-Agent`; `username` = A-KID; senha = base64(K_AF); `realm` = `3gpp-akma@<FQDN>`; verificação do FQDN | TS 24.109 K.2 |
| Divergência E: K.2 b)2) `3GPP-bootstrapping-akma` × K.2 c)3) `3gpp-akma@FQDN` | TS 24.109 K.2 |
| TLS 1.3: `"3GPP-AKMA" ";" A-KID`; TLS 1.2: prefixo `3GPP-AKMA` e A-KID | TS 24.109 K.3.3, K.3.2 |
| Sem derivação de chave fresca à moda do GBA; a cargo do protocolo Ua\* | TS 24.109 K.2, K.3.2, K.3.3 (NOTAS 2) |
| `AKId`: `"rid<value>.atid<value>@<realm>"`; exemplo com `akid` (divergência A) | TS 29.522 Tab. 5.14.5.4.2-1 |
| CRs 0269 (17.1.0), 0559 (17.6.0), 1784 (19.6.0, CP-260067) sobre `AKId` | TS 29.522, histórico |
| `AkmaAfKeyRequest`, `AkmaAfKeyData`, `ServiceDisableNotif`: atributos, presença e notas | TS 29.522 §5.14.5.3, Anexo A.12 |
| `kaf`, `AfId`, `AKId` apenas `type: string`, sem `pattern` | TS 29.522 Anexo A.12 |
| `notifUri` (`Uri`) e `expiry` (`DateTime`) referenciam TS 29.122 | TS 29.522 Anexo A.12 |
| `DateTime` idêntico em TS 29.122 e TS 29.571 (`format: date-time`) | YAML de ambas |
| `Bytes` = base64 em TS 29.122 e TS 29.571 | YAML de ambas |
| RID: 1 a 4 dígitos, cada dígito significativo, padrão `"0"` | TS 23.003 §2.2B item 3 |
| Domínio doméstico `5gc.mnc<MNC>.mcc<MCC>.3gppnetwork.org`, 3 dígitos, zero à esquerda | TS 23.003 §28.2 |
| TS 23.003 não menciona AKMA, A-KID nem A-TID | busca no texto |
| UDM *shall* incluir indicação AKMA e RID na resposta de `Nudm_UEAuthentication_Get` (EAP-AKA' e 5G-AKA) | TS 33.501 §6.1.3.1 passo 2, §6.1.3.2.0 |
| Saída opcional "AKMA Indication and Routing indicator" | TS 33.501 §14.2.2 |
| `Nudm_UECM_AuthTrigger`: entrada SUPI, sem saídas | TS 33.501 §14.2.6 |
| P0 = IMSI de TS 23.003 §2.2 (ou NAI/GLI/GCI), cadeia conforme B.2.1.2 | TS 33.501 A.7.0 |
| NOTA A de §6.1.5.2 sobre o AAnF e a expiração da chave de AF | TS 33.501 §6.1.5.2 |
| `akmaInd` (booleano, padrão `false`) e `routingId` (`^[0-9]{1,4}$`, condicional a `akmaInd`) | TS 29.503 Tab. 6.3.6.2.3-1; `TS29503_Nudm_UEAU.yaml` |
| AUSF configurado para guardar K_AUSF a preserva após a autenticação primária | TS 29.503 §5.4.2.2.2 |
| AAnF entre os consumidores de `Nudm_SDM`; `id-translation-result` com `af-id`; só `supi` obrigatório | TS 29.503 §5.2.2.1, §5.2.2.2.10; `TS29503_Nudm_SDM.yaml` |
| `trigger-auth`: POST no corpo; `get` com corpo obrigatório no YAML (divergência D); 204, 403 `REAUTHENTICATION_NOT_ALLOWED`, 404 | TS 29.503 Tab. 6.2.3.13.4.1-1 e 6.2.3.13.4.2.2-2; `TS29503_Nudm_UECM.yaml` |
| `AuthTriggerInfo` sem `required` no YAML | `TS29503_Nudm_UECM.yaml` |
| `RoamingStatusReport`: `roaming` e `newServingPlmn` obrigatórios, `accessType`, `purged` | TS 29.503 Tab. 6.4.6.2.8-1 |
| CRs 0523, 1133, 1180, 1253 | TS 29.503, histórico |
| `Kseaf`, `KnrProSe` `[A-Fa-f0-9]{64}`, `Msk` `{128}`, sem âncora; `SorMac`/`UpuMac` ancorados | YAML de TS 29.509 |
| `AANF`, `naanf-akma`, `AanfInfo.routingIndicators` (`minItems: 1`), `aanfInfoList`; `routing-indicator` para AUSF, AANF e UDM | TS 29.510 Tab. 6.1.6.2.73-1 e de descoberta; YAML |
| Campos obrigatórios de `NFProfile` e `NFService` | `TS29510_Nnrf_NFManagement.yaml` |
| `Supi` = `imsi-<imsi>`; `VarUeId`; `SupportedFeatures` hexadecimal | `TS29571_CommonData.yaml` |
| POST obrigatórios: 400, 401, 403, 404, 411, 413, 415, 429, 500, 502, 503 | TS 29.500 Tab. 5.2.7.1-1 |
| URI de custom operation sem recurso; figuras com POST; arquétipo *unsafe and non-idempotent* | TS 29.501 §4.4.2, §4.6.1.2, Anexo C.4 |
| AAnF autônomo ou colocalizado com AUSF ou NEF | TS 33.535 §4.1 NOTA |
| Interfaces AAnF↔AUSF, AAnF↔AF/NEF e AAnF↔UDM com proteção de confidencialidade, integridade e repetição; K_AF com tempo de vida máximo pela política de autenticação do operador | TS 33.535 §4.4.0 |
| "AKMA subscription data": dado que indica se o assinante pode usar AKMA; o UDM o guarda | TS 33.535 §3.1, §4.2.5 |
| UDM *"may"* indicar; se indicar, *"shall"* incluir o RID; AUSF guarda K_AUSF e gera K_AKMA e A-TID após a autenticação primária concluída com sucesso; registro no AAnF com SUPI; AAnF guarda o mais recente | TS 33.535 §6.1 passos 2 a 5 |
| NOTA 1a: em reautenticação o AAnF apaga o par anterior; NOTA 1: o AUSF não precisa guardar material AKMA após a entrega | TS 33.535 §6.1 |
| A-KID em NAI, `username` com RID e A-TID, `realm` com o identificador da rede doméstica; *"The AUSF shall use the RID received from the UDM"* | TS 33.535 §6.1 |
| K_AKMA e A-KID com tempo de vida implícito; K_AF com tempo de vida explícito; K_AF pode continuar até o fim da sessão ou da validade após nova K_AKMA | TS 33.535 §5.2 |
| Autorização do AF por política local **ou** token OAuth2; presença de K_AKMA pelo A-KID; sem K_AKMA, resposta de erro | TS 33.535 §6.2.1 passo 2 |
| Subscrição de `ROAMING_STATUS` com *"shall"* (questão 1); regra de registro duplo | TS 33.535 §6.2.1 passos 5 a 7 |
| Guardar o A-KID transmitido a cada AF (NOTA 1 de §6.8.1; passo 8 de §6.2.1) | TS 33.535 §6.2.1 passo 8, §6.8.1 NOTA 1 |
| Sem procedimento explícito de *refresh* de K_AF; *refresh* de K_AKMA por `AuthTrigger` com linguagem permissiva | TS 33.535 §6.4.3, §6.4.4 |
| Seleção do AAnF pelo NRF ou por configuração local; RID como fator; RID padrão → qualquer instância | TS 33.535 §6.7 |
| Remissão da KDF a "Annex B.2.2 of TS 33.220" (deriva de referência) | TS 33.535 A.1.1 |
| Parâmetros de A.2, A.3 e A.4; SUPI igual ao P0 de TS 33.501 A.7.0; `AF_ID = FQDN ‖ identificador Ua*` | TS 33.535 A.2 a A.4 |
| `psk_identity` do TLS 1.3 descrita só como prefixo `3GPP-AKMA` e A-KID, sem separador (defeito D-E) | TS 33.535 B.1.3.2.2 |
| Operações `register-anchorkey`, `retrieve-applicationkey`, `remove-context`, todas POST, sem recursos | TS 29.535 §5.1.3, §5.1.4 |
| `register-anchorkey`: 200 com `AkmaKeyInfo` no corpo | TS 29.535 §4.2.2.2.2, Tab. 5.1.4.2.2-2 |
| `retrieve-applicationkey`: sem K_AKMA → 403 `K_AKMA_NOT_PRESENT`; 200 com `kaf`, `expiry` e `supi` (sem `anonInd`) | TS 29.535 §4.2.2.3.2 |
| `remove-context`: 204; contexto inexistente → 404 `AKMA_CONTEXT_NOT_FOUND` | TS 29.535 §4.2.2.4.2 |
| `AkmaKeyInfo`: `supi`/`gpsi` mutuamente exclusivos; sem `AKMA_GPSI_Support`, só `supi`; YAML com `oneOf`; Anexo A prevalece | TS 29.535 Tab. 5.1.6.2.2-1, Anexo A.1; YAML |
| *Features* `AKMA_GPSI_Support` (1) e `RoamingRestriction` (2) | TS 29.535 Tab. 5.1.8-1 |
| Erros `K_AKMA_NOT_PRESENT` (403), `ROAMING_AKMA_SERVICE_DENIED` (403, `RoamingRestriction`), `AKMA_CONTEXT_NOT_FOUND` (404) | TS 29.535 Tab. 5.1.7.3-1 |
| OAuth2 *"may"* ser usado, por configuração local; quatro escopos | TS 29.535 §5.1.9; YAML |
| YAML v1.2.0 corresponde a V19.5.0; códigos de POST = 400, 401, 403, 404, 411, 413, 415, 429, 500, 502, 503 (+307, 308) | `TS29535_Naanf_AKMA.yaml` |
| TS 23.502 não lista `AuthTrigger` em `Nudm_UECM` nem o AAnF entre os consumidores de `Nudm_EventExposure_Subscribe` — também na V19.9.0 | TS 23.502 V19.9.0 §5.2.3.2, §5.2.3.5.2 |
| AKMA e AKMA_TLS são itens de trabalho do Release 17 (UID 890030 e 950043) | TR 121 917 V17.0.1 §16.2, §16.3 |

---

## 3. Correções

### 3.1 Precedente mais forte para `kAkma` e `kaf`

`[3GPP]` `TS29503_Nudm_UEAU.yaml`:

```yaml
Kausf:
  type: string
  pattern: '^[A-Fa-f0-9]{64}$'
```

É material de chave de 256 bits, transportado do UDM ao AUSF, com padrão **ancorado** — e é a chave da qual K_AKMA deriva. O Adendo D, o ADR-0002 e o defeito D-D citam apenas TS 29.509, cujos padrões não são ancorados. A decisão de projeto (hexadecimal de 64 caracteres, padrão ancorado) não muda; ganha um precedente direto na mesma cadeia de chaves.

### 3.2 Cláusula da guarda de K_AUSF

`[3GPP]` TS 33.501 §6.1.1.1: *"The KAUSF may be securely stored in the AUSF based on the home operator's policy on using such key e.g. if the control plane solution for Steering of Roaming […] or UE Parameter Update procedures […] or AKMA are supported by the HPLMN."*

`[3GPP]` TS 33.501 §6.2.2.1, "Keys in the AUSF": o AUSF guarda (ou substitui) a K_AUSF quando RES\* e XRES\* coincidem (5G-AKA, §6.1.3.2.0 passo 11) ou quando envia EAP-Success (EAP-AKA', §6.1.3.1 passo 10).

O documento `00` §5.1 e o material de partida atribuem a primeira afirmação a §6.2.2. A segunda é a base normativa dos dois pontos de acionamento do AKMA no AUSF.

### 3.3 `immediateFlag`

`[3GPP]` TS 29.503 §5.5.2.2.2: *"If the NF Service Consumer has included the immediateFlag with value as "true" […] and the event requested for immediate reporting is reported by the UDM (e.g. […] "ROAMING_STATUS"), the UDM **may** include the current status of the event if available in the response."*

O Adendo A §A.5 cita "§5.2.3.5.2" e conclui que `immediateFlag` torna o passo 6 de TS 33.535 §6.2.1 implementável de forma síncrona. A cláusula correta é §5.5.2.2.2, e o texto é permissivo: o estado corrente pode não vir na resposta. Afeta a questão 1, fora do MVP.

### 3.4 Cabeçalhos obrigatórios

`[3GPP]` TS 29.500 V19.7.0 Tab. 5.2.3.2.1-1 lista 21 cabeçalhos customizados que as NFs devem suportar: `3gpp-Sbi-Message-Priority`, `-Callback`, `-Target-apiRoot`, `-Scp-apiRoot`, `-Routing-Binding`, `-Binding`, `-Discovery-*`, `-Producer-Id`, `-Oci`, `-Lci`, `-Client-Credentials`, `-Source-NF-Client-Credentials`, `-Nrf-Uri`, `-Target-Nf-Id`, `-Max-Forward-Hops`, `-Originating-Network-Id`, `-Access-Scope`, `-Access-Token`, `-Target-Nf-Group-Id`, `-Nrf-Uri-Callback`, `-NF-Peer-Info`. O Adendo B §B.3 lista quatro. A maior parte serve à comunicação indireta (SCP), fora do MVP.

### 3.5 Identificadores Ua\* registrados

`[3GPP]` TS 33.220 H.3 inclui também `(0x01,0x00,0x00,0x00,0x00)` a `(…,0x06)`, `(0x01,0x00,0x00,0x01,0x00)` (TMPI) e `(0x01,0x00,0x02,yy,zz)`. Em particular, `(0x01,0x00,0x00,0x00,0x02)`: *"Ua security protocol HTTP digest authentication according to TS 24.109, unless HTTP digest authentication is used in the context of another Ua security protocol, which is already covered elsewhere in this Annex."* O documento `00` §2.4 omite esses valores, e o ADR-0003 não registra por que `…02` não foi escolhido para um Ua\* que é, na forma, HTTP Digest.

### 3.6 Forma do identificador Ua\* no `AfId`

`[3GPP]` TS 29.522 Tab. 5.14.5.4.2-1: *"<Ua\* security protocol id> is a string of 5 octet"*; exemplo `www.app1.com.0100BC0001`. A representação por dez caracteres hexadecimais é inferida do exemplo, que usa maiúsculas; o documento `00` §2.4 a apresenta como `[3GPP]`. `[INTERPRETAÇÃO]` Correta, mas é interpretação. (O valor do exemplo, `01 00 BC 00 01`, não está registrado em H.3 — é só ilustrativo.)

### 3.7 Digest de K.2: RFC 3310 × RFC 7616

`[DIVERGÊNCIA]` TS 24.109 K.2 c): *"HTTP Digest authentication (see RFC 3310 [6]) shall be used"*. RFC 3310 é o HTTP Digest AKA, em que a senha vem do RES. O próprio K.2 b) remete à §5.2, e §5.2.2.1 manda usar RFC 7616 com a chave de aplicação em base64 como senha; TS 33.222 §5.3.0 passo 3 também cita RFC 7235 e RFC 7616, e a NOTA 1b desaconselha MD5. Com senha = base64(K_AF), só RFC 7616 é coerente. Candidato a defeito (CT1). O Adendo D §D.4.1 registra a citação a RFC 3310 sem notar o conflito.

### 3.8 Origem do "204 se não existir"

`[3GPP]` TS 29.522 Tab. 5.14.3.2.2-2 (API `3gpp-akma` do NEF): *"If the requested data does not exist, the NEF shall respond with "204 No Content""*; a mesma API define 403 `K_AKMA_NOT_PRESENT` (Tab. 5.14.7.3-1). O material de partida atribui o 204 ao `ApplicationKey_Get` do Naanf. TS 29.535 confirma que o Naanf tem as duas respostas — e com o mesmo problema: ver 3.11.

### 3.9 Editoriais em TS 29.522

Tab. 5.14.5.1-1 remete `ServiceDisableNotif` à cláusula "5.15.5.3.4" (é 5.14.5.3.4); o texto cita "table 514.5.2-1"; "Routing Indicatorand". Sem efeito técnico.

### 3.10 Divergência E: o Stage 2 pesa a favor de `3GPP-bootstrapping-akma`

`[3GPP]` TS 33.535 B.1.2.2, que adapta TS 33.222 §5.3.0 ao AKMA: *"At step 3, if the AF selects AKMA for deriving the key, then the AF shall include the "3GPP-bootstrapping-akma" within the WWW-Authenticate header field."*; no passo 4, o cliente verifica *"the FQDN in the realm attribute"*.

O quadro completo passa a ser:

| Fonte | `realm` |
|---|---|
| TS 33.535 B.1.2.2 passo 3 (Stage 2, normativo) | contém `3GPP-bootstrapping-akma`, e o FQDN do AF (passo 4) |
| TS 24.109 K.2 b)2) | *"shall contain the constant string "3GPP-bootstrapping-akma""* |
| TS 24.109 §5.2.2.1 item 3 (modelo GBA que K.2 adapta) | `<constante>@<FQDN>`, com a mesma constante que sinaliza o modo, ex. `3GPP-bootstrapping@naf1.operator.com` |
| TS 24.109 K.2 c)3) | primeira parte `3gpp-akma`, ex. `3gpp-akma@af1.operator.com` |

O registro de defeitos (D-C) e o Adendo D decidiram seguir K.2 c)3), argumentando que é o valor que entra no cálculo e o que a alínea d) manda verificar. A alínea d) verifica só a segunda parte (o FQDN), que as duas leituras satisfazem; e o valor que entra no cálculo é o que o AF envia no desafio. Três fontes apontam para `3GPP-bootstrapping-akma@<FQDN>` e uma para `3gpp-akma@<FQDN>`. A decisão foi revista em 16/10/2026: o projeto segue TS 33.535 (`docs/conformidade/desvios.md`, DV-6).

### 3.11 `204 No Content` em `retrieve-applicationkey` sem condição alcançável

`[DIVERGÊNCIA]` TS 29.535 §4.2.2.3.2 e Tab. 5.1.4.3.2-2: *"If the requested AKMA Application Key information for the UE does not exist, the AAnF shall respond with "204 No Content""*. A mesma cláusula manda responder 403 `K_AKMA_NOT_PRESENT` quando a K_AKMA não está presente, e K_AF é sempre derivável quando a K_AKMA está presente (TS 33.535 §6.2.1 passo 7). Não há condição distinta que leve ao 204. Candidato a defeito (CT3); o material de partida registrava as duas respostas lado a lado. Decisão em 16/10/2026: nunca emitir 204 (DV-7).

### 3.12 Guarda do A-KID transmitido a cada AF

`[3GPP]` TS 33.535 §6.2.1 passo 8: *"If AAnF has subscribed the event for RoamingStatusReport, then the AAnF is expected to keep track of the transmitted A-KIDs and the recipient AFs."* O documento `00` §6.3 e o material de partida apresentam essa guarda como parte incondicional do contexto (citando só §6.8.1 NOTA 1). Ela está ligada à subscrição de *roaming*, desabilitada no MVP. Guardá-la não viola nada; deixa de ser requisito do MVP.

### 3.13 Remissão cruzada entre TS 33.535 e TS 23.502

`[DIVERGÊNCIA]` TS 33.535 §4.3.0 NOTA 1 e §7.4 apontam TS 23.502 §5.2.3.5.2 (`Nudm_EventExposure_Subscribe`) como serviço do UDM usado pelo AKMA; essa cláusula, na V19.9.0, lista como consumidores apenas NEF e SMS-GMSC. Reforça a divergência de enumeração já registrada. TS 23.502 V19.9.0 não contém nenhuma ocorrência de AKMA ou AAnF.

### 3.14 Editoriais em TS 29.535

O `CtxRemove` do YAML não declara `required: [supi]`, embora o texto diga que `supi` *"shall be included"* (mesmo padrão de `AuthTriggerInfo` em TS 29.503); `info.title` do YAML é `3gpp-akma` (título da API do NEF); descrições de 307/308 de `register-anchorkey` falam em *"retrieval notification"*; "manadatory".

### 3.15 Correção ao próprio plano (revisão 2)

A revisão 2 do plano apontou `AuthenticationInfo.routingIndicator` (TS 29.509, enviado pelo AMF) como fonte candidata do RID para o AUSF. TS 33.535 §6.1 determina que *"The AUSF shall use the RID received from the UDM"*, e §6.7 NOTA 1 repete que o AUSF obtém o RID na resposta do UDM. Essa fonte está excluída.

---

## 4. Achados novos, sem correção a fazer

| Achado | Fonte | Uso |
|---|---|---|
| O AMF pode enviar o RID ao AUSF: `AuthenticationInfo.routingIndicator`, opcional, `^[0-9]{1,4}$` | TS 29.509 Tab. 6.1.6.2.2-1 | **não** serve ao A-KID: TS 33.535 §6.1 manda usar o RID do UDM (3.15) |
| O AF *"shall be authenticated and authorized by the operator network before"* receber K_AF | TS 33.535 §4.2.2; TS 29.535 §4.2.1.3.2 | DV-1 |
| Aplicações no UE não acessam K_AKMA e só obtêm K_AF dos AF_ID a que estão autorizadas (como, fora do escopo 3GPP) | TS 33.535 §4.4.3 | MVP-4 |
| Iniciação do AKMA: o AF responde ao UE com mensagem de iniciação própria do Ua\* | TS 33.535 §6.5 | MVP-3 |
| Expiração de K_AF: o AF pode recusar ou renovar; com K_AUSF nova, o UE pode tentar com o A-KID novo | TS 33.535 §6.4.2 | MVP-3, MVP-5 |
| No perfil OSCORE (CBOR), o A-KID é estruturado como `[RID : tstr, A-TID : bstr, HPLMN-ID : tstr]`: o A-TID é cadeia de octetos. A forma textual na NAI continua sem codificação definida (lacuna C) | TS 33.535 D.3.5 | ADR-0002 |
| A K_AKMA e o A-TID são gerados *"after the primary authentication procedure is successfully completed"*; a norma não condiciona a autenticação primária ao resultado do registro no AAnF | TS 33.535 §6.1 passo 3 | D9 |
| Corpo `application/json`; erros `application/problem+json` (RFC 9457) | TS 29.535 §5.1.2.2.2 | MVP-1 |
| Erros de protocolo e de aplicação de TS 29.500 Tab. 5.2.7.2-1 suportados para os códigos obrigatórios; tratamento conforme TS 29.501 §4.8 | TS 29.535 §5.1.7.1 | MVP-1 |
| Nenhuma causa específica para AF não autorizado; os 403 genéricos de TS 29.500 (`ACCESS_TOKEN_CLAIM_MISMATCH`, `MISSING_PARAMETER` etc.) não descrevem o caso; `ProblemDetails` é opcional no 403 | TS 29.535 Tab. 5.1.4.3.2-2; TS 29.500 Tab. 5.2.7.2-1 | P3 |
| TS 33.535 V19.0.0 é a 18.8.0 republicada no Release 19 (*"Update to Rel-19 version (MCC)"*), sem mudança técnica | TS 33.535 Anexo E | baseline |
| Pelo título, o CR 0207 (18.3.0) *"AKMA roaming policy control in AAnF"* introduziu o controle de *roaming*; o conteúdo do CR não está no acervo | TS 33.535 Anexo E | questão 1, apenas indício |
| `AuthenticationInfoRequest` (AUSF→UDM) não tem RID; o SUCI vai no URI | TS 29.503 Tab. 6.3.6.2.2-1 | D10 |
| TLS dentro da PLMN *"unless network security is provided by other means"*; NOTA 2 sobre interfaces confiáveis | TS 33.501 §13.1.0 | DV-1, ADR-0001-R2 §5.5 |
| OAuth2 obrigatório de suportar; escopo em nível de serviço como extensão básica; não se aplica a notificações | TS 33.501 §13.4.1.0 | MVP-1 |
| Método não suportado → 501 | TS 29.500 §5.2.7.2 | testes negativos |
| Corpo JSON até 16 milhões de octetos | TS 29.501 §6.2 | teste de 413 |
| Sem `routingIndicators`, o AAnF atende qualquer RID; chave do `aanfInfoList` com até 32 caracteres | TS 29.510 Tab. 6.1.6.2.73-1, Tab. de `NFProfile` | MVP-1 |
| O comprimento do MNC depende do MCC | TS 23.003 §2.2 | domínio doméstico nunca inferido do IMSI |
| O UE envia o nome do servidor no `Host`; consistência do FQDN na derivação | TS 33.222 §5.3.0 passo 2; TS 33.220 §4.5.2 | MVP-3 |
| Proteção de integridade do Digest por `qop=auth-int`, opcional | TS 24.109 §5.2.6 | DV-5 |
| `trigger-auth` 404 com causas `USER_NOT_FOUND` e `CONTEXT_NOT_FOUND` | TS 29.503 Tab. 6.2.3.13.4.2.2-2 | fora do MVP |
| Restrição "P0 é cadeia ASCII estática" vale só para TS 33.220 | TS 33.220 B.2.0 | P10 |
| `AKMA_GPSI_Support` e `RoamingRestriction` têm números diferentes em cada API: na `3gpp-akma` do NEF, `RoamingRestriction` é a *feature* 1 | TS 29.522 Tab. 5.14.6-1 | não confundir com a numeração de TS 29.535 |

---

## 5. Pendente de conferência

Apenas a semântica de `akmaAllowed` e `routingId` no `AuthenticationSubscription` do UDR, que depende de TS 29.505 (fora do acervo). Tudo o mais que a documentação atribui às normas do baseline foi conferido.
