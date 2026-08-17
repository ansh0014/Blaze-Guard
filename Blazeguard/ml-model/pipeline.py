from detection import detect_fire
from severity import compute_severity
from spread import predict_spread

def analyze_fire(image_path, latitude=28.6139, longitude=77.2090, wind_speed=18, wind_deg=180, vegetation_density=0.7):

    fire_event = detect_fire(image_path, latitude=latitude, longitude=longitude)

    if fire_event is None:
        return {"message": "No fire detected"}

    features = [
        fire_event["brightness"],
        fire_event["bright_t31"],
        fire_event["scan"],
        fire_event["track"],
        fire_event["latitude"],
        fire_event["longitude"],
        fire_event["confidence"]
    ]
    severity = compute_severity(features)

    radius, corridors = predict_spread(
        severity,
        wind_speed = wind_speed,
        wind_deg = wind_deg,
        vegetation_density = vegetation_density
    )

    return {
        "severity": severity,
        "radius_meters": radius,
        "corridors": corridors,
        "location": fire_event["location"],
        "confidence": fire_event["confidence"]
    }
