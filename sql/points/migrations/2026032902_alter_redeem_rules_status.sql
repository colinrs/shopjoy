-- Alter redeem_rules.status from VARCHAR(20) to TINYINT for iota enum alignment
-- Status: inactive=0, active=1
ALTER TABLE redeem_rules MODIFY COLUMN status TINYINT NOT NULL DEFAULT 0;