# pshinta
Pörssisähkön hinta

Tulostaa vuorokauden pörssisähkön hinnan (sisältäen ALV:n) komentoriville graafina. Klo 14:00 jälkeen tulostaa myös seuraavan vuorokauden hinnan.

Työkalu käyttää lähteenä tätä: [Api.spot-hinta.fi](https://spot-hinta.fi)

Esimerkki:

```
 0.0096 |
 0.0086 |                                            *
 0.0076 |                                           ***
 0.0066 |                                         *******
 0.0056 |           *                         *   ********
 0.0046 |          ***                      **************
 0.0036 |          *****                   ***************
 0.0026 |         *********                ***************
 0.0016 |        ***************          ****************
 0.0006 |       ******************************************
-0.0004 | *    *******************************************
-0.0014 | ************************************************
€/kWh ''|'''''|'''''|'''''|'''''|'''''|'''''|'''''|'''''|'''
        ^           ^           ^           ^           ^
      22.09       22.09       23.09       23.09       24.09
      00:00       12:00       00:00       12:00       00:00
```
