# super-fantasy-api
API service for calculating, retrieving stats, and other items

# Baseball

## Prerequisites

- Have a running mongo on your localhost
- You'll need these env variables:
  - MONGO_URI=mongodb://localhost:27017
  - DB_NAME=super-fantasy
  - COLLECTION_NAME=Baseball

### Start Mongo

Run:

```sh
mongosh
show dbs
db.createCollection("Baseball")
```

## Build and Run

To build run:

```sh
go build -o super-fantasy-api
```

To run:

```sh
./super-fantasy-api 
```

## Example Commands (baseball)

### Upload

Upload csvs in `csv/`

ATC

```sh
curl -X POST http://localhost:8080/api/v1/baseball/upload \
  -F "csv=@ATC-2026-Pitcher-Projections.csv" \
  -F "settings={\"source\": \"fangraphs\", \"position\": \"pitcher\", \"year\": \"2026\", \"suffix\": \"atc\"}" \
  -H "Content-Type: multipart/form-data"
```

BatX

```sh
curl -X POST http://localhost:8080/api/v1/baseball/upload \
  -F "csv=@BatX-2026-Batter-Projections.csv" \
  -F "settings={\"source\": \"fangraphs\", \"position\": \"batter\", \"year\": \"2026\", \"suffix\": \"batx\"}" \
  -H "Content-Type: multipart/form-data"
```

Steamer (batter)

```sh
curl -X POST http://localhost:8080/api/v1/baseball/upload \
  -F "csv=@Steamer-2026-Batter-Projections.csv" \
  -F "settings={\"source\": \"fangraphs\", \"position\": \"batter\", \"year\": \"2026\", \"suffix\": \"steamer\"}" \
  -H "Content-Type: multipart/form-data"
```

Steamer (pitcher)

```sh
curl -X POST http://localhost:8080/api/v1/baseball/upload \
  -F "csv=@Steamer-2026-Pitcher-Projections.csv" \
  -F "settings={\"source\": \"fangraphs\", \"position\": \"pitcher\", \"year\": \"2026\", \"suffix\": \"steamer\"}" \
  -H "Content-Type: multipart/form-data"
```

FantasyPros (batter)

```sh
curl -X POST http://localhost:8080/api/v1/baseball/upload \
  -F "csv=@FantasyPros-2026-Batter-Projections.csv" \
  -F "settings={\"source\": \"fantasypros\", \"position\": \"batter\", \"year\": \"2026\"}" \
  -H "Content-Type: multipart/form-data"
```

FantasyPros (pitcher)

```sh
curl -X POST http://localhost:8080/api/v1/baseball/upload \
  -F "csv=@FantasyPros-2026-Pitcher-Projections.csv" \
  -F "settings={\"source\": \"fantasypros\", \"position\": \"pitcher\", \"year\": \"2026\"}" \
  -H "Content-Type: multipart/form-data"
```

After upload, all data should be in mongodb

### Export

Export data into csv file using league settings and all documents available

```sh
curl -X POST http://localhost:8080/api/v1/baseball/export \
  -F "settings={\"projection_name\": \"aggregate\", \"settings\": {\"batting\": {\"runs_scored\": 1, \"total_bases\": 1, \"runs_batted_in\": 1, \"walks\": 1, \"strikeouts\": -1, \"stolen_bases\": 1, \"hitting_for_cycle\": 15}, \"pitching\": {\"innings_pitched\": 3, \"hits_allowed\": -1, \"earned_runs\": -2, \"walks_issued\": -1, \"strikeouts\": 1, \"no_hitters\": 5, \"perfect_games\": 10, \"wins\": 5, \"losses\": -5, \"saves\": 5, \"holds\": 3}}}" \
  -H "Content-Type: multipart/form-data" \
  -o player_points.csv
```