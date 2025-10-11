-- Create waiting_list table for car service queue management
-- This table manages the queue system where customers take a number and wait for service

CREATE TABLE waiting_list (
    id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    created_at DATETIME2 NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME2 NOT NULL DEFAULT GETDATE(),
    deleted_at DATETIME2 NULL,

    queue_number INT NOT NULL,
    vehicle_id UNIQUEIDENTIFIER NOT NULL,
    customer_id UNIQUEIDENTIFIER NOT NULL,
    service_date DATETIME2 NOT NULL,
    service_type VARCHAR(100) NOT NULL,
    estimated_time INT NULL, -- estimated service time in minutes
    status VARCHAR(30) NOT NULL DEFAULT 'waiting',
    called_at DATETIME2 NULL,
    service_start_at DATETIME2 NULL,
    service_end_at DATETIME2 NULL,
    notes TEXT NULL,

    -- Foreign key constraints
    CONSTRAINT FK_waiting_list_vehicle FOREIGN KEY (vehicle_id) REFERENCES vehicles(id),
    CONSTRAINT FK_waiting_list_customer FOREIGN KEY (customer_id) REFERENCES users(id),

    -- Unique constraint to ensure queue number is unique per service date
    CONSTRAINT UQ_waiting_list_queue_date UNIQUE (queue_number, service_date),

    -- Check constraint for valid status values
    CONSTRAINT CHK_waiting_list_status CHECK (
        status IN ('waiting', 'called', 'in_service', 'completed', 'canceled', 'no_show')
    )
);

-- Create indexes for better query performance
CREATE INDEX IX_waiting_list_customer_id ON waiting_list(customer_id);
CREATE INDEX IX_waiting_list_vehicle_id ON waiting_list(vehicle_id);
CREATE INDEX IX_waiting_list_service_date ON waiting_list(service_date);
CREATE INDEX IX_waiting_list_status ON waiting_list(status);
CREATE INDEX IX_waiting_list_deleted_at ON waiting_list(deleted_at);

-- Create a trigger to update the updated_at column
CREATE TRIGGER TR_waiting_list_updated_at
ON waiting_list
AFTER UPDATE
AS
BEGIN
    SET NOCOUNT ON;
    UPDATE waiting_list
    SET updated_at = GETDATE()
    FROM waiting_list wl
    INNER JOIN inserted i ON wl.id = i.id;
END;
GO

-- Insert sample data (optional, for testing)
-- Uncomment if you want to insert test data
/*
-- Assuming you have existing users and vehicles
DECLARE @customerId UNIQUEIDENTIFIER = (SELECT TOP 1 id FROM users WHERE role = 'customer');
DECLARE @vehicleId UNIQUEIDENTIFIER = (SELECT TOP 1 id FROM vehicles);

IF @customerId IS NOT NULL AND @vehicleId IS NOT NULL
BEGIN
    INSERT INTO waiting_list (queue_number, vehicle_id, customer_id, service_date, service_type, estimated_time, status, notes)
    VALUES
        (1, @vehicleId, @customerId, CAST(GETDATE() AS DATE), 'Oil Change', 30, 'waiting', 'Regular oil change service'),
        (2, @vehicleId, @customerId, CAST(GETDATE() AS DATE), 'Tire Rotation', 45, 'waiting', 'Standard tire rotation');
END;
*/
