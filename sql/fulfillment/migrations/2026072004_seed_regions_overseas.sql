-- 海外区域种子数据：US, EU(主要国家), JP, KR, SEA
-- 平台预置数据，tenant_id=0
-- id 从 1000 起分配，避免与未来其他预置数据冲突

-- ============================================================
-- United States (51 行：1 root + 50 states)
-- ============================================================
INSERT INTO regions (id, tenant_id, code, name, level, parent_code, country_code, postal_pattern, sort, is_active, created_at, updated_at) VALUES
(1000, 0, 'US', 'United States', 1, '', 'US', '^[0-9]{5}(-[0-9]{4})?$', 0, 1, NOW(), NOW()),
(1001, 0, 'US-AL', 'Alabama',           2, 'US', 'US', NULL,  1, 1, NOW(), NOW()),
(1002, 0, 'US-AK', 'Alaska',            2, 'US', 'US', NULL,  2, 1, NOW(), NOW()),
(1003, 0, 'US-AZ', 'Arizona',           2, 'US', 'US', NULL,  3, 1, NOW(), NOW()),
(1004, 0, 'US-AR', 'Arkansas',          2, 'US', 'US', NULL,  4, 1, NOW(), NOW()),
(1005, 0, 'US-CA', 'California',        2, 'US', 'US', NULL,  5, 1, NOW(), NOW()),
(1006, 0, 'US-CO', 'Colorado',          2, 'US', 'US', NULL,  6, 1, NOW(), NOW()),
(1007, 0, 'US-CT', 'Connecticut',       2, 'US', 'US', NULL,  7, 1, NOW(), NOW()),
(1008, 0, 'US-DE', 'Delaware',          2, 'US', 'US', NULL,  8, 1, NOW(), NOW()),
(1009, 0, 'US-FL', 'Florida',           2, 'US', 'US', NULL,  9, 1, NOW(), NOW()),
(1010, 0, 'US-GA', 'Georgia',           2, 'US', 'US', NULL, 10, 1, NOW(), NOW()),
(1011, 0, 'US-HI', 'Hawaii',            2, 'US', 'US', NULL, 11, 1, NOW(), NOW()),
(1012, 0, 'US-ID', 'Idaho',             2, 'US', 'US', NULL, 12, 1, NOW(), NOW()),
(1013, 0, 'US-IL', 'Illinois',          2, 'US', 'US', NULL, 13, 1, NOW(), NOW()),
(1014, 0, 'US-IN', 'Indiana',           2, 'US', 'US', NULL, 14, 1, NOW(), NOW()),
(1015, 0, 'US-IA', 'Iowa',              2, 'US', 'US', NULL, 15, 1, NOW(), NOW()),
(1016, 0, 'US-KS', 'Kansas',            2, 'US', 'US', NULL, 16, 1, NOW(), NOW()),
(1017, 0, 'US-KY', 'Kentucky',          2, 'US', 'US', NULL, 17, 1, NOW(), NOW()),
(1018, 0, 'US-LA', 'Louisiana',         2, 'US', 'US', NULL, 18, 1, NOW(), NOW()),
(1019, 0, 'US-ME', 'Maine',             2, 'US', 'US', NULL, 19, 1, NOW(), NOW()),
(1020, 0, 'US-MD', 'Maryland',          2, 'US', 'US', NULL, 20, 1, NOW(), NOW()),
(1021, 0, 'US-MA', 'Massachusetts',     2, 'US', 'US', NULL, 21, 1, NOW(), NOW()),
(1022, 0, 'US-MI', 'Michigan',          2, 'US', 'US', NULL, 22, 1, NOW(), NOW()),
(1023, 0, 'US-MN', 'Minnesota',         2, 'US', 'US', NULL, 23, 1, NOW(), NOW()),
(1024, 0, 'US-MS', 'Mississippi',       2, 'US', 'US', NULL, 24, 1, NOW(), NOW()),
(1025, 0, 'US-MO', 'Missouri',          2, 'US', 'US', NULL, 25, 1, NOW(), NOW()),
(1026, 0, 'US-MT', 'Montana',           2, 'US', 'US', NULL, 26, 1, NOW(), NOW()),
(1027, 0, 'US-NE', 'Nebraska',          2, 'US', 'US', NULL, 27, 1, NOW(), NOW()),
(1028, 0, 'US-NV', 'Nevada',            2, 'US', 'US', NULL, 28, 1, NOW(), NOW()),
(1029, 0, 'US-NH', 'New Hampshire',     2, 'US', 'US', NULL, 29, 1, NOW(), NOW()),
(1030, 0, 'US-NJ', 'New Jersey',        2, 'US', 'US', NULL, 30, 1, NOW(), NOW()),
(1031, 0, 'US-NM', 'New Mexico',        2, 'US', 'US', NULL, 31, 1, NOW(), NOW()),
(1032, 0, 'US-NY', 'New York',          2, 'US', 'US', NULL, 32, 1, NOW(), NOW()),
(1033, 0, 'US-NC', 'North Carolina',    2, 'US', 'US', NULL, 33, 1, NOW(), NOW()),
(1034, 0, 'US-ND', 'North Dakota',      2, 'US', 'US', NULL, 34, 1, NOW(), NOW()),
(1035, 0, 'US-OH', 'Ohio',              2, 'US', 'US', NULL, 35, 1, NOW(), NOW()),
(1036, 0, 'US-OK', 'Oklahoma',          2, 'US', 'US', NULL, 36, 1, NOW(), NOW()),
(1037, 0, 'US-OR', 'Oregon',            2, 'US', 'US', NULL, 37, 1, NOW(), NOW()),
(1038, 0, 'US-PA', 'Pennsylvania',      2, 'US', 'US', NULL, 38, 1, NOW(), NOW()),
(1039, 0, 'US-RI', 'Rhode Island',      2, 'US', 'US', NULL, 39, 1, NOW(), NOW()),
(1040, 0, 'US-SC', 'South Carolina',    2, 'US', 'US', NULL, 40, 1, NOW(), NOW()),
(1041, 0, 'US-SD', 'South Dakota',      2, 'US', 'US', NULL, 41, 1, NOW(), NOW()),
(1042, 0, 'US-TN', 'Tennessee',         2, 'US', 'US', NULL, 42, 1, NOW(), NOW()),
(1043, 0, 'US-TX', 'Texas',             2, 'US', 'US', NULL, 43, 1, NOW(), NOW()),
(1044, 0, 'US-UT', 'Utah',              2, 'US', 'US', NULL, 44, 1, NOW(), NOW()),
(1045, 0, 'US-VT', 'Vermont',           2, 'US', 'US', NULL, 45, 1, NOW(), NOW()),
(1046, 0, 'US-VA', 'Virginia',          2, 'US', 'US', NULL, 46, 1, NOW(), NOW()),
(1047, 0, 'US-WA', 'Washington',        2, 'US', 'US', NULL, 47, 1, NOW(), NOW()),
(1048, 0, 'US-WV', 'West Virginia',     2, 'US', 'US', NULL, 48, 1, NOW(), NOW()),
(1049, 0, 'US-WI', 'Wisconsin',         2, 'US', 'US', NULL, 49, 1, NOW(), NOW()),
(1050, 0, 'US-WY', 'Wyoming',           2, 'US', 'US', NULL, 50, 1, NOW(), NOW());

-- ============================================================
-- European Union 主要国家 (5 行)
-- ============================================================
INSERT INTO regions (id, tenant_id, code, name, level, parent_code, country_code, postal_pattern, sort, is_active, created_at, updated_at) VALUES
(2000, 0, 'DE', 'Germany',     1, '', 'DE', '^[0-9]{5}$',            0, 1, NOW(), NOW()),
(2001, 0, 'FR', 'France',      1, '', 'FR', '^[0-9]{5}$',            0, 1, NOW(), NOW()),
(2002, 0, 'IT', 'Italy',       1, '', 'IT', '^[0-9]{5}$',            0, 1, NOW(), NOW()),
(2003, 0, 'ES', 'Spain',       1, '', 'ES', '^[0-9]{5}$',            0, 1, NOW(), NOW()),
(2004, 0, 'NL', 'Netherlands', 1, '', 'NL', '^[0-9]{4}\\s?[A-Z]{2}$', 0, 1, NOW(), NOW());

-- ============================================================
-- Japan (48 行：1 root + 47 都道府県)
-- ============================================================
INSERT INTO regions (id, tenant_id, code, name, level, parent_code, country_code, postal_pattern, sort, is_active, created_at, updated_at) VALUES
(3000, 0, 'JP', 'Japan', 1, '', 'JP', '^[0-9]{3}-[0-9]{4}$', 0, 1, NOW(), NOW()),
(3001, 0, 'JP-01', 'Hokkaido',         2, 'JP', 'JP', NULL,  1, 1, NOW(), NOW()),
(3002, 0, 'JP-02', 'Aomori',           2, 'JP', 'JP', NULL,  2, 1, NOW(), NOW()),
(3003, 0, 'JP-03', 'Iwate',            2, 'JP', 'JP', NULL,  3, 1, NOW(), NOW()),
(3004, 0, 'JP-04', 'Miyagi',           2, 'JP', 'JP', NULL,  4, 1, NOW(), NOW()),
(3005, 0, 'JP-05', 'Akita',            2, 'JP', 'JP', NULL,  5, 1, NOW(), NOW()),
(3006, 0, 'JP-06', 'Yamagata',         2, 'JP', 'JP', NULL,  6, 1, NOW(), NOW()),
(3007, 0, 'JP-07', 'Fukushima',        2, 'JP', 'JP', NULL,  7, 1, NOW(), NOW()),
(3008, 0, 'JP-08', 'Ibaraki',          2, 'JP', 'JP', NULL,  8, 1, NOW(), NOW()),
(3009, 0, 'JP-09', 'Tochigi',          2, 'JP', 'JP', NULL,  9, 1, NOW(), NOW()),
(3010, 0, 'JP-10', 'Gunma',            2, 'JP', 'JP', NULL, 10, 1, NOW(), NOW()),
(3011, 0, 'JP-11', 'Saitama',          2, 'JP', 'JP', NULL, 11, 1, NOW(), NOW()),
(3012, 0, 'JP-12', 'Chiba',            2, 'JP', 'JP', NULL, 12, 1, NOW(), NOW()),
(3013, 0, 'JP-13', 'Tokyo',            2, 'JP', 'JP', NULL, 13, 1, NOW(), NOW()),
(3014, 0, 'JP-14', 'Kanagawa',         2, 'JP', 'JP', NULL, 14, 1, NOW(), NOW()),
(3015, 0, 'JP-15', 'Niigata',          2, 'JP', 'JP', NULL, 15, 1, NOW(), NOW()),
(3016, 0, 'JP-16', 'Toyama',           2, 'JP', 'JP', NULL, 16, 1, NOW(), NOW()),
(3017, 0, 'JP-17', 'Ishikawa',         2, 'JP', 'JP', NULL, 17, 1, NOW(), NOW()),
(3018, 0, 'JP-18', 'Fukui',            2, 'JP', 'JP', NULL, 18, 1, NOW(), NOW()),
(3019, 0, 'JP-19', 'Yamanashi',        2, 'JP', 'JP', NULL, 19, 1, NOW(), NOW()),
(3020, 0, 'JP-20', 'Nagano',           2, 'JP', 'JP', NULL, 20, 1, NOW(), NOW()),
(3021, 0, 'JP-21', 'Gifu',             2, 'JP', 'JP', NULL, 21, 1, NOW(), NOW()),
(3022, 0, 'JP-22', 'Shizuoka',         2, 'JP', 'JP', NULL, 22, 1, NOW(), NOW()),
(3023, 0, 'JP-23', 'Aichi',            2, 'JP', 'JP', NULL, 23, 1, NOW(), NOW()),
(3024, 0, 'JP-24', 'Mie',              2, 'JP', 'JP', NULL, 24, 1, NOW(), NOW()),
(3025, 0, 'JP-25', 'Shiga',            2, 'JP', 'JP', NULL, 25, 1, NOW(), NOW()),
(3026, 0, 'JP-26', 'Kyoto',            2, 'JP', 'JP', NULL, 26, 1, NOW(), NOW()),
(3027, 0, 'JP-27', 'Osaka',            2, 'JP', 'JP', NULL, 27, 1, NOW(), NOW()),
(3028, 0, 'JP-28', 'Hyogo',            2, 'JP', 'JP', NULL, 28, 1, NOW(), NOW()),
(3029, 0, 'JP-29', 'Nara',             2, 'JP', 'JP', NULL, 29, 1, NOW(), NOW()),
(3030, 0, 'JP-30', 'Wakayama',         2, 'JP', 'JP', NULL, 30, 1, NOW(), NOW()),
(3031, 0, 'JP-31', 'Tottori',          2, 'JP', 'JP', NULL, 31, 1, NOW(), NOW()),
(3032, 0, 'JP-32', 'Shimane',          2, 'JP', 'JP', NULL, 32, 1, NOW(), NOW()),
(3033, 0, 'JP-33', 'Okayama',          2, 'JP', 'JP', NULL, 33, 1, NOW(), NOW()),
(3034, 0, 'JP-34', 'Hiroshima',        2, 'JP', 'JP', NULL, 34, 1, NOW(), NOW()),
(3035, 0, 'JP-35', 'Yamaguchi',        2, 'JP', 'JP', NULL, 35, 1, NOW(), NOW()),
(3036, 0, 'JP-36', 'Tokushima',        2, 'JP', 'JP', NULL, 36, 1, NOW(), NOW()),
(3037, 0, 'JP-37', 'Kagawa',           2, 'JP', 'JP', NULL, 37, 1, NOW(), NOW()),
(3038, 0, 'JP-38', 'Ehime',            2, 'JP', 'JP', NULL, 38, 1, NOW(), NOW()),
(3039, 0, 'JP-39', 'Kochi',            2, 'JP', 'JP', NULL, 39, 1, NOW(), NOW()),
(3040, 0, 'JP-40', 'Fukuoka',          2, 'JP', 'JP', NULL, 40, 1, NOW(), NOW()),
(3041, 0, 'JP-41', 'Saga',             2, 'JP', 'JP', NULL, 41, 1, NOW(), NOW()),
(3042, 0, 'JP-42', 'Nagasaki',         2, 'JP', 'JP', NULL, 42, 1, NOW(), NOW()),
(3043, 0, 'JP-43', 'Kumamoto',         2, 'JP', 'JP', NULL, 43, 1, NOW(), NOW()),
(3044, 0, 'JP-44', 'Oita',             2, 'JP', 'JP', NULL, 44, 1, NOW(), NOW()),
(3045, 0, 'JP-45', 'Miyazaki',         2, 'JP', 'JP', NULL, 45, 1, NOW(), NOW()),
(3046, 0, 'JP-46', 'Kagoshima',        2, 'JP', 'JP', NULL, 46, 1, NOW(), NOW()),
(3047, 0, 'JP-47', 'Okinawa',          2, 'JP', 'JP', NULL, 47, 1, NOW(), NOW());

-- ============================================================
-- South Korea (18 行：1 root + 17 시/도)
-- ============================================================
INSERT INTO regions (id, tenant_id, code, name, level, parent_code, country_code, postal_pattern, sort, is_active, created_at, updated_at) VALUES
(4000, 0, 'KR', 'South Korea', 1, '', 'KR', '^[0-9]{5}$', 0, 1, NOW(), NOW()),
(4001, 0, 'KR-11', 'Seoul',                          2, 'KR', 'KR', NULL,  1, 1, NOW(), NOW()),
(4002, 0, 'KR-21', 'Busan',                          2, 'KR', 'KR', NULL,  2, 1, NOW(), NOW()),
(4003, 0, 'KR-22', 'Daegu',                          2, 'KR', 'KR', NULL,  3, 1, NOW(), NOW()),
(4004, 0, 'KR-23', 'Incheon',                        2, 'KR', 'KR', NULL,  4, 1, NOW(), NOW()),
(4005, 0, 'KR-24', 'Gwangju',                        2, 'KR', 'KR', NULL,  5, 1, NOW(), NOW()),
(4006, 0, 'KR-25', 'Daejeon',                        2, 'KR', 'KR', NULL,  6, 1, NOW(), NOW()),
(4007, 0, 'KR-26', 'Ulsan',                          2, 'KR', 'KR', NULL,  7, 1, NOW(), NOW()),
(4008, 0, 'KR-29', 'Sejong',                         2, 'KR', 'KR', NULL,  8, 1, NOW(), NOW()),
(4009, 0, 'KR-31', 'Gyeonggi-do',                    2, 'KR', 'KR', NULL,  9, 1, NOW(), NOW()),
(4010, 0, 'KR-32', 'Gangwon-do',                     2, 'KR', 'KR', NULL, 10, 1, NOW(), NOW()),
(4011, 0, 'KR-33', 'Chungcheongbuk-do',              2, 'KR', 'KR', NULL, 11, 1, NOW(), NOW()),
(4012, 0, 'KR-34', 'Chungcheongnam-do',              2, 'KR', 'KR', NULL, 12, 1, NOW(), NOW()),
(4013, 0, 'KR-35', 'Jeollabuk-do',                   2, 'KR', 'KR', NULL, 13, 1, NOW(), NOW()),
(4014, 0, 'KR-36', 'Jeollanam-do',                   2, 'KR', 'KR', NULL, 14, 1, NOW(), NOW()),
(4015, 0, 'KR-37', 'Gyeongsangbuk-do',               2, 'KR', 'KR', NULL, 15, 1, NOW(), NOW()),
(4016, 0, 'KR-38', 'Gyeongsangnam-do',               2, 'KR', 'KR', NULL, 16, 1, NOW(), NOW()),
(4017, 0, 'KR-50', 'Jeju-do',                        2, 'KR', 'KR', NULL, 17, 1, NOW(), NOW());

-- ============================================================
-- Southeast Asia (6 行)
-- ============================================================
INSERT INTO regions (id, tenant_id, code, name, level, parent_code, country_code, postal_pattern, sort, is_active, created_at, updated_at) VALUES
(5000, 0, 'SG', 'Singapore',   1, '', 'SG', '^[0-9]{6}$', 0, 1, NOW(), NOW()),
(5001, 0, 'MY', 'Malaysia',    1, '', 'MY', '^[0-9]{5}$', 0, 1, NOW(), NOW()),
(5002, 0, 'TH', 'Thailand',    1, '', 'TH', '^[0-9]{5}$', 0, 1, NOW(), NOW()),
(5003, 0, 'PH', 'Philippines', 1, '', 'PH', '^[0-9]{4}$', 0, 1, NOW(), NOW()),
(5004, 0, 'ID', 'Indonesia',   1, '', 'ID', '^[0-9]{5}$', 0, 1, NOW(), NOW()),
(5005, 0, 'VN', 'Vietnam',     1, '', 'VN', '^[0-9]{6}$', 0, 1, NOW(), NOW());