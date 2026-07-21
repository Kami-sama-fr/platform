# Roadmap API — Plateforme d'Administration Kami-Sama

> Inventaire complet des endpoints API nécessaires pour chaque page du dashboard admin.
> Statut : `✅` = implémenté | `🔶` = modèle existe, endpoint manque | `❌` = à créer

---

## Table des matières

1. [Overview `/dash`](#1-overview--dash)
2. [Content](#2-content)
3. [Community](#3-community)
4. [Scheduling](#4-scheduling)
5. [Analytics](#5-analytics)
6. [Media](#6-media)
7. [Support](#7-support)
8. [System](#8-system)
9. [Settings](#9-settings)

---

## 1. Overview `/dash`

Page : `apps/app/(platform)/dash/page.tsx`
Composant : `Overview` — données actuellement mockées dans `overview-data.ts`

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/dashboard/stats` | KPIs : abonnés actifs, vues du jour, revenus mensuels, nb anime catalogue + deltas | ❌ |
| 2 | GET | `/api/v1/dashboard/weekly-views` | Vues quotidiennes sur 7 jours (pour graphique barres) | ❌ |
| 3 | GET | `/api/v1/dashboard/subscription-distribution` | Répartition abonnements : plan, nb membres, pourcentage | ❌ |
| 4 | GET | `/api/v1/dashboard/top-anime` | Top anime par vues avec note (classement) | ❌ |
| 5 | GET | `/api/v1/dashboard/recent-uploads` | Derniers uploads : titre, épisode, langue, status, timestamp | ❌ |

**Modèles existants exploitables :** `WatchProgress`, `WatchHistory`, `MediaAsset`, `Anime`

---

## 2. Content

### 2.1 Anime — `/dash/content/anime`

Page : `apps/app/(platform)/dash/content/anime/page.tsx`
Composant : `AnimeManager` — CRUD complet côté client, données mockées

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/anime` | Liste paginée + filtres (status, genre, studio, année, tri, recherche) | ✅ |
| 2 | POST | `/api/v1/anime` | Créer un anime | ✅ |
| 3 | GET | `/api/v1/anime/:animeId` | Détail d'un anime | ✅ |
| 4 | PATCH | `/api/v1/anime/:animeId` | Modifier un anime | ✅ |
| 5 | DELETE | `/api/v1/anime/:animeId` | Supprimer un anime | ✅ |
| 6 | POST | `/api/v1/anime/:animeId/duplicate` | Dupliquer un anime | ❌ |
| 7 | POST | `/api/v1/anime/bulk-delete` | Suppression multiple par IDs | ❌ |
| 8 | GET | `/api/v1/genres` | Liste des genres (pour filtres/formulaire) | ✅ |

**Note :** Les endpoints 1-5 et 8 sont déjà enregistrés dans `routes.go`. Le frontend doit être wiré pour les appeler.

### 2.2 Library — `/dash/content/library`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/libraries` | Liste des bibliothèques média (Jellyfin + local) | 🔶 `source/libraries` existe |
| 2 | POST | `/api/v1/libraries` | Créer/importer une bibliothèque | ❌ |
| 3 | GET | `/api/v1/libraries/:libraryId` | Détail d'une bibliothèque | 🔶 `source/libraries/:libraryId` existe |
| 4 | PATCH | `/api/v1/libraries/:libraryId` | Modifier une bibliothèque | ❌ |
| 5 | DELETE | `/api/v1/libraries/:libraryId` | Supprimer une bibliothèque | ❌ |
| 6 | POST | `/api/v1/libraries/:libraryId/sync` | Synchroniser une bibliothèque | ✅ `source/libraries/:libraryId/sync` |
| 7 | GET | `/api/v1/libraries/:libraryId/sync` | Statut de synchronisation | ✅ |

### 2.3 Categories — `/dash/content/categories`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/categories` | Liste des catégories | ❌ |
| 2 | POST | `/api/v1/categories` | Créer une catégorie | ❌ |
| 3 | GET | `/api/v1/categories/:categoryId` | Détail d'une catégorie | ❌ |
| 4 | PATCH | `/api/v1/categories/:categoryId` | Modifier une catégorie | ❌ |
| 5 | DELETE | `/api/v1/categories/:categoryId` | Supprimer une catégorie | ❌ |

**Note :** Aucun modèle `Category` n'existe. Le modèle `Anime` a un champ `Source` mais pas de catégories. Il faudra créer le modèle + repo + service + routes.

### 2.4 Genres — `/dash/content/genres`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/genres` | Liste des genres | ✅ |
| 2 | POST | `/api/v1/genres` | Créer un genre | ✅ |
| 3 | GET | `/api/v1/genres/:genreId` | Détail d'un genre | ✅ |
| 4 | PATCH | `/api/v1/genres/:genreId` | Modifier un genre | ✅ |
| 5 | DELETE | `/api/v1/genres/:genreId` | Supprimer un genre | ✅ |

**Status :** Tous les endpoints existent. La page frontend doit être construite.

### 2.5 Tags — `/dash/content/tags`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/tags` | Liste des tags | ❌ |
| 2 | POST | `/api/v1/tags` | Créer un tag | ❌ |
| 3 | GET | `/api/v1/tags/:tagId` | Détail d'un tag | ❌ |
| 4 | PATCH | `/api/v1/tags/:tagId` | Modifier un tag | ❌ |
| 5 | DELETE | `/api/v1/tags/:tagId` | Supprimer un tag | ❌ |

**Note :** Le modèle `Tag` existe dans `models/anime.go` mais n'a pas de repo/service/routes.

### 2.6 Providers — `/dash/content/providers`

Page référencée dans `nav-config.ts` mais aucun fichier `page.tsx` n'existe

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/providers` | Liste des fournisseurs de contenu | ❌ |
| 2 | POST | `/api/v1/providers` | Ajouter un fournisseur | ❌ |
| 3 | GET | `/api/v1/providers/:providerId` | Détail d'un fournisseur | ❌ |
| 4 | PATCH | `/api/v1/providers/:providerId` | Modifier un fournisseur | ❌ |
| 5 | DELETE | `/api/v1/providers/:providerId` | Supprimer un fournisseur | ❌ |

**Note :** Modèle `SourceConfig` existe pour stocker les configs de sources (Jellyfin, etc.) mais pas de CRUD admin dédié.

---

## 3. Community

### 3.1 Users — `/dash/community/users`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/admin/users` | Liste paginée des utilisateurs (recherche, filtres status/role) | ❌ |
| 2 | GET | `/api/v1/admin/users/:userId` | Détail d'un utilisateur | ❌ |
| 3 | PATCH | `/api/v1/admin/users/:userId` | Modifier un utilisateur (status, rôles, bio) | ❌ |
| 4 | DELETE | `/api/v1/admin/users/:userId` | Désactiver/supprimer un utilisateur | ❌ |
| 5 | POST | `/api/v1/admin/users/:userId/disable` | Désactiver un compte | ❌ |
| 6 | POST | `/api/v1/admin/users/:userId/enable` | Réactiver un compte | ❌ |
| 7 | GET | `/api/v1/admin/users/:userId/sessions` | Sessions actives d'un utilisateur | 🔶 `auth/sessions` existe (self) |
| 8 | DELETE | `/api/v1/admin/users/:userId/sessions/:sessionId` | Révoquer une session | 🔶 `auth/sessions/:sessionId` existe (self) |

**Modèles existants :** `User`, `AuthSession`, `AuthAccount`
**Note :** Les endpoints `GET/PATCH /api/v1/me` existent mais sont user-scoped. Il faut des endpoints admin.

### 3.2 Profiles — `/dash/community/profiles`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/admin/profiles` | Liste des profils publics | ❌ |
| 2 | GET | `/api/v1/admin/profiles/:userId` | Profil public d'un utilisateur | ❌ |
| 3 | PATCH | `/api/v1/admin/profiles/:userId` | Modifier un profil (avatar, bio, display name) | 🔶 `PATCH /api/v1/me` existe (self) |
| 4 | DELETE | `/api/v1/admin/profiles/:userId` | Supprimer un profil | ❌ |

### 3.3 Roles — `/dash/community/roles`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/admin/roles` | Liste des rôles | ❌ |
| 2 | POST | `/api/v1/admin/roles` | Créer un rôle | ❌ |
| 3 | GET | `/api/v1/admin/roles/:roleId` | Détail d'un rôle + permissions | ❌ |
| 4 | PATCH | `/api/v1/admin/roles/:roleId` | Modifier un rôle | ❌ |
| 5 | DELETE | `/api/v1/admin/roles/:roleId` | Supprimer un rôle | ❌ |
| 6 | POST | `/api/v1/admin/roles/:roleId/assign` | Assigner un rôle à un utilisateur | ❌ |
| 7 | DELETE | `/api/v1/admin/roles/:roleId/assign/:userId` | Retirer un rôle | ❌ |

**Note :** Les rôles sont actuellement gérés en dur dans le code (`superadmin`, `admin`, `user`, `workspace:owner`). Pas de modèle DB dédié pour les définitions de rôles.

### 3.4 Comments — `/dash/community/comments`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/admin/comments` | Liste paginée des commentaires (filtres: status, auteur, anime) | ❌ |
| 2 | GET | `/api/v1/admin/comments/:commentId` | Détail d'un commentaire | 🔶 via `community/reviews/:reviewId/comments` |
| 3 | PATCH | `/api/v1/admin/comments/:commentId` | Modérer un commentaire (approuver, masquer, supprimer) | 🔶 `community/comments/:commentId` existe |
| 4 | DELETE | `/api/v1/admin/comments/:commentId` | Supprimer un commentaire | 🔶 `community/comments/:commentId` existe |
| 5 | POST | `/api/v1/admin/comments/:commentId/approve` | Approuver un commentaire | ❌ |
| 6 | POST | `/api/v1/admin/comments/:commentId/flag` | Marquer comme signalé | ❌ |

**Modèle existant :** `Comment` avec champ `DeletedAt` (soft delete)

### 3.5 Reviews — `/dash/community/reviews`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/admin/reviews` | Liste paginée des reviews (filtres: rating, auteur, anime) | ❌ |
| 2 | GET | `/api/v1/admin/reviews/:reviewId` | Détail d'une review | ✅ `community/reviews/:reviewId` |
| 3 | PATCH | `/api/v1/admin/reviews/:reviewId` | Modérer une review | ✅ `community/reviews/:reviewId` |
| 4 | DELETE | `/api/v1/admin/reviews/:reviewId` | Supprimer une review | ✅ `community/reviews/:reviewId` |
| 5 | POST | `/api/v1/admin/reviews/:reviewId/feature` | Mettre en avant une review | ❌ |

### 3.6 Reports — `/dash/community/reports`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/admin/reports` | Liste paginée des signalements (filtres: status, type) | 🔶 `community/reviews` existe |
| 2 | GET | `/api/v1/admin/reports/:reportId` | Détail d'un signalement | ❌ |
| 3 | PATCH | `/api/v1/admin/reports/:reportId` | Traiter un signalement (résoudre, rejeter) | 🔶 `community/reports/:reportId` existe |
| 4 | POST | `/api/v1/admin/reports/:reportId/resolve` | Marquer comme résolu | ❌ |
| 5 | POST | `/api/v1/admin/reports/:reportId/dismiss` | Rejeter le signalement | ❌ |

**Modèle existant :** `Report` avec `Status`, `ResolvedBy`, `ResolvedAt`

### 3.7 Moderations — `/dash/community/moderations`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/admin/moderations` | File d'attente de modération | ❌ |
| 2 | GET | `/api/v1/admin/moderations/:moderationId` | Détail d'un item en modération | ❌ |
| 3 | POST | `/api/v1/admin/moderations/:moderationId/approve` | Approuver | ❌ |
| 4 | POST | `/api/v1/admin/moderations/:moderationId/reject` | Rejeter | ❌ |
| 5 | POST | `/api/v1/admin/moderations/:moderationId/escalate` | Escalader | ❌ |

**Note :** Pas de modèle dédié. La modération est distribuée entre `Comment`, `Review`, `Report`. Il faudra un modèle `ModerationQueue` ou unifier via `Report`.

### 3.8 Watchlists — `/dash/community/watchlists`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/admin/watchlists` | Liste des watchlists (analytics) | ✅ `community/watchlists` |
| 2 | GET | `/api/v1/admin/watchlists/:watchlistId` | Détail d'une watchlist | ✅ `community/watchlists/:watchlistId` |
| 3 | GET | `/api/v1/admin/watchlists/stats` | Statistiques watchlists (populaire, tendances) | ❌ |

### 3.9 Permissions — `/dash/community/permissions`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/admin/permissions` | Matrice des permissions (rôle × ressource × action) | ❌ |
| 2 | PATCH | `/api/v1/admin/permissions` | Mettre à jour les permissions d'un rôle | ❌ |
| 3 | GET | `/api/v1/admin/permissions/effective/:userId` | Permissions effectives d'un utilisateur | ❌ |

---

## 4. Scheduling

### 4.1 Calendar — `/dash/scheduling/calendar`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/scheduling/calendar` | Événements du calendrier (plage de dates) | ❌ |
| 2 | POST | `/api/v1/scheduling/calendar` | Créer un événement calendrier | ❌ |
| 3 | GET | `/api/v1/scheduling/calendar/:eventId` | Détail d'un événement | ❌ |
| 4 | PATCH | `/api/v1/scheduling/calendar/:eventId` | Modifier un événement | ❌ |
| 5 | DELETE | `/api/v1/scheduling/calendar/:eventId` | Supprimer un événement | ❌ |

**Note :** Le modèle `ReleaseSchedule` peut servir de base. Il faudra probablement un modèle `CalendarEvent` plus générique.

### 4.2 Releases — `/dash/scheduling/releases`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/scheduling/releases` | Liste des sorties planifiées | 🔶 `scheduling/upcoming` existe |
| 2 | POST | `/api/v1/scheduling/releases` | Planifier une sortie | 🔶 `scheduling/schedules` existe |
| 3 | GET | `/api/v1/scheduling/releases/:releaseId` | Détail d'une sortie | ❌ |
| 4 | PATCH | `/api/v1/scheduling/releases/:releaseId` | Modifier une sortie | 🔶 `scheduling/schedules/:scheduleId` |
| 5 | DELETE | `/api/v1/scheduling/releases/:releaseId` | Annuler une sortie | ❌ |
| 6 | POST | `/api/v1/scheduling/releases/:releaseId/publish` | Publier immédiatement | ❌ |

### 4.3 Premieres — `/dash/scheduling/premieres`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/scheduling/premieres` | Liste des premières | ❌ |
| 2 | POST | `/api/v1/scheduling/premieres` | Planifier une première | ❌ |
| 3 | GET | `/api/v1/scheduling/premieres/:premiereId` | Détail d'une première | ❌ |
| 4 | PATCH | `/api/v1/scheduling/premieres/:premiereId` | Modifier une première | ❌ |
| 5 | DELETE | `/api/v1/scheduling/premieres/:premiereId` | Annuler une première | ❌ |

**Note :** Pas de modèle dédié. Pourrait être une extension de `ReleaseSchedule` avec un type `premiere`.

### 4.4 Simulcasts — `/dash/scheduling/simulcasts`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/scheduling/simulcasts` | Liste des simulcasts | ✅ |
| 2 | POST | `/api/v1/scheduling/simulcasts` | Créer un simulcast | ✅ |
| 3 | GET | `/api/v1/scheduling/simulcasts/:simulcastId` | Détail d'un simulcast | ✅ |
| 4 | PATCH | `/api/v1/scheduling/simulcasts/:simulcastId` | Modifier un simulcast | ✅ |
| 5 | DELETE | `/api/v1/scheduling/simulcasts/:simulcastId` | Supprimer un simulcast | ✅ |
| 6 | GET | `/api/v1/scheduling/simulcasts/week/:weekday` | Simulcasts par jour | ✅ |

**Status :** Tous les endpoints existent. La page frontend doit être construite.

### 4.5 Notifications — `/dash/scheduling/notifications`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/admin/notifications/templates` | Liste des templates de notification | ❌ |
| 2 | POST | `/api/v1/admin/notifications/templates` | Créer un template | ❌ |
| 3 | PATCH | `/api/v1/admin/notifications/templates/:templateId` | Modifier un template | ❌ |
| 4 | DELETE | `/api/v1/admin/notifications/templates/:templateId` | Supprimer un template | ❌ |
| 5 | POST | `/api/v1/admin/notifications/send` | Envoyer une notification (broadcast ou ciblée) | ❌ |
| 6 | GET | `/api/v1/admin/notifications/history` | Historique des notifications envoyées | 🔶 `notifications` existe (user-scoped) |

**Modèle existant :** `NotificationTemplate` (pas de repo/routes)

---

## 5. Analytics

### 5.1 Overview — `/dash/analytics/overview`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/analytics/overview` | Métriques agrégées : vues totales, utilisateurs actifs, temps moyen, taux de complétion | ❌ |
| 2 | GET | `/api/v1/analytics/overview/period` | Métriques par période (jour/semaine/mois) avec tendances | ❌ |

### 5.2 Watch Time — `/dash/analytics/watch-time`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/analytics/watch-time` | Données de temps de visionnage par période | ❌ |
| 2 | GET | `/api/v1/analytics/watch-time/by-anime` | Temps moyen par anime | ❌ |
| 3 | GET | `/api/v1/analytics/watch-time/by-episode` | Temps moyen par épisode | ❌ |
| 4 | GET | `/api/v1/analytics/watch-time/histogram` | Distribution des sessions de visionnage | ❌ |

**Modèles existants :** `WatchProgress`, `WatchHistory`

### 5.3 Devices — `/dash/analytics/devices`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/analytics/devices` | Répartition par type d'appareil | ❌ |
| 2 | GET | `/api/v1/analytics/devices/browsers` | Répartition par navigateur | ❌ |
| 3 | GET | `/api/v1/analytics/devices/os` | Répartition par OS | ❌ |

**Modèle existant :** `Device` (dans `platform.go`) — mais pas de service analytics.

### 5.4 Popular — `/dash/analytics/popular`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/analytics/popular` | Contenu le plus populaire (vues, temps moyen) | ❌ |
| 2 | GET | `/api/v1/analytics/popular/trending` | Tendances actuelles | ❌ |
| 3 | GET | `/api/v1/analytics/popular/new` | Nouvelles entrées populaires | ❌ |

### 5.5 Geography — `/dash/analytics/geography`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/analytics/geography` | Distribution des vues par pays/région | ❌ |
| 2 | GET | `/api/v1/analytics/geography/top-countries` | Top pays par nombre de vues | ❌ |

**Modèle existant :** `Device` a `Country`, `City` mais pas d'agrégation analytics.

### 5.6 Active Users — `/dash/analytics/actives`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/analytics/active-users` | Nombre d'utilisateurs actifs par période (DAU/WAU/MAU) | ❌ |
| 2 | GET | `/api/v1/analytics/active-users/retention` | Taux de rétention | ❌ |
| 3 | GET | `/api/v1/analytics/active-users/sessions` | Nombre moyen de sessions par utilisateur | ❌ |

---

## 6. Media

### 6.1 Videos — `/dash/media/videos`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/media/videos` | Liste des assets vidéo (filtres: type, status, anime) | ✅ `media?type=video` |
| 2 | POST | `/api/v1/media/videos` | Uploader un vidéo | ✅ `media` |
| 3 | GET | `/api/v1/media/videos/:videoId` | Détail d'un asset vidéo | ✅ `media/:mediaId` |
| 4 | PATCH | `/api/v1/media/videos/:videoId` | Modifier les métadonnées | ✅ `media/:mediaId` |
| 5 | DELETE | `/api/v1/media/videos/:videoId` | Supprimer un asset vidéo | ✅ `media/:mediaId` |

### 6.2 Thumbnails — `/dash/media/thumbnails`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/media/thumbnails` | Liste des thumbnails | ✅ `media?type=thumbnail` |
| 2 | POST | `/api/v1/media/thumbnails` | Uploader un thumbnail | ✅ `media` |
| 3 | GET | `/api/v1/media/thumbnails/:thumbnailId` | Détail d'un thumbnail | ✅ `media/:mediaId` |
| 4 | DELETE | `/api/v1/media/thumbnails/:thumbnailId` | Supprimer un thumbnail | ✅ `media/:mediaId` |
| 5 | POST | `/api/v1/media/thumbnails/generate` | Générer automatiquement un thumbnail | ❌ |

### 6.3 Posters — `/dash/media/posters`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/media/posters` | Liste des posters | ✅ `media?type=poster` |
| 2 | POST | `/api/v1/media/posters` | Uploader un poster | ✅ `media` |
| 3 | GET | `/api/v1/media/posters/:posterId` | Détail d'un poster | ✅ `media/:mediaId` |
| 4 | DELETE | `/api/v1/media/posters/:posterId` | Supprimer un poster | ✅ `media/:mediaId` |

### 6.4 Trailers — `/dash/media/trailers`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/media/trailers` | Liste des bandes-annonces | ✅ `media?type=trailer` |
| 2 | POST | `/api/v1/media/trailers` | Uploader une BA | ✅ `media` |
| 3 | GET | `/api/v1/media/trailers/:trailerId` | Détail d'une BA | ✅ `media/:mediaId` |
| 4 | DELETE | `/api/v1/media/trailers/:trailerId` | Supprimer une BA | ✅ `media/:mediaId` |

### 6.5 Audio — `/dash/media/audio`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/media/audio` | Liste des pistes audio | ✅ `media?type=audio` |
| 2 | POST | `/api/v1/media/audio` | Uploader une piste audio | ✅ `media` |
| 3 | GET | `/api/v1/media/audio/:audioId` | Détail d'une piste audio | ✅ `media/:mediaId` |
| 4 | DELETE | `/api/v1/media/audio/:audioId` | Supprimer une piste audio | ✅ `media/:mediaId` |

### 6.6 Subtitles — `/dash/media/subtitles`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/media/subtitles` | Liste des sous-titres | ✅ `media?type=subtitle` |
| 2 | POST | `/api/v1/media/subtitles` | Uploader un fichier sous-titre | ✅ `media` |
| 3 | GET | `/api/v1/media/subtitles/:subtitleId` | Détail d'un sous-titre | ✅ `media/:mediaId` |
| 4 | DELETE | `/api/v1/media/subtitles/:subtitleId` | Supprimer un sous-titre | ✅ `media/:mediaId` |

### 6.7 Encoding — `/dash/media/encoding`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/media/encoding-jobs` | Liste des jobs d'encodage (filtres: status) | ✅ |
| 2 | GET | `/api/v1/media/encoding-jobs/:jobId` | Détail d'un job | ✅ |
| 3 | POST | `/api/v1/media/encoding-jobs/:jobId/retry` | Relancer un job échoué | ❌ |
| 4 | POST | `/api/v1/media/encoding-jobs/:jobId/cancel` | Annuler un job en cours | ❌ |
| 5 | GET | `/api/v1/media/encoding/profiles` | Profils d'encodage disponibles | ❌ |

### 6.8 Uploads — `/dash/media/uploads`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/media/uploads` | File d'upload en cours + historique | ❌ |
| 2 | POST | `/api/v1/media/uploads` | Initier un upload (multipart ou presigned URL) | ❌ |
| 3 | GET | `/api/v1/media/uploads/:uploadId` | Progression d'un upload | ❌ |
| 4 | DELETE | `/api/v1/media/uploads/:uploadId` | Annuler un upload | ❌ |
| 5 | POST | `/api/v1/media/uploads/:uploadId/complete` | Finaliser un upload | ❌ |

---

## 7. Support

### 7.1 Tickets — `/dash/support/tickets`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/support/tickets` | Liste des tickets (filtres: status, priorité, catégorie) | ❌ |
| 2 | POST | `/api/v1/support/tickets` | Créer un ticket | ❌ |
| 3 | GET | `/api/v1/support/tickets/:ticketId` | Détail d'un ticket | ❌ |
| 4 | PATCH | `/api/v1/support/tickets/:ticketId` | Modifier un ticket (status, assignation) | ❌ |
| 5 | POST | `/api/v1/support/tickets/:ticketId/reply` | Répondre à un ticket | ❌ |
| 6 | POST | `/api/v1/support/tickets/:ticketId/close` | Fermer un ticket | ❌ |
| 7 | POST | `/api/v1/support/tickets/:ticketId/escalate` | Escalader un ticket | ❌ |

**Note :** Aucun modèle `Ticket` n'existe. À créer complète.

### 7.2 Contact — `/dash/support/contact`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/support/contact` | Liste des messages de contact | ❌ |
| 2 | GET | `/api/v1/support/contact/:messageId` | Détail d'un message | ❌ |
| 3 | PATCH | `/api/v1/support/contact/:messageId` | Marquer comme lu/répondu | ❌ |
| 4 | POST | `/api/v1/support/contact/:messageId/reply` | Répondre | ❌ |
| 5 | DELETE | `/api/v1/support/contact/:messageId` | Supprimer | ❌ |

**Note :** Le modèle `Contact` existe mais est un modèle de carnet d'adresses, pas un formulaire de contact.

### 7.3 FAQ — `/dash/support/faq`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/support/faq` | Liste des FAQ | ❌ |
| 2 | POST | `/api/v1/support/faq` | Créer une entrée FAQ | ❌ |
| 3 | GET | `/api/v1/support/faq/:faqId` | Détail d'une entrée | ❌ |
| 4 | PATCH | `/api/v1/support/faq/:faqId` | Modifier une entrée | ❌ |
| 5 | DELETE | `/api/v1/support/faq/:faqId` | Supprimer une entrée | ❌ |
| 6 | PUT | `/api/v1/support/faq/reorder` | Réordonner les entrées | ❌ |

### 7.4 Abuse — `/dash/support/abuse`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/support/abuse` | Liste des signalements d'abus | 🔶 `community/reports` existe |
| 2 | GET | `/api/v1/support/abuse/:abuseId` | Détail d'un rapport d'abus | 🔶 `community/reports/:reportId` |
| 3 | PATCH | `/api/v1/support/abuse/:abuseId` | Traiter (escalader, bloquer, rejeter) | 🔶 `community/reports/:reportId` |
| 4 | POST | `/api/v1/support/abuse/:abuseId/block` | Bloquer l'utilisateur signalé | ❌ |
| 5 | POST | `/api/v1/support/abuse/:abuseId/warn` | Avertir l'utilisateur signalé | ❌ |

### 7.5 Logs — `/dash/support/logs`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/support/logs` | Journaux d'activité support (filtres: date, action, utilisateur) | ❌ |
| 2 | GET | `/api/v1/support/logs/export` | Exporter les logs (CSV/JSON) | ❌ |

**Modèle existant :** `AuditLog` (mais dans le contexte workspace, pas support)

---

## 8. System

### 8.1 Health — `/dash/system/health`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/system/health` | État de santé de tous les services | ✅ `/health/live` + `/health/ready` |
| 2 | GET | `/api/v1/system/health/services` | État détaillé par service (DB, Redis, Jellyfin, CDN) | 🔶 basique |
| 3 | GET | `/api/v1/system/health/uptime` | Uptime du serveur | ❌ |
| 4 | GET | `/api/v1/system/health/metrics` | Métriques système (CPU, RAM, disque) | 🔶 `/metrics` existe (Prometheus) |

### 8.2 Logs — `/dash/system/logs`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/system/logs` | Journaux système (filtres: niveau, source, date) | ❌ |
| 2 | GET | `/api/v1/system/logs/:logId` | Détail d'une entrée log | ❌ |
| 3 | GET | `/api/v1/system/logs/search` | Recherche dans les logs | ❌ |
| 4 | GET | `/api/v1/system/logs/stream` | Stream temps réel (SSE/WebSocket) | ❌ |

**Modèle existant :** `AuthAuditEvent` — mais limité aux événements d'auth.

### 8.3 Queue — `/dash/system/queue`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/system/queue` | État de la file d'attente (nb jobs, en attente, en cours, échoués) | ❌ |
| 2 | GET | `/api/v1/system/queue/jobs` | Liste des jobs (filtres: status, type) | ❌ |
| 3 | POST | `/api/v1/system/queue/jobs/:jobId/retry` | Relancer un job échoué | ❌ |
| 4 | POST | `/api/v1/system/queue/jobs/:jobId/cancel` | Annuler un job | ❌ |
| 5 | POST | `/api/v1/system/queue/flush` | Vider la file | ❌ |

### 8.4 Cache — `/dash/system/cache`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/system/cache` | État du cache (hit rate, taille, keys) | ❌ |
| 2 | POST | `/api/v1/system/cache/flush` | Vider tout le cache | ❌ |
| 3 | POST | `/api/v1/system/cache/flush/:pattern` | Vider le cache par pattern | ❌ |
| 4 | GET | `/api/v1/system/cache/keys` | Lister les clés de cache | ❌ |
| 5 | DELETE | `/api/v1/system/cache/keys/:key` | Supprimer une clé du cache | ❌ |

### 8.5 Search — `/dash/system/search`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/system/search` | État du moteur de recherche (Meilisearch) | ❌ |
| 2 | POST | `/api/v1/system/search/reindex` | Relancer l'indexation complète | ❌ |
| 3 | GET | `/api/v1/system/search/indexes` | Liste des indexes | ❌ |
| 4 | GET | `/api/v1/system/search/indexes/:indexName` | Stats d'un index | ❌ |
| 5 | POST | `/api/v1/system/search/indexes/:indexName/update` | Mettre à jour un index | ❌ |

**Note :** Le frontend a des routes `/api/v1/search/*` mais aucun endpoint system admin.

### 8.6 Background Jobs — `/dash/system/background`

Page vide — à implémenter (note : le dossier a une typo `backroung`)

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/system/background-jobs` | Liste des tâches planifiées (cron) | ❌ |
| 2 | POST | `/api/v1/system/background-jobs/:jobId/run` | Exécuter manuellement une tâche | ❌ |
| 3 | POST | `/api/v1/system/background-jobs/:jobId/pause` | Pause une tâche | ❌ |
| 4 | POST | `/api/v1/system/background-jobs/:jobId/resume` | Reprendre une tâche | ❌ |
| 5 | GET | `/api/v1/system/background-jobs/history` | Historique d'exécution | ❌ |

---

## 9. Settings

### 9.1 General — `/dash/settings`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/settings` | Tous les paramètres système | ✅ |
| 2 | PUT | `/api/v1/settings/:key` | Modifier un paramètre | ✅ |
| 3 | GET | `/api/v1/settings/general` | Paramètres généraux (nom, description, locale) | 🔶 via `/api/v1/settings` |
| 4 | PUT | `/api/v1/settings/general` | Modifier les paramètres généraux | 🔶 via `/api/v1/settings/:key` |

### 9.2 Security — `/dash/settings/security`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/settings/security` | Paramètres de sécurité | ❌ |
| 2 | PUT | `/api/v1/settings/security` | Modifier les paramètres de sécurité | ❌ |
| 3 | GET | `/api/v1/settings/security/sessions` | Config des sessions (TTL, max) | ❌ |
| 4 | PUT | `/api/v1/settings/security/sessions` | Modifier la config sessions | ❌ |
| 5 | GET | `/api/v1/settings/security/rate-limit` | Config rate limiting | ❌ |
| 6 | PUT | `/api/v1/settings/security/rate-limit` | Modifier le rate limiting | ❌ |
| 7 | GET | `/api/v1/settings/security/2fa` | Config 2FA obligatoire | ❌ |
| 8 | PUT | `/api/v1/settings/security/2fa` | Modifier la config 2FA | ❌ |

### 9.3 Branding — `/dash/settings/branding`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/settings/branding` | Paramètres visuels (logo, couleurs, favicon) | ❌ |
| 2 | PUT | `/api/v1/settings/branding` | Modifier le branding | ❌ |
| 3 | POST | `/api/v1/settings/branding/logo` | Uploader le logo | ❌ |
| 4 | POST | `/api/v1/settings/branding/favicon` | Uploader le favicon | ❌ |
| 5 | POST | `/api/v1/settings/branding/preview` | Prévisualiser les changements | ❌ |

### 9.4 Email — `/dash/settings/email`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/settings/email` | Config SMTP / provider email | ❌ |
| 2 | PUT | `/api/v1/settings/email` | Modifier la config email | ❌ |
| 3 | POST | `/api/v1/settings/email/test` | Envoyer un email de test | ❌ |
| 4 | GET | `/api/v1/settings/email/templates` | Templates d'email | ❌ |
| 5 | PATCH | `/api/v1/settings/email/templates/:templateId` | Modifier un template | ❌ |
| 6 | GET | `/api/v1/settings/email/logs` | Historique des emails envoyés | ❌ |

### 9.5 SEO — `/dash/settings/seo`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/settings/seo` | Métadonnées SEO globales | ✅ |
| 2 | PUT | `/api/v1/settings/seo` | Modifier les métadonnées SEO | ✅ |
| 3 | GET | `/api/v1/settings/seo/pages` | SEO par page (paths dynamiques) | ❌ |
| 4 | PUT | `/api/v1/settings/seo/pages/:pagePath` | Modifier le SEO d'une page | ❌ |
| 5 | GET | `/api/v1/settings/seo/sitemap` | Générer/visualiser le sitemap | ❌ |
| 6 | GET | `/api/v1/settings/seo/robots` | Visualiser le robots.txt | ❌ |

**Modèle existant :** `SeoMeta` (dans `settings.go`)

### 9.6 Storage — `/dash/settings/storage`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/settings/storage` | Config stockage (S3/GCS/Azure/Local) | ❌ |
| 2 | PUT | `/api/v1/settings/storage` | Modifier la config stockage | ❌ |
| 3 | POST | `/api/v1/settings/storage/test` | Tester la connexion stockage | ❌ |
| 4 | GET | `/api/v1/settings/storage/usage` | Utilisation du stockage | ❌ |
| 5 | GET | `/api/v1/settings/storage/buckets` | Lister les buckets | ❌ |

**Note :** L'interface `ObjectStorage` existe dans `interfaces/services.go` mais pas de config admin.

### 9.7 CDN — `/dash/settings/cdn`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/settings/cdn` | Config CDN | ❌ |
| 2 | PUT | `/api/v1/settings/cdn` | Modifier la config CDN | ❌ |
| 3 | POST | `/api/v1/settings/cdn/purge` | Purger le cache CDN | ❌ |
| 4 | GET | `/api/v1/settings/cdn/stats` | Statistiques CDN (bande passante, requêtes) | ❌ |

**Modèle existant :** `CdnAsset` (dans `media.go`) mais pas de config admin.

### 9.8 Domains — `/dash/settings/domains`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/settings/domains` | Liste des domaines personnalisés | ❌ |
| 2 | POST | `/api/v1/settings/domains` | Ajouter un domaine | ❌ |
| 3 | GET | `/api/v1/settings/domains/:domainId` | Détail d'un domaine | ❌ |
| 4 | DELETE | `/api/v1/settings/domains/:domainId` | Supprimer un domaine | ❌ |
| 5 | POST | `/api/v1/settings/domains/:domainId/verify` | Vérifier un domaine | ❌ |
| 6 | POST | `/api/v1/settings/domains/:domainId/ssl` | Générer le certificat SSL | ❌ |

### 9.9 APIs — `/dash/settings/apis`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/settings/apis` | Liste des clés API | ❌ |
| 2 | POST | `/api/v1/settings/apis` | Générer une clé API | ❌ |
| 3 | GET | `/api/v1/settings/apis/:keyId` | Détail d'une clé | ❌ |
| 4 | PATCH | `/api/v1/settings/apis/:keyId` | Modifier les permissions d'une clé | ❌ |
| 5 | DELETE | `/api/v1/settings/apis/:keyId` | Révoquer une clé | ❌ |
| 6 | GET | `/api/v1/settings/apis/:keyId/usage` | Utilisation d'une clé | ❌ |

**Note :** Le modèle `Application` dans `platform.go` pourrait servir. Le champ `SYSTEM_KEY` dans `.env` est une clé système hardcodée.

### 9.10 OAuth — `/dash/settings/oauth`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/settings/oauth` | Config des providers OAuth | ✅ (côté auth) |
| 2 | PUT | `/api/v1/settings/oauth/:provider` | Configurer un provider OAuth | ❌ |
| 3 | POST | `/api/v1/settings/oauth/:provider/test` | Tester la config OAuth | ❌ |
| 4 | GET | `/api/v1/settings/oauth/:provider/callback-url` | Obtenir l'URL de callback | ❌ |

### 9.11 Integrations — `/dash/settings/integrations`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/settings/integrations` | Liste des intégrations tierces | ❌ |
| 2 | POST | `/api/v1/settings/integrations` | Ajouter une intégration | ❌ |
| 3 | GET | `/api/v1/settings/integrations/:integrationId` | Détail d'une intégration | ❌ |
| 4 | PATCH | `/api/v1/settings/integrations/:integrationId` | Modifier une intégration | ❌ |
| 5 | DELETE | `/api/v1/settings/integrations/:integrationId` | Supprimer une intégration | ❌ |
| 6 | POST | `/api/v1/settings/integrations/:integrationId/test` | Tester la connexion | ❌ |
| 7 | GET | `/api/v1/settings/integrations/:integrationId/logs` | Logs d'une intégration | ❌ |

**Note :** L'endpoint webhook `POST /api/v1/webhooks/:provider/:integrationId` existe.

### 9.12 Maintenance — `/dash/settings/maintenance`

Page vide — à implémenter

| # | Méthode | Endpoint | Description | Statut |
|---|---------|----------|-------------|--------|
| 1 | GET | `/api/v1/settings/maintenance` | État de la maintenance (mode on/off, message) | ❌ |
| 2 | PUT | `/api/v1/settings/maintenance` | Activer/désactiver le mode maintenance | ❌ |
| 3 | POST | `/api/v1/settings/maintenance/cache-clear` | Vider tous les caches | ❌ |
| 4 | POST | `/api/v1/settings/maintenance/db-optimize` | Optimiser la base de données | ❌ |
| 5 | GET | `/api/v1/settings/maintenance/jobs` | Tâches de maintenance planifiées | ❌ |

---

## Résumé global

| Section | Endpoints totaux | ✅ Implémentés | 🔶 Partiels | ❌ À créer |
|---------|:---:|:---:|:---:|:---:|
| **Overview** | 5 | 0 | 0 | 5 |
| **Content** | 34 | 14 | 3 | 17 |
| **Community** | 40 | 9 | 5 | 26 |
| **Scheduling** | 27 | 6 | 3 | 18 |
| **Analytics** | 18 | 0 | 0 | 18 |
| **Media** | 37 | 24 | 0 | 13 |
| **Support** | 22 | 0 | 3 | 19 |
| **System** | 25 | 2 | 2 | 21 |
| **Settings** | 52 | 3 | 4 | 45 |
| **TOTAL** | **260** | **58** | **20** | **182** |

### Priorités recommandées

**Phase 1 — Fondations (Models + Repos manquants)**
- Créer les modèles : `Category`, `FAQ`, `Ticket`, `ModerationQueue`, `CalendarEvent`, `Premiere`, `Integration`, `ApiKey`, `EmailTemplate`, `MaintenanceSetting`
- Ajouter les repos/services/routes pour : Tags, Categories, Audit Logs (admin-scoped)

**Phase 2 — Pages existantes à wirer**
- Overview : remplacer les données mock par de vrais appels API
- Anime Manager : wirer le frontend aux endpoints existants
- Genres : construire la page (endpoints déjà prêts)
- Simulcasts : construire la page (endpoints déjà prêts)

**Phase 3 — Admin User Management**
- Endpoints admin pour users, sessions, rôles, permissions
- Pages Community : Users, Roles, Permissions

**Phase 4 — Analytics & Monitoring**
- Endpoints analytics agrégés (watch time, devices, geography, active users)
- System health détaillé, logs, queue, cache

**Phase 5 — Settings & Configuration**
- CRUD settings par catégorie (security, branding, email, SEO, storage, CDN, domains, APIs, OAuth, integrations, maintenance)

**Phase 6 — Support & Modération**
- Ticketing system complet
- FAQ management
- Modération unifiée
