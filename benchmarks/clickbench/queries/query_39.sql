SELECT URL, COUNT(*) AS PageViews
FROM hits
WHERE CounterID = 62
  AND EventDate >= DATE('2013-07-01')
  AND EventDate <= DATE('2013-07-31')
  AND IsRefresh = 0
  AND IsLink <> 0
  AND IsDownload = 0
GROUP BY URL
ORDER BY PageViews DESC OFFSET 1000 LIMIT 10;
