-- M3 payment: record channel on recharge orders
ALTER TABLE recharge_orders
    ADD COLUMN IF NOT EXISTS payment_channel VARCHAR(20) DEFAULT 'mock';

CREATE INDEX IF NOT EXISTS idx_recharge_orders_status ON recharge_orders(status);
