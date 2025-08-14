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
7. Apabila backend sudah berhasil running, hit API Start Publish Data, API ini akan mengirim data secara terus-menerus dengan jeda 2 detik dan akan berhenti jika radius antara latitude dan longitude dengan titik destination sudah 0 atau bisa juga dihentikan dengan hit API Stop Publish Data 
8. Ketika Start Publish Data berjalan maka proses yang akan dilakukan meliputi
   pengiriman data ke MQTT, menerima data dari MQTT, menyimpan data ke PostgreSQL dan apabila radius sudah berada dalam 50 m maka akan mengirim data ke RabbitMQ. Adapun untuk mengecek data yang diterima worker service dilakukan dengan log print payload
