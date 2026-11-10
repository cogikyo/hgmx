import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M13 22"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const emblem = el.querySelector("g");
  animate(emblem, { scale: [1, 0.85, 1.1, 1] }, { duration: τmid, delay: τmin });
}
