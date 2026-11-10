import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const lines = [el.querySelector('path[d^="M12 6"]'), el.querySelector('path[d="M6 12v9"]')];
  animate(lines, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const x = el.querySelectorAll('path[d^="m8.5"]');
  animate(x, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
