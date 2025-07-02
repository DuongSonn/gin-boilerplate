-- Create "users" table
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL,
  "created_at" integer NOT NULL,
  "updated_at" integer NOT NULL,
  "deleted_at" timestamptz NULL,
  "email" character varying(100) NULL,
  "phone_number" character varying(20) NULL,
  "password" text NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  PRIMARY KEY ("id")
);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "public"."users" ("deleted_at");
-- Create "oauth" table
CREATE TABLE "public"."oauth" (
  "id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "ip" text NOT NULL,
  "platform" character varying(10) NOT NULL,
  "token" text NOT NULL,
  "status" text NOT NULL,
  "expire_at" integer NOT NULL,
  "created_at" bigint NULL,
  "updated_at" bigint NULL,
  "login_at" integer NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_o_auths" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Drop enum type "status"
DROP TYPE "public"."status";
