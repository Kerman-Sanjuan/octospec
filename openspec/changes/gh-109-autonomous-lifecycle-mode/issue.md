# Issue #109: Modo autónomo: ejecutar el ciclo de vida completo sin supervisión humana

https://github.com/Kerman-Sanjuan/octospec/issues/109

## User story

Como mantenedor/desarrollador de octospec, quiero un modo autónomo donde todo el ciclo de vida (idea → spec → apply → ship → archive) se ejecute solo con agentes, para poder dejarlo trabajando durante la noche sin tener que invocar y aprobar cada etapa manualmente.

## Context / problem

Hoy cada etapa del ciclo de vida (`idea`, `spec`, `apply`, `ship`, `archive`) se ejecuta manualmente: un humano tiene que invocar cada comando y aprobar el paso a la siguiente etapa. No existe un modo donde el sistema recorra el ciclo completo sin supervisión, así que no se puede dejar trabajando durante la noche sin intervención.

## Requirements

- El modo autónomo SHALL recorrer las 5 etapas completas (idea → spec → apply → ship → archive) para un issue sin pedir confirmación humana en ningún paso intermedio.
- El modo autónomo SHALL operar sobre las issues en curso (`status:in-progress` u otros estados activos) del backlog de GitHub, identificando además bugs que surjan durante el trabajo, siempre respetando la definición de producto existente (specs en `openspec/specs/`).
- Si una etapa falla o un gate no pasa, el sistema SHALL reintentar hasta que la etapa complete correctamente, antes de continuar con la siguiente.
- Todo el rastro del trabajo (comentarios, commits, pull requests, actualizaciones de estado) SHALL quedar registrado en GitHub, de la misma forma en que se hace hoy manualmente.

## Success criteria

- Un desarrollador deja el sistema corriendo por la noche y, al día siguiente, encuentra uno o más issues movidos de `backlog`/`in-progress` hasta `done` (o al menos hasta pull request abierto), con todos los pasos (spec, apply, ship) documentados en GitHub.
- Si una etapa falla repetidamente, el issue queda en un estado claro (p. ej. comentario de error) en vez de quedar bloqueado en silencio.

## Out of scope

- No incluye un scheduler/cron que lance el modo automáticamente; se invoca manualmente para que corra toda la noche.
- No incluye cambiar las reglas de los gates existentes (G2a, G2b, G3, G4, G6).
- No incluye desplegar o mergear a `main` sin pull request; sigue siendo PR + revisión, aunque sea revisión asíncrona.
