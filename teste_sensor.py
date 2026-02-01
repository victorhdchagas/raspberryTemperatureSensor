import time
import board
import adafruit_dht

dht_device = adafruit_dht.DHT11(board.D4)

while True:
   try:
      temperature_c = dht_device.temperature
      humidity = dht_device.humidity
      print(f"Temp: {temperature_c:.1f} C   Umidade:{humidity}%")
   except RuntimeError as error:
      print(error.args[0])
      time.sleep(2.0)
      continue
   except Exception as error:
      dht_device.exit()
      raise error
   time.sleep(2.0)
