CREATE INDEX "weather.weather_data_latitude_longitude_idx"
ON "weather"."weather_data"(latitude, longitude);
