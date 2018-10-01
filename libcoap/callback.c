
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
                           str *tok,
                           str *query,
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
                    str *token,
                    str *queryString,
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

    if (session->tls == NULL) {
        return -1;
    }

    ssl = (SSL *)session->tls;
    if (X509_V_OK != SSL_get_verify_result(ssl)) {
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
    if(*key == '\0' && length == 0){
        coap_resource_set_dirty(resource, NULL);
    } else {
        str *query = coap_new_string(length);
        query->s = key;
        query->length = length;
        coap_resource_set_dirty(resource, query);
    }
}