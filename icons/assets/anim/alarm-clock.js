import { animate, stagger, τmax } from "../motion.js";

export default function (el) {
  const clock = el.querySelector("g");
  animate(clock, { rotate: [0, -12, 12, -8, 8, 0] }, { duration: τmax });
  const bells = [el.querySelector('path[d^="M5 3"]'), el.querySelector('path[d^="m22 6"]')];
  animate(bells, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09) });
}
