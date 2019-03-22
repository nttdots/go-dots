
#include "callback.h"

extern void export_response_handler(coap_context_t *ctx,
                            coap_session_t *sess,
                            coap_pdu_t *sent,
                            coap_pdu_t *received,
                            coap_tid_t id);

extern void export_method_handler(coap_context_t *ctx,
                           coap_resource_t *rsrc,
                           coap_session_t *sess,
                           coap_pdu_t *req,
                           coap_string_t *tok,
                           coap_string_t *query,
                           coap_pdu_t *resp);

extern void export_nack_handler(coap_context_t *ctx,
                    coap_session_t *sess,
                    coap_pdu_t *sent,
                    coap_nack_reason_t reason,
                    coap_tid_t id);

extern void export_ping_handler(coap_context_t *ctx,
                    coap_session_t *sess,
                    coap_pdu_t *sent,
                    coap_tid_t id);

extern void export_event_handler(coap_context_t *ctx,
                    coap_event_t event,
                    coap_session_t *sess);

void response_handler(coap_context_t *context,
                      coap_session_t *session,
                      coap_pdu_t *sent,
                      coap_pdu_t *received,
                      const coap_tid_t id) {

    export_response_handler(context, session, sent, received, id);
}

void method_handler(coap_context_t *context,
                    coap_resource_t *resource,
                    coap_session_t *session,
                    coap_pdu_t *request,
                    coap_string_t *token,
                    coap_string_t *queryString,
                    coap_pdu_t *response) {

    export_method_handler(context, resource, session, request, token, queryString, response);
}

void nack_handler(coap_context_t *context,
                    coap_session_t *session,
                    coap_pdu_t *sent,
                    coap_nack_reason_t reason,
                    const coap_tid_t id){

    export_nack_handler(context, session, sent, reason, id);
}

void ping_handler(coap_context_t *context,
                      coap_session_t *session,
                      coap_pdu_t *sent,
                      const coap_tid_t id) {

    export_ping_handler(context, session, sent, id);
}

void event_handler(coap_context_t *context,
                      coap_event_t event,
                      void *data) {

    export_event_handler(context, event, (coap_session_t *)data);
}

int coap_dtls_get_peer_common_name(coap_session_t *session,
                                    char *buf,
                                    size_t buf_len){
    SSL *ssl;
    X509 *cert;
    X509_NAME *name;
    int cn_len;

    coap_openssl_context_t *ctx = (coap_openssl_context_t *)session->context->dtls_context;
    coap_dtls_pki_t *setup_data = &ctx->setup_data;

    if (session->tls == NULL) {
        return -1;
    }
    ssl = (SSL *)session->tls;

    long verify_result = SSL_get_verify_result(ssl);
    switch (verify_result) {
    case X509_V_ERR_CERT_NOT_YET_VALID:
    case X509_V_ERR_CERT_HAS_EXPIRED:
        if (setup_data->allow_expired_certs)
            verify_result = X509_V_OK;
        break;
    case X509_V_ERR_SELF_SIGNED_CERT_IN_CHAIN:
        if (setup_data->allow_self_signed)
            verify_result = X509_V_OK;
        break;
    case X509_V_ERR_UNABLE_TO_GET_CRL:
        if (setup_data->allow_no_crl)
            verify_result = X509_V_OK;
        break;
    case X509_V_ERR_CRL_NOT_YET_VALID:
    case X509_V_ERR_CRL_HAS_EXPIRED:
        if (setup_data->allow_expired_crl)
            verify_result = X509_V_OK;
        break;
    default:
        break;
    }
    if (X509_V_OK != verify_result) {
        coap_log(LOG_WARNING, "    %s\n", X509_verify_cert_error_string(verify_result));
        return -1;
    }
    cert = SSL_get_peer_certificate(ssl);
    if (cert == NULL) {
        return -1;
    }

    name = X509_get_subject_name(cert);
    cn_len = X509_NAME_get_text_by_NID(name, NID_commonName, NULL, 0);
    if (cn_len < 0) {
        return -1;
    }
    if (buf_len < (size_t)cn_len + 1) {
        return -1;
    }
    return X509_NAME_get_text_by_NID(name, NID_commonName, buf, buf_len);

}

void coap_set_dirty(coap_resource_t *resource, char *key, int length) {
    if (*key == '\0' && length == 0) {
        coap_resource_notify_observers(resource, NULL);
    } else {
        coap_string_t *query = coap_new_string(length);
        query->s = (uint8_t*)key;
        query->length = (size_t)length;
        coap_resource_notify_observers(resource, query);
    }
}

int coap_check_subscribers(coap_resource_t *resource) {
    return !(resource->subscribers == NULL);
}

int coap_check_dirty(coap_resource_t *resource) {
    return resource->dirty;
}

// Get token from subcribers
char* coap_get_token_subscribers(coap_resource_t *resource) {
    coap_subscription_t *subscriber = resource->subscribers;
    if (subscriber != NULL) {
        return subscriber->token;
    }
    return (char*)0;
}

// Get size of block 2 from subcribers
int coap_get_size_block2_subscribers(coap_resource_t *resource) {
    coap_subscription_t *subscriber = resource->subscribers;
    if (subscriber != NULL) {
        coap_block_t block2 = subscriber->block2;
        return block2.szx;
    }
    return 0;
}

// create coap_block_t
coap_block_t coap_create_block(unsigned int num, unsigned int m, unsigned int size) {
   coap_block_t block = { num, m, size };
   return block;
}