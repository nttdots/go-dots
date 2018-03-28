#pragma once

#include <coap/coap.h>

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

void pong_handler(struct coap_context_t *,
                      coap_session_t *,
                      coap_pdu_t *,
                      const coap_tid_t);

void set_server_common_name (coap_session_t*,
                             const char* );

int coap_dtls_get_peer_common_name1(coap_session_t *,
                                    char *,
                                    size_t);