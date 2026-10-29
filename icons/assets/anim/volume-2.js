import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const speaker = el.querySelector('path[d^="M11 4.702"]');
  animate(speaker, { scale: [1, 0.85, 1] }, { duration: τmid });

  const waves = el.querySelectorAll('path:not([d^="M11 4.702"])');
  animate(waves, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09) });
}
