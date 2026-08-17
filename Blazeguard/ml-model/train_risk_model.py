import numpy as np
import pandas as pd
from sklearn.ensemble import RandomForestRegressor
import joblib
import os

# Ensure the models directory exists
os.makedirs("models", exist_ok=True)

# Generate synthetic dataset matching fire risk conditions (10,000 observations)
np.random.seed(42)
n_samples = 10000

temperature = np.random.uniform(0, 45, n_samples)       # Temperature in Celsius
humidity = np.random.uniform(10, 100, n_samples)        # Relative humidity %
wind_speed = np.random.uniform(0, 50, n_samples)        # Wind speed in km/h

# Custom fire weather risk heuristic (representing FWI)
risk_score = (1.0 - (humidity / 100.0)) * 0.45 + (temperature / 45.0) * 0.35 + (wind_speed / 50.0) * 0.20
# Add Gaussian noise
risk_score += np.random.normal(0, 0.04, n_samples)
risk_score = np.clip(risk_score, 0.0, 1.0)

df = pd.DataFrame({
    'temperature': temperature,
    'humidity': humidity,
    'wind_speed': wind_speed,
    'risk_score': risk_score
})

X = df[['temperature', 'humidity', 'wind_speed']]
y = df['risk_score']

print("Training Fire Weather Risk Regressor model...")
model = RandomForestRegressor(n_estimators=50, random_state=42)
model.fit(X, y)

joblib.dump(model, "models/risk_model.pkl")
print("Model successfully saved to models/risk_model.pkl")
