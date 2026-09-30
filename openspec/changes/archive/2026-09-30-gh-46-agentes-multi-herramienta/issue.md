# Issue #46: Modelo agentico multi-herramienta con modelos configurables por agente

https://github.com/Kerman-Sanjuan/octospec/issues/46

## User story

Como maintainer de un repo que usa octospec, quiero que cada etapa del flujo (idea, spec, apply, ship, archive) corra sobre un agente especializado y configurable, apoyado por skills del proceso, para producir mejores artefactos con menos fricción.

## Context / problem

El flujo depende de que el humano redacte bien cada artefacto y la calidad es irregular. Un agente monolítico carga todo el manual y diluye la atención. No permite elegir modelo por etapa, aunque unas tareas piden razonamiento y otras velocidad o coste. El estándar multi-herramienta no está definido, así que cada herramienta arrastra su propia copia. El alcance y el contexto de cada etapa son grandes, así que conviene partir el trabajo en agentes con capacidades específicas.

## Requirements

- octospec SHALL definir un agente por etapa del flujo (idea, spec, apply, ship, archive), cada uno con un prompt y una misión acotados.
- octospec MUST permitir elegir el modelo de cada agente desde un único punto de configuración.
- octospec SHALL generar las definiciones de agente para el estándar multi-herramienta (por ejemplo pi, Claude Code, OpenCode, Codex) desde una única fuente canónica.
- octospec MUST inyectar a cada agente solo el contexto y las skills de su etapa.
- octospec SHALL mantener el CLI como orquestador: valida entradas, invoca al agente y escribe el artefacto.
- octospec MUST definir el estándar multi-herramienta (formato, campos obligatorios y cómo se renderiza por herramienta) antes de implementar los agentes.
- octospec SHALL validar que el conjunto de agentes no rompe los gates existentes.
- octospec MUST exponer un flujo interactivo, con interfaz TUI, para definir el modelo de cada rol (pensante, implementador, etc.).

## Success criteria

- Existe una fuente canónica de agentes y `octospec update` la renderiza para cada herramienta soportada, sin copias por tool.
- Se puede cambiar el modelo de un agente en un solo archivo y el cambio se propaga a todas las herramientas.
- Al ejecutar cada etapa, el agente recibe solo su contexto y sus skills, y se puede comprobar en el prompt generado.
- El CLI ofrece un flujo interactivo, con interfaz TUI, para definir el modelo de cada rol sin editar archivos a mano.
- Los gates actuales (`cli` y `gates`) siguen en verde y `openspec validate --all --strict` pasa.
- La documentación (README y docs) describe cómo elegir modelo y qué agente cubre cada etapa.
- Añadir una herramienta nueva no obliga a tocar la definición de los agentes.
- Las instrucciones para desarrollar octospec (`AGENTS.md` y el flujo de comandos del repo) siguen las nuevas reglas.

## Out of scope

- Redefinir las reglas de los gates ni `scripts/check-gates.sh`.
- Proveedores de modelo locales o no soportados por el harness.
- Reescribir todas las skills existentes, solo las referencias que cada agente necesita.
- Migrar repos consumidores ya instalados, salvo que la actualización sea automática.
