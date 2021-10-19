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
