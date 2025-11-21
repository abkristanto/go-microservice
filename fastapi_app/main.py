from datetime import datetime, timezone
from typing import List

from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI()


class Event(BaseModel):
    id: str
    title: str
    description: str
    starts_at: datetime


@app.get("/events", response_model=List[Event])
def get_events():
    now = datetime.now(timezone.utc)
    return [
        Event(
            id="evt_1",
            title="Test Event 1",
            description="First test event",
            starts_at=now,
        ),
        Event(
            id="evt_2",
            title="Test Event 2",
            description="Second test event",
            starts_at=now.replace(hour=(now.hour + 1) % 24),
        ),
    ]
