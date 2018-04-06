
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

int verify_certificate(coap_context_t *ctx, coap_dtls_pki_t * setup_data) {
    char* ciphers = "TLSv1.2:TLSv1.0:!PSK";
    coap_openssl_context_t *context = (coap_openssl_context_t *)(ctx->dtls_context);
    if (context->dtls.ctx) {
        if (setup_data->ca_file) {
            SSL_CTX_set_verify(context->dtls.ctx, SSL_VERIFY_PEER | SSL_VERIFY_FAIL_IF_NO_PEER_CERT, NULL);
            if (0 == SSL_CTX_load_verify_locations(context->dtls.ctx, setup_data->ca_file, NULL)) {
                coap_log(LOG_WARNING, "*** verify_certificate: DTLS: %s: Unable to load verify locations\n", setup_data->ca_file);
                return 0;
            }
        }

        if (setup_data->public_cert && setup_data->public_cert[0]) {
            if (0 == SSL_CTX_set_cipher_list(context->dtls.ctx, ciphers)){
                coap_log(LOG_WARNING, "*** verify_certificate: Unable to set ciphers %s \n",  ciphers);
                return 0;
            }
        }
    }

    if (context->tls.ctx) {
        if (setup_data->ca_file) {
            SSL_CTX_set_verify(context->tls.ctx, SSL_VERIFY_PEER | SSL_VERIFY_FAIL_IF_NO_PEER_CERT, NULL);
            if (0 == SSL_CTX_load_verify_locations(context->tls.ctx, setup_data->ca_file, NULL)) {
                coap_log(LOG_WARNING, "*** verify_certificate: TLS: %s: Unable to load verify locations\n", setup_data->ca_file);
                return 0;
            }
        }
        if (setup_data->public_cert && setup_data->public_cert[0]) {
            if (0 == SSL_CTX_set_cipher_list(context->tls.ctx, ciphers)){
                coap_log(LOG_WARNING, "*** verify_certificate: Unable to set ciphers %s \n",  ciphers);
                return 0;
            }
        }
    }
    return 1;
}