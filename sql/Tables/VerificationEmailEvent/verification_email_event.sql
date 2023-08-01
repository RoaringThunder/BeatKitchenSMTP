-- Table: smtp.verification_email_event
-- DROP TABLE IF EXISTS smtp.verification_email_event;
CREATE TABLE IF NOT EXISTS smtp.verification_email_event (
    id bigint NOT NULL DEFAULT nextval('smtp.verification_email_event_id_seq' :: regclass),
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    recipient text COLLATE pg_catalog."default",
    status text COLLATE pg_catalog."default" DEFAULT 'SENT' :: text,
    error_msg text COLLATE pg_catalog."default",
    CONSTRAINT verification_email_event_pkey PRIMARY KEY (id)
) TABLESPACE pg_default;

ALTER TABLE
    IF EXISTS smtp.verification_email_event OWNER to postgres;

-- Index: idx_smtp_verification_email_event_deleted_at
-- DROP INDEX IF EXISTS smtp.idx_smtp_verification_email_event_deleted_at;
CREATE INDEX IF NOT EXISTS idx_smtp_verification_email_event_deleted_at ON smtp.verification_email_event USING btree (deleted_at ASC NULLS LAST) TABLESPACE pg_default;

CREATE TRIGGER tr_smtp_verification_email_event_hist
AFTER
INSERT
    OR DELETE
    OR
UPDATE
    ON smtp.verification_email_event FOR EACH ROW EXECUTE FUNCTION smtp.verification_email_event();