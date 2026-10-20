# Lado AF e Ua*

Tudo o que é feito **fora** do núcleo para usar AKMA: o AF de teste, o protocolo Ua\* simplificado, o lado UE (ME no UERANSIM e cliente Ua\*) e os roteiros que exercitam o fluxo.

| Documento | MVP | Estado |
|---|---|---|
| `ua-star.md` — especificação do Ua\* simplificado | MVP-3 | aprovada e implementada (18/10/2026) |
| `af.md` — o AF de teste: arquivos, configuração, log | MVP-3 | escrito (18/10/2026) |
| `teste-af.md` — testes, integração com o AAnF e critérios do MVP-3 | MVP-3 | escrito (18/10/2026) |
| `ue-ueransim.md` — K_AKMA, A-KID e K_AF no ME do UERANSIM; `nr-cli`; autoteste; *namespace* do UE; log | MVP-3a, MVP-4 | escrito (18/10/2026) |
| `ua-client.md` — cliente Ua\* do UE, com a K_AF do ME e pelo plano de usuário | MVP-3, MVP-4 | escrito (18/10/2026) |
| `teste-ue.md` — autoteste, ponta a ponta, nova autenticação e critérios do MVP-4 | MVP-4 | escrito (18/10/2026) |
| `e2e.md` — roteiro ponta a ponta, testes negativos e resultado | MVP-5, MVP-6 | escrito (18/10/2026); grupo do NEF e R2n em 19/10/2026 |
| `teste-nef.md` — AF via NEF, acesso anônimo e UERANSIM em outra máquina | MVP-6 | escrito (19/10/2026) |
| `teste-instalacao-github.md` — instalação numa VM limpa só com o repositório público; critério 6 do MVP-7 | MVP-7 | escrito (20/10/2026) |

O que é comum aos dois lados — KDF, formatos de identificador, codificação, vetores — fica em `docs/shared/`.
