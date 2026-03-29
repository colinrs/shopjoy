-- Alter earn_rules.status from VARCHAR(20) to TINYINT for iota enum alignment
-- Status: draft=0, active=1, inactive=2
ALTER TABLE earn_rules MODIFY COLUMN status TINYINT NOT NULL DEFAULT 0;