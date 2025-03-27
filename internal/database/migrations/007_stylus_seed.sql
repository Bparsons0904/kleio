INSERT OR IGNORE INTO styluses 
(name, manufacturer, model_number, expected_lifespan_hours, active, primary_stylus) 
VALUES
-- Ortofon
('2M Red', 'Ortofon', '2M Red', 1000, 0, 0),
('2M Blue', 'Ortofon', '2M Blue', 1500, 0, 0),
('2M Bronze', 'Ortofon', '2M Bronze', 2000, 0, 0),
('2M Black', 'Ortofon', '2M Black', 2000, 0, 0),
('Super OM 10', 'Ortofon', 'Super OM 10', 1000, 0, 0),
('Concorde MkII Club', 'Ortofon', 'Concorde MkII Club', 1000, 0, 0),

-- Audio-Technica
('AT-VM95E', 'Audio-Technica', 'AT-VM95E', 800, 0, 0),
('AT-VM95EN', 'Audio-Technica', 'AT-VM95EN', 1000, 0, 0),
('AT-VM95ML', 'Audio-Technica', 'AT-VM95ML', 1500, 0, 0),
('AT-VM95SH', 'Audio-Technica', 'AT-VM95SH', 2000, 0, 0),
('AT33PTG/II', 'Audio-Technica', 'AT33PTG/II', 1500, 0, 0),

-- Nagaoka
('MP-110', 'Nagaoka', 'MP-110', 1000, 0, 0),
('MP-150', 'Nagaoka', 'MP-150', 1500, 0, 0),
('MP-200', 'Nagaoka', 'MP-200', 2000, 0, 0),

-- Grado
('Prestige Black3', 'Grado', 'Prestige Black3', 1000, 0, 0),
('Prestige Gold3', 'Grado', 'Prestige Gold3', 1500, 0, 0),
('Timbre Series Platinum3', 'Grado', 'Timbre Platinum3', 2000, 0, 0),

-- Denon
('DL-103', 'Denon', 'DL-103', 1500, 0, 0),
('DL-110', 'Denon', 'DL-110', 1500, 0, 0);
