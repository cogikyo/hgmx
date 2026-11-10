import { animate, stagger } from "../motion.js";

export default function (el) {
  const pieces = [el.querySelector('path[d^="M4 16"]'), el.querySelector('path[d^="M16 12"]')];
  animate(pieces, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12) });
}
