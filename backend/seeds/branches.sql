-- Seed Operational Units / Branches for Dual-Brand Clinic:
-- 1. Kidinest: Pediatric & Child Development Center
-- 2. Semarang Skin Center (SSC): Aesthetic Dermatology & Laser Center

INSERT INTO branches (id, code, name, address, phone, is_active, created_at, updated_at)
VALUES 
(
    'a0000000-0000-0000-0000-000000000001',
    'KIDINEST',
    'Kidinest Pediatric & Child Development Center',
    'Jl. Dokter Cipto No. 100, Semarang, Jawa Tengah',
    '024-8888001',
    TRUE,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
),
(
    'a0000000-0000-0000-0000-000000000002',
    'SSC',
    'Semarang Skin Center (Estetika & Dermatologi)',
    'Jl. Erlangga Raya No. 45, Semarang, Jawa Tengah',
    '024-8888002',
    TRUE,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
)
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    address = EXCLUDED.address,
    phone = EXCLUDED.phone;
