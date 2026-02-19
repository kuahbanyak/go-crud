-- Migration: Add mechanic_id column to waiting_lists table
-- Description: Adds mechanic_id field to track which mechanic is assigned to service each queue entry
-- Date: 2026-02-07

-- Add mechanic_id column to waiting_lists table
ALTER TABLE waiting_lists
ADD mechanic_id uniqueidentifier NULL;

-- Add foreign key constraint to users table
ALTER TABLE waiting_lists
ADD CONSTRAINT FK_waiting_lists_mechanic
FOREIGN KEY (mechanic_id) REFERENCES users(id)
ON DELETE SET NULL;

-- Add index for better query performance
CREATE INDEX IX_waiting_lists_mechanic_id ON waiting_lists(mechanic_id);

-- Comment: mechanic_id is nullable because not all queues will have an assigned mechanic immediately
-- Mechanics can assign themselves to queues using the /api/v1/mechanic/waiting-list/assign endpoint

