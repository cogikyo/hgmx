import { animate, τmid } from "../motion.js";

export default function (el) {
  const speaker = el.querySelector('path[d^="M11 4.702"]');
  animate(speaker, { scale: [1, 0.85, 1] }, { duration: τmid });

  const wave = el.querySelector('path[d^="M16 9"]');
  animate(wave, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3 });
}
