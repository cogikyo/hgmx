import { animate, stagger } from "../motion.js";

export default function (el) {
  const waves = [el.querySelector('path[d^="M5 9"]'), el.querySelector('path[d^="M5 15"]')];
  animate(waves, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12) });
}
