CREATE TABLE IF NOT EXISTS teams (
                                     team_id SERIAL PRIMARY KEY,
                                     team_name TEXT NOT NULL,
                                     team_status BOOLEAN
);


-- Инициализация таблицы logs
CREATE TABLE IF NOT EXISTS log (
                                           id SERIAL PRIMARY KEY,
                                           date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                           action TEXT NOT NULL,
                                           status TEXT NOT NULL,
                                           details TEXT NOT NULL

);

-- Инициализация таблицы events
CREATE TABLE IF NOT EXISTS events (
                                      event_id SERIAL PRIMARY KEY,
                                      event_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                      event_tournament TEXT,
                                      team_home TEXT,
                                      team_away TEXT,
                                      goals_home INTEGER,
                                      goals_away INTEGER,
                                      pen_home INTEGER,
                                      pen_away INTEGER,
                                      rc_home INTEGER,
                                      rc_away INTEGER,
                                      importance BOOLEAN,
                                      event_status TEXT,
                                      published_status TEXT
);



