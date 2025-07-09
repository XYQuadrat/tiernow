-- Create tierlist table
CREATE TABLE tierlist (
    uuid TEXT PRIMARY KEY,
    name TEXT NOT NULL
);

-- Create tier table
CREATE TABLE tier (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tierlist_uuid TEXT NOT NULL,
    name TEXT NOT NULL,
    "order" INTEGER NOT NULL,
    FOREIGN KEY (tierlist_uuid) REFERENCES tierlist(uuid) ON DELETE CASCADE,
    UNIQUE (tierlist_uuid, "order") -- Unique order per tierlist
);

-- Create entry table
CREATE TABLE entry (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tierlist_uuid TEXT NOT NULL,
    tier_id INTEGER,
    file_key TEXT NOT NULL,
    "order" INTEGER NOT NULL,
    FOREIGN KEY (tierlist_uuid) REFERENCES tierlist(uuid) ON DELETE CASCADE,
    UNIQUE (tier_id, "order") -- Unique order per tier
);
