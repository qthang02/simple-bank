ALTER TABLE IF EXISTS "account" DROP CONSTRAINT "account_owner_currency_key";

ALTER TABLE IF EXISTS "account" DROP CONSTRAINT "account_owner_fkey";

DROP TABLE IF EXISTS "users";