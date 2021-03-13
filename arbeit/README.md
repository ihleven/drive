


Job
 * 20.8.2001 - 31.11.2007 Uni
 * 1.4.2008 - 31.1.2013 Uniklinik
 * 1.2.2013 - 31.7.2018 Averbis
 * 1.8.2018 -           Interchalet


Arbeitsjahr ()
 * kann mehrere pro Kalenderjahr geben (Averbis bis Juli 2018 / IC ab August oder auch )
 - Account, Job, Kalenderjahr
 - Tage Urlaub, Sonderaurlaub,...

Arbeitsmonat

 - Arbeitsjahr (Account, Job), Kalendermonat
 - Tage Urlaub, Sonderaurlaub,...

Arbeitstag

    ID          int
	Account     *domain.Account
	Jahr *Arbeitsjahr // 
	Job         *Job // Ref sollte auf Arbeitsjahr gehen, das dann auf denJob geht  
	Tag         *Kalendertag
	TagStatus   *string // A:Arb.tag, H: HalbAT, N: Werkt (nicht-Arbeitst.), S: Sonntag & F: Feiertag
	Status      *string // B: Büro, D: Dienstreise, H: Homeoffice, Z: Zeitausgleich, K: krank, U: urlaub
    Urlaubstage float // 0, 0.5, 1 - wieviele Tage vom Urlaubskontingent abgehen
	Typ         *string 
	Soll        *float64
	Start       *time.Time // beim HauptJob
	Ende        *time.Time
	Brutto      *float64
	Pausen      *float64
	Extra       *float64
	Netto       *float64 // Brutto + Extra - Pausen
	Differenz   *float64 // Netto - Soll
	Saldo       *float64 // ergibt sich aus Saldo Vortag + Differenz
	Zeitspannen []Zeitspanne

Zeitspanne

 - Person, Job, Arbeitstag?/Kalendertag?