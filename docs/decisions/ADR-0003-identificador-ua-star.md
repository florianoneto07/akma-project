# ADR-0003 — Identificador de protocolo de segurança Ua* do MVP

**Estado:** Aceito — com a aprovação do plano de desenvolvimento, em 16/10/2026
**Data:** 02 de setembro de 2026

---

## 1. Questão

O `AF_ID` usado como P0 na derivação de K_AF é composto pelo FQDN do AF concatenado ao identificador de protocolo de segurança Ua* (TS 33.535 Anexo A.4). O MVP não implementará integralmente nenhum dos perfis normativos do Anexo K de TS 24.109. Que identificador usar?

---

## 2. Decisão

Adotar o identificador de faixa privada **`(0xFF, 0x00, 0x00, 0x00, 0x01)`** para o protocolo Ua* simplificado do MVP.

Forma de transporte do `afId`: `"<FQDN do AF>.FF00000001"` *(errata 16/10/2026: maiúsculas na escrita, como no exemplo de TS 29.522 — decisão P9 do plano; a leitura aceita as duas caixas)*.
Forma de P0 da KDF: `UTF-8_NFKC(FQDN) ‖ 0xFF 0x00 0x00 0x00 0x01`.

---

## 3. Fundamentação

### 3.1 A faixa privada é normativa

TS 33.220 V19.1.0, Anexo H.1: o identificador é uma cadeia de cinco octetos; o primeiro identifica a organização, os quatro restantes o protocolo.

Anexo H.2: o valor `0xFF` no primeiro octeto designa a **faixa privada**. A NOTA 2 registra que identificadores com `0xFF` no primeiro octeto podem ser usados para definir protocolos locais ou experimentais **sem necessidade de registro**, com a ressalva de que uma falha de segurança em um protocolo sobre Ua* pode ser explorada contra outro.

O uso é, portanto, expressamente previsto. **Não constitui desvio da norma.**

### 3.2 Por que não o valor padrão

O identificador `(0x00,0x00,0x00,0x00,0x00)` também é normativo, mas a NOTA 1 do Anexo H.2 registra que protocolos que o utilizam não podem ser separados entre si. Isto elimina a separação criptográfica de chaves entre protocolos, propriedade que o projeto pretende preservar e eventualmente verificar formalmente.

### 3.3 Por que um identificador fixo

TS 33.222 V19.0.0, §5.4.0.2, NOTA 0: é necessário que o UE conheça o identificador Ua* para derivar a chave **antes** de enviar o `ClientHello`; por isso o TLS 1.3 usa identificador fixo, independente da *ciphersuite*, que só seria determinada em estágio posterior do protocolo.

Isto explica a estrutura do Anexo H de TS 33.220: o TLS 1.2 com chave compartilhada usa `(0x01,0x00,0x01,yy,zz)`, dependente da *ciphersuite*; o TLS 1.3 usa `(0x01,0x00,0x00,0x00,0x07)`, fixo.

Um identificador dependente da *ciphersuite* cria dependência circular na derivação de K_AF. A escolha de um identificador fixo evita o problema por construção.

---

### 3.4 Por que não `(0x01,0x00,0x00,0x00,0x02)`

*Acrescentado em 16/10/2026.*

`[3GPP]` TS 33.220 V19.1.0, H.3, registra `(0x01,0x00,0x00,0x00,0x02)` como *"Ua security protocol HTTP digest authentication according to TS 24.109, unless HTTP digest authentication is used in the context of another Ua security protocol, which is already covered elsewhere in this Annex"*. O Ua\* simplificado do MVP tem a forma de um HTTP Digest: o de TS 24.109 §5.2, com os parâmetros AKMA de TS 24.109 K.2 e TS 33.535 B.1.2.2.

`[3GPP]` Para o AKMA, TS 24.109 Anexo K especifica o Digest apenas dentro de um túnel TLS com autenticação do AF por certificado (K.2), adaptação de TS 33.222 §5.3; em H.3, esse protocolo do GBA tem outro identificador, `(0x01,0x00,0x01,yy,zz)`.

`[INTERPRETAÇÃO]` O identificador `…02` não é usado porque:

1. o Ua\* do MVP não tem TLS, e o Anexo K não especifica Digest AKMA sem TLS — usar um identificador registrado afirmaria um perfil que o projeto não implementa;
2. a faixa privada mantém as chaves do protocolo experimental separadas das de qualquer protocolo registrado, com a ressalva de H.2 NOTA 2.

A decisão da seção 2 não muda.

## 4. Consequências

- A migração para o perfil TLS 1.3 PSK do Anexo K.3.3 de TS 24.109 exige apenas a troca do identificador por `(0x01,0x00,0x00,0x00,0x07)` e a adoção do formato de `psk_identity` `"3GPP-AKMA" ";" A-KID`. Nenhuma alteração arquitetural.
- A migração para TLS 1.2 exigiria tratar a dependência da *ciphersuite* na derivação, e não é recomendada.
- O Ua\* do MVP segue RFC 7616 (desvio DV-5) e usa `realm = "3GPP-bootstrapping-akma@<FQDN>"`, conforme TS 33.535 B.1.2.2 (desvio DV-6); ver `docs/conformidade/desvios.md`.
- O identificador escolhido é uma **decisão do projeto**, e o Ua* do MVP é um protocolo experimental. Ambos devem ser declarados como tais no trabalho. ~~As cláusulas 6.1 a 6.3 de TS 33.535 são cumpridas integralmente; apenas o perfil Ua* fica fora do escopo do MVP.~~ *(errata 16/10/2026)* As etapas obrigatórias de TS 33.535 §6.1 e §6.2.1 são cumpridas, com os desvios registrados em `docs/conformidade/desvios.md` (em particular DV-2, subscrição de `ROAMING_STATUS` desabilitada); §6.3 (via NEF) fica para o MVP-6+; o perfil Ua\* normativo fica fora do escopo do MVP.

---

## 5. Referências

- TS 33.220 V19.1.0, Anexo H.1, H.2 (NOTA 1 e NOTA 2), H.3
- TS 33.222 V19.0.0, §5.4.0.2, NOTA 0
- TS 33.535 V19.0.0, Anexos A.4, B.1.2.2 e B.1.3
- TS 24.109 V19.0.0, §5.2 e Anexo K
