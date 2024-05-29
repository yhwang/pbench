SELECT TraficSourceID,
       SearchEngineID,
       AdvEngineID,
       CASE WHEN (SearchEngineID = 0 AND AdvEngineID = 0) THEN Referer ELSE '' END AS Src,
       URL                                                                         AS Dst,
       COUNT(*)                                                                    AS PageViews
FROM hits
WHERE CounterID = 62
  AND EventDate >= DATE('2013-07-01')
  AND EventDate <= DATE('2013-07-31')
  AND IsRefresh = 0
GROUP BY TraficSourceID,
         SearchEngineID,
         AdvEngineID,
         CASE WHEN (SearchEngineID = 0 AND AdvEngineID = 0) THEN Referer ELSE '' END,
         URL
ORDER BY PageViews DESC OFFSET 1000 LIMIT 10;
