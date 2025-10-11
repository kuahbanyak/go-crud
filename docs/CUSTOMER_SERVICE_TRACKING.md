# Customer Service Progress Tracking

## Overview
Customers can now track the real-time progress of their car service using the waiting list system. This feature provides detailed information about queue position, estimated wait time, and service status.

## API Endpoint

### Track Service Progress
**Endpoint:** `GET /api/v1/waiting-list/{id}/progress`

**Authentication:** Required (Customer must be authenticated)

**Description:** Allows customers to track the real-time progress of their car service

**Parameters:**
- `id` (path parameter) - The UUID of the waiting list ticket

**Response Example:**
```json
{
  "success": true,
  "message": "Service progress retrieved successfully",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "queue_number": 5,
    "status": "waiting",
    "status_message": "⏳ 2 customer(s) ahead of you. Currently serving queue #3",
    "vehicle_brand": "Toyota",
    "vehicle_model": "Camry",
    "license_plate": "ABC-1234",
    "service_type": "Oil Change",
    "service_date": "2025-10-09T00:00:00Z",
    "estimated_time_minutes": 45,
    "queue_position": 5,
    "people_ahead": 2,
    "estimated_wait_minutes": 60,
    "timeline": {
      "queue_taken_at": "2025-10-09T08:30:00Z",
      "called_at": null,
      "service_start_at": null,
      "service_end_at": null
    },
    "notes": "Regular maintenance"
  }
}
```

## Service Status Progression

The service goes through the following statuses:

1. **waiting** - Customer is waiting in queue
   - Status message shows how many people are ahead
   - Estimated wait time is calculated

2. **called** - Customer has been called to the service area
   - Status message: "📢 You've been called! Please proceed to the service area immediately."

3. **in_service** - Vehicle is currently being serviced
   - Status message: "🔧 Your vehicle is currently being serviced. Please wait in the customer lounge."
   - Estimated wait time becomes 0

4. **completed** - Service is finished
   - Status message: "✅ Your service has been completed! Thank you for choosing our service."
   - Timeline shows the complete service duration

5. **canceled** - Service was canceled
6. **no_show** - Customer didn't show up

## How to Use

### Step 1: Take a Queue Number
First, the customer needs to take a queue number:
```bash
POST /api/v1/waiting-list/take
Authorization: Bearer <token>

{
  "vehicle_id": "your-vehicle-uuid",
  "service_type": "Oil Change",
  "service_date": "2025-10-09",
  "estimated_time": 45,
  "notes": "Regular maintenance"
}
```

### Step 2: Get Your Ticket ID
The response will include your ticket ID:
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "queue_number": 5,
  ...
}
```

### Step 3: Track Progress
Use the ticket ID to track progress:
```bash
GET /api/v1/waiting-list/{id}/progress
Authorization: Bearer <token>
```

### Alternative: View All Your Tickets
You can also view all your queue entries:
```bash
GET /api/v1/waiting-list/my-queue
Authorization: Bearer <token>
```

## Real-Time Updates

For a real-time tracking experience, your frontend application should:

1. **Poll the endpoint** every 30-60 seconds to get updated progress
2. **Display a progress bar** showing:
   - Current queue position
   - Number of people ahead
   - Estimated wait time
   - Current status with friendly messages

3. **Show notifications** when status changes:
   - When you're next in line
   - When you've been called
   - When service starts
   - When service is completed

## Example Frontend Implementation

### JavaScript/React Example
```javascript
const trackServiceProgress = async (ticketId) => {
  try {
    const response = await fetch(`/api/v1/waiting-list/${ticketId}/progress`, {
      headers: {
        'Authorization': `Bearer ${authToken}`
      }
    });
    
    const data = await response.json();
    
    // Update UI with progress information
    console.log('Queue Position:', data.data.queue_position);
    console.log('People Ahead:', data.data.people_ahead);
    console.log('Status:', data.data.status_message);
    console.log('Estimated Wait:', data.data.estimated_wait_minutes, 'minutes');
    
    return data.data;
  } catch (error) {
    console.error('Failed to fetch progress:', error);
  }
};

// Poll every 30 seconds
setInterval(() => {
  trackServiceProgress('your-ticket-id');
}, 30000);
```

## Key Features

✅ **Real-time Queue Position** - Know exactly where you are in line
✅ **Estimated Wait Time** - Get accurate wait time estimates (30 min per customer)
✅ **Friendly Status Messages** - Easy-to-understand progress updates with emojis
✅ **Service Timeline** - Track when you were called, when service started, and when it completed
✅ **Vehicle Information** - See which vehicle is being serviced
✅ **Security** - Customers can only view their own service tickets

## Error Handling

- **401 Unauthorized**: Customer is not authenticated
- **403 Forbidden**: Trying to view another customer's ticket
- **404 Not Found**: Ticket ID doesn't exist
- **400 Bad Request**: Invalid ticket ID format

## Tips for Best User Experience

1. **Refresh regularly** - The queue status updates as admin staff progress through services
2. **Enable notifications** - Get alerted when your status changes

