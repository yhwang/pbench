SELECT UserID, extract(minute FROM EventTime), SearchPhrase, COUNT(*)
FROM hits
GROUP BY UserID, extract(minute FROM EventTime), SearchPhrase
ORDER BY COUNT(*) DESC LIMIT 10;
