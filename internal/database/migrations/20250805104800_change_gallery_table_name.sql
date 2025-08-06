-- +goose Up
ALTER TABLE gallery
RENAME TO galleries;

-- +goose Down
ALTER TABLE galleries
RENAME TO gallery;
