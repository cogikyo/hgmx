import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const shaft = el.querySelector('path[d^="M4 10"]');
  animate(shaft, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const head = el.querySelector('path[d^="m8 22"]');
  animate(head, { translateX: [4, -2, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
