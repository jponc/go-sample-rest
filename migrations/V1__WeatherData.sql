CREATE TABLE "weather"."weather_data" (
    "id" character varying NOT NULL,
    "latitude" float NOT NULL,
    "longitude" float NOT NULL,
    "temperature" float NOT NULL,
    "wind_direction" float NOT NULL,
    "wind_speed" float NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT "weather_data_pk" PRIMARY KEY ("id")
);
