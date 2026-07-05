# TEPEGÖZ - Log Analiz ve Tehdit Tespit Aracı

Tepegöz log dosyalarını analiz eden, kural tabanlı canlı tespit yapan ve raporlayan bir siber güvenlik aracıdır.


---

## Proje Özellikleri


* **Hibrit Arayüz:** Hem menü tabanlı hem de komut satırı argümanları ile çalışabilir.
* **Statik Analiz:** Geçmiş log dosyalarını okuyarak IP istatistikleri ve zaman yoğunluk grafiği çıkarır.
* **Canlı İzleme :** Log dosyasını anlık takip eder  ve şüpheli durumları yakalar.
* **Raporlama:** Tespit edilen tehditleri `reports/` klasörüne CSV formatında, analiz özetlerini ise TXT formatında kaydeder.

---

## Kurulum ve Çalıştırma

Projenin sorunsuz çalışması için **Docker** ve **Docker Compose** kurulu olmalıdır.

**Adım 1:** `emirhan_ozgen.zip` dosyasını klasöre çıkartın.

**Adım 2:** Proje dizininde terminali açın ve şu komutu çalıştırın:

```bash
docker compose up --build
```

**Adım 3:** Uygulama başarıyla başlatıldığında karşınıza etkileşimli menü gelecektir.


## Dosya Yapısı 

Proje içerisindeki kritik dosyaların görevleri şöyledir:

* **`test.log`:** Projenin test edilmesi amacıyla ana dizine yerleştirilmiş örnek log dosyası.
* **`configs/rules.yaml`:** Tehdit tespitinde kullanılan Regex kurallarını barındırır. Yeni saldırı imzaları buraya eklenebilir.
* **`reports/`:** Uygulamanın ürettiği tüm analiz çıktıları bu klasörde otomatik olarak depolanır.

### Example (rules.yaml)
```yaml
rules:
  - id: "R001"
    name: "Critical System Fail"
    regex: "CRITICAL"
    level: "HIGH"
  
  - id: "R002"
    name: "Error SSH"
    regex: "Failed password"
    level: "MEDIUM"
```


## Çıktılar ve Raporlama

Araç Projenin isterlerine uygun olarak iki farklı formatta rapor üretir:

1.  **Log ve Alarm Kayıtları (CSV):**
    * Canlı izleme  modunda tespit edilen tehditler, **CSV formatında** kaydedilir.
    * **Dosya:** `reports/alerts_YYYY-MM-DD.csv`

2.  **Analiz Özeti (TXT):**
    * Statik analiz sonucunda elde edilen istatistiksel veriler ve **ASCII Zaman Grafikleri**, görsel bütünlüğün ve formatın korunması amacıyla **TXT formatında** raporlanır.
    * **Dosya:** `reports/summary_report_YYYY-MM-DD_Saat.txt`


