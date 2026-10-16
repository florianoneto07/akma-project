# Adendo D — Codificação de material de chave e perfis Ua* Stage 3

**Documento:** complemento ao *Baseline Normativo* v2.0 e aos Adendos A, B e C
**Versão:** D.1
**Data:** 02 de setembro de 2026
**Fontes incorporadas:** 3GPP TS 29.509 V19.6.0, TS 33.222 V19.0.0, TS 24.109 V19.0.0
**Aviso:** este adendo **revê uma recomendação anterior do projeto**. Ver seção D.3.

---

## D.1 Encerramento da LACUNA B

### D.1.1 A convenção da série 29

`[3GPP]` `TS29509_Nausf_UEAuthentication.yaml` v19.6.0 define o material de chave transportado sobre SBI da seguinte forma:

| Tipo | Definição | Comprimento |
|---|---|---|
| `Kseaf` | `type: string`, `pattern: '[A-Fa-f0-9]{64}'` | 256 bits |
| `KnrProSe` | `type: string`, `pattern: '[A-Fa-f0-9]{64}'` | 256 bits |
| `Msk` | `type: string`, `pattern: '[A-Fa-f0-9]{128}'` | 512 bits |
| `ResStar` | `type: string`, `pattern: '[A-Fa-f0-9]{32}'` | 128 bits |
| `HxresStar` | `type: string`, `pattern: '[A-Fa-f0-9]{32}'` | 128 bits |

`[3GPP]` `TS29509_Nausf_SoRProtection.yaml` e `TS29509_Nausf_UPUProtection.yaml` seguem o mesmo padrão para valores de comprimento fixo — `SorMac` e `UpuMac` com `^[A-Fa-f0-9]{32}$`, `CounterSor` e `CounterUpu` com `^[A-Fa-f0-9]{4}$`, `UpuHeader` com `^[A-Fa-f0-9]{2}$` e descrição explícita *"encoded as 2 hexadecimal characters"*.

`[3GPP]` Nas mesmas especificações, quando o conteúdo é um contêiner opaco de comprimento variável, o tipo utilizado é outro: `SorHeader`, `SorTransparentInfo` e `UpuTransparentInfo` referenciam `Bytes` de TS 29.571; `SecuredPacket` e `EapPayload` usam `format: byte`. Ambos correspondem a base64.

**A convenção da série 29 é, portanto, sistemática e discriminante:**

- material criptográfico de comprimento fixo → cadeia **hexadecimal**, com `pattern` explícito;
- contêiner binário opaco de comprimento variável → **base64**, via `Bytes` ou `format: byte`.

### D.1.2 Aplicação ao AKMA

`[INTERPRETAÇÃO]` K_AKMA e K_AF são chaves de comprimento fixo de 256 bits, produzidas pela mesma KDF e com o mesmo comprimento de saída que K_SEAF e K_NR_ProSe, ambas declaradas como `[A-Fa-f0-9]{64}` em TS 29.509. Por precedente direto e pela convenção acima, os atributos `kAkma` (TS 29.535) e `kaf` (TS 29.522) devem ser codificados como **cadeias hexadecimais de 64 caracteres**.

A ausência de `pattern` nesses dois atributos é, à luz desta comparação, um defeito editorial das especificações TS 29.535 e TS 29.522, e não uma indicação de codificação distinta. O argumento por ausência registrado no Adendo C §C.3 — de que `Bytes` estava disponível e não foi usado — fica confirmado por evidência positiva.

**A LACUNA B passa a estar encerrada por precedente.** Permanece, contudo, verdadeiro que nenhuma cláusula normativa fixa a codificação para o AKMA; registrar como observação no trabalho e como candidata a submissão editorial ao 3GPP, junto com as DIVERGÊNCIAS A e D.

`[3GPP]` **Observação de implementação.** Os padrões em `TS29509_Nausf_UEAuthentication.yaml` **não são ancorados** — `'[A-Fa-f0-9]{64}'`, sem `^` e `$`. Uma validação que aplique o padrão literalmente aceita cadeias com prefixo e sufixo arbitrários. Os arquivos de SoR e UPU, por contraste, usam a forma ancorada. A implementação do AKMA deve ancorar seus padrões, independentemente da forma adotada pelo precedente.

---

## D.2 Uma distinção de camada não registrada anteriormente

`[3GPP]` TS 24.109 V19.0.0, Anexo K, cláusula K.2 c) 2): no perfil de autenticação HTTP Digest do AKMA, *"a senha usada nos cálculos do digest deve ser K_AF, com K_AF codificada em Base64 conforme especificado na RFC 4648"*.

Isto **não contradiz** a seção D.1. Trata-se de duas camadas distintas:

| Camada | Objeto | Codificação | Fonte |
|---|---|---|---|
| SBI (JSON, N61/N62/N63) | atributos `kAkma`, `kaf` | hexadecimal, 64 caracteres | Precedente TS 29.509 |
| Ua* (HTTP Digest) | K_AF como senha do digest | base64, RFC 4648 | TS 24.109 K.2 c) 2), normativo |

A implementação deve manter K_AF internamente como 32 octetos e aplicar a codificação apropriada em cada fronteira. Confundir as duas produz falha de autenticação no Ua* sem erro visível no SBI.

---

## D.3 LACUNA C: revisão da recomendação do projeto

### D.3.1 O que as novas fontes mostram

`[3GPP]` TS 33.222 V19.0.0 **não contém nenhuma ocorrência do termo AKMA**. A especificação trata exclusivamente do GBA. A remissão do Anexo H de TS 33.220 a *"TLS 1.3 com chave compartilhada, TS 33.222 cláusula 5.4.0.2"* identifica o protocolo ao qual o identificador Ua* `(0x01,0x00,0x00,0x00,0x07)` é alocado; o perfil correspondente para o AKMA está em TS 24.109, Anexo K.

`[3GPP]` TS 33.222 V19.0.0, §5.4.0.2, declara que *"o formato preciso da PSK identity é especificado em TS 24.109"*. TS 24.109 é, portanto, a fonte Stage 3 para o formato dos identificadores no Ua*.

`[3GPP]` TS 24.109 Anexo K **não especifica a codificação do valor do A-TID**. A LACUNA C não é encerrada por evidência direta.

### D.3.2 O precedente do B-TID

`[3GPP]` TS 33.220 V19.1.0, cláusula 4.5.2, passo 7: *"O valor do B-TID deve ser gerado em formato NAI tomando o valor RAND codificado em base64 (cf. RFC 4648) e o nome do servidor BSF, isto é, `base64encode(RAND)@BSF_servers_domain_name`."* A mesma construção aparece nas cláusulas correspondentes de GBA_U e GBA_Digest.

`[3GPP]` A mesma especificação emprega base64 no TMPI: *"TEMP é a codificação base64 dos 24 octetos mais significativos da saída da KDF."*

### D.3.3 Reversão da recomendação

Nos documentos v2.0 e Adendo C, o projeto adotou **hexadecimal minúsculo** como codificação padrão para o valor do A-TID dentro do A-KID. **Essa recomendação é revista.**

Argumentos que passam a favorecer **base64 conforme RFC 4648**:

1. **Analogia estrutural direta.** O A-KID é um identificador em formato NAI, com parte `username` derivada de material criptográfico e parte `realm` contendo o domínio da rede doméstica. O B-TID tem exatamente a mesma estrutura e função, e usa base64 na parte `username`.
2. **Herança declarada.** TS 24.109 Anexo K constrói os perfis Ua* do AKMA por modificação explícita dos perfis do GBA, com *"o AF do AKMA assumindo o papel do NAF do GBA"*. A herança é assumida pela própria especificação.
3. **Precedente interno ao AKMA.** TS 24.109 K.2 c) 2) codifica K_AF em base64 RFC 4648 no perfil Digest do AKMA, demonstrando que a especificação Stage 3 do Ua* do AKMA recorre a base64 quando precisa representar material binário como texto.
4. **Ausência de precedente contrário.** Nenhuma fonte do acervo emprega hexadecimal em parte `username` de identificador NAI.

Argumentos remanescentes a favor de hexadecimal, e sua avaliação:

- *O alfabeto base64 padrão inclui `+`, `/` e `=`, potencialmente problemáticos em um NAI.* Avaliação: o B-TID do GBA contém esses caracteres — RAND de 16 octetos produz 24 caracteres com um octeto de preenchimento — e é usado em produção há duas décadas. A cláusula 2.2 da RFC 7542 admite esses caracteres na parte `utf8-username`. O argumento não se sustenta.
- *Consistência com a codificação hexadecimal de `kAkma` e `kaf` sobre SBI.* Avaliação: a seção D.2 demonstra que o AKMA já emprega codificações distintas por camada. A consistência entre camadas não é um princípio observado pelas especificações.

`[INTERPRETAÇÃO]` **Recomendação revista:** adotar **base64 conforme RFC 4648** como codificação padrão do valor do A-TID na parte `username` do A-KID. Um A-TID de 32 octetos produz 44 caracteres, dos quais o último é de preenchimento. Manter a parametrização já prevista (`akid.Encoding = base64 | base64url | hex`), agora com `base64` como padrão.

`[LACUNA]` A LACUNA C **permanece formalmente aberta**. A recomendação acima é decisão de projeto fundamentada em analogia, não requisito 3GPP, e deve ser assim declarada no trabalho. A questão deve ser submetida à discussão na *issue* free5GC #1146 e, se possível, ao 3GPP CT1, responsável por TS 24.109.

---

## D.4 Perfis Ua* Stage 3 — material novo

TS 24.109 V19.0.0, Anexo K (normativo), *"TLS with AKMA profiles"*, fornece o Stage 3 do ponto de referência Ua*, que não estava disponível nos documentos anteriores. Este material é diretamente relevante para as fases F5 e F6.

### D.4.1 K.2 — Autenticação do UE por chave compartilhada com AF autenticado por certificado

`[3GPP]`:

- O UE estabelece um túnel TLS com o AF; o AF é autenticado por certificado de chave pública; o UE verifica que o certificado corresponde ao FQDN do AF. Não há autenticação de cliente no TLS.
- O UE indica suporte a Digest baseado em AKMA incluindo o *product token* constante **`"3gpp-akma"`** no cabeçalho `User-Agent`.
- O AF indica a escolha pelo segredo compartilhado AKMA por meio do parâmetro `realm` do cabeçalho `WWW-Authenticate`, que deve conter a cadeia constante **`"3GPP-bootstrapping-akma"`**.
- Havendo escolha entre GBA_Digest e AKMA, o AF deve selecionar AKMA.
- No HTTP Digest (RFC 3310): `username` = **A-KID**; senha = **base64(K_AF)**; `realm` com duas partes delimitadas por `@`, sendo a primeira a cadeia constante **`"3gpp-akma"`** e a segunda o **FQDN do AF**, por exemplo `3gpp-akma@af1.operator.com`.
- UE e AF devem verificar, em cada requisição e resposta, que a segunda parte do atributo `realm` é igual ao FQDN do AF.

`[DIVERGÊNCIA] — DIVERGÊNCIA E.` A cláusula K.2 especifica **dois valores distintos** para o parâmetro `realm` do mesmo intercâmbio Digest: `"3GPP-bootstrapping-akma"` no item b) 2), referente ao desafio `WWW-Authenticate`, e `"3gpp-akma@<FQDN>"` no item c) 3), referente ao cálculo do digest. Em HTTP Digest, o cliente ecoa no cabeçalho `Authorization` o `realm` recebido no desafio, e o valor participa do cálculo de `HA1`; os dois valores devem, portanto, coincidir.

A redação de b) 2) assemelha-se à do perfil GBA correspondente e pode ser resíduo editorial da adaptação. `[INTERPRETAÇÃO]` Seguir o item c) 3) — `"3gpp-akma@<FQDN do AF>"` —, que é o que participa do cálculo e o que a alínea d) manda verificar. Registrar a divergência.

### D.4.2 K.3.2 — TLS 1.2 com PSK

`[3GPP]`:

- `ClientHello` com uma ou mais *ciphersuites* PSK e com a extensão `server_name` contendo o *hostname* do AF.
- `ServerKeyExchange` com `psk_identity_hint` contendo a cadeia estática **`"3GPP-AKMA"`**.
- `ClientKeyExchange` com `psk_identity` contendo o prefixo **`"3GPP-AKMA"`** e o **A-KID**.
- UE e AF derivam o *premaster secret* do TLS a partir de K_AF.

### D.4.3 K.3.3 — TLS 1.3 com PSK

`[3GPP]`:

- O UE indica suporte a PSK pela extensão `psk_key_exchange_modes`, envia o *hostname* do AF na extensão `server_name` e indica também métodos de autenticação distintos de PSK.
- Os parâmetros `psk_identity` dentro do campo `psk_identities` devem conter o prefixo de espaço de nomes **`"3GPP-AKMA"`**, um **caractere separador `";"`** e o **A-KID**. Múltiplos parâmetros são separados por vírgula.
- UE e AF derivam a **PSK externa** do TLS a partir de K_AF.
- O AF responde com `ServerHello` indicando o índice do parâmetro `psk_identity` do AKMA.

`[3GPP]` **Ponto de implementação relevante.** TS 33.535 V19.0.0, Anexo B.1.3, descreve a `psk_identity` apenas como composta do prefixo `"3GPP-AKMA"` e do A-KID, **sem mencionar o caractere separador**. O separador `";"` aparece somente em TS 24.109 K.3.3. Uma implementação que siga apenas TS 33.535 produzirá `psk_identity` em formato incompatível. O formato normativo é:

```
"3GPP-AKMA" ";" A-KID
```

### D.4.4 Ausência de renovação de chave — confirmação Stage 3

`[3GPP]` As cláusulas K.2 e K.3 declaram que as cláusulas 5.3.2.3, 5.3.2.4, 5.3.3.3 e 5.3.3.4 de TS 24.109 — relativas à derivação de chave fresca no GBA — **não são suportadas**, *"pois o AKMA não suporta derivar uma chave fresca do mesmo modo que o GBA"*. As NOTAS 2 de ambas as cláusulas registram que *"o modo como uma chave fresca é derivada para o AKMA fica a cargo da implementação do protocolo Ua*"*.

Isto confirma, no Stage 3, o que TS 33.535 §6.4.3 estabelece no Stage 2: não há procedimento explícito de renovação de K_AF. As duas fontes são coerentes, e a observação da cláusula 5.2 de TS 33.535 — de que uma K_AF derivada de uma K_AKMA inalterada não é uma chave nova — permanece o comportamento normativo.

### D.4.5 O identificador Ua* fixo do TLS 1.3

`[3GPP]` TS 33.222 V19.0.0, §5.4.0.2, NOTA 0: *"É necessário que o UE conheça o identificador de protocolo de segurança Ua para derivar Ks_(ext)_NAF ou Ks_int_NAF antes de enviar a mensagem ClientHello. Por isso, o TLS 1.3 usa um identificador de protocolo de segurança Ua fixo, independente da ciphersuite que só seria determinada em estágio posterior do protocolo TLS."*

Isto explica a estrutura do Anexo H de TS 33.220 registrada no documento v2.0 §2.4: o TLS 1.2 possui identificador dependente da *ciphersuite*, `(0x01,0x00,0x01,yy,zz)`, enquanto o TLS 1.3 possui identificador fixo, `(0x01,0x00,0x00,0x00,0x07)`.

`[MVP]` **Consequência para o projeto.** Um identificador Ua* dependente da *ciphersuite* cria dependência circular na derivação de K_AF: o UE precisa da chave antes de conhecer a *ciphersuite* negociada. A decisão registrada no documento v2.0 §6.4 — adotar um identificador **fixo** da faixa privada, `(0xFF,0x00,0x00,0x00,0x01)` — evita esse problema por construção. A decisão estava correta; agora possui fundamento explícito. A evolução para TLS 1.3 PSK, que também usa identificador fixo, é direta; a evolução para TLS 1.2 exigiria tratar a dependência da *ciphersuite*.

---

## D.5 Estado consolidado

### D.5.1 Questões em aberto

| Nº | Assunto | Estado |
|---|---|---|
| 1 | Obrigatoriedade da subscrição de `ROAMING_STATUS` | Aberta |
| 5 | `AkmaInd`, `RoutingId`, `AANF`, `naanf-akma` no `free5gc/openapi` | Aberta — inspeção do repositório |
| A | Rótulo `atid` versus `akid` em TS 29.522 | Aberta — defeito editorial |
| B | Codificação de `kAkma` e `kaf` | **Encerrada por precedente.** Hexadecimal, 64 caracteres (TS 29.509) |
| C | Codificação do valor do A-TID | Aberta — **recomendação revista para base64 RFC 4648** |
| D | Método HTTP de `trigger-auth` | Aberta — três fontes sustentam POST |
| **E** | Dois valores de `realm` em TS 24.109 K.2 | **Nova, aberta** — seguir K.2 c) 3) |

### D.5.2 Defeitos editoriais candidatos a submissão

| Ref. | Especificação | Descrição |
|---|---|---|
| A | TS 29.522 §5.14.5.4.2 | Texto normativo prescreve rótulo `atid`; exemplo usa `akid` |
| D | `TS29503_Nudm_UECM.yaml` | `trigger-auth` declarada como `get` com `requestBody` obrigatório; corpo da TS 29.503, TS 29.501 §4.6.1.2 e Anexo C.4 indicam POST |
| E | TS 24.109 K.2 | Dois valores distintos de `realm` para o mesmo intercâmbio Digest |
| — | TS 29.522, TS 29.535 | `kaf` e `kAkma` sem `pattern`, ao contrário do precedente de TS 29.509 |
| — | TS 33.535 Anexo B.1.3 | Omite o separador `";"` da `psk_identity` especificado em TS 24.109 K.3.3 |

Recomenda-se consolidar estes cinco itens em um documento único do repositório. O conjunto tem valor próprio como contribuição do trabalho, independentemente da implementação.

### D.5.3 Baseline documental

O acervo cobre agora todas as especificações necessárias ao AKMA, incluindo o Stage 3 do ponto de referência Ua*:

| Camada | Especificações |
|---|---|
| Stage 2 do AKMA | TS 33.535 V19.0.0 |
| Segurança 5G | TS 33.501 V19.7.0 |
| KDF e identificadores Ua* | TS 33.220 V19.1.0 |
| Identificadores | TS 23.003 V19.7.0 |
| Stage 3 do AAnF | TS 29.535 V19.5.0 |
| Stage 3 do NEF e tipos AKMA | TS 29.522 V19.7.0 |
| Stage 3 do UDM | TS 29.503 V19.7.0 |
| Stage 3 do AUSF (precedente) | TS 29.509 V19.6.0 |
| NRF | TS 29.510 V19.7.0 |
| Tipos comuns | TS 29.571 V19.7.0, TS 29.122 V19.7.0 |
| SBA | TS 29.500 V19.7.0, TS 29.501 V19.5.0 |
| **Stage 3 do Ua\*** | **TS 24.109 V19.0.0** |
| Contexto GBA/TLS | TS 33.222 V19.0.0 |
| Procedimentos 5GS | TS 23.502 V19.8.0 |

**Nenhuma lacuna documental bloqueante permanece.** As três questões abertas de serialização — A, C e E — dizem respeito a formatos de identificador e admitem tratamento parametrizado.

---

## D.6 Síntese

1. **A LACUNA B está encerrada.** TS 29.509 declara `Kseaf` e `KnrProSe`, ambas de 256 bits, como `[A-Fa-f0-9]{64}`, e reserva base64 para contêineres opacos. `kAkma` e `kaf` seguem a mesma convenção: hexadecimal, 64 caracteres.
2. **Uma distinção de camada foi identificada.** K_AF é hexadecimal sobre SBI e base64 no perfil Digest do Ua* (TS 24.109 K.2). As duas codificações coexistem por especificação.
3. **A recomendação sobre a LACUNA C é revista de hexadecimal para base64.** O B-TID do GBA, estruturalmente idêntico ao A-KID, usa `base64encode(RAND)@domínio`; TS 24.109 Anexo K constrói o Ua* do AKMA por modificação explícita do GBA; e o próprio AKMA usa base64 no Digest. A questão permanece formalmente aberta e a decisão é do projeto.
4. **O Stage 3 do Ua* está agora disponível** em TS 24.109 Anexo K, com três perfis — Digest, TLS 1.2 PSK e TLS 1.3 PSK — integralmente especificados. Isto habilita as fases F5 e F6 com fundamento normativo, em vez de um protocolo improvisado.
5. **Duas omissões de TS 33.535 foram identificadas**: o separador `";"` da `psk_identity` e a codificação base64 de K_AF no Digest, ambas presentes apenas em TS 24.109.
6. **A decisão do MVP de usar identificador Ua* fixo da faixa privada está fundamentada**: identificadores dependentes de *ciphersuite* criam dependência circular na derivação, conforme a NOTA 0 de TS 33.222 §5.4.0.2.
