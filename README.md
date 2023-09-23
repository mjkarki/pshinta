# pshinta
Pörssisähkön hinta

Tulostaa vuorokauden pörssisähkön hinnan (sisältäen ALV:n) komentoriville graafina. Klo 14:00 jälkeen tulostaa myös seuraavan vuorokauden hinnan.

Työkalu käyttää lähteenä tätä: [Api.spot-hinta.fi](https://spot-hinta.fi)

Esimerkki:

```
 0.0102 |                                            *
 0.0092 |                                           ***
 0.0082 |                    *                      ***
 0.0072 |                   ***                     ***
 0.0062 |                  *****                    ****
 0.0052 |                 *******                  ******
 0.0042 |            *** *********                 ******
 0.0032 |           ***************        **      *******
 0.0022 |          *****************      ***      *******
 0.0012 |         ****************************     *******
 0.0002 |       *******************************    *******
-0.0008 | **************************************  ********
-0.0018 | **************************************  ********
-0.0028 | ************************************** *********
-0.0038 | ************************************************
  €/kWh ''|'''''|'''''|'''''|'''''|'''''|'''''|'''''|'''''|'''
          ^           ^           ^           ^           ^
        23.09       23.09       24.09       24.09       25.09
        00:00       12:00       00:00       12:00       00:00
```
