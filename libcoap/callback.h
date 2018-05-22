#pragma once

#include <coap/coap.h>
#include <openssl/ssl.h>

void response_handler(struct coap_context_t *,
                      coap_session_t *,
                      coap_pdu_t *,
                      coap_pdu_t *,
                      const coap_tid_t);

void method_handler(coap_context_t *,
                    coap_resource_t *,
                    coap_session_t *,
                    coap_pdu_t *,
                    str *,
                    str *,
                    coap_pdu_t *);

void nack_handler(struct coap_context_t *,
                    coap_session_t *,
                    coap_pdu_t *,
                    coap_nack_reason_t,
                    const coap_tid_t);
/**
 * Get peer common name (from certificate issuer names)
 * @param session   The CoAP session
 * @param buf       The return value
 * @param buflen    The max length of the return value 
 * @return          Return the index of the next matching entry -1 if not found
 */
int coap_dtls_get_peer_common_name(coap_session_t *session,
                                    char *buf,
                                    size_t buf_len);

/**
 * Verify certificate data and set list of available ciphers for context
 * @param ctx     The CoAP context
 * @param setup_data  certificate data
 * @return            Return 1 for success, 0 for failure
 */
int verify_certificate(coap_context_t *ctx, coap_dtls_pki_t *setup_data);


typedef struct coap_dtls_context_t {
  SSL_CTX *ctx;
  SSL *ssl;	/* OpenSSL object for listening to connection requests */
  HMAC_CTX *cookie_hmac;
  BIO_METHOD *meth;
  BIO_ADDR *bio_addr;
} coap_dtls_context_t;

typedef struct coap_tls_context_t {
  SSL_CTX *ctx;
  BIO_METHOD *meth;
} coap_tls_context_t;

typedef struct coap_openssl_context_t {
  coap_dtls_context_t dtls;
  coap_tls_context_t tls;
  int psk_pki_enabled;
} coap_openssl_context_t;

void coap_set_dirty(coap_resource_t *resource, char *query, int length);
coap_resource_t *coap_get_resource(coap_context_t *context, char *key, int length);