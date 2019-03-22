#pragma once

#include <coap2/coap.h>
#include <openssl/ssl.h>
#include <openssl/err.h>

void response_handler(struct coap_context_t *,
                      coap_session_t *,
                      coap_pdu_t *,
                      coap_pdu_t *,
                      const coap_tid_t);

void method_handler(coap_context_t *,
                    coap_resource_t *,
                    coap_session_t *,
                    coap_pdu_t *,
                    coap_string_t *,
                    coap_string_t *,
                    coap_pdu_t *);

void nack_handler(struct coap_context_t *,
                    coap_session_t *,
                    coap_pdu_t *,
                    coap_nack_reason_t,
                    const coap_tid_t);

void ping_handler(struct coap_context_t *,
                    coap_session_t *,
                    coap_pdu_t *,
                    const coap_tid_t);

void event_handler(coap_context_t *context,
                    coap_event_t event,
                    void *data);
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

typedef struct sni_entry {
  char *sni;
  SSL_CTX *ctx;
} sni_entry;

typedef struct coap_openssl_context_t {
  coap_dtls_context_t dtls;
  coap_tls_context_t tls;
  coap_dtls_pki_t setup_data;
  int psk_pki_enabled;
  size_t sni_count;
  sni_entry *sni_entry_list;
} coap_openssl_context_t;

void coap_set_dirty(coap_resource_t *resource, char *query, int length);

int coap_check_subscribers(coap_resource_t *resource);
int coap_check_dirty(coap_resource_t *resource);
char* coap_get_token_subscribers(coap_resource_t *resource);
int coap_get_size_block2_subscribers(coap_resource_t *resource);
coap_block_t coap_create_block(unsigned int num, unsigned int m, unsigned int size);