CREATE TABLE IF NOT EXISTS card_transactions (
    id INT PRIMARY KEY, -- autoincrement
    card_id TEXT NOT NULL,
    created DATETIME NOT NULL,
    amount INT NOT NULL,
    memo TEXT NOT NULL,
    FOREIGN KEY (card_id) REFERENCES cards(id)
);