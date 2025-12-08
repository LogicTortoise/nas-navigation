from fastapi import FastAPI
from datetime import datetime

app = FastAPI(title="Test Server")

start_time = datetime.now()

@app.get("/")
def root():
    return {
        "message": "Hello from Test Server!",
        "start_time": start_time.isoformat(),
        "uptime_seconds": (datetime.now() - start_time).total_seconds()
    }

@app.get("/health")
def health():
    return {"status": "healthy", "timestamp": datetime.now().isoformat()}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8899)
