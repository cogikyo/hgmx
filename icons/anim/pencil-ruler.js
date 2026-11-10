import { animate, stagger, τmax } from "../motion.js";

export default function (el) {
  const pencil = el.querySelector("g");
  animate(pencil, { rotate: [0, -12, 8, -4, 0] }, { duration: τmax });

  const ticks = el.querySelectorAll('path[d^="m8 6"], path[d^="m18 16"]');
  animate(ticks, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.21) });
}
