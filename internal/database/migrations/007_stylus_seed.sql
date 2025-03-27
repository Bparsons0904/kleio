INSERT OR IGNORE INTO styluses 
(name, manufacturer, expected_lifespan_hours, active, primary_stylus, base_model) 
VALUES
-- Ortofon
('2M Red', 'Ortofon', 1000, 0, 0, true),
('2M Blue', 'Ortofon', 1500, 0, 0, true),
('2M Bronze', 'Ortofon', 2000, 0, 0, true),
('2M Black', 'Ortofon', 2000, 0, 0, true),
('Super OM 10', 'Ortofon', 1000, 0, 0, true),
('Concorde MkII Club', 'Ortofon', 1000, 0, 0, true),
-- Audio-Technica
('AT-VM95E', 'Audio-Technica', 800, 0, 0, true),
('AT-VM95EN', 'Audio-Technica', 1000, 0, 0, true),
('AT-VM95ML', 'Audio-Technica', 1500, 0, 0, true),
('AT-VM95SH', 'Audio-Technica', 2000, 0, 0, true),
('AT33PTG/II', 'Audio-Technica', 1500, 0, 0, true),
-- Nagaoka
('MP-110', 'Nagaoka', 1000, 0, 0, true),
('MP-150', 'Nagaoka', 1500, 0, 0, true),
('MP-200', 'Nagaoka', 2000, 0, 0, true),
-- Grado
('Prestige Black3', 'Grado', 1000, 0, 0, true),
('Prestige Gold3', 'Grado', 1500, 0, 0, true),
('Timbre Series Platinum3', 'Grado', 2000, 0, 0, true),
-- Denon
('DL-103', 'Denon', 1500, 0, 0, true),
('DL-110', 'Denon', 1500, 0, 0, true);
