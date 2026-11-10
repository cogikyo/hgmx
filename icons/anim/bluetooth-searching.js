import { animate, stagger, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const rune = el.querySelector('path[d^="m7 7"]');
  animate(rune, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
  animate(rune, { scale: [1, 0.85, 1.15, 1] }, { duration: τmax });

  const signal = [el.querySelector('path[d^="M18 12"]'), el.querySelector('path[d^="M20.83"]')];
  animate(signal, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09, { startDelay: τmin }) });
}
