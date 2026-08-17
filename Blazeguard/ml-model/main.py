import os
import shutil
import tempfile
import joblib
from fastapi import FastAPI, UploadFile, File, Form, HTTPException
from pydantic import BaseModel, Field
from typing import Optional, List

from detection import detect_fire
from severity import compute_severity
from spread import predict_spread
from pipeline import analyze_fire

app = FastAPI(
    title="Blaze-Guard Machine Learning Service",
    description="FastAPI wrappers for YOLOv8 fire detection, NASA FIRMS severity classification, and Random Forest weather risk predictions.",
    version="1.1.0"
)

# Load weather risk model trained via train_risk_model.py
risk_model_path = "models/risk_model.pkl"
if os.path.exists(risk_model_path):
    risk_model = joblib.load(risk_model_path)
else:
    risk_model = None
    print(f"Warning: Risk model not found at {risk_model_path}. Run train_risk_model.py first.")

# Request schema for severity calculation
class SeverityRequest(BaseModel):
    brightness: float = Field(..., description="Fire brightness temperature (K)")
    bright_t31: float = Field(..., description="Channel 31 brightness temperature (K)")
    scan: float = Field(..., description="Along-scan pixel size")
    track: float = Field(..., description="Along-track pixel size")
    latitude: float = Field(..., description="Fire center latitude")
    longitude: float = Field(..., description="Fire center longitude")
    confidence: float = Field(..., description="Detection confidence score (0 to 1)")
    wind_speed: float = Field(18.0, description="Wind speed in km/h or m/s")
    wind_deg: float = Field(180.0, description="Wind direction in degrees (0 to 360)")
    vegetation_density: float = Field(0.7, description="Vegetation fuel load density (0 to 1)")

# Request schema for weather risk index prediction
class RiskRequest(BaseModel):
    temperature: float = Field(..., description="Temperature in Celsius")
    humidity: float = Field(..., description="Relative humidity %")
    wind_speed: float = Field(..., description="Wind speed in km/h or m/s")


@app.get("/health")
def health_check():
    return {
        "status": "healthy",
        "service": "blazeguard-ml-model",
        "risk_model_loaded": risk_model is not None
    }


@app.post("/detect")
async def detect(
    file: UploadFile = File(...),
    latitude: float = Form(28.6139),
    longitude: float = Form(77.2090)
):
    """
    Upload an image to run the YOLO detection model.
    """
    if not file.content_type.startswith("image/"):
        raise HTTPException(status_code=400, detail="Uploaded file must be an image.")

    # Save to a temporary file
    temp_dir = tempfile.gettempdir()
    temp_file_path = os.path.join(temp_dir, file.filename)
    
    try:
        with open(temp_file_path, "wb") as buffer:
            shutil.copyfileobj(file.file, buffer)
        
        result = detect_fire(temp_file_path, latitude=latitude, longitude=longitude)
        
        if result is None:
            return {"fire_detected": False, "message": "No fire detected in the image."}
        
        return {
            "fire_detected": True,
            "metrics": result
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Inference error: {str(e)}")
    finally:
        if os.path.exists(temp_file_path):
            os.remove(temp_file_path)


@app.post("/predict-severity")
def predict_severity(payload: SeverityRequest):
    """
    Directly run the tabular severity and spread model using NASA FIRMS features.
    This is used for incoming JSON events (e.g. from the NASA FIRMS live feed).
    """
    features = [
        payload.brightness,
        payload.bright_t31,
        payload.scan,
        payload.track,
        payload.latitude,
        payload.longitude,
        payload.confidence
    ]
    
    try:
        severity = compute_severity(features)
        radius, corridors = predict_spread(
            severity,
            wind_speed=payload.wind_speed,
            wind_deg=payload.wind_deg,
            vegetation_density=payload.vegetation_density
        )
        
        return {
            "severity": severity,
            "radius_meters": radius,
            "corridors": corridors,
            "location": [payload.latitude, payload.longitude]
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Severity prediction error: {str(e)}")


@app.post("/predict-risk")
def predict_risk(payload: RiskRequest):
    """
    Predict weather-based fire risk score (0.0 to 1.0) using the trained Random Forest model.
    """
    if risk_model is None:
        raise HTTPException(status_code=503, detail="Risk model not loaded on server.")
    
    try:
        X = [[payload.temperature, payload.humidity, payload.wind_speed]]
        risk_score = float(risk_model.predict(X)[0])
        return {
            "risk_score": risk_score,
            "level": "High" if risk_score > 0.7 else "Medium" if risk_score > 0.4 else "Low"
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Risk prediction error: {str(e)}")


@app.post("/analyze")
async def analyze(
    file: UploadFile = File(...),
    latitude: float = Form(28.6139),
    longitude: float = Form(77.2090),
    wind_speed: float = Form(18.0),
    wind_deg: float = Form(180.0),
    vegetation_density: float = Form(0.7)
):
    """
    Run end-to-end image-based fire detection, severity classification, and spread radius calculation.
    """
    if not file.content_type.startswith("image/"):
        raise HTTPException(status_code=400, detail="Uploaded file must be an image.")

    temp_dir = tempfile.gettempdir()
    temp_file_path = os.path.join(temp_dir, file.filename)
    
    try:
        with open(temp_file_path, "wb") as buffer:
            shutil.copyfileobj(file.file, buffer)
        
        result = analyze_fire(
            temp_file_path,
            latitude=latitude,
            longitude=longitude,
            wind_speed=wind_speed,
            wind_deg=wind_deg,
            vegetation_density=vegetation_density
        )
        
        return result
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Analysis pipeline error: {str(e)}")
    finally:
        if os.path.exists(temp_file_path):
            os.remove(temp_file_path)


if __name__ == "__main__":
    import uvicorn
    uvicorn.run("main:app", host="0.0.0.0", port=9000, reload=True)
