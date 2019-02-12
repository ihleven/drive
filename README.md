# drive

Go virtual drive

type Storage interface {
Exists(path)
Listdir
Save
}

# Album

- ist Verzeichnis
- URL: alben/<pfad zu verzeichnis>
- Metadaten in Datei meta.json => wird von struct AlbumMeta geparst

## Unterverzeichnisse

- import: enthält einen Teil der Bilder des Albums, der mit der gleichen Kamera gemacht wurde (RX100, iPhone6s | matthias, iPhoneSE:Vicky) - sind nicht direkt

Augenblick: ## zeitlicher Aspekt, gruppiert zeitlich zusammenhängende Fotos

    # Momentaufnahme

    class AugenblickType(DjangoChoices):
        Transfer = ChoiceItem('TRANSFER') # # typ = [urlaub, ausflug, ,
        Moment   = ChoiceItem('MOMENT') # Ausflug
        Thema    = ChoiceItem('THEMA') # Info?' thema(Freiburg, Averbis), ort]
        Location = ChoiceItem('LOCATION') # Info?' thema(Freiburg, Averbis), ort]
        Info     = ChoiceItem('INFO')  # info(einkaufen, )

    id = models.AutoField(primary_key=True)
    typ = models.CharField(max_length=8, choices=AugenblickType.choices, null=True, blank=True)
    slug = models.CharField(max_length=50, unique=True, db_index=True, null=True, blank=True)
    album = models.ForeignKey(Album, on_delete=models.PROTECT, null=True, blank=True, related_name='augenblicke')

    title = models.CharField(max_length=100, blank=True)
    abstract = models.TextField(blank=True)  # vorspann
    description = models.TextField(blank=True)  # content

    #ausgangspunkt = models.ForeignKey(Location, to_field="slug", blank=True, null=True, on_delete=models.SET_NULL)
    ziel = models.ForeignKey(Location, to_field="slug", blank=True, null=True,
                             on_delete=models.SET_NULL)

    # day = models.ForeignKey('kalender.Kalendertag', null=True, blank=True, related_name='moments')
    date = models.DateField(null=True, blank=True)  # redundant
    onset = models.TimeField(null=True, blank=True)
    offset = models.TimeField(null=True, blank=True)
    end_date = models.DateField(null=True, blank=True)

    # location
    gps = models.TextField(blank=True)

    num_fotos = models.IntegerField(null=True, blank=True)
    index_foto = models.IntegerField(null=True, blank=True)  # Indexfoto, das für den Moment steht, nicht FK wegen cycle

    folder = models.ForeignKey(Folder, null=True, blank=True, related_name='augenblicke')

    wer = models.TextField(blank=True)
    was = models.TextField(blank=True)
    wann = models.TextField(blank=True)
    wo = models.TextField(blank=True)
    wie = models.TextField(blank=True)
    warum = models.TextField(blank=True)

    class Meta:
        db_table = 'fotos_augenblick'
        verbose_name_plural = 'Augenblicke'

    def __str__(self):
        return self.title
