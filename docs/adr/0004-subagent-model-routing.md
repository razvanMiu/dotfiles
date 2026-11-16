---
status: accepted
---

# Subagent model routing

The Pi parent session runs `ollama-cloud/deepseek-v4.1-flash` at thinking `high`, and `subagents.agentOverrides` pins one model per builtin role. The parent choice is driven by what a long tool session costs: the parent re-reads its whole context every turn, so cached-prefix price dominates, and `deepseek-v4.1-flash` combines a 1M window, a 384K output ceiling, vision for screenshots, and a $0.006 cache-read that no other capable model matches. Roles are pinned in `~/.config/pi/agent/settings.json` rather than written into skills, so the workflow documents never name a model and swapping providers is a one-place edit.

Roles split by what they actually spend. `scout` and `researcher` run `minimax-m3.1-flash`, because reconnaissance and web synthesis are input-heavy with compressed output. `worker` and `reviewer` run `deepseek-v4.1-flash` at `high`, since both need sustained reasoning over large inputs where the cached-input price compounds. `oracle` runs `glm-5.3-flash` at `high`: it keeps a model family independent of the parent, which is the point of a drift audit, while staying 9.3x cheaper than the `glm-5.3` flagship; it is also multimodal, and the oracle inherits a forked parent context that contains screenshots the text-only flagship would receive as omitted images. `evidence-auditor` runs `glm-5.3-flash` at `low`, and `delegate` runs `gpt-oss:20b` for bounded mechanical work.

The free in-house `eea-inhouse/Inhouse-LLM/qwen3.8-27b` (262K context, 8192 max output, no thinking, unreachable off the EEA network) is the observational-memory worker model, sized with `agentMaxTokens: 8192`. That work is frequent, tool-free summarization, which fits a free endpoint and a small output ceiling. The same model also serves reconnaissance when the operator is on that network, by passing it per run: `subagent({agent: "scout", model: "eea-inhouse/Inhouse-LLM/qwen3.8-27b"})`. It is not a default role because no fallback exists at any layer: `pi-subagents` resolves one model per launch and returns that model's provider errors instead of switching, while `pi-observational-memory` falls back to the session model only when the configured model is absent from the registry, never when a request fails. Off the EEA network the memory run therefore warns and backs off rather than degrading silently, and the session keeps working without new observations.

## Considered Options

- Run the orchestrator on `glm-5.3` or `kimi-k3`; rejected because the parent's cost scales with every turn's re-read input, where those models cost 4.7x to 10x more for a modest quality gain.
- Keep `glm-5.3` for the oracle; rejected because `glm-5.3-flash` costs 9.3x less on input, scores 84.3 versus 88.2 on Terminal-Bench 2.1, and is the only one of the two that can read images in the inherited context.
- Put `deepseek-v4.1-flash` in the oracle too; rejected because same-family review of the parent's decisions shares the parent's blind spots.
- Make the free in-house model a default subagent role; rejected because it is network-gated and no layer implements model fallback.
- Put the role pins in the orchestration skill; rejected because model choice is configuration, and skills stay portable.

## Consequences

Every role pin must name a vision model if that role can receive images: the oracle inherits forked context and the worker verifies screenshots. Off the EEA network the observational-memory worker cannot run, so observations accumulate only while connected, and `agentMaxTokens` must stay at or below 8192 or the run aborts. Cache pricing is a claim about the provider, not a guarantee; if cached-prefix hits stay low in a long session, the parent's cost advantage narrows and the crossover against `minimax-m3.1-flash` is roughly 94% cached tokens.
