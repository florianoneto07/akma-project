# Registro de defeitos e divergências em especificações 3GPP

**Escopo:** inconsistências verificáveis identificadas durante a análise normativa do AKMA Release 19.
**Estado:** aberto para submissão. Nenhum item foi ainda submetido ao 3GPP.

Cada item registra a inconsistência, as fontes que a evidenciam, a decisão adotada pelo projeto e o grupo de trabalho responsável.

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

**Decisão do projeto.** Adotar POST. Item fora do escopo do MVP, pois TS 33.535 §6.4.4 emprega linguagem permissiva.

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

**Impacto.** Em HTTP Digest, o cliente ecoa no cabeçalho `Authorization` o `realm` recebido no desafio, e o valor participa do cálculo de `HA1`. Os dois valores devem coincidir. A redação de b) 2) assemelha-se à do perfil GBA correspondente e aparenta ser resíduo editorial da adaptação.

**Decisão do projeto.** Seguir o item c) 3) — `"3gpp-akma@<FQDN do AF>"` —, que é o valor que participa do cálculo do digest e o que a alínea d) da mesma cláusula manda verificar.

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

## Divergências de enumeração em TS 23.502

Registradas para rastreabilidade; não constituem defeito submetível, por serem listas de contexto em especificação de Stage 2.

| Item | Descrição |
|---|---|
| TS 23.502 V19.8.0 §5.2.3.2 | A lista de operações do serviço `Nudm_UECM` não inclui `AuthTrigger`, introduzida em TS 29.503 pelo CR 1133 (18.3.0) |
| TS 23.502 V19.8.0 §5.2.3.5.2 | A lista de consumidores de `Nudm_EventExposure_Subscribe` inclui apenas NEF e SMS-GMSC; TS 29.503 §5.2.2.1 lista o AAnF entre os consumidores do UDM, incluído pelo CR 1180 (18.4.0) |
