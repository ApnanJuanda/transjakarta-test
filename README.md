1. Execute query pada file 20250812001_create_table_vehicle_locations.up.sql 
2. Execute docker-compose up --build
3. import file collection dan env ke postman
4. hit API Start Publish Data, API ini akan mengirim data secara terus-menerus dengan jeda 2 detik dan akan berhenti jika radius antara latitude dan longitude dengan titik destination sudah 0
   atau bisa juga dihentikan dengan hit API Stop Publish Data
5. Ketika Start Publish Data berjalan maka proses yang akan dilakukan meliputi
pengiriman data ke MQTT, menerima data dari MQTT, menyimpan data ke PostgreSQL dan apabila radius sudah berada dalam 50 m maka akan mengirim data ke RabbitMQ. Adapun untuk mengecek data yang diterima worker service dilakukan dengan log print payload
