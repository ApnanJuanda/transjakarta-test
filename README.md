1. Jalankan perintah docker-compose up --build 
2. Check pada docker desktop di container mqtt, pilih files -> mosquitto -> config -> edit file mosquitto.conf
   tambahkan konfigurasi berikut:
   listener 1883
   listener 9001
   protocol websockets
   allow_anonymous true
3. Buka cmd, lalu jalankan docker exec -it transjakarta-postgres psql -U postgres -d db_transjakarta
4. Salin dan Jalankan query dari file 20250812001_create_table_vehicle_locations.up.sql
5. Jika container backend masih belum bisa berjalan, hentikan semua Docker container yang sedang aktif.
6. Execute kembali docker-compose up --build