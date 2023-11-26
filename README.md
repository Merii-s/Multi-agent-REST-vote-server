# AI30 - Systeme de vote



## Lancer le serveur

Pour lancer le serveur : go run .\ia04\cmd\launch_server.go

## Fonctionnalités implémentées

Les méthodes de vote possibles sur le serveur sont :
- Majorité simple (`"majority"`)
- Borda (borda) (`"borda"`)
- Elire le gagnant de Condorcet s'il existe (`"condorcet"`)
- Copeland (`"copeland"`)
- Aproval (`"approval"`)

Pour chacune de ces méthodes, /results renvoie le gagnant et un classement (propriétés winner et ranking)

## Fonctionnalités non implémentées

voteragent.go n'a pas été implémenté, seul le service de vote l'est