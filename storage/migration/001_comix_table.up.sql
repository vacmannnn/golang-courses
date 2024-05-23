CREATE TABLE comics
(
    url      TEXT,
    keywords TEXT,
    comicsID INTEGER,
    UNIQUE (comicsID, url, keywords)
);
