delete q, r
from pbench_queries q
         join pbench_runs r on q.run_id = r.run_id
where r.run_name LIKE 'json041724-0426';
