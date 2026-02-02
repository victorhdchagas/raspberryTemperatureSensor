import sys
import board
import adafruit_dht
import json

def get_reading(pin_number):
    # Map string "17" to board.D17
    pin_attr = f"D{pin_number}"
    if not hasattr(board, pin_attr):
        return {"temperature": None, "humidity": None, "error": f"Invalid pin: {pin_number}"}
    
    pin = getattr(board, pin_attr)
    dht_device = adafruit_dht.DHT11(pin)
    
    try:
        temperature_c = dht_device.temperature
        humidity = dht_device.humidity
        if temperature_c is not None and humidity is not None:
            return {
                "temperature": float(temperature_c),
                "humidity": float(humidity),
                "error": None
            }
        else:
            return {"temperature": None, "humidity": None, "error": "Read timeout or empty result"}
    except RuntimeError as error:
        # Errors happen fairly often with DHT sensors, just report it
        return {"temperature": None, "humidity": None, "error": str(error)}
    except Exception as error:
        return {"temperature": None, "humidity": None, "error": str(error)}
    finally:
        dht_device.exit()

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print(json.dumps({"temperature": None, "humidity": None, "error": "Missing GPIO pin argument"}))
        sys.exit(1)
    
    pin_arg = sys.argv[1]
    result = get_reading(pin_arg)
    print(json.dumps(result))
