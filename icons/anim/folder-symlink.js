import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const folder = el.querySelector('path[d^="M2 9.35"]');
  animate(folder, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const head = el.querySelector('path[d^="m8 16"]');
  animate(head, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
