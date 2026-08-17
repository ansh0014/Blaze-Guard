def compute_spread_radius(severity, wind_speed, vegetation_density):
    if severity == "Low":
        base = 300
    elif severity == "Medium":
        base = 800
    else:
        base = 1500

    wind_multiplier = 1 + (wind_speed / 40)
    veg_multiplier = 1 + vegetation_density

    return int(base * wind_multiplier * veg_multiplier)

def predict_spread(severity, wind_speed, wind_deg, vegetation_density):
    radius = compute_spread_radius(severity, wind_speed, vegetation_density)

    # wind_deg is the angle the wind blows FROM.
    # The fire spreads in the direction the wind blows TO: (wind_deg + 180) % 360
    blow_deg = (wind_deg + 180) % 360

    corridors = []
    
    # 0/360 is North, 90 is East, 180 is South, 270 is West
    if 315 <= blow_deg or blow_deg < 45:
        corridors.append("corridor_north")
    elif 45 <= blow_deg < 135:
        corridors.append("corridor_east")
    elif 135 <= blow_deg < 225:
        corridors.append("corridor_south")
    else:
        corridors.append("corridor_west")

    # If wind speed is substantial (> 15 km/h), add adjacent corridors
    if wind_speed > 15:
        if "corridor_north" in corridors:
            corridors.extend(["corridor_east", "corridor_west"])
        elif "corridor_south" in corridors:
            corridors.extend(["corridor_east", "corridor_west"])
        elif "corridor_east" in corridors:
            corridors.extend(["corridor_north", "corridor_south"])
        elif "corridor_west" in corridors:
            corridors.extend(["corridor_north", "corridor_south"])

    # High severity fires spread rapidly in all directions; include a buffer
    if severity == "High":
        for c in ["corridor_north", "corridor_east", "corridor_south", "corridor_west"]:
            if c not in corridors:
                corridors.append(c)

    return radius, list(set(corridors))
