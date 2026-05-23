-- Seed: borrowers и loan_programs

INSERT INTO borrowers (id, user_id, inn, bin, org_name, activity_type, farm_id) VALUES
  ('b1000000-0000-0000-0000-000000000001',
   '00000000-0000-0000-0000-000000000001', '850312300123', NULL,
   NULL, 'crop_farming', 'f1000000-0000-0000-0000-000000000001'),
  ('b1000000-0000-0000-0000-000000000002',
   '00000000-0000-0000-0000-000000000002', '780605400456', NULL,
   'КХ "Жаксыбеков А."', 'crop_farming', 'f1000000-0000-0000-0000-000000000001'),
  ('b1000000-0000-0000-0000-000000000003',
   '00000000-0000-0000-0000-000000000003', '900120500789', '170340012345',
   'КХ "Сейткали Н."', 'livestock', 'f1000000-0000-0000-0000-000000000002')
ON CONFLICT (inn) DO NOTHING;

INSERT INTO loan_programs (id, name, name_kz, name_en, rate, min_amount, max_amount,
                           min_term_months, max_term_months, activity_types, is_active) VALUES
  ('a1000000-0000-0000-0000-000000000001',
   'Весенне-полевые работы', 'Көктемгі егіс жұмыстары', 'Spring Field Works',
   9.00, 500000, 50000000, 6, 36, ARRAY['crop_farming','mixed'], TRUE),
  ('a1000000-0000-0000-0000-000000000002',
   'Развитие животноводства', 'Мал шаруашылығын дамыту', 'Livestock Development',
   8.50, 1000000, 100000000, 12, 60, ARRAY['livestock','mixed'], TRUE),
  ('a1000000-0000-0000-0000-000000000003',
   'Техническое переоснащение', 'Техникалық жаңарту', 'Technical Re-equipment',
   10.00, 2000000, 150000000, 12, 84, ARRAY['crop_farming','livestock','mixed'], TRUE),
  ('a1000000-0000-0000-0000-000000000004',
   'Малый агробизнес', 'Шағын агробизнес', 'Small Agribusiness',
   11.00, 200000, 10000000, 3, 24, ARRAY['crop_farming','livestock','mixed'], FALSE)
ON CONFLICT (id) DO NOTHING;
