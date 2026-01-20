-- name: GetTierlist :one
SELECT
    tierlist.uuid,
    tierlist.name,
    json_group_array(
        json_object(
            'id', tier.id,
            'name', tier.name,
            'order', tier."order",
            'entries', (
                SELECT
                    json_group_array(
                        json_object(
                            'id', entry.id,
                            'file_key', entry.file_key,
                            'order', entry."order"
                        )
                    )
                FROM entry
                WHERE entry.tier_id = tier.id
                ORDER BY entry."order"
            )
        )
    ) AS tiers,
    (
        SELECT
        json_group_array(
            json_object(
                'id', entry.id,
                'file_key', entry.file_key,
                'order', entry."order"
            )
        )
        FROM entry
        WHERE entry.tier_id IS NULL AND entry.tierlist_uuid = ?1
        ORDER BY entry."order"
    )
    AS unassigned_entries
FROM tierlist
LEFT JOIN tier ON tierlist.uuid = tier.tierlist_uuid
WHERE tierlist.uuid = ?1
GROUP BY tierlist.uuid;

-- name: CreateTierlist :one
INSERT INTO tierlist (
    uuid,
    name
) VALUES (
    ?,
    ?
)
RETURNING *;

-- name: CreateTier :one
INSERT INTO tier (
    tierlist_uuid,
    name,
    "order"
) VALUES (
    ?,
    ?,
    ?
)
RETURNING *;

-- name: UploadImageMetadata :one
INSERT INTO entry (
    tierlist_uuid,
    file_key,
    "order"
)
SELECT ?1, ?2, IFNULL(MAX("order"), 0) + 1
FROM entry 
WHERE entry.tierlist_uuid = ?1
RETURNING *;

-- name: SetImageTier :one
UPDATE entry
SET tier_id = ?1
WHERE id = ?2
RETURNING *