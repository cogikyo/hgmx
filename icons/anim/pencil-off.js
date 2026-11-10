import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const pencil = el.querySelector("g");
  animate(pencil, { rotate: [0, -8, 4, 0] }, { duration: τmid });

  const slash = el.querySelector('path[d^="m2 2"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
