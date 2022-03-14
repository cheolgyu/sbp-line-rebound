/*
    두점으로 y=mx+b 구하기
-- example.line_point

select y3 ,qu.price  --,ROW_NUMBER() OVER( order BY qu.price  ) AS ROW_NUM
from
(
	select 
		((m*(x2+1))+b)::integer  as y3 
	from (
		select x2  --2-x1 as run, y2-y1 as rise
			,round((y2-y1 )/( x2-x1)::numeric ,2 )as m 
			,round((y2-y1 )/( x2-x1)::numeric ,2 )*x2 *(-1) + y2 as b 
		from  (
				select count(*) as x1 from hist.price_stock where code = '000020'   and p_date <=20210621
			)x1 ,(
				select count(*) as x2 from hist.price_stock where code = '000020'   and p_date <=20210629
			)x2,
			(
				select count(*) as y1 from utils.quote_unit where market = 1   and price <= 14350
			)y1 ,(
				select count(*) as y2 from utils.quote_unit where market = 1   and price <= 15600
			)y2

	)t
)t , utils.quote_unit qu
where qu.market = 1  offset 6115  limit 1
;
*/
----------------------------------------
----------------------------------------
----------------------------------------
----------------------------------------
----------------------------------------
----------------------------------------
 -- 53분걸림 ?? 
 -- ==> CREATE INDEX ON hist.rebound (code_id);
 -- ==> 317 msec 

--  select 1 from project.func_lines();
-- Successfully run. Total query runtime: 5 min 41 secs.
-- Successfully run. Total query runtime: 7 min 49 secs.

----------------------------------------
----------------------------------------
----------------------------------------
----------------------------------------
----------------------------------------
CREATE or replace FUNCTION project.func_lines()
 RETURNS void AS $$
declare
    row record;
	r_p record;
	x3 numeric;
	y3 numeric;

	count integer;
	r2 record;
	y3_price integer;
	arr_g_type  TEXT[]  :=  array['open','close','low','high'];
	i_price_type integer;   
	cnt integer := 0;   
	log_start_dt timestamp;
	log_end_sec integer;
	log_txt TEXT;
BEGIN
 -- code get
	select array_agg(id) as id
	 	,array_agg(code) as code
		,array_agg(name) as name 
	from meta.config where upper_code='price_type'  INTO r_p;
	
	FOR row in 
		select 
			t.*,
			case when t.market_name ='KOSPI' THEN 1 ELSE 2 END AS MARKET_NUM 
		from(
			select  
				code_id ,code, name, market_type
				,(select  name from meta.config where upper_code='market_type' and id = market_type) as market_name
			from only public.company pc 
			where stop is not true
				and code_type = (select  id from meta.config where upper_code='code_type' and code ='stock')
			order by code asc
		) t
	
	loop
		FOREACH i_price_type IN ARRAY r_p.id
		LOOP
			
			EXECUTE format(' SELECT * FROM  project.func_lines_get_last_point(%L ,%L) ', row.code_id, i_price_type ) INTO r2;
			
			select now() into log_start_dt;
			
			

			IF r2.x1 > 0 then

				select cnt +1 into cnt;
				EXECUTE format(' SELECT * FROM  project.func_lines_get_next_point(%L ,%L    ,%L ,%L ,%L ,%L ) ', row.code_id ,row.MARKET_NUM   ,
					r2.x1 ,r2.y1 ,r2.x2 ,r2.y2 ) INTO x3 ,y3;

				IF y3 > 0 then

					EXECUTE format(' SELECT * FROM  project.func_lines_convert_y3(%L ,%L ) ', row.MARKET_NUM ,y3::integer ) INTO y3_price;
					EXECUTE format(' SELECT * FROM  project.func_lines_update_tb_daily_line( %L,%L,%L,  %L ,%L ,%L ,%L ,%L ,%L) '
						, row.code_id , row.code ,i_price_type 
						,r2.x1::numeric  ,r2.y1::numeric  ,r2.x2::numeric  ,r2.y2::numeric  ,x3::numeric ,y3_price::numeric ) ;

				END IF;

			END IF;

			EXECUTE format(' select  CAST(EXTRACT(MINUTE FROM (%L - now()) ) AS INTEGER) ', log_start_dt ) INTO log_end_sec;
			IF log_end_sec < 0 then
				-- 9748
				-- RAISE NOTICE 'Iterator: ====================================================== '
				-- RAISE NOTICE 'Iterator: % ', row.code;
				-- RAISE NOTICE 'Iterator: % ', row.code_id;
				-- RAISE NOTICE 'Iterator: % ', i_price_type;
				-- RAISE NOTICE 'Iterator: % ', log_end_sec;
				-- EXIT ;
			END IF;

		END LOOP;
	end loop;	
	RAISE NOTICE 'Iterator: %', cnt ;
  
END;  $$
LANGUAGE plpgsql;
 select 1 from project.func_lines();

----------------------------------------
----------------------------------------
----------------------------------------
----------------------------------------
----------------------------------------
----------------------------------------
-- project.func_lines_tb_daily_line 업데이트
-- 
--  select 1 from project.func_lines_update_tb_daily_line();

----------------------------------------
----------------------------------------
----------------------------------------
----------------------------------------
----------------------------------------
CREATE OR REPLACE FUNCTION  project.func_lines_update_tb_daily_line(
	inp_code_id integer, inp_code VARCHAR , inp_price_type integer
	,x1 numeric  ,y1 numeric ,x2 numeric ,y2  numeric,x3 numeric ,y3  numeric)
 RETURNS void AS
 $$
BEGIN

 -- code  
--------------------
EXECUTE format('INSERT INTO project.tb_line(
	code_id, code,  price_type, x1, y1, x2, y2, x3, y3)
	VALUES (%L, %L,%L,
	%L, %L, %L, %L, %L, %L) on CONFLICT ("code_id","price_type") DO UPDATE  
	SET  x1=%L, y1=%L, x2=%L, y2=%L, x3=%L, y3=%L ;
	 ', inp_code_id, inp_code , inp_price_type
	,x1   ,y1  ,x2  ,y2  ,x3  ,y3 
	,x1   ,y1  ,x2  ,y2  ,x3  ,y3  
	) 
; 
--------------------

-------------
END;
$$
LANGUAGE plpgsql;
--select * from project.func_lines_update_tb_daily_line(1,'ss',1,1,2,3,4,5,6);



----------------------------------------
----------------------------------------
----------------------------------------
----------------------------------------
----------------------------------------
----------------------------------------
-- 
-- hist.bound의 마지막  점 두개 찾기  
--  select project.func_lines_get_last_point();

----------------------------------------
----------------------------------------
----------------------------------------
----------------------------------------
----------------------------------------
CREATE OR REPLACE FUNCTION  project.func_lines_get_last_point(inp_code_id integer, inp_price_type integer )
 RETURNS TABLE (count integer  ,x1 numeric  ,y1 numeric ,x2 numeric ,y2  numeric ) AS
 $$
declare
    r record;
	x1 numeric;
	x2 numeric;
	y1 numeric;
	y2 numeric;
	count integer := 0;
BEGIN

 -- code  
--------------------
FOR r in 
		SELECT t.X2 AS X, t.Y2 AS Y
			FROM HIST.rebound t
			WHERE t.CODE_id = inp_code_id
			AND t.price_type = inp_price_type
			ORDER  BY t.X2 DESC
			limit 2 
LOOP
	--------------------
	IF count = 1 THEN
		x1 := r.x ;
		y1 := r.y ;
	END IF ;

	IF count = 0 THEN
		x2 := r.x ;
		y2 := r.y ;
	END IF ;
	--------------------
	select count+1 into count;
END LOOP;
--------------------

RETURN QUERY select count, x1,y1,x2,y2;
-------------
END;
$$
LANGUAGE plpgsql;
--   select * from project.func_lines_get_last_point(1,3);

----------------------------------------
----------------------------------------
----------------------------------------
----------------------------------------
----------------------------------------
----------------------------------------
----------------------------------------
--  이전 두개의 점을 가지고 y = mx+b 그래프의 m,b를 구한 후 x+1 의 y값을 찾음.
--  select project.func_lines_get_next_point();

----------------------------------------
----------------------------------------
----------------------------------------
----------------------------------------
----------------------------------------
CREATE OR REPLACE FUNCTION  project.func_lines_get_next_point(
  inp_code_id integer ,inp_market_num integer
  ,inp_x1 numeric, inp_y1 numeric
  ,inp_x2 numeric, inp_y2 numeric
  )
 RETURNS TABLE (x3 numeric, y3 numeric) AS
 $$
declare
	x numeric;
	y numeric; 
BEGIN
 -- code  
RETURN QUERY
	select (inp_x2+1 )::numeric as x3, ((m*(x2+1))+b)::numeric  as y3  
	from (		
		select x2
			,round((y2-y1 )/( x2-x1)::numeric ,2 )as m 
			,round((y2-y1 )/( x2-x1)::numeric ,2 )*x2 *(-1) + y2 as b 
		from  (	
				select count(*) as x1 from hist.price where code_id = inp_code_id   and dt <= inp_x1
			)x1 ,(
				select count(*) as x2 from hist.price where code_id = inp_code_id  and dt <= inp_x2
			)x2,
			(
				select count(*) as y1 from utils.quote_unit where market = inp_market_num   and price <= inp_y1
			)y1 ,(
				select count(*) as y2 from utils.quote_unit where market = inp_market_num  and price <= inp_y2
			)y2
	)t;

END;
$$
LANGUAGE plpgsql;
-- select * from  project.func_lines_get_next_point(1,1, 20210621,14350  ,20210629,15600 );
----------------------------------------

----------------------------------------
----------------------------------------
----------------------------------------
----------------------------------------
----------------------------------------
----------------------------------------
----------------------------------------
-- inp_y 번째에 해당하는 마켓의 호가단위의 가격을 반환함.
--  select project.func_lines_convert_y3();

----------------------------------------
----------------------------------------
----------------------------------------
----------------------------------------
----------------------------------------
CREATE OR REPLACE FUNCTION  project.func_lines_convert_y3( inp_market_num integer , inp_y integer  )
 RETURNS TABLE ( price integer) AS
 $$
declare
	x integer;
	y integer; 
BEGIN
 -- code  
RETURN QUERY
	select qu.price 
		from  utils.quote_unit qu
		where qu.market = inp_market_num  offset inp_y  limit 1 ;

END;
$$
LANGUAGE plpgsql;
--select * from  project.func_lines_convert_y3(1, 6116  );
----------------------------------------




