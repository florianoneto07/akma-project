/*
 * Copyright 2026 Floriano Neto.
 * Licensed under the GNU General Public License v3.0 or later; see COPYING.
 *
 * Callbacks of wolfSSL that call back into Go. Every WOLFSSL object carries
 * the cgo.Handle of its Conn as I/O context and as PSK callback context.
 */
#ifndef WOLFTLS_SHIM_H
#define WOLFTLS_SHIM_H

#include <stdint.h>
#include <stdlib.h>
#include <wolfssl/options.h>
#include <wolfssl/ssl.h>

/* Configures a context made by the Go side: callbacks, PSK, tickets. */
void wolftls_setup_server(WOLFSSL_CTX *ctx, int tls13, int tickets);
void wolftls_setup_client(WOLFSSL_CTX *ctx, int tls13);

/* Attaches the cgo.Handle h of a Conn to ssl. */
void wolftls_attach(WOLFSSL *ssl, uintptr_t h);

/* Returns the reason text of a wolfSSL error code into buf. */
void wolftls_error_string(int err, char *buf, int len);

#endif
