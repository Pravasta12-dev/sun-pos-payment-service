ALTER TABLE transactions
ADD COLUMN transaction_scope VARCHAR(50) NOT NULL DEFAULT 'merchant';