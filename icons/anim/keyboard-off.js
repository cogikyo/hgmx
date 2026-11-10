import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const board = el.querySelectorAll('path[d^="M 20 4"], path[d^="M 22 6"], path[d^="M 9.69"], path[d^="M20 20"], path[d^="M 7 16"]');
  animate(board, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const slash = el.querySelector('path[d="m2 2 20 20"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
