INSERT INTO shipping_rate
    (name, province_id, min_weight_g, max_weight_g, fee, is_active)
VALUES
    ('Default', NULL, 0, NULL, 50, true);
    -- ('Default 0-1kg',  NULL, 0,    1000, 45, true),
    -- ('Default 1-3kg',  NULL, 1000, 3000, 60, true),
    -- ('Default 3kg+',   NULL, 3000, NULL, 80, true),
    -- ('กรุงเทพมหานคร 0-1kg', 1, 0, 1000, 40, TRUE),
    -- ('กรุงเทพมหานคร 1-3kg', 1, 1000, 3000, 50, TRUE),
    -- ('กรุงเทพมหานคร 3kg+', 1, 3000, NULL, 65, TRUE);