# Autopart-Service

### Create migrate file 
```bash
migrate create -ext sql -dir migrations -seq <name>
```

### Database Migration (golang-migrate)

```bash
migrate -path migrations -database "mysql://<user>:<password>@tcp(localhost:3306)/<dbname>" up

migrate -path migrations -database "mysql://<user>:<password>@tcp(localhost:3306)/<dbname>" down

migrate -path migrations -database "mysql://<user>:<password>@tcp(localhost:3306)/<dbname>" version

migrate -path migrations -database "mysql://<user>:<password>@tcp(localhost:3306)/<dbname>" force <version>
```
### Sqlc generate
```bash
sqlc generate -f internal/infrastructure/database/sqlc.yaml
```

### Database Migration (false case)
```bash
migrate -path migrations -database "mysql://<user>:<password>@tcp(localhost:3306)/<dbname>" version
Output เช่น 3 (dirty) → หมายถึง version ปัจจุบัน = 3 และอยู่ในสถานะ dirty

2. กำหนด version ให้ clean (force)
bash
Copy code
migrate -path migrations -database "mysql://<user>:<password>@tcp(localhost:3306)/<dbname>" force <version>
<version> = version ล่าสุดที่เสร็จสมบูรณ์ (ก่อนเกิด dirty)

ตัวอย่าง:

bash
Copy code
migrate -path migrations -database "mysql://admin:ec%3FZI121@tcp(localhost:3306)/database_server" force 2
⚠️ force ไม่แก้ไข table ใด ๆ แค่ปรับตัวเลข version ในตาราง schema_migrations

3. Run down migrations
bash
Copy code
migrate -path migrations -database "mysql://<user>:<password>@tcp(localhost:3306)/<dbname>" down
จะลบทุก migration ที่ run แล้ว (จาก version ล่าสุดลงไป)

ถ้า DB มี foreign key หรือ hierarchy ต้องลบ ลำดับจาก child → parent

ถ้ามี error ให้แก้ไข SQL ในไฟล์ down migration ก่อน เช่น ลบ sub-category ก่อน parent

4. Run up migrations ใหม่
bash
Copy code
migrate -path migrations -database "mysql://<user>:<password>@tcp(localhost:3306)/<dbname>" up
จะ run migration ตั้งแต่ version 1 → ล่าสุด

Insert seed data ใหม่ทั้งหมด

5. ตรวจสอบสถานะหลัง migration
bash
Copy code
migrate -path migrations -database "mysql://<user>:<password>@tcp(localhost:3306)/<dbname>" version
Output ควรเป็น version ล่าสุด dirty = false