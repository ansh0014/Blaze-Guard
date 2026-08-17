from ultralytics import YOLO

model = YOLO("models/best.pt")

def detect_fire(img_path, latitude=28.6139, longitude=77.2090):

    results = model(img_path)

    for r in results:
        for box in r.boxes:

            x1, y1, x2, y2 = box.xyxy[0].tolist()
            confidence = float(box.conf[0])

            area = (x2 - x1) * (y2 - y1)

            brightness = 300 + (confidence * 70)
            bright_t31 = 290 + (confidence * 40)
            scan = min(1.0, area / 100000)
            track = min(1.0, area / 120000)

            return {
                "brightness": brightness,
                "bright_t31": bright_t31,
                "scan": scan,
                "track": track,
                "latitude": latitude,
                "longitude": longitude,
                "confidence": confidence,
                "location": (latitude, longitude)
            }

    return None
