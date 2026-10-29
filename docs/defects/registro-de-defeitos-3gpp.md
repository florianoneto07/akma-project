# Registro de defeitos e divergências em especificações 3GPP

**Escopo:** inconsistências verificáveis identificadas durante a análise normativa do AKMA Release 19.
**Estado:** aberto para submissão. Nenhum item foi ainda submetido ao 3GPP.

Cada item registra a inconsistência, as fontes que a evidenciam, a decisão adotada pelo projeto e o grupo de trabalho responsável.

**Revisão:** 16 de outubro de 2026 — conferência contra o texto das normas (`docs/baseline/conferencia-normativa.md`); D-C revisto; D-D complementado; D-F a D-I acrescentados. 22 de outubro de 2026 — D-J a D-N, do OSCORE do AKMA (TS 33.535 Anexo D; `docs/af-side/ua-star-norma.md` §7.2); D-A a D-I conferidos de novo contra os textos de `specs/` (abaixo); D-B passa a MVP próprio. 28 de outubro de 2026 — D-B: especificação do MVP-12 aprovada. Estes defeitos são a base das recomendações à 3GPP da V2 do repositório público (plano, MVP-13). Na implementação do MVP-12, D-B implementado (DV-12) e D-O e D-P acrescentados, de TS 29.503. No MVP-13: todos os itens conferidos nas versões mais recentes das normas (seção seguinte), e as lacunas que o projeto teve de decidir acrescentadas como sugestões, D-Q a D-Z. 29 de outubro de 2026 — D-AA a D-AC, do CAPIF, achados na proposta do MVP-9, aprovada nessa data.

### Situação no projeto (22/10/2026)

Todos os itens estão no texto das normas; nenhum é defeito do projeto. A decisão de cada um:

| Item | Implementação |
|---|---|
| D-A | `shared/akma/akid`: lê `atid` e `akid`, escreve `atid`; o AAnF, o AF e o NEF usam essa leitura |
| D-B | DV-12: POST no UDM (`free5gc/patches/udm/0002`), no cliente do `openapi` (`0003`) e no AAnF; GET recebe 405. MVP-12 (28/10/2026, `docs/free5gc-side/renovacao-kakma.md`) |
| D-C | DV-6: `realm="3GPP-bootstrapping-akma@<FQDN>"` no AF e no cliente |
| D-D | `shared/akma/codec`: `^[A-Fa-f0-9]{64}$` ao ler, minúsculas ao escrever, no AUSF, no AAnF, no AF e no cliente |
| D-E | `3GPP-AKMA;<A-KID>` nos perfis PSK (U2) |
| D-F | DV-5: RFC 7616 com SHA-256 |
| D-G | DV-7: o AAnF nunca responde 204; o AF e o NEF tratam o 204 |
| D-H | editorial, nada a implementar |
| D-I | o AAnF recusa `CtxRemove` sem `supi` (400) |
| D-J a D-N | fase (d) do MVP-8 (U24 a U26, U29) |
| D-O | editorial, nada a implementar |
| D-P | o AMF do projeto só registra o *callback* no acesso 3GPP, que o YAML declara com 202 |
| D-Q | A-TID em base64 no A-KID (ADR-0002), configurável no AAnF |
| D-R | a K_AF é a PSK do TLS 1.3 e o PSK do *premaster* do TLS 1.2, sem KDF a mais (U1) |
| D-S | `3GPP-AKMA;<A-KID>` também no TLS 1.2 (U2) |
| D-T | DTLS 1.2 com o procedimento e o identificador Ua\* do TLS 1.2 (U13) |
| D-U | CoAP dentro do DTLS (U14) |
| D-V | códigos CoAP com a causa no diagnóstico (U30) |
| D-W | DV-11: só as suítes de TS 33.210 |
| D-X | o UDM responde ao `trigger-auth` depois do AMF (K4 do MVP-12) |
| D-Y | DV-9: o NRF emite o *token* do AF externo |
| D-Z | o AMF ignora o conteúdo fora de um registro (MVP-12) |
| D-AA, D-AB | DV-13: `iss` e `client_id` com o identificador do *invoker*; `exp` como instante (MVP-9, ponto C9) |
| D-AC | nada no MVP-9: o *token* do CCF usa o escopo do CAPIF; a autorização do GPSI é do MVP-10 |

### Conferência nas versões mais recentes (28/10/2026)

Pedido do responsável (MVP-13, ponto L9): antes de publicar as recomendações, cada item conferido na versão mais recente da norma no arquivo da 3GPP (`https://www.3gpp.org/ftp/Specs/archive/`), baixada em `specs/` fora do controle de versão. As versões de TS 33.535 (V19.0.0), TS 24.109 (V19.0.0) e TS 33.220 (V19.1.0) continuam as mais recentes. TS 29.522, TS 29.535 e TS 29.503 têm a V20.1.0 (Release 20, setembro de 2026), com os YAML `TS29522_AKMA.yaml`, `TS29535_Naanf_AKMA.yaml` v1.3.0-alpha.1 e `TS29503_Nudm_UECM.yaml` v1.5.0-alpha.2.

| Item | Versão conferida | Situação |
|---|---|---|
| D-A | TS 29.522 V20.1.0 | presente: o texto diz `"rid<value>.atid<value>"`, e o exemplo, `rid012.akid019345346@…` |
| D-B | `TS29503_Nudm_UECM.yaml` v1.5.0-alpha.2 | presente: `get:` |
| D-C, D-F | TS 24.109 V19.0.0; TS 33.535 V19.0.0 | presentes (sem versão nova) |
| D-D | `TS29535_Naanf_AKMA.yaml` v1.3.0-alpha.1; `TS29522_AKMA.yaml` da V20.1.0 | presente: `kAkma` e `kaf` só `type: string` |
| D-E, D-J a D-N | TS 33.535 V19.0.0; TS 33.220 V19.1.0 | presentes (sem versão nova) |
| D-G | TS 29.535 V20.1.0 | presente: o 204 e o 403 `K_AKMA_NOT_PRESENT` no mesmo §4.2.2.3.2 |
| D-H | TS 29.522 V20.1.0 | em parte: a remissão a 5.14.5.3.4 e a "table 5.14.5.2-1" foram corrigidas; "Routing Indicatorand" continua. Novo: "retrievel" na descrição de `notifUri` |
| D-I | TS 29.535 V20.1.0 | em parte: "retrieval notification" corrigido; continuam o `CtxRemove` sem `required`, o título `3gpp-akma` do YAML e um "manadatory" |
| D-O | TS 29.503 V20.1.0 | presente: passo 2a de §5.3.2.16.2 e "manadatory" em várias notas |
| D-P | `TS29503_Nudm_UECM.yaml` v1.5.0-alpha.2 | presente |
| D-Q a D-Z | as versões citadas em cada item | presentes |
| D-AA, D-AB | TS 29.222 V20.1.0 e o YAML da mesma versão; TS 33.122 V19.4.0 (sem versão nova) | presentes (conferido em 29/10/2026) |
| D-AC | `TS29522_AKMA.yaml` v1.2.0-alpha.1 (TS 29.522 V20.0.0 e V20.1.0) | corrigido na V20.0.0: `scopes: {}` (conferido em 29/10/2026) |
| TS 23.502 | V20.3.0 | presentes: §5.2.3.2 sem `AuthTrigger`; consumidores de `Nudm_EventExposure_Subscribe` só NEF e SMS-GMSC |

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
| — | `trigger-auth`: quem responde 204; "manadatory" | D-O |
| — | 202 ausente na notificação de reautenticação do acesso não 3GPP | D-P |
| — | Sugestões: lacunas que o projeto teve de decidir | D-Q a D-Z |

---

## D-A — Rótulo `atid` versus `akid` no formato do A-KID

| | |
|---|---|
| **Especificação** | TS 29.522 V19.7.0, Tabela 5.14.5.4.2-1, tipo `AKId` |
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

## D-O — Editoriais no `trigger-auth` de TS 29.503 (prioridade baixa)

| | |
|---|---|
| **Especificação** | TS 29.503 V19.7.0, §5.3.2.16.2 e tabela 6.2.3.13.4.2.2-2 |
| **Grupo responsável** | 3GPP CT4 |

**Descrição.**

1. §5.3.2.16.2 passo 2a: *"On success, the NF Service Consumer responds with "204 No Content""*. Quem responde é o UDM, o produtor; o consumidor (o AAnF, por exemplo) é quem manda o pedido no passo 1. O passo 2b do mesmo procedimento já fala do UDM implicitamente (*"If the user does not exist, HTTP status code "404 Not Found" shall be returned"*).
2. NOTA da tabela 6.2.3.13.4.2.2-2: *"The manadatory HTTP error status code"*, em vez de *"mandatory"*. É o mesmo erro do D-I, em TS 29.535.

**Impacto.** Nenhum na implementação: as tabelas deixam claro quem responde. O passo 2a pode confundir quem lê só o procedimento.

**Proposta.** No passo 2a, *"the UDM responds with "204 No Content""*; na NOTA, *"mandatory"*.

**Decisão do projeto.** Nada a implementar (ponto K12 do MVP-12).

---

## D-P — 202 ausente na notificação de reautenticação do acesso não 3GPP

| | |
|---|---|
| **Especificação** | `TS29503_Nudm_UECM.yaml` v1.4.2 (Anexo A de TS 29.503 V19.7.0) × TS 29.503 §5.3.2.15.2 e tabela 6.2.5.6-3 |
| **Grupo responsável** | 3GPP CT4 |

**Descrição.** O YAML declara o *callback* `reauthNotification` duas vezes: no registro do AMF para o acesso 3GPP (`PUT /{ueId}/registrations/amf-3gpp-access`) e no do acesso não 3GPP (`PUT .../amf-non-3gpp-access`). O primeiro lista a resposta `202` (*"Re-Authentication Notification is accepted and reauthentication pending flag has been set"*); o segundo vai de `204` direto a `307`, sem o `202`. O texto tem uma só notificação: §5.3.2.15.2 passo 2b e a tabela 6.2.5.6-3 definem o `202` para o UE não alcançável, sem distinguir o acesso.

**Impacto.** Pela precedência do YAML (Anexo A.1), um AMF registrado para o acesso não 3GPP não poderia responder `202`, e não teria como dizer ao UDM que marcou a autenticação como pendente. Um UDM gerado do YAML trataria o `202` desse acesso como resposta não declarada.

**Proposta.** Acrescentar o `202`, com a mesma descrição, ao `reauthNotification` do registro para o acesso não 3GPP.

**Decisão do projeto.** Sem efeito no laboratório: o AMF do projeto só registra o `reauthNotifyCallbackUri` no acesso 3GPP (`free5gc/patches/amf/0001`), e o cliente do UDM aceita o `202` de qualquer acesso.

---

## Sugestões: lacunas que o projeto teve de decidir

Acrescentadas em 28/10/2026 (MVP-13, ponto L7), a pedido do responsável ("erros ou sugestões"). Não são contradições do texto, como D-A a D-P. São pontos em que a norma não diz o bastante para duas implementações interoperarem, e o projeto teve de escolher.

## D-Q — Codificação do A-TID no A-KID

| | |
|---|---|
| **Especificações** | TS 33.535 V19.0.0, §6.1 e Anexo A.3; TS 29.522 V20.1.0, Tabela 5.14.5.4.2-1, tipo `AKId` |
| **Grupo responsável** | 3GPP SA3 e CT3 |

**Descrição.** TS 33.535 §6.1: *"A-KID shall be in NAI format as specified in clause 2.2 of IETF RFC 7542 [6], i.e. username@realm. The username part shall include the RID and the A-TID"*. O A-TID é a saída da KDF do Anexo A.3, 256 bits. Nenhuma das normas diz como esses bits viram os caracteres do *username*: hexadecimal, base64, base64url. TS 29.522 dá o formato `"rid<value>.atid<value>"`, também sem a codificação. É a questão C do projeto.

**Impacto.** UE, AUSF, AAnF e AF precisam do mesmo texto do A-KID. Duas codificações diferentes dão A-KIDs diferentes para o mesmo A-TID, e o AF não obtém a K_AF.

**Proposta.** Fixar a codificação, por exemplo a hexadecimal minúscula de 64 caracteres, a mesma dos outros valores de 256 bits nas APIs (D-D), ou base64 com o alfabeto permitido no *username* de RFC 7542.

**Decisão do projeto** (ADR-0002). Base64, configurável no AAnF (`akma.aKid`).

## D-R — Derivação da PSK e do *premaster* a partir da K_AF

| | |
|---|---|
| **Especificações** | TS 24.109 V19.0.0, K.3.2 f) e K.3.3; TS 33.535 V19.0.0, B.1.3 |
| **Grupo responsável** | 3GPP CT1 e SA3 |

**Descrição.** TS 24.109 K.3.2 f): *"the UE and AF shall derive the TLS premaster secret from KAF (AKMA Application Key)"*. K.3.3: *"The UE shall derive the TLS external PSK from the AF specific key KAF"* e *"The AF shall derive the TLS external PSK from the AF specific key KAF"*. Nenhuma das cláusulas, nem TS 33.535 B.1.3, nem TS 33.210, dá a função de derivação: a K_AF pode ser a própria PSK, ou entrar numa KDF.

**Impacto.** UE e AF com derivações diferentes não completam o *handshake*.

**Proposta.** Dizer que a K_AF é usada diretamente: no TLS 1.2, como o PSK de RFC 4279 e RFC 5489 no *premaster secret*; no TLS 1.3, como a PSK externa, com SHA-256. A separação de chaves por protocolo já está no AF_ID (TS 33.220 H.2). Ou definir a KDF no Anexo A de TS 33.535.

**Decisão do projeto** (U1). A K_AF diretamente, nos dois casos.

## D-S — Separador da identidade PSK no TLS 1.2

| | |
|---|---|
| **Especificação** | TS 24.109 V19.0.0, K.3.2 e) |
| **Grupo responsável** | 3GPP CT1 |

**Descrição.** K.3.2 e): o `psk_identity` *"shall contain a prefix "3GPP-AKMA" and the A-KID"*, sem separador. Para o TLS 1.3, K.3.3 1) d) diz *"a prefix indicating the PSK identity name space, i.e. "3GPP-AKMA", a separator character ";" and the A-KID"*. O GBA usa `;` (TS 24.109 §5.3.3.1.2). D-E trata da mesma omissão no TLS 1.3 de TS 33.535.

**Impacto.** `3GPP-AKMA<A-KID>` e `3GPP-AKMA;<A-KID>` são identidades diferentes.

**Proposta.** Acrescentar o separador `";"` a K.3.2 e), como em K.3.3.

**Decisão do projeto** (U2). `3GPP-AKMA;<A-KID>` nos dois TLS.

## D-T — DTLS 1.2 no AKMA

| | |
|---|---|
| **Especificações** | TS 33.535 V19.0.0, C.1.1 e C.2; TS 33.220 V19.1.0, Anexo H.3; TS 33.210 V20.0.0, §6.2.1 |
| **Grupo responsável** | 3GPP SA3 |

**Descrição.** TS 33.210 §6.2.1: *"If DTLS is supported then DTLS 1.2 as specified in RFC 6347 [49] shall be supported and DTLS 1.3 as specified in RFC 9147 [75] should be supported"*. TS 33.535 C.1.1: *"The PSK identity specified in B.1 for TLS is also applicable for DTLS"*, e B.1 inclui o TLS 1.2. Mas C.2 só tem os procedimentos do DTLS 1.3, e TS 33.220 H.3 não dá identificador Ua\* ao DTLS 1.2: `(0x01,0x00,0x00,0x00,0x09)` é o *"Ua security protocol for DTLS according to Annex O"*, que só descreve o DTLS 1.3, e `(0x01,0x00,0x01,yy,zz)` é dito *"for TLS 1.2"*.

**Impacto.** Quem segue TS 33.210 e oferece DTLS 1.2 não tem procedimento nem AF_ID definidos para ele.

**Proposta.** Acrescentar a C.2 o procedimento do DTLS 1.2 e dizer qual identificador Ua\* ele usa; ou dizer que o AKMA só usa DTLS 1.3.

**Decisão do projeto** (U13). DTLS 1.2 com o procedimento do TLS 1.2 e `01 00 01 yy zz` da suíte; o `…09` fica só para o DTLS 1.3, para não repartir a mesma K_AF entre as duas versões (RFC 8446 Apêndice E.7).

## D-U — Protocolo de aplicação dentro do DTLS

| | |
|---|---|
| **Especificações** | TS 33.535 V19.0.0, Anexo C; TS 33.220 V19.1.0, Anexo O.2.2 |
| **Grupo responsável** | 3GPP SA3 |

**Descrição.** O Anexo C de TS 33.535 não diz o que vai dentro do DTLS. O Anexo O de TS 33.220, que ele adapta, diz: *"It does not have any dependency on the application layer protocol."*

**Impacto.** UE e AF precisam do mesmo protocolo de aplicação para trocar qualquer coisa. Com CoAP, o protocolo usual sobre DTLS, vem a suíte obrigatória da RFC 7252 (D-W).

**Proposta.** Nomear o protocolo, ou dizer que ele é escolhido pela aplicação e como o UE e o AF o combinam.

**Decisão do projeto** (U14). CoAP (RFC 7252), esquema `coaps`.

## D-V — Erros do perfil OSCORE

| | |
|---|---|
| **Especificações** | TS 33.535 V19.0.0, Anexo D; TS 33.220 V19.1.0, Anexo P |
| **Grupo responsável** | 3GPP SA3 |

**Descrição.** O Anexo D define o pedido `POST /akma` e a resposta de sucesso (D.3.2, D.3.5), mas nenhuma resposta de erro: A-KID desconhecido, K_AF indisponível, AF não autorizado, `OSC-INP` não suportado. O Anexo P de TS 33.220 também não. O perfil HTTP do Anexo B tem os seus erros.

**Impacto.** Cada AF responde de um jeito, e o UE não sabe se deve refazer a autenticação primária, tentar de novo ou desistir.

**Proposta.** Definir as respostas de erro do recurso `/akma`, com códigos CoAP e as mesmas causas do Anexo B.

**Decisão do projeto** (U30). Códigos CoAP com a causa no *payload* de diagnóstico; os erros do OSCORE pela RFC 8613 §8.2.

## D-W — Suíte obrigatória do CoAP com PSK contra o perfil de TS 33.210

| | |
|---|---|
| **Especificações** | TS 33.535 V19.0.0, C.2.1; TS 33.210 V20.0.0, §6.2; RFC 7252 §9.1.3.1 |
| **Grupo responsável** | 3GPP SA3 |

**Descrição.** TS 33.535 C.2.1: *"The TLS profile specified in TS 33.210 [13] clause 6.2 apply to DTLS 1.3"*. TS 33.210 §6.2.3: *"Only cipher suites with AEAD (e.g. GCM) and PFS (i.e., ECDHE) shall be supported."* RFC 7252 §9.1.3.1, para o CoAP no modo PSK: *"Implementations in these modes MUST support the mandatory-to-implement cipher suite TLS_PSK_WITH_AES_128_CCM_8"*, uma suíte de DTLS 1.2 sem ECDHE.

**Impacto.** Um UE ou AF que use CoAP sobre DTLS não cumpre as duas regras ao mesmo tempo.

**Proposta.** Dizer, junto com D-U, que no Ua\* valem as suítes de TS 33.210 e não a suíte obrigatória da RFC 7252.

**Decisão do projeto** (DV-11). Só as suítes de TS 33.210.

## D-X — Quando o UDM responde ao `trigger-auth`

| | |
|---|---|
| **Especificações** | TS 29.503 V20.1.0, §5.3.2.16.2; TS 33.501 V20.3.0, §6.1.5.2 |
| **Grupo responsável** | 3GPP CT4 |

**Descrição.** TS 33.501 §6.1.5.2 passo 1: *"The UDM may acknowledge the request with an Nudm_UECM_AuthTrigger Response to the NF"*. TS 29.503 §5.3.2.16.2 tem o `204` de sucesso e o `404` de usuário inexistente. Nenhuma das duas diz se o `204` vem antes ou depois da `Nudm_UECM_Re-AuthenticationNotification` ao AMF, nem o que o UDM responde quando o AMF recusa (`403`, `409`) ou marca a autenticação como pendente (`202`).

**Impacto.** O consumidor (o AAnF) não sabe se um `204` quer dizer que a autenticação vai acontecer. A política do AAnF para pedir de novo depende disso.

**Proposta.** Dizer se o UDM responde depois do AMF e, nesse caso, como mapeia as respostas da tabela 6.2.5.6-3 nas do `trigger-auth`, que não tem `409`.

**Decisão do projeto** (K4 do MVP-12). O UDM responde depois do AMF: `204` para `204` e `202`; `403 REAUTHENTICATION_NOT_ALLOWED` para `403` e `409`; `404 CONTEXT_NOT_FOUND` para `404`.

## D-Y — Servidor de autorização do NEF sem CAPIF

| | |
|---|---|
| **Especificações** | TS 29.522 V20.1.0, §6; TS 33.501 V20.3.0, §12.4 |
| **Grupo responsável** | 3GPP CT3 |

**Descrição.** TS 29.522 §6: *"When CAPIF is not supported and/or not used, the access to the NEF Northbound APIs shall be authorized by means of OAuth2 protocol (see IETF RFC 6749 [13]), based on local configuration, using the "Client Credentials" authorization grant as defined in clause 12.4 of 3GPP TS 33.501 [6]. When OAuth2 is used, a client shall obtain a "token" from the authorization server prior to consuming services offered by any NEF Northbound API."* Nem essa cláusula nem TS 33.501 §12.4 dizem qual é o servidor de autorização. Com CAPIF, é o CAPIF Core Function (TS 29.522 §7.2).

**Impacto.** Sem CAPIF, cada operador escolhe: o NRF, um servidor próprio. O AF não sabe onde pedir o *token*.

**Proposta.** Nomear o servidor de autorização quando o CAPIF não é usado, ou dizer que ele é dado por configuração ao AF.

**Decisão do projeto** (DV-9). O NRF, por configuração local.

## D-Z — Alcance do `NAS message container` no `Security Mode Complete`

| | |
|---|---|
| **Especificação** | TS 24.501 V20.1.0, §5.4.2.3 e §4.4.6 |
| **Grupo responsável** | 3GPP CT1 |

**Descrição.** §5.4.2.3: *"If, prior to receiving the SECURITY MODE COMMAND message, the UE without a valid 5G NAS security context had sent a REGISTRATION REQUEST message the UE shall include the entire REGISTRATION REQUEST message in the NAS message container IE of the SECURITY MODE COMPLETE message as described in clause 4.4.6."* A frase não diz que vale só para o procedimento de registro em curso. §4.4.6 a), a proteção da mensagem NAS inicial, fala do *"security mode control procedure"* que segue aquele `REGISTRATION REQUEST`. Do lado da rede, §5.4.2.3 diz que o AMF *"shall complete the ongoing registration procedure by considering the REGISTRATION REQUEST message contained in the NAS message container IE as the message that triggered the procedure"*, e não diz o que fazer quando nenhum registro está em curso.

**Impacto.** Numa autenticação primária pedida pela rede doméstica (TS 33.501 §6.1.5), depois do registro, um UE que leia a frase ao pé da letra põe no `Security Mode Complete` o `REGISTRATION REQUEST` antigo. O UERANSIM faz assim, e o AMF do free5GC o tratava como um registro novo (MVP-12).

**Proposta.** Limitar a regra ao procedimento de controle do modo de segurança iniciado durante o registro em que o `REGISTRATION REQUEST` foi enviado, e dizer que o AMF ignora o conteúdo fora de um procedimento em curso.

**Decisão do projeto.** O AMF ignora o conteúdo durante a autenticação pedida pela rede (`free5gc/patches/amf/0001`); o UE do UERANSIM não muda.

---

## Defeitos do CAPIF

Acrescentados em 29/10/2026, da leitura de TS 29.222, TS 33.122 e TS 29.522 para o MVP-9 (plano, seção 4). São contradições do texto, como D-A a D-P.

## D-AA — `iss` do *token* do CAPIF

| | |
|---|---|
| **Especificações** | TS 29.222 V19.8.0, tabela 8.5.4.2.8-1, e `TS29222_CAPIF_Security_API.yaml`; TS 33.122 V19.4.0, C.2.2 |
| **Grupo responsável** | 3GPP CT3; SA3 para TS 33.122 |

**Descrição.** TS 29.222 exige a *claim* `iss` no *token* do CCF: *"This attribute shall contain the API invoker Identifier as defined in Annex C.2.2 of 3GPP TS 33.122 [16] and in clause 4.1.1 of IETF RFC 7519 [24]."* As duas referências não sustentam isso. TS 33.122 C.2.2 não tem `iss` e põe o identificador do *invoker* em `client_id`, *"REQUIRED"*. A RFC 7519 §4.1.1 define `iss` como *"the principal that issued the JWT"*, que é o CCF.

**Impacto.** Um CCF que siga TS 29.222 e um AEF que siga TS 33.122, ou o contrário, não concordam sobre onde está o *invoker*: o AEF procura `client_id` e não acha, ou lê em `iss` o emissor.

**Proposta.** Alinhar as duas normas: o identificador do *invoker* em `client_id`, como em TS 33.122; e `iss`, se ficar, com o identificador do CCF, como na RFC 7519. Ajustar o YAML.

**Decisão do projeto** (DV-13). As duas *claims*, `iss` e `client_id`, com o identificador do *invoker*.

---

## D-AB — `exp` do *token* do CAPIF como duração

| | |
|---|---|
| **Especificações** | TS 29.222 V19.8.0, tabela 8.5.4.2.8-1, e `TS29222_CAPIF_Security_API.yaml`; TS 33.122 V19.4.0, C.2.2 e C.6 |
| **Grupo responsável** | 3GPP CT3 |

**Descrição.** TS 29.222 tipa a *claim* `exp` como `DurationSec`: *"This attribute shall contain the number of seconds after which the access_token is considered to be expired as defined in clause 4.1.4 of IETF RFC 7519 [24]."* A RFC 7519 §4.1.4, que a própria tabela cita, define `exp` como instante (*NumericDate*): *"the expiration time on or after which the JWT MUST NOT be accepted for processing"*. TS 33.122 C.6 também o trata como instante: *"If the current time is beyond the time specified by the "exp" parameter, the associated token shall no longer be considered valid"*. O `AccessTokenClaims` não tem `iat`, e uma duração sem o instante da emissão não diz quando o *token* vence.

**Impacto.** Um CCF que ponha em `exp` a duração, por exemplo 3600, emite *tokens* que um AEF fiel à RFC 7519 julga vencidos desde 1970.

**Proposta.** Tipar `exp` como instante (*NumericDate*, os segundos desde 1970-01-01 em UTC), com a descrição da RFC 7519. A duração continua em `expires_in` da resposta.

**Decisão do projeto** (DV-13). `exp` como instante.

---

## D-AC — Escopo `nnef-akma:gpsi-access` no YAML da API AKMA do NEF

| | |
|---|---|
| **Especificações** | TS 29.522 V19.7.0, §7.2, e `TS29522_AKMA.yaml` v1.1.0 (Anexo A) |
| **Grupo responsável** | 3GPP CT3 |

**Descrição.** §7.2: *"The NEF northbound APIs do not define any scopes for OAuth2 authorization in the present specification."* O YAML da API AKMA que vem com a V19.7.0 (v1.1.0, cuja `externalDocs` cita a V18.6.0) declara em `oAuth2ClientCredentials` o escopo `nnef-akma:gpsi-access` (*"Return GPSI in the AKMA Application Key information for the UE."*). Nenhuma cláusula do texto define esse escopo.

**Impacto.** Quem gera o servidor ou o cliente pelo YAML pode exigir ou pedir um escopo que o texto diz não existir.

**Situação.** Corrigido na V20.0.0: o `TS29522_AKMA.yaml` da Release 20 (v1.2.0-alpha.1) tem `scopes: {}`. O histórico dessa versão traz o CR 1861, *"Security in OpenAPI definitions"*. Fica no registro, com esta nota, e não entra na lista pública (ponto L9 do MVP-13).

**Decisão do projeto.** Nada no MVP-9: o *token* do CCF usa o escopo do CAPIF. A autorização do GPSI é do MVP-10.

---

## Divergências de enumeração em TS 23.502

Registradas para rastreabilidade; não constituem defeito submetível, por serem listas de contexto em especificação de Stage 2.

| Item | Descrição |
|---|---|
| TS 23.502 V19.8.0 §5.2.3.2 | A lista de operações do serviço `Nudm_UECM` não inclui `AuthTrigger`, introduzida em TS 29.503 pelo CR 1133 (18.3.0) |
| TS 23.502 V19.8.0 §5.2.3.5.2 | A lista de consumidores de `Nudm_EventExposure_Subscribe` inclui apenas NEF e SMS-GMSC; TS 29.503 §5.2.2.1 lista o AAnF entre os consumidores do UDM, incluído pelo CR 1180 (18.4.0) |
| TS 33.535 V19.0.0 §4.3.0 NOTA 1 e §7.4 × TS 23.502 §5.2.3.5.2 | A TS 33.535 aponta `Nudm_EventExposure_Subscribe` (TS 23.502 §5.2.3.5.2) como serviço do UDM usado pelo AKMA; a cláusula apontada não lista o AAnF como consumidor. Inconsistência entre as duas normas (acrescentado em 16/10/2026) |

As duas primeiras linhas foram conferidas também na V19.9.0, cópia disponível no acervo (Adendo E), e as três na V20.3.0, em 28/10/2026.
