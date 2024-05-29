SELECT WindowClientWidth, WindowClientHeight, COUNT(*) AS PageViews
FROM hits
WHERE CounterID = 62
  AND EventDate >= DATE('2013-07-01')
  AND EventDate <= DATE('2013-07-31')
  AND IsRefresh = 0
  AND DontCountHits = 0
  AND URLHash = 2868770270353813622
GROUP BY WindowClientWidth, WindowClientHeight
ORDER BY PageViews DESC OFFSET 10000 LIMIT 10;
