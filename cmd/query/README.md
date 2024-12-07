# Сервис "Query"

В демонстрационном решении используется СУБД Postgresql с БД sandbox, в которой имеется единственная таблица, созданная с помощью простейшего sql-кода ниже:

```
create table query
(
	name text not null
);
insert into query(name) values ('');
```