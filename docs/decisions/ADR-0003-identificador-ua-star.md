# ADR-0003 — Identificador de protocolo de segurança Ua* do MVP

**Estado:** Proposto
**Data:** 02 de setembro de 2026

---

## 1. Questão

O `AF_ID` usado como P0 na derivação de K_AF é composto pelo FQDN do AF concatenado ao identificador de protocolo de segurança Ua* (TS 33.535 Anexo A.4). O MVP não implementará integralmente nenhum dos perfis normativos do Anexo K de TS 24.109. Que identificador usar?

---

## 2. Decisão

Adotar o identificador de faixa privada **`(0xFF, 0x00, 0x00, 0x00, 0x01)`** para o protocolo Ua* simplificado do MVP.

Forma de transporte do `afId`: `"<FQDN do AF>.ff00000001"`.
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

## 4. Consequências

- A migração para o perfil TLS 1.3 PSK do Anexo K.3.3 de TS 24.109 exige apenas a troca do identificador por `(0x01,0x00,0x00,0x00,0x07)` e a adoção do formato de `psk_identity` `"3GPP-AKMA" ";" A-KID`. Nenhuma alteração arquitetural.
- A migração para TLS 1.2 exigiria tratar a dependência da *ciphersuite* na derivação, e não é recomendada.
- O identificador escolhido é uma **decisão do projeto**, e o Ua* do MVP é um protocolo experimental. Ambos devem ser declarados como tais no trabalho. As cláusulas 6.1 a 6.3 de TS 33.535 são cumpridas integralmente; apenas o perfil Ua* fica fora do escopo do MVP.

---

## 5. Referências

- TS 33.220 V19.1.0, Anexo H.1, H.2 (NOTA 1 e NOTA 2), H.3
- TS 33.222 V19.0.0, §5.4.0.2, NOTA 0
- TS 33.535 V19.0.0, Anexos A.4 e B.1.3
- TS 24.109 V19.0.0, Anexo K
