# Registro de defeitos e divergências em especificações 3GPP

**Escopo:** inconsistências verificáveis identificadas durante a análise normativa do AKMA Release 19.
**Estado:** aberto para submissão. Nenhum item foi ainda submetido ao 3GPP.

Cada item registra a inconsistência, as fontes que a evidenciam, a decisão adotada pelo projeto e o grupo de trabalho responsável.

**Revisão:** 16 de outubro de 2026 — conferência contra o texto das normas (`docs/baseline/conferencia-normativa.md`); D-C revisto; D-D complementado; D-F a D-I acrescentados. 22 de outubro de 2026 — D-J a D-N, do OSCORE do AKMA (TS 33.535 Anexo D; `docs/af-side/ua-star-norma.md` §7.2); D-A a D-I conferidos de novo contra os textos de `specs/` (abaixo); D-B passa a MVP próprio. 28 de outubro de 2026 — D-B: especificação do MVP-12 aprovada. Estes defeitos são a base das recomendações à 3GPP da V2 do repositório público (plano, MVP-13).

### Situação no projeto (22/10/2026)

Todos os itens estão no texto das normas; nenhum é defeito do projeto. A decisão de cada um:

| Item | Implementação |
|---|---|
| D-A | `shared/akma/akid`: lê `atid` e `akid`, escreve `atid`; o AAnF, o AF e o NEF usam essa leitura |
| D-B | ainda não: o projeto não usa o `trigger-auth`, e o UDM e o AMF do free5GC não têm a operação. MVP-12 (decisão do responsável de 22/10/2026), com a especificação aprovada em 28/10/2026 (`docs/free5gc-side/renovacao-kakma.md`) |
| D-C | DV-6: `realm="3GPP-bootstrapping-akma@<FQDN>"` no AF e no cliente |
| D-D | `shared/akma/codec`: `^[A-Fa-f0-9]{64}$` ao ler, minúsculas ao escrever, no AUSF, no AAnF, no AF e no cliente |
| D-E | `3GPP-AKMA;<A-KID>` nos perfis PSK (U2) |
| D-F | DV-5: RFC 7616 com SHA-256 |
| D-G | DV-7: o AAnF nunca responde 204; o AF e o NEF tratam o 204 |
| D-H | editorial, nada a implementar |
| D-I | o AAnF recusa `CtxRemove` sem `supi` (400) |
| D-J a D-N | fase (d) do MVP-8 (U24 a U26, U29) |

### Correspondência com as letras das questões em aberto

Os documentos `00` a `04` e o README usam letras A–E para questões em aberto; este registro usa D-A a D-N para defeitos. As letras não coincidem:

| Questão (README, adendos) | Assunto | Item deste registro |
|---|---|---|
| A | Rótulo `atid` × `akid` | D-A |
| B | Codificação de `kAkma` e `kaf` | D-D |
| C | Codificação do valor do A-TID | — (lacuna, sem defeito; decisão no ADR-0002) |
| D | Método HTTP de `trigger-auth` | D-B |
| E | Dois valores de `realm` | D-C |
| — | Separador da `psk_identity` | D-E |
| — | RFC 3310 no Digest do AKMA | D-F |
| — | 204 inalcançável em `retrieve-applicationkey` | D-G |
| — | Editoriais de TS 29.522 e TS 29.535 | D-H, D-I |
| — | OSCORE do AKMA: remissões, campo sem codificação, OSC-INP, HKDF, HPLMN-ID | D-J a D-N |

---

## D-A — Rótulo `atid` versus `akid` no formato do A-KID

| | |
|---|---|
| **Especificação** | TS 29.522 V19.7.0, Tabela 5.14.5.4-1, tipo `AKId` |
| **Grupo responsável** | 3GPP CT3 |

**Descrição.** No mesmo parágrafo, o texto normativo prescreve que a parte `username` do A-KID tenha o formato `"rid<value>.atid<value>"`, com `"atid"` como rótulo do A-TID, enquanto o exemplo imediatamente subsequente emprega `"akid"`:

```
rid012.akid019345346@5gc.mnc012.mcc345.3gppnetwork.org
```

**Impacto.** As duas formas são mutuamente incompatíveis para interoperabilidade. O AAnF que analisa o A-KID recebido do AF precisa reconhecer um rótulo determinado.

**Histórico relevante.** A mesma especificação registra três correções sobre este tipo: CR 0269 (17.1.0), CR 0559 (17.6.0) e CR 1784 (19.6.0, CP-260067). A persistência da inconsistência sugere que o exemplo não foi atualizado junto com o texto.

**Decisão do projeto.** Adotar o rótulo `atid`, prescrito no texto normativo. O analisador sintático do AAnF será tolerante na leitura, aceitando `atid` e `akid`, e estrito na escrita, emitindo apenas `atid`.

---

## D-B — Método HTTP da operação `trigger-auth`

| | |
|---|---|
| **Especificação** | `TS29503_Nudm_UECM.yaml` v1.4.2, Anexo A de TS 29.503 V19.7.0 |
| **Grupo responsável** | 3GPP CT4 |

**Descrição.** O corpo da TS 29.503, Tabela 6.2.3.13.4.1-1, mapeia a operação customizada `trigger-auth` para o método **POST**. O arquivo OpenAPI declara `get:`, com `requestBody` marcado `required: true`.

**Impacto.** Um `GET` com corpo de requisição obrigatório contraria a semântica de HTTP e não é tratado de forma previsível por servidores e clientes HTTP/2. O Anexo A.1 de TS 29.503 contém cláusula de precedência que cobre explicitamente *"methods"*, de modo que a aplicação literal da regra faria prevalecer uma construção não implementável.

**Fontes que sustentam POST.**

1. TS 29.503 V19.7.0, Tabela 6.2.3.13.4.1-1, e a Tabela 6.2.3.13.4.2.2-1, intitulada *"Data structures supported by the POST Request Body on this resource"*.
2. TS 29.501 V19.5.0, §4.6.1.2 — as duas figuras normativas de operações customizadas são intituladas *"Custom Operation … using HTTP POST"*.
3. TS 29.501 V19.5.0, Anexo C.4 — o arquétipo de operação customizada modela operação *"unsafe and non-idempotent"*. O método GET é seguro e idempotente por definição de IETF RFC 9110.

**Decisão do projeto.** Adotar POST. A operação (TS 33.535 §6.4.4, opcional) entra no MVP-12 (decisão do responsável de 22/10/2026; especificação aprovada em 28/10/2026, `docs/free5gc-side/renovacao-kakma.md`): o UDM e o AMF do free5GC não a têm, e ela exige a autenticação primária disparada pela rede doméstica (TS 33.501 §6.1.5).

**Observação secundária.** O esquema `AuthTriggerInfo` não declara `required`, embora TS 33.501 §14.2.6 declare *"Input, Required: SUPI"* e a Tabela 6.2.3.13.4.2.2-1 classifique a estrutura como `M`.

---

## D-C — Dois valores de `realm` no perfil HTTP Digest do AKMA

| | |
|---|---|
| **Especificação** | TS 24.109 V19.0.0, Anexo K, cláusula K.2 |
| **Grupo responsável** | 3GPP CT1 |

**Descrição.** A cláusula especifica dois valores distintos para o parâmetro `realm` do mesmo intercâmbio HTTP Digest:

- item b) 2): o `realm` do cabeçalho `WWW-Authenticate` deve conter a cadeia constante `"3GPP-bootstrapping-akma"`;
- item c) 3): o parâmetro `realm` deve conter duas partes delimitadas por `@`, sendo a primeira a cadeia constante `"3gpp-akma"` e a segunda o FQDN do AF.

**Impacto.** Em HTTP Digest, o cliente ecoa no cabeçalho `Authorization` o `realm` recebido no desafio, e o valor participa do cálculo de `HA1`. Os dois valores devem coincidir.

**Evidência acrescentada em 16/10/2026.** TS 33.535 V19.0.0 B.1.2.2 (Stage 2, normativo), que adapta TS 33.222 §5.3.0: no passo 3, *"the AF shall include the "3GPP-bootstrapping-akma" within the WWW-Authenticate header field"*; no passo 4, o cliente verifica *"the FQDN in the realm attribute"*. O modelo GBA de TS 24.109 §5.2.2.1 item 3 usa `<constante>@<FQDN>` com a mesma constante que sinaliza o modo (`3GPP-bootstrapping@naf1.operator.com`). Três fontes apoiam `3GPP-bootstrapping-akma@<FQDN>`; só c) 3) usa `3gpp-akma`, que coincide com o *product token* do `User-Agent` — o resíduo editorial parece estar em c) 3), não em b) 2).

**Decisão do projeto** (revista em 16/10/2026, desvio DV-6). Seguir TS 33.535: `realm = "3GPP-bootstrapping-akma@<FQDN do AF>"`. ~~Seguir o item c) 3) — `"3gpp-akma@<FQDN do AF>"` —, que é o valor que participa do cálculo do digest e o que a alínea d) da mesma cláusula manda verificar.~~ A alínea d) verifica só o FQDN, que as duas leituras satisfazem.

---

## D-D — Ausência de `pattern` nos atributos de material de chave

| | |
|---|---|
| **Especificações** | TS 29.535 V19.5.0 (`kAkma`) e TS 29.522 V19.7.0 (`kaf`) |
| **Grupo responsável** | 3GPP CT3 e CT4 |

**Descrição.** Os atributos `kAkma` e `kaf` são declarados como `type: string`, sem `pattern`, `format` ou nota sobre codificação, embora ambos transportem chaves de 256 bits.

**Precedente contrário.** TS 29.509 V19.6.0 declara material de chave de comprimento fixo com padrão hexadecimal explícito: `Kseaf` e `KnrProSe` como `[A-Fa-f0-9]{64}`, `Msk` como `[A-Fa-f0-9]{128}`, `ResStar` e `HxresStar` como `[A-Fa-f0-9]{32}`. As mesmas especificações reservam base64 — via `Bytes` ou `format: byte` — para contêineres opacos de comprimento variável.

**Impacto.** A ausência de padrão deixa a codificação indeterminada, com risco de divergência entre implementações.

**Decisão do projeto.** Codificar `kAkma` e `kaf` como cadeias hexadecimais minúsculas de 64 caracteres, por precedente direto de TS 29.509. Aplicar padrão ancorado, `^[A-Fa-f0-9]{64}$`.

**Precedente acrescentado em 16/10/2026.** `TS29503_Nudm_UEAU.yaml` define `Kausf` como `^[A-Fa-f0-9]{64}$` — ancorado, na interface UDM→AUSF, e é a chave da qual K_AKMA deriva. É o precedente mais próximo.

**Observação correlata.** Os padrões em `TS29509_Nausf_UEAuthentication.yaml` não são ancorados — `'[A-Fa-f0-9]{64}'`, sem `^` e `$` —, ao contrário dos arquivos de SoR e UPU da mesma especificação. Uma validação literal aceitaria cadeias com prefixo e sufixo arbitrários.

---

## D-E — Omissão do separador na `psk_identity` do TLS 1.3

| | |
|---|---|
| **Especificação** | TS 33.535 V19.0.0, Anexo B.1.3 |
| **Grupo responsável** | 3GPP SA3 |

**Descrição.** O Anexo B.1.3 de TS 33.535 descreve a `psk_identity` como composta do prefixo `"3GPP-AKMA"` e do A-KID, sem mencionar caractere separador. TS 24.109 V19.0.0, cláusula K.3.3, especifica que os parâmetros `psk_identity` devem conter o prefixo, **um caractere separador `";"`** e o A-KID.

**Impacto.** Uma implementação que siga apenas TS 33.535 produz `psk_identity` em formato incompatível com o Stage 3. O formato normativo é:

```
"3GPP-AKMA" ";" A-KID
```

**Observação correlata.** TS 33.535 também não menciona que a senha do perfil HTTP Digest é K_AF codificada em base64 conforme IETF RFC 4648, requisito estabelecido em TS 24.109 K.2 c) 2).

**Decisão do projeto.** Seguir TS 24.109, que é a fonte Stage 3 do ponto de referência Ua*. TS 33.222 §5.4.0.2 confirma essa atribuição ao declarar que o formato preciso da PSK identity é especificado em TS 24.109.

---

## D-F — Referência a RFC 3310 no Digest do AKMA

| | |
|---|---|
| **Especificação** | TS 24.109 V19.0.0, Anexo K, cláusula K.2 c) |
| **Grupo responsável** | 3GPP CT1 |

**Descrição.** K.2 c): *"HTTP Digest authentication (see RFC 3310 [6]) shall be used"*. RFC 3310 é o HTTP Digest AKA, em que a senha deriva do RES do AKA. O mesmo K.2 b) remete à cláusula 5.2, cuja §5.2.2.1 manda usar RFC 7616 com a chave de aplicação em base64 como senha; TS 33.535 B.1.2.2 segue TS 33.222 §5.3.0, que cita RFC 7235 e RFC 7616 e desaconselha MD5 (NOTA 1b).

**Impacto.** Com senha = base64(K_AF), apenas RFC 7616 é aplicável.

**Decisão do projeto** (desvio DV-5). Digest conforme RFC 7616, `algorithm=SHA-256`.

---

## D-G — `204 No Content` sem condição alcançável em `retrieve-applicationkey`

| | |
|---|---|
| **Especificação** | TS 29.535 V19.5.0, §4.2.2.3.2 e Tabela 5.1.4.3.2-2; `TS29535_Naanf_AKMA.yaml` v1.2.0 |
| **Grupo responsável** | 3GPP CT3 |

**Descrição.** A cláusula manda responder 403 `K_AKMA_NOT_PRESENT` quando a K_AKMA identificada pelo A-KID não está presente e, também, *"If the requested AKMA Application Key information for the UE does not exist, the AAnF shall respond with "204 No Content""*. Com a K_AKMA presente, o AAnF deriva a K_AF (TS 33.535 §6.2.1 passo 7). Não resta condição que leve ao 204. A mesma construção aparece na API `3gpp-akma` do NEF (TS 29.522 Tabela 5.14.3.2.2-2).

**Impacto.** Consumidores diferentes podem esperar respostas diferentes para o mesmo caso.

**Decisão do projeto** (desvio DV-7). Nunca emitir 204; responder 403 `K_AKMA_NOT_PRESENT` quando não há K_AKMA. Os clientes do projeto tratam 204 mesmo assim.

---

## D-H — Editoriais em TS 29.522 (prioridade baixa)

| | |
|---|---|
| **Especificação** | TS 29.522 V19.7.0 |
| **Grupo responsável** | 3GPP CT3 |

Tabela 5.14.5.1-1 remete `ServiceDisableNotif` à cláusula "5.15.5.3.4" (é 5.14.5.3.4); o texto cita "table 514.5.2-1"; "Routing Indicatorand" na Tabela 5.14.5.4.2-1. Sem efeito técnico.

---

## D-I — Editoriais em TS 29.535 (prioridade baixa)

| | |
|---|---|
| **Especificação** | TS 29.535 V19.5.0; `TS29535_Naanf_AKMA.yaml` v1.2.0 |
| **Grupo responsável** | 3GPP CT3 |

O esquema `CtxRemove` não declara `required: [supi]`, embora o texto (Tabela 5.1.6.2.3-1) diga que `supi` *"shall be included"*; `info.title` do YAML é `3gpp-akma` (título da API do NEF); as respostas 307/308 de `register-anchorkey` falam em *"retrieval notification"*; "manadatory" nas notas das tabelas. Impacto técnico só no primeiro: a validação pelo esquema aceita `CtxRemove` vazio — a implementação deve exigir `supi`.

---

## D-J — Remissões erradas no OSCORE do AKMA

| | |
|---|---|
| **Especificação** | TS 33.535 V19.0.0, Anexo D, cláusulas D.3.2 e D.3.3 |
| **Grupo responsável** | 3GPP SA3 |

**Descrição.** D.3.2 dá a semântica dos parâmetros do pedido por *"TS 33.220 [4], clause Y.2.3. Step 1"* e a de N2 e UE-SID por *"clause P.2.3. Step 3"*; D.3.3 parte do contexto de *"TS 33.220 [4], clause Y.3.3"*. TS 33.220 V19.1.0 não tem Anexo Y, e P.2.3 é *"Requirements on the NAF"*, sem passos. O conteúdo citado está em P.3.2 (passos 1 e 3) e P.3.3.

**Impacto.** Editorial: o conteúdo se acha pelo assunto. Mas D.3.3 só dá parte do contexto e deixa o resto (versão, HKDF, AEAD, *ID Context*) para a cláusula citada, que precisa ser a certa.

**Decisão do projeto.** Ler as remissões como P.3.2 e P.3.3.

---

## D-K — Campo do pedido OSCORE sem codificação

| | |
|---|---|
| **Especificação** | TS 33.535 V19.0.0, D.3.2 iii) e D.3.5; TS 33.220 V19.1.0, P.3.2 e P.3.5 |
| **Grupo responsável** | 3GPP SA3 |

**Descrição.** D.3.2 iii) lista no pedido *"CoAP Security protocol identifier, A-KID, N1, AF-SID, ?OSC-INP"*. O arranjo CBOR de D.3.5 é `[A-KID, N1, AF-SID, ? OSC-INP]`, sem o identificador. O Anexo P do GBA tem a mesma diferença: P.3.2 dá ao campo o valor `"01"` para o OSCORE, e P.3.5 não o codifica.

**Impacto.** Uma implementação que siga D.3.2 manda um arranjo que outra, que siga D.3.5, não lê.

**Decisão do projeto** (U24). Seguir o arranjo de D.3.5, sem o campo.

---

## D-L — Tipo e negociação do OSC-INP

| | |
|---|---|
| **Especificação** | TS 33.535 V19.0.0, D.3.5; TS 33.220 V19.1.0, P.3.3 e P.3.5 |
| **Grupo responsável** | 3GPP SA3 |

**Descrição.** No arranjo, `? OSC-INP: bstr`; na definição, `OSC-INP = {…}`, um mapa CBOR. Os rótulos (1 versão, 3 HKDF, 4 AEAD, 5 sal, 6 *ID Context*) são os do `OSCORE_Input_Material` da RFC 9203 §3.2.1, que admite texto ou inteiro para o HKDF e o AEAD; o Anexo D só admite inteiro. O sal (5) concorre com o *Master Salt* que D.3.3 define. A resposta, `[N2, UE-SID]`, não tem como dizer ao UE qual valor o AF adotou.

**Impacto.** Um OSC-INP como mapa e outro como `bstr` com o mapa dentro não interoperam. Um AF que não aceite o algoritmo pedido não tem como recusá-lo dentro do procedimento.

**Decisão do projeto** (U29). `bstr` com o mapa dentro; o UE não manda OSC-INP; o AF só aceita os valores padrão e recusa os demais com `4.00`.

---

## D-M — Parâmetros do HKDF do *Master Secret*

| | |
|---|---|
| **Especificação** | TS 33.535 V19.0.0, D.3.3; TS 33.220 V19.1.0, P.3.3 |
| **Grupo responsável** | 3GPP SA3 |

**Descrição.** *"OMS = OSCORE Master Secret = HKDF(KAF, "AKMA-OSCORE")"*, com HKDF *"specified in IETF RFC 5869"*. A RFC 5869 define HKDF(sal, IKM, info, L) em duas etapas. O texto não diz o sal, a codificação do rótulo nem o comprimento L do *Master Secret*. O Anexo P tem a mesma forma, com `"GBA-OSCORE"`.

**Impacto.** Escolhas diferentes de sal, de codificação ou de L dão *Master Secrets* diferentes, e o primeiro pedido protegido não decifra.

**Decisão do projeto** (U26). HKDF-SHA-256 completo, sal vazio, IKM = K_AF, `info` = `AKMA-OSCORE` em ASCII, L = 32.

---

## D-N — Forma do HPLMN-ID no A-KID em CBOR

| | |
|---|---|
| **Especificação** | TS 33.535 V19.0.0, D.3.5 |
| **Grupo responsável** | 3GPP SA3 |

**Descrição.** `A-KID = [RID : tstr, A-TID : bstr, HPLMN-ID : tstr]`, sem a forma do HPLMN-ID: o MCC e o MNC (o *Home Network Identifier* de TS 23.003 §2.2B), o domínio do *realm* ou outra. Com o A-TID em octetos, o AF tem de refazer o texto do A-KID para consultar o AAnF, e a codificação textual do A-TID também não está fixada (questão C, ADR-0002).

**Impacto.** UE e AF com formas diferentes não chegam ao mesmo A-KID. Um AF e um AAnF com codificações diferentes do A-TID também não.

**Decisão do projeto** (U25). HPLMN-ID = MCC seguido do MNC; o AF refaz o *realm* com o MNC em três dígitos e o A-TID em base64, como o resto do projeto.

---

## Divergências de enumeração em TS 23.502

Registradas para rastreabilidade; não constituem defeito submetível, por serem listas de contexto em especificação de Stage 2.

| Item | Descrição |
|---|---|
| TS 23.502 V19.8.0 §5.2.3.2 | A lista de operações do serviço `Nudm_UECM` não inclui `AuthTrigger`, introduzida em TS 29.503 pelo CR 1133 (18.3.0) |
| TS 23.502 V19.8.0 §5.2.3.5.2 | A lista de consumidores de `Nudm_EventExposure_Subscribe` inclui apenas NEF e SMS-GMSC; TS 29.503 §5.2.2.1 lista o AAnF entre os consumidores do UDM, incluído pelo CR 1180 (18.4.0) |
| TS 33.535 V19.0.0 §4.3.0 NOTA 1 e §7.4 × TS 23.502 §5.2.3.5.2 | A TS 33.535 aponta `Nudm_EventExposure_Subscribe` (TS 23.502 §5.2.3.5.2) como serviço do UDM usado pelo AKMA; a cláusula apontada não lista o AAnF como consumidor. Inconsistência entre as duas normas (acrescentado em 16/10/2026) |

As duas primeiras linhas foram conferidas também na V19.9.0, cópia disponível no acervo (Adendo E).
