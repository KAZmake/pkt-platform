-- Seed: users — зеркало тестовых пользователей Keycloak
-- keycloak_id совпадает с sub в JWT (для тестов используем фиксированные UUID)

INSERT INTO users (id, keycloak_id, email, role, first_name, last_name, phone) VALUES
  ('00000000-0000-0000-0000-000000000001', 'kc-borrower-001', 'borrower@pkt.test',  'borrower', 'Тестовый',   'Заёмщик',        '+7 701 111 0001'),
  ('00000000-0000-0000-0000-000000000002', 'kc-borrower-002', 'borrower2@pkt.test', 'borrower', 'Асет',       'Жаксыбеков',     '+7 702 222 0002'),
  ('00000000-0000-0000-0000-000000000003', 'kc-borrower-003', 'borrower3@pkt.test', 'borrower', 'Нурлан',     'Сейткали',       '+7 705 333 0003'),
  ('00000000-0000-0000-0000-000000000010', 'kc-employee-001', 'employee@pkt.test',  'employee', 'Тестовый',   'Сотрудник',      '+7 701 555 0010'),
  ('00000000-0000-0000-0000-000000000011', 'kc-employee-002', 'employee2@pkt.test', 'employee', 'Гульнара',   'Имангалиева',    '+7 707 555 0011'),
  ('00000000-0000-0000-0000-000000000020', 'kc-expert-001',   'expert@pkt.test',    'expert',   'Тестовый',   'Эксперт',        '+7 701 777 0020'),
  ('00000000-0000-0000-0000-000000000021', 'kc-expert-002',   'expert2@pkt.test',   'expert',   'Бауыржан',   'Ахметов',        '+7 778 777 0021'),
  ('00000000-0000-0000-0000-000000000030', 'kc-admin-001',    'admin@pkt.test',     'admin',    'Тестовый',   'Администратор',  '+7 701 999 0030')
ON CONFLICT (keycloak_id) DO NOTHING;
