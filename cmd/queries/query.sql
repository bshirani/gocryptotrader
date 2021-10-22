-- select count(*) from candle where timestamp > '09-19-2021';


select dhj.base, dhj.quote, count(*), ccount
from datahistoryjob dhj
JOIN (
    select c.base, c.quote, count(*) as ccount
    from candle c
    where timestamp > '08-19-2021'
    group by c.base, c.quote
    ORDER BY count(*)
) t on t.base = dhj.base and t.quote = dhj.quote
where status = 0
group by dhj.base, dhj.quote, ccount
order by count(*);




select t.date, base, quote, count(c.*)
FROM (
    select generate_series(
    (date '2013-10-01')::timestamp,
    (date '2013-10-02')::timestamp,
    interval '1 day') AS "date"
    ) t
LEFT OUTER JOIN candle c ON t.date = date_trunc('day', c.timestamp)
where
(timestamp >= '2013-10-01' and timestamp <= '2013-10-02') or timestamp is null
group by t.date, base, quote
order by t.date, base, quote;
