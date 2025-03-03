CREATE VIEW IF NOT EXISTS card_balances AS SELECT cards.id AS card_id, IFNULL(SUM(amount), 0) AS balance FROM cards LEFT JOIN card_transactions ON cards.id = card_transactions.card_id;
