import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const speaker = el.querySelector('path[d^="M11 4.702"]');
  animate(speaker, { scale: [1, 0.85, 1] }, { duration: τmid });

  const bars = el.querySelectorAll('path[d^="m16.5"]');
  animate(bars, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: stagger(τmin) });
}
