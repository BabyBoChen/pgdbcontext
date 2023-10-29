SELECT A.column_name::varchar as fieldname
,pg_catalog.col_description(E.oid,A.ordinal_position) as shortdesc
,coalesce(A.character_maximum_length,0)+coalesce(A.numeric_precision,0) as datalength
,coalesce(A.numeric_scale,0) as numericscale 
,CASE WHEN A.is_nullable='YES'
	THEN true
	ELSE false
	END AS isallownull
,A.column_default as defaultvalue
,CASE WHEN constraint_type='PRIMARY KEY'
	THEN true
	ELSE false
	END AS isprimarykey
,case when A.is_identity='YES'
	then true
	else false
	end as isidentity
FROM information_schema.columns AS A
LEFT JOIN information_schema.constraint_column_usage AS B ON A.table_schema=B.table_schema AND A.table_name=B.table_name AND A.column_name=B.column_name
LEFT JOIN information_schema.table_constraints AS C ON B.table_schema=C.table_schema AND B.table_name=C.table_name AND B.constraint_name=C.constraint_name
LEFT JOIN pg_catalog.pg_namespace AS D ON D.nspname=A.table_schema
LEFT JOIN pg_catalog.pg_class AS E ON E.relnamespace=D.oid AND E.relname=A.table_name
WHERE A.table_schema='public' AND A.table_name='dish';