-- Seed: application_history, collaterals, application_collaterals,
--        expert_conclusions, committee_votes

-- Temporarily disable UPDATE rule to allow idempotent INSERT
ALTER TABLE application_history DISABLE RULE application_history_no_update;

INSERT INTO application_history
  (id, application_id, from_status, to_status, actor_id, comment, created_at)
SELECT id, application_id, from_status, to_status, actor_id, comment, created_at
FROM (VALUES
  ('a3000000-0000-0000-0000-000000000001'::uuid,
   'a2000000-0000-0000-0000-000000000002'::uuid, NULL::varchar, 'received',
   '00000000-0000-0000-0000-000000000001'::uuid, 'Заявка подана заёмщиком',
   NOW() - INTERVAL '2 days'),
  ('a3000000-0000-0000-0000-000000000002'::uuid,
   'a2000000-0000-0000-0000-000000000002'::uuid, 'received', 'primary_scoring',
   '00000000-0000-0000-0000-000000000010'::uuid, 'Принята в работу',
   NOW() - INTERVAL '6 hours'),
  ('a3000000-0000-0000-0000-000000000003'::uuid,
   'a2000000-0000-0000-0000-000000000012'::uuid, NULL::varchar, 'received',
   '00000000-0000-0000-0000-000000000003'::uuid, 'Заявка подана',
   NOW() - INTERVAL '45 days'),
  ('a3000000-0000-0000-0000-000000000004'::uuid,
   'a2000000-0000-0000-0000-000000000012'::uuid, 'received', 'primary_scoring',
   '00000000-0000-0000-0000-000000000010'::uuid, NULL::text,
   NOW() - INTERVAL '44 days'),
  ('a3000000-0000-0000-0000-000000000005'::uuid,
   'a2000000-0000-0000-0000-000000000012'::uuid, 'primary_scoring', 'security_check',
   '00000000-0000-0000-0000-000000000010'::uuid, NULL::text,
   NOW() - INTERVAL '42 days'),
  ('a3000000-0000-0000-0000-000000000006'::uuid,
   'a2000000-0000-0000-0000-000000000012'::uuid, 'security_check', 'collateral_expertise',
   '00000000-0000-0000-0000-000000000010'::uuid, NULL::text,
   NOW() - INTERVAL '38 days'),
  ('a3000000-0000-0000-0000-000000000007'::uuid,
   'a2000000-0000-0000-0000-000000000012'::uuid, 'collateral_expertise', 'legal_check',
   '00000000-0000-0000-0000-000000000020'::uuid, 'Залог проверен, оценка подтверждена',
   NOW() - INTERVAL '32 days'),
  ('a3000000-0000-0000-0000-000000000008'::uuid,
   'a2000000-0000-0000-0000-000000000012'::uuid, 'legal_check', 'credit_analysis',
   '00000000-0000-0000-0000-000000000021'::uuid, 'Юр. проверка пройдена',
   NOW() - INTERVAL '28 days'),
  ('a3000000-0000-0000-0000-000000000009'::uuid,
   'a2000000-0000-0000-0000-000000000012'::uuid, 'credit_analysis', 'credit_committee',
   '00000000-0000-0000-0000-000000000011'::uuid, NULL::text,
   NOW() - INTERVAL '20 days'),
  ('a3000000-0000-0000-0000-000000000010'::uuid,
   'a2000000-0000-0000-0000-000000000012'::uuid, 'credit_committee', 'approved',
   '00000000-0000-0000-0000-000000000030'::uuid, 'КК одобрил единогласно',
   NOW() - INTERVAL '18 days'),
  ('a3000000-0000-0000-0000-000000000011'::uuid,
   'a2000000-0000-0000-0000-000000000012'::uuid, 'approved', 'documentation',
   '00000000-0000-0000-0000-000000000010'::uuid, NULL::text,
   NOW() - INTERVAL '15 days'),
  ('a3000000-0000-0000-0000-000000000012'::uuid,
   'a2000000-0000-0000-0000-000000000012'::uuid, 'documentation', 'issued',
   '00000000-0000-0000-0000-000000000010'::uuid, 'Договор подписан, средства перечислены',
   NOW() - INTERVAL '10 days')
) AS v(id, application_id, from_status, to_status, actor_id, comment, created_at)
WHERE NOT EXISTS (SELECT 1 FROM application_history ah WHERE ah.id = v.id);

ALTER TABLE application_history ENABLE RULE application_history_no_update;

INSERT INTO collaterals (id, type, description, estimated_value, cadastral_number,
                         insurance_expiry, last_inventory_date, is_released) VALUES
  ('d2000000-0000-0000-0000-000000000001',
   'land', 'Земельный участок с/х назначения, 300 га, Акжаикский р-н',
   45000000, '027-051-123-456', '2026-12-31', '2024-10-15', FALSE),
  ('d2000000-0000-0000-0000-000000000002',
   'equipment', 'Комбайн John Deere S680, 2021 г.в.',
   18500000, NULL, '2025-06-30', '2024-08-20', FALSE),
  ('d2000000-0000-0000-0000-000000000003',
   'livestock', 'КРС 120 голов, герефордская порода',
   12000000, NULL, '2025-09-30', '2024-09-01', FALSE),
  ('d2000000-0000-0000-0000-000000000004',
   'land', 'Пастбищный участок 280 га, Бурлинский р-н',
   22400000, '027-052-234-100', '2025-12-31', '2024-07-10', TRUE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO application_collaterals (application_id, collateral_id, attached_at, released_at) VALUES
  ('a2000000-0000-0000-0000-000000000004', 'd2000000-0000-0000-0000-000000000001',
   NOW() - INTERVAL '5 days', NULL),
  ('a2000000-0000-0000-0000-000000000007', 'd2000000-0000-0000-0000-000000000002',
   NOW() - INTERVAL '14 days', NULL),
  ('a2000000-0000-0000-0000-000000000007', 'd2000000-0000-0000-0000-000000000003',
   NOW() - INTERVAL '14 days', NULL),
  ('a2000000-0000-0000-0000-000000000012', 'd2000000-0000-0000-0000-000000000004',
   NOW() - INTERVAL '40 days', NOW() - INTERVAL '10 days')
ON CONFLICT (application_id, collateral_id) DO NOTHING;

INSERT INTO expert_conclusions (id, application_id, expert_id, stage, risks,
                                conclusion_text, result) VALUES
  ('ec000000-0000-0000-0000-000000000001',
   'a2000000-0000-0000-0000-000000000012',
   '00000000-0000-0000-0000-000000000020',
   'collateral_expertise',
   '{"land_title_risk": false, "insurance_risk": false, "valuation_risk": false}',
   'Залог соответствует требованиям. Рыночная стоимость подтверждена.',
   'approved'),
  ('ec000000-0000-0000-0000-000000000002',
   'a2000000-0000-0000-0000-000000000012',
   '00000000-0000-0000-0000-000000000021',
   'legal_check',
   '{"encumbrance_risk": false, "ownership_risk": false}',
   'Правоустанавливающие документы в порядке. Обременений не выявлено.',
   'approved'),
  ('ec000000-0000-0000-0000-000000000003',
   'a2000000-0000-0000-0000-000000000012',
   '00000000-0000-0000-0000-000000000011',
   'credit_analysis',
   '{"repayment_risk": false, "market_risk": true, "seasonal_risk": true}',
   'Финансовое состояние удовлетворительное. Сезонные риски приняты во внимание. Рекомендую одобрить.',
   'approved')
ON CONFLICT (id) DO NOTHING;

INSERT INTO committee_votes (id, application_id, expert_id, vote, comment, signed_at) VALUES
  ('d3000000-0000-0000-0000-000000000001',
   'a2000000-0000-0000-0000-000000000012',
   '00000000-0000-0000-0000-000000000020', 'approved',
   'Залог ликвидный, риски приемлемые', NOW() - INTERVAL '18 days'),
  ('d3000000-0000-0000-0000-000000000002',
   'a2000000-0000-0000-0000-000000000012',
   '00000000-0000-0000-0000-000000000021', 'approved',
   NULL, NOW() - INTERVAL '18 days'),
  ('d3000000-0000-0000-0000-000000000003',
   'a2000000-0000-0000-0000-000000000012',
   '00000000-0000-0000-0000-000000000030', 'approved',
   'Одобрено', NOW() - INTERVAL '18 days')
ON CONFLICT (id) DO NOTHING;
