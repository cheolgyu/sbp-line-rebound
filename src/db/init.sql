
CREATE SCHEMA "hist";
DROP TABLE IF  EXISTS hist.rebound  CASCADE;
CREATE TABLE IF NOT EXISTS hist.rebound  (
	"code_id" integer not null REFERENCES "meta"."code"(id),
	"price_type" integer not null REFERENCES "meta"."config"(id),
   	"x1" numeric(8, 0) not null REFERENCES "meta"."opening"(dt) CHECK (x1 < x2),
    "y1" numeric(20, 2),
    "x2" numeric(8, 0) not null REFERENCES "meta"."opening"(dt) CHECK (x1 < x2),
    "y2" numeric(20, 2),
    "x_tick" numeric(20, 0),
    "y_minus" numeric(20, 2),
    "y_percent" numeric(20, 2)
) PARTITION  BY  RANGE (x1); 
CREATE INDEX ON hist.rebound (code_id);
ALTER TABLE hist.rebound ADD CONSTRAINT hist_rebound_code_id_price_type_x1_key PRIMARY KEY (code_id, price_type,x1);


-- 1956년 부터 2021+10 년까지

CREATE OR REPLACE FUNCTION hist.create_table_rebound()
 RETURNS void
 LANGUAGE plpgsql
AS $function$
declare
	start_y integer  := 1956;
	end_y integer  := 2031;
BEGIN
	 LOOP
        
        EXIT WHEN start_y = end_y;
        
		   execute format('CREATE TABLE IF NOT EXISTS hist.rebound_%s PARTITION OF hist.rebound FOR VALUES FROM (%s0101) TO (%s1231); ', start_y , start_y , start_y );
		SELECT start_y+1 INTO start_y;
    END LOOP;

END;
$function$
;
select * from hist.create_table_rebound();




--------------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------------


DROP TABLE IF  EXISTS public.tb_rebound;
CREATE TABLE IF NOT EXISTS public.tb_rebound (
    "code_id" integer PRIMARY KEY REFERENCES "meta"."code"(id),
    "cp_x1" numeric(8, 0),
    "cp_y1" numeric(20, 0),
    "cp_x2" numeric(8, 0),
    "cp_y2" numeric(20, 0),
    "cp_x_tick" numeric(20, 0),
    "cp_y_minus" numeric(20, 2),
    "cp_y_percent" numeric(20, 2),

    "op_x1" numeric(8, 0),
    "op_y1" numeric(20, 0),
    "op_x2" numeric(8, 0),
    "op_y2" numeric(20, 0),
    "op_x_tick" numeric(20, 0),
    "op_y_minus" numeric(20, 2),
    "op_y_percent" numeric(20, 2),

    "lp_x1" numeric(8, 0),
    "lp_y1" numeric(20, 0),
    "lp_x2" numeric(8, 0),
    "lp_y2" numeric(20, 0),
    "lp_x_tick" numeric(20, 0),
    "lp_y_minus" numeric(20, 2),
    "lp_y_percent" numeric(20, 2),

    "hp_x1" numeric(8, 0),
    "hp_y1" numeric(20, 0),
    "hp_x2" numeric(8, 0),
    "hp_y2" numeric(20, 0),
    "hp_x_tick" numeric(20, 0),
    "hp_y_minus" numeric(20, 2),
    "hp_y_percent" numeric(20, 2)
);