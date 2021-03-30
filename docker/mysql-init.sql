DROP SCHEMA IF EXISTS hybs_sample_game;
CREATE SCHEMA hybs_sample_game;
USE hybs_sample_game;

DROP TABLE IF EXISTS player_example_data;

CREATE TABLE player_example_data
(
    hayabusa_id    VARCHAR(40),
    int_field      INT,
    string_field   VARCHAR(40)
);

CREATE INDEX idx_player ON player_example_data (hayabusa_id);