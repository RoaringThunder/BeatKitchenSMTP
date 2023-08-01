-- Table: smtp.verification_email_event_hist
-- DROP TABLE IF EXISTS smtp.verification_email_event_hist;
CREATE TABLE IF NOT EXISTS smtp.verification_email_event_hist (
    operation character(1) COLLATE pg_catalog."default",
    stamp timestamp with time zone,
    id bigint,
    error_msg text COLLATE pg_catalog."default",
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    recipient text COLLATE pg_catalog."default",
    status text COLLATE pg_catalog."default"
) TABLESPACE pg_default;

ALTER TABLE
    IF EXISTS smtp.verification_email_event_hist OWNER to postgres;

GRANT ALL ON TABLE smtp.verification_email_event_hist TO postgres;