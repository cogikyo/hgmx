import { animate, stagger, τmax } from "../motion.js";

export default function (el) {
  const screen = el.querySelector('path[d^="M2 8"]');
  animate(screen, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });

  const waves = [el.querySelector("line"), el.querySelector('path[d^="M2 16"]'), el.querySelector('path[d^="M2 12"]')];
  animate(waves, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09) });
}
