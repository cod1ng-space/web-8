# Сервис "Count"


В демонстрационном решении используется СУБД Postgresql с БД sandbox, в которой имеется единственная таблица, созданная с помощью простейшего sql-кода ниже:

```
create table counter
(
	num integer not null
);
insert into counter(num) values (0);
```