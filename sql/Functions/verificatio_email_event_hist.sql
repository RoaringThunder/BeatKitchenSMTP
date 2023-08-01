-- FUNCTION: sd_01.endpoint_log_hist()
-- DROP FUNCTION IF EXISTS sd_01.endpoint_log_hist();
CREATE
OR REPLACE FUNCTION smtp.verification_email_event() RETURNS trigger LANGUAGE 'plpgsql' COST 100 VOLATILE NOT LEAKPROOF AS $ BODY $ DECLARE BEGIN IF (TG_OP = 'DELETE') THEN
INSERT INTO
    smtp.verification_email_event
SELECT
    'D',
    now(),
    OLD.*;

RETURN OLD;

ELSEIF (TG_OP = 'UPDATE') THEN
INSERT INTO
    smtp.verification_email_event
SELECT
    'U',
    now(),
    NEW.*;

RETURN NEW;

ELSEIF (TG_OP = 'INSERT') THEN
INSERT INTO
    smtp.verification_email_event
SELECT
    'I',
    now(),
    NEW.*;

RETURN NEW;

END IF;

RETURN NULL;

END;

$ BODY $;

ALTER FUNCTION smtp.verification_email_event() OWNER TO postgres;