# Watchman

It started as a concrete idea that came up in a specific context.

Not from a project I was deeply involved in at the code level, but from design discussions around a related problem. My contribution was conceptual rather than practical, focused on how to limit deviations in LLM-driven generation once constraints become relevant.

One observation was consistent: soft constraints tend to erode over time. Systems adapt around them, especially when the objective is task completion rather than alignment with intent.

Separately, I maintain a code generator that is fully deterministic and heavily template-based. Iterating on it often means updating templates, adjusting assumptions, and revalidating behavior as [patterns and guarantees evolve](https://github.com/hatmaxkit/hatmax-legacy/blob/main/docs/project-direction.md). Over time, this reduces the practical benefit of generation and increases maintenance effort.

Moving that generator toward an LLM-driven approach seems likely. At the same time, unconstrained generation produces results that do not match the patterns and constraints I want to enforce.

Watchman is an attempt to connect these concerns.

The intent is to allow generation and iteration, while enforcing fixed constraints at the execution level, in a mechanical and predictable way.

This repository reflects an early stage of that work. The shape of the solution is still evolving, and its usefulness beyond my own workflows is not yet clear.

