BEGIN;

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Abbotic', 'KFA-PROD-0002', 'Pcs', 'Farmasi Pediatri', 10, 271617, 300000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Acethion 600mg', 'KFA-PROD-0003', 'Pcs', 'Farmasi Pediatri', 10, 10212, 11500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Acetin 600 Tablet Effervecent', 'KFA-PROD-0004', 'Strip', 'Tablet & Kaplet', 10, 11655, 13000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Acetylcysteine 200mg', 'KFA-PROD-0005', 'Pcs', 'Farmasi Pediatri', 10, 1110, 1500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Acifar Krim', 'KFA-PROD-0006', 'Tube', 'Krim & Salep', 10, 1500, 2000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Acmist Moisturising Cream Gel 50 Gram', 'KFA-PROD-0007', 'Tube', 'Krim & Salep', 10, 216450, 233766, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Acnacare', 'KFA-PROD-0008', 'Pcs', 'Farmasi Pediatri', 10, 7700, 7700, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Aerochamber Plus Z L Mask', 'KFA-PROD-0009', 'Pcs', 'Alat Kesehatan', 10, 388500, 428000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Aerochamber Plus Z M Mask', 'KFA-PROD-0010', 'Pcs', 'Alat Kesehatan', 10, 431235, 475000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Aerochamber Plus Z S Mask', 'KFA-PROD-0011', 'Pcs', 'Alat Kesehatan', 10, 431235, 475000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Alco Flu Syrup', 'KFA-PROD-0012', 'Botol', 'Sirup & Suspensi', 10, 80000, 90000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Alerfed Plus Expectorant Syrup', 'KFA-PROD-0013', 'Botol', 'Sirup & Suspensi', 10, 33300, 37000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Alganax 0,5mg Tab', 'KFA-PROD-0014', 'Strip', 'Tablet & Kaplet', 10, 3885, 4500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Alkohol 70% 1liter', 'KFA-PROD-0015', 'Pcs', 'Farmasi Pediatri', 10, 30000, 49500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Alkohol swab', 'KFA-PROD-0016', 'Pcs', 'Farmasi Pediatri', 10, 100, 165, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Aloclair Plus Spray 15 ml', 'KFA-PROD-0017', 'Pcs', 'Farmasi Pediatri', 10, 109000, 120000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Amoxan Dry Syrup', 'KFA-PROD-0018', 'Botol', 'Sirup & Suspensi', 10, 25641, 28500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Amoxan Pedriatic Drop', 'KFA-PROD-0019', 'Botol', 'Drop Pediatri', 10, 25530, 28500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Amoxsan 500 mg Kaps', 'KFA-PROD-0020', 'Pcs', 'Farmasi Pediatri', 10, 3663, 4000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Apetic Suspensi 60ml', 'KFA-PROD-0021', 'Botol', 'Sirup & Suspensi', 10, 38850, 43000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Apialys Drop 10 ml', 'KFA-PROD-0022', 'Botol', 'Drop Pediatri', 10, 58275, 64000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Apialys Syr 100 ml', 'KFA-PROD-0023', 'Pcs', 'Farmasi Pediatri', 10, 48285, 53500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Apiban Drops', 'KFA-PROD-0024', 'Botol', 'Drop Pediatri', 10, 149850, 165000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Astharol Tablet', 'KFA-PROD-0025', 'Strip', 'Tablet & Kaplet', 10, 1610, 1900, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Atomist Cream', 'KFA-PROD-0026', 'Tube', 'Krim & Salep', 10, 1998, 2200, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Avamyst Nasal Spray', 'KFA-PROD-0027', 'Pcs', 'Farmasi Pediatri', 10, 341000, 375000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Avaxim 80 Hep A', 'KFA-PROD-0028', 'Pcs', 'Farmasi Pediatri', 10, 357605, 440000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Bariesun SPF 50 tinted', 'KFA-PROD-0029', 'Pcs', 'Farmasi Pediatri', 10, 5550, 6200, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Becomzet Capsul', 'KFA-PROD-0030', 'Strip', 'Kapsul', 10, 2609, 3000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Benoson Cream 5 gram', 'KFA-PROD-0031', 'Tube', 'Krim & Salep', 10, 3774, 4500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Bepanthen Baby 20gr', 'KFA-PROD-0032', 'Pcs', 'Farmasi Pediatri', 10, 67044, 74000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Bephanten Baby 10g', 'KFA-PROD-0033', 'Pcs', 'Farmasi Pediatri', 10, 41847, 46500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Bioplacenton jelly 15 gram', 'KFA-PROD-0034', 'Pcs', 'Farmasi Pediatri', 10, 33000, 36500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Bisolvon Solution 50 ml', 'KFA-PROD-0035', 'Botol', 'Sirup & Suspensi', 10, 95867, 106000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Bivalent Oral Poliomyelitis Vaksin', 'KFA-PROD-0036', 'Pcs', 'Farmasi Pediatri', 10, 11177, 15000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'BNS Breathy Nasal Spray', 'KFA-PROD-0037', 'Pcs', 'Farmasi Pediatri', 10, 97700, 108000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Bright Glow skin', 'KFA-PROD-0038', 'Pcs', 'Farmasi Pediatri', 10, 11000, 11000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Buscopan Plus Tablet', 'KFA-PROD-0039', 'Strip', 'Tablet & Kaplet', 10, 5533, 6500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Bye-Bye Fever Anak', 'KFA-PROD-0040', 'Pcs', 'Farmasi Pediatri', 10, 10700, 12500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Bye-Bye Fever Bayi', 'KFA-PROD-0041', 'Pcs', 'Farmasi Pediatri', 10, 7600, 9000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'CALENDULA OIL', 'KFA-PROD-0042', 'Pcs', 'Farmasi Pediatri', 10, 2000, 20000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Cefila Dry Syrup', 'KFA-PROD-0043', 'Botol', 'Sirup & Suspensi', 10, 114000, 130000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Cendo Cenfresh Minidose', 'KFA-PROD-0044', 'Pcs', 'Farmasi Pediatri', 10, 29800, 33000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Ceradan Advanced Mouist skin', 'KFA-PROD-0045', 'Pcs', 'Farmasi Pediatri', 10, 3182, 3700, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Ceradan Hydra 80gram', 'KFA-PROD-0046', 'Pcs', 'Farmasi Pediatri', 10, 2429, 2700, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Ceradan Shooting Gel', 'KFA-PROD-0047', 'Tube', 'Krim & Salep', 10, 2638, 2850, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Cetinal Syrup', 'KFA-PROD-0048', 'Botol', 'Sirup & Suspensi', 10, 97125, 107000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Cetinal Tab 10mg', 'KFA-PROD-0049', 'Strip', 'Tablet & Kaplet', 10, 7770, 8700, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Claneksi Dry Sirup', 'KFA-PROD-0050', 'Botol', 'Sirup & Suspensi', 10, 68265, 74000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Claneksi Forte Dry Syrup', 'KFA-PROD-0051', 'Botol', 'Sirup & Suspensi', 10, 90631, 99000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Clinovir 200 mg Tablet', 'KFA-PROD-0052', 'Strip', 'Tablet & Kaplet', 10, 5650, 6500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Clinovir Cream 5gr', 'KFA-PROD-0053', 'Tube', 'Krim & Salep', 10, 8560, 9500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Clobetasol Cream 0,05%', 'KFA-PROD-0054', 'Tube', 'Krim & Salep', 10, 1443, 2000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Cobazim 1000mg Cap', 'KFA-PROD-0055', 'Strip', 'Kapsul', 10, 6178, 7000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Cobazim 3000mg Cap', 'KFA-PROD-0056', 'Strip', 'Kapsul', 10, 13875, 16000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Combivent 2.5ml UDV', 'KFA-PROD-0057', 'Pcs', 'Farmasi Pediatri', 10, 9750, 13500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Comtusi Caramel Syrup 60 ml', 'KFA-PROD-0058', 'Botol', 'Sirup & Suspensi', 10, 43300, 47000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Comtusi Strawberry Syrup 60ml', 'KFA-PROD-0059', 'Botol', 'Sirup & Suspensi', 10, 43290, 47700, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Cortamine Syrup', 'KFA-PROD-0060', 'Botol', 'Sirup & Suspensi', 10, 71762, 80000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Cortamine Tablet', 'KFA-PROD-0061', 'Strip', 'Tablet & Kaplet', 10, 4431, 4900, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Cortidex Tab', 'KFA-PROD-0062', 'Strip', 'Tablet & Kaplet', 10, 416, 600, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'CTDM et N251754', 'KFA-PROD-0063', 'Pcs', 'Farmasi Pediatri', 10, 10000, 10000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'CTM 4 mg Tablet', 'KFA-PROD-0064', 'Strip', 'Tablet & Kaplet', 10, 700, 1000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Curcuma Force Tab', 'KFA-PROD-0065', 'Strip', 'Tablet & Kaplet', 10, 1083, 1500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Curvit Kaplet', 'KFA-PROD-0066', 'Strip', 'Tablet & Kaplet', 10, 2120, 3000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'D3 More 5000', 'KFA-PROD-0067', 'Pcs', 'Farmasi Pediatri', 10, 4440, 5500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'D40', 'KFA-PROD-0068', 'Pcs', 'Farmasi Pediatri', 10, 8975, 10000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'D5', 'KFA-PROD-0069', 'Pcs', 'Farmasi Pediatri', 10, 8423, 9500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Daily D3 Kids Drops', 'KFA-PROD-0070', 'Botol', 'Drop Pediatri', 10, 162060, 193000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Daylax Sachet', 'KFA-PROD-0071', 'Sachet', 'Serbuk & Sachet', 10, 25227, 30000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Dehidralite', 'KFA-PROD-0072', 'Pcs', 'Farmasi Pediatri', 10, 19536, 22000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Dermacure Cream', 'KFA-PROD-0073', 'Tube', 'Krim & Salep', 10, 8547, 9500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Dexamethasone Injeksi', 'KFA-PROD-0074', 'Vial', 'Injeksi & Cairan', 10, 2500, 3500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Dexomet Krim', 'KFA-PROD-0075', 'Tube', 'Krim & Salep', 10, 2775, 3300, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Digenta Cream', 'KFA-PROD-0076', 'Tube', 'Krim & Salep', 10, 11877, 13500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Dipenhydramine Injeksi', 'KFA-PROD-0077', 'Vial', 'Injeksi & Cairan', 10, 4200, 5000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Eau Thermale Reparatrice 150 ml', 'KFA-PROD-0078', 'Pcs', 'Farmasi Pediatri', 10, 138750, 160000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Elementre Cleansing Gel 250 ml', 'KFA-PROD-0079', 'Tube', 'Krim & Salep', 10, 510000, 510000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Elementre dermo  Nourishing cr', 'KFA-PROD-0080', 'Pcs', 'Farmasi Pediatri', 10, 16800, 18500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Elox Cream', 'KFA-PROD-0081', 'Tube', 'Krim & Salep', 10, 8658, 10500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Engerix Vaksin', 'KFA-PROD-0082', 'Pcs', 'Farmasi Pediatri', 10, 145000, 160000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Enzyplex', 'KFA-PROD-0083', 'Pcs', 'Farmasi Pediatri', 10, 2931, 3200, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Erysanbe 200mg Chewable Tablet', 'KFA-PROD-0084', 'Strip', 'Tablet & Kaplet', 10, 2331, 2517, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Estesia Cream 90 gram', 'KFA-PROD-0085', 'Tube', 'Krim & Salep', 10, 7523, 8500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Ferlin Drops 30 ml', 'KFA-PROD-0086', 'Botol', 'Drop Pediatri', 10, 36700, 41000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Ferro K Drop', 'KFA-PROD-0087', 'Botol', 'Drop Pediatri', 10, 43290, 46753, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Flamicort Tablet', 'KFA-PROD-0088', 'Strip', 'Tablet & Kaplet', 10, 4440, 5000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Folda Kaplet', 'KFA-PROD-0089', 'Strip', 'Tablet & Kaplet', 10, 5000, 6000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Fluhalt Powder for Suspension', 'KFA-PROD-0090', 'Botol', 'Sirup & Suspensi', 10, 960000, 970000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Folavit 1000mcg Tab', 'KFA-PROD-0091', 'Strip', 'Tablet & Kaplet', 10, 2442, 2637, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Forti D 1000', 'KFA-PROD-0092', 'Pcs', 'Farmasi Pediatri', 10, 2735, 2955, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Fortiboost D3+K2 Drops', 'KFA-PROD-0093', 'Botol', 'Drop Pediatri', 10, 75480, 83000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Forumen Ear Drop 10 ml', 'KFA-PROD-0094', 'Botol', 'Drop Pediatri', 10, 33700, 38000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Fuladic Cream', 'KFA-PROD-0095', 'Tube', 'Krim & Salep', 10, 12876, 15000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Fungasol Cream Tube', 'KFA-PROD-0096', 'Tube', 'Krim & Salep', 10, 3108, 3500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Fungasol SS 2% 80ml', 'KFA-PROD-0097', 'Pcs', 'Farmasi Pediatri', 10, 97125, 104895, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Glizigen Spray 30 ml', 'KFA-PROD-0098', 'Pcs', 'Farmasi Pediatri', 10, 488500, 538000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Glusela Tablet', 'KFA-PROD-0099', 'Strip', 'Tablet & Kaplet', 10, 16650, 18500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Havrix 720 Junior', 'KFA-PROD-0100', 'Pcs', 'Farmasi Pediatri', 10, 431256, 505000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Healive Vaksin', 'KFA-PROD-0101', 'Pcs', 'Farmasi Pediatri', 10, 255300, 420000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Heptasan Tab', 'KFA-PROD-0102', 'Strip', 'Tablet & Kaplet', 10, 512, 1000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Hexaxim Vaksin', 'KFA-PROD-0103', 'Pcs', 'Farmasi Pediatri', 10, 880785, 960000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Iliadin drop', 'KFA-PROD-0104', 'Botol', 'Drop Pediatri', 10, 84844, 95000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'IMOJEV Vaksin', 'KFA-PROD-0105', 'Pcs', 'Farmasi Pediatri', 10, 444000, 500000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Imunkid Drop', 'KFA-PROD-0106', 'Botol', 'Drop Pediatri', 10, 61050, 67500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Imunped Drop', 'KFA-PROD-0107', 'Botol', 'Drop Pediatri', 10, 58275, 64000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Imunped Syrup', 'KFA-PROD-0108', 'Botol', 'Sirup & Suspensi', 10, 56610, 63000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Inadine Solution 1Liter', 'KFA-PROD-0109', 'Botol', 'Sirup & Suspensi', 10, 222, 240, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Inerson cream 15 gram', 'KFA-PROD-0110', 'Tube', 'Krim & Salep', 10, 11200, 12500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Infanrix Hexa', 'KFA-PROD-0111', 'Pcs', 'Farmasi Pediatri', 10, 777000, 910000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Inlive Vaksin', 'KFA-PROD-0112', 'Pcs', 'Farmasi Pediatri', 10, 876900, 950000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Inpepsa Suspensi 100 ml', 'KFA-PROD-0113', 'Botol', 'Sirup & Suspensi', 10, 98600, 108500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Interlac Drops 5ml', 'KFA-PROD-0114', 'Botol', 'Drop Pediatri', 10, 366300, 403000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Interlac Sachet', 'KFA-PROD-0115', 'Sachet', 'Serbuk & Sachet', 10, 13875, 14985, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Interlac Tab Kunyah', 'KFA-PROD-0116', 'Strip', 'Tablet & Kaplet', 10, 13209, 14700, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Intrizin Syrup', 'KFA-PROD-0117', 'Botol', 'Sirup & Suspensi', 10, 90576, 100000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Intrizin Tablet', 'KFA-PROD-0118', 'Strip', 'Tablet & Kaplet', 10, 9426, 10500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Isoprinosine Syrup', 'KFA-PROD-0119', 'Botol', 'Sirup & Suspensi', 10, 111451, 122500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Isprinol Tab', 'KFA-PROD-0120', 'Strip', 'Tablet & Kaplet', 10, 13431, 15000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Kenacort Tab 4 mg', 'KFA-PROD-0121', 'Strip', 'Tablet & Kaplet', 10, 5850, 6500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Kenalog in oral base 5g', 'KFA-PROD-0122', 'Pcs', 'Farmasi Pediatri', 10, 72150, 79500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Koentoem sabun mandi', 'KFA-PROD-0123', 'Pcs', 'Farmasi Pediatri', 10, 560, 700, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Lacto B', 'KFA-PROD-0124', 'Pcs', 'Farmasi Pediatri', 10, 9380, 11000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Lameson 4 mg Tablet', 'KFA-PROD-0125', 'Strip', 'Tablet & Kaplet', 10, 5000, 6000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Lamodex Cream', 'KFA-PROD-0126', 'Tube', 'Krim & Salep', 10, 4995, 5500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'LAPIFED DM  SIRUP', 'KFA-PROD-0127', 'Botol', 'Sirup & Suspensi', 10, 45000, 50000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Lapifed Expectorant', 'KFA-PROD-0128', 'Pcs', 'Farmasi Pediatri', 10, 39600, 44000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Lasal tab 2 mg', 'KFA-PROD-0129', 'Strip', 'Tablet & Kaplet', 10, 2137, 2308, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Lasal tab 4 mg', 'KFA-PROD-0130', 'Strip', 'Tablet & Kaplet', 10, 2747, 2967, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Lidocain Hcl 2%', 'KFA-PROD-0131', 'Pcs', 'Farmasi Pediatri', 10, 1665, 2500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Lidocain Injeksi', 'KFA-PROD-0132', 'Vial', 'Injeksi & Cairan', 10, 2000, 3500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Liprolac Baby Drop', 'KFA-PROD-0133', 'Botol', 'Drop Pediatri', 10, 316350, 348000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Liprolac Vanilla Sachet', 'KFA-PROD-0134', 'Sachet', 'Serbuk & Sachet', 10, 10545, 12000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Lodia Tab', 'KFA-PROD-0135', 'Strip', 'Tablet & Kaplet', 10, 1388, 1600, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Lymax Sirup', 'KFA-PROD-0136', 'Botol', 'Sirup & Suspensi', 10, 47500, 53000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Lymax Drops', 'KFA-PROD-0137', 'Botol', 'Drop Pediatri', 10, 55500, 62000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Lupred Tablet', 'KFA-PROD-0138', 'Strip', 'Tablet & Kaplet', 10, 4773, 5500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Maganol Forte Tab', 'KFA-PROD-0139', 'Strip', 'Tablet & Kaplet', 10, 555, 1000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Maltofer Drops 30ml', 'KFA-PROD-0140', 'Botol', 'Drop Pediatri', 10, 72317, 80000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Maltofer Fol Chewable Tab', 'KFA-PROD-0141', 'Strip', 'Tablet & Kaplet', 10, 3700, 4100, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Maxprinol', 'KFA-PROD-0142', 'Pcs', 'Farmasi Pediatri', 10, 12000, 13000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Mazthion tab', 'KFA-PROD-0143', 'Strip', 'Tablet & Kaplet', 10, 25000, 35000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'MenQuadfi Vaksin', 'KFA-PROD-0144', 'Pcs', 'Farmasi Pediatri', 10, 703000, 760000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Meptin Mini  25 mcg Tablet', 'KFA-PROD-0145', 'Strip', 'Tablet & Kaplet', 10, 3000, 3500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Memucil Kapsul', 'KFA-PROD-0146', 'Strip', 'Kapsul', 10, 4000, 5000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Mesone Cream', 'KFA-PROD-0147', 'Tube', 'Krim & Salep', 10, 7215, 8000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'MMR-II VAKSIN', 'KFA-PROD-0148', 'Pcs', 'Farmasi Pediatri', 10, 449550, 550000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Moisturizer Uriage hysec', 'KFA-PROD-0149', 'Pcs', 'Farmasi Pediatri', 10, 290000, 290000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Mucos Drop', 'KFA-PROD-0150', 'Botol', 'Drop Pediatri', 10, 38295, 42500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Mucos Syrup', 'KFA-PROD-0151', 'Botol', 'Sirup & Suspensi', 10, 18315, 20500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Mucos Tablet', 'KFA-PROD-0152', 'Strip', 'Tablet & Kaplet', 10, 888, 1200, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Mupirocin Cream', 'KFA-PROD-0153', 'Tube', 'Krim & Salep', 10, 4440, 5500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Myconys Drops 12 ml', 'KFA-PROD-0154', 'Botol', 'Drop Pediatri', 10, 54500, 60000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Myco Z Cream', 'KFA-PROD-0155', 'Tube', 'Krim & Salep', 10, 98790, 109000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Narfoz Syrup 30 ml', 'KFA-PROD-0156', 'Botol', 'Sirup & Suspensi', 10, 71200, 78500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Narfoz Syrup 60 ml', 'KFA-PROD-0157', 'Botol', 'Sirup & Suspensi', 10, 110000, 121000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Nebulizer Mask & Tubing', 'KFA-PROD-0158', 'Pcs', 'Alat Kesehatan', 10, 15000, 16500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Neo Kaolana Syrup', 'KFA-PROD-0159', 'Botol', 'Sirup & Suspensi', 10, 15540, 18000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Neoral SGC 100 mg', 'KFA-PROD-0160', 'Pcs', 'Farmasi Pediatri', 10, 78333, 86200, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Neostrata Clarifying Peel', 'KFA-PROD-0161', 'Pcs', 'Farmasi Pediatri', 10, 38850, 43000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Neostrata Rejuvenating Peel 35%', 'KFA-PROD-0162', 'Pcs', 'Farmasi Pediatri', 10, 31783, 35000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Nexium 10 mg Sachet', 'KFA-PROD-0163', 'Sachet', 'Serbuk & Sachet', 10, 22200, 25000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Nifural Syrup 60ml', 'KFA-PROD-0164', 'Botol', 'Sirup & Suspensi', 10, 68607, 76000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Nipe Drops 15 ml', 'KFA-PROD-0165', 'Botol', 'Drop Pediatri', 10, 96000, 106000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Norages Drops', 'KFA-PROD-0166', 'Botol', 'Drop Pediatri', 10, 43290, 48000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Norages Syrup', 'KFA-PROD-0167', 'Botol', 'Sirup & Suspensi', 10, 35520, 40000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'normagut', 'KFA-PROD-0168', 'Pcs', 'Farmasi Pediatri', 10, 9364, 11000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'NRC (Krim Malam)', 'KFA-PROD-0169', 'Tube', 'Krim & Salep', 10, 100000, 100000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Nyndia Drops 12ml', 'KFA-PROD-0170', 'Botol', 'Drop Pediatri', 10, 48285, 53500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Nystatin Drops', 'KFA-PROD-0171', 'Botol', 'Drop Pediatri', 10, 25530, 29000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Omeprazole 20 mg Kapsul', 'KFA-PROD-0172', 'Strip', 'Kapsul', 10, 800, 1200, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Ondane 4mg', 'KFA-PROD-0173', 'Pcs', 'Farmasi Pediatri', 10, 14985, 17000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Ondansetron Injeksi', 'KFA-PROD-0174', 'Vial', 'Injeksi & Cairan', 10, 4300, 5000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Ossovit Kapsul', 'KFA-PROD-0175', 'Strip', 'Kapsul', 10, 7000, 8000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Ossovit Syrup', 'KFA-PROD-0176', 'Botol', 'Sirup & Suspensi', 10, 101898, 115000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Pamol 125 mg Suppo', 'KFA-PROD-0177', 'Pcs', 'Farmasi Pediatri', 10, 18800, 21000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Pamol 250 mg Suppo', 'KFA-PROD-0178', 'Pcs', 'Farmasi Pediatri', 10, 27800, 32000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Papulex Cleansing Gel 150 ml', 'KFA-PROD-0179', 'Tube', 'Krim & Salep', 10, 1764, 3000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Paratusin DMP Tab', 'KFA-PROD-0180', 'Strip', 'Tablet & Kaplet', 10, 1498, 1700, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'PCV 20 Vaksin', 'KFA-PROD-0181', 'Pcs', 'Farmasi Pediatri', 10, 945000, 1050000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Pedilax Sachet', 'KFA-PROD-0182', 'Sachet', 'Serbuk & Sachet', 10, 33300, 37000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Pehacain Injeksi', 'KFA-PROD-0183', 'Vial', 'Injeksi & Cairan', 10, 4000, 6600, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Pehacain Injeksi Sampel', 'KFA-PROD-0184', 'Vial', 'Injeksi & Cairan', 10, 4000, 6000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Pepsamax Susp', 'KFA-PROD-0185', 'Pcs', 'Farmasi Pediatri', 10, 75500, 83000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Plasterin', 'KFA-PROD-0186', 'Pcs', 'Farmasi Pediatri', 10, 14, 23, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Polysilane Suspensi 100 ml', 'KFA-PROD-0187', 'Botol', 'Sirup & Suspensi', 10, 24730, 27500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Pot Acrylic 10gram', 'KFA-PROD-0188', 'Pcs', 'Farmasi Pediatri', 10, 3608, 3900, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Pot Acrylic 5gram', 'KFA-PROD-0189', 'Pcs', 'Farmasi Pediatri', 10, 3330, 3600, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Pot Airless Jar 15 gram', 'KFA-PROD-0190', 'Pcs', 'Farmasi Pediatri', 10, 16650, 18000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Prevenar 20 Vaksin', 'KFA-PROD-0191', 'Pcs', 'Farmasi Pediatri', 10, 950000, 1050000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Priorix', 'KFA-PROD-0192', 'Pcs', 'Farmasi Pediatri', 10, 464445, 577000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Probiokid', 'KFA-PROD-0193', 'Pcs', 'Farmasi Pediatri', 10, 17650, 19500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Prolacta With DHA for Baby', 'KFA-PROD-0194', 'Pcs', 'Farmasi Pediatri', 10, 6600, 8000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Promuba Suspensi', 'KFA-PROD-0195', 'Botol', 'Sirup & Suspensi', 10, 43290, 48000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Propyretic 160 Suppositoria', 'KFA-PROD-0196', 'Pcs', 'Farmasi Pediatri', 10, 9018, 10000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'PROQUAD VAKSIN', 'KFA-PROD-0197', 'Pcs', 'Farmasi Pediatri', 10, 1176178, 1270000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Proris 200 mg Caplet', 'KFA-PROD-0198', 'Strip', 'Kapsul', 10, 889, 1000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Proris Forte Suspensi', 'KFA-PROD-0199', 'Botol', 'Sirup & Suspensi', 10, 31169, 34500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Proris Suppositoria', 'KFA-PROD-0200', 'Pcs', 'Farmasi Pediatri', 10, 7937, 8800, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Proris Susp 60 ml Hijau', 'KFA-PROD-0201', 'Pcs', 'Farmasi Pediatri', 10, 31500, 35000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Prosentials', 'KFA-PROD-0202', 'Pcs', 'Farmasi Pediatri', 10, 7260, 9000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Prove D3 1000iu Tab', 'KFA-PROD-0203', 'Strip', 'Tablet & Kaplet', 10, 3219, 3600, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Prove D3 5000iu Tab', 'KFA-PROD-0204', 'Strip', 'Tablet & Kaplet', 10, 4718, 5200, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Prove D3 Drops', 'KFA-PROD-0205', 'Botol', 'Drop Pediatri', 10, 244200, 270000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'PROVE D3 DROPS SAMPEL', 'KFA-PROD-0206', 'Botol', 'Drop Pediatri', 10, 130000, 250000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Prove D3+K2', 'KFA-PROD-0207', 'Pcs', 'Farmasi Pediatri', 10, 3885, 4500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Psoramist Cream 150 gram', 'KFA-PROD-0208', 'Tube', 'Krim & Salep', 10, 771450, 833166, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Psoramist Cream 150 gram Sampel', 'KFA-PROD-0209', 'Tube', 'Krim & Salep', 10, 771450, 833166, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Pulmicort 0,25 mg/ml', 'KFA-PROD-0210', 'Pcs', 'Farmasi Pediatri', 10, 22500, 24000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Pulmicort 0,5 mg/ml', 'KFA-PROD-0211', 'Pcs', 'Farmasi Pediatri', 10, 27500, 31000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'QDenga Vaksin', 'KFA-PROD-0212', 'Pcs', 'Farmasi Pediatri', 10, 468007, 605000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Ranitidin Injeksi', 'KFA-PROD-0213', 'Vial', 'Injeksi & Cairan', 10, 4500, 5500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Rhinofed tab', 'KFA-PROD-0214', 'Strip', 'Tablet & Kaplet', 10, 3000, 3500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Rhinos Neo Drops', 'KFA-PROD-0215', 'Botol', 'Drop Pediatri', 10, 73260, 81000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Rillus Tab', 'KFA-PROD-0216', 'Strip', 'Tablet & Kaplet', 10, 10545, 11600, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Rotarix Vaksin', 'KFA-PROD-0217', 'Pcs', 'Farmasi Pediatri', 10, 386414, 463000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Rotateq Vaksin', 'KFA-PROD-0218', 'Pcs', 'Farmasi Pediatri', 10, 300307, 365000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Rymont 10 mg', 'KFA-PROD-0219', 'Pcs', 'Farmasi Pediatri', 10, 9397, 10500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Rymont 4 mg ODF', 'KFA-PROD-0220', 'Pcs', 'Farmasi Pediatri', 10, 9990, 11000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Ryvel Drops', 'KFA-PROD-0221', 'Botol', 'Drop Pediatri', 10, 101010, 112000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'SABUN BABY', 'KFA-PROD-0222', 'Pcs', 'Farmasi Pediatri', 10, 2000, 25000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Salep antibiotik 5 gr', 'KFA-PROD-0223', 'Tube', 'Krim & Salep', 10, 30000, 30000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Salgen Cr', 'KFA-PROD-0224', 'Pcs', 'Farmasi Pediatri', 10, 1364, 3000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Salofalk Granules', 'KFA-PROD-0225', 'Pcs', 'Farmasi Pediatri', 10, 63500, 70000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Sanmol Drop', 'KFA-PROD-0226', 'Botol', 'Drop Pediatri', 10, 20868, 23000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Sanmol Forte Syrup', 'KFA-PROD-0227', 'Botol', 'Sirup & Suspensi', 10, 34077, 37500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Sanmol Syrup', 'KFA-PROD-0228', 'Botol', 'Sirup & Suspensi', 10, 19092, 22000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'SANMOL SYRUP SAMPEL', 'KFA-PROD-0229', 'Botol', 'Sirup & Suspensi', 10, 50000, 80000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Sanprima Suspensi', 'KFA-PROD-0230', 'Botol', 'Sirup & Suspensi', 10, 32911, 36500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Sanprinol Syrup', 'KFA-PROD-0231', 'Botol', 'Sirup & Suspensi', 10, 83250, 92000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Sarung Tangan Latex', 'KFA-PROD-0232', 'Pcs', 'Farmasi Pediatri', 10, 400, 432, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Sarung Tangan Nitril Altamed', 'KFA-PROD-0233', 'Pcs', 'Farmasi Pediatri', 10, 580, 1000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Scacid Cream 30 Gram', 'KFA-PROD-0234', 'Tube', 'Krim & Salep', 10, 86580, 95500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Seloxy AA', 'KFA-PROD-0235', 'Pcs', 'Farmasi Pediatri', 10, 7770, 10000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Seloxy Premium', 'KFA-PROD-0236', 'Pcs', 'Farmasi Pediatri', 10, 9768, 12000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Serum basic (neotone isispharma)', 'KFA-PROD-0237', 'Pcs', 'Farmasi Pediatri', 10, 15873, 17500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Silex Syrup', 'KFA-PROD-0238', 'Botol', 'Sirup & Suspensi', 10, 82140, 90500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'SIM - DHA Syr 100 ml', 'KFA-PROD-0239', 'Pcs', 'Farmasi Pediatri', 10, 61000, 67500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Sirplus Rasa Anggur', 'KFA-PROD-0240', 'Pcs', 'Farmasi Pediatri', 10, 15540, 17000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Sirplus Rasa Jeruk', 'KFA-PROD-0241', 'Pcs', 'Farmasi Pediatri', 10, 15540, 17000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Sirplus Rasa Stroberi', 'KFA-PROD-0242', 'Pcs', 'Farmasi Pediatri', 10, 15540, 17000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Sirplus Stroberi Tab', 'KFA-PROD-0243', 'Strip', 'Tablet & Kaplet', 10, 233, 252, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Sporetik 100 mg Caps', 'KFA-PROD-0244', 'Strip', 'Kapsul', 10, 25419, 28000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Sporetik 200 mg Caps', 'KFA-PROD-0245', 'Strip', 'Kapsul', 10, 36750, 40500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Sporetik Dry Syrup', 'KFA-PROD-0246', 'Botol', 'Sirup & Suspensi', 10, 107392, 118000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Spuit 10cc Terumo', 'KFA-PROD-0247', 'Pcs', 'Alat Kesehatan', 10, 3250, 5400, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Spuit 1cc Onemed', 'KFA-PROD-0248', 'Pcs', 'Alat Kesehatan', 10, 800, 1320, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Spuit 1cc Terumo', 'KFA-PROD-0249', 'Pcs', 'Alat Kesehatan', 10, 2500, 4200, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Spuit 20cc Terumo', 'KFA-PROD-0250', 'Pcs', 'Alat Kesehatan', 10, 6500, 10800, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Spuit 3cc Onemed', 'KFA-PROD-0251', 'Pcs', 'Alat Kesehatan', 10, 600, 1000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Spuit 3cc Terumo', 'KFA-PROD-0252', 'Pcs', 'Alat Kesehatan', 10, 2350, 3900, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Spuit 5cc Onemed', 'KFA-PROD-0253', 'Pcs', 'Alat Kesehatan', 10, 620, 1100, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Spuit 5cc Terumo', 'KFA-PROD-0254', 'Pcs', 'Alat Kesehatan', 10, 3000, 4950, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'SUNSCREEN', 'KFA-PROD-0255', 'Pcs', 'Farmasi Pediatri', 10, 70000, 115000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Sunscreen Baby', 'KFA-PROD-0256', 'Pcs', 'Farmasi Pediatri', 10, 0, 100000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Sunscreen Baby BHP', 'KFA-PROD-0257', 'Pcs', 'Farmasi Pediatri', 10, 0, 0, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Synchro Anti iritasi', 'KFA-PROD-0258', 'Pcs', 'Farmasi Pediatri', 10, 2930, 5500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Syringe Reeme Pro Injector 3 CC (Berpu)', 'KFA-PROD-0259', 'Pcs', 'Farmasi Pediatri', 10, 5500, 9075, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Syringe Reeme Pro Injector 5 CC (Berpu)', 'KFA-PROD-0260', 'Pcs', 'Farmasi Pediatri', 10, 5500, 9075, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Syuga Olive Oil', 'KFA-PROD-0261', 'Pcs', 'Farmasi Pediatri', 10, 408, 500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Syuga Olive Oil', 'KFA-PROD-0262', 'Pcs', 'Farmasi Pediatri', 10, 408, 500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Syuga Olive Oil SAMPLE', 'KFA-PROD-0263', 'Pcs', 'Farmasi Pediatri', 10, 408, 500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Tabir surya', 'KFA-PROD-0264', 'Strip', 'Tablet & Kaplet', 10, 65000, 65000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Teen Derm Gel Sensitive 250 ml', 'KFA-PROD-0265', 'Tube', 'Krim & Salep', 10, 295260, 350000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'TEMPRA DROP SAMPEL', 'KFA-PROD-0266', 'Botol', 'Drop Pediatri', 10, 44000, 70400, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Tempra Drops Grape 15 ml', 'KFA-PROD-0267', 'Botol', 'Drop Pediatri', 10, 47730, 53000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Tempra Drops Strawberry 15 ml', 'KFA-PROD-0268', 'Botol', 'Drop Pediatri', 10, 47424, 53000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Tempra Forte Grape Syrup 30 ml', 'KFA-PROD-0269', 'Botol', 'Sirup & Suspensi', 10, 25530, 28000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Tempra Syrup Grape 30 ml', 'KFA-PROD-0270', 'Botol', 'Sirup & Suspensi', 10, 23088, 25500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Teorol Drop', 'KFA-PROD-0271', 'Botol', 'Drop Pediatri', 10, 137640, 151500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Termisil Cr', 'KFA-PROD-0272', 'Pcs', 'Farmasi Pediatri', 10, 5883, 6500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Thrombogel 10 gram', 'KFA-PROD-0273', 'Tube', 'Krim & Salep', 10, 46620, 52000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Tiriz Oral Drops 15 ml', 'KFA-PROD-0274', 'Botol', 'Drop Pediatri', 10, 87690, 96500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Tranex Serum Innoceutics', 'KFA-PROD-0275', 'Pcs', 'Farmasi Pediatri', 10, 100000, 200000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Transpulmin Baby 10 g', 'KFA-PROD-0276', 'Pcs', 'Farmasi Pediatri', 10, 49173, 55000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Transpulmin Baby 20 gr', 'KFA-PROD-0277', 'Pcs', 'Farmasi Pediatri', 10, 73038, 80500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Tremenza tab', 'KFA-PROD-0278', 'Strip', 'Tablet & Kaplet', 10, 1925, 2500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Trilac Injeksi', 'KFA-PROD-0279', 'Vial', 'Injeksi & Cairan', 10, 28000, 31000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Typhim VI', 'KFA-PROD-0280', 'Pcs', 'Farmasi Pediatri', 10, 278610, 300000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Uriage bariderm cica', 'KFA-PROD-0281', 'Pcs', 'Farmasi Pediatri', 10, 5855, 6500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Uriage cica daily', 'KFA-PROD-0282', 'Pcs', 'Farmasi Pediatri', 10, 0, 320000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Urine collector 100 ml', 'KFA-PROD-0283', 'Pcs', 'Farmasi Pediatri', 10, 5000, 5500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'UV Doux Sunscreen Gel SPF 50', 'KFA-PROD-0284', 'Tube', 'Krim & Salep', 10, 249750, 269730, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'UV E SPF 50 FLUID INVISIBLE', 'KFA-PROD-0285', 'Pcs', 'Farmasi Pediatri', 10, 7992, 9000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Vaksin BCG Bayi', 'KFA-PROD-0286', 'Pcs', 'Farmasi Pediatri', 10, 38295, 275000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Vaksin Campak', 'KFA-PROD-0287', 'Pcs', 'Farmasi Pediatri', 10, 20757, 50000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Vaksin HPV Gardasil 9', 'KFA-PROD-0288', 'Pcs', 'Farmasi Pediatri', 10, 1800000, 2275000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Vaksin Polio (BOPV)', 'KFA-PROD-0289', 'Pcs', 'Farmasi Pediatri', 10, 11, 15000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Vaksin Varicella', 'KFA-PROD-0290', 'Pcs', 'Farmasi Pediatri', 10, 474414, 525000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Valisanbe 2 mg Tab', 'KFA-PROD-0291', 'Strip', 'Tablet & Kaplet', 10, 361, 500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Varivax Vaksin', 'KFA-PROD-0292', 'Pcs', 'Farmasi Pediatri', 10, 626342, 676000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'VAXIGRIP TETRA SAMPEL', 'KFA-PROD-0293', 'Pcs', 'Farmasi Pediatri', 10, 700000, 900000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Vaxigrip Tetra Vaksin', 'KFA-PROD-0294', 'Pcs', 'Farmasi Pediatri', 10, 199800, 285000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Vaxneuvance PCV 15', 'KFA-PROD-0295', 'Pcs', 'Farmasi Pediatri', 10, 880202, 1005000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'VCO Nutiver', 'KFA-PROD-0296', 'Pcs', 'Farmasi Pediatri', 10, 204, 500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Vectrine Dry Syrup', 'KFA-PROD-0297', 'Botol', 'Sirup & Suspensi', 10, 75500, 84000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Velutine 2,5 mg Nebule', 'KFA-PROD-0298', 'Pcs', 'Farmasi Pediatri', 10, 12000, 13500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Ventolin Inhaler', 'KFA-PROD-0299', 'Pcs', 'Farmasi Pediatri', 10, 162142, 180000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Ventolin Nebul 2,5 mg', 'KFA-PROD-0300', 'Pcs', 'Farmasi Pediatri', 10, 18262, 20000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Vermoran 500 mg Chewable Tablet', 'KFA-PROD-0301', 'Strip', 'Tablet & Kaplet', 10, 8400, 9300, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Vestein 300 mg Kap', 'KFA-PROD-0302', 'Pcs', 'Farmasi Pediatri', 10, 7548, 8500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Vestein Dry Syrup 60 ml', 'KFA-PROD-0303', 'Botol', 'Sirup & Suspensi', 10, 74370, 82000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Viusid Liquid 100 ml', 'KFA-PROD-0304', 'Pcs', 'Farmasi Pediatri', 10, 333000, 366500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Vometa Drops', 'KFA-PROD-0305', 'Botol', 'Drop Pediatri', 10, 66600, 75000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Vosedon Suspensi', 'KFA-PROD-0306', 'Botol', 'Sirup & Suspensi', 10, 38350, 42500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Vosedon Tablet', 'KFA-PROD-0307', 'Strip', 'Tablet & Kaplet', 10, 2997, 3400, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Xemose shooting balm', 'KFA-PROD-0308', 'Pcs', 'Farmasi Pediatri', 10, 1809, 2000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Xepazym', 'KFA-PROD-0309', 'Pcs', 'Farmasi Pediatri', 10, 3885, 4300, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Xybat FC 30tab', 'KFA-PROD-0310', 'Strip', 'Tablet & Kaplet', 10, 6119, 6800, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Xybat Syrup', 'KFA-PROD-0311', 'Botol', 'Sirup & Suspensi', 10, 95460, 105000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Zamel Drop', 'KFA-PROD-0312', 'Botol', 'Drop Pediatri', 10, 54390, 60000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Zamel Syrup', 'KFA-PROD-0313', 'Botol', 'Sirup & Suspensi', 10, 49000, 54000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Zegavit Kaplet', 'KFA-PROD-0314', 'Strip', 'Tablet & Kaplet', 10, 4717, 5300, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Zelface Cr', 'KFA-PROD-0315', 'Pcs', 'Farmasi Pediatri', 10, 5106, 6000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Zibramax 500 mg Tab', 'KFA-PROD-0316', 'Strip', 'Tablet & Kaplet', 10, 44400, 49000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Zibramax Sirup Kering', 'KFA-PROD-0317', 'Botol', 'Sirup & Suspensi', 10, 127650, 140500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'ZInc Sulfat 20mg', 'KFA-PROD-0318', 'Pcs', 'Farmasi Pediatri', 10, 250, 350, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Zincpro Drops 15ml', 'KFA-PROD-0319', 'Botol', 'Drop Pediatri', 10, 33470, 37500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Zincpro Syrup', 'KFA-PROD-0320', 'Botol', 'Sirup & Suspensi', 10, 39627, 43600, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-PRD-' || LPAD(id::text, 4, '0'), '2028-12-31', 30, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Avaxim 80 Hep A', 'KFA-VAX-0002', 'Vial', 'Vaksinasi Anak', 5, 357605, 440000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-VAX-' || LPAD(id::text, 4, '0'), '2027-12-31', 20, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Engerix Vaksin', 'KFA-VAX-0003', 'Vial', 'Vaksinasi Anak', 5, 145000, 160000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-VAX-' || LPAD(id::text, 4, '0'), '2027-12-31', 20, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Havrix 720 Junior', 'KFA-VAX-0004', 'Vial', 'Vaksinasi Anak', 5, 431256, 505000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-VAX-' || LPAD(id::text, 4, '0'), '2027-12-31', 20, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Healive Vaksin', 'KFA-VAX-0005', 'Vial', 'Vaksinasi Anak', 5, 255300, 420000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-VAX-' || LPAD(id::text, 4, '0'), '2027-12-31', 20, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Hexaxim Vaksin', 'KFA-VAX-0006', 'Vial', 'Vaksinasi Anak', 5, 880785, 960000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-VAX-' || LPAD(id::text, 4, '0'), '2027-12-31', 20, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'IMOJEV Vaksin', 'KFA-VAX-0007', 'Vial', 'Vaksinasi Anak', 5, 444000, 500000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-VAX-' || LPAD(id::text, 4, '0'), '2027-12-31', 20, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Infarix Hexa', 'KFA-VAX-0008', 'Vial', 'Vaksinasi Anak', 5, 777000, 910000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-VAX-' || LPAD(id::text, 4, '0'), '2027-12-31', 20, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Inlive Vaksin', 'KFA-VAX-0009', 'Vial', 'Vaksinasi Anak', 5, 876900, 950000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-VAX-' || LPAD(id::text, 4, '0'), '2027-12-31', 20, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'MenQuadfi Vaksin', 'KFA-VAX-0010', 'Vial', 'Vaksinasi Anak', 5, 703000, 760000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-VAX-' || LPAD(id::text, 4, '0'), '2027-12-31', 20, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'MMR-II VAKSIN', 'KFA-VAX-0011', 'Vial', 'Vaksinasi Anak', 5, 449550, 550000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-VAX-' || LPAD(id::text, 4, '0'), '2027-12-31', 20, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'PCV 20 Vaksin', 'KFA-VAX-0012', 'Vial', 'Vaksinasi Anak', 5, 945000, 1050000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-VAX-' || LPAD(id::text, 4, '0'), '2027-12-31', 20, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Prevenar 20 Vaksin', 'KFA-VAX-0013', 'Vial', 'Vaksinasi Anak', 5, 950000, 1050000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-VAX-' || LPAD(id::text, 4, '0'), '2027-12-31', 20, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Priorix', 'KFA-VAX-0014', 'Vial', 'Vaksinasi Anak', 5, 464445, 577000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-VAX-' || LPAD(id::text, 4, '0'), '2027-12-31', 20, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'PROQUAD VAKSIN', 'KFA-VAX-0015', 'Vial', 'Vaksinasi Anak', 5, 1176178, 1270000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-VAX-' || LPAD(id::text, 4, '0'), '2027-12-31', 20, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'QDenga Vaksin', 'KFA-VAX-0016', 'Vial', 'Vaksinasi Anak', 5, 468007, 605000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-VAX-' || LPAD(id::text, 4, '0'), '2027-12-31', 20, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Rotarix Vaksin', 'KFA-VAX-0017', 'Vial', 'Vaksinasi Anak', 5, 386414, 463000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-VAX-' || LPAD(id::text, 4, '0'), '2027-12-31', 20, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Rotateq Vaksin', 'KFA-VAX-0018', 'Vial', 'Vaksinasi Anak', 5, 300307, 365000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-VAX-' || LPAD(id::text, 4, '0'), '2027-12-31', 20, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Typhim VI', 'KFA-VAX-0019', 'Vial', 'Vaksinasi Anak', 5, 278610, 300000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-VAX-' || LPAD(id::text, 4, '0'), '2027-12-31', 20, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Vaksin BCG Bayi', 'KFA-VAX-0020', 'Vial', 'Vaksinasi Anak', 5, 38295, 275000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-VAX-' || LPAD(id::text, 4, '0'), '2027-12-31', 20, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Vaksin Campak', 'KFA-VAX-0021', 'Vial', 'Vaksinasi Anak', 5, 20757, 50000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-VAX-' || LPAD(id::text, 4, '0'), '2027-12-31', 20, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Vaksin Polio (BOPV)', 'KFA-VAX-0022', 'Vial', 'Vaksinasi Anak', 5, 11, 15000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-VAX-' || LPAD(id::text, 4, '0'), '2027-12-31', 20, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Vaksin Varicella', 'KFA-VAX-0023', 'Vial', 'Vaksinasi Anak', 5, 474414, 525000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-VAX-' || LPAD(id::text, 4, '0'), '2027-12-31', 20, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Varivax Vaksin', 'KFA-VAX-0024', 'Vial', 'Vaksinasi Anak', 5, 626342, 676000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-VAX-' || LPAD(id::text, 4, '0'), '2027-12-31', 20, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Vaxigrip Tetra Vaksin', 'KFA-VAX-0025', 'Vial', 'Vaksinasi Anak', 5, 199800, 285000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-VAX-' || LPAD(id::text, 4, '0'), '2027-12-31', 20, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    

    WITH d AS (
        INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, min_stock_alert, buy_price, sell_price, is_active, created_at, updated_at)
        VALUES ('KIDINEST', 'Vaxneuvance PCV15', 'KFA-VAX-0026', 'Vial', 'Vaksinasi Anak', 5, 880202, 1005000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (name, clinic_unit) DO UPDATE SET
            unit = EXCLUDED.unit,
            category = EXCLUDED.category,
            buy_price = EXCLUDED.buy_price,
            sell_price = EXCLUDED.sell_price,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id
    )
    INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at)
    SELECT id, 'BCH-VAX-' || LPAD(id::text, 4, '0'), '2027-12-31', 20, 0, CURRENT_TIMESTAMP
    FROM d
    WHERE NOT EXISTS (SELECT 1 FROM drug_batches WHERE drug_batches.drug_id = d.id);
    
COMMIT;